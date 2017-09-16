package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/influx6/box"
	"github.com/moby/moby/client"
)

// NetworkInspect returns a new NetworkInspectSpell instance to be executed on the client.
func (d *DockerCaster) NetworkInspect(netOp types.NetworkInspectOptions) (*NetworkInspectSpell, error) {
	var spell NetworkInspectSpell

	spell.netOp = netOp

	return &spell, nil
}

// NetworkInspectSpell defines a function type to modify internal fields of the NetworkInspectSpell.
type NetworkInspectOptions func(*NetworkInspectSpell)

// NetworkInspectResponseCallback defines a function type for NetworkInspectSpell response.
type NetworkInspectResponseCallback func(types.NetworkResource) error

// NetworkInspectSpell defines a structure which implements the Spell interface
// for executing of docker based commands for NetworkInspect.
type NetworkInspectSpell struct {
	client *client.Client

	netOp types.NetworkInspectOptions
}

// Spell returns a object implementing the box.Shell interface.
func (cm *NetworkInspectSpell) Spell(callback NetworkInspectResponseCallback) box.Spell {
	return &onceNetworkInspectSpell{spell: cm, callback: cb}
}

type onceNetworkInspectSpell struct {
	callback NetworkInspectResponseCallback
	spell    *NetworkInspectSpell
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceNetworkInspectSpell) Exec(ctx box.CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *NetworkInspectSpell) Exec(ctx box.CancelContext, callback NetworkInspectResponseCallback) error {
	// Execute client NetworkInspect method.
	ret0, err := cm.client.NetworkInspect(cm.netOp)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback(ret0)
	}

	return nil
}
