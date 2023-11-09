# Smart Core BOS Hub Example

This example helps you to start the Smart Core BOS in a hub configuration.
The setup includes two area controllers, one building controller, and two edge gateways.

Each node starts out as an independent node, to combine them together into a cohort execute the `HubApi.EnrollHubNode`
gRPC api on the building controller, passing the ip:port of each AC and EG. See [enroll.http](./enroll.http) for
details.

As each BOS node exposes a gRPC and HTTPS API their ports would conflict so they have been configured with the following
pattern:

```plaintext
23203
^^^\\- Node number: 03
^^\--- Node type: 1 = BC, 2 = EG, 3 = AC
\\---- API type: 23 = gRPC, 8 = HTTPS
```

For example the HTTPS API on BC-01 uses port 8101, the gRPC API on AC-02 uses 23302, etc.

The building controller is configured to store enrollment data in the local bolt db. In production you should probably
use the postgres storage option.

An intellij 'Hub' run configuration (and group) have been setup to run all the nodes.

In case you want to run the nodes individually, the following commands can be used:

```shell
go run ./cmd/bos --policy-mode=off --data-dir example/config/hub/eg-02
```

Each node has a `--data-dir` flag that points to the directory containing the node's configuration and data.
