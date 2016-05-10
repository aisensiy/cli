package cmd

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/sjkyspa/stacks/controller/api/api"
	"github.com/sjkyspa/stacks/controller/api/net"
	"github.com/sjkyspa/stacks/client/config"
	"io/ioutil"
	"encoding/json"
)

func StackCreate(filename string) error {
	configRepository := config.NewConfigRepository(func(err error) {})
	stackRepository := api.NewStackRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))
	content, err := getStackFileContent(filename)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	content, err = yaml.YAMLToJSON(content)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	stackDefinition := make(map[string]interface{})
	if err := json.Unmarshal(content, &stackDefinition); err != nil {
		return err
	}

	stackModel, err := stackRepository.Create(stackDefinition)
	if err != nil {
		return err
	}
	fmt.Printf("create stack %s with uuid %s\n", stackModel.Name(), stackModel.Id())
	return nil
}

func getStackFileContent(filename string) (content []byte, err error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return []byte{}, err
	}
	return contents, err
}

func StacksList() error {
	configRepository := config.NewConfigRepository(func(error) {})
	stackRepository := api.NewStackRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))
	stacks, err := stackRepository.GetStacks()
	if err != nil {
		return err
	}
	fmt.Printf("=== Stacks: [%d]\n", len(stacks.Items()))

	for _, stack := range stacks.Items() {
		fmt.Printf("name: %s id: %s\n", stack.Name(), stack.Id())
	}
	return nil
}

func StackRemove(name string) error {
	configRepository := config.NewConfigRepository(func(err error) {})
	stackRepository := api.NewStackRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))
	stacks, err := stackRepository.GetStackByName(name)
	if err != nil {
		return err
	}
	stackId := stacks.Items()[0].Id()
	err = stackRepository.Delete(stackId)
	if err != nil {
		return err
	}
	fmt.Printf("delete stack successfully\n")
	return nil
}

func StackUpdate(id string, filename string) (error) {
	configRepository := config.NewConfigRepository(func(err error) {})
	stackRepository := api.NewStackRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))
	content, err := getStackFileContent(filename)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	content, err = yaml.YAMLToJSON(content)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	stackDefinition := make(map[string]interface{})
	if err := json.Unmarshal(content, &stackDefinition); err != nil {
		return err
	}

	stackModel, err := stackRepository.GetStack(id)
	if err != nil {
		return err
	}

	err = stackModel.Update(stackDefinition)
	if err != nil {
		return err
	} else {
		fmt.Printf("updated stack %s with uuid %s\n", stackModel.Name(), stackModel.Id())
	}
	return nil
}

func StackPublish(id string) (error) {
	configRepository := config.NewConfigRepository(func(err error) {})
	stackRepository := api.NewStackRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))

	stackModel, err := stackRepository.GetStack(id)
	if err != nil {
		return err
	}

	err = stackModel.Publish()
	if err != nil {
		return err
	} else {
		fmt.Printf("publish stack %s with uuid %s\n", stackModel.Name(), stackModel.Id())
	}
	return nil
}
