package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/influx6/box"
	"github.com/moby/moby/client"
)

// CheckpointDelete returns a new CheckpointDeleteSpell instance to be executed on the client.
func (d *DockerCaster) CheckpointDelete(container string, chop types.CheckpointDeleteOptions) (*CheckpointDeleteSpell, error) {
	var spell CheckpointDeleteSpell

	spell.container = container

	spell.chop = chop

	return &spell, nil
}

// CheckpointDeleteSpell defines a function type to modify internal fields of the CheckpointDeleteSpell.
type CheckpointDeleteOptions func(*CheckpointDeleteSpell)

// CheckpointDeleteResponseCallback defines a function type for CheckpointDeleteSpell response.
type CheckpointDeleteResponseCallback func() error

// CheckpointDeleteSpell defines a structure which implements the Spell interface
// for executing of docker based commands for CheckpointDelete.
type CheckpointDeleteSpell struct {
	client *client.Client

	container string

	chop types.CheckpointDeleteOptions
}

// Spell returns a object implementing the box.Shell interface.
func (cm *CheckpointDeleteSpell) Spell(callback CheckpointDeleteResponseCallback) box.Spell {
	return &onceCheckpointDeleteSpell{spell: cm, callback: cb}
}

type onceCheckpointDeleteSpell struct {
	callback CheckpointDeleteResponseCallback
	spell    *CheckpointDeleteSpell
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceCheckpointDeleteSpell) Exec(ctx box.CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *CheckpointDeleteSpell) Exec(ctx box.CancelContext, callback CheckpointDeleteResponseCallback) error {
	// Execute client CheckpointDelete method.
	err := cm.client.CheckpointDelete(cm.container, cm.chop)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback()
	}

	return nil
}
