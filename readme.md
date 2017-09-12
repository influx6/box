Box
-----

![Box](./media/pkgbox.png)


Box provides a commandline to for deploying binaries files using docker.

# What we want

```bash
> box create io ./cmd/geth/geth-bin
Creating docker image for binary: "geth-bin"
Tagging docker image as "geth-bin"
Using docker host: /var/socket/docker.socket

Run `box push` to deliver image to docker repo.
Run `docker push REMOTE_DOCKER_HOST` to deliver image to run docker host.
Run `docker live` to get access to a local server waiting request to run binary and dashboard
```


# Notes
Alot of the initial code for this project was generated using [Moz](https://github.com/influx6/moz).
It uses a template file and specific directives found in [doc.go](./doc.go) to create the initial base.
