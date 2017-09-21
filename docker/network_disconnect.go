package docker

import (
	"context"

	"github.com/influx6/faux/context"
	"github.com/influx6/faux/ops"
	"github.com/moby/moby/client"
)

// NetworkDisconnect returns a new NetworkDisconnectOp instance to be executed on the client.
func (d *DockerCaster) NetworkDisconnect(networkID string) (*NetworkDisconnectOp, error) {
	var spell NetworkDisconnectOp

	spell.networkID = networkID

	return &spell, nil
}

// NetworkDisconnectOptions defines a function type to modify internal fields of the NetworkDisconnectOp.
type NetworkDisconnectOptions func(*NetworkDisconnectOp)

// NetworkDisconnectResponseCallback defines a function type for NetworkDisconnectOp response.
type NetworkDisconnectResponseCallback func() error

// NetworkDisconnectOp defines a structure which implements the Op interface
// for executing of docker based commands for NetworkDisconnect.
type NetworkDisconnectOp struct {
	client *client.Client

	networkID string
}

// Op returns a object implementing the ops.Op interface.
func (cm *NetworkDisconnectOp) Op(callback NetworkDisconnectResponseCallback) ops.Op {
	return &onceNetworkDisconnectOp{spell: cm, callback: cb}
}

type onceNetworkDisconnectOp struct {
	callback NetworkDisconnectResponseCallback
	spell    *NetworkDisconnectOp
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceNetworkDisconnectOp) Exec(ctx context.CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *NetworkDisconnectOp) Exec(ctx context.CancelContext, callback NetworkDisconnectResponseCallback) error {
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
