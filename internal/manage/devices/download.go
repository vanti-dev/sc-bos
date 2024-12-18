package devices

import (
	"context"
	"encoding/base64"
	"encoding/csv"
	"errors"
	"fmt"
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
	tokenStr, err := downloadTokenToString(tokenData)
	if err != nil {
		return nil, err
	}
	u := s.downloadUrlBase
	if err := s.writeDownloadToken(&u, tokenStr); err != nil {
		return nil, err
	}
	return &gen.DownloadDevicesUrl{
		Url: u.String(),
	}, nil
}

// DownloadDevicesHTTPHandler responds to HTTP request urls returned by GetDownloadDevicesUrl.
func (s *Server) DownloadDevicesHTTPHandler(w http.ResponseWriter, r *http.Request) {
	// 1. parse and validate the request
	tokenStr, err := s.readDownloadToken(r)
	if err != nil {
		http.Error(w, "invalid download token", http.StatusBadRequest)
		return
	}
	token, err := downloadTokenFromString(tokenStr)
	if err != nil {
		http.Error(w, "corrupted download token", http.StatusBadRequest)
		return
	}

	// 2. work out which devices to return and collect data headers
	traitInfo := s.getTraitInfo()
	deviceList, headers, err := s.listDevicesAndHeaders(token, traitInfo)
	if err != nil {
		var httpErr httpError
		if !errors.As(err, &httpErr) {
			http.Error(w, httpErr.msg, httpErr.code)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	headerIndex := make(map[string]int)
	for i, h := range headers {
		headerIndex[h] = i
	}

	// 3. start collecting the data and streaming it to the client
	// note: we only set the headers here (rather than earlier) to allow bad status codes to be returned if needed
	w.Header().Set("Content-Disposition", "attachment; filename=devices.csv")
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")

	csvOut := csv.NewWriter(w)
	if err := csvOut.Write(headers); err != nil {
		return
	}

	for _, d := range deviceList {
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

func (s *Server) listDevicesAndHeaders(token *DownloadToken, traitInfo map[string]traitInfo) (devices []*traits.Metadata, headers []string, err error) {
	devices = s.node.ListAllMetadata(
		resource.WithInclude(func(id string, item proto.Message) bool {
			if item == nil {
				return false
			}
			md := item.(*traits.Metadata)
			device := &gen.Device{
				Name:     id,
				Metadata: md,
			}
			if len(md.Traits) == 0 {
				return false
			}
			// Skip boring devices, aka those that have no metadata or other trait data.
			// They'd just show up as name=md.Name and a bunch of empty columns anyway.
			if proto.Equal(md, &traits.Metadata{Name: md.Name, Traits: []*traits.TraitMetadata{{Name: string(trait.Metadata)}}}) {
				return false
			}
			return deviceMatchesQuery(token.Request.Query, device)
		}),
	)

	headerSet := make(map[string]struct{})
	if err := collectMetadataHeaders(headerSet, devices); err != nil {
		return nil, nil, err
	}

	collectTraitHeaders(headerSet, devices, traitInfo)

	headers = slices.Sorted(maps.Keys(headerSet))
	return devices, headers, nil
}

func collectMetadataHeaders(dst map[string]struct{}, deviceList []*traits.Metadata) error {
	// collect headers for all populated metadata fields
	for _, d := range deviceList {
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
			dst[header] = struct{}{}
			return nil
		})
		if err != nil {
			return httpError{http.StatusInternalServerError, "failed to collect headers"}
		}
	}
	delete(dst, "traits.name") // we process these separately
	return nil
}

func collectTraitHeaders(dst map[string]struct{}, deviceList []*traits.Metadata, traitInfo map[string]traitInfo) {
	// capture trait headers
	traitNameSet := make(map[string]struct{})
	for _, d := range deviceList {
		for _, t := range d.Traits {
			traitNameSet[t.Name] = struct{}{}
		}
	}

	for traitName := range traitNameSet {
		info, ok := traitInfo[traitName]
		if !ok {
			continue
		}
		for _, header := range info.headers {
			dst[header] = struct{}{}
		}
	}
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
	if len(parts) == 0 {
		return ""
	}
	return "md." + strings.Join(parts, ".")
}

func (s *Server) RegisterHTTPMux(mux *http.ServeMux) {
	mux.HandleFunc(s.downloadUrlBase.Path, s.DownloadDevicesHTTPHandler)
}

type DownloadTokenWriter func(dst *url.URL, token string) error
type DownloadTokenReader func(*http.Request) (string, error)

func WithDownloadTokenCodec(w DownloadTokenWriter, r DownloadTokenReader) Option {
	return func(s *Server) {
		s.downloadTokenWriter = w
		s.downloadTokenReader = r
	}
}

func WithDownloadUrlBase(base url.URL) Option {
	return func(s *Server) {
		s.downloadUrlBase = base
	}
}

func ReadDownloadToken(r *http.Request) (string, error) {
	return r.URL.Query().Get("ddt"), nil
}

func (s *Server) readDownloadToken(r *http.Request) (string, error) {
	if s.downloadTokenReader == nil {
		return ReadDownloadToken(r)
	}
	return s.downloadTokenReader(r)
}

func WriteDownloadToken(dst *url.URL, token string) error {
	q := dst.Query()
	q.Set("ddt", token)
	dst.RawQuery = q.Encode()
	return nil
}

func (s *Server) writeDownloadToken(dst *url.URL, token string) error {
	if s.downloadTokenWriter == nil {
		return WriteDownloadToken(dst, token)
	}
	return s.downloadTokenWriter(dst, token)
}

func downloadTokenToString(token *DownloadToken) (string, error) {
	data, err := proto.Marshal(token)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data), nil
}

func downloadTokenFromString(token string) (*DownloadToken, error) {
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

type httpError struct {
	code int
	msg  string
}

func (h httpError) Error() string {
	return http.StatusText(h.code) + ": " + h.msg
}
