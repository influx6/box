package dockish

import (
	"github.com/docker/docker/api/types"
	"github.com/moby/moby/client"
)

// NetworkCreate returns a new NetworkCreateSpell instance to be executed on the client.
func (d *DockerCaster) NetworkCreate(network types.NetworkCreate) (*NetworkCreateSpell, error) {
	var spell NetworkCreateSpell

	spell.network = network

	return &spell, nil
}

// NetworkCreateSpell defines a function type to modify internal fields of the NetworkCreateSpell.
type NetworkCreateOptions func(*NetworkCreateSpell)

// NetworkCreateResponseCallback defines a function type for NetworkCreateSpell response.
type NetworkCreateResponseCallback func(types.NetworkCreateResponse) error

// AlwaysNetworkCreateSpellWith returns a object that always executes the provided NetworkCreateSpell with the provided callback.
func AlwaysNetworkCreateSpellWith(bm *NetworkCreateSpell, cb NetworkCreateResponseCallback) Spell {
	return &onceNetworkCreateSpell{spell: bm, callback: cb}
}

type onceNetworkCreateSpell struct {
	callback NetworkCreateResponseCallback
	spell    *NetworkCreateSpell
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceNetworkCreateSpell) Exec(ctx CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// NetworkCreateSpell defines a structure which implements the Spell interface
// for executing of docker based commands for NetworkCreate.
type NetworkCreateSpell struct {
	client *client.Client

	network types.NetworkCreate
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *NetworkCreateSpell) Exec(ctx CancelContext, callback NetworkCreateResponseCallback) error {
	// Execute client NetworkCreate method.
	ret0, err := cm.client.NetworkCreate(cm.network)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback(ret0)
	}

	return nil
}
