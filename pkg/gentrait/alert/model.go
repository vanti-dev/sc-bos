package alert

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

type Model struct {
	alerts []*resource.Value // of *gen.Alert
}

func NewModel(opts ...resource.Option) *Model {
	defaultOptions := []resource.Option{resource.WithInitialValue(&gen.Alert{})}
	value := resource.NewValue(append(defaultOptions, opts...)...)
	_, _ = value.Set(&gen.Alert{}, resource.InterceptBefore(func(old, new proto.Message) {
		oldVal := old.(*gen.Alert)
		newVal := new.(*gen.Alert)
		now := value.Clock().Now()
		if oldVal.CreateTime == nil {
			newVal.CreateTime = timestamppb.New(now)
		}
		if newVal.ResolveTime == nil {
			newVal.ResolveTime = timestamppb.New(now)
		}
	}))
	return &Model{}
}

func (m *Model) GetAllAlerts() []*gen.Alert {
	var results []*gen.Alert

	for _, v := range m.alerts {
		results = append(results, v.Get().(*gen.Alert))
	}
	return results
}

func (m *Model) AddAlert(a *gen.Alert, opts ...resource.WriteOption) {
	value := resource.NewValue(resource.WithInitialValue(&gen.Alert{}))
	value.Set(a, opts...)
	m.alerts = append(m.alerts, value)
}
