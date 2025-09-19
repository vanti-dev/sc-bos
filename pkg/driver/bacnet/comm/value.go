package comm

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"sync"

	"go.uber.org/multierr"

	"github.com/vanti-dev/gobacnet"
	"github.com/vanti-dev/gobacnet/property"
	bactypes "github.com/vanti-dev/gobacnet/types"
	"github.com/vanti-dev/gobacnet/types/objecttype"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/ctxerr"

	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/config"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/known"
)

func ReadProperty(ctx context.Context, client *gobacnet.Client, known known.Context, value config.ValueSource) (any, error) {
	device, object, property, err := value.Lookup(known)
	if err != nil {
		return nil, err
	}

	req := bactypes.ReadPropertyData{
		Object: bactypes.Object{
			ID: object.ID,
			Properties: []bactypes.Property{
				{ID: property, ArrayIndex: bactypes.ArrayAll},
			},
		},
	}
	res, err := client.ReadProperty(ctx, device, req)
	if err != nil {
		return nil, ctxerr.Cause(ctx, err)
	}
	if len(res.Object.Properties) == 0 {
		// Shouldn't happen, but has on occasion. I guess it depends how the device responds to our request
		return nil, errors.New("zero length object properties")
	}
	return value.Scaled(res.Object.Properties[0].Data), nil
}

type key struct {
	did bactypes.ObjectInstance
	oid bactypes.ObjectID
	pid property.ID
}

// ReadPropertiesChunked is like readProperties but splits values into chunks of at most chunkSize that are executed in parallel.
func ReadPropertiesChunked(ctx context.Context, client *gobacnet.Client, known known.Context, chunkSize int, values ...config.ValueSource) []any {
	if chunkSize == 0 {
		return ReadProperties(ctx, client, known, values...)
	}

	var wg sync.WaitGroup
	chunkCount := int(math.Ceil(float64(len(values)) / float64(chunkSize)))
	wg.Add(chunkCount)
	n := int(math.Ceil(float64(len(values)) / float64(chunkCount)))

	results := make([]any, len(values))

	for i := range chunkCount {
		from, to := i*n, (i+1)*n
		if to > len(values) {
			to = len(values)
		}
		go func() {
			defer wg.Done()
			props := ReadProperties(ctx, client, known, values[from:to]...)
			copy(results[from:to], props)
		}()
	}

	wg.Wait()
	return results
}

func ReadProperties(ctx context.Context, client *gobacnet.Client, known known.Context, values ...config.ValueSource) []any {
	res := make([]any, len(values))
	for i := range res {
		res[i] = ErrPropNotFound
	}

	resIndexes := make(map[key][]int)

	devices := make(map[bactypes.ObjectInstance]bactypes.Device)
	reqsPerDevice := make(map[bactypes.ObjectInstance]*bactypes.ReadMultipleProperty)

	for i, value := range values {
		device, object, prop, err := value.Lookup(known)
		if err != nil {
			res[i] = err
			continue
		}

		req, ok := reqsPerDevice[device.ID.Instance]
		if !ok {
			req = &bactypes.ReadMultipleProperty{}
			reqsPerDevice[device.ID.Instance] = req
			devices[device.ID.Instance] = device
		}

		// it's really unlikely that you're asking for multiple properties of the same object, but if you are,
		// the following should work anyway

		k := key{device.ID.Instance, object.ID, prop}
		resIndexes[k] = append(resIndexes[k], i)
		req.Objects = append(req.Objects, bactypes.Object{
			ID: object.ID,
			Properties: []bactypes.Property{
				{ID: prop, ArrayIndex: bactypes.ArrayAll},
			},
		})
	}

	for id, req := range reqsPerDevice {
		readMultiProperties(ctx, client, devices[id], *req, resIndexes, res)
	}

	for i, v := range res {
		res[i] = values[i].Scaled(v)
	}

	return res
}

func readMultiProperties(ctx context.Context, client *gobacnet.Client, device bactypes.Device, req bactypes.ReadMultipleProperty, resIndexes map[key][]int, res []any) {
	multiRes, err := client.ReadMultiProperty(ctx, device, req)
	if err != nil {
		// todo: be more conservative about which errors we try individual property reads for
		err = ctxerr.Cause(ctx, err)
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			// stop early as ctx is done anyway
			for _, object := range req.Objects {
				for _, prop := range object.Properties {
					k := key{device.ID.Instance, object.ID, prop.ID}
					for _, i := range resIndexes[k] {
						res[i] = err
					}
				}
			}
			return
		}

		// check if there are errors recorded for individual properties
		var propErrLen int
		errs := multierr.Errors(err)
		for _, err := range errs {
			var propErr bactypes.PropertyAccessError
			if errors.As(err, &propErr) {
				propErrLen++
				k := key{device.ID.Instance, propErr.ObjectID, propErr.Property}
				for _, i := range resIndexes[k] {
					res[i] = err
				}
			}
		}
		if propErrLen != len(errs) {
			// read the properties one at a time as the multi read failed
			for _, object := range req.Objects {
				for _, prop := range object.Properties {
					oneRes, err := client.ReadProperty(ctx, device, bactypes.ReadPropertyData{
						Object: bactypes.Object{
							ID:         object.ID,
							Properties: []bactypes.Property{prop},
						},
					})
					if err != nil {
						k := key{device.ID.Instance, object.ID, prop.ID}
						for _, i := range resIndexes[k] {
							res[i] = ctxerr.Cause(ctx, err)
						}
						continue
					}
					multiRes.Objects = append(multiRes.Objects, oneRes.Object)
				}
			}
		}
	}

	for _, object := range multiRes.Objects {
		for _, prop := range object.Properties {
			k := key{device.ID.Instance, object.ID, prop.ID}
			for _, i := range resIndexes[k] {
				res[i] = prop.Data
			}
		}
	}
}

func AnyValue(data any) (any, error) {
	if err, ok := data.(error); ok {
		return nil, err
	}
	return data, nil
}

func Float64Value(data any) (float64, error) {
	switch v := data.(type) {
	case error:
		return 0, v
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case uint32:
		return float64(v), nil
	case int32:
		return float64(v), nil
	}

	return 0, fmt.Errorf("unsupported conversion %T -> float64 for val %v", data, data)
}

func Float32Value(data any) (float32, error) {
	switch v := data.(type) {
	case error:
		return 0, v
	case float32:
		return v, nil
	case float64:
		return float32(v), nil
	case uint32:
		return float32(v), nil
	case int32:
		return float32(v), nil
	}

	return 0, fmt.Errorf("unsupported conversion %T -> float32 for val %v", data, data)
}

func BoolValue(data any) (bool, error) {
	switch v := data.(type) {
	case error:
		return false, v
	case bool:
		return v, nil
	case int, int8, int16, int32, uint, uint8, uint16, uint32:
		return v == 1, nil
	}

	return false, fmt.Errorf("unsupported conversion %T -> bool for val %v", data, data)
}

func IntValue(data any) (int64, error) {
	switch v := data.(type) {
	case error:
		return 0, v
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	}

	return 0, fmt.Errorf("unsupported conversion %T -> int for val %v", data, data)
}

func EnumValue(data any) (bactypes.Enumerated, error) {
	switch v := data.(type) {
	case error:
		return 0, v
	case uint8:
		return bactypes.Enumerated(v), nil
	case uint16:
		return bactypes.Enumerated(v), nil
	case uint32:
		return bactypes.Enumerated(v), nil
	case int8:
		return bactypes.Enumerated(v), nil
	case int16:
		return bactypes.Enumerated(v), nil
	case int32:
		return bactypes.Enumerated(v), nil
	}

	return 0, fmt.Errorf("unsupported conversion %T -> bactypes.Enumerated for val %v", data, data)
}

func StringValue(data any) (string, error) {
	switch v := data.(type) {
	case error:
		return "", v
	case string:
		return v, nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

func BitStringValue(data any) (bactypes.BitString, error) {
	switch v := data.(type) {
	case error:
		return bactypes.BitString{}, v
	case bactypes.BitString:
		return v, nil
	}

	return bactypes.BitString{}, fmt.Errorf("unsupported conversion %T -> bactypes.BitString for val %v", data, data)
}

func WriteProperty(ctx context.Context, client *gobacnet.Client, known known.Context, value config.ValueSource, data any, priority uint) error {
	device, object, property, err := value.Lookup(known)
	if err != nil {
		return err
	}

	data = massageValueForWrite(device, object, property, data)

	req := bactypes.ReadPropertyData{
		Object: bactypes.Object{
			ID: object.ID,
			Properties: []bactypes.Property{
				{
					ID:         property,
					ArrayIndex: bactypes.ArrayAll,
					Data:       data,
				},
			},
		},
	}
	writePriority := known.GetDeviceDefaultWritePriority(device.ID.Instance)

	if priority > 0 {
		writePriority = priority // allow overriding the device default priority with one given if non-zero
	}
	err = client.WriteProperty(ctx, device, req, writePriority)
	return ctxerr.Cause(ctx, err)
}

// massageValueForWrite converts value to a more correct type for the given BACnet object and property.
// For example BinaryValue.PresentValue is defined to be of type enumeration, which looks very much like a uint32, but not exactly so we convert it here.
func massageValueForWrite(_ bactypes.Device, obj bactypes.Object, prop property.ID, value any) any {
	switch obj.ID.Type {
	case objecttype.BinaryValue, objecttype.BinaryOutput:
		switch prop {
		case property.PresentValue:
			switch v := value.(type) {
			case bool:
				if v {
					return bactypes.Enumerated(1)
				} else {
					return bactypes.Enumerated(0)
				}
			case float32:
				return bactypes.Enumerated(v)
			case float64:
				return bactypes.Enumerated(v)
			case int:
				return bactypes.Enumerated(v)
			case int8:
				return bactypes.Enumerated(v)
			case int16:
				return bactypes.Enumerated(v)
			case int32:
				return bactypes.Enumerated(v)
			case int64:
				return bactypes.Enumerated(v)
			case uint:
				return bactypes.Enumerated(v)
			case uint32:
				return bactypes.Enumerated(v)
			case uint64:
				return bactypes.Enumerated(v)
			}
		}
	}
	return value
}

// EngineeringUnits is the BACnet engineering units type (uint16).
type EngineeringUnits uint16

// BACnet Engineering Units 0...129 (from ASHRAE 135 / OPC UA mapping).
const (
	EngineeringUnitsSquareMetres             EngineeringUnits = 0
	EngineeringUnitsSquareFeet               EngineeringUnits = 1
	EngineeringUnitsMilliAmperes             EngineeringUnits = 2
	EngineeringUnitsAmperes                  EngineeringUnits = 3
	EngineeringUnitsOhms                     EngineeringUnits = 4
	EngineeringUnitsVolts                    EngineeringUnits = 5
	EngineeringUnitsKilovolts                EngineeringUnits = 6
	EngineeringUnitsMegavolts                EngineeringUnits = 7
	EngineeringUnitsVoltAmperes              EngineeringUnits = 8
	EngineeringUnitsKilovoltAmperes          EngineeringUnits = 9
	EngineeringUnitsMegavoltAmperes          EngineeringUnits = 10
	EngineeringUnitsDegreesPhase             EngineeringUnits = 14
	EngineeringUnitsPowerFactor              EngineeringUnits = 15
	EngineeringUnitsJoules                   EngineeringUnits = 16
	EngineeringUnitsKilojoules               EngineeringUnits = 17
	EngineeringUnitsWattHours                EngineeringUnits = 18
	EngineeringUnitsKilowattHours            EngineeringUnits = 19
	EngineeringUnitsBtus                     EngineeringUnits = 20
	EngineeringUnitsTherms                   EngineeringUnits = 21
	EngineeringUnitsTonsPerHour              EngineeringUnits = 22
	EngineeringUnitsHertz                    EngineeringUnits = 27
	EngineeringUnitsMillimetres              EngineeringUnits = 30
	EngineeringUnitsMetres                   EngineeringUnits = 31
	EngineeringUnitsInches                   EngineeringUnits = 32
	EngineeringUnitsFeet                     EngineeringUnits = 33
	EngineeringUnitsKilograms                EngineeringUnits = 39
	EngineeringUnitsPoundsMass               EngineeringUnits = 40
	EngineeringUnitsTonsMass                 EngineeringUnits = 41
	EngineeringUnitsWatts                    EngineeringUnits = 73
	EngineeringUnitsKilowatts                EngineeringUnits = 74
	EngineeringUnitsMegawatts                EngineeringUnits = 75
	EngineeringUnitsHorsepower               EngineeringUnits = 78
	EngineeringUnitsPascals                  EngineeringUnits = 80
	EngineeringUnitsKilopascals              EngineeringUnits = 81
	EngineeringUnitsBars                     EngineeringUnits = 82
	EngineeringUnitsPoundsForcePerSquareInch EngineeringUnits = 83
	EngineeringUnitsCentimetresOfWater       EngineeringUnits = 84
	EngineeringUnitsInchesOfWater            EngineeringUnits = 85
	EngineeringUnitsMillimetresOfMercury     EngineeringUnits = 86
	EngineeringUnitsInchesOfMercury          EngineeringUnits = 88
	EngineeringUnitsDegreesCelsius           EngineeringUnits = 89
	EngineeringUnitsDegreesKelvin            EngineeringUnits = 90
	EngineeringUnitsDegreesFahrenheit        EngineeringUnits = 91
	EngineeringUnitsHours                    EngineeringUnits = 98
	EngineeringUnitsMinutes                  EngineeringUnits = 99
	EngineeringUnitsSeconds                  EngineeringUnits = 100
	EngineeringUnitsMetresPerSecond          EngineeringUnits = 101
	EngineeringUnitsFeetPerSecond            EngineeringUnits = 103
	EngineeringUnitsFeetPerMinute            EngineeringUnits = 104
	EngineeringUnitsMilesPerHour             EngineeringUnits = 105
	EngineeringUnitsCubicFeet                EngineeringUnits = 106
	EngineeringUnitsCubicMetres              EngineeringUnits = 107
	EngineeringUnitsImperialGallons          EngineeringUnits = 108
	EngineeringUnitsLitres                   EngineeringUnits = 109
	EngineeringUnitsUsGallons                EngineeringUnits = 110
	EngineeringUnitsCubicFeetPerMinute       EngineeringUnits = 112
	EngineeringUnitsCubicMetresPerHour       EngineeringUnits = 115
	EngineeringUnitsImperialGallonsPerMinute EngineeringUnits = 116
	EngineeringUnitsLitresPerSecond          EngineeringUnits = 117
	EngineeringUnitsLitresPerMinute          EngineeringUnits = 118
	EngineeringUnitsLitresPerHour            EngineeringUnits = 119
	EngineeringUnitsUsGallonsPerMinute       EngineeringUnits = 120
	EngineeringUnitsPercent                  EngineeringUnits = 129
)

var unitSymbols = map[EngineeringUnits]string{
	EngineeringUnitsSquareMetres:             "m²",
	EngineeringUnitsSquareFeet:               "ft²",
	EngineeringUnitsMilliAmperes:             "mA",
	EngineeringUnitsAmperes:                  "A",
	EngineeringUnitsOhms:                     "Ω",
	EngineeringUnitsVolts:                    "V",
	EngineeringUnitsKilovolts:                "kV",
	EngineeringUnitsMegavolts:                "MV",
	EngineeringUnitsVoltAmperes:              "VA",
	EngineeringUnitsKilovoltAmperes:          "kVA",
	EngineeringUnitsMegavoltAmperes:          "MVA",
	EngineeringUnitsDegreesPhase:             "° phase",
	EngineeringUnitsPowerFactor:              "pf",
	EngineeringUnitsJoules:                   "J",
	EngineeringUnitsKilojoules:               "kJ",
	EngineeringUnitsWattHours:                "Wh",
	EngineeringUnitsKilowattHours:            "kWh",
	EngineeringUnitsBtus:                     "BTU",
	EngineeringUnitsTherms:                   "therm",
	EngineeringUnitsTonsPerHour:              "ton/h",
	EngineeringUnitsHertz:                    "Hz",
	EngineeringUnitsPercent:                  "%",
	EngineeringUnitsMetres:                   "m",
	EngineeringUnitsFeet:                     "ft",
	EngineeringUnitsInches:                   "in",
	EngineeringUnitsMillimetres:              "mm",
	EngineeringUnitsWatts:                    "W",
	EngineeringUnitsKilowatts:                "kW",
	EngineeringUnitsMegawatts:                "MW",
	EngineeringUnitsHorsepower:               "hp",
	EngineeringUnitsPascals:                  "Pa",
	EngineeringUnitsKilopascals:              "kPa",
	EngineeringUnitsBars:                     "bar",
	EngineeringUnitsPoundsForcePerSquareInch: "psi",
	EngineeringUnitsCentimetresOfWater:       "cmH₂O",
	EngineeringUnitsInchesOfWater:            "inH₂O",
	EngineeringUnitsMillimetresOfMercury:     "mmHg",
	EngineeringUnitsInchesOfMercury:          "inHg",
	EngineeringUnitsDegreesCelsius:           "°C",
	EngineeringUnitsDegreesFahrenheit:        "°F",
	EngineeringUnitsDegreesKelvin:            "K",
	EngineeringUnitsHours:                    "h",
	EngineeringUnitsMinutes:                  "min",
	EngineeringUnitsSeconds:                  "s",
	EngineeringUnitsLitres:                   "L",
	EngineeringUnitsUsGallons:                "gal (US)",
	EngineeringUnitsImperialGallons:          "gal (Imp)",
	EngineeringUnitsLitresPerSecond:          "L/s",
	EngineeringUnitsLitresPerMinute:          "L/min",
	EngineeringUnitsLitresPerHour:            "L/h",
	EngineeringUnitsUsGallonsPerMinute:       "gpm (US)",
	EngineeringUnitsImperialGallonsPerMinute: "gpm (Imp)",
	EngineeringUnitsCubicMetres:              "m³",
	EngineeringUnitsCubicFeet:                "ft³",
	EngineeringUnitsCubicMetresPerHour:       "m³/h",
	EngineeringUnitsCubicFeetPerMinute:       "cfm",
	EngineeringUnitsMetresPerSecond:          "m/s",
	EngineeringUnitsFeetPerSecond:            "ft/s",
	EngineeringUnitsFeetPerMinute:            "ft/min",
	EngineeringUnitsMilesPerHour:             "mph",
	EngineeringUnitsKilograms:                "kg",
	EngineeringUnitsPoundsMass:               "lb",
	EngineeringUnitsTonsMass:                 "ton",
}

// String implements fmt.Stringer.
func (u EngineeringUnits) String() string {
	if s, ok := unitSymbols[u]; ok {
		return s
	}
	return "Unknown(" + strconv.Itoa(int(u)) + ")"
}
