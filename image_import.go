package dockish

import (
	"io"

	"github.com/docker/docker/api/types"
	"github.com/moby/moby/client"
)

// ImageImport returns a new ImageImportSpell instance to be executed on the client.
func (d *DockerCaster) ImageImport(impOp types.ImageImportOptions) (*ImageImportSpell, error) {
	var spell ImageImportSpell

	spell.impOp = impOp

	return &spell, nil
}

// ImageImportSpell defines a function type to modify internal fields of the ImageImportSpell.
type ImageImportOptions func(*ImageImportSpell)

// ImageImportResponseCallback defines a function type for ImageImportSpell response.
type ImageImportResponseCallback func(io.ReadCloser) error

// AlwaysImageImportSpellWith returns a object that always executes the provided ImageImportSpell with the provided callback.
func AlwaysImageImportSpellWith(bm *ImageImportSpell, cb ImageImportResponseCallback) Spell {
	return &onceImageImportSpell{spell: bm, callback: cb}
}

type onceImageImportSpell struct {
	callback ImageImportResponseCallback
	spell    *ImageImportSpell
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceImageImportSpell) Exec(ctx CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// ImageImportSpell defines a structure which implements the Spell interface
// for executing of docker based commands for ImageImport.
type ImageImportSpell struct {
	client *client.Client

	impOp types.ImageImportOptions
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *ImageImportSpell) Exec(ctx CancelContext, callback ImageImportResponseCallback) error {
	// Execute client ImageImport method.
	ret0, err := cm.client.ImageImport(cm.impOp)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback(ret0)
	}

	return nil
}
