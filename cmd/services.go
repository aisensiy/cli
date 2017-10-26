package cmd

import (
	"fmt"
	"github.com/gambol99/go-marathon"
	"github.com/olekukonko/tablewriter"
	deployApi "github.com/sjkyspa/stacks/launcher/api/api"
	deployNet "github.com/sjkyspa/stacks/launcher/api/net"
	"os"
)

func ServiceCreate() error {
	fmt.Println("create service")

	marathonURL := "http://marathon.tw.com"
	config := marathon.NewDefaultConfig()
	config.URL = marathonURL
	client, err := marathon.NewClient(config)
	if err != nil {
		return err
	}

	fmt.Println("connected to marathon")

	application := marathon.NewDockerApplication()
	application.Name("gocd111").CPU(1).Memory(441).Storage(0)
	application.Container.Docker.Container("registry.tw.com:80/gocd-server").Bridged().ExposePort(marathon.PortMapping{
		ContainerPort: 8153,
		HostPort:      0,
		Protocol:      "tcp",
		ServicePort:   0,
	})
	applicationCreated, err := client.CreateApplication(application)

	if err != nil {
		return err
	}
	fmt.Println(applicationCreated.ID)
	return nil
}

func ServiceInfo(appName, serviceName string) (apiErr error) {
	service, apiErr := GetService(appName, serviceName)
	if apiErr != nil {
		return apiErr
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	fmt.Printf("--- %s Service\n", service.Name())
	table.Append([]string{"ID", service.Id()})
	table.Append([]string{"instances", fmt.Sprintf("%d", service.Instance())})
	table.Append([]string{"memory", fmt.Sprintf("%v", service.Memory())})
	table.Append([]string{"cpus", fmt.Sprintf("%v", service.CPU())})

	table.Render() // Send output

	return
}

func ServiceUpdate(appId, serviceName string, params deployApi.ServiceConfigParams) (apiErr error) {
	service, apiErr := GetService(appId, serviceName)
	if apiErr != nil {
		return apiErr
	}
	apiErr = service.Update(params)
	return
}

func GetService(appId, serviceName string) (service deployApi.LauncherService, apiErr error) {
	configRepository, appId, err := load(appId)
	if err != nil {
		return nil, err
	}
	deployRepo := deployApi.NewDeploymentRepository(configRepository, deployNet.NewCloudControllerGateway(configRepository))
	deployment, apiErr := deployRepo.GetDeploymentByAppName(appId)
	if apiErr != nil {
		return
	}
	service, apiErr = deployment.GetService(serviceName)
	return
}
