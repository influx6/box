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

// ImageImport returns a new ImageImportOp instance to be executed on the client.
func (d *DockerCaster) ImageImport(impOp types.ImageImportOptions) (*ImageImportOp, error) {
	var spell ImageImportOp

	spell.impOp = impOp

	return &spell, nil
}

// ImageImportOptions defines a function type to modify internal fields of the ImageImportOp.
type ImageImportOptions func(*ImageImportOp)

// ImageImportResponseCallback defines a function type for ImageImportOp response.
type ImageImportResponseCallback func(io.ReadCloser) error

// ImageImportOp defines a structure which implements the Op interface
// for executing of docker based commands for ImageImport.
type ImageImportOp struct {
	client *client.Client

	impOp types.ImageImportOptions
}

// Op returns a object implementing the ops.Op interface.
func (cm *ImageImportOp) Op(callback ImageImportResponseCallback) ops.Op {
	return &onceImageImportOp{spell: cm, callback: cb}
}

type onceImageImportOp struct {
	callback ImageImportResponseCallback
	spell    *ImageImportOp
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceImageImportOp) Exec(ctx context.CancelContext, m metrics.Metrics) error {
	return cm.spell.Exec(ctx, m, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *ImageImportOp) Exec(ctx context.CancelContext, m metrics.Metrics, callback ImageImportResponseCallback) error {
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

	// Execute client ImageImport method.
	ret0, err := cm.client.ImageImport(cm.impOp)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback(ret0)
	}

	return nil
}
