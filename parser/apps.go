package parser

import (
	"errors"
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/sjkyspa/stacks/client/cmd"
	"os"
	"strconv"
)

func Apps(argv []string) error {
	usage := `
Valid commands for apps:

apps:create             create a new application
apps:list               list accessible applications
apps:info               view info about an application
apps:destroy            destroy an application and stop application instance in deployment environment
apps:stack-update       change to use another stack
apps:logs               view logs
apps:collaborators     	view collaborators
apps:add-collaborator	add collaborator
apps:rm-collaborator    remove collaborator
apps:transfer           transfer app to others, user or organization
apps:launch             launch non build app
apps:localization	get codebase for an app

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
	case "apps:collaborators":
		return appCollaborators(argv)
	case "apps:add-collaborator":
		return appAddCollaborator(argv)
	case "apps:rm-collaborator":
		return appRmCollaborator(argv)
	case "apps:transfer":
		return appTransfer(argv)
	case "apps:launch":
		return appLaunch(argv)
	case "apps:localization":
		return appLocalization(argv)
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
func appLocalization(argv []string) error {
	usage := `
Get codebase for an app in sub directory.

Usage: cde apps:localization <name> [options]

Arguments:
  <name>
  	the uniquely identifiable name for the application

Options:
  -d --directory=<directory>
  	default sub directory name
	`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	appName := safeGetValue(args, "<name>")
	directory := safeGetValue(args, "--directory")

	if appName == "" {
		return errors.New("<name> are essential parameters")
	}

	return cmd.AppLocalization(appName, directory)
}

func appCreate(argv []string) error {
	usage := `
Creates a new application.

Usage: cde apps:create <name> [options]

Arguments:
  <name>
  	a uniquely identifiable name for the application. No other app can already
    exist with this name.
Options:
  -d --deploy=<deploy>
    tell system to deploy this app or not, 1 means need, 0 mean no, default 1
  -s --stack=<stack>
  	a stack name
  -p --unified_procedure=<unified_procedure>
	a unified procedure name
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	name := safeGetValue(args, "<name>")
	stack := safeGetValue(args, "--stack")
	unifiedProcedure := safeGetValue(args, "--unified_procedure")
	needDeploy := safeGetOrDefault(args, "--deploy", "1")

	if (stack == "" && unifiedProcedure == "") || name == "" {
		return errors.New("<name> <stack> are essential parameters")
	}

	if !cmd.IsAppNameInvalid(name) {
		return fmt.Errorf("'%s' does not match the pattern '[a-z0-9-]+'\n", name)
	}

	return cmd.AppCreate(name, stack, unifiedProcedure, needDeploy)
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
    the uniquely identifiable name for the application.
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

func appCollaborators(argv []string) error {
	usage := `
Prints collaborators in app

Usage: cde apps:collaborators [options]

Options:
  -a --app=<app>
    the uniquely identifiable id for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	appId := safeGetValue(args, "--app")
	return cmd.AppCollaborators(appId)
}

func appAddCollaborator(argv []string) error {
	usage := `
Add collaborator for app

Usage: cde apps:add-collaborator <email> [options]

Arguments:
  <email>
    email of collaborator

Options:
  -a --app=<app>
    the uniquely identifiable id for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	email := safeGetValue(args, "<email>")
	if email == "" {
		return errors.New("<email> is essential parameters")
	}

	appId := safeGetValue(args, "--app")

	return cmd.AppAddCollaborator(appId, email)
}

func appRmCollaborator(argv []string) error {
	usage := `
Remove collaborator for app

Usage: cde apps:rm-collaborator <email> [options]

Arguments:
  <email>
    email of collaborator

Options:
  -a --app=<app>
    the uniquely identifiable id for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	email := safeGetValue(args, "<email>")
	if email == "" {
		return errors.New("<email> is essential parameters")
	}

	appId := safeGetValue(args, "--app")

	return cmd.AppRmCollaborator(appId, email)
}

func appTransfer(argv []string) error {
	usage := `
Transfer app to others, user or organization

Usage: cde apps:transfer [options]

Options:
  -a --app=<app>
    the uniquely identifiable id for the application.
  -o --org=<org>
    name of the organization transfer to
  -e --email=<email>
    email of the user transfer to
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	appId := safeGetValue(args, "--app")
	email := safeGetValue(args, "--email")
	org := safeGetValue(args, "--org")
	return cmd.AppTransfer(appId, email, org)
}

func appLaunch(argv []string) error {
	usage := `
Launch non build app

Usage: cde apps:launch [options]

Options:
  -a --app=<app>
    the uniquely identifiable id for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	appId := safeGetValue(args, "--app")
	return cmd.AppLaunch(appId)
}
