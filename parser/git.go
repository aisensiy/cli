package parser

import (
	"fmt"

	docopt "github.com/docopt/docopt-go"
	"github.com/sjkyspa/stacks/client/cmd"
	"github.com/urfave/cli"
)

func GitCommands() cli.Command {
	return cli.Command {
		Name: "git",
		Usage: "Git Commands",
		Subcommands: []cli.Command {
			{
				Name: "remote",
				Usage: "Adds git remote of application to repository.",
				ArgsUsage: " ",
				Flags: []cli.Flag {
					cli.StringFlag {
						Name: "app, a",
						Usage: "Specify app with name",
					},
					cli.StringFlag {
						Name: "remote, r",
						Usage: "Name of remote to create. [default: cde]",
					},
				},
				Action: func(c *cli.Context) error {
					remote := c.String("remote")
					if remote == "" {
						remote = "cde"
					}
					if err := cmd.GitRemote(c.String("app"), remote); err != nil {
						return cli.NewExitError(err, 1)
					}
					return nil
				},
			},
		},
	}
}

// Git routes git commands to their specific function.
func Git(argv []string) error {
	usage := `
Valid commands for git:

git:remote          Adds git remote of application to repository

Use 'deis help [command]' to learn more.
`

	switch argv[0] {
	case "git:remote":
		return gitRemote(argv)
	case "git":
		fmt.Print(usage)
		return nil
	default:
		PrintUsage()
		return nil
	}
}

func gitRemote(argv []string) error {
	usage := `
Adds git remote of application to repository

Usage: cde git:remote <app> [options]

Options:
  app
    the uniquely identifiable name for the application.
  -r --remote=REMOTE
    name of remote to create. [default: cde]
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	err = cmd.GitRemote(safeGetValue(args, "<app>"), args["--remote"].(string))
	//if(err != nil ){
	//	if strings.Contains(err.Error(), "exit status 128") {
	//		fmt.Println("Please use another remote name with `cde git:remote <app> -r <remote name>`")
	//	}
	//}
	return err
}
