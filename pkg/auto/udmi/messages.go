package udmi

// PointsEvent presents the JSON payload of a UDMI Event message
// see https://faucetsdn.github.io/udmi/docs/messages/pointset.html#event
type PointsEvent map[string]PointValue

func (f PointsEvent) Equal(other PointsEvent) bool {
	if f == nil && other == nil {
		return true
	}
	if f == nil || other == nil {
		return false
	}
	if len(f) != len(other) {
		return false
	}
	for key, value := range f {
		v, ok := other[key]
		if !ok {
			return false
		}
		if value == v {
			continue
		}
		if value.PresentValue != v.PresentValue {
			return false
		}
	}
	return true
}

// PointValue is a single UDMI point value
// see https://faucetsdn.github.io/udmi/docs/messages/pointset.html#event
type PointValue struct {
	// should be a primitive value: string, bool, float...
	PresentValue any `json:"present_value"`
}

// ConfigMessage is a UDMI config message, used for control/settings
// https://faucetsdn.github.io/udmi/docs/messages/pointset.html#config
type ConfigMessage struct {
	PointSet ConfigPointSet `json:"pointset"`
}

// ConfigPointSet is a UDMI point set, for config messages
// https://faucetsdn.github.io/udmi/docs/messages/pointset.html#config
type ConfigPointSet struct {
	Points map[string]PointSetValue `json:"points"`
}

// PointSetValue is a single UDMI point with set value, for config messages
// https://faucetsdn.github.io/udmi/docs/messages/pointset.html#config
type PointSetValue struct {
	SetValue any `json:"set_value"`
}
