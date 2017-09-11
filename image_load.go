package dockish

import (
	"io"

	"github.com/docker/docker/api/types"
	"github.com/moby/moby/client"
)

// ImageLoad returns a new ImageLoadSpell instance to be executed on the client.
func (d *DockerCaster) ImageLoad(reader io.Reader) (*ImageLoadSpell, error) {
	var spell ImageLoadSpell

	spell.reader = reader

	return &spell, nil
}

// ImageLoadSpell defines a function type to modify internal fields of the ImageLoadSpell.
type ImageLoadOptions func(*ImageLoadSpell)

// ImageLoadResponseCallback defines a function type for ImageLoadSpell response.
type ImageLoadResponseCallback func(types.ImageLoadResponse) error

// AlwaysImageLoadSpellWith returns a object that always executes the provided ImageLoadSpell with the provided callback.
func AlwaysImageLoadSpellWith(bm *ImageLoadSpell, cb ImageLoadResponseCallback) Spell {
	return &onceImageLoadSpell{spell: bm, callback: cb}
}

type onceImageLoadSpell struct {
	callback ImageLoadResponseCallback
	spell    *ImageLoadSpell
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceImageLoadSpell) Exec(ctx CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// ImageLoadSpell defines a structure which implements the Spell interface
// for executing of docker based commands for ImageLoad.
type ImageLoadSpell struct {
	client *client.Client

	reader io.Reader
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *ImageLoadSpell) Exec(ctx CancelContext, callback ImageLoadResponseCallback) error {
	// Execute client ImageLoad method.
	ret0, err := cm.client.ImageLoad(cm.reader)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback(ret0)
	}

	return nil
}
