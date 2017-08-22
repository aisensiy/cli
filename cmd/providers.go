package cmd

import (
	"github.com/sjkyspa/stacks/client/config"
	"github.com/sjkyspa/stacks/launcher/api/api"
	"github.com/sjkyspa/stacks/launcher/api/net"
	"fmt"
)

func ProviderCreate(providerName string, providerType string, configMap map[string]interface{}) error {
	configRepository := config.NewConfigRepository(func(error) {})
	providerRepository := api.NewProviderRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))

	provider, err := providerRepository.Enroll(api.ProviderParams{
		Name:   providerName,
		Type:   providerType,
		Config: configMap,
	})

	if err != nil {
		return err
	}
	fmt.Printf("create provider %s\n", provider.Name())
	return nil
}
