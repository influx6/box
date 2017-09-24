package debian

import (
	"fmt"

	"github.com/influx6/box/recipes/exec"
	"github.com/influx6/box/shared/pkgaction"
)

// InstallPackage returns a exec.Command that is executed to install/remove a giving ubuntu package.
func InstallPackage(pkgName string, action pkgaction.PackageAction) *exec.Commander {
	command := fmt.Sprintf("DEBIAN_FRONTEND=noninteractive sudo -E apt-get %+s -y  %s", action, pkgName)
	return exec.New(exec.Command(command), exec.Async())
}
