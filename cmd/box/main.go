package main

import (
	"fmt"
	"runtime"

	"github.com/fatih/color"
	"github.com/influx6/faux/metrics"
	"github.com/influx6/faux/metrics/sentries/stdout"
	"github.com/minio/cli"
)

var (
	events = metrics.New(stdout.Stdout{})
)

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

}

// versionFn defines the action called when seeking the Version detail.
func versionFn(c *cli.Context) {
	fmt.Println(color.BlueString(fmt.Sprintf("box version %s %s/%s", Version, runtime.GOOS, runtime.GOARCH)))
}
