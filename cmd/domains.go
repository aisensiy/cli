package cmd

import (
	"fmt"
	"github.com/cde/apisdk/api"
	"github.com/cde/apisdk/net"
	"github.com/cde/client/config"
)

func DomainsAdd(name string) error {
	configRepository := config.NewConfigRepository(func(err error) {})
	domainRepository := api.NewDomainRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))
	domainParams := api.DomainParams{
		Name: name,
	}
	_, err := domainRepository.Create(domainParams)
	if err != nil {
		return err
	}
	fmt.Printf("create domain %s\n", name)
	return nil
}
