package docker

import (
	"io"

	"github.com/docker/docker/api/types"
	"github.com/influx6/box"
	"github.com/moby/moby/client"
)

// ImagePull returns a new ImagePullSpell instance to be executed on the client.
func (d *DockerCaster) ImagePull(imgOp types.ImagePullOptions) (*ImagePullSpell, error) {
	var spell ImagePullSpell

	spell.imgOp = imgOp

	return &spell, nil
}

// ImagePullSpell defines a function type to modify internal fields of the ImagePullSpell.
type ImagePullOptions func(*ImagePullSpell)

// ImagePullResponseCallback defines a function type for ImagePullSpell response.
type ImagePullResponseCallback func(io.ReadCloser) error

// ImagePullSpell defines a structure which implements the Spell interface
// for executing of docker based commands for ImagePull.
type ImagePullSpell struct {
	client *client.Client

	imgOp types.ImagePullOptions
}

// Spell returns a object implementing the box.Shell interface.
func (cm *ImagePullSpell) Spell(callback ImagePullResponseCallback) box.Spell {
	return &onceImagePullSpell{spell: cm, callback: cb}
}

type onceImagePullSpell struct {
	callback ImagePullResponseCallback
	spell    *ImagePullSpell
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceImagePullSpell) Exec(ctx box.CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *ImagePullSpell) Exec(ctx box.CancelContext, callback ImagePullResponseCallback) error {
	// Execute client ImagePull method.
	ret0, err := cm.client.ImagePull(cm.imgOp)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback(ret0)
	}

	return nil
}
