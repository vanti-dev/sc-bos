## Gallagher access control driver 

## API Docs
https://gallaghersecurity.github.io/cc-rest-docs/ref/index.html

## Set up

The command centre requires API keys and also client certificates to communicate. 
This requires access to the command centre configuration client, the location of which will be project specific.
- In addition to the API key, Gallagher command centre uses mutual TLS and so requires a client certificate. Generate a cert and key on the client communicating with the Gallagher API, 
    this will usually be the building controller.
- Then on the Gallagher configuration client, go to configure->services & workstations->Rest Client 1->API Key
- Copy the thumbprint of the client certificate into the Client Certificate Thumbprint box
- Make sure you are sending the client cert with the request, in addition to providing the API key in the Authorization Header

## Configuration

There is an example configuration file in the `pkg/driver/gallagher` directory. Use this as a reference for how
to configure the driver for use in and sc-bos instance.

