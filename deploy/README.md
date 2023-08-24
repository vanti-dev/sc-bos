Deployment and Packaging
========================

This directory contains scripts, manifests etc. relating to the packaging and deployment of Smart Core.

## Build Area Controller container image
Run `./build-container.sh` which will build and tag a container image containing an area controller executable built
from `cmd/bos` and the conductor UI (stored in `/static`).

Requires `podman`, `go` to be installed on your machine. Your Go installation must be able to authenticate against
private repos on Github. You must have a `~/.npmrc` file which will be used inside the container to fetch private
NPM dependencies.