package merge

import (
	"context"
	"errors"
	"fmt"

	"github.com/vanti-dev/gobacnet"
	"github.com/vanti-dev/gobacnet/property"
	bactypes "github.com/vanti-dev/gobacnet/types"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/ctxerr"

	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/config"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/known"
)

func readProperty(ctx context.Context, client *gobacnet.Client, known known.Context, value config.ValueSource) (any, error) {
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
	return res.Object.Properties[0].Data, nil
}

type key struct {
	did bactypes.ObjectInstance
	oid bactypes.ObjectID
	pid property.ID
}

func readProperties(ctx context.Context, client *gobacnet.Client, known known.Context, values ...config.ValueSource) []any {
	res := make([]any, len(values))
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

	for _, object := range multiRes.Objects {
		for _, prop := range object.Properties {
			k := key{device.ID.Instance, object.ID, prop.ID}
			for _, i := range resIndexes[k] {
				res[i] = prop.Data
			}
		}
	}
}

func float64Value(data any) (float64, error) {
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

func float32Value(data any) (float32, error) {
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

func stringValue(data any) (string, error) {
	switch v := data.(type) {
	case error:
		return "", v
	case string:
		return v, nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

func valuesEquivalent(a, b any) bool {
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

func readPropertyFloat32(ctx context.Context, client *gobacnet.Client, known known.Context, value config.ValueSource) (float32, error) {
	data, err := readProperty(ctx, client, known, value)
	if err != nil {
		return 0, err
	}
	return float32Value(data)
}

func writeProperty(ctx context.Context, client *gobacnet.Client, known known.Context, value config.ValueSource, data any, priority uint) error {
	device, object, property, err := value.Lookup(known)
	if err != nil {
		return err
	}

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
	err = client.WriteProperty(ctx, device, req, priority)
	return ctxerr.Cause(ctx, err)
}
