package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/ops"
	"github.com/moby/moby/client"
)

// ImagePull returns a new ImagePullOp instance to be executed on the client.
func (d *DockerCaster) ImagePull(imgOp types.ImagePullOptions) (*ImagePullOp, error) {
	var spell ImagePullOp

	spell.imgOp = imgOp

	return &spell, nil
}

// ImagePullOptions defines a function type to modify internal fields of the ImagePullOp.
type ImagePullOptions func(*ImagePullOp)

// ImagePullResponseCallback defines a function type for ImagePullOp response.
type ImagePullResponseCallback func(io.ReadCloser) error

// ImagePullOp defines a structure which implements the Op interface
// for executing of docker based commands for ImagePull.
type ImagePullOp struct {
	client *client.Client

	imgOp types.ImagePullOptions
}

// Op returns a object implementing the ops.Op interface.
func (cm *ImagePullOp) Op(callback ImagePullResponseCallback) ops.Op {
	return &onceImagePullOp{spell: cm, callback: cb}
}

type onceImagePullOp struct {
	callback ImagePullResponseCallback
	spell    *ImagePullOp
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceImagePullOp) Exec(ctx context.CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *ImagePullOp) Exec(ctx context.CancelContext, callback ImagePullResponseCallback) error {
	// Execute client ImagePull method.
	ret0, err := cm.client.ImagePull(cm.imgOp)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback(ret0)
	}

	return nil
}
