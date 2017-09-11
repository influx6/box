Dockish
-----------
Dockish provides a higher level API abstraction for executing docker related operations (creating images, containers) with functionally composed filesystems.

## Project Idea
Dockish is rather a simple idea, which sterns from the need to natively deploy code without over dependence on varying devops tools like puppet and chef. These tools are awesome and definitely meet a need, but with ever growing numbers of tools, deployment is become harder.  

Docker however does provide an alternative approach that can be mixed with different tools, or used natively with it's Swarm API, to deploy code has containers with ease.

Dockish provides a tiny but useful twist in this approach, I believe the more compact a deployment strategy is the better. Hence, dockish is meant to provide a higher level abstraction on top of the docker engine API, to interface with a docker client locally or remotely installed
, both to build and interact with images and containers.

More so, Dockish adopts the idea that all files related to the building process should be as innately available within a single binary as possible(including scripts, files, binaries), which would surprise you that docker does the same thing by zipping your files to it's daemon to then process.

Hence all docker commands that may deal with the filesystem will be interacted with through a in-memory filesystem, which would allow us easily move our deployment as a single binary file without worry.

More so, allow us bundle all necessary pieces for the creation of docker images, which in turn will be used to create docker containers or swarms right within code.

Dockish is a very simple idea, but if well combined can ease alot of pain in deployment using docker and hopefully kubernetes in the future.

*Dockish is never meant to break ground or replace the massive tools already available in Devops, but if one is able to simply deploy a fleet of services with a single binary using dockers swarm(and in future Kubernetes), then life can be far more easier.*

# API Idea
Experimenting with different API structures has proven useful, below is a sample of the possible look and feel of a dockish API to build an image. See [ImageBuild](./example/imagebuild) that actually works.

```go
client, err := client.NewEnvClient()
if err != nil {
	panic(err)
}

docker := dockish.New(client)

base := filesystem.FileSystem(
	filesystem.File(
		"Dockerfile",
		filesystem.Content(`
			FROM alpine:latest
			MAINTAINER Alexander Ewetumo <trinoxf@gmail.com>

			CMD ["/bin/ls"]
		`),
	),
)

builder, err := docker.BuildImage("wombat", filesystem.GzipTarFS(base), nil)
if err != nil {
  //...
}

err := builder.Exec(context.Background());
```


# Notes
Alot of the initial code for this project was generated using [Moz](https://github.com/influx6/moz).
It uses a template file and specific directives found in [doc.go](./doc.go) to create the initial base.
