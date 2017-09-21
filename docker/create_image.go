package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/ops"
	"github.com/moby/moby/client"
)

// CreateImage returns a new CreateImageOp instance to be executed on the client.
func (d *DockerCaster) CreateImage(reader io.ReadCloser) (*CreateImageOp, error) {
	var spell CreateImageOp

	spell.reader = reader

	return &spell, nil
}

// CreateImageOptions defines a function type to modify internal fields of the CreateImageOp.
type CreateImageOptions func(*CreateImageOp)

// CreateImageResponseCallback defines a function type for CreateImageOp response.
type CreateImageResponseCallback func(types.ImageLoadResponse) error

// CreateImageOp defines a structure which implements the Op interface
// for executing of docker based commands for CreateImage.
type CreateImageOp struct {
	client *client.Client

	reader io.ReadCloser
}

// Op returns a object implementing the ops.Op interface.
func (cm *CreateImageOp) Op(callback CreateImageResponseCallback) ops.Op {
	return &onceCreateImageOp{spell: cm, callback: cb}
}

type onceCreateImageOp struct {
	callback CreateImageResponseCallback
	spell    *CreateImageOp
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceCreateImageOp) Exec(ctx context.CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *CreateImageOp) Exec(ctx context.CancelContext, callback CreateImageResponseCallback) error {
	// Execute client CreateImage method.
	ret0, err := cm.client.CreateImage(cm.reader)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback(ret0)
	}

	return nil
}
