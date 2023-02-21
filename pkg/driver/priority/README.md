# Priority Device Driver

This driver adapts an existing device to add configurable priority levels to writes. Each priority level is announced as
a new named device implementing the traits mentioned in the config. An additional `priority` device will also be created
whose writes write to the default priority level, and reads come from the real device. This is the device you should
typically use when connecting to the device.

The driver config looks like this:

```json5
{
  "drivers": [
    {
      "type": "priority",
      "suffix": "priority", // announces {name}/priority
      "separator": "/", // separator used between {name}, {suffix}, and {slots}
      // announced priority levels, highest priority first. Announced names are {name}/priority/{slot}
      "slots": ["1", "2", "3", "4", "5"],
      "defaultSlot": "3", // the slot the suffix device uses for writes, defaults to the middle slot
      "devices": [
        // These values are Metadata. The metadata will be associated with the priority device {name}/priority
        {"name": "MyDevice", "traits": [{"name": "smartcore.traits.Light"}]}
      ]
    }
  ]
}
```

## Caveats

Not all traits are implemented, see driver.go/newPriorityClient for the trait clients that are available.
