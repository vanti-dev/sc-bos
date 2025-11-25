package lighttest

import (
	"context"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/smart-core-os/sc-bos/pkg/gen"
)

var errNotEnabled = status.Error(codes.FailedPrecondition, "lighttest - not enabled")

type Holder struct {
	gen.UnimplementedLightingTestApiServer

	client gen.LightingTestApiClient
	mu     sync.Mutex
}

func (h *Holder) Register(server *grpc.Server) {
	gen.RegisterLightingTestApiServer(server, h)
}

func (h *Holder) Fill(client gen.LightingTestApiClient) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.client = client
}

func (h *Holder) Empty() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.client = nil
}

func (h *Holder) GetLightHealth(ctx context.Context, request *gen.GetLightHealthRequest) (*gen.LightHealth, error) {
	c, err := h.getClient()
	if err != nil {
		return nil, err
	}
	return c.GetLightHealth(ctx, request)
}

func (h *Holder) ListLightHealth(ctx context.Context, request *gen.ListLightHealthRequest) (*gen.ListLightHealthResponse, error) {
	c, err := h.getClient()
	if err != nil {
		return nil, err
	}
	return c.ListLightHealth(ctx, request)
}

func (h *Holder) ListLightEvents(ctx context.Context, request *gen.ListLightEventsRequest) (*gen.ListLightEventsResponse, error) {
	c, err := h.getClient()
	if err != nil {
		return nil, err
	}
	return c.ListLightEvents(ctx, request)
}

func (h *Holder) GetReportCSV(ctx context.Context, request *gen.GetReportCSVRequest) (*gen.ReportCSV, error) {
	c, err := h.getClient()
	if err != nil {
		return nil, err
	}
	return c.GetReportCSV(ctx, request)
}

func (h *Holder) getClient() (gen.LightingTestApiClient, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.client == nil {
		return nil, errNotEnabled
	}
	return h.client, nil
}
