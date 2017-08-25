package cmd

import (
	"github.com/sjkyspa/stacks/client/config"
	"github.com/sjkyspa/stacks/launcher/api/api"
	"github.com/sjkyspa/stacks/launcher/api/net"
	"fmt"
	"strconv"
	"github.com/olekukonko/tablewriter"
	"os"
)

func ProviderCreate(providerName string, providerType string, consumer string, configMap map[string]interface{}) error {
	configRepository := config.NewConfigRepository(func(error) {})
	providerRepository := api.NewProviderRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))

	provider, err := providerRepository.Enroll(api.ProviderParams{
		Name:     providerName,
		Type:     providerType,
		Config:   configMap,
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

	outputProvidersListInfo(providers)

	return nil
}

func outputProvidersListInfo(providers api.Providers) {
	fmt.Printf("=== Providers [%d]\n", len(providers.Items()))
	var data [][]string
	data = append(data, []string{"name", "type", "owner", "for", "created_at"})

	for _, provider := range providers.Items() {
		data = append(data, []string{provider.Name(), provider.Type(), provider.Owner(), provider.Consumer(), strconv.Itoa(provider.CreatedAt())})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetRowSeparator("-")
	table.SetRowLine(true)
	table.AppendBulk(data)
	table.Render()
}
