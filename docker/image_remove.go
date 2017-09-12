package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/influx6/box"
	"github.com/moby/moby/client"
)

// ImageRemove returns a new ImageRemoveSpell instance to be executed on the client.
func (d *DockerCaster) ImageRemove(removeOps types.ImageRemoveOptions) (*ImageRemoveSpell, error) {
	var spell ImageRemoveSpell

	spell.removeOps = removeOps

	return &spell, nil
}

// ImageRemoveSpell defines a function type to modify internal fields of the ImageRemoveSpell.
type ImageRemoveOptions func(*ImageRemoveSpell)

// ImageRemoveResponseCallback defines a function type for ImageRemoveSpell response.
type ImageRemoveResponseCallback func([]types.ImageDeleteResponseItem) error

// AlwaysImageRemoveSpellWith returns a object that always executes the provided ImageRemoveSpell with the provided callback.
func AlwaysImageRemoveSpellWith(bm *ImageRemoveSpell, cb ImageRemoveResponseCallback) box.Spell {
	return &onceImageRemoveSpell{spell: bm, callback: cb}
}

type onceImageRemoveSpell struct {
	callback ImageRemoveResponseCallback
	spell    *ImageRemoveSpell
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceImageRemoveSpell) Exec(ctx box.CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// ImageRemoveSpell defines a structure which implements the Spell interface
// for executing of docker based commands for ImageRemove.
type ImageRemoveSpell struct {
	client *client.Client

	removeOps types.ImageRemoveOptions
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *ImageRemoveSpell) Exec(ctx box.CancelContext, callback ImageRemoveResponseCallback) error {
	// Execute client ImageRemove method.
	ret0, err := cm.client.ImageRemove(cm.removeOps)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback(ret0)
	}

	return nil
}
