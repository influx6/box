package docker

import (
	"context"

	"github.com/influx6/box"
	"github.com/influx6/faux/context"
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

// ImageTagSpell defines a structure which implements the Spell interface
// for executing of docker based commands for ImageTag.
type ImageTagSpell struct {
	client *client.Client

	tag string
}

// Spell returns a object implementing the box.Shell interface.
func (cm *ImageTagSpell) Spell(callback ImageTagResponseCallback) box.Spell {
	return &onceImageTagSpell{spell: cm, callback: cb}
}

type onceImageTagSpell struct {
	callback ImageTagResponseCallback
	spell    *ImageTagSpell
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceImageTagSpell) Exec(ctx context.CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *ImageTagSpell) Exec(ctx context.CancelContext, callback ImageTagResponseCallback) error {
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
