package ubuntu

import (
	"fmt"

	"github.com/influx6/box/recipes/exec"
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/metrics"
)

// custom package installers
var (
	OpenSSHInstall           = PkgPartial(PkgCommand("openssh", InstallAction), Debian())
	GitInstall               = PkgPartial(PkgCommand("git", InstallAction), Debian())
	CurlInstall              = PkgPartial(PkgCommand("curl", InstallAction), Debian())
	WgetInstall              = PkgPartial(PkgCommand("wget", InstallAction), Debian())
	AxelInstall              = PkgPartial(PkgCommand("axel", InstallAction), Debian())
	AptTransportHTTPSInstall = PkgPartial(PkgCommand("apt-transport-https", InstallAction), Debian())
	DockerEngineInstall      = PkgPartial(PkgCommand("docker-engine", InstallAction), Debian())
	DockerCEInstall          = PkgPartial(PkgCommand("docker-ce", InstallAction), UbuntuSystemd())
)

// custom executors.
var (
	SudoInstaller         = exec.New(exec.Command("if ! type sudo; then apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y sudo; fi"))
	DockerSourceInstaller = exec.New(exec.Command("wget -nv -O - https://get.docker.com/ | sh"))
)

//===============================================================================================================

// PackageSourceUpdate will run necessary commands to update apt-get for a debian/ubuntu system.
type PackageSourceUpdate struct {
	DoWithCmd exec.CommanderOption
}

// Exec executes giving recipe for building giving docker image for the provided binary.
func (pkgSrc PackageSourceUpdate) Exec(ctx context.CancelContext, m metrics.Metrics) error {
	cmd := exec.New(exec.Command("sudo apt-get -y update"), exec.Async())

	if pkgSrc.DoWithCmd != nil {
		pkgSrc.DoWithCmd(cmd)
	}

	return cmd.Exec(ctx, m)
}

//===============================================================================================================

// PacakgeAction defines a int type to represent a package action for a package installer.
type PackageAction int

// PackageOption defines a package option to set option for packageInstaller.
type PackageOption func(*PackageInstaller)

// PackagePartial defines a function type that returns a PkgInstaller from a series of PackageOptions.
type PackagePartial func(...PackageOption) *PackageInstaller

// CommandOption sets the function to be used to be used as function for cmd.
func CommandOption(cmd exec.CommanderOption) PackageOption {
	return func(pkg *PackageInstaller) {
		pkg.DoWithCmd = cmd
	}
}

// UbuntuSystemd sets the flag to be used as systemd ubuntu.
func UbuntuSystemd() PackageOption {
	return func(pkg *PackageInstaller) {
		pkg.systemdUbuntu = true
		pkg.upstartUbuntu = false
		pkg.debian = false
	}
}

// PkgCommand sets the command and action to be used as pkg operation.
func PkgCommand(name string, action PackageAction) PackageOption {
	return func(pkg *PackageInstaller) {
		pkg.Name = name
		pkg.Action = action
	}
}

// PkgApply applies all options to provided PackageInstaller.
func PkgApply(ops ...PackageOption) PackageOption {
	return func(pkg *PackageInstaller) {
		for _, op := range ops {
			op(pkg)
		}
	}
}

// PkgPartial applies all options to the returned PackageInstaller.
func PkgPartial(ops ...PackageOption) PackagePartial {
	return func(more ...PackageOption) *PackageInstaller {
		pkg := new(PackageInstaller)

		for _, op := range append(ops, more...) {
			op(pkg)
		}

		return pkg
	}
}

// Debian sets the flag to be used as upstart ubuntu.
func Debian() PackageOption {
	return func(pkg *PackageInstaller) {
		pkg.systemdUbuntu = false
		pkg.upstartUbuntu = true
		pkg.debian = true
	}
}

// UbuntuUpstart sets the flag to be used as upstart ubuntu.
func UbuntuUpstart() PackageOption {
	return func(pkg *PackageInstaller) {
		pkg.systemdUbuntu = false
		pkg.upstartUbuntu = true
		pkg.debian = false
	}
}

// String returns the name of the action.
func (ap PackageAction) String() string {
	switch ap {
	case InstallAction:
		return "install"
	case RemoveAction:
		return "remove"
	case PurgeAction:
		return "purge"
	}

	return "unknown"
}

// pkg constant types
const (
	InstallAction PackageAction = iota
	RemoveAction
	PurgeAction
)

// PackageInstaller will run necessary commands to install giving package.
type PackageInstaller struct {
	Action        PackageAction
	Name          string
	debian        bool // set to true if debian system
	upstartUbuntu bool //set true if its ubuntu with upstart
	systemdUbuntu bool //set true if ubuntu with systemd
	DoWithCmd     exec.CommanderOption
}

// Exec executes giving recipe for building giving docker image for the provided binary.
func (pkg *PackageInstaller) Exec(ctx context.CancelContext, m metrics.Metrics) error {
	var command string

	switch {
	case pkg.debian:
		command = fmt.Sprintf("DEBIAN_FRONTEND=noninteractive sudo -E apt-get %+s -y  %s", pkg.Name, pkg.Name)
	case pkg.upstartUbuntu:
		if pkg.Action == PurgeAction {
			pkg.Action = RemoveAction
		}

		command = fmt.Sprintf("DEBIAN_FRONTEND=noninteractive sudo -E apt-get %+s -y -o Dpkg::Options::=\"--force-confnew\" %s", pkg.Action, pkg.Name)
	case pkg.systemdUbuntu:
		command = fmt.Sprintf("DEBIAN_FRONTEND=noninteractive sudo -E apt-get %+s -y  %s", pkg.Action, pkg.Name)
	}

	cmd := exec.New(exec.Command(command), exec.Async())

	if pkg.DoWithCmd != nil {
		pkg.DoWithCmd(cmd)
	}

	return cmd.Exec(ctx, m)
}
