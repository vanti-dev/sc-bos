package auto

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vanti-dev/sc-bos/pkg/gentrait/securityevent"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
)

func SecurityEventAuto(model *securityevent.Model) *service.Service[string] {
	slc := service.New(service.MonoApply(func(ctx context.Context, _ string) error {
		go func() {
			timer := time.NewTimer(30 * time.Second)
			for {
				select {
				case <-ctx.Done():
					return
				case <-timer.C:
					_, _ = model.GenerateSecurityEvent(timestamppb.Now())
					timer.Reset(30 * time.Second)
				}
			}
		}()
		return nil
	}), service.WithParser(func(data []byte) (string, error) {
		return string(data), nil
	}))
	_, _ = slc.Configure([]byte{})
	return slc
}
