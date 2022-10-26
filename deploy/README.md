# Deploying the Project

Deployment of applications and configuration, setup of the Beckhoff PLC OS, and generally getting things to a running
state is handled by the Ansible playbooks found in this directory.

## Prerequisite tools and knowledge

This project uses Community Ansible to deploy code, config, and setup the target machine. Install it by following
[these instructions](https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html#control-node-requirements)
.

The list of area controllers and other machines that we deploy to are listed in [inventory.yaml](./inventory.yaml) along
with any unique variable they have.
[This example](https://github.com/ansible/ansible/blob/devel/examples/hosts.yaml) is a good resource to learn
about the inventory.

You run ansible playbooks via a command like this:

```shell
ansible-playbook -i inventory.yaml playbooks/deploy-config.yaml
```

If you want to deploy to a single host you can use the limit arg `-l <host pattern>` option:

```shell
# Deploy to all prod hosts except floor-03. The \ is interpreted by the shell to escape the !
ansible-playbook -i inventory.yaml -l "prod:\!floor-03" playbooks/deploy-config.yaml
```

## Production deployment

### One time

Each area controller needs to have some initial setup before it will run our software. The playbook
[prepare-os.yaml](./playbooks/prepare-os.yaml) does this setup for you. It needs to be run once per Beckhoff PLC.

### Each time

TL;DR

1. `cd ui/conductor && yarn build` - Build the UI production package
   - Produces `dist` directory
2. `cd deploy && ./package.sh` - Bundle the go and ui into archives for copying to the target machines
   - Produces archives in `/tmp/bsp-ew-build`
3. `cd deploy && ansible-playbook -i inventory.yaml playbooks/deploy-config.yaml playbooks/deploy-area-controllers.yaml`
   - Transfer the bundles, build go code, start services, etc

Standard deployments involve these pieces:

1. The area controller process
2. Configuration files for the area controller
3. The area controller UI
4. Daytime lighting scheduler process (temporary)

The Go applications currently need to be built (via `go build`) on the target machine. Ansible does this for you, but
you have to prepare some things before ansible does its thing, namely building the ui and archiving the go code. There's
no reason ansible can't do these steps for you, those steps just haven't been written yet.

The reason we archive, transfer, extract, and build on the target machine is down to our use of CGO to link the twincat
libraries into the go program. CGO + FreeBSD means a load of pain which is easily avoided by building the go code on the
machine it will be running on.
