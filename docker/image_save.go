package docker

import (
	"context"
	"io"

	"github.com/influx6/faux/context"
	"github.com/influx6/faux/ops"
	"github.com/moby/moby/client"
)

// ImageSave returns a new ImageSaveOp instance to be executed on the client.
func (d *DockerCaster) ImageSave(ops []string) (*ImageSaveOp, error) {
	var spell ImageSaveOp

	spell.ops = ops

	return &spell, nil
}

// ImageSaveOptions defines a function type to modify internal fields of the ImageSaveOp.
type ImageSaveOptions func(*ImageSaveOp)

// ImageSaveResponseCallback defines a function type for ImageSaveOp response.
type ImageSaveResponseCallback func(io.ReadCloser) error

// ImageSaveOp defines a structure which implements the Op interface
// for executing of docker based commands for ImageSave.
type ImageSaveOp struct {
	client *client.Client

	ops []string
}

// Op returns a object implementing the ops.Op interface.
func (cm *ImageSaveOp) Op(callback ImageSaveResponseCallback) ops.Op {
	return &onceImageSaveOp{spell: cm, callback: cb}
}

type onceImageSaveOp struct {
	callback ImageSaveResponseCallback
	spell    *ImageSaveOp
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceImageSaveOp) Exec(ctx context.CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *ImageSaveOp) Exec(ctx context.CancelContext, callback ImageSaveResponseCallback) error {
	// Execute client ImageSave method.
	ret0, err := cm.client.ImageSave(cm.ops)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback(ret0)
	}

	return nil
}
