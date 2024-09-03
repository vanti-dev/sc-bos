# Overview of proxy routing

Here are the routing rules we expect the gateway to support:

| Service                | Payload                        | Target                              | Comments                                                                        |
|------------------------|--------------------------------|-------------------------------------|---------------------------------------------------------------------------------|
| `{gateway}/ServiceApi` | `{name: "{name}/{type}", ...}` | `{nodeByHubName[name]}/ServiceApi`  | Adjusting the payload to `{name: "{type}", ... }`                               |
| `{gateway}/ServiceApi` | `{name: "{type}", ...}`        | `{proxy}/ServiceApi`                | The gateway should report its own services.                                     |
| `{gateway}/DeviceApi`  | `*`                            | `{proxy}/DeviceApi`                 | The gateway should report on all known devices.                                 |
| `{gateway}/ParentApi`  | `{name: "{proxy}"}`            | `{proxy}/ParentApi`                 | The gateway should report on all known children.                                |
| `{gateway}/{api}`      | `{name: "{name}", ...}`        | `{nodeByAnnouncedName[name]}/{api}` | All named API requests should be forwarded to the node that announced the name. |

All other non-named API requests should be forwarded to the hub:

- HubApi
- TenantApi
- LightTestApi
- Other future APIs

`nodeByHubName` maps from node name (as recorded on the hub) to a connection to that node.
This should also include the hub and the gateway. The hub needs special processing because a hub does not return itself
when asked for enrolled nodes.
`nodeByAnnouncedName` maps from the name of an announced device/child on a node to a connection to that node.
`{proxy}/DeviceApi` means the `smartcore.bos.DeviceApi` hosted on the `proxy` server.
