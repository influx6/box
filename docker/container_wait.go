package docker

import (
	"context"

	"github.com/influx6/faux/context"
	"github.com/influx6/faux/metrics"
	"github.com/influx6/faux/ops"
	"github.com/moby/moby/client"
	"github.com/moby/moby/container"
)

// ContainerWait returns a new ContainerWaitOp instance to be executed on the client.
func (d *DockerCaster) ContainerWait(containerID string, container container.WaitCondition) (*ContainerWaitOp, error) {
	var spell ContainerWaitOp

	spell.containerID = containerID

	spell.container = container

	return &spell, nil
}

// ContainerWaitOptions defines a function type to modify internal fields of the ContainerWaitOp.
type ContainerWaitOptions func(*ContainerWaitOp)

// ContainerWaitResponseCallback defines a function type for ContainerWaitOp response.
type ContainerWaitResponseCallback func() error

// ContainerWaitOp defines a structure which implements the Op interface
// for executing of docker based commands for ContainerWait.
type ContainerWaitOp struct {
	client *client.Client

	containerID string

	container container.WaitCondition
}

// Op returns a object implementing the ops.Op interface.
func (cm *ContainerWaitOp) Op(callback ContainerWaitResponseCallback) ops.Op {
	return &onceContainerWaitOp{spell: cm, callback: cb}
}

type onceContainerWaitOp struct {
	callback ContainerWaitResponseCallback
	spell    *ContainerWaitOp
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceContainerWaitOp) Exec(ctx context.CancelContext, m metrics.Metrics) error {
	return cm.spell.Exec(ctx, m, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *ContainerWaitOp) Exec(ctx context.CancelContext, m metrics.Metrics, callback ContainerWaitResponseCallback) error {
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

	// Execute client ContainerWait method.
	err := cm.client.ContainerWait(reqCtx, cm.containerID, cm.container)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback()
	}

	return nil
}
