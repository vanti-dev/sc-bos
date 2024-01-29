package temperaturepb

import (
	"context"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

// PriorityServer adds priority support to a TemperatureApiServer.
// See [gen.Priority] for more information about how priority is applied to reads and writes.
type PriorityServer struct {
	gen.TemperatureApiServer
	// todo: replace this array with something more memory efficient
	priorityArray [256]*gen.UpdateTemperatureRequest
}

func NewPriorityServer(s gen.TemperatureApiServer) *PriorityServer {
	return &PriorityServer{TemperatureApiServer: s}
}

func (s *PriorityServer) GetTemperature(ctx context.Context, request *gen.GetTemperatureRequest) (*gen.Temperature, error) {
	if request.Priority == 0 {
		return s.TemperatureApiServer.GetTemperature(ctx, request)
	}
	pl, err := getPriority(request)
	if err != nil {
		return nil, err
	}
	if err := setTopPriorityHeader(ctx, s.priorityArray[:]); err != nil {
		return nil, err
	}
	d := s.priorityArray[pl]
	if d == nil {
		return nil, nil
	}
	return d.Temperature, nil
}

func (s *PriorityServer) UpdateTemperature(ctx context.Context, request *gen.UpdateTemperatureRequest) (res *gen.Temperature, err error) {
	pl, err := getPriority(request)
	if err != nil {
		return nil, err
	}
	old := s.priorityArray[pl]
	defer func() {
		if err != nil {
			s.priorityArray[pl] = old
		}
	}()

	if request.GetPriority() > 0 && request.GetTemperature() == nil {
		// todo: lock
		s.priorityArray[pl] = nil // pl is ok here, shouldClear returns false if we are using defaults implicitly.
		return s.applyChange(ctx, pl)
	}

	// We need to read the current value if
	// a. the write is a delta or
	// b. the write uses a mask
	if incompleteWrite(request) {
		err := s.hydrateWrite(ctx, request)
		if err != nil {
			return nil, err
		}
	}

	// todo: lock
	s.priorityArray[pl] = request
	return s.applyChange(ctx, pl)
}

// hydrateWrite converts an incomplete write into a complete write.
// A complete write is independent of the current value of the resource.
func (s *PriorityServer) hydrateWrite(ctx context.Context, request *gen.UpdateTemperatureRequest) error {
	rd, err := s.TemperatureApiServer.GetTemperature(ctx, &gen.GetTemperatureRequest{Name: request.Name})
	if err != nil {
		return err // todo: better error messages and codes
	}
	// todo: this doesn't feel performant, try and think of a different way to do it
	rv := resource.NewValue(resource.WithInitialValue(rd))
	nv, err := rv.Set(request.GetTemperature(), resource.WithUpdateMask(request.GetUpdateMask()), resource.InterceptBefore(func(old, new proto.Message) {
		ov, nv := old.(*gen.Temperature), new.(*gen.Temperature)
		if request.Delta {
			switch {
			case nv.SetPoint == nil:
				nv.SetPoint = ov.SetPoint
			case ov.SetPoint == nil:
				// do nothing, new value will replace old value anyway
			default:
				nv.SetPoint.ValueCelsius += ov.SetPoint.ValueCelsius
			}
		}
	}))
	if err != nil {
		return err
	}
	request.Temperature = nv.(*gen.Temperature)
	request.Delta = false
	request.UpdateMask = nil // todo: work out if this is safe to do
	return nil
}

func (s *PriorityServer) applyChange(ctx context.Context, p gen.Priority_Level) (*gen.Temperature, error) {
	// todo: lock
	v, top, ok := topPriority(s.priorityArray[:])
	if !ok {
		// Nothing to write, nothing to do, no values in the p-array.
		// Get the value from the peer to make the response accurate.
		return s.GetTemperature(ctx, &gen.GetTemperatureRequest{})
	}
	if err := setPriorityLevelHeader(ctx, top); err != nil {
		return nil, err
	}
	if top < p {
		return v.GetTemperature(), nil // the updated priority is lower than the top priority, nothing to do
	}
	// need to write the new value
	return s.TemperatureApiServer.UpdateTemperature(ctx, v)
}

func topPriority[T any](a []*T) (*T, gen.Priority_Level, bool) {
	for i, v := range a {
		if v != nil {
			return v, gen.Priority_Level(i), true
		}
	}
	return nil, 0, false
}

func setTopPriorityHeader[T any](ctx context.Context, a []*T) error {
	_, l, ok := topPriority(a)
	if !ok {
		return nil
	}
	return setPriorityLevelHeader(ctx, l)
}

func setPriorityLevelHeader(ctx context.Context, p gen.Priority_Level) error {
	md := metadata.MD{}
	md.Append("PriorityLevel", strconv.Itoa(int(p)))
	return grpc.SetHeader(ctx, md)
}

func getPriority(p interface{ GetPriority() gen.Priority_Level }) (gen.Priority_Level, error) {
	pl := p.GetPriority()
	if pl == 0 {
		return gen.Priority_DEFAULT, nil
	}
	if pl < 0 || pl > 255 {
		return pl, status.Errorf(codes.InvalidArgument, "invalid priority: %d, want [0,255]", p)
	}
	return pl, nil
}

func shouldClear(req interface{ GetPriority() gen.Priority_Level }, res proto.Message) bool {
	p := req.GetPriority()
	return p > 0 && res == nil
}

func incompleteWrite(r any) bool {
	if i, ok := r.(interface{ GetDelta() bool }); ok {
		if i.GetDelta() {
			return true
		}
	}
	if i, ok := r.(interface{ GetUpdateMask() *fieldmaskpb.FieldMask }); ok {
		if len(i.GetUpdateMask().GetPaths()) > 0 {
			return true
		}
	}

	return false
}
