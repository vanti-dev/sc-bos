# Smart Core BACnet/IP driver

This package implements integration between BACnet/IP and Smart Core. The driver uses a
Vanti [fork of gobacnet](https://github.com/vanti-dev/gobacnet/tree/write) maintained on the `write` branch.

The driver interprets [config](config/root.go) and sets up connections to BACnet devices and associates BACnet objects
with Smart Core traits and properties.

There's a definitive sample config file in [config/testdata](config/testdata/sample.json5). The driver only supports
JSON (not json5), but the json5 helps with documenting the properties available.

**Sample BACnet driver config**

```json
{
  "type": "bacnet", "name": "MyDriverImpl",
  "devices": [
    {
      "id": 10002,
      "objects": [
        {"id": "BinaryValue:1", "trait": "smartcore.traits.OnOff"}
      ]
    }
  ],
  "traits": [
    {
      "name": "thermostat", "kind": "smartcore.traits.AirTemperature",
      "setPoint": {"device": 10002, "object": "AnalogInput:0"},
      "ambientTemperature": {"device": 10002, "object": "AnalogOutput:0"}
    }
  ]
}
```

## BACnet - Smart Core Mapping

The driver does not make any assumptions about which objects implement which traits. If an object does map well to the
semantics of a trait then you can specify this in the `devices.objects.trait` property as shown in the example above. If
an Object->Trait mapping isn't implemented it can be added into the [adapt](adapt) package where each go file contains
mappings from that BACnet object type into relevant traits.

More complicated mappings from multiple objects to a single named trait are configured via the `traits` config property
and the [merge](merge) package.

The driver also publishes a non-Smart Core gRPC API described in [bacnet.proto](rpc/bacnet.proto) that provides low
level access to BACnet services like ReadProperty and WriteProperty against configured devices.
