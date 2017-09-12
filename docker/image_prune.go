package docker

import (
	"github.com/docker/docker/api/types"
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

// AlwaysImagePruneSpellWith returns a object that always executes the provided ImagePruneSpell with the provided callback.
func AlwaysImagePruneSpellWith(bm *ImagePruneSpell, cb ImagePruneResponseCallback) Spell {
	return &onceImagePruneSpell{spell: bm, callback: cb}
}

type onceImagePruneSpell struct {
	callback ImagePruneResponseCallback
	spell    *ImagePruneSpell
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceImagePruneSpell) Exec(ctx CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// ImagePruneSpell defines a structure which implements the Spell interface
// for executing of docker based commands for ImagePrune.
type ImagePruneSpell struct {
	client *client.Client

	args filters.Args
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *ImagePruneSpell) Exec(ctx CancelContext, callback ImagePruneResponseCallback) error {
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
