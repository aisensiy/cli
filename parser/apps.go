package parser

import (
	"errors"
	"fmt"
	docopt "github.com/sjkyspa/stacks/Godeps/_workspace/src/github.com/docopt/docopt-go"
	"github.com/sjkyspa/stacks/client/cmd"
	"os"
	"strconv"
	"regexp"
)

func Apps(argv []string) error {
	usage := `
Valid commands for apps:

apps:create        	create a new application
apps:list          	list accessible applications
apps:info          	view info about an application
apps:destroy       	destroy an application and stop application instance in deployment environment
apps:stack-update	change to use another stack
apps:logs       	view logs

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
	case "apps:logs":
		return appLogs(argv)
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

	regex := regexp.MustCompile(`^[a-z0-9\-]+$`)
	if !regex.MatchString(name) {
		return fmt.Errorf("'%s' does not match the pattern '[a-z0-9-]+'\n", name)
	}

	return cmd.AppCreate(name, stack)
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

func appLogs(argv []string) error {
	usage := `
Prints info about the current application.

Usage: cde apps:logs [options]

Options:
  -a --app=<app>
    the uniquely identifiable id for the application.
  -n --lines=<lines>
    the number of lines to display
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	appId := safeGetValue(args, "--app")
	lines := safeGetOrDefault(args, "--lines", "100")
	var lineNum int
	if lineNum, err = strconv.Atoi(lines); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}

	return cmd.AppLog(appId, lineNum)
}
