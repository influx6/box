Dockish
-----------
Dockish provides a commandline to for deploying binaries files using docker.

# API Idea

# What we want

```bash
> dockish create io ./cmd/geth/geth-bin
Creating docker image for binary: "geth-bin"
Tagging docker image as "geth-bin"
Using docker host: /var/socket/docker.socket

Run `dockish push` to deliver image to docker repo.
```

```bash

```

# Notes
Alot of the initial code for this project was generated using [Moz](https://github.com/influx6/moz).
It uses a template file and specific directives found in [doc.go](./doc.go) to create the initial base.
