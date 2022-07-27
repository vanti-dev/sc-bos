Node-to-Node authentication & authorization
===========================================

Authentication & Authorization between nodes is handled by mutual TLS.

## Enrollment Procedure
When an Area Controller is newly installed, it cannot communicate with the rest of the Smart Core network. It must first
be commissioned. 
The system uses a Trust on First Use security model, similar to SSH.
Area Controllers are enrolled from the App Server. The administrator will specify the IP address and intended name of
the Area Controller. The administrator will then be presented with a key signature to confirm out-of-band. When the
administrator approves the enrollment, a CSR from the Area Controller will be signed and returned, and enrollment is
complete.

The App Server has its own simple CA, which is used only for Smart Core.

## Nodes - Server Ports
### App Server
 - Internal Network; Port 443
   - HTTPS: web content & gRPC-Web
   - Internally trusted certificate (commissioned externally)
 - Internal Network; Port 23557
   - gRPC
   - Self-signed certificate (the app server is root of trust)

### Smart Core Gateway
 - Tenant Network; Port 23557
   - gRPC
   - Publicly trusted certificate

### Area Controller
 - Internal Network; Port 443
   - HTTPS: web content & gRPC-Web (hosts the local commissioning interface)
   - On first boot: self-signed certificate
   - After enrollment: certificate signed by App Server
 - Internal Network; Port 23558
   - gRPC
   - On first boot: self-signed certificate
   - After enrollment: certificate signed by App Server

## Certificate Rotation

TODO

Requirements:
  1. The services shall be able to swap certificates without a restart.
  2. Certificate rotation shall be completely automatic, assuming the required services are operational.