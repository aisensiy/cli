package parser

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/sjkyspa/stacks/client/cmd"
)

// Config routes config commands to their specific function.
func Dev(argv []string) error {
	usage := `
Valid commands for config:

dev:up        	 start up the dev env
dev:down         shutdown the dev env
dev:destroy

Use 'cde help [command]' to learn more.
`
	switch argv[0] {
	case "dev:up":
		return devUp(argv)
	case "dev:down":
		return devDown(argv)
	case "dev:destroy":
		return devDestroy(argv)
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

Usage: cde dev:up [options]
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

Usage: cde dev:down [options]

Options:
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

Usage: cde dev:destroy [options]

Options:
`

	_, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmd.DevDestroy()
}
