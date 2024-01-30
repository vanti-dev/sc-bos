# Keycloak Login Prototype

Apps to test KeyCloak setup and client login/auth mechanisms.
This app consists of three parts:

1. A simple Vue web ui (in this dir) with features for logging in/out and performing authenticated requests (read and write) to the server. See the `[keycloak-login-prototype] UI` run configuration.
2. Configuration for setting up a server to accept and check read/write requests from this UI, see [bos-data](./bos-data). See the `[keycloak-login-prototype] BOS` run configuration.
3. A containerised KeyCloak server setup with the relevant realms, clients, and users. This exists at the root of the repo.

## Quick Start

1. Start the KeyCloak server: `podman compose start` (from root of repo)
2. Start the SC BOS server: Run the `[keycloak-login-prototype] BOS` run configuration or 
   ```shell
    go run ./cmd/bos \
      --data .data/keycloak-login-prototype \
      --appconf ./ui/keycloak-login-prototype/bos-data/appconf.json \
      --sysconf ./ui/keycloak-login-prototype/bos-data/sysconf.json
   ```
3. Start the UI: Run the `[keycloak-login-prototype] UI` run configuration or
   ```shell
    cd ui/keycloak-login-prototype
    yarn run dev
   ```
4. Visit https://localhost:8000 to accept the self-signed certificate.
5. Visit http://localhost:5173 to see the UI.

The first time you'll likely need to `yarn install` in the UI dir to grab the dependencies.

## Using the app

The UI takes you through a simplified (but complete) login flow. First you login, which invokes the OAuth flow, then you use the retrieved token as part of read and/or write requests which the server authorises.
