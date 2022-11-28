package messageports

import "fmt"

type Field string

const (
	Timestamp   Field = "TIMESTAMP"
	EventID     Field = "EVENTID"
	EventDesc   Field = "EVENTDESC"
	NetworkID   Field = "NETWORKID"
	NetworkDesc Field = "NETWORKDESC"
	Nc100ID     Field = "NC100ID"
	Nc100Desc   Field = "NC100DESC"
	DeviceID    Field = "DEVICEID"
	DeviceDesc  Field = "DEVICEDESC"
	CardID      Field = "CARDID"
	CardNumber  Field = "CARDNUMBER"
	Cardholder  Field = "CARDHOLDER"
	UsageCount  Field = "USAGECOUNT"
)

func (f Field) String() string {
	return string(f)
}

type Fields struct {
	Timestamp   *Time
	EventID     *uint
	EventDesc   *string
	NetworkID   *uint
	NetworkDesc *string
	Nc100ID     *uint
	Nc100Desc   *string
	DeviceID    *uint
	DeviceDesc  *string
	CardID      *uint
	CardNumber  *uint
	Cardholder  *string
	UsageCount  *uint
}

// field returns a pointer to the field in f that corresponds to field.
func (f *Fields) field(field Field) (any, error) {
	switch field {
	case Timestamp:
		return f.Timestamp, nil
	case EventID:
		return f.EventID, nil
	case EventDesc:
		return f.EventDesc, nil
	case NetworkID:
		return f.NetworkID, nil
	case NetworkDesc:
		return f.NetworkDesc, nil
	case Nc100ID:
		return f.Nc100ID, nil
	case Nc100Desc:
		return f.Nc100Desc, nil
	case DeviceID:
		return f.DeviceID, nil
	case DeviceDesc:
		return f.DeviceDesc, nil
	case CardID:
		return f.CardID, nil
	case CardNumber:
		return f.CardNumber, nil
	case Cardholder:
		return f.Cardholder, nil
	case UsageCount:
		return f.UsageCount, nil
	default:
		return nil, fmt.Errorf("unknow field %s", field)
	}
}
