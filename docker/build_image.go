package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/influx6/moz/gen/filesystem"
	"github.com/moby/moby/client"
)

// BuildImageWith returns a new BuildImageSpell instance to be executed on the client.
func (d *DockerCaster) BuildImageWith(name string, fs io.Reader, ops ...BuildImageOptions) (*BuildImageSpell, error) {
	var spell BuildImageSpell
	spell.name = name
	spell.filesystem = fs
	spell.client = d.client

	for _, op := range ops {
		op(&spell)
	}

	if spell.filesystem == nil {
		return nil, ErrNoFilesystemProvided
	}

	return &spell, nil
}

// BuildImage returns a new BuildImageSpell instance to be executed on the client.
func (d *DockerCaster) BuildImage(name string, fs filesystem.Filesystem, ops ...BuildImageOptions) (*BuildImageSpell, error) {
	buildFSReader, err := fs.ToReader()
	if err != nil {
		return nil, err
	}

	var spell BuildImageSpell
	spell.name = name
	spell.filesystem = buildFSReader
	spell.client = d.client

	for _, op := range ops {
		op(&spell)
	}

	if spell.filesystem == nil {
		return nil, ErrNoFilesystemProvided
	}

	return &spell, nil
}

// BuildImageSpell defines a function type to modify internal fields of the BuildImageSpell.
type BuildImageOptions func(*BuildImageSpell)

// BuildImageResponseCallback defines a function type for BuildImageSpell.
type BuildImageResponseCallback func(types.ImageBuildResponse) error

// ImageSupplementaryTags sets the types.ImageBuildOptions for the BuildImageSpell.
func ImageSupplementaryTags(tags ...string) BuildImageOptions {
	return func(im *BuildImageSpell) {
		im.tags = tags
	}
}

// ImageBuildOptions sets the types.ImageBuildOptions for the BuildImageSpell.
func ImageBuildOptions(op types.ImageBuildOptions) BuildImageOptions {
	return func(im *BuildImageSpell) {
		im.imageOps = &op
	}
}

// AlwaysBuildImageSpellWith returns a object that always executes the provided BuildImageSpell with the provided callback.
func AlwaysBuildImageSpellWith(bm *BuildImageSpell, cb BuildImageResponseCallback) Spell {
	return &oncebuildImageSpell{spell: bm, callback: cb}
}

type oncebuildImageSpell struct {
	callback BuildImageResponseCallback
	spell    *BuildImageSpell
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *oncebuildImageSpell) Exec(ctx CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// BuildImageSpell defines a structure which implements the Spell interface
// for building of docker image.
// To use the BuildImageSpell the following is required:
// 1. A filesystem to be supplied has the build context
// 2. A name for the image to be build
// As optional, a types.ImageBuildOptions if you wish to override internally generated as
// a means of setting more fields for the docker clients request.
type BuildImageSpell struct {
	name       string
	dockerfile string
	tags       []string
	client     *client.Client
	filesystem io.Reader
	imageOps   *types.ImageBuildOptions
}

// Exec executes the image creation request through the provided docker client.
// If the spell has the types.ImageBuildOptions without a filesystem, then the
// types.ImageBuildOptions will be used as is, else the filesystem will be included
// has the underline BuildContext.
func (cm *BuildImageSpell) Exec(ctx CancelContext, callback BuildImageResponseCallback) error {
	if cm.client == nil {
		return ErrNoDockerClientProvided
	}

	done := make(chan struct{})
	defer close(done)

	// Cancel context if are done or if context has expired.
	reqCtx, cancel := context.WithCancel(context.Background())
	go func() {
		select {
		case <-ctx.Done():
			cancel()
			return
		case <-done:
			return
		}
	}()

	var imageOps types.ImageBuildOptions

	if cm.imageOps != nil {
		imageOps = *cm.imageOps
	} else {
		imageOps.PullParent = true
		imageOps.SuppressOutput = false
	}

	imageOps.Tags = append(imageOps.Tags, cm.name)
	imageOps.Tags = append(imageOps.Tags, cm.tags...)

	res, err := cm.client.ImageBuild(reqCtx, cm.filesystem, imageOps)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback(res)
	}

	defer res.Body.Close()

	return nil
}
