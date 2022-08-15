App Server executable
=====================

## Configuration
The App Server takes its configuration from a combination of command line arguments and configuration files.

### Config Directory
 - `system.json` - core configuration for the App Server - e.g. logging, database connection and so on
 - `pki/`
   - `root.cert.pem` - Root CA certificate(s) for this Smart Core installation. Must be manually provided.
   - `enrollment-ca.cert.pem` - Intermediate CA certificate for enrolling controllers. Must be part of a chain starting
      at `root.cert.pem`. Must be manually provided.
   - `enrollment-ca.key.pem` - Private key file for `enrollment-ca.cert.pem`. 