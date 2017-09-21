package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/ops"
	"github.com/moby/moby/client"
)

// NetworkCreate returns a new NetworkCreateOp instance to be executed on the client.
func (d *DockerCaster) NetworkCreate(network types.NetworkCreate) (*NetworkCreateOp, error) {
	var spell NetworkCreateOp

	spell.network = network

	return &spell, nil
}

// NetworkCreateOptions defines a function type to modify internal fields of the NetworkCreateOp.
type NetworkCreateOptions func(*NetworkCreateOp)

// NetworkCreateResponseCallback defines a function type for NetworkCreateOp response.
type NetworkCreateResponseCallback func(types.NetworkCreateResponse) error

// NetworkCreateOp defines a structure which implements the Op interface
// for executing of docker based commands for NetworkCreate.
type NetworkCreateOp struct {
	client *client.Client

	network types.NetworkCreate
}

// Op returns a object implementing the ops.Op interface.
func (cm *NetworkCreateOp) Op(callback NetworkCreateResponseCallback) ops.Op {
	return &onceNetworkCreateOp{spell: cm, callback: cb}
}

type onceNetworkCreateOp struct {
	callback NetworkCreateResponseCallback
	spell    *NetworkCreateOp
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceNetworkCreateOp) Exec(ctx context.CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *NetworkCreateOp) Exec(ctx context.CancelContext, callback NetworkCreateResponseCallback) error {
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
