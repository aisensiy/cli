package cmd

import (
	"fmt"
	"github.com/sjkyspa/stacks/client/config"
	"github.com/sjkyspa/stacks/controller/api/api"
	"github.com/sjkyspa/stacks/controller/api/net"
	launcherApi "github.com/sjkyspa/stacks/launcher/api/api"
	deploymentNet "github.com/sjkyspa/stacks/launcher/api/net"
	"io/ioutil"
	"os/exec"
	"os"
	"errors"
)

func ScaffoldCreate(stackName, unifiedProcedure, provider, owner string, directory string, appName string, needDeploy string) error {

	if appName != "" {
		if isApplicationExist(appName) {
			return fmt.Errorf("Application %s already exists.", appName)
		}
		if directory == "" {
			directory = appName
		}
	}

	if directory == "" {
		directory = stackName
	}
	currentDir, _ := os.Getwd()
	target := fmt.Sprintf("%s//%s", currentDir, directory)
	if IsDirectoryExist(target) {
		return fmt.Errorf("directory %s already exists", directory);
	}

	if stackName != "" {
		stack, err := getStack(stackName)
		if err != nil {

			return err
		}
		gitRepo := stack.GetTemplateCode()

		if gitRepo == "" {
			return fmt.Errorf("git repositry is no valid, please check the definition of stack '%s' to make sure it contains valid template code.", stackName)
		}

		cmdString := fmt.Sprintf("git clone %s %s; cd %s; git remote remove origin; rm -rf .git; git init", gitRepo, directory, directory)

		err = ExecuteCmd(cmdString)

		if err != nil {
			return err
		}

		if appName != "" {
			os.Chdir(target)
			err = AppCreate(appName, stackName, unifiedProcedure, provider, owner, needDeploy)
			if err != nil {
				return err
			}
		}
		return nil
	} else if unifiedProcedure != "" {
		configRepository := config.NewConfigRepository(func(error) {})
		gateway := deploymentNet.NewCloudControllerGateway(configRepository)
		ups := launcherApi.NewUpsRepository(configRepository, gateway)

		unifiedProcedures, err := ups.GetUPByName(unifiedProcedure)
		if err != nil {
			return err
		}

		if unifiedProcedures.Count() != 1 {
			return fmt.Errorf("can not find the unified procedure by name given")
		}

		template := unifiedProcedures.Items()[0].Template()
		if template.URI == "" {
			return fmt.Errorf("git repositry is no valid, please check the definition of stack '%s' to make sure it contains valid template code.", stackName)
		}

		cmdString := fmt.Sprintf("git clone %s %s; cd %s; git remote remove origin; rm -rf .git; git init", template.URI, directory, directory)

		err = ExecuteCmd(cmdString)

		if err != nil {
			return err
		}

		if appName != "" {
			os.Chdir(target)
			err = AppCreate(appName, stackName, unifiedProcedure, provider, owner, needDeploy)
			if err != nil {
				return err
			}
		}
		return nil
	} else {
		return errors.New("please use -s to specify stack or use the -p to specify unified procedure")
	}
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

func IsDirectoryExist(directory string) (bool) {
	if _, notExistErr := os.Stat(directory); notExistErr != nil {
		return false
	}

	return true
}

func isApplicationExist(appName string) bool {
	configRepository := config.NewConfigRepository(func(error) {})
	appRepository := api.NewAppRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))
	app, _ := appRepository.GetApp(appName)
	return app != nil
}
