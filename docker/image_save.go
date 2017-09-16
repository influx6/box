package docker

import (
	"io"

	"github.com/influx6/box"
	"github.com/moby/moby/client"
)

// ImageSave returns a new ImageSaveSpell instance to be executed on the client.
func (d *DockerCaster) ImageSave(ops []string) (*ImageSaveSpell, error) {
	var spell ImageSaveSpell

	spell.ops = ops

	return &spell, nil
}

// ImageSaveSpell defines a function type to modify internal fields of the ImageSaveSpell.
type ImageSaveOptions func(*ImageSaveSpell)

// ImageSaveResponseCallback defines a function type for ImageSaveSpell response.
type ImageSaveResponseCallback func(io.ReadCloser) error

// ImageSaveSpell defines a structure which implements the Spell interface
// for executing of docker based commands for ImageSave.
type ImageSaveSpell struct {
	client *client.Client

	ops []string
}

// Spell returns a object implementing the box.Shell interface.
func (cm *ImageSaveSpell) Spell(callback ImageSaveResponseCallback) box.Spell {
	return &onceImageSaveSpell{spell: cm, callback: cb}
}

type onceImageSaveSpell struct {
	callback ImageSaveResponseCallback
	spell    *ImageSaveSpell
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceImageSaveSpell) Exec(ctx box.CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *ImageSaveSpell) Exec(ctx box.CancelContext, callback ImageSaveResponseCallback) error {
	// Execute client ImageSave method.
	ret0, err := cm.client.ImageSave(cm.ops)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback(ret0)
	}

	return nil
}
