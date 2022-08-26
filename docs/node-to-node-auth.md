Node-to-Node authentication & authorization
===========================================

Authentication & Authorization between nodes is handled by mutual TLS.

## Enrollment Procedure
When an Area Controller is newly installed, it cannot communicate with the rest of the Smart Core network. It must first
be commissioned.
The system uses a Trust on First Use security model, similar to SSH.
Area Controllers are enrolled from the Building Controller. The administrator will specify the IP address and intended
name of
the Area Controller. The administrator will then be presented with a key signature to confirm out-of-band. When the
administrator approves the enrollment, a CSR from the Area Controller will be signed and returned, and enrollment is
complete.

The Building Controller has its own simple CA, which is used only for Smart Core.

## Nodes - Server Ports

The blow table lists all the exposed ports and their purpose along with how certificates are created for each node.

| Node                |  Port | Network  | Purpose              | Cert signer                                                |
|:--------------------|------:|----------|----------------------|------------------------------------------------------------|
| Building Controller |   443 | Internal | UI hosting, grpc-web | Self or via local filesystem                               |
| Building Controller | 23557 | Internal | Smart Core (gRPC)    | Self or via local filesystem                               |
| Edge Gateway        | 23557 | Tenant   | Smart Core (gRPC)    | Publicly trusted CA                                        |
| Area Controller     |   443 | Internal | UI hosting, grpc-web | First boot: self<br/>After enrollment: Building Controller |
| Area Controller     | 23557 | Internal | Smart Core (gRPC)    | First boot: self<br/>After enrollment: Building Controller |

## Certificate Rotation

TODO

Requirements:

1. The services shall be able to swap certificates without a restart.
2. Certificate rotation shall be completely automatic, assuming the required services are operational.
