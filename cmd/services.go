package cmd

import (
	"fmt"
	"github.com/sjkyspa/stacks/Godeps/_workspace/src/github.com/gambol99/go-marathon"
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
