package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/ops"
	"github.com/moby/moby/client"
)

// CopyToContainer returns a new CopyToContainerOp instance to be executed on the client.
func (d *DockerCaster) CopyToContainer(container string, topath string, reader io.ReadCloser, cops types.CopyToContainerOptions) (*CopyToContainerOp, error) {
	var spell CopyToContainerOp

	spell.container = container

	spell.topath = topath

	spell.reader = reader

	spell.cops = cops

	return &spell, nil
}

// CopyToContainerOptions defines a function type to modify internal fields of the CopyToContainerOp.
type CopyToContainerOptions func(*CopyToContainerOp)

// CopyToContainerResponseCallback defines a function type for CopyToContainerOp response.
type CopyToContainerResponseCallback func() error

// CopyToContainerOp defines a structure which implements the Op interface
// for executing of docker based commands for CopyToContainer.
type CopyToContainerOp struct {
	client *client.Client

	container string

	topath string

	reader io.ReadCloser

	cops types.CopyToContainerOptions
}

// Op returns a object implementing the ops.Op interface.
func (cm *CopyToContainerOp) Op(callback CopyToContainerResponseCallback) ops.Op {
	return &onceCopyToContainerOp{spell: cm, callback: cb}
}

type onceCopyToContainerOp struct {
	callback CopyToContainerResponseCallback
	spell    *CopyToContainerOp
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceCopyToContainerOp) Exec(ctx context.CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *CopyToContainerOp) Exec(ctx context.CancelContext, callback CopyToContainerResponseCallback) error {
	// Execute client CopyToContainer method.
	err := cm.client.CopyToContainer(cm.container, cm.topath, cm.reader, cm.cops)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback()
	}

	return nil
}
