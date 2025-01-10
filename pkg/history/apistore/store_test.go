package apistore

import (
	"context"
	"testing"

	"github.com/smart-core-os/sc-golang/pkg/wrap"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

func TestSlice_Len_panicOnFail(t *testing.T) {
	client := gen.NewHistoryAdminApiClient(wrap.ServerToClient(gen.HistoryAdminApi_ServiceDesc, &testFailingHistoryAdminServer{}))
	// the below can panic if we aren't handling sync.Once retries correctly.
	_, _ = New(client, "name", "source").Len(context.Background())
}

type testFailingHistoryAdminServer struct {
	gen.UnimplementedHistoryAdminApiServer
}
