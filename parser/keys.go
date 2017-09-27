package parser

import (
	"fmt"
	docopt "github.com/docopt/docopt-go"
	"github.com/sjkyspa/stacks/client/cmd"
	"github.com/urfave/cli"
)

func KeysCommands() cli.Command {
	return cli.Command {
		Name: "keys",
		Usage: "Keys Commands",
		Subcommands: []cli.Command {
			{
				Name: "list",
				Usage: "List SSH keys for the logged in user",
				ArgsUsage: " ",
				Action: func(c *cli.Context) error{
					if err := cmd.ListKeys(); err != nil {
						return cli.NewExitError(err, 1)
					}
					return nil
				},
			},
			{
				Name: "add",
				Usage: "Add an SSH key",
				ArgsUsage: "[ssh-file-path]",
				Action: func(c *cli.Context) error{
					if c.Args().Get(0) == "" {
						return cli.NewExitError(fmt.Sprintf("USAGE: %s %s", c.Command.HelpName, c.Command.ArgsUsage), 1)
					}
					if err := cmd.AddKey(c.Args().Get(0)); err != nil {
						return cli.NewExitError(err, 1)
					}
					return nil
				},
			},
			{
				Name: "remove",
				Usage: "Remove an SSH key",
				ArgsUsage: "[key-id]",
				Action: func(c *cli.Context) error{
					if c.Args().Get(0) == "" {
						return cli.NewExitError(fmt.Sprintf("USAGE: %s %s", c.Command.HelpName, c.Command.ArgsUsage), 1)
					}
					if err := cmd.RemoveKey(c.Args().Get(0)); err != nil {
						return cli.NewExitError(err, 1)
					}
					return nil
				},
			},
		},
	}
}

func Keys(argv []string) error {
	usage := `
Valid commands for keys:

keys:list        list SSH keys for the logged in user
keys:add         add an SSH key
keys:remove      remove an SSH key

Use 'cde help [command]' to learn more.
`
	switch argv[0] {
	case "keys:list":
		return keyList(argv)
	case "keys:add":
		return addKey(argv)
	case "keys:remove":
		return removeKey(argv)
	case "keys":
		fmt.Print(usage)
	default:
		if printHelp(argv, usage) {
			return nil
		}
		PrintUsage()
		return nil
	}
	return nil
}

func keyList(argv []string) error {
	usage := `
List keys.

Usage: cde keys:list

`
	_, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmd.ListKeys()
}

func addKey(argv []string) error {
	usage := `
Add a key.

Usage: cde keys:add <ssh>

Arguments:
  <ssh>
  	the ssh file path
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	sshFilePath := safeGetValue(args, "<ssh>")
	return cmd.AddKey(sshFilePath)
}

func removeKey(argv []string) error {
	usage := `
Remove a key.

Usage: cde keys:remove <key>

Arguments:
  <key>
  	the key id
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	keyId := safeGetValue(args, "<key>")
	return cmd.RemoveKey(keyId)
}
