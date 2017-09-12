package main

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/influx6/box/docker"
	"github.com/influx6/moz/gen/filesystem"
	"github.com/moby/moby/client"
)

func main() {

	client, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	image := docker.New(client)

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

	builder, err := image.BuildImage("wombat", filesystem.GzipTarFS(base))

	if err != nil {
		panic(err)
	}

	doFunc := func(res types.ImageBuildResponse) error {
		defer res.Body.Close()

		var buf bytes.Buffer
		io.Copy(&buf, res.Body)

		fmt.Printf("Response: %+s\n", buf.Bytes())
		return nil
	}

	if err := builder.Exec(context.Background(), doFunc); err != nil {
		panic(err)
	}
}
