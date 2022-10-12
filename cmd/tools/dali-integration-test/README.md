# cmd/tools/dali-integration-test

This program can be used as an end-to-end test of the Smart Core to DALI bridging functionality.

It connects to a single DALI bus bridge in TwinCAT 3, and toggles all luminaires on the bus once a second
by sending broadcast commands. This can be used to test communication between Smart Core and DALI.
Internally, this program starts a gRPC server and connects to it, ensuring that as much of the chain
as possible is tested.

## Prerequisites
  - Your local machine has the TwinCAT 3 runtime installed.
  - You have a Beckhoff Embedded PC running the PLC project from `bsp-ew-plc`
  - The Embedded PC has one or more KL6821 DALI slices installed, and these are linked to the appropriate bus IO
    variables in the PLC project.
  - The Embedded PC is routable from your local machine over ADS (check in the ADS route editor). Make sure to note down
    its Net ID.
  - You have the correct environment variable to build and run Go programs that link to the ADS DLL using CGo -
    see the [twincat3-ads-go](https://github.com/vanti-dev/twincat3-ads-go) README for more info.

## Arguments
  - `ams-net-id` - (Required) The TwinCAT 3 AMS NetID of the PLC running the bridge, in dotted-decimal form
    e.g. `1.2.3.4.1.1`
  - `ads-port` - TwinCAT 3 Port number of the PLC instance running the bridge. Defaults to 851, the first PLC instance.
  - `bus-prefix` - The common prefix of the bridge variables used for communicating with the PLC code. Expects to find
     certain PLC variables present:
    - An `FB_DALIBrudge` at `<prefix>_bridge`
    - An `ST_DALIBridgeResponse` at `<prefix>_response`
    - An `ST_DALIBridgeNotification` at `<prefix>_notification`