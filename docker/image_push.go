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

// ImagePush returns a new ImagePushOp instance to be executed on the client.
func (d *DockerCaster) ImagePush(imp types.ImagePushOptions) (*ImagePushOp, error) {
	var spell ImagePushOp

	spell.imp = imp

	return &spell, nil
}

// ImagePushOptions defines a function type to modify internal fields of the ImagePushOp.
type ImagePushOptions func(*ImagePushOp)

// ImagePushResponseCallback defines a function type for ImagePushOp response.
type ImagePushResponseCallback func(io.ReadCloser) error

// ImagePushOp defines a structure which implements the Op interface
// for executing of docker based commands for ImagePush.
type ImagePushOp struct {
	client *client.Client

	imp types.ImagePushOptions
}

// Op returns a object implementing the ops.Op interface.
func (cm *ImagePushOp) Op(callback ImagePushResponseCallback) ops.Op {
	return &onceImagePushOp{spell: cm, callback: cb}
}

type onceImagePushOp struct {
	callback ImagePushResponseCallback
	spell    *ImagePushOp
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceImagePushOp) Exec(ctx context.CancelContext, m metrics.Metrics) error {
	return cm.spell.Exec(ctx, m, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *ImagePushOp) Exec(ctx context.CancelContext, m metrics.Metrics, callback ImagePushResponseCallback) error {
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

	// Execute client ImagePush method.
	ret0, err := cm.client.ImagePush(cm.imp)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback(ret0)
	}

	return nil
}
