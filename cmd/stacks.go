package cmd

import (
	"fmt"
	"github.com/cde/apisdk/api"
	"github.com/cde/apisdk/net"
	"github.com/cde/client/config"
)

func StackCreate(name string) error {
	configRepository := config.NewConfigRepository(func(err error) {})
	stackRepository := api.NewStackRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))
	stackParams := api.StackParams{
		Name: name,
	}
	stackModel, err := stackRepository.Create(stackParams)
	if err != nil {
		return err
	}
	fmt.Printf("create stack %s with uuid %s\n", stackModel.Name(), stackModel.Id())
	return nil
}
