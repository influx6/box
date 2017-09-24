package pkgaction

// pkg constant types
const (
	InstallAction PackageAction = iota
	RemoveAction
	PurgeAction
)

// PacakgeAction defines a int type to represent a package action for a package installer.
type PackageAction int

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
