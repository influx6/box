package main

import (
	"fmt"

	"github.com/influx6/faux/context"
	"github.com/influx6/faux/flags"
)

// consts ...
const (
	entryPort   = 34567
	securedPort = 34568
	boxDir      = ".box"
)

func main() {
	flags.Run("box", flags.Command{
		Name:      "Serve",
		ShortDesc: "runs a box server on the host.",
		Desc:      "sets up a box server, turning the system into a box container server",
		Action: func(ctx context.Context) error {
			fmt.Println("Welcome to box serve!")

			secport, _ := ctx.Bag().GetInt("secport")
			fmt.Printf("SecPort: %d\n", secport)

			cmdport, _ := ctx.Bag().GetInt("cmdport")
			fmt.Printf("CmdPort: %d\n", cmdport)
			return nil
		},
		Flags: []flags.Flag{
			&flags.IntFlag{
				Name:    "secport",
				Default: 34567,
				Desc:    "Port for security tcp server",
			},
			&flags.IntFlag{
				Name:    "cmdport",
				Default: 34568,
				Desc:    "Port for commander tcp server",
			},
		},
	})
}

// serve intiailizesthe box service and creates appropriate profile and
// folder to service incoming request.
func serve(secret string, serverName string, company string) error {
	return nil
}
