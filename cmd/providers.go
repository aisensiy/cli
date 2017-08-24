package cmd

import (
	"github.com/sjkyspa/stacks/client/config"
	"github.com/sjkyspa/stacks/launcher/api/api"
	"github.com/sjkyspa/stacks/launcher/api/net"
	"fmt"
)

func ProviderCreate(providerName string, providerType string,consumer string, configMap map[string]interface{}) error {
	configRepository := config.NewConfigRepository(func(error) {})
	providerRepository := api.NewProviderRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))

	provider, err := providerRepository.Enroll(api.ProviderParams{
		Name:   providerName,
		Type:   providerType,
		Config: configMap,
		Consumer: consumer,
	})

	if err != nil {
		return err
	}
	fmt.Printf("create provider %s\n", provider.Name())
	return nil
}

func ProviderList() error {
	configRepository := config.NewConfigRepository(func(error) {})
	providerRepository := api.NewProviderRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))

	providers, err := providerRepository.GetProviders()
	if err != nil {
		return err
	}

	fmt.Printf("=== Providers [%d]\n", len(providers.Items()))

	for _, provider:= range providers.Items() {
		fmt.Printf("name: %s\n", provider.Name())
	}

	return nil
}
