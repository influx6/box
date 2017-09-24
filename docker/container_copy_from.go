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

// CopyFromContainer returns a new CopyFromContainerOp instance to be executed on the client.
func (d *DockerCaster) CopyFromContainer(container string, srcPath string) (*CopyFromContainerOp, error) {
	var spell CopyFromContainerOp

	spell.container = container

	spell.srcPath = srcPath

	return &spell, nil
}

// CopyFromContainerOptions defines a function type to modify internal fields of the CopyFromContainerOp.
type CopyFromContainerOptions func(*CopyFromContainerOp)

// CopyFromContainerResponseCallback defines a function type for CopyFromContainerOp response.
type CopyFromContainerResponseCallback func(io.ReadCloser, types.ContainerPathStat) error

// CopyFromContainerOp defines a structure which implements the Op interface
// for executing of docker based commands for CopyFromContainer.
type CopyFromContainerOp struct {
	client *client.Client

	container string

	srcPath string
}

// Op returns a object implementing the ops.Op interface.
func (cm *CopyFromContainerOp) Op(callback CopyFromContainerResponseCallback) ops.Op {
	return &onceCopyFromContainerOp{spell: cm, callback: cb}
}

type onceCopyFromContainerOp struct {
	callback CopyFromContainerResponseCallback
	spell    *CopyFromContainerOp
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceCopyFromContainerOp) Exec(ctx context.CancelContext, m metrics.Metrics) error {
	return cm.spell.Exec(ctx, m, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *CopyFromContainerOp) Exec(ctx context.CancelContext, m metrics.Metrics, callback CopyFromContainerResponseCallback) error {
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

	// Execute client CopyFromContainer method.
	ret0, ret1, err := cm.client.CopyFromContainer(cm.container, cm.srcPath)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback(ret0, ret1)
	}

	return nil
}
