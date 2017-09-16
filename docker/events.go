package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/influx6/box"
	"github.com/moby/moby/client"
)

// Events returns a new EventsSpell instance to be executed on the client.
func (d *DockerCaster) Events(eventOp types.EventsOptions) (*EventsSpell, error) {
	var spell EventsSpell

	spell.eventOp = eventOp

	return &spell, nil
}

// EventsSpell defines a function type to modify internal fields of the EventsSpell.
type EventsOptions func(*EventsSpell)

// EventsResponseCallback defines a function type for EventsSpell response.
type EventsResponseCallback func() error

// EventsSpell defines a structure which implements the Spell interface
// for executing of docker based commands for Events.
type EventsSpell struct {
	client *client.Client

	eventOp types.EventsOptions
}

// Spell returns a object implementing the box.Shell interface.
func (cm *EventsSpell) Spell(callback EventsResponseCallback) box.Spell {
	return &onceEventsSpell{spell: cm, callback: cb}
}

type onceEventsSpell struct {
	callback EventsResponseCallback
	spell    *EventsSpell
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceEventsSpell) Exec(ctx box.CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *EventsSpell) Exec(ctx box.CancelContext, callback EventsResponseCallback) error {
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
