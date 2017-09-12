package docker

import (
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

// AlwaysNetworkDisconnectSpellWith returns a object that always executes the provided NetworkDisconnectSpell with the provided callback.
func AlwaysNetworkDisconnectSpellWith(bm *NetworkDisconnectSpell, cb NetworkDisconnectResponseCallback) Spell {
	return &onceNetworkDisconnectSpell{spell: bm, callback: cb}
}

type onceNetworkDisconnectSpell struct {
	callback NetworkDisconnectResponseCallback
	spell    *NetworkDisconnectSpell
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceNetworkDisconnectSpell) Exec(ctx CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// NetworkDisconnectSpell defines a structure which implements the Spell interface
// for executing of docker based commands for NetworkDisconnect.
type NetworkDisconnectSpell struct {
	client *client.Client

	networkID string
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *NetworkDisconnectSpell) Exec(ctx CancelContext, callback NetworkDisconnectResponseCallback) error {
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
