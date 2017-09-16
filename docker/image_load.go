package docker

import (
	"io"

	"github.com/docker/docker/api/types"
	"github.com/influx6/box"
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

// ImageLoadSpell defines a structure which implements the Spell interface
// for executing of docker based commands for ImageLoad.
type ImageLoadSpell struct {
	client *client.Client

	reader io.Reader
}

// Spell returns a object implementing the box.Shell interface.
func (cm *ImageLoadSpell) Spell(callback ImageLoadResponseCallback) box.Spell {
	return &onceImageLoadSpell{spell: cm, callback: cb}
}

type onceImageLoadSpell struct {
	callback ImageLoadResponseCallback
	spell    *ImageLoadSpell
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceImageLoadSpell) Exec(ctx box.CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *ImageLoadSpell) Exec(ctx box.CancelContext, callback ImageLoadResponseCallback) error {
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
