package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/olekukonko/tablewriter"
	"github.com/cnupp/cli/config"
	"github.com/sjkyspa/stacks/controller/api/api"
	"github.com/sjkyspa/stacks/controller/api/net"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
)

func StackCreate(filename string) error {
	configRepository := config.NewConfigRepository(func(err error) {})
	stackRepository := api.NewStackRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))
	content, err := getStackFileContent(filename)
	if err != nil {
		return err
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
	if err != nil || stacks.Count() == 0 {
		err = fmt.Errorf("no stack found")
		return err
	}
	fmt.Printf("=== Stacks: [%d]\n", len(stacks.Items()))

	for _, stack := range stacks.Items() {
		fmt.Printf("name: %s id: %s\n", stack.Name(), stack.Id())
	}
	return nil
}

func GetStack(stackName string) error {
	configRepository := config.NewConfigRepository(func(error) {})
	stackRepository := api.NewStackRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))
	stacks, err := stackRepository.GetStackByName(stackName)

	if err != nil || stacks.Count() == 0 {
		err = fmt.Errorf("stack not found")
		return err
	}

	stackId := stacks.Items()[0].Id()
	stackObj, err := stackRepository.GetStack(stackId)

	outputStackDescription(stackObj)
	outputStackTemplate(stackObj.GetTemplate())
	outputStackLanguages(stackObj.GetLanguages())
	outputStackFrameworks(stackObj.GetFrameworks())
	outputStackServices(stackObj.GetServices())
	return nil
}

func outputStackDescription(stack api.Stack) {
	fmt.Printf("--- %s Stack\n", stack.Name())
	data := make([][]string, 3)
	data[0] = []string{"Type", stack.Type()}
	data[1] = []string{"Status", stack.GetStatus()}
	data[2] = []string{"Description", stack.GetDescription()}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}

func outputStackTemplate(template api.Template) {
	fmt.Printf("--- Template\n")
	data := make([][]string, 2)
	data[0] = []string{"Type", template.Type}
	data[1] = []string{"Uri", template.URI}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}

func outputStackLanguages(languages []api.Language) {
	fmt.Printf("--- Languages\n")
	data := make([][]string, len(languages))
	for index, language := range languages {
		data[index] = []string{strconv.Itoa(index), language.Name, language.Version}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeader([]string{"", "Name", "Version"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}

func outputStackFrameworks(frameworks []api.Framework) {
	fmt.Printf("--- Frameworks\n")
	data := make([][]string, len(frameworks))
	for index, framework := range frameworks {
		data[index] = []string{strconv.Itoa(index), framework.Name, framework.Version}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeader([]string{"", "Name", "Version"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}

func outputStackServices(services map[string]api.Service) {
	fmt.Println("--- Services")
	keys := reflect.ValueOf(services).MapKeys()

	for _, key := range keys {
		var data [][]string
		service := services[key.String()]
		data = append(data, []string{key.String(), "name", service.GetName()})
		data = append(data, []string{key.String(), "image", service.GetImage()})
		data = append(data, []string{key.String(), "mem", fmt.Sprintf("%f", service.GetMem())})
		data = append(data, []string{key.String(), "cpus", fmt.Sprintf("%f", service.GetCpu())})
		data = append(data, []string{key.String(), "instances", fmt.Sprintf("%d", service.GetInstances())})
		data = append(data, []string{key.String(), "expose", fmt.Sprintf("%d", service.GetExpose())})
		if api.NotEmptyImage(service.GetBuild()) {
			data = append(data, []string{key.String(), "build", fmt.Sprintf("image:%s", service.GetBuild().Name)})
			data = append(data, []string{key.String(), "build", fmt.Sprintf("mem:%d", service.GetBuild().Mem)})
			data = append(data, []string{key.String(), "build", fmt.Sprintf("cpus:%f", service.GetBuild().Cpus)})
		}
		if api.NotEmptyImage(service.GetVerify()) {
			data = append(data, []string{key.String(), "verify", fmt.Sprintf("image:%s", service.GetVerify().Name)})
			data = append(data, []string{key.String(), "verify", fmt.Sprintf("mem:%d", service.GetVerify().Mem)})
			data = append(data, []string{key.String(), "verify", fmt.Sprintf("cpus:%f", service.GetVerify().Cpus)})
		}
		healthChecks := service.GetHealthChecks()
		for index, healthcheck := range healthChecks {
			data = append(data, []string{key.String(), fmt.Sprintf("healthcheck %d", index), fmt.Sprintf("protocol:%s", healthcheck.Protocol)})
			data = append(data, []string{key.String(), fmt.Sprintf("healthcheck %d", index), fmt.Sprintf("command:%s", healthcheck.Command)})
			data = append(data, []string{key.String(), fmt.Sprintf("healthcheck %d", index), fmt.Sprintf("path:%s", healthcheck.Path)})
			data = append(data, []string{key.String(), fmt.Sprintf("healthcheck %d", index), fmt.Sprintf("grace:%d", healthcheck.Grace)})
			data = append(data, []string{key.String(), fmt.Sprintf("healthcheck %d", index), fmt.Sprintf("timeout:%d", healthcheck.Timeout)})
			data = append(data, []string{key.String(), fmt.Sprintf("healthcheck %d", index), fmt.Sprintf("interval:%d", healthcheck.Interval)})
			data = append(data, []string{key.String(), fmt.Sprintf("healthcheck %d", index), fmt.Sprintf("port:%d", healthcheck.Port)})
			data = append(data, []string{key.String(), fmt.Sprintf("healthcheck %d", index), fmt.Sprintf("portIndex:%d", healthcheck.PortIndex)})
			data = append(data, []string{key.String(), fmt.Sprintf("healthcheck %d", index), fmt.Sprintf("maxConsecutiveFailures:%d", healthcheck.MaxConsecutiveFailures)})
		}
		envKeys := reflect.ValueOf(service.GetEnv()).MapKeys()
		if len(envKeys) > 0 {
			envs := service.GetEnv()
			for _, envKey := range envKeys {
				data = append(data, []string{key.String(), "environment", fmt.Sprintf("%s:%s", envKey.String(), envs[envKey.String()])})
			}
		}
		volumes := service.GetVolumes()
		for index, volume := range volumes {
			data = append(data, []string{key.String(), fmt.Sprintf("volume %d", index), fmt.Sprintf("container:%s", volume.ContainerPath)})
			data = append(data, []string{key.String(), fmt.Sprintf("volume %d", index), fmt.Sprintf("host:%s", volume.HostPath)})
			data = append(data, []string{key.String(), fmt.Sprintf("volume %d", index), fmt.Sprintf("mode:%s", volume.Mode)})
		}
		data = append(data, []string{key.String(), "buildable", fmt.Sprintf("%s", strconv.FormatBool(service.IsBuildable()))})

		table := tablewriter.NewWriter(os.Stdout)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetRowSeparator("-")
		table.SetAutoMergeCells(true)
		table.SetRowLine(true)
		table.AppendBulk(data)
		table.Render()
	}
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

func StackUpdate(id string, filename string) error {
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

func StackPublish(id string) error {
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

func StackUnPublish(id string) error {
	configRepository := config.NewConfigRepository(func(err error) {})
	stackRepository := api.NewStackRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))

	stackModel, err := stackRepository.GetStack(id)
	if err != nil {
		return err
	}

	err = stackModel.UnPublish()
	if err != nil {
		return err
	} else {
		fmt.Printf("unpublish stack %s with uuid %s\n", stackModel.Name(), stackModel.Id())
	}
	return nil
}
