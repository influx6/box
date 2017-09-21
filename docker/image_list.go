package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/ops"
	"github.com/moby/moby/client"
)

// ImageList returns a new ImageListOp instance to be executed on the client.
func (d *DockerCaster) ImageList(listOps types.ImageListOptions) (*ImageListOp, error) {
	var spell ImageListOp

	spell.listOps = listOps

	return &spell, nil
}

// ImageListOptions defines a function type to modify internal fields of the ImageListOp.
type ImageListOptions func(*ImageListOp)

// ImageListResponseCallback defines a function type for ImageListOp response.
type ImageListResponseCallback func([]types.ImageSummary) error

// ImageListOp defines a structure which implements the Op interface
// for executing of docker based commands for ImageList.
type ImageListOp struct {
	client *client.Client

	listOps types.ImageListOptions
}

// Op returns a object implementing the ops.Op interface.
func (cm *ImageListOp) Op(callback ImageListResponseCallback) ops.Op {
	return &onceImageListOp{spell: cm, callback: cb}
}

type onceImageListOp struct {
	callback ImageListResponseCallback
	spell    *ImageListOp
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceImageListOp) Exec(ctx context.CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *ImageListOp) Exec(ctx context.CancelContext, callback ImageListResponseCallback) error {
	// Execute client ImageList method.
	ret0, err := cm.client.ImageList(cm.listOps)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback(ret0)
	}

	return nil
}
