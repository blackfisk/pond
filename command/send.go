package command

import (
	"flag"
	"github.com/blackfisk/pond/client"
	"github.com/mitchellh/cli"
	"strings"
)

type SendCommand struct {
	Ui cli.Ui
}

func (c *SendCommand) Help() string {
	helpText := `
        Usage: pond send email [options]

        Fetches messages from pond

        Options:

        `
	return strings.TrimSpace(helpText)
}

func (c *SendCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("send", flag.ContinueOnError)
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}
	args = cmdFlags.Args()

	if len(args) < 1 {
		c.Ui.Error("You need to tell me to who should I send it")
		c.Ui.Error("")
		c.Ui.Error(c.Help())
		return 1
	}

	email := args[0]
	message := args[1]

	pc := client.NewPondClient("http://localhost:12345")
        sent := pc.Send(email, message)
        if sent {
                c.Ui.Output("Message sent!")
                return 1
        } else {
                c.Ui.Output("Problems sending the message")
                return 1
        }
}

func (c *SendCommand) Synopsis() string {
	return "Fetches messages from a pond server"
}
