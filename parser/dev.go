package parser

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/sjkyspa/stacks/client/cmd"
	"github.com/urfave/cli"
)

func DevCommands() cli.Command {
	return cli.Command{
		Name:  "dev",
		Usage: "Dev Commands",
		Subcommands: []cli.Command{
			{
				Name:      "up",
				Usage:     "Start up the local dev env.",
				ArgsUsage: " ",
				Action: func(c *cli.Context) error {
					if err := cmd.DevUp(); err != nil {
						return cli.NewExitError(err, 1)
					}
					return nil
				},
			},
			{
				Name:      "down",
				Usage:     "Shutdown the local dev env.",
				ArgsUsage: " ",
				Action: func(c *cli.Context) error {
					if err := cmd.DevDown(); err != nil {
						return cli.NewExitError(err, 1)
					}
					return nil
				},
			},
			{
				Name:      "destroy",
				Usage:     "Destroy the local dev env.",
				ArgsUsage: " ",
				Action: func(c *cli.Context) error {
					if err := cmd.DevDestroy(); err != nil {
						return cli.NewExitError(err, 1)
					}
					return nil
				},
			},
			{
				Name:      "env",
				Usage:     "Display the env variables.",
				ArgsUsage: " ",
				Action: func(c *cli.Context) error {
					if err := cmd.DevEnv(); err != nil {
						return cli.NewExitError(err, 1)
					}
					return nil
				},
			},
		},
	}
}

// Config routes config commands to their specific function.
func Dev(argv []string) error {
	usage := `
Valid commands for config:

dev:up        	 start up the dev env
dev:down         shutdown the dev env
dev:destroy      destroy the dev env
dev:env          display the env variables

Use 'cde help [command]' to learn more.
`
	switch argv[0] {
	case "dev:up":
		return devUp(argv)
	case "dev:down":
		return devDown(argv)
	case "dev:destroy":
		return devDestroy(argv)
	case "dev:env":
		return devEnv(argv)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "dev" {
			argv[0] = "dev:up"
			return devUp(argv)
		}

		PrintUsage()
		return nil
	}
}

func devUp(argv []string) error {
	usage := `
Start up the local dev env

Usage: cde dev:up
`

	_, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmd.DevUp()
}

func devDown(argv []string) error {
	usage := `
Shutdown the local dev env

Usage: cde dev:down
`

	_, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmd.DevDown()
}

func devDestroy(argv []string) error {
	usage := `
Destroy the local dev env

Usage: cde dev:destroy
`

	_, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmd.DevDestroy()
}

func devEnv(argv []string) error {
	usage := `
Display the env variables

Usage: cde dev:env
	`

	_, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmd.DevEnv()
}
