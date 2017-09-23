package exec

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/influx6/faux/context"
	"github.com/influx6/faux/metrics"
)

const (
	ExecLogKey = "recipe:exec"
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

// Binary sets the binary command for the Commander.
func Binary(bin string, flag string) CommanderOption {
	return func(cm *Commander) {
		cm.Binary = bin
		cm.Flag = flag
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

// Metric sets the metric field for logging for the Commander.
func Metric(m metrics.Metrics) CommanderOption {
	return func(cm *Commander) {
		cm.Metrics = m
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

// Apply takes the giving series of CommandOption returning a function that always applies them to passed in commanders.
func Apply(ops ...CommanderOption) CommanderOption {
	return func(cm *Commander) {
		for _, op := range ops {
			op(cm)
		}
	}
}

// Commander runs provided command within a /bin/sh -c "{COMMAND}", returning
// response associatedly. It also attaches if provided stdin, stdout and stderr readers/writers.
// Commander allows you to set the binary to use and flag, where each defaults to /bin/sh for binary
// and -c for flag respectively.
type Commander struct {
	Async   bool
	Command string
	Binary  string
	Flag    string
	Envs    map[string]string
	In      io.Reader
	Out     io.Writer
	Err     io.Writer
	Metrics metrics.Metrics
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
	if c.Metrics == nil {
		c.Metrics = metrics.New()
	}

	if c.Binary == "" {
		c.Binary = "/bin/sh"
	}

	if c.Flag == "" {
		c.Flag = "-c"
	}

	var execCommand []string

	switch {
	case c.Command == "":
		execCommand = append(execCommand, c.Binary)
	case c.Command != "":
		execCommand = append(execCommand, c.Binary, c.Flag, c.Command)
	}

	cmder := exec.Command(execCommand[0], execCommand[1:]...)
	cmder.Stderr = c.Err
	cmder.Stdin = c.In
	cmder.Stdout = c.Out
	cmder.Env = os.Environ()

	if c.Envs != nil {
		for name, val := range c.Envs {
			cmder.Env = append(cmder.Env, fmt.Sprintf("%s=%s", name, val))
		}
	}

	c.Metrics.Emit(metrics.WithFields(metrics.Fields{
		"opid":    ExecLogKey,
		"command": execCommand,
		"envs":    cmder.Env,
	}).WithMessage("Executing native commands"))

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
