package merge

import (
	"context"
	"errors"
	"fmt"

	"github.com/vanti-dev/gobacnet"
	bactypes "github.com/vanti-dev/gobacnet/types"

	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/config"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/known"
)

func readProperty(_ context.Context, client *gobacnet.Client, known known.Context, value config.ValueSource) (any, error) {
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
	res, err := client.ReadProperty(device, req)
	if err != nil {
		return nil, err
	}
	if len(res.Object.Properties) == 0 {
		// Shouldn't happen, but has on occasion. I guess it depends how the device responds to our request
		return nil, errors.New("zero length object properties")
	}
	return res.Object.Properties[0].Data, nil
}

func readPropertyFloat64(ctx context.Context, client *gobacnet.Client, known known.Context, value config.ValueSource) (float64, error) {
	data, err := readProperty(ctx, client, known, value)
	if err != nil {
		return 0, err
	}
	switch v := data.(type) {
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

func readPropertyFloat32(ctx context.Context, client *gobacnet.Client, known known.Context, value config.ValueSource) (float32, error) {
	data, err := readProperty(ctx, client, known, value)
	if err != nil {
		return 0, err
	}
	switch v := data.(type) {
	case float32:
		return v, nil
	case float64:
		return float32(v), nil
	case uint32:
		return float32(v), nil
	case int32:
		return float32(v), nil
	}

	return 0, fmt.Errorf("unsupported conversion %T -> float64 for val %v", data, data)
}

func writeProperty(_ context.Context, client *gobacnet.Client, known known.Context, value config.ValueSource, data any, priority uint) error {
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
	return client.WriteProperty(device, req, priority)
}
