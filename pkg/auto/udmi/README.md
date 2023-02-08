# Auto - UDMI

This auto implements (some of) the [UDMI](https://faucetsdn.github.io/udmi/) spec for data export and control, via MQTT. Drivers are responsible for conversion to [UDMI data structures](https://faucetsdn.github.io/udmi/gencode/docs/), and expose this by implementing the [UdmiService](../../../proto/udmi.proto). That same service also allows for the [config](https://faucetsdn.github.io/udmi/docs/messages/config.html) UDMI flow, which is how control is implemented.
