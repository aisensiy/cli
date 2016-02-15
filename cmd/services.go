package cmd

import (
	"fmt"
	"github.com/sjkyspa/stacks/Godeps/_workspace/src/github.com/gambol99/go-marathon"
	deployApi "github.com/sjkyspa/stacks/deploymentsdk/api"
	deployNet "github.com/sjkyspa/stacks/deploymentsdk/net"
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
	application.Container.Docker.Container("registry.tw.com:80/gocd-server").Bridged().ExposePort(8153, 0, 0, "tcp")
	applicationCreated, err := client.CreateApplication(application)

	if err != nil {
		return err
	}
	fmt.Println(applicationCreated.ID)
	return nil
}

func ServiceInfo(appName, serviceName string) (apiErr error) {
	service, apiErr := GetService(appName,serviceName)
	if apiErr != nil {
		return apiErr
	}

	fmt.Printf("=== %s Service\n", service.Name())
	fmt.Println("id:        ", service.Id())
	fmt.Println("instances: ", service.Instance())
	fmt.Println("memory:    ", service.Memory())
	fmt.Println("cpus:    ", service.CPU())
	fmt.Println("disk:      ", service.Disk())
	return
}

func ServiceUpdate(appId, serviceName string, params deployApi.ServiceConfigParams) ( apiErr error) {
	service, apiErr := GetService(appId, serviceName)
	if apiErr != nil {
		return apiErr
	}
	apiErr = service.Update(params)
	return
}

func GetService(appId, serviceName string) (service deployApi.Service, apiErr error) {
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
