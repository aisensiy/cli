package parser

import (
	"fmt"
	docopt "github.com/docopt/docopt-go"
	"github.com/sjkyspa/stacks/client/cmd"
	"regexp"
)

func Scaffold(argv []string) error {
	usage := `
Creates a new scaffold in current directory.

Usage: cde scaffold <stackName> [options]

Arguments:
  <stackName>
  	a stack name

Options:
  -d --dir=<dir> default sub directory name
  -a --app=<app-name> create a new scaffold and create a new app in sub directory
  --deploy tell system to deploy this app or not, 1 means need, 0 mean no, default 1
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	stackName := safeGetValue(args, "<stackName>")

	dir := safeGetOrDefault(args, "--dir", "")

	appName := safeGetOrDefault(args, "--app", "")

	needDeploy := safeGetOrDefault(args, "--deploy", "1")

	if appName!="" && !cmd.IsAppNameInvalid(appName) {
		return fmt.Errorf("'%s' does not match the pattern '[a-z0-9-]+'\n", appName)
	}

	return cmd.ScaffoldCreate(stackName, dir, appName, needDeploy)
}

func retrieveGitName(gitUrl string) (string, error) {
	regex := regexp.MustCompile(`^.*/(.+).git$`)
	if regex.MatchString(gitUrl) {
		captures := regex.FindStringSubmatch(gitUrl)
		return captures[1], nil
	}
	return "", fmt.Errorf("Invalid Git URL")
}
