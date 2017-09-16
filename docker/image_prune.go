package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/influx6/box"
	"github.com/influx6/faux/context"
	"github.com/moby/moby/api/types/filters"
	"github.com/moby/moby/client"
)

// ImagePrune returns a new ImagePruneSpell instance to be executed on the client.
func (d *DockerCaster) ImagePrune(args filters.Args) (*ImagePruneSpell, error) {
	var spell ImagePruneSpell

	spell.args = args

	return &spell, nil
}

// ImagePruneSpell defines a function type to modify internal fields of the ImagePruneSpell.
type ImagePruneOptions func(*ImagePruneSpell)

// ImagePruneResponseCallback defines a function type for ImagePruneSpell response.
type ImagePruneResponseCallback func(types.ImagesPruneReport) error

// ImagePruneSpell defines a structure which implements the Spell interface
// for executing of docker based commands for ImagePrune.
type ImagePruneSpell struct {
	client *client.Client

	args filters.Args
}

// Spell returns a object implementing the box.Shell interface.
func (cm *ImagePruneSpell) Spell(callback ImagePruneResponseCallback) box.Spell {
	return &onceImagePruneSpell{spell: cm, callback: cb}
}

type onceImagePruneSpell struct {
	callback ImagePruneResponseCallback
	spell    *ImagePruneSpell
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceImagePruneSpell) Exec(ctx context.CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *ImagePruneSpell) Exec(ctx context.CancelContext, callback ImagePruneResponseCallback) error {
	// Execute client ImagePrune method.
	ret0, err := cm.client.ImagePrune(cm.args)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback(ret0)
	}

	return nil
}
