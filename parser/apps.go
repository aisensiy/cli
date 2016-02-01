package parser

import (
	"errors"
	"fmt"
	docopt "github.com/sjkyspa/stacks/Godeps/_workspace/src/github.com/docopt/docopt-go"
	"github.com/sjkyspa/stacks/client/cmd"
	"os"
	"strconv"
)

func Apps(argv []string) error {
	usage := `
Valid commands for apps:

apps:create        create a new application
apps:list          list accessible applications
apps:info          view info about an application
apps:destroy       destroy an application and stop application instance in deployment environment
apps:stack-update       change to use another stack

Use 'cde help [command]' to learn more.
`
	switch argv[0] {
	case "apps:create":
		return appCreate(argv)
	case "apps:list":
		return appList()
	case "apps:info":
		return appInfo(argv)
	case "apps:destroy":
		return appDestroy(argv)
	case "apps:stack-update":
		return appStackUpdate(argv)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "apps" {
			argv[0] = "apps:list"
			return appList()
		}

		PrintUsage()
		return nil
	}
	return nil
}

func appCreate(argv []string) error {
	usage := `
Creates a new application.

Usage: cde apps:create <name> <stack> [options]

Arguments:
  <name>
  	a uniquely identifiable name for the application. No other app can already
    exist with this name.
  <stack>
  	a stack name

Options:
  --mem=<mem>
  	allocated memory for this app. [default: 512]
  --disk=<disk>
  	max allocated disk size. [default: 20]
  --instances=<instances>
  	default started instance number. [default: 1]
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	name := safeGetValue(args, "<name>")
	stack := safeGetValue(args, "<stack>")
	if stack == "" || name == "" {
		return errors.New("<name> <stack> are essential parameters")
	}
	memory := safeGetOrDefault(args, "--mem", "512")
	disk := safeGetOrDefault(args, "--disk", "20")
	instances := safeGetOrDefault(args, "--instances", "1")

	var mem, ins, diskSize int

	if mem, err = strconv.Atoi(memory); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}

	if ins, err = strconv.Atoi(instances); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}

	if diskSize, err = strconv.Atoi(disk); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}
	return cmd.AppCreate(name, stack, mem, diskSize, ins)
}

func appList() error {
	return cmd.AppsList()
}

func appInfo(argv []string) error {
	usage := `
Prints info about the current application.

Usage: cde apps:info [options]

Options:
  -a --app=<app>
    the uniquely identifiable id for the application.
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	appId := safeGetValue(args, "--app")

	return cmd.GetApp(appId)

}

func appDestroy(argv []string) error {
	usage := `
Destroy an application and stop application instance in deployment environment.
Usage: cde apps:destroy [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	appId := safeGetValue(args, "--app")

	return cmd.DestroyApp(appId)

}

func appStackUpdate(argv []string) error {
	usage := `
Change to use another stack.
Usage: cde apps:stack-update [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  -s --stack=<stack>
    another existing stack name.
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	appName := safeGetValue(args, "--app")
	stackName := safeGetValue(args, "--stack")

	return cmd.SwitchStack(appName, stackName)

}
