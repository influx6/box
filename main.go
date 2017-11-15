package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	version   = "0.0.1" // rely on linker -ldflags -X main.version"
	gitCommit = ""      // rely on linker: -ldflags -X main.gitCommit"
)

var (
	getVersion   = flag.Bool("v", false, "Print version")
	forceRebuild = flag.Bool("f", false, "force rebuild")
)

func main() {
	flag.Usage = printUsage
	flag.Parse()

	// if we are to print getVersion.
	if *getVersion {
		printVersion()
		return
	}

	command := flag.Arg(0)
	name := flag.Arg(1)

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get directory path: %+q", err)
		return
	}

	switch command {
	case "command":
	default:
		printUsage()
	}

	log.Println("Trail asset bundling ready!")
}

// printVersion prints corresponding build getVersion with associated build stamp and git commit if provided.
func printVersion() {
	fragments := []string{version}

	if gitCommit != "" {
		fragments = append(fragments, fmt.Sprintf("git#%s", gitCommit))
	}

	fmt.Fprint(os.Stdout, strings.Join(fragments, " "))
}

// printUsage prints out usage message for CLI tool.
func printUsage() {
	fmt.Fprintf(os.Stdout, `Usage: box [options]
Box handles the migration of your app to your server.

COMMANDS:


EXAMPLES:


FLAGS:
	-v       Print version.
	-f 	     Force re-generation of all files
`)
}
