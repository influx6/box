package docker

import (
	gctx "context"
	"errors"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/influx6/faux/context"
	"github.com/influx6/moz/gen/filesystem"
)

// errors ...
var (
	ErrNoDockerClientProvided = errors.New("required a instance of docker.Client")
	ErrNoFilesystemProvided   = errors.New("Filesystem required")
)

// BuildImageSpell defines a structure which implements the Spell interface
// for building of docker image.
// To use the BuildImageSpell the following is required:
// 1. A filesystem to be supplied has the build context
// 2. A name for the image to be build
// As optional, a types.ImageBuildOptions if you wish to override internally generated as
// a means of setting more fields for the docker clients request.
type BuildImageSpell struct {
	Name       string
	Dockerfile string
	Tags       []string
	Client     *client.Client
	Filesystem filesystem.Filesystem
	Options    *types.ImageBuildOptions
	Callback   func(types.ImageBuildResponse)
}

// Exec executes the image creation request through the provided docker client.
// If the spell has the types.ImageBuildOptions without a filesystem, then the
// types.ImageBuildOptions will be used as is, else the filesystem will be included
// has the underline BuildContext.
func (cm *BuildImageSpell) Exec(ctx context.CancelContext) error {
	if cm.Client == nil {
		return ErrNoDockerClientProvided
	}

	if cm.Filesystem == nil {
		return ErrNoFilesystemProvided
	}

	done := make(chan struct{})
	defer close(done)

	// Cancel context if are done or if context has expired.
	mctx, cancel := gctx.WithCancel(gctx.Background())
	go func() {
		select {
		case <-mctx.Done():
			cancel()
			return
		case <-done:
			return
		}
	}()

	reader, err := cm.Filesystem.ToReader()
	if err != nil {
		return err
	}

	var imageOps types.ImageBuildOptions

	if cm.Options != nil {
		imageOps = *cm.Options
	} else {
		imageOps.PullParent = true
		imageOps.SuppressOutput = false
	}

	imageOps.Tags = append(imageOps.Tags, cm.Name)
	imageOps.Tags = append(imageOps.Tags, cm.Tags...)

	res, err := cm.Client.ImageBuild(mctx, reader, imageOps)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if cm.Callback != nil {
		cm.Callback(res)
	}

	return nil
}
