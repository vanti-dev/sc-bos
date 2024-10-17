# Smart Core Demos

This demo directory contains a collection of demos that showcase the capabilities of Smart Core.

In each folder, there is a demo for a specific site, each with its own config and ui assets. This readme will use the 
`vanti-ugs` demo as an example to explain how to run the demos. `vanti-ugs` is the Vanti Birmingham office instance of 
Smart Core.

Each folder has one or more Dockerfile-* files which are necessary to build the Docker images required for 
running that specific example. This only needs to be done on first run and if the demo is updated. 
The images are then pushed to a Docker registry, so they can be accessed from anywhere the demo needs to be set up.  
For example, to build the `vanti-ugs` images, run these commands from the repo root dir:

`docker build --tag "smartcoredemo.azurecr.io/sc-ugs" -f .\demo\vanti-ugs\Dockerfile-Ugs .`
`docker build --tag "smartcoredemo.azurecr.io/seed-db" -f demo/vanti-ugs/Dockerfile-SeedDb .`

Then they need to be pushed to a Docker registry, at the moment this is on Azure (which requires the Azure CLI tool to log in). 

`az acr login --name smartcoredemo`
`docker push smartcoredemo.azurecr.io/sc-ugs`
`docker push smartcoredemo.azurecr.io/seed-db`

To run the demo, you need to have Docker installed on your machine. That should be the only prerequisite as we want non-devs
to be able to run this demo. Also, the Docker daemon must be running, on Windows you can just press the Windows key and 
type docker and start Docker desktop.

Then, when Docker is running, start the container using the compose file in the `vanti-ugs` folder:

`docker compose -f docker-compose-ugs.yml up`

Then you can access the site at `https://localhost:8443/` in your browser. (for the Vanti example, for any other examples
refer to the `listenHTTPS` setting in the relevant system.conf.json)

You can also use this demo to test against the API, for example, to get the list of children of the root node, 
run the `client-parent` tool in `cmd/tools` and it will print the devices in `vanti-ugs` 
set the insecure-skip-verify flag to true as the demo uses a self-signed certificate. 