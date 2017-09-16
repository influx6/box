package exec

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/influx6/faux/context"
)

var (
	ErrCommandFailed = errors.New("Command failed to execute succcesfully")
)

// CommanderOption defines a function type that aguments a commander's field.
type CommanderOption func(*Commander)

// Command sets the command for the Commander.
func Command(c string) CommanderOption {
	return func(cm *Commander) {
		cm.Command = c
	}
}

// Sync sets the commander to run in synchronouse mode.
func Sync() CommanderOption {
	return SetAsync(false)
}

// ASync sets the commander to run in asynchronouse mode.
func Async() CommanderOption {
	return SetAsync(true)
}

// SetAsync sets the command for the Commander.
func SetAsync(b bool) CommanderOption {
	return func(cm *Commander) {
		cm.Async = b
	}
}

// Input sets the input reader for the Commander.
func Input(in io.Reader) CommanderOption {
	return func(cm *Commander) {
		cm.In = in
	}
}

// Output sets the output writer for the Commander.
func Output(out io.Writer) CommanderOption {
	return func(cm *Commander) {
		cm.Out = out
	}
}

// Err sets the error writer for the Commander.
func Err(err io.Writer) CommanderOption {
	return func(cm *Commander) {
		cm.Err = err
	}
}

// Envs sets the map of environment for the Commander.
func Envs(envs map[string]string) CommanderOption {
	return func(cm *Commander) {
		cm.Envs = envs
	}
}

// Commander runs provided command within a /bin/sh -c "{COMMAND}", returning
// response associatedly. It also attaches if provided stdin, stdout and stderr readers/writers.
type Commander struct {
	Async   bool
	Command string
	Envs    map[string]string
	In      io.Reader
	Out     io.Writer
	Err     io.Writer
}

// New returns a new Commander instance.
func New(ops ...CommanderOption) *Commander {
	cm := new(Commander)

	for _, op := range ops {
		op(cm)
	}

	return cm
}

// Exec executes giving command associated within the command with os/exec.
func (c *Commander) Exec(ctx context.CancelContext) error {
	var cmder *exec.Cmd

	switch {
	case c.Command == "":
		cmder = exec.Command("/bin/sh")
	case c.Command != "":
		cmder = exec.Command("/bin/sh", "-c", c.Command)
	}

	cmder.Stderr = c.Err
	cmder.Stdin = c.In
	cmder.Stdout = c.Out
	cmder.Env = os.Environ()

	if c.Envs != nil {
		for name, val := range c.Envs {
			cmder.Env = append(cmder.Env, fmt.Sprintf("%s=%s", name, val))
		}
	}
	if !c.Async {
		return cmder.Run()
	}

	if err := cmder.Start(); err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		if cmder.Process == nil {
			return
		}

		cmder.Process.Kill()
	}()

	if err := cmder.Wait(); err != nil {
		return err
	}

	if cmder.ProcessState == nil {
		return nil
	}

	if !cmder.ProcessState.Success() {
		return ErrCommandFailed
	}

	return nil
}
