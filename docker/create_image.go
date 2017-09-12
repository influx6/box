package docker

import (
	"io"

	"github.com/docker/docker/api/types"
	"github.com/influx6/box"
	"github.com/moby/moby/client"
)

// CreateImage returns a new CreateImageSpell instance to be executed on the client.
func (d *DockerCaster) CreateImage(reader io.ReadCloser) (*CreateImageSpell, error) {
	var spell CreateImageSpell

	spell.reader = reader

	return &spell, nil
}

// CreateImageSpell defines a function type to modify internal fields of the CreateImageSpell.
type CreateImageOptions func(*CreateImageSpell)

// CreateImageResponseCallback defines a function type for CreateImageSpell response.
type CreateImageResponseCallback func(types.ImageLoadResponse) error

// AlwaysCreateImageSpellWith returns a object that always executes the provided CreateImageSpell with the provided callback.
func AlwaysCreateImageSpellWith(bm *CreateImageSpell, cb CreateImageResponseCallback) box.Spell {
	return &onceCreateImageSpell{spell: bm, callback: cb}
}

type onceCreateImageSpell struct {
	callback CreateImageResponseCallback
	spell    *CreateImageSpell
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceCreateImageSpell) Exec(ctx box.CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// CreateImageSpell defines a structure which implements the Spell interface
// for executing of docker based commands for CreateImage.
type CreateImageSpell struct {
	client *client.Client

	reader io.ReadCloser
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *CreateImageSpell) Exec(ctx box.CancelContext, callback CreateImageResponseCallback) error {
	// Execute client CreateImage method.
	ret0, err := cm.client.CreateImage(cm.reader)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback(ret0)
	}

	return nil
}
