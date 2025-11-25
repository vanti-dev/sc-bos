package pressurepb

import (
	"github.com/smart-core-os/sc-bos/pkg/gen"
)

type InfoServer struct {
	gen.UnimplementedPressureInfoServer
	PressureSupport *gen.PressureSupport
}

func (i *InfoServer) DescribePressure(_ *gen.DescribePressureRequest) (*gen.PressureSupport, error) {
	return i.PressureSupport, nil
}
