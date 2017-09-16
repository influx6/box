package docker

import (
	"io"

	"github.com/docker/docker/api/types"
	"github.com/influx6/box"
	"github.com/moby/moby/client"
)

// CopyToContainer returns a new CopyToContainerSpell instance to be executed on the client.
func (d *DockerCaster) CopyToContainer(container string, topath string, reader io.ReadCloser, cops types.CopyToContainerOptions) (*CopyToContainerSpell, error) {
	var spell CopyToContainerSpell

	spell.container = container

	spell.topath = topath

	spell.reader = reader

	spell.cops = cops

	return &spell, nil
}

// CopyToContainerSpell defines a function type to modify internal fields of the CopyToContainerSpell.
type CopyToContainerOptions func(*CopyToContainerSpell)

// CopyToContainerResponseCallback defines a function type for CopyToContainerSpell response.
type CopyToContainerResponseCallback func() error

// CopyToContainerSpell defines a structure which implements the Spell interface
// for executing of docker based commands for CopyToContainer.
type CopyToContainerSpell struct {
	client *client.Client

	container string

	topath string

	reader io.ReadCloser

	cops types.CopyToContainerOptions
}

// Spell returns a object implementing the box.Shell interface.
func (cm *CopyToContainerSpell) Spell(callback CopyToContainerResponseCallback) box.Spell {
	return &onceCopyToContainerSpell{spell: cm, callback: cb}
}

type onceCopyToContainerSpell struct {
	callback CopyToContainerResponseCallback
	spell    *CopyToContainerSpell
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceCopyToContainerSpell) Exec(ctx box.CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *CopyToContainerSpell) Exec(ctx box.CancelContext, callback CopyToContainerResponseCallback) error {
	// Execute client CopyToContainer method.
	err := cm.client.CopyToContainer(cm.container, cm.topath, cm.reader, cm.cops)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback()
	}

	return nil
}
