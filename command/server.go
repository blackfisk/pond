package command

import (
	"net/http"
	"os"
	"fmt"

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

func (c *ServerCommand) Run(args []string) int {
	p := pond.NewPond()
	port := ":" + os.Getenv("PORT")

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
