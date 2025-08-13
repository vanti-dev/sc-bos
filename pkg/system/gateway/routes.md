# Overview of gateway routing

Here are the routing rules we expect the gateway to support:

| Service                 | Payload                 | Target                              | Comments                                                                        |
|-------------------------|-------------------------|-------------------------------------|---------------------------------------------------------------------------------|
| `{gateway}/DevicesApi`  | `*`                     | `{gateway}/DevicesApi`              | The gateway should report on all known devices.                                 |
| `{gateway}/ParentApi`   | `{name: "{gateway}"}`   | `{gateway}/ParentApi`               | The gateway should report on all known children.                                |
| `{gateway}/MetadataApi` | `{name: "{name}"}`      | `{gateway}/MetadataApi`             | The gateway reports the metadata for all devices from its own cache.            |
| `{gateway}/{api}`       | `{name: "{name}", ...}` | `{nodeByAnnouncedName[name]}/{api}` | All named API requests should be forwarded to the node that announced the name. |

All other non-named API requests should be forwarded to the hub:

- HubApi
- TenantApi
- LightTestApi
- Other future APIs

`nodeByAnnouncedName` maps from the name of an announced device/child on a node to a connection to that node.
`{gateway}/DevicesApi` means the `smartcore.bos.DevicesApi` hosted on the `proxy` server.
