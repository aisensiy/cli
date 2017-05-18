package cmd

import (
	"fmt"
	"github.com/sjkyspa/stacks/client/config"
	"github.com/sjkyspa/stacks/controller/api/api"
	"github.com/sjkyspa/stacks/controller/api/net"
	"io/ioutil"
	"os/exec"
)

func ScaffoldCreate(stackName string, directory string) error {
	if directory == "" {
		directory = stackName
	}

	stack, err := getStack(stackName)
	if err != nil {
		return err
	}

	gitRepo := stack.GetTemplateCode()

	if gitRepo == "" {
		return fmt.Errorf("git repositry is no valid, please check the definition of stack '%s' to make sure it contains valid template code.", stackName)
	}

	cmdString := fmt.Sprintf("git clone %s %s; cd %s; git remote remove origin; rm -rf .git; git init", gitRepo, directory, directory)

	return ExecuteCmd(cmdString)
}

func ExecuteCmd(cmdString string) error {
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

func getStack(stackName string) (stackObj api.Stack, err error) {
	configRepository := config.NewConfigRepository(func(err error) {})
	stackRepository := api.NewStackRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))
	stacks, err := stackRepository.GetStackByName(stackName)

	if err != nil || stacks.Count() == 0 {
		err = fmt.Errorf("stack not found")
		return
	}

	stackId := stacks.Items()[0].Id()
	stackObj, err = stackRepository.GetStack(stackId)
	return
}