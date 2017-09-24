package docker

import (
	"context"

	"github.com/docker/docker/api/types/image"
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/metrics"
	"github.com/influx6/faux/ops"
	"github.com/moby/moby/client"
)

// ImageHistory returns a new ImageHistoryOp instance to be executed on the client.
func (d *DockerCaster) ImageHistory() (*ImageHistoryOp, error) {
	var spell ImageHistoryOp

	return &spell, nil
}

// ImageHistoryOptions defines a function type to modify internal fields of the ImageHistoryOp.
type ImageHistoryOptions func(*ImageHistoryOp)

// ImageHistoryResponseCallback defines a function type for ImageHistoryOp response.
type ImageHistoryResponseCallback func(image.HistoryResponseItem) error

// ImageHistoryOp defines a structure which implements the Op interface
// for executing of docker based commands for ImageHistory.
type ImageHistoryOp struct {
	client *client.Client
}

// Op returns a object implementing the ops.Op interface.
func (cm *ImageHistoryOp) Op(callback ImageHistoryResponseCallback) ops.Op {
	return &onceImageHistoryOp{spell: cm, callback: cb}
}

type onceImageHistoryOp struct {
	callback ImageHistoryResponseCallback
	spell    *ImageHistoryOp
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceImageHistoryOp) Exec(ctx context.CancelContext, m metrics.Metrics) error {
	return cm.spell.Exec(ctx, m, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *ImageHistoryOp) Exec(ctx context.CancelContext, m metrics.Metrics, callback ImageHistoryResponseCallback) error {
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

	// Execute client ImageHistory method.
	ret0, err := cm.client.ImageHistory()
	if err != nil {
		return err
	}

	if callback != nil {
		return callback(ret0)
	}

	return nil
}
