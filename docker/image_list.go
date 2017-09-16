package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/influx6/box"
	"github.com/moby/moby/client"
)

// ImageList returns a new ImageListSpell instance to be executed on the client.
func (d *DockerCaster) ImageList(listOps types.ImageListOptions) (*ImageListSpell, error) {
	var spell ImageListSpell

	spell.listOps = listOps

	return &spell, nil
}

// ImageListSpell defines a function type to modify internal fields of the ImageListSpell.
type ImageListOptions func(*ImageListSpell)

// ImageListResponseCallback defines a function type for ImageListSpell response.
type ImageListResponseCallback func([]types.ImageSummary) error

// ImageListSpell defines a structure which implements the Spell interface
// for executing of docker based commands for ImageList.
type ImageListSpell struct {
	client *client.Client

	listOps types.ImageListOptions
}

// Spell returns a object implementing the box.Shell interface.
func (cm *ImageListSpell) Spell(callback ImageListResponseCallback) box.Spell {
	return &onceImageListSpell{spell: cm, callback: cb}
}

type onceImageListSpell struct {
	callback ImageListResponseCallback
	spell    *ImageListSpell
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceImageListSpell) Exec(ctx box.CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *ImageListSpell) Exec(ctx box.CancelContext, callback ImageListResponseCallback) error {
	// Execute client ImageList method.
	ret0, err := cm.client.ImageList(cm.listOps)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback(ret0)
	}

	return nil
}
