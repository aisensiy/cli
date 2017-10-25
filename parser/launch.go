package parser

import (
	"fmt"

	"github.com/docopt/docopt-go"
	"github.com/sjkyspa/stacks/client/cmd"
	"gopkg.in/urfave/cli.v2"
)

func LaunchCommands() *cli.Command {
	return &cli.Command{
		Name:  "launch",
		Usage: "Launch Commands",
		Subcommands: []*cli.Command{
			{
				Name:      "build",
				Usage:     "Launch a build procedure.",
				ArgsUsage: "<filename> <app-name>",
				Action: func(c *cli.Context) error {
					filename := c.Args().Get(0)
					appName := c.Args().Get(1)
					if filename == "" || appName == "" {
						return cli.Exit(fmt.Sprintf("USAGE: %s %s", c.Command.HelpName, c.Command.ArgsUsage), 1)
					}
					if err := cmd.LaunchBuild(filename, appName); err != nil {
						return cli.Exit(fmt.Sprintf("%v", err), 1)
					}
					return nil
				},
			},
		},
	}
}

func Launch(argv []string) error {
	usage := `
Valid commands for launch:

launch:build  launch a build procedure

Use 'cde help [command]' to learn more.
`
	switch argv[0] {
	case "launch:build":
		return launchBuild(argv)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		PrintUsage()
		return nil
	}
	return nil
}

func launchBuild(argv []string) error {
	usage := `
Launch a build procedure.

Usage: cde launch:build (-f <filename>) (-a <app-name>)

Arguments:
  <filename>
  the code base to build with
  <app-name>
  the app to build with
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}
	filename := safeGetValue(args, "<filename>")
	appName := safeGetValue(args, "<app-name>")

	return cmd.LaunchBuild(filename, appName)
}
