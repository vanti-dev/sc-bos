# Smart Core OPC UA driver

This package implements integration between OPC UA and Smart Core. 
The driver uses the [gopcua](https://github.com/gopcua/opcua) to communicate with OPC UA servers.

## How it works

Everything in OPC UA is a node, the Nodes we are most interested in are the Variable Nodes.
These variables represent a read/writable value in the OPC UA server.
Each Variable Node has a NodeID, which is a unique identifier for the node in the server.
In the config, the device defines which variables it wants to subscribe to.
When the underlying value of the Variable Node changes, we get an event through the channel and act, 
depending on which traits this device has configured to support.
2 different devices can subscribe to the same nodeID without an issue. 

## Traits

In the config, each device configures which traits it supports. (see config/sample.json for an example)
Each trait has its own configuration, which is used to map the OPC UA Variable Node to the trait.
