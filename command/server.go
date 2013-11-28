package command

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/blackfisk/pond/pond"
	"github.com/mitchellh/cli"
	"strings"
)

type ServerCommand struct {
	Ui cli.Ui
}

func (c *ServerCommand) Help() string {
	helpText := `
        Usage: pond server [options]

        Starts a pond server

        Options:

        `
	return strings.TrimSpace(helpText)
}

type ponds []string

func (p *ponds) String() string {
	return fmt.Sprintf("%d", *p)
}

// The second method is Set(value string) error
func (p *ponds) Set(value string) error {
	*p = append(*p, value)
	return nil
}

func (c *ServerCommand) Run(args []string) int {
	var ponds ponds
        var port_to_use string

	cmdFlags := flag.NewFlagSet("server", flag.ContinueOnError)
	cmdFlags.Var(&ponds, "pond", "Define the address of the pond")

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	p := pond.NewPond(ponds)

        if os.Getenv("PORT") != "" {
                port_to_use = os.Getenv("PORT")
        } else {
                port_to_use = "4242"
        }

	port := ":" + port_to_use

	http.Handle("/", p)

	c.Ui.Output(fmt.Sprintf("--> Listening in %s", port))
	err := http.ListenAndServe(port, nil)

	if err != nil {
		c.Ui.Error(fmt.Sprintf("ListenAndServe: ", err))
	}

	return 0
}

func (c *ServerCommand) Synopsis() string {
	return "Starts a Pond server"
}
