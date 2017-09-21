package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/ops"
	"github.com/moby/moby/client"
)

// ImageInspectWithRaw returns a new ImageInspectWithRawOp instance to be executed on the client.
func (d *DockerCaster) ImageInspectWithRaw() (*ImageInspectWithRawOp, error) {
	var spell ImageInspectWithRawOp

	return &spell, nil
}

// ImageInspectWithRawOptions defines a function type to modify internal fields of the ImageInspectWithRawOp.
type ImageInspectWithRawOptions func(*ImageInspectWithRawOp)

// ImageInspectWithRawResponseCallback defines a function type for ImageInspectWithRawOp response.
type ImageInspectWithRawResponseCallback func(types.ImageInspect) error

// ImageInspectWithRawOp defines a structure which implements the Op interface
// for executing of docker based commands for ImageInspectWithRaw.
type ImageInspectWithRawOp struct {
	client *client.Client
}

// Op returns a object implementing the ops.Op interface.
func (cm *ImageInspectWithRawOp) Op(callback ImageInspectWithRawResponseCallback) ops.Op {
	return &onceImageInspectWithRawOp{spell: cm, callback: cb}
}

type onceImageInspectWithRawOp struct {
	callback ImageInspectWithRawResponseCallback
	spell    *ImageInspectWithRawOp
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceImageInspectWithRawOp) Exec(ctx context.CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *ImageInspectWithRawOp) Exec(ctx context.CancelContext, callback ImageInspectWithRawResponseCallback) error {
	// Execute client ImageInspectWithRaw method.
	ret0, err := cm.client.ImageInspectWithRaw()
	if err != nil {
		return err
	}

	if callback != nil {
		return callback(ret0)
	}

	return nil
}
