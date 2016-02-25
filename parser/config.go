package parser

import (
	docopt "github.com/sjkyspa/stacks/Godeps/_workspace/src/github.com/docopt/docopt-go"
	"github.com/sjkyspa/stacks/client/cmd"
)

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

Usage: cde config:unset <key>

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

	return cmd.ConfigUnset(safeGetValue(args, "--app"), safeGetValue(args, "<key>"))
}
