package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/metrics"
	"github.com/influx6/faux/ops"
	"github.com/moby/moby/client"
)

// ImageLoad returns a new ImageLoadOp instance to be executed on the client.
func (d *DockerCaster) ImageLoad(reader io.Reader) (*ImageLoadOp, error) {
	var spell ImageLoadOp

	spell.reader = reader

	return &spell, nil
}

// ImageLoadOptions defines a function type to modify internal fields of the ImageLoadOp.
type ImageLoadOptions func(*ImageLoadOp)

// ImageLoadResponseCallback defines a function type for ImageLoadOp response.
type ImageLoadResponseCallback func(types.ImageLoadResponse) error

// ImageLoadOp defines a structure which implements the Op interface
// for executing of docker based commands for ImageLoad.
type ImageLoadOp struct {
	client *client.Client

	reader io.Reader
}

// Op returns a object implementing the ops.Op interface.
func (cm *ImageLoadOp) Op(callback ImageLoadResponseCallback) ops.Op {
	return &onceImageLoadOp{spell: cm, callback: cb}
}

type onceImageLoadOp struct {
	callback ImageLoadResponseCallback
	spell    *ImageLoadOp
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceImageLoadOp) Exec(ctx context.CancelContext, m metrics.Metrics) error {
	return cm.spell.Exec(ctx, m, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *ImageLoadOp) Exec(ctx context.CancelContext, m metrics.Metrics, callback ImageLoadResponseCallback) error {
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

	// Execute client ImageLoad method.
	ret0, err := cm.client.ImageLoad(cm.reader)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback(ret0)
	}

	return nil
}
