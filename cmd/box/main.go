package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/influx6/box"
	"github.com/influx6/box/shared/consts"
	"github.com/influx6/faux/metrics"
	"github.com/influx6/faux/metrics/custom"
	"github.com/minio/cli"

	_ "github.com/influx6/box/recipes/darwin"
	_ "github.com/influx6/box/recipes/linux"
	_ "github.com/influx6/box/recipes/windows"
)

var (
	red    = color.New(color.FgRed)
	yellow = color.New(color.FgYellow)
	green  = color.New(color.FgGreen)
	events = metrics.Switch(metrics.MetricKey, map[string]metrics.Metrics{
		consts.CliOpLog:              custom.StackDisplayWith(os.Stdout, "[Op]", "-", nil),
		consts.RecipeExecLogErrorKey: custom.NewCustomEmitter(os.Stdout, execCommandError),
		consts.RecipeExecLogKey:      custom.NewCustomEmitter(os.Stdout, execCommandOperation),
		consts.CliErrLog:             metrics.New(errorDisplay, custom.StackDisplayWith(os.Stdout, red.Sprint("X"), red.Sprint("-"), nil)),
	})
)

func execCommandOperation(m metrics.Entry) []byte {
	var bu bytes.Buffer

	if commands, ok := m.Get("command"); ok {
		green.Fprint(&bu, "- Executing Shell Command\n")
		yellow.Fprintf(&bu, "Command: ")
		fmt.Fprintf(&bu, "%+q\n", commands)
		fmt.Fprint(&bu, "\n")
	}

	return bu.Bytes()
}

func execCommandError(m metrics.Entry) []byte {
	var bu bytes.Buffer

	if commands, ok := m.Get("command"); ok {
		red.Fprint(&bu, "X Executed Shell Command\n")
		yellow.Fprintf(&bu, "Command: ")
		fmt.Fprintf(&bu, "%+q\n", commands)

		if commandErr, ok := m.Get("error"); ok {
			red.Fprintf(&bu, "Error: ")
			fmt.Fprintf(&bu, "%+q\n", commandErr)
		}

		if commandErrData, ok := m.Get("error_data"); ok {
			red.Fprintf(&bu, "Error Data: ")
			fmt.Fprintf(&bu, "%+q\n", commandErrData)
		}

		fmt.Fprint(&bu, "\n")
	}

	return bu.Bytes()
}

func errorDisplay(m metrics.Entry) metrics.Entry {
	m.Message = red.Sprintf(m.Message)
	return m
}

// Version defines the version number for the cli.
var Version = "0.1"

var helpTemplate = `NAME:
{{.Name}} - {{.Usage}}

DESCRIPTION:
{{.Description}}

USAGE:
{{.Name}} {{if .Flags}}[flags] {{end}}command{{if .Flags}}{{end}} [arguments...]

COMMANDS:
	{{range .Commands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
	{{end}}{{if .Flags}}
FLAGS:
	{{range .Flags}}{{.}}
	{{end}}{{end}}
VERSION:
` + Version +
	`{{ "\n"}}`

// Cmd defines a struct for defining a command.
type Cmd struct {
	*cli.App
}

func main() {
	app := cli.NewApp()
	app.Name = "box"
	app.Author = ""
	app.Usage = "box {{command}}"
	app.Flags = []cli.Flag{}
	app.Description = "box: One cli to rule them all with docker."
	app.CustomAppHelpTemplate = helpTemplate

	app.Commands = []cli.Command{
		{
			Name:   "version",
			Action: versionFn,
			Flags:  []cli.Flag{},
		},
		{
			Name:        "init",
			Action:      initFn,
			Description: "Runs all needed actions to install and provision the host for hosting docker containers",
			Flags:       []cli.Flag{},
		},
	}

	app.Before = func(c *cli.Context) error {
		return nil
	}

	app.RunAndExitOnError()
}

func initFn(c *cli.Context) {
	exec, err := box.CreateWithJSON(strings.ToLower(runtime.GOOS), map[string]interface{}{})
	if err != nil {
		events.Emit(metrics.WithKey(consts.CliErrLog).With("error", err).WithMessage("Failed to find provisioner for %q", runtime.GOOS))
		return
	}

	if err := exec.Exec(context.Background(), events); err != nil {
		events.Emit(metrics.WithKey(consts.CliErrLog).With("error", err).WithMessage("Failed to run provisioner for %q", runtime.GOOS))
		return
	}
}

// versionFn defines the action called when seeking the Version detail.
func versionFn(c *cli.Context) {
	fmt.Println(color.BlueString(fmt.Sprintf("box v%s %s/%s", Version, runtime.GOOS, runtime.GOARCH)))
}
