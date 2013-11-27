package main

import (
	"fmt"
	"github.com/blackfisk/pond/command"
	"github.com/mitchellh/cli"
	"io/ioutil"
	"log"
	"os"
)

// Commands is the mapping of all the available Pond commands.
var Commands map[string]cli.CommandFactory

func main() {
	os.Exit(realMain())
}

func realMain() int {
	log.SetOutput(ioutil.Discard)

	// Get the command line args. We shortcut "--version" and "-v" to
	// just show the version.
	args := os.Args[1:]
	for _, arg := range args {
		if arg == "-v" || arg == "--version" {
			newArgs := make([]string, len(args)+1)
			newArgs[0] = "version"
			copy(newArgs[1:], args)
			args = newArgs
			break
		}
	}

	cli := &cli.CLI{
		Args:     args,
		Commands: Commands,
	}

	exitCode, err := cli.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		return 1
	}

	return exitCode

}

func init() {
	ui := &cli.BasicUi{Writer: os.Stdout}
	Commands = map[string]cli.CommandFactory{
		"server": func() (cli.Command, error) {
			return &command.ServerCommand{Ui: ui}, nil
		},

		"fetch": func() (cli.Command, error) {
			return &command.FetchCommand{Ui: ui}, nil
		},

		"send": func() (cli.Command, error) {
			return &command.SendCommand{Ui: ui}, nil
		},

		// send in anonymous
		// daemon

	}
}
