package config

import (
	"encoding/json"
	"errors"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/vanti-dev/bsp-ew/internal/driver/bacnet/known"
	"github.com/vanti-dev/gobacnet/property"
	bactypes "github.com/vanti-dev/gobacnet/types"
)

// Trait is the common configuration for bacnet device traits.
// Specific implementations that pull objects together into a trait should embed this type into their specific config
// types.
//
//	type OnOffConfig struct {
//	  Trait
//	  Value *ValueSource `json:"value,omitempty"`
//	}
type Trait struct {
	Name string     `json:"name,omitempty"`
	Kind trait.Name `json:"kind,omitempty"`
}

type RawTrait struct {
	Trait
	Raw json.RawMessage `json:"-"`
}

func (c *RawTrait) MarshalJSON() ([]byte, error) {
	return c.Raw, nil
}

func (c *RawTrait) UnmarshalJSON(buf []byte) error {
	if c == nil {
		*c = RawTrait{}
	}
	c.Raw = buf
	return json.Unmarshal(buf, &c.Trait)
}

// ValueSource configures a single object property as the source of some trait value.
type ValueSource struct {
	Device   *DeviceRef  `json:"device,omitempty"`
	Object   *ObjectRef  `json:"object,omitempty"`
	Property *PropertyID `json:"property,omitempty"`
}

// Lookup finds the gobacnet device, object, and property this ValueSource refers to.
func (vs ValueSource) Lookup(ctx known.Context) (bactypes.Device, bactypes.Object, property.ID, error) {
	p := property.PresentValue
	if vs.Property != nil {
		p = property.ID(*vs.Property)
	}

	if vs.Device == nil || vs.Object == nil {
		return bactypes.Device{}, bactypes.Object{}, p, errors.New("missing device or object")
	}

	device, err := vs.Device.Lookup(ctx)
	if err != nil {
		return device, bactypes.Object{}, p, err
	}
	object, err := vs.Object.Lookup(device, ctx)
	return device, object, p, err
}

type DeviceRef struct {
	id   bactypes.ObjectInstance
	name string
}

func (d DeviceRef) Lookup(ctx known.Context) (bactypes.Device, error) {
	if d.name != "" {
		return ctx.LookupDeviceByName(d.name)
	}
	return ctx.LookupDeviceByID(d.id)
}

func (d DeviceRef) MarshalJSON() ([]byte, error) {
	if d.name != "" {
		return json.Marshal(d.name)
	}
	return json.Marshal(d.id)
}

func (d *DeviceRef) UnmarshalJSON(bytes []byte) error {
	var val any
	if err := json.Unmarshal(bytes, &val); err != nil {
		return err
	}

	switch v := val.(type) {
	case float64:
		*d = DeviceRef{id: bactypes.ObjectInstance(v)}
		return nil
	case string:
		*d = DeviceRef{name: v}
		return nil
	default:
		return errors.New("invalid device ref")
	}
}

type ObjectRef struct {
	id   ObjectID
	name string
}

func (o ObjectRef) Lookup(device bactypes.Device, ctx known.Context) (bactypes.Object, error) {
	if o.name != "" {
		return ctx.LookupObjectByName(device, o.name)
	}
	return ctx.LookupObjectByID(device, bactypes.ObjectID(o.id))
}

func (o ObjectRef) MarshalJSON() ([]byte, error) {
	if o.name != "" {
		return json.Marshal(o.name)
	}
	return json.Marshal(o.id)
}

func (o *ObjectRef) UnmarshalJSON(bytes []byte) error {
	var s string
	if err := json.Unmarshal(bytes, &s); err != nil {
		return err
	}

	oid, err := ObjectIDFromString(s)
	if err == nil {
		*o = ObjectRef{id: oid}
	}
	*o = ObjectRef{name: s}
	return nil
}
