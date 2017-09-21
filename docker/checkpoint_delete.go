package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/ops"
	"github.com/moby/moby/client"
)

// CheckpointDelete returns a new CheckpointDeleteOp instance to be executed on the client.
func (d *DockerCaster) CheckpointDelete(container string, chop types.CheckpointDeleteOptions) (*CheckpointDeleteOp, error) {
	var spell CheckpointDeleteOp

	spell.container = container

	spell.chop = chop

	return &spell, nil
}

// CheckpointDeleteOptions defines a function type to modify internal fields of the CheckpointDeleteOp.
type CheckpointDeleteOptions func(*CheckpointDeleteOp)

// CheckpointDeleteResponseCallback defines a function type for CheckpointDeleteOp response.
type CheckpointDeleteResponseCallback func() error

// CheckpointDeleteOp defines a structure which implements the Op interface
// for executing of docker based commands for CheckpointDelete.
type CheckpointDeleteOp struct {
	client *client.Client

	container string

	chop types.CheckpointDeleteOptions
}

// Op returns a object implementing the ops.Op interface.
func (cm *CheckpointDeleteOp) Op(callback CheckpointDeleteResponseCallback) ops.Op {
	return &onceCheckpointDeleteOp{spell: cm, callback: cb}
}

type onceCheckpointDeleteOp struct {
	callback CheckpointDeleteResponseCallback
	spell    *CheckpointDeleteOp
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceCheckpointDeleteOp) Exec(ctx context.CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *CheckpointDeleteOp) Exec(ctx context.CancelContext, callback CheckpointDeleteResponseCallback) error {
	// Execute client CheckpointDelete method.
	err := cm.client.CheckpointDelete(cm.container, cm.chop)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback()
	}

	return nil
}
