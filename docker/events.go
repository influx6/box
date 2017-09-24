package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/metrics"
	"github.com/influx6/faux/ops"
	"github.com/moby/moby/client"
)

// Events returns a new EventsOp instance to be executed on the client.
func (d *DockerCaster) Events(eventOp types.EventsOptions) (*EventsOp, error) {
	var spell EventsOp

	spell.eventOp = eventOp

	return &spell, nil
}

// EventsOptions defines a function type to modify internal fields of the EventsOp.
type EventsOptions func(*EventsOp)

// EventsResponseCallback defines a function type for EventsOp response.
type EventsResponseCallback func() error

// EventsOp defines a structure which implements the Op interface
// for executing of docker based commands for Events.
type EventsOp struct {
	client *client.Client

	eventOp types.EventsOptions
}

// Op returns a object implementing the ops.Op interface.
func (cm *EventsOp) Op(callback EventsResponseCallback) ops.Op {
	return &onceEventsOp{spell: cm, callback: cb}
}

type onceEventsOp struct {
	callback EventsResponseCallback
	spell    *EventsOp
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceEventsOp) Exec(ctx context.CancelContext, m metrics.Metrics) error {
	return cm.spell.Exec(ctx, m, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *EventsOp) Exec(ctx context.CancelContext, m metrics.Metrics, callback EventsResponseCallback) error {
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

	// Execute client Events method.
	err := cm.client.Events(reqCtx, cm.eventOp)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback()
	}

	return nil
}
