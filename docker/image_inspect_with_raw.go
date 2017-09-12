package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/moby/moby/client"
)

// ImageInspectWithRaw returns a new ImageInspectWithRawSpell instance to be executed on the client.
func (d *DockerCaster) ImageInspectWithRaw() (*ImageInspectWithRawSpell, error) {
	var spell ImageInspectWithRawSpell

	return &spell, nil
}

// ImageInspectWithRawSpell defines a function type to modify internal fields of the ImageInspectWithRawSpell.
type ImageInspectWithRawOptions func(*ImageInspectWithRawSpell)

// ImageInspectWithRawResponseCallback defines a function type for ImageInspectWithRawSpell response.
type ImageInspectWithRawResponseCallback func(types.ImageInspect) error

// AlwaysImageInspectWithRawSpellWith returns a object that always executes the provided ImageInspectWithRawSpell with the provided callback.
func AlwaysImageInspectWithRawSpellWith(bm *ImageInspectWithRawSpell, cb ImageInspectWithRawResponseCallback) Spell {
	return &onceImageInspectWithRawSpell{spell: bm, callback: cb}
}

type onceImageInspectWithRawSpell struct {
	callback ImageInspectWithRawResponseCallback
	spell    *ImageInspectWithRawSpell
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceImageInspectWithRawSpell) Exec(ctx CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// ImageInspectWithRawSpell defines a structure which implements the Spell interface
// for executing of docker based commands for ImageInspectWithRaw.
type ImageInspectWithRawSpell struct {
	client *client.Client
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *ImageInspectWithRawSpell) Exec(ctx CancelContext, callback ImageInspectWithRawResponseCallback) error {
	// Execute client ImageInspectWithRaw method.
	ret0, err := cm.client.ImageInspectWithRaw()
	if err != nil {
		return err
	}

	if callback != nil {
		return callback(ret0)
	}

	return nil
}
