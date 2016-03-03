package parser

import (
	docopt "github.com/sjkyspa/stacks/Godeps/_workspace/src/github.com/docopt/docopt-go"
	"github.com/sjkyspa/stacks/client/config"
	"github.com/sjkyspa/stacks/apisdk/api"
	"github.com/sjkyspa/stacks/apisdk/net"
	"os/exec"
	"fmt"
	"regexp"
	"io/ioutil"
)

func Scaffold(argv []string) error {
	usage := `
Creates a new application.

Usage: cde scaffold <stack> [options]

Arguments:
  <stack>
  	a stack name

Options:
  -d --dir=<dir> default sub directory name
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	stack := safeGetValue(args, "<stack>")
	dir := safeGetOrDefault(args, "--dir", "")
	configRepository := config.NewConfigRepository(func(err error) {})
	stackRepository := api.NewStackRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))
	stacks, err := stackRepository.GetStackByName(stack)
	if (err != nil || stacks.Count() == 0) {
		return fmt.Errorf("stack not found")
	}
	stackId := stacks.Items()[0].Id();
	stackObj, err := stackRepository.GetStack(stackId)
	gitRepo := stackObj.GetTemplateCode()
	if(dir == ""){
		dir = retrieveGitName(gitRepo)
		if(dir == "") {
			return fmt.Errorf("directory name invalid")
		}
	}
	cmdString := fmt.Sprintf("git clone %s %s; cd %s; git remote remove origin", gitRepo, dir, dir)
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

func retrieveGitName(gitName string) string {
	regex := regexp.MustCompile(`^.*/(.+).git$`)
	if regex.MatchString(gitName) {
		captures := regex.FindStringSubmatch(gitName)
		return captures[1];
	}
	return "";
}
