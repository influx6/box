package main

import "flag"

// consts ...
const (
	entryPort   = 34567
	securedPort = 34568
	boxDir      = ".box"
)

func main() {
	flag.Parse()

}

// serve intiailizesthe box service and creates appropriate profile and
// folder to service incoming request.
func serve(secret string, serverName string, company string) error {
	return nil
}
