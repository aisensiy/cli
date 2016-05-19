package cmd
import (
	"github.com/sjkyspa/stacks/client/config"
"github.com/sjkyspa/stacks/controller/api/net"
	"github.com/sjkyspa/stacks/controller/api/api"
	"fmt"
	"errors"
)

func OrgCreate(orgName string) error {
	configRepository := config.NewConfigRepository(func(error) {})
	orgRepo := api.NewOrgRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))

	orgParams := api.OrgParams{
		Name: orgName,
	}
	createdOrg, err := orgRepo.Create(orgParams)
	if err != nil {
		return err
	}
	fmt.Println("create org success", createdOrg.Name())
	return err
}

func GetOrg(orgName string) error {
	configRepository := config.NewConfigRepository(func(error) {})
	orgRepo := api.NewOrgRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))

	org, err := orgRepo.GetOrg(orgName)
	if err != nil {
		return err
	}
	fmt.Println("get org success", org.Name())
	return err
}

func SetCurrentOrg(orgName string) error {
	configRepository := config.NewConfigRepository(func(error) {})
	orgRepo := api.NewOrgRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))

	_, err := orgRepo.GetOrg(orgName)
	if err != nil {
		return err
	}
	configRepository.SetCurrentOrg(orgName)
	fmt.Println("set current org as", orgName)
	return err
}

func ListMembers(orgName string) error {
	configRepository, orgName := loadOrg(orgName)

	if (orgName == "") {
		return errors.New("can not find default org")
	}

	orgRepo := api.NewOrgRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))

	users, err := orgRepo.GetOrgMembers(orgName)
	if err != nil {
		return err
	}

	fmt.Printf("=== Members: [%d]\n", len(users))

	for _, user := range users {
		fmt.Printf("email: %s\n", user.Email())
	}
	return nil
}
