package command

import (
	"github.com/blackfisk/pond/client"
	"github.com/mitchellh/cli"
	"strings"
	"flag"
)

type FetchCommand struct {
	Ui cli.Ui
}

func (c *FetchCommand) Help() string {
	helpText := `
        Usage: pond send [options] to

        Send messages to a pond

        Options:

        `
	return strings.TrimSpace(helpText)
}

func (c *FetchCommand) Run(args []string) int {
        var address string

	cmdFlags := flag.NewFlagSet("fetch", flag.ContinueOnError)
        cmdFlags.StringVar(&address, "pond", "", "Define the address of the pond")

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

        if address == "" {
		c.Ui.Error("You need to define a destination -pond")
		c.Ui.Error("")

                return 1
        }

	pc := client.NewPondClient(address)
	pc.Fetch()

	return 0
}

func (c *FetchCommand) Synopsis() string {
	return "Fetches messages from a pond server"
}
