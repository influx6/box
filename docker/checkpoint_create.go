package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/metrics"
	"github.com/influx6/faux/ops"
	"github.com/moby/moby/client"
)

// CheckpointCreate returns a new CheckpointCreateOp instance to be executed on the client.
func (d *DockerCaster) CheckpointCreate(container string, chop types.CheckpointCreateOptions) (*CheckpointCreateOp, error) {
	var spell CheckpointCreateOp

	spell.container = container

	spell.chop = chop

	return &spell, nil
}

// CheckpointCreateOptions defines a function type to modify internal fields of the CheckpointCreateOp.
type CheckpointCreateOptions func(*CheckpointCreateOp)

// CheckpointCreateResponseCallback defines a function type for CheckpointCreateOp response.
type CheckpointCreateResponseCallback func() error

// CheckpointCreateOp defines a structure which implements the Op interface
// for executing of docker based commands for CheckpointCreate.
type CheckpointCreateOp struct {
	client *client.Client

	container string

	chop types.CheckpointCreateOptions
}

// Op returns a object implementing the ops.Op interface.
func (cm *CheckpointCreateOp) Op(callback CheckpointCreateResponseCallback) ops.Op {
	return &onceCheckpointCreateOp{spell: cm, callback: cb}
}

type onceCheckpointCreateOp struct {
	callback CheckpointCreateResponseCallback
	spell    *CheckpointCreateOp
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceCheckpointCreateOp) Exec(ctx context.CancelContext, m metrics.Metrics) error {
	return cm.spell.Exec(ctx, m, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *CheckpointCreateOp) Exec(ctx context.CancelContext, m metrics.Metrics, callback CheckpointCreateResponseCallback) error {
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

	// Execute client CheckpointCreate method.
	err := cm.client.CheckpointCreate(reqCtx, cm.container, cm.chop)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback()
	}

	return nil
}
