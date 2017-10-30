package cmd

import (
	"errors"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/sjkyspa/cde/config"
	"github.com/sjkyspa/stacks/launcher/api/api"
	"github.com/sjkyspa/stacks/launcher/api/net"
	"os"
	"time"
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
		data = append(data, []string{provider.Name(), provider.Type(), provider.Owner(), provider.Consumer(), time.Unix(int64(provider.CreatedAt()/1000), 0).String()})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetRowSeparator("-")
	table.SetRowLine(true)
	table.AppendBulk(data)
	table.Render()
}

func GetProviderByName(name string) error {
	configRepository := config.NewConfigRepository(func(error) {})
	providerRepository := api.NewProviderRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))

	provider, err := providerRepository.GetProviderByName(name)
	if err != nil {
		return err
	}

	outputProviderInfo(provider)

	return nil
}

func outputProviderInfo(provider api.Provider) {
	var data [][]string
	data = append(data, []string{"name", provider.Name(), ""})
	data = append(data, []string{"type", provider.Type(), ""})
	data = append(data, []string{"owner", provider.Owner(), ""})
	data = append(data, []string{"for", provider.Consumer(), ""})
	data = append(data, []string{"created_at", time.Unix(int64(provider.CreatedAt()/1000), 0).String(), ""})

	if len(provider.Config()) > 0 {
		for key, value := range provider.Config() {
			data = append(data, []string{"config", key, value.(string)})
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.AppendBulk(data)
	table.Render()
}

func ProviderUpdate(providerName string, updateConfigMap map[string]interface{}) error {
	configRepository := config.NewConfigRepository(func(error) {})
	providerRepository := api.NewProviderRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))

	provider, err := providerRepository.GetProviderByName(providerName)
	if err != nil {
		return err
	}

	configMap := provider.Config()
	for k, v := range updateConfigMap {
		if v == "" {
			_, ok := configMap[k]
			if ok {
				delete(configMap, k)
			} else {
				return errors.New(fmt.Sprintf("could not remove not existed key %s in config", k))
			}
		} else {
			configMap[k] = v
		}
	}

	if err := providerRepository.UpdateProvider(provider.ID(), map[string]interface{}{"config": configMap}); err != nil {
		return err
	}

	fmt.Printf("Provider %s config updated\n", provider.Name())

	return nil
}
