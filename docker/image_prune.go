package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/metrics"
	"github.com/influx6/faux/ops"
	"github.com/moby/moby/api/types/filters"
	"github.com/moby/moby/client"
)

// ImagePrune returns a new ImagePruneOp instance to be executed on the client.
func (d *DockerCaster) ImagePrune(args filters.Args) (*ImagePruneOp, error) {
	var spell ImagePruneOp

	spell.args = args

	return &spell, nil
}

// ImagePruneOptions defines a function type to modify internal fields of the ImagePruneOp.
type ImagePruneOptions func(*ImagePruneOp)

// ImagePruneResponseCallback defines a function type for ImagePruneOp response.
type ImagePruneResponseCallback func(types.ImagesPruneReport) error

// ImagePruneOp defines a structure which implements the Op interface
// for executing of docker based commands for ImagePrune.
type ImagePruneOp struct {
	client *client.Client

	args filters.Args
}

// Op returns a object implementing the ops.Op interface.
func (cm *ImagePruneOp) Op(callback ImagePruneResponseCallback) ops.Op {
	return &onceImagePruneOp{spell: cm, callback: cb}
}

type onceImagePruneOp struct {
	callback ImagePruneResponseCallback
	spell    *ImagePruneOp
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceImagePruneOp) Exec(ctx context.CancelContext, m metrics.Metrics) error {
	return cm.spell.Exec(ctx, m, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *ImagePruneOp) Exec(ctx context.CancelContext, m metrics.Metrics, callback ImagePruneResponseCallback) error {
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

	// Execute client ImagePrune method.
	ret0, err := cm.client.ImagePrune(cm.args)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback(ret0)
	}

	return nil
}
