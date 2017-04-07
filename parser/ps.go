package parser

import (
	"fmt"
	docopt "github.com/docopt/docopt-go"
	"github.com/sjkyspa/stacks/client/cmd"
	"github.com/sjkyspa/stacks/launcher/api/api"
	"strconv"
)

func Ps(argv []string) error {
	usage := `
Valid commands for keys:

ps:restart	restart a service process (without restarting dependent services)
ps:scale	scale a service process
ps:list		list all dependent services for an application(service)

Use 'cde help [command]' to learn more.
`
	switch argv[0] {
	case "ps:restart":
		return processRestart(argv)
	case "ps:scale":
		return processScale(argv)
	case "ps:list":
		return listDependentServices(argv)
	case "ps":
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
Usage: cde ps:restart [options]

Options:
  -a --app=<app>
    the uniquely identifiable id for the application.
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}
	appId := safeGetValue(args, "--app")

	return cmd.RestartApp(appId)
}

func processScale(argv []string) error {
	usage := `
Scale a service for an application.
Usage: cde ps:scale <service-name> <num> [options]

Arguments:
  <service-name>
  	the service name
  <num>
    instance count

Options:
  -a --app=<app>
    the uniquely identifiable id for the application.
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}
	appId := safeGetValue(args, "--app")
	serviceName := safeGetValue(args, "<service-name>")
	instances := safeGetValue(args, "<num>")

	var instanceNum int = 0
	if instanceNum, err = strconv.Atoi(instances); err != nil {
		fmt.Sprintf("Error: %v\n", err)
		return err
	}

	originService, err := cmd.GetService(appId, serviceName)
	if err != nil {
		return err
	}

	params := api.ServiceConfigParams{
		Instance: instanceNum,
		CPUS:     originService.CPU(),
		Memory:   originService.Memory(),
	}

	return cmd.Scale(appId, serviceName, params)
}

func listDependentServices(argv []string) error {
	usage := `
List dependent services for an application.
Usage: cde ps:list [options]

Options:
  -a --app=<app>
    the uniquely identifiable id for the application.
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}
	appName := safeGetValue(args, "--app")

	return cmd.ListDependentServices(appName)
}
