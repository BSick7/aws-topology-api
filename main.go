package main

import (
	"fmt"
	"github.com/BSick7/aws-topology-api/command"
	"github.com/mitchellh/cli"
	"os"
)

var Version string

func main() {
	c := cli.NewCLI("aws-topology-api", Version)
	c.Args = os.Args[1:]
	metaPtr := &command.Meta{
		Ui: &cli.BasicUi{
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		},
		Version: Version,
	}
	meta := *metaPtr

	c.Commands = map[string]cli.CommandFactory{
		"serve": func() (cli.Command, error) {
			return &command.ServeCommand{
				Meta: meta,
			}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		fmt.Println(err)
	}

	os.Exit(exitStatus)
}
