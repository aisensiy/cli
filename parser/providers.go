package parser

import (
	"errors"
	"github.com/sjkyspa/stacks/client/cmd"
	"strings"
	cli "gopkg.in/urfave/cli.v2"
	"fmt"
)

func ProvidersCommand() *cli.Command {
	return &cli.Command{
		Name:  "providers",
		Usage: "Providers Commands",
		Subcommands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List all Providers",
				ArgsUsage: " ",
				Action: func(c *cli.Context) error {
					err := cmd.ProviderList()
					if err != nil {
						return cli.Exit(fmt.Sprintf("%v", err), 1)
					}
					return nil
				},
			},
			{
				Name:      "info",
				Usage:     "Get info of a Provider",
				ArgsUsage: "<provider-name>",
				Action: func(c *cli.Context) error {
					if (c.Args().Get(0) == "") {
						return cli.Exit(fmt.Sprintf("USAGE: %s %s", c.Command.HelpName, c.Command.ArgsUsage), 1)
					}
					err := cmd.GetProviderByName(c.Args().Get(0))
					if err != nil {
						return cli.Exit(fmt.Sprintf("%v", err), 1)
					}
					return nil
				},
			},
			{
				Name:      "enroll",
				Usage:     "Enroll a new Provider",
				ArgsUsage: "<name] <type>",
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:  "config, c",
						Usage: "Set provider's configuration. Key \"endpoint\" is required. (Tips: String \"$$\" needs to be escaped in shell).",
					},
					&cli.StringFlag{
						Name:  "for, f",
						Usage: "Specify an organization for the provider.",
					},
				},
				Action: func(c *cli.Context) error {
					if c.Args().Get(0) == "" || c.Args().Get(1) == "" {
						return cli.Exit(fmt.Sprintf("USAGE: %s %s", c.Command.HelpName, c.Command.ArgsUsage), 1)
					}
					configMap, err := enrollConfigConvert(c.StringSlice("config"))
					if err != nil {
						return cli.Exit(fmt.Sprintf("%v", err), 1)
					}
					consumer := c.String("for")
					err = cmd.ProviderCreate(c.Args().Get(0), c.Args().Get(1), consumer, configMap)
					if err != nil {
						return cli.Exit(fmt.Sprintf("%v", err), 1)
					}
					return nil
				},
			},
			{
				Name:      "update",
				Usage:     "Update an existing Provider",
				ArgsUsage: "<name>",
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:  "config, c",
						Usage: "Set provider's configuration. Key \"endpoint\" is required.",
					},
				},
				Action: func(c *cli.Context) error {
					if c.Args().Get(0) == "" {
						return cli.Exit(fmt.Sprintf("USAGE: %s %s", c.Command.HelpName, c.Command.ArgsUsage), 1)
					}
					configMap, err := updateConfigConvert(c.StringSlice("config"))
					if err != nil {
						return cli.Exit(fmt.Sprintf("%v", err), 1)
					}
					err = cmd.ProviderUpdate(c.Args().Get(0), configMap)
					if err != nil {
						return cli.Exit(fmt.Sprintf("%v", err), 1)
					}
					return nil
				},
			},
		},
	}
}

func updateConfigConvert(config []string) (map[string]interface{}, error) {
	if len(config) == 0 {
		return nil, errors.New("please input at least one config")
	}

	configMap := map[string]interface{}{}
	for _, v := range config {
		pair := strings.Split(v, "=")
		if len(pair) != 2 {
			pair = append(pair, "")
		}
		configMap[pair[0]] = pair[1]
	}
	return configMap, nil
}

func enrollConfigConvert(config []string) (map[string]interface{}, error) {
	if len(config) == 0 {
		return nil, errors.New("please input at least one config")
	}

	configMap := map[string]interface{}{}
	for _, v := range config {
		pair := strings.Split(v, "=")
		if len(pair) != 2 {
			return nil, errors.New("invalid config format")
		}
		if pair[1] == "" {
			return nil, errors.New("config value should not be empty")
		}
		configMap[pair[0]] = pair[1]
	}
	return configMap, nil
}
