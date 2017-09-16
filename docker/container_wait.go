package docker

import (
	"github.com/influx6/box"
	"github.com/moby/moby/client"
	"github.com/moby/moby/container"
)

// ContainerWait returns a new ContainerWaitSpell instance to be executed on the client.
func (d *DockerCaster) ContainerWait(containerID string, container container.WaitCondition) (*ContainerWaitSpell, error) {
	var spell ContainerWaitSpell

	spell.containerID = containerID

	spell.container = container

	return &spell, nil
}

// ContainerWaitSpell defines a function type to modify internal fields of the ContainerWaitSpell.
type ContainerWaitOptions func(*ContainerWaitSpell)

// ContainerWaitResponseCallback defines a function type for ContainerWaitSpell response.
type ContainerWaitResponseCallback func() error

// ContainerWaitSpell defines a structure which implements the Spell interface
// for executing of docker based commands for ContainerWait.
type ContainerWaitSpell struct {
	client *client.Client

	containerID string

	container container.WaitCondition
}

// Spell returns a object implementing the box.Shell interface.
func (cm *ContainerWaitSpell) Spell(callback ContainerWaitResponseCallback) box.Spell {
	return &onceContainerWaitSpell{spell: cm, callback: cb}
}

type onceContainerWaitSpell struct {
	callback ContainerWaitResponseCallback
	spell    *ContainerWaitSpell
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceContainerWaitSpell) Exec(ctx box.CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *ContainerWaitSpell) Exec(ctx box.CancelContext, callback ContainerWaitResponseCallback) error {
	// Execute client ContainerWait method.
	err := cm.client.ContainerWait(cm.containerID, cm.container)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback()
	}

	return nil
}
