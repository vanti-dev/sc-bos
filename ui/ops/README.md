# Smart Core Ops UI

This ui allows you to setup and manage smart core nodes (area controllers, app servers, edge gateways, etc). This is the
default ui you get when you visit any of these boxes IPs.

It's located here in the project because we don't yet know how much of it will be 'Smart Core' and how much will be
'Project', so best to be specific and pull out the commonality when we actually have proof it'll be needed.

## Getting started

Run the UI in dev mode (hot reloading, etc) with `yarn run dev` - assuming you have run `yarn install` and have the
dependencies ready. You will likely also need to `yarn install` in [../ui-gen](../ui-gen) as this package depends on
that file-based package.

We don't really know what the final scope of this app is, it started as the admin page for controllers but has evolved
into a unified operations and commissioning interface. The app allows you to manage API access, see what devices you
have, control the lights and see alerts.

The design
follows [this Figma design](https://www.figma.com/proto/5wfaoD7k13k1g0XTbdoc3q/SmartCore-Design-System-v1.0?page-id=420%3A2128&node-id=495%3A2440&viewport=202%2C130%2C0.32&scaling=min-zoom&starting-point-node-id=420%3A5995).

The different features of the UI require different capabilities from the server, here are some examples:

- Logging in with keycloak requires keycloak to be running - see [docs/install/](../../docs/install/dev.md)
- Logging in with username + password requires an area controller configured with local user accounts -
  see [area-controller/README.md](../../cmd/area-controller/README.md#local-authentication) for details
- Notifications requires a controller configured with the [alert system](../../pkg/system/alerts/README.md)
- System pages, like lights, require - or at least work better - with some actual devices to view and manipulate. These
  can be configured on the controller via the `"drivers"` config, samples of which can be found
  in [config/samples](../../config/samples). For detailed config options see the `internal/driver` sub-packages, each of
  which should provide a readme and config structure.

To connect the UI to a local area controller, create a `.env.local` file in this directory with the following:
```properties
VITE_GRPC_ENDPOINT=https://localhost:23557
```

When testing device or system pages it can be easier to point the UI at
the [Smart Core Playground](https://github.com/smart-core-os/sc-playground), updating the above config to:

```properties
VITE_GRPC_ENDPOINT=https://localhost:8443
```
