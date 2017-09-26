package parser

import (
	"fmt"
	docopt "github.com/docopt/docopt-go"
	"github.com/sjkyspa/stacks/client/cmd"
	"regexp"
	"github.com/urfave/cli"
)

func ScaffoldCommands() cli.Command {
	return cli.Command{
		Name: "scaffold",
		Usage: "Creates a new scaffold in current directory.",
		ArgsUsage: "[]",
		Flags: []cli.Flag {
			cli.StringFlag{
				Name: "stack, s",
				Usage: "Specify stack with name",
			},
			cli.StringFlag{
				Name: "unified_procedure, u",
				Usage: "Specify unified procedure with name",
			},
			cli.StringFlag{
				Name: "provider, p",
				Usage: "Specify provider for provide the app runtime",
			},
			cli.StringFlag{
				Name: "owner, o",
				Usage: "Specify owner for the app",
			},
			cli.StringFlag{
				Name: "dir, d",
				Usage: "Specify default sub directory with name",
			},
			cli.StringFlag{
				Name: "app, a",
				Usage: "Create a new scaffold and create a new app in sub directory",
			},
			cli.StringFlag{
				Name: "deploy",
				Usage: "Tell system to deploy this app or not, 1 means need, 0 mean no, default 1",
			},
		},
		Action: func(c *cli.Context) error{
			appName := c.String("app")
			if appName!="" && !cmd.IsAppNameInvalid(appName) {
				return cli.NewExitError(fmt.Errorf("'%s' does not match the pattern '[a-z0-9-]+'\n", appName), 1)
			}

			needDeploy := c.String("deploy")
			if needDeploy == "" {
				needDeploy = "1"
			}

			if err := cmd.ScaffoldCreate(c.String("stack"), c.String("unified_procedure"),  c.String("provider"), c.String("owner"),  c.String("dir"), appName, needDeploy); err != nil {
				return cli.NewExitError(err, 1)
			}

			return nil
		},
	}
}

func Scaffold(argv []string) error {
	usage := `
Creates a new scaffold in current directory.

Usage: cde scaffold [options]

Arguments:
Options:
  -s --stack=<stackName>
  	a stack name
  -u --unified_procedure=<unified_procedure>
	a unified procedure name
  -p --provider=<unified_procedure>
	the provider for provide the app runtime
  -o --owner=<owner>
	the owner for the app
  -d --dir=<dir>
	default sub directory name
  -a --app=<app-name>
	create a new scaffold and create a new app in sub directory
  --deploy tell system to deploy this app or not, 1 means need, 0 mean no, default 1
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	unifiedProcedure := safeGetValue(args, "--unified_procedure")
	stackName := safeGetValue(args, "--stack")

	dir := safeGetOrDefault(args, "--dir", "")

	appName := safeGetOrDefault(args, "--app", "")
	provider := safeGetOrDefault(args, "--provider", "")
	owner := safeGetOrDefault(args, "--owner", "")

	needDeploy := safeGetOrDefault(args, "--deploy", "1")

	if appName!="" && !cmd.IsAppNameInvalid(appName) {
		return fmt.Errorf("'%s' does not match the pattern '[a-z0-9-]+'\n", appName)
	}

	return cmd.ScaffoldCreate(stackName, unifiedProcedure,  provider, owner,  dir, appName, needDeploy)
}

func retrieveGitName(gitUrl string) (string, error) {
	regex := regexp.MustCompile(`^.*/(.+).git$`)
	if regex.MatchString(gitUrl) {
		captures := regex.FindStringSubmatch(gitUrl)
		return captures[1], nil
	}
	return "", fmt.Errorf("Invalid Git URL")
}
