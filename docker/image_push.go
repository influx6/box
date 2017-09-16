package docker

import (
	"io"

	"github.com/docker/docker/api/types"
	"github.com/influx6/box"
	"github.com/moby/moby/client"
)

// ImagePush returns a new ImagePushSpell instance to be executed on the client.
func (d *DockerCaster) ImagePush(imp types.ImagePushOptions) (*ImagePushSpell, error) {
	var spell ImagePushSpell

	spell.imp = imp

	return &spell, nil
}

// ImagePushSpell defines a function type to modify internal fields of the ImagePushSpell.
type ImagePushOptions func(*ImagePushSpell)

// ImagePushResponseCallback defines a function type for ImagePushSpell response.
type ImagePushResponseCallback func(io.ReadCloser) error

// ImagePushSpell defines a structure which implements the Spell interface
// for executing of docker based commands for ImagePush.
type ImagePushSpell struct {
	client *client.Client

	imp types.ImagePushOptions
}

// Spell returns a object implementing the box.Shell interface.
func (cm *ImagePushSpell) Spell(callback ImagePushResponseCallback) box.Spell {
	return &onceImagePushSpell{spell: cm, callback: cb}
}

type onceImagePushSpell struct {
	callback ImagePushResponseCallback
	spell    *ImagePushSpell
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceImagePushSpell) Exec(ctx box.CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *ImagePushSpell) Exec(ctx box.CancelContext, callback ImagePushResponseCallback) error {
	// Execute client ImagePush method.
	ret0, err := cm.client.ImagePush(cm.imp)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback(ret0)
	}

	return nil
}
