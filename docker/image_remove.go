package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/metrics"
	"github.com/influx6/faux/ops"
	"github.com/moby/moby/client"
)

// ImageRemove returns a new ImageRemoveOp instance to be executed on the client.
func (d *DockerCaster) ImageRemove(removeOps types.ImageRemoveOptions) (*ImageRemoveOp, error) {
	var spell ImageRemoveOp

	spell.removeOps = removeOps

	return &spell, nil
}

// ImageRemoveOptions defines a function type to modify internal fields of the ImageRemoveOp.
type ImageRemoveOptions func(*ImageRemoveOp)

// ImageRemoveResponseCallback defines a function type for ImageRemoveOp response.
type ImageRemoveResponseCallback func([]types.ImageDeleteResponseItem) error

// ImageRemoveOp defines a structure which implements the Op interface
// for executing of docker based commands for ImageRemove.
type ImageRemoveOp struct {
	client *client.Client

	removeOps types.ImageRemoveOptions
}

// Op returns a object implementing the ops.Op interface.
func (cm *ImageRemoveOp) Op(callback ImageRemoveResponseCallback) ops.Op {
	return &onceImageRemoveOp{spell: cm, callback: cb}
}

type onceImageRemoveOp struct {
	callback ImageRemoveResponseCallback
	spell    *ImageRemoveOp
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceImageRemoveOp) Exec(ctx context.CancelContext, m metrics.Metrics) error {
	return cm.spell.Exec(ctx, m, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *ImageRemoveOp) Exec(ctx context.CancelContext, m metrics.Metrics, callback ImageRemoveResponseCallback) error {
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

	// Execute client ImageRemove method.
	ret0, err := cm.client.ImageRemove(cm.removeOps)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback(ret0)
	}

	return nil
}
