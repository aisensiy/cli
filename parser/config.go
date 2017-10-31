package parser

import (
	"fmt"
	docopt "github.com/docopt/docopt-go"
	"github.com/cnupp/cli/cmd"
	cli "gopkg.in/urfave/cli.v2"
)

func ConfigCommands() *cli.Command {
	return &cli.Command{
		Name:  "config",
		Usage: "Config Commands",
		Subcommands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List environment variables for an app.",
				ArgsUsage: " ",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "oneline",
						Usage: "Print output on one line",
					},
					&cli.StringFlag{
						Name:  "app, a",
						Usage: "Specify app with name",
					},
				},
				Action: func(c *cli.Context) error {
					if err := cmd.ConfigList(c.String("app"), c.Bool("oneline")); err != nil {
						return cli.Exit(fmt.Sprintf("%v", err), 1)
					}
					return nil
				},
			},
			{
				Name:      "set",
				Usage:     "Set environment variables for an app",
				ArgsUsage: "<key>=<value> [<key>=<value>...]",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "app, a",
						Usage: "Specify app with name",
					},
				},
				Action: func(c *cli.Context) error {
					if !c.Args().Present() {
						return cli.Exit(fmt.Sprintf("USAGE: %s %s", c.Command.HelpName, c.Command.ArgsUsage), 1)
					}
					envs := append(c.Args().Tail(), c.Args().First())
					if err := cmd.ConfigSet(c.String("app"), envs); err != nil {
						return cli.Exit(fmt.Sprintf("%v", err), 1)
					}
					return nil
				},
			},
			{
				Name:      "unset",
				Usage:     "Unset environment variables for an app",
				ArgsUsage: "<key> [<key>...]",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "app, a",
						Usage: "Specify app with name",
					},
				},
				Action: func(c *cli.Context) error {
					if !c.Args().Present() {
						return cli.Exit(fmt.Sprintf("USAGE: %s %s", c.Command.HelpName, c.Command.ArgsUsage), 1)
					}
					keys := append(c.Args().Tail(), c.Args().First())
					if err := cmd.ConfigUnset(c.String("app"), keys); err != nil {
						return cli.Exit(fmt.Sprintf("%v", err), 1)
					}
					return nil
				},
			},
		},
	}
}

// Config routes config commands to their specific function.
func Config(argv []string) error {
	usage := `
Valid commands for config:

config:list        list environment variables for an app
config:set         set environment variables for an app
config:unset       unset environment variables for an app

Use 'cde help [command]' to learn more.
`

	switch argv[0] {
	case "config:list":
		return configList(argv)
	case "config:set":
		return configSet(argv)
	case "config:unset":
		return configUnset(argv)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "config" {
			argv[0] = "config:list"
			return configList(argv)
		}

		PrintUsage()
		return nil
	}
}

func configList(argv []string) error {
	usage := `
Lists environment variables for an application.

Usage: cde config:list [options]

Options:
  --oneline
    print output on one line.
  -a --app=<app>
    the uniquely identifiable name of the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmd.ConfigList(safeGetValue(args, "--app"), args["--oneline"].(bool))
}

func configSet(argv []string) error {
	usage := `
Sets environment variables for an application.

Usage: cde config:set <key>=<value> [<key>=<value>...] [options]

Arguments:
  <key>
    the uniquely identifiable name for the environment variable.
  <value>
    the value of said environment variable.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmd.ConfigSet(safeGetValue(args, "--app"), args["<key>=<value>"].([]string))
}

func configUnset(argv []string) error {
	usage := `
Unsets an environment variable for an application.

Usage: cde config:unset <key> [<key>...] [options]

Arguments:
  <key>
    the variable to remove from the application's environment.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmd.ConfigUnset(safeGetValue(args, "--app"), args["<key>"].([]string))
}
