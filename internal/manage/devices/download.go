package devices

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/csv"
	"errors"
	"fmt"
	"iter"
	"maps"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protopath"
	"google.golang.org/protobuf/reflect/protorange"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-api/go/traits"
	timepb "github.com/smart-core-os/sc-api/go/types/time"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

//go:generate protomod protoc -- -I . -I ../../../proto --go_out=paths=source_relative:. download.proto

type tokenClaims struct {
	Body string `json:"b"`
}

func (s *Server) GetDownloadDevicesUrl(_ context.Context, request *gen.GetDownloadDevicesUrlRequest) (*gen.DownloadDevicesUrl, error) {
	// validate
	switch request.MediaType {
	case "":
		request.MediaType = "text/csv"
	case "text/csv": // supported
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unsupported media type %q", request.MediaType)
	}

	expireAfter := s.now().Add(s.downloadExpiry)
	tokenStr, err := s.signAndSerializeDownloadToken(&DownloadToken{Request: request}, expireAfter)
	if err != nil {
		return nil, err
	}

	u := s.downloadUrlBase
	if err := s.writeDownloadToken(&u, tokenStr); err != nil {
		return nil, err
	}
	return &gen.DownloadDevicesUrl{
		Url:             u.String(),
		MediaType:       request.MediaType,
		ExpireAfterTime: timestamppb.New(expireAfter),
	}, nil
}

// DownloadDevicesHTTPHandler responds to HTTP request urls returned by GetDownloadDevicesUrl.
// Requests must include a valid download token.
//
// CSV responses will include a header as the first row, for which columns are sorted alphabetically and grouped by md.name then md.* then *.
// For metadata columns each device is inspected to find all non-empty fields, each of which is included as a column in the md.* group.
// For trait values the supported traits and included columns is defined by the traitInfo map, see [Server.getTraitInfo] for details.
// Trait columns are fixed based on the advertised traits a device supports via its metadata.
// Typical column names are dot separated property paths, e.g. md.location.floor, access.grant, or meter.usage.
//
// Devices that have no traits or metadata (excluding their name and that they implement the metadata trait) are excluded from the response.
func (s *Server) DownloadDevicesHTTPHandler(w http.ResponseWriter, r *http.Request) {
	// 1. parse and validate the request
	tokenStr, err := s.readDownloadToken(r)
	if err != nil {
		http.Error(w, "invalid download token", http.StatusUnauthorized)
		return
	}

	writeErr := func(err error) {
		var httpErr httpError
		if errors.As(err, &httpErr) {
			http.Error(w, httpErr.msg, httpErr.code)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	token, err := s.parseAndValidateDownloadToken(tokenStr)
	if err != nil {
		writeErr(err)
		return
	}

	// 2. work out which devices to return and collect data headers
	traitInfo := s.getTraitInfo()
	deviceList, headers, err := s.listDevicesAndHeaders(token, traitInfo)
	if err != nil {
		writeErr(err)
		return
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
	out := newCSVWriter(csvOut, headerIndex)

	if token.Request.History == nil {
		s.writeLiveData(r.Context(), out, deviceList, traitInfo)
	} else {
		s.writeHistoricalData(r.Context(), out, deviceList, traitInfo, token.Request.History)
	}
	csvOut.Flush()
}

type writer interface {
	Write(row map[string]string)
}

func newCSVWriter(out *csv.Writer, headerIndex map[string]int) *csvWriter {
	return &csvWriter{out, headerIndex, make([]string, len(headerIndex))}
}

type csvWriter struct {
	out         *csv.Writer
	headerIndex map[string]int
	rowBuf      []string
}

func (c *csvWriter) Write(row map[string]string) {
	for h, v := range row {
		index, ok := c.headerIndex[h]
		if !ok {
			continue
		}
		c.rowBuf[index] = v
	}
	_ = c.out.Write(c.rowBuf)
	clear(c.rowBuf)
}

func (s *Server) writeLiveData(ctx context.Context, out writer, devices []*traits.Metadata, traitInfo map[string]traitInfo) {
	for _, d := range devices {
		row := make(map[string]string)
		captureMDValues(d, row)
		for _, t := range d.Traits {
			info, ok := traitInfo[t.Name]
			if !ok {
				continue
			}
			pageCtx, cleanup := context.WithTimeout(ctx, s.downloadPageTimeout)
			values, err := info.get(pageCtx, d.Name)
			cleanup()
			if err != nil {
				row[info.headers[0]] = fmt.Sprintf("ERR: %v", err)
				continue
			}
			maps.Copy(row, values)
		}
		out.Write(row)
	}
}

type source struct {
	device *traits.Metadata
	cursor *historyCursor
	skip   bool
}

func (s *Server) writeHistoricalData(ctx context.Context, out writer, devices []*traits.Metadata, traitInfo map[string]traitInfo, period *timepb.Period) {
	// we might want to allow the user to specify the order in the future
	const (
		// pageSize = 100
		// order    = "source"
		// pageSize = 10 // time order reads from all sources at once, limit memory use
		// order    = "time"
		pageSize = 100 // for device order we only keep the traits for a single device in memory at a time
		order    = "device"
	)

	var sources []*source
	for _, d := range devices {
		for _, t := range d.Traits {
			info, ok := traitInfo[t.Name]
			if !ok || info.history == nil {
				continue
			}
			cursor := info.history(d.Name, period, pageSize)
			if cursor == nil {
				continue
			}
			sources = append(sources, &source{device: d, cursor: cursor})
		}
	}

	switch order {
	case "source":
		s.writeHistoryDataBySource(ctx, out, sources)
	case "device":
		s.writeHistoryDataByDevice(ctx, out, sources)
	case "time":
		s.writeHistoryDataByTime(ctx, out, sources)
	}
}

// writeHistoryDataBySource writes history data ordered by source, then time.
// This means each devices trait records are written together.
func (s *Server) writeHistoryDataBySource(ctx context.Context, out writer, sources []*source) {
	for _, source := range sources {
		mdVals := make(map[string]string)
		captureMDValues(source.device, mdVals)
		for {
			pageCtx, cleanup := context.WithTimeout(ctx, s.downloadPageTimeout)
			head, err := source.cursor.Head(pageCtx)
			cleanup()
			if err != nil {
				break
			}
			head.use()
			vals := head.vals
			maps.Copy(vals, mdVals)
			vals["timestamp"] = head.at.Format(time.DateTime)
			out.Write(vals)
		}
	}
}

// writeHistoryDataByDevice writes history data ordered by device, then time.
func (s *Server) writeHistoryDataByDevice(ctx context.Context, out writer, sources []*source) {
	sourcesByDevice := make(map[string][]*source)
	for _, source := range sources {
		sourcesByDevice[source.device.Name] = append(sourcesByDevice[source.device.Name], source)
	}

	for _, sources := range sourcesByDevice {
		s.writeHistoryDataByTime(ctx, out, sources)
	}
}

// writeHistoryDataByTime writes history data ordered by time.
// This means device trait records are interleaved in time order.
func (s *Server) writeHistoryDataByTime(ctx context.Context, out writer, sources []*source) {
	// cache of metadata values we can reuse during the main loop
	mds := make(map[string]map[string]string, len(sources))
	for _, source := range sources {
		if _, ok := mds[source.device.Name]; ok {
			continue
		}
		mdVals := make(map[string]string)
		captureMDValues(source.device, mdVals)
		mds[source.device.Name] = mdVals
	}

	for len(sources) > 0 {
		var (
			oldestRecord *historyRecord
			oldestSource *source
			anySkipped   bool
		)
		for _, source := range sources {
			pageCtx, cleanup := context.WithTimeout(ctx, s.downloadPageTimeout)
			head, err := source.cursor.Head(pageCtx)
			cleanup()
			if err != nil {
				anySkipped = true
				source.skip = true
				continue
			}
			switch {
			case oldestRecord == nil:
				oldestRecord, oldestSource = &head, source
			case head.at.Before(oldestRecord.at):
				oldestRecord, oldestSource = &head, source
			}
		}

		// both checks aren't strictly necessary as both are set at the same time,
		// however the static checker doesn't know that
		if oldestRecord == nil || oldestSource == nil {
			return // no records were processed
		}

		oldestRecord.use()
		vals := mds[oldestSource.device.Name]
		maps.Copy(vals, oldestRecord.vals)
		vals["timestamp"] = oldestRecord.at.Format(time.DateTime)
		out.Write(vals)

		if anySkipped {
			sources = slices.DeleteFunc(sources, func(source *source) bool {
				return source.skip
			})
		}
	}
}

func (s *Server) signAndSerializeDownloadToken(tokenBody *DownloadToken, expireAfter time.Time) (string, error) {
	tokenBodyStr, err := downloadTokenToString(tokenBody)
	if err != nil {
		return "", status.Errorf(codes.Unavailable, "token body creation error: %v", err)
	}

	key, err := s.downloadKey()
	if err != nil {
		return "", status.Errorf(codes.Unavailable, "token key creation error: %v", err)
	}

	// using JWT/JOSE here for the signing/key gen means we can also use it later for validation
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: key}, (&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		return "", status.Errorf(codes.Unavailable, "token signer creation error: %v", err)
	}

	jwtClaims := &jwt.Claims{Expiry: jwt.NewNumericDate(expireAfter)}
	tokenClaims := tokenClaims{Body: tokenBodyStr}
	tokenStr, err := jwt.Signed(signer).Claims(tokenClaims).Claims(jwtClaims).Serialize()
	if err != nil {
		return "", status.Errorf(codes.Unavailable, "token serialization error: %v", err)
	}
	return tokenStr, nil
}

func (s *Server) parseAndValidateDownloadToken(tokenStr string) (*DownloadToken, error) {
	key, err := s.downloadKey()
	if err != nil {
		return nil, httpError{code: http.StatusInternalServerError, msg: "failed to read signing key"}
	}

	jwtToken, err := jwt.ParseSigned(tokenStr, []jose.SignatureAlgorithm{jose.HS256})
	if err != nil {
		return nil, httpError{code: http.StatusUnauthorized, msg: "invalid token"}
	}
	jwtClaims := &jwt.Claims{}
	var tokenClaims tokenClaims
	if err := jwtToken.Claims(key, jwtClaims, &tokenClaims); err != nil {
		return nil, httpError{code: http.StatusUnauthorized, msg: "untrusted token"}
	}
	if err := jwtClaims.ValidateWithLeeway(jwt.Expected{Time: s.now()}, s.downloadExpiryLeeway); err != nil {
		return nil, httpError{code: http.StatusUnauthorized, msg: "token expired"}
	}

	token, err := downloadTokenFromString(tokenClaims.Body)
	if err != nil {
		return nil, httpError{code: http.StatusUnauthorized, msg: "corrupted download token"}
	}
	return token, nil
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

	if token.Request.History != nil {
		// delete headers for traits that don't support history
		for _, info := range traitInfo {
			if info.history == nil {
				for _, header := range info.headers {
					delete(headerSet, header)
				}
			}
		}
	}

	headers = sortHeaders(maps.Keys(headerSet))

	if token.Request.History != nil {
		// timestamp should be the first column
		headers = append([]string{"timestamp"}, headers...)
	}

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
	delete(dst, "md.traits.name") // we process these separately
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

func captureMDValues(md *traits.Metadata, row map[string]string) {
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
		row[header] = values.Values[len(values.Values)-1].String()
		return nil
	})
}

func sortHeaders(headers iter.Seq[string]) []string {
	return slices.SortedFunc(headers, func(a string, b string) int {
		// sort metadata fields first
		aIsMD := strings.HasPrefix(a, "md.")
		bIsMD := strings.HasPrefix(b, "md.")
		switch {
		case aIsMD && !bIsMD:
			return -1
		case !aIsMD && bIsMD:
			return 1
		case aIsMD && bIsMD:
			// make sure md.name is first
			switch {
			case a == "md.name" && b == "md.name":
				return 0
			case a == "md.name":
				return -1
			case b == "md.name":
				return 1
			default:
				return strings.Compare(a, b)
			}
		default:
			return strings.Compare(a, b)
		}
	})
}

type traitInfo struct {
	headers []string
	get     func(ctx context.Context, name string) (map[string]string, error)
	history func(name string, period *timepb.Period, pageSize int32) *historyCursor
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
	mux.Handle(s.downloadUrlBase.Path, gziphandler.GzipHandler(http.HandlerFunc(s.DownloadDevicesHTTPHandler)))
}

func newHMACKeyGen(size int) func() ([]byte, error) {
	// todo: support key rotation that doesn't invalidate unexpired tokens,
	//  I expect the sig will need to change to func(string)(string, []byte, error)
	key := make([]byte, size)
	var err error
	_, err = rand.Read(key)
	return func() ([]byte, error) {
		return key, err
	}
}

// DownloadTokenWriter is a function that writes a download token to a URL.
type DownloadTokenWriter = func(dst *url.URL, token string) error

// DownloadTokenReader is a function that reads a download token from a URL.
type DownloadTokenReader = func(*http.Request) (string, error)

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

// readDownloadToken is the default implementation of DownloadTokenReader.
func readDownloadToken(r *http.Request) (string, error) {
	return r.URL.Query().Get("ddt"), nil
}

func (s *Server) readDownloadToken(r *http.Request) (string, error) {
	if s.downloadTokenReader == nil {
		return readDownloadToken(r)
	}
	return s.downloadTokenReader(r)
}

// writeDownloadToken is the default implementation of DownloadTokenWriter.
func writeDownloadToken(dst *url.URL, token string) error {
	q := dst.Query()
	q.Set("ddt", token)
	dst.RawQuery = q.Encode()
	return nil
}

func (s *Server) writeDownloadToken(dst *url.URL, token string) error {
	if s.downloadTokenWriter == nil {
		return writeDownloadToken(dst, token)
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
