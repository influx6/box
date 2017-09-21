package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/influx6/faux/context"
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
func (cm *onceEventsOp) Exec(ctx context.CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *EventsOp) Exec(ctx context.CancelContext, callback EventsResponseCallback) error {
	// Execute client Events method.
	err := cm.client.Events(cm.eventOp)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback()
	}

	return nil
}
