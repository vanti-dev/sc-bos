package meter

import (
	"context"

	"github.com/smart-core-os/sc-bos/pkg/gen"
)

type InfoServer struct {
	gen.UnimplementedMeterInfoServer
	MeterReading *gen.MeterReadingSupport
}

func (i *InfoServer) DescribeMeterReading(_ context.Context, _ *gen.DescribeMeterReadingRequest) (*gen.MeterReadingSupport, error) {
	return i.MeterReading, nil
}
