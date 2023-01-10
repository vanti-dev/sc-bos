# Mock Device Driver

This driver allows you to configure in-memory mock devices that implement specific traits, useful for testing.

The driver config looks like this:

```json
{
  "drivers": [
    {
      "type": "mock",
      "devices": [
        {"name": "MyDevice", "traits": [{"name": "smartcore.traits.Light"}]}
      ]
    }
  ]
}
```

The `devices` array in the config is modelled after the [Metadata] trait resource type, any metadata you want can be specified:

```json
{"name": "MyDevice", "membership": {"subsystem": "Lighting"}, ...}
```

## Caveats

Not all traits are implemented, see driver.go/newMockClient for the trait clients that are available.

[Metadata]: https://smart-core-os.github.io/api/traits/metadata.html#traits-metadata-proto-2
