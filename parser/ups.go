package parser

import (
	"github.com/sjkyspa/stacks/client/cmd"
	cli "gopkg.in/urfave/cli.v2"
)

func UpsCommand() *cli.Command {
	return &cli.Command{
		Name:  "ups",
		Usage: "Unified Procedures Commands",
		Subcommands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List all Unified Procedures",
				ArgsUsage: " ",
				Action: func(c *cli.Context) error {
					return cmd.UpsList()
				},
			},
			{
				Name:      "info",
				Usage:     "Get info of an Unified Procedure",
				ArgsUsage: "<up-name>",
				Action: func(c *cli.Context) error {
					return cmd.UpsInfo(c.Args().First())
				},
			},
			{
				Name:      "draft",
				Usage:     "Create a new Unified Procedure",
				ArgsUsage: "<up-file>",
				Action: func(c *cli.Context) error {
					return cmd.UpCreate(c.Args().First())
				},
			},
			{
				Name:      "update",
				Usage:     "Update an existing Unified Procedure",
				ArgsUsage: "<up-id> <up-file>",
				Action: func(c *cli.Context) error {
					return cmd.UpUpdate(c.Args().First(), c.Args().Get(1))
				},
			},
			{
				Name:      "remove",
				Usage:     "Delete an Unified Procedure",
				ArgsUsage: "<up-id>",
				Action: func(c *cli.Context) error {
					return cmd.UpRemove(c.Args().First())
				},
			},
		},
	}
}
