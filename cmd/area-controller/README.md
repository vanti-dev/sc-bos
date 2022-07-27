Area Controller Command
=======================

This command runs an Area Controller instance. 

## Data Directory
 - `private-key.pem` - private key used for Smart Core, PKCS#8 wrapped in PEM
 - `self-signed.crt` - self-signed X.509 certificate for `private-key.pem` - used before controller has enrolled
 - `enrollment/`
   - `enrollment.json` - data file generated upon enrollment
   - `ca.crt` - Root CA for the Smart Core installation
   - `cert.crt` - X.509 certificate for `../private-key.pem` signed by the Root CA
 - `cache/`
   - `publications/` - cache of management server publications, including configuration