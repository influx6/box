package ubuntu

import (
	"fmt"
	"os"

	"github.com/influx6/box"
	"github.com/influx6/box/recipes/osinfo"
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/exec"
	"github.com/influx6/faux/metrics"
	"github.com/influx6/faux/ops"

	"github.com/influx6/box/shared/pkgaction"
)

var (
	_ = box.RegisterJSON("linux/ubuntu", func() ops.Op {
		return &ubuntuProvisioner{}
	})
)

// custom package installers
var (
	sudo         = exec.New(exec.Command("apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y sudo"), exec.Async())
	dockerSource = exec.New(exec.Command("wget -nv -O - https://get.docker.com/ | sh"), exec.Async())
	aptUpdate    = exec.New(exec.Async(), exec.Command("sudo apt-get -y update"))
)

// ubuntuProvisioner implements ops.Op interface and contains necessary procedures to provision a
// ubuntu linux vm/system for app deployment with box.
type ubuntuProvisioner struct {
	Info         osinfo.Info `json:"os_info"`
	UpstartBased bool        `json:"upstart"`
}

func (ubp *ubuntuProvisioner) Exec(ctx context.CancelContext, m metrics.Metrics) error {

	if err := exec.New(exec.Command("if ! type sudo; then exit 1; fi")).Exec(ctx, m); err != nil {
		// Attempt to install sudo
		if err := exec.ApplyImmediate(sudo, exec.Output(os.Stdout)).Exec(ctx, m); err != nil {
			return err
		}
	} else {
		// Update apt-get to ensure latest
		if err := exec.ApplyImmediate(aptUpdate, exec.Output(os.Stdout)).Exec(ctx, m); err != nil {
			return err
		}
	}

	// Attempt to install git packages.
	if err := InstallPackage("git", pkgaction.InstallAction, ubp.UpstartBased, exec.Output(os.Stdout)).Exec(ctx, m); err != nil {
		return err
	}

	// Attempt to install wget, curl, ufw packages.
	if err := InstallPackage("curl wget ufw", pkgaction.InstallAction, ubp.UpstartBased, exec.Output(os.Stdout)).Exec(ctx, m); err != nil {
		return err
	}

	// Attempt to install openssh packages.
	if err := InstallPackage("openssh-client sshcommand openssl", pkgaction.InstallAction, ubp.UpstartBased, exec.Output(os.Stdout)).Exec(ctx, m); err != nil {
		return err
	}

	// Attempt to install apt-transport-https packages.
	if err := InstallPackage("apt-transport-https", pkgaction.InstallAction, ubp.UpstartBased, exec.Output(os.Stdout)).Exec(ctx, m); err != nil {
		return err
	}

	// Attempt to install apt-transport-https packages.
	if err := InstallPackage("make autoconf build-essential software-properties-common", pkgaction.InstallAction, ubp.UpstartBased, exec.Output(os.Stdout)).Exec(ctx, m); err != nil {
		return err
	}

	// Call docker installation from https://get.docker.com.
	// TODO(infux6): If this fails then consider manually hooking into https://download.docker.com/.
	if err := exec.New(exec.Command("wget -nv -O - https://get.docker.com/ | sh"), exec.Async(), exec.Output(os.Stdout)).Exec(ctx, m); err != nil {
		return err
	}

	return nil
}

// InstallPackage returns a exec.Command that is executed to install/remove a giving ubuntu package.
func InstallPackage(pkgName string, action pkgaction.PackageAction, upstartbased bool, cmds ...exec.CommanderOption) *exec.Commander {
	var command string

	if action == pkgaction.PurgeAction {
		action = pkgaction.RemoveAction
	}

	switch upstartbased {
	case true:
		command = fmt.Sprintf("DEBIAN_FRONTEND=noninteractive sudo -E apt-get %+s -y -o Dpkg::Options::=\"--force-confnew\" %s", action, pkgName)
	case false:
		command = fmt.Sprintf("DEBIAN_FRONTEND=noninteractive sudo -E apt-get %+s -y %s", action, pkgName)
	}

	return exec.ApplyImmediate(exec.New(exec.Command(command), exec.Async()), cmds...)
}
