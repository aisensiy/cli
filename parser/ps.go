package parser

import (
	"fmt"
	docopt "github.com/sjkyspa/stacks/Godeps/_workspace/src/github.com/docopt/docopt-go"
	"github.com/sjkyspa/stacks/client/cmd"
)

func Ps(argv []string) error {
	usage := `
Valid commands for keys:

ps:restart	restart a service process (without restarting dependent services)
ps:scale    scale a service process
ps:list     list all dependent services for an application(service)

Use 'cde help [command]' to learn more.
`
	switch argv[0] {
	case "ps:restart":
		return processRestart(argv)
	case "ps:scale":
		return processScale(argv)
	case "ps:list":
		return listDependentServices(argv)
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

func processRestart(argv []string) error {
	usage := `
Restart a service process (without restarting dependent services).
Usage: cde ps:restart <app>

Arguments:
  <app>
  	the application name
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}
	appName := safeGetValue(args, "<app>")

	return cmd.RestartApp(appName)
}

func processScale(argv []string) error {
	return nil
}

func listDependentServices(argv []string) error {
	return nil
}