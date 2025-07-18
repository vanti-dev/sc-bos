package opcua

import (
	"context"
	"encoding/json"

	"github.com/gopcua/opcua/ua"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/vanti-dev/sc-bos/pkg/driver/opcua/config"
	"github.com/vanti-dev/sc-bos/pkg/driver/opcua/conv"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

type Meter struct {
	gen.UnimplementedMeterApiServer
	gen.UnimplementedMeterInfoServer

	energyValue *resource.Value // *gen.MeterReading
	logger      *zap.Logger
	meterConfig config.MeterConfig
	scName      string
}

func readMeterConfig(raw []byte) (cfg config.MeterConfig, err error) {
	err = json.Unmarshal(raw, &cfg)
	return
}

func newMeter(n string, config config.RawTrait, l *zap.Logger) (*Meter, error) {
	cfg, err := readMeterConfig(config.Raw)
	if err != nil {
		return nil, err
	}
	return &Meter{
		energyValue: resource.NewValue(resource.WithInitialValue(&gen.MeterReading{}), resource.WithNoDuplicates()),
		logger:      l,
		meterConfig: cfg,
		scName:      n,
	}, nil
}

func (m *Meter) GetMeterReading(_ context.Context, _ *gen.GetMeterReadingRequest) (*gen.MeterReading, error) {
	return m.energyValue.Get().(*gen.MeterReading), nil
}

func (m *Meter) PullMeterReadings(_ *gen.PullMeterReadingsRequest, server gen.MeterApi_PullMeterReadingsServer) error {
	for value := range m.energyValue.Pull(server.Context()) {
		err := server.Send(&gen.PullMeterReadingsResponse{Changes: []*gen.PullMeterReadingsResponse_Change{
			{
				Name:         m.scName,
				ChangeTime:   timestamppb.New(value.ChangeTime),
				MeterReading: m.energyValue.Get().(*gen.MeterReading),
			},
		}})
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Meter) DescribeMeterReading(context.Context, *gen.DescribeMeterReadingRequest) (*gen.MeterReadingSupport, error) {
	return &gen.MeterReadingSupport{
		UsageUnit: m.meterConfig.Unit,
	}, nil
}

func (m *Meter) handleMeterEvent(node *ua.NodeID, value any) {

	if m.meterConfig.Usage != nil && NodeIdsAreEqual(m.meterConfig.Usage.NodeId, node) {
		v, err := conv.Float32Value(value)
		if err != nil {
			m.logger.Error("failed to convert value", zap.Error(err))
		}
		_, _ = m.energyValue.Set(&gen.MeterReading{
			Usage:   v,
			EndTime: timestamppb.Now(),
		})
	}
}
