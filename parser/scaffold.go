package parser

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/sjkyspa/stacks/client/config"
	"github.com/sjkyspa/stacks/controller/api/api"
	"github.com/sjkyspa/stacks/controller/api/net"
	"os/exec"
	"fmt"
	"regexp"
	"io/ioutil"
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
	if (dir == "") {
		dir = stackName
	}

	stack, err := getStack(stackName)
	if err != nil {
		return err
	}

	gitRepo := stack.GetTemplateCode()

	if (gitRepo == "") {
		return fmt.Errorf("git repositry is no valid, please check the definition of stack '%s' to make sure it contains valid template code.", stackName)
	}

	cmdString := fmt.Sprintf("git clone %s %s; cd %s; git remote remove origin; rm -rf .git; git init", gitRepo, dir, dir)

	return executeCmd(cmdString)
}

func getStack(stackName string) (stackObj api.Stack, err error) {
	configRepository := config.NewConfigRepository(func(err error) {})
	stackRepository := api.NewStackRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))
	stacks, err := stackRepository.GetStackByName(stackName)

	if (err != nil || stacks.Count() == 0) {
		err = fmt.Errorf("stack not found")
		return
	}

	stackId := stacks.Items()[0].Id();
	stackObj, err = stackRepository.GetStack(stackId)
	return
}

func executeCmd(cmdString string) error {
	cmd := exec.Command("/bin/sh", "-c", cmdString)
	stderr, err := cmd.StderrPipe()

	if err != nil {
		return err
	}

	if err = cmd.Start(); err != nil {
		return err
	}

	output, _ := ioutil.ReadAll(stderr)
	fmt.Print(string(output))

	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func retrieveGitName(gitUrl string) (string, error) {
	regex := regexp.MustCompile(`^.*/(.+).git$`)
	if regex.MatchString(gitUrl) {
		captures := regex.FindStringSubmatch(gitUrl)
		return captures[1], nil;
	}
	return "", fmt.Errorf("Invalid Git URL");
}
