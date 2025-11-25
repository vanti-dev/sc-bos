package alert

import (
	"github.com/smart-core-os/sc-bos/pkg/gen"
	"github.com/smart-core-os/sc-golang/pkg/resource"
)

type Model struct {
	alerts []*resource.Value // of *gen.Alert
}

func NewModel() *Model {
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
