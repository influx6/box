package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/influx6/box"
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

// ImageInspectWithRawSpell defines a structure which implements the Spell interface
// for executing of docker based commands for ImageInspectWithRaw.
type ImageInspectWithRawSpell struct {
	client *client.Client
}

// Spell returns a object implementing the box.Shell interface.
func (cm *ImageInspectWithRawSpell) Spell(callback ImageInspectWithRawResponseCallback) box.Spell {
	return &onceImageInspectWithRawSpell{spell: cm, callback: cb}
}

type onceImageInspectWithRawSpell struct {
	callback ImageInspectWithRawResponseCallback
	spell    *ImageInspectWithRawSpell
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceImageInspectWithRawSpell) Exec(ctx box.CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *ImageInspectWithRawSpell) Exec(ctx box.CancelContext, callback ImageInspectWithRawResponseCallback) error {
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
