package batch

import (
	"google.golang.org/grpc"

	"github.com/vanti-dev/sc-bos/pkg/gen"
)

type Server struct {
	gen.UnimplementedBatchDataApiServer

	source grpc.ClientConnInterface
}
