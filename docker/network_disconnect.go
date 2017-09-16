package docker

import (
	"github.com/influx6/box"
	"github.com/moby/moby/client"
)

// NetworkDisconnect returns a new NetworkDisconnectSpell instance to be executed on the client.
func (d *DockerCaster) NetworkDisconnect(networkID string) (*NetworkDisconnectSpell, error) {
	var spell NetworkDisconnectSpell

	spell.networkID = networkID

	return &spell, nil
}

// NetworkDisconnectSpell defines a function type to modify internal fields of the NetworkDisconnectSpell.
type NetworkDisconnectOptions func(*NetworkDisconnectSpell)

// NetworkDisconnectResponseCallback defines a function type for NetworkDisconnectSpell response.
type NetworkDisconnectResponseCallback func() error

// NetworkDisconnectSpell defines a structure which implements the Spell interface
// for executing of docker based commands for NetworkDisconnect.
type NetworkDisconnectSpell struct {
	client *client.Client

	networkID string
}

// Spell returns a object implementing the box.Shell interface.
func (cm *NetworkDisconnectSpell) Spell(callback NetworkDisconnectResponseCallback) box.Spell {
	return &onceNetworkDisconnectSpell{spell: cm, callback: cb}
}

type onceNetworkDisconnectSpell struct {
	callback NetworkDisconnectResponseCallback
	spell    *NetworkDisconnectSpell
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceNetworkDisconnectSpell) Exec(ctx box.CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *NetworkDisconnectSpell) Exec(ctx box.CancelContext, callback NetworkDisconnectResponseCallback) error {
	// Execute client NetworkDisconnect method.
	err := cm.client.NetworkDisconnect(cm.networkID)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback()
	}

	return nil
}
