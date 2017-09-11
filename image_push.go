package dockish

import (
	"io"

	"github.com/docker/docker/api/types"
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

// AlwaysImagePushSpellWith returns a object that always executes the provided ImagePushSpell with the provided callback.
func AlwaysImagePushSpellWith(bm *ImagePushSpell, cb ImagePushResponseCallback) Spell {
	return &onceImagePushSpell{spell: bm, callback: cb}
}

type onceImagePushSpell struct {
	callback ImagePushResponseCallback
	spell    *ImagePushSpell
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceImagePushSpell) Exec(ctx CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// ImagePushSpell defines a structure which implements the Spell interface
// for executing of docker based commands for ImagePush.
type ImagePushSpell struct {
	client *client.Client

	imp types.ImagePushOptions
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *ImagePushSpell) Exec(ctx CancelContext, callback ImagePushResponseCallback) error {
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
