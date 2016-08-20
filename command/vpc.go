package command

import (
	"fmt"
	"github.com/BSick7/aws-topology-api/api"
	"github.com/BSick7/aws-topology-api/services"
	"github.com/mitchellh/cli"
	"strings"
)

type VpcCommand struct {
	Meta
}

func (c *VpcCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("vpc")
	flags.Usage = func() {
		c.Ui.Error(c.Help())
	}
	if err := flags.Parse(args); err != nil {
		return 1
	}

	vpcId := args[len(args)-1]
	if vpcId == "" {
		c.Ui.Error("vpc-id is required")
		return cli.RunResultHelp
	}

	b := services.NewBroker(nil)
	root, err := api.GetVpcTopology(b, vpcId)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	c.Ui.Output(fmt.Sprintf("%+v\n", root))
	return 0
}

func (c *VpcCommand) Synopsis() string {
	return "Collect vpc topology"
}

func (c *VpcCommand) Help() string {
	helpText := `
Usage: aws-topology-api vpc [options] <vpc-id>

Agent Options:

`
	return strings.TrimSpace(helpText)
}
