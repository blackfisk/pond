package command

import (
        "strings"
        "github.com/mitchellh/cli"
        "github.com/blackfisk/pond/client"
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
        pc := client.NewPondClient("http://localhost:12345")
        pc.Fetch()

        return 0
}

func (c *FetchCommand) Synopsis() string {
        return "Fetches messages from a pond server"
}
