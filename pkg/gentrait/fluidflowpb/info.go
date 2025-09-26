package fluidflowpb

import (
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

type InfoServer struct {
	gen.UnimplementedFluidFlowInfoServer
	FluidFlowSupport *gen.FluidFlowSupport
}

func (i *InfoServer) DescribeFluidFlow(_ *gen.DescribeFluidFlowRequest) (*gen.FluidFlowSupport, error) {
	return i.FluidFlowSupport, nil
}
