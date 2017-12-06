package actions

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

// BuildFileSystemImage defines a structure which implements the Spell interface
// for building of docker image.
// To use the BuildImageSpell the following is required:
// 1. A filesystem to be supplied has the build context
// 2. A name for the image to be build
// As optional, a types.ImageBuildOptions if you wish to override internally generated as
// a means of setting more fields for the docker clients request.
type BuildFileSystemImage struct {
	Name       string
	Dockerfile string
	Tags       []string
	Client     *client.Client
	Filesystem filesystem.Filesystem
	Options    *types.ImageBuildOptions
	Callback   func(types.ImageBuildResponse) error
}

// Exec executes the image creation request through the provided docker client.
// If the spell has the types.ImageBuildOptions without a filesystem, then the
// types.ImageBuildOptions will be used as is, else the filesystem will be included
// has the underline BuildContext.
func (bm BuildFileSystemImage) Exec(ctx context.CancelContext) error {
	if bm.Client == nil {
		return ErrNoDockerClientProvided
	}

	if bm.Filesystem == nil {
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

	reader, err := bm.Filesystem.ToReader()
	if err != nil {
		return err
	}

	var imageOps types.ImageBuildOptions

	if bm.Options != nil {
		imageOps = *bm.Options
	} else {
		imageOps.PullParent = true
		imageOps.SuppressOutput = false
	}

	imageOps.Tags = append(imageOps.Tags, bm.Name)
	imageOps.Tags = append(imageOps.Tags, bm.Tags...)

	res, err := bm.Client.ImageBuild(mctx, reader, imageOps)
	if err != nil {
		return err
	}

	// defer res.Body.Close()

	if bm.Callback != nil {
		return bm.Callback(res)
	}

	return nil
}
