# Set up for local development

Local development of sc-bos requires a PostgreSQL db and Keycloak server to be running.
We've provided a compose file in the root of the repo which includes these services and some setup (but not all).

## Using Docker Compose

In the root of this repo, run:

```shell
podman compose up -d
# or docker-compose up -d
```

Which will start the DB and Keycloak services.
We also include a pgAdmin service to help with db inspection.

The compose file will automatically configure the DB with the relevant tables and extensions
and configure Keycloak with a Smart Core realm containing users, applications, and OAuth settings.

## Additional Setup

### PgAdmin

PgAdmin needs to be told which database to admin:

1. Login to [pg-admin](http://localhost:8881) using the username and password defined in the docker-compose file (
   admin@example.com/postgres)
2. Create a connection to the local postgres server. This needs to be done using the IP of the docker bridge network
   rather than localhost.

   It might be `172.18.0.4`, but if that doesn't work you can find the correct IP by doing:
   `docker network inspect bsp-ew_default` and looking for `bsp-ew-db-1`
