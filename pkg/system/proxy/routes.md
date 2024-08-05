# Overview of proxy routing

Here are the routing rules we expect the proxy to support:

| Service              | Payload                        | Target                              | Comments                                                                        |
|----------------------|--------------------------------|-------------------------------------|---------------------------------------------------------------------------------|
| `{proxy}/ServiceApi` | `{name: "{name}/{type}", ...}` | `{nodeByHubName[name]}/ServiceApi`  | Adjusting the payload to `{name: "{type}", ... }`                               |
| `{proxy}/ServiceApi` | `{name: "{type}", ...}`        | `{proxy}/ServiceApi`                | The proxy should report its own services.                                       |
| `{proxy}/DeviceApi`  | `*`                            | `{proxy}/DeviceApi`                 | The proxy should report on all known devices.                                   |
| `{proxy}/ParentApi`  | `{name: "{proxy}"}`            | `{proxy}/ParentApi`                 | The proxy should report on all known children.                                  |
| `{proxy}/{api}`      | `{name: "{name}", ...}`        | `{nodeByAnnouncedName[name]}/{api}` | All named API requests should be forwarded to the node that announced the name. |

All other non-named API requests should be forwarded to the hub:

- HubApi
- TenantApi
- LightTestApi
- Other future APIs

`nodeByHubName` maps from node name (as recorded on the hub) to a connection to that node.
This should also include the hub and the proxy. The hub needs special processing because a hub does not return itself
when asked for enrolled nodes.
`nodeByAnnouncedName` maps from the name of an announced device/child on a node to a connection to that node.
`{proxy}/DeviceApi` means the `smartcore.bos.DeviceApi` hosted on the `proxy` server.
