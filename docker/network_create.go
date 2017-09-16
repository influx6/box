package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/influx6/box"
	"github.com/influx6/faux/context"
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

// NetworkCreateSpell defines a structure which implements the Spell interface
// for executing of docker based commands for NetworkCreate.
type NetworkCreateSpell struct {
	client *client.Client

	network types.NetworkCreate
}

// Spell returns a object implementing the box.Shell interface.
func (cm *NetworkCreateSpell) Spell(callback NetworkCreateResponseCallback) box.Spell {
	return &onceNetworkCreateSpell{spell: cm, callback: cb}
}

type onceNetworkCreateSpell struct {
	callback NetworkCreateResponseCallback
	spell    *NetworkCreateSpell
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceNetworkCreateSpell) Exec(ctx context.CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *NetworkCreateSpell) Exec(ctx context.CancelContext, callback NetworkCreateResponseCallback) error {
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
