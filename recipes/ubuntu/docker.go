package recipes

import "github.com/influx6/box"
import "github.com/influx6/faux/context"

// InstallDocker handles installation of docker into a ubuntu vm/host.
type InstallDocker struct {
	PreCalls  []box.Spell
	PostCalls []box.Spell
}

// Exec executes giving recipe for building giving docker image for the provided binary.
func (docker InstallDocker) Exec(ctx context.CancelContext) error {

	return nil
}

// installOpenssh will run necessary commands to install openssh.
type installOpenssh struct{}

// Exec executes giving recipe for building giving docker image for the provided binary.
func (ssh installOpenssh) Exec(ctx context.CancelContext) error {

	return nil
}
