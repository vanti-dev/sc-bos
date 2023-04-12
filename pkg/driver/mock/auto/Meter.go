package auto

import (
	"context"
	"time"

	"golang.org/x/exp/rand"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/meter"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
)

func MeterAuto(model *meter.Model) *service.Service[string] {
	slc := service.New(service.MonoApply(func(ctx context.Context, _ string) error {
		go func() {
			timer := time.NewTimer((30 * time.Second) + time.Duration(rand.Float32())*time.Minute)
			start := timestamppb.Now()
			value := rand.Float32() * 100
			for {
				select {
				case <-ctx.Done():
					return
				case <-timer.C:
					value += rand.Float32() * 100
					state := gen.MeterReading{
						Usage:     value,
						StartTime: start,
						EndTime:   timestamppb.Now(),
					}
					_, _ = model.UpdateMeterReading(&state)
					timer = time.NewTimer((30 * time.Second) + time.Duration(rand.Float32())*time.Minute)
				}
			}
		}()
		return nil
	}), service.WithParser(func(data []byte) (string, error) {
		return string(data), nil
	}))
	_, _ = slc.Configure([]byte{}) // call configure to ensure we load when start is called.
	return slc
}
