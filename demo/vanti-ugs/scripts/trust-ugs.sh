#!/usr/bin/env bash

docker cp vanti-ugs-sc-ugs-1:/data/grpc-self-signed.cert.pem /tmp/vanti-ugs-sc-ugs-1-grpc-self-signed.cert.pem
sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain /tmp/vanti-ugs-sc-ugs-1-grpc-self-signed.cert.pem