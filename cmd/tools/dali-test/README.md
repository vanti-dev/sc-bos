# cmd/tools/dali-test

This program can be used as an end-to-end test of the Smart Core to DALI bridging functionality.

It connects to a single DALI bus bridge in TwinCAT 3, and toggles all luminaires on the bus once a second
by sending broadcast commands. This can be used to test communication between Smart Core and DALI.
Internally, this program starts a gRPC server and connects to it, ensuring that as much of the chain
as possible is tested.

## Arguments
  - `ams-net-id` - (Required) The TwinCAT 3 AMS NetID of the PLC running the bridge, in dotted-decimal form
    e.g. `1.2.3.4.1.1`
  - `ads-port` - TwinCAT 3 Port number of the PLC instance running the bridge. Defaults to 851, the first PLC instance.
  - `bus-prefix` - The common prefix of the bridge variables used for communicating with the PLC code. Expects to find
     certain PLC variables present:
    - An `FB_DALIBrudge` at `<prefix>_bridge`
    - An `ST_DALIBridgeResponse` at `<prefix>_response`
    - An `ST_DALIBridgeNotification` at `<prefix>_notification`