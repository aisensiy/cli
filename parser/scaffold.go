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
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	stackName := safeGetValue(args, "<stackName>")

	dir := safeGetOrDefault(args, "--dir", "")

	return cmd.ScaffoldCreate(stackName, dir)
}

func retrieveGitName(gitUrl string) (string, error) {
	regex := regexp.MustCompile(`^.*/(.+).git$`)
	if regex.MatchString(gitUrl) {
		captures := regex.FindStringSubmatch(gitUrl)
		return captures[1], nil
	}
	return "", fmt.Errorf("Invalid Git URL")
}
