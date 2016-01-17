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

func DomainsList() error {
	configRepository := config.NewConfigRepository(func(err error) {})
	domainRepository := api.NewDomainRepository(configRepository, net.NewCloudControllerGateway(configRepository))
	domains, err := domainRepository.GetDomains()
	if err != nil {
		return err
	}
	fmt.Printf("=== Domains [%d]\n", len(domains.Items()))

	for _, domain := range domains.Items() {
		fmt.Printf("%s %s\n", domain.Id(), domain.Name())
	}
	return nil
}

func DomainsRemove(domainName string) error {
	configRepository := config.NewConfigRepository(func(err error) {})
	domainRepository := api.NewDomainRepository(configRepository, net.NewCloudControllerGateway(configRepository))
	domains, err := domainRepository.GetDomainByName(domainName)
	if err != nil {
		return err
	}
	domainId := domains.Items()[0].Id()
	err = domainRepository.Delete(domainId)
	if err != nil {
		return err
	}
	fmt.Printf("=== Domain [%s] deleted\n", domainName)
	return nil
}