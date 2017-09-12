package docker

import (
	"github.com/influx6/box"
	"github.com/moby/moby/client"
)

// ImageTag returns a new ImageTagSpell instance to be executed on the client.
func (d *DockerCaster) ImageTag(tag string) (*ImageTagSpell, error) {
	var spell ImageTagSpell

	spell.tag = tag

	return &spell, nil
}

// ImageTagSpell defines a function type to modify internal fields of the ImageTagSpell.
type ImageTagOptions func(*ImageTagSpell)

// ImageTagResponseCallback defines a function type for ImageTagSpell response.
type ImageTagResponseCallback func() error

// AlwaysImageTagSpellWith returns a object that always executes the provided ImageTagSpell with the provided callback.
func AlwaysImageTagSpellWith(bm *ImageTagSpell, cb ImageTagResponseCallback) box.Spell {
	return &onceImageTagSpell{spell: bm, callback: cb}
}

type onceImageTagSpell struct {
	callback ImageTagResponseCallback
	spell    *ImageTagSpell
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceImageTagSpell) Exec(ctx box.CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// ImageTagSpell defines a structure which implements the Spell interface
// for executing of docker based commands for ImageTag.
type ImageTagSpell struct {
	client *client.Client

	tag string
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *ImageTagSpell) Exec(ctx box.CancelContext, callback ImageTagResponseCallback) error {
	// Execute client ImageTag method.
	err := cm.client.ImageTag(cm.tag)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback()
	}

	return nil
}
