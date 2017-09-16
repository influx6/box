Box[Experiemental]
-----
![Box](./media/pkgbox.png)

Box provides a commandline to for deploying binaries files using docker installed on remote systems.

Please note box is an idea, what you see below you is a story of an idea for what box should be like interacting with the service.

## Install
Install linux server into your linux host with `wget -o https://box.io/install.sh | sh`

Supported OS deployments: `Ubuntu i386/amd64`, `ArchLinux i386/amd64`.


# Standalone Box service

## What we want

```bash
> box register -n thunder.io 192.40.30.90:5060
Adding service `thunder.io` has box provider
Service registered to remote host `192.40.30.90:5060`.
Negotiate ssh encrytion keys with remote box service.
Service `thunder.io` ready for operation. :)
```

### What should happen

- Local box binary adds hosts into list of remote box host with name
- Local box negotiate pub keys with remote host

# Standalone Stateless Binary

## What we want

```bash
> box func io ./cmd/geth/geth-bin
Creating docker image for binary: "geth-bin"
Tagging docker image as "geth-bin"
Using docker host: /var/socket/docker.socket
```

Run `box push geth-bin` to deliver image to docker repo.

Run `docker push -r REMOTE_DOCKER_HOST geth-bin` to deliver image to run docker host.

Run `docker live geth-bin` to get access to a local server waiting request to run binary and dashboard

### What should happen

- Create custom binary that wraps `geth-bin` binary with server binary that provides:
  1. `geth-bin/exec` to process incoming request.
  2. `geth-bin/stat` to return stats on requests executed.

- Create docker image with `geth-bin` tag for deployment.


# Stateless Binary Service

## What we want

```bash
> box service new thunder.io
Create service for binary functions
Service launch on https://127.0.0.1:3000/

> box service add -n thunder.io ./cmd/geth/geth-bin
Push binary `geth-bin` to service `thunder.io(https://127.0.0.1)`
Registering binary has `geth-bin` with route `/service/geth-bin`
Preparing data collectors for `geth-bin` requests
Ready for requests at `thunder.io/service/geth-bin`.
```

### What should happen

- Create push binary `geth-bin` to service local directory:
- Add binary path for `geth-bin` as `/services/geth-bin`
- Setup underline records for `geth-bin` request.

- Create docker image with `geth-bin` tag for deployment.

# Notes
Alot of the initial code for this project was generated using [Moz](https://github.com/influx6/moz).
It uses a template file and specific directives found in [doc.go](./doc.go) to create the initial base.
