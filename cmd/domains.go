package cmd

import (
	"fmt"
	"github.com/cnupp/cli/config"
	"github.com/cnupp/appssdk/api"
	"github.com/cnupp/appssdk/net"
	"io/ioutil"
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
	err := domainRepository.Delete(domainName)
	if err != nil {
		return err
	}
	fmt.Printf("=== Domain [%s] deleted\n", domainName)
	return nil
}

func DomainsCert(domainName, crtFile, privateKeyFile string) error {
	if "" == crtFile {
		return fmt.Errorf("%s", "The crt can not be empty")
	}

	if "" == privateKeyFile {
		return fmt.Errorf("%s", "The crt can not be empty")
	}

	crt, err := ioutil.ReadFile(crtFile)
	if err != nil {
		return fmt.Errorf("Please ensure %s exist and can be accessed", crtFile)
	}

	privateKey, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return fmt.Errorf("Please ensure %s exist and can be accessed", privateKeyFile)
	}

	configRepository := config.NewConfigRepository(func(err error) {})
	domainRepository := api.NewDomainRepository(configRepository, net.NewCloudControllerGateway(configRepository))
	domain, err := domainRepository.GetDomain(domainName)
	if err != nil {
		return err
	}

	err = domain.AttachCert(api.CertParams{
		Crt: string(crt),
		Key: string(privateKey),
	})

	if err != nil {
		return err
	}

	fmt.Printf("=== Domain(%s) crt [%s] key[%s] attached\n", domainName, crtFile, privateKeyFile)
	return nil
}
