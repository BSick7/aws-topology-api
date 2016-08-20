package command

import (
	"strings"
)

type ServeCommand struct {
	Meta
}

func (c *ServeCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("agent")
	var configLoc string
	var bindAddr string
	var bindPort int
	flags.StringVar(&configLoc, "config", "", "")
	flags.StringVar(&bindAddr, "bind-addr", "", "")
	flags.IntVar(&bindPort, "bind-port", 0, "")
	flags.Usage = func() {
		c.Ui.Error(c.Help())
	}
	if err := flags.Parse(args); err != nil {
		return 1
	}

	return 0
}

func (c *ServeCommand) Synopsis() string {
	return "Run aws-topology-api server"
}

func (c *ServeCommand) Help() string {
	helpText := `
Usage: aws-topology-api serve [options]

Agent Options:

  -config=file|directory   Loads configuration from directory of
                           *.hcl files or individual file.hcl.

  -bind-addr=address       Uses the address to bind the API server.

  -bind-port=port          Uses the port to bind the API server.
                           Defaults to 8080
`
	return strings.TrimSpace(helpText)
}
