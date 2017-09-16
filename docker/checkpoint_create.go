package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/influx6/box"
	"github.com/moby/moby/client"
)

// CheckpointCreate returns a new CheckpointCreateSpell instance to be executed on the client.
func (d *DockerCaster) CheckpointCreate(container string, chop types.CheckpointCreateOptions) (*CheckpointCreateSpell, error) {
	var spell CheckpointCreateSpell

	spell.container = container

	spell.chop = chop

	return &spell, nil
}

// CheckpointCreateSpell defines a function type to modify internal fields of the CheckpointCreateSpell.
type CheckpointCreateOptions func(*CheckpointCreateSpell)

// CheckpointCreateResponseCallback defines a function type for CheckpointCreateSpell response.
type CheckpointCreateResponseCallback func() error

// CheckpointCreateSpell defines a structure which implements the Spell interface
// for executing of docker based commands for CheckpointCreate.
type CheckpointCreateSpell struct {
	client *client.Client

	container string

	chop types.CheckpointCreateOptions
}

// Spell returns a object implementing the box.Shell interface.
func (cm *CheckpointCreateSpell) Spell(callback CheckpointCreateResponseCallback) box.Spell {
	return &onceCheckpointCreateSpell{spell: cm, callback: cb}
}

type onceCheckpointCreateSpell struct {
	callback CheckpointCreateResponseCallback
	spell    *CheckpointCreateSpell
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceCheckpointCreateSpell) Exec(ctx box.CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *CheckpointCreateSpell) Exec(ctx box.CancelContext, callback CheckpointCreateResponseCallback) error {
	// Execute client CheckpointCreate method.
	err := cm.client.CheckpointCreate(cm.container, cm.chop)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback()
	}

	return nil
}
