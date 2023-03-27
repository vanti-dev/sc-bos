# Edge Gateways and Central Site Management

Smart Core is a distributed building operating system. What this means is that it is composed of many loosely coupled
pieces all working together and independently to satisfy all the features of a building. While this is all well and
good, it is an implementation detail of the architecture and not something operations teams or tenants of the building
should need to worry about on a day to day basis. With this in mind we aim for any features of the BOS to be accessible
from a central location including operations features like dashboards or config updated to technical features like a
buildings API. This document explains how we do that.

The problem we're trying to solve is more or less the same for both APIs and for user interfaces. Those UIs access their
data via the same APIs that we expose to tenants and other integration teams so solving this issue for the API should
solve it for UIs too.

What we do is provide a specialised controller instance that we call the gateway. This controller typically doesn't host
any drivers or automations, doesn't connect to databases or have any zones configured, but can talk to all the other
controllers in the building and proxies or consolidates those controllers apis into one place for easier consumption.

## TL;DR

Some APIs are forwarded unmodified, some are implemented locally using caches, and some are modified slightly on the way
through the gateway. These can be broken down as follows:

1. Traits and other apis that use naming are routed without modification apart from
    1. Metadata, Parent and Devices APIs are implemented via a local cache
    2. Services api modifies the names on the way through to add the node name: `drivers` on AC1 gets exposed on the
       gateway as `AC1/drivers`
2. Non-named APIs are proxied through to the hub unmodified except for the EnrollmentApi which shouldn't be proxied.

A gateway can be configured by enabling the proxy system in `system.json`

```json
{
  "systems": {
    "proxy": {
      "ignore": ["1.2.3.4:1234"]
    }
  }
}
```

## Naming and traits

One of the core principles in Smart Core is the name, every device, every controller, zones, automations, concepts, they
all have a name. We can use this to fairly easily implement proxying on the gateway, if AC1 has a device `light1` and
AC2 has a device `light2` then the gateway when asked to turn light1 off knows to redirect that request to AC1.

The gateway populates its routing table when it is enrolled with a hub. The gateway queries the hub for the list of
enrolled nodes then asks each node for the list of named children. The traits and other apis these children implement
become the routing table for the gateway.

It's worth noting here that only if a name advertises, via the parent and/or metadata traits, that they implements a
particular api will the gateway add that route to its table. There are some exceptions for named apis that don't fit
with traits but not many. If your device isn't available on the gateway then it's possible that it doesn't announce its
traits correctly on its own controller.

As the routing table is auto-populated from the cohort and parent apis it is possible to enter into an infinite loop
when more than one gateway is enrolled with a hub. In the future we may detect this situation and avoid it but for now
you'll need to configure the `"ignore"` property of your `proxy` system to include all the addresses for nodes that
should not be proxied, i.e. every gateway. To be safe, include the current gateway too.

## Non-named APIS

Some APIs don't follow the standard naming pattern of Smart Core, each of these is implemented in a specific way.

### Alerts

The alert api is exposed via the hub name on the gateway, for example `gateway.ListAlerts({name: "hub"})` should return
all the alerts stored on the hub. This is different to normal operation as usually you'd use the empty name to get the
controller alerts, but in this case the empty name resolves to the gateway which may or may not have its own alerts,
either stored locally or remotely. To allow for this you should use the name of the hub when requesting global alerts.

Other nodes alerts will not be surfaced via the gateway, these nodes should be configured to store alerts in the hub to
have them show via the gateway API

### Hub

The hub api is exposed by the gateway backed by the cohort hub, asking the gateway for the list of nodes will cause the
gateway to ask the hub and return those results. This implies that a gateway cannot also be a hub. To enable this
behaviour configure the storage type for the hub system to `"proxy"` in `system.json`

```json
{
  "name": "EG-01",
  "systems": {
    "hub": {"storage": {"type": "proxy"}},
    ...
  }
}
```

### Tenant

The tenant api follows a similar pattern to the hub proxy, you configure the storage engine for the tenant system as
type `proxy`.

```json
{
  "name": "EG-01",
  "systems": {
    "tenants": {"storage": {"type": "proxy"}},
    ...
  }
}
```
