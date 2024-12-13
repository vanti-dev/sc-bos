package devices

import (
	"context"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"iter"
	"maps"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protopath"
	"google.golang.org/protobuf/reflect/protorange"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

//go:generate protomod protoc -- -I . -I ../../../proto --go_out=paths=source_relative:. download.proto

func (s *Server) GetDownloadDevicesUrl(_ context.Context, request *gen.GetDownloadDevicesUrlRequest) (*gen.DownloadDevicesUrl, error) {
	// validate
	switch request.MediaType {
	case "":
		request.MediaType = "text/csv"
	case "text/csv": // supported
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unsupported media type %q", request.MediaType)
	}

	tokenData := &DownloadToken{
		Request: request,
	}
	tokenStr, err := encodeDownloadToken(tokenData)
	if err != nil {
		return nil, err
	}
	u := s.downloadUrlBase
	if err := s.encodeDownloadUrlToken(&u, tokenStr); err != nil {
		return nil, err
	}
	return &gen.DownloadDevicesUrl{
		Url: u.String(),
	}, nil
}

// DownloadDevicesHTTPHandler responds to HTTP request urls returned by GetDownloadDevicesUrl.
func (s *Server) DownloadDevicesHTTPHandler(w http.ResponseWriter, r *http.Request) {
	// 1. parse and validate the request
	tokenStr, err := s.decodeDownloadUrlToken(r)
	if err != nil {
		http.Error(w, "invalid download token", http.StatusBadRequest)
		return
	}
	token, err := decodeDownloadToken(tokenStr)
	if err != nil {
		http.Error(w, "corrupted download token", http.StatusBadRequest)
		return
	}

	// 2. collect header related information
	deviceList := s.node.ListAllMetadata(
		resource.WithInclude(func(id string, item proto.Message) bool {
			if item == nil {
				return false
			}
			device := &gen.Device{
				Name:     id,
				Metadata: item.(*traits.Metadata),
			}
			return deviceMatchesQuery(token.Request.Query, device)
		}),
	)
	usefulDevices := func() iter.Seq[*traits.Metadata] {
		return func(yield func(*traits.Metadata) bool) {
			for _, d := range deviceList {
				if len(d.Traits) == 0 {
					continue
				}
				if proto.Equal(d, &traits.Metadata{Name: d.Name, Traits: []*traits.TraitMetadata{{Name: string(trait.Metadata)}}}) {
					continue // skip boring devices
				}
				if !yield(d) {
					return
				}
			}
		}
	}

	headerSet := make(map[string]struct{})
	// collect headers for all populated metadata fields
	for d := range usefulDevices() {
		err := protorange.Range(d.ProtoReflect(), func(values protopath.Values) error {
			p := values.Path
			leafStep := p.Index(-1)
			switch leafStep.Kind() {
			case protopath.FieldAccessStep:
				fd := leafStep.FieldDescriptor()
				if fd == nil {
					return nil
				}
				if fd.Cardinality() == protoreflect.Repeated {
					return nil // no headers for lists
				}
				if fd.Kind() == protoreflect.MessageKind {
					return nil // no headers for messages
				}
			case protopath.ListIndexStep:
				fd := leafStep.FieldDescriptor()
				if fd == nil {
					return nil
				}
				if fd.Cardinality() == protoreflect.Repeated {
					return nil // no headers for lists
				}
				if fd.Kind() == protoreflect.MessageKind {
					return nil // no headers for messages
				}
			default:
			}

			header := protoPathToHeader(p)
			if header == "" {
				return nil
			}
			headerSet[header] = struct{}{}
			return nil
		})
		if err != nil {
			http.Error(w, "failed to collect headers", http.StatusInternalServerError)
			return
		}
	}
	delete(headerSet, "traits.name") // we process these separately

	// capture trait headers
	traitNameSet := make(map[string]struct{})
	for d := range usefulDevices() {
		for _, t := range d.Traits {
			traitNameSet[t.Name] = struct{}{}
		}
	}

	traitInfo := s.getTraitInfo()
	for traitName := range traitNameSet {
		info, ok := traitInfo[traitName]
		if !ok {
			continue
		}
		for _, header := range info.headers {
			headerSet[header] = struct{}{}
		}
	}

	headers := slices.Sorted(maps.Keys(headerSet))
	headerIndex := make(map[string]int)
	for i, h := range headers {
		headerIndex[h] = i
	}

	// 3. start collecting the data and streaming it to the client
	w.Header().Set("Content-Disposition", "attachment; filename=devices.csv")
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")

	csvOut := csv.NewWriter(w)
	if err := csvOut.Write(headers); err != nil {
		return
	}

	for d := range usefulDevices() {
		row := make([]string, len(headers))
		captureMDValues(d, headerIndex, row)
		for _, t := range d.Traits {
			info, ok := traitInfo[t.Name]
			if !ok {
				continue
			}
			values, err := info.get(r.Context(), d.Name)
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to get trait info: %q %q %v", d.Name, t.Name, err), http.StatusInternalServerError)
				continue
			}
			for h, v := range values {
				index, ok := headerIndex[h]
				if !ok {
					continue
				}
				row[index] = v
			}
		}
		if err := csvOut.Write(row); err != nil {
			return
		}
	}
	csvOut.Flush()
}

func captureMDValues(md *traits.Metadata, headerIndex map[string]int, row []string) {
	_ = protorange.Range(md.ProtoReflect(), func(values protopath.Values) error {
		p := values.Path
		leafStep := p.Index(-1)
		switch leafStep.Kind() {
		case protopath.FieldAccessStep:
			fd := leafStep.FieldDescriptor()
			if fd == nil {
				return nil
			}
			if fd.Cardinality() == protoreflect.Repeated {
				return nil // no headers for lists
			}
			if fd.Kind() == protoreflect.MessageKind {
				return nil // no headers for messages
			}
		case protopath.ListIndexStep:
			fd := leafStep.FieldDescriptor()
			if fd == nil {
				return nil
			}
			if fd.Cardinality() == protoreflect.Repeated {
				return nil // no headers for lists
			}
			if fd.Kind() == protoreflect.MessageKind {
				return nil // no headers for messages
			}
		default:
		}

		header := protoPathToHeader(p)
		if header == "" {
			return nil
		}
		index, ok := headerIndex[header]
		if !ok {
			return nil // skip
		}
		row[index] = values.Values[len(values.Values)-1].String()
		return nil
	})
}

type traitInfo struct {
	headers []string
	get     func(ctx context.Context, name string) (map[string]string, error)
}

func protoPathToHeader(p protopath.Path) string {
	var parts []string
	for _, step := range p {
		switch step.Kind() {
		case protopath.FieldAccessStep:
			parts = append(parts, string(step.FieldDescriptor().Name()))
		case protopath.MapIndexStep:
			parts = append(parts, step.MapIndex().String())
		case protopath.ListIndexStep:
			// skip writing the index, {bar: [{foo}]} -> bar.foo instead of bar[0].foo
		case protopath.RootStep:
			// skip writing the root
		case protopath.UnknownAccessStep:
		case protopath.AnyExpandStep:
		}
	}
	return strings.Join(parts, ".")
}

func (s *Server) RegisterHTTPMux(mux *http.ServeMux) {
	mux.HandleFunc(s.downloadUrlBase.Path, s.DownloadDevicesHTTPHandler)
}

type DownloadUrlEncoder func(dst *url.URL, token string) error
type DownloadUrlDecoder func(*http.Request) (string, error)

func WithDownloadUrlCodec(urlEncoder DownloadUrlEncoder, urlDecoder DownloadUrlDecoder) Option {
	return func(s *Server) {
		s.downloadUrlEncoder = urlEncoder
		s.downloadUrlDecoder = urlDecoder
	}
}

func WithDownloadUrlBase(base url.URL) Option {
	return func(s *Server) {
		s.downloadUrlBase = base
	}
}

func DecodeDownloadUrlToken(r *http.Request) (string, error) {
	return r.URL.Query().Get("ddt"), nil
}

func (s *Server) decodeDownloadUrlToken(r *http.Request) (string, error) {
	if s.downloadUrlDecoder == nil {
		return DecodeDownloadUrlToken(r)
	}
	return s.downloadUrlDecoder(r)
}

func EncodeDownloadUrlToken(dst *url.URL, token string) error {
	q := dst.Query()
	q.Set("ddt", token)
	dst.RawQuery = q.Encode()
	return nil
}

func (s *Server) encodeDownloadUrlToken(dst *url.URL, token string) error {
	if s.downloadUrlEncoder == nil {
		return EncodeDownloadUrlToken(dst, token)
	}
	return s.downloadUrlEncoder(dst, token)
}

func encodeDownloadToken(token *DownloadToken) (string, error) {
	data, err := proto.Marshal(token)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data), nil
}

func decodeDownloadToken(token string) (*DownloadToken, error) {
	data, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}
	var dt DownloadToken
	if err := proto.Unmarshal(data, &dt); err != nil {
		return nil, err
	}
	return &dt, nil
}
