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

	return nil
}
