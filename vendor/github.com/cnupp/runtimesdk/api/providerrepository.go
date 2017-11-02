package api

import (
	"github.com/cnupp/runtimesdk/config"
	"github.com/cnupp/runtimesdk/net"
	"encoding/json"
	"fmt"
)

type ProviderRepository interface {
	Enroll(params ProviderParams) (Provider, error)
	GetProviders() (Providers, error)
	GetProviderByName(name string) (Provider, error)
	UpdateProvider(id string, config map[string]interface{}) (error)
	GetProvidersByURL(uri string) (Providers, error)
}

type DefaultProviderRepository struct {
	config  config.Reader
	gateway net.Gateway
}

func (dpr DefaultProviderRepository) Enroll(params ProviderParams) (provider Provider, apiErr error) {
	data, err := json.Marshal(params)
	if err != nil {
		apiErr = err
		return
	}

	res, err := dpr.gateway.Request("POST", "/providers", data)
	if err != nil {
		apiErr = err
		return
	}

	location, err := res.Location()
	if err != nil {
		apiErr = err
		return
	}

	var providerModel ProviderModel
	err = dpr.gateway.Get(location.String(), &providerModel)
	if err != nil {
		apiErr = err
		return
	}

	provider = providerModel
	return
}

func (dpr DefaultProviderRepository) GetProviders() (providers Providers, err error) {
	providersModel := ProvidersModel{
		CountField: 0,
		ItemsField: []ProviderModel{},
	}

	page, err := dpr.GetProvidersByURL("/providers")
	if err != nil {
		return
	}
	providersModel.CountField = page.Count()

	for len(page.Items()) > 0 {
		providersModel.ItemsField = append(providersModel.ItemsField, page.Items()...)
		page, err = page.Next()
		if err != nil {
			return
		}
	}
	providers = providersModel

	return
}

func (dpr DefaultProviderRepository) GetProvidersByURL(url string) (providers Providers, err error) {
	var providersModel ProvidersModel

	err = dpr.gateway.Get(url, &providersModel)
	providersModel.ProviderRepo = dpr
	providers = providersModel
	return
}

func (dpr DefaultProviderRepository) GetProviderByName(name string) (provider Provider, err error) {
	var providersModel ProvidersModel

	err = dpr.gateway.Get(fmt.Sprintf("/providers?name=%s", name), &providersModel)
	if err != nil {
		return
	}
	if providersModel.Count() < 1 {
		err = fmt.Errorf("provider \"%s\" not found", name)
		return
	}

	var providerModel ProviderModel
	err = dpr.gateway.Get(fmt.Sprintf("/providers/%s", providersModel.Items()[0].ID()), &providerModel)
	if err != nil {
		return
	}
	provider = providerModel
	return
}

func (dpr DefaultProviderRepository) UpdateProvider(id string, config map[string]interface{}) (err error) {
	data, err := json.Marshal(config)
	if err != nil {
		return
	}

	_, err = dpr.gateway.Request("PUT", fmt.Sprintf("/providers/%s", id), data)

	return
}

func NewProviderRepository(config config.Reader, gateway net.Gateway) ProviderRepository {
	return DefaultProviderRepository{config: config, gateway: gateway}
}
