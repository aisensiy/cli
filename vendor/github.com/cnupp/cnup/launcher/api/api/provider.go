package api

import (
	"github.com/cnupp/cnup/controller/api/api"
)

type Provider interface {
	ID() string
	Name() string
	Owner() string
	Consumer() string
	Type() string
	CreatedAt() int
	Config() map[string]interface{}
	Links() api.Links
}

type ProviderModel struct {
	IDField        string                    `json:"id"`
	NameField      string                    `json:"name"`
	TypeField      string                    `json:"type"`
	OwnerField     string                    `json:"owner"`
	ConsumerField  string                    `json:"consumer"`
	ConfigField    map[string]interface{}    `json:"config"`
	CreatedAtField int 			 `json:"created_at"`
	LinksArray     []api.Link                `json:"links"`
}

type ProviderParams struct {
	Name   string 			`json:"name"`
	Type   string 			`json:"type"`
	Config map[string]interface{} 	`json:"config"`
	Consumer string			`json:"consumer,omitempty"`
}

func (pm ProviderModel) ID() string {
	return pm.IDField
}

func (pm ProviderModel) Name() string {
	return pm.NameField
}

func (pm ProviderModel) Owner() string {
	return pm.OwnerField
}

func (pm ProviderModel) Consumer() string {
	return pm.ConsumerField
}

func (pm ProviderModel) CreatedAt() int {
	return pm.CreatedAtField
}

func (pm ProviderModel) Links() api.Links {
	return api.LinksModel{
		Links: pm.LinksArray,
	}
}

func (pm ProviderModel) Type() string {
	return pm.TypeField
}

func (pm ProviderModel) Config() map[string]interface{} {
	return pm.ConfigField
}

type Providers interface {
	Items() []ProviderModel
	Count() int
	Next() (Providers, error)
}

type ProvidersModel struct {
	FirstField   string                `json:"first"`
	LastField    string                `json:"last"`
	PrevField    string                `json:"prev"`
	NextField    string                `json:"next"`
	CountField   int                   `json:"count"`
	SelfField    string                `json:"self"`
	ItemsField   []ProviderModel       `json:"items"`
	ProviderRepo ProviderRepository
}

func (psm ProvidersModel) Items() []ProviderModel {
	return psm.ItemsField
}

func (psm ProvidersModel) Count() int {
	return psm.CountField
}

func (psm ProvidersModel) Next() (providers Providers, err error) {
	if psm.NextField == "" {
		providers = ProvidersModel{
			ItemsField: []ProviderModel{},
		}
		return
	}

	providers, err = psm.ProviderRepo.GetProvidersByURL(psm.NextField)
	return
}
