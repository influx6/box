package docker

import (
	"io"

	"github.com/docker/docker/api/types"
	"github.com/influx6/box"
	"github.com/moby/moby/client"
)

// CopyFromContainer returns a new CopyFromContainerSpell instance to be executed on the client.
func (d *DockerCaster) CopyFromContainer(container string, srcPath string) (*CopyFromContainerSpell, error) {
	var spell CopyFromContainerSpell

	spell.container = container

	spell.srcPath = srcPath

	return &spell, nil
}

// CopyFromContainerSpell defines a function type to modify internal fields of the CopyFromContainerSpell.
type CopyFromContainerOptions func(*CopyFromContainerSpell)

// CopyFromContainerResponseCallback defines a function type for CopyFromContainerSpell response.
type CopyFromContainerResponseCallback func(io.ReadCloser, types.ContainerPathStat) error

// AlwaysCopyFromContainerSpellWith returns a object that always executes the provided CopyFromContainerSpell with the provided callback.
func AlwaysCopyFromContainerSpellWith(bm *CopyFromContainerSpell, cb CopyFromContainerResponseCallback) box.Spell {
	return &onceCopyFromContainerSpell{spell: bm, callback: cb}
}

type onceCopyFromContainerSpell struct {
	callback CopyFromContainerResponseCallback
	spell    *CopyFromContainerSpell
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceCopyFromContainerSpell) Exec(ctx box.CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// CopyFromContainerSpell defines a structure which implements the Spell interface
// for executing of docker based commands for CopyFromContainer.
type CopyFromContainerSpell struct {
	client *client.Client

	container string

	srcPath string
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *CopyFromContainerSpell) Exec(ctx box.CancelContext, callback CopyFromContainerResponseCallback) error {
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
