package command

import (
	"fmt"
	"github.mdl.zone/deployer/deploy/agent"
	"github.mdl.zone/deployer/deploy/types"
	"strings"
)

type AgentCommand struct {
	Meta
}

func (c *AgentCommand) Run(args []string) int {
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

	ac, err := types.NewAgentConfigFromLocation(configLoc)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("error reading configuration [%s]: %s", configLoc, err))
		return 1
	}
	ac.Merge(&types.AgentConfig{
		Bind: types.AgentConfigBind{
			Address: bindAddr,
			Port:    bindPort,
		},
	})
	if err := ac.Validate(); err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	if err := agent.Run(ac, c.Meta.Version); err != nil {
		c.Ui.Error(fmt.Sprintf("error running agent: %s", err))
		return 1
	}
	return 0
}

func (c *AgentCommand) Synopsis() string {
	return "Run deploy agent"
}

func (c *AgentCommand) Help() string {
	helpText := `
Usage: deploy agent [options]

Agent Options:

  -config=file|directory   Loads configuration from directory of
                           *.hcl files or individual file.hcl.

  -bind-addr=address       Uses the address to bind the API server.

  -bind-port=port          Uses the port to bind the API server.
                           Defaults to 8080
`
	return strings.TrimSpace(helpText)
}
