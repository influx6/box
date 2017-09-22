package ubuntu

import (
	"github.com/influx6/box"
	"github.com/influx6/box/recipes/exec/osinfo"
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/ops"
)

var (
	_ = box.RegisterJSON("linux/ubuntu", func() ops.Op {
		return &ubuntuProvisioner{}
	})
)

// ubuntuProvisioner implements ops.Op interface and contains necessary procedures to provision a
// ubuntu linux vm/system for app deployment with box.
type ubuntuProvisioner struct {
	Info osinfo.Info `json:"os_info"`
}

func (ubp *ubuntuProvisioner) Exec(ctx context.CancelContext) error {

	// Attempt to install necessary packages packages.
	if err := GitInstall().Exec(ctx); err != nil {
		return err
	}

	if err := CurlInstall().Exec(ctx); err != nil {
		return err
	}

	if err := WgetInstall().Exec(ctx); err != nil {
		return err
	}

	if err := OpenSSHInstall().Exec(ctx); err != nil {
		return err
	}

	if err := AptTransportHTTPSInstall().Exec(ctx); err != nil {
		return err
	}

	// Call docker installation from https://get.docker.com
	if err := DockerSourceInstaller.Exec(ctx); err != nil {
		return err
	}

	return nil
}
