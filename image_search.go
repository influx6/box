package dockish

import (
	"github.com/docker/docker/api/types"
	"github.com/moby/moby/api/types/registry"
	"github.com/moby/moby/client"
)

// ImageSearch returns a new ImageSearchSpell instance to be executed on the client.
func (d *DockerCaster) ImageSearch(searchOps types.ImageSearchOptions) (*ImageSearchSpell, error) {
	var spell ImageSearchSpell

	spell.searchOps = searchOps

	return &spell, nil
}

// ImageSearchSpell defines a function type to modify internal fields of the ImageSearchSpell.
type ImageSearchOptions func(*ImageSearchSpell)

// ImageSearchResponseCallback defines a function type for ImageSearchSpell response.
type ImageSearchResponseCallback func([]registry.SearchResult) error

// AlwaysImageSearchSpellWith returns a object that always executes the provided ImageSearchSpell with the provided callback.
func AlwaysImageSearchSpellWith(bm *ImageSearchSpell, cb ImageSearchResponseCallback) Spell {
	return &onceImageSearchSpell{spell: bm, callback: cb}
}

type onceImageSearchSpell struct {
	callback ImageSearchResponseCallback
	spell    *ImageSearchSpell
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceImageSearchSpell) Exec(ctx CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// ImageSearchSpell defines a structure which implements the Spell interface
// for executing of docker based commands for ImageSearch.
type ImageSearchSpell struct {
	client *client.Client

	searchOps types.ImageSearchOptions
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *ImageSearchSpell) Exec(ctx CancelContext, callback ImageSearchResponseCallback) error {
	// Execute client ImageSearch method.
	ret0, err := cm.client.ImageSearch(cm.searchOps)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback(ret0)
	}

	return nil
}
