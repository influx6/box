package exec_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/influx6/box/recipes/exec"
	"github.com/influx6/faux/tests"
)

func TestLsCommand(t *testing.T) {
	var outs, errs bytes.Buffer
	lsCmd := exec.New(exec.Command("ls ./.."), exec.Sync(), exec.Output(&outs), exec.Err(&errs))
	ctx, cn := context.WithTimeout(context.Background(), 20*time.Second)
	defer cn()

	if err := lsCmd.Exec(ctx); err != nil {
		tests.Info("Output: %+q", outs.Bytes())
		tests.Info("Errs: %+q", errs.Bytes())
		tests.Failed("Should have succcesfully executed command: %+q", err)
	}
	tests.Passed("Should have succcesfully executed command")
	tests.Info("Output: %+q", outs.Bytes())
}
