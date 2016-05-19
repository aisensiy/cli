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

func AddMember(orgName string, email string) error {
	configRepository, orgName := loadOrg(orgName)

	if (orgName == "") {
		return errors.New("can not find default org")
	}

	orgRepo := api.NewOrgRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))

	err := orgRepo.AddMember(orgName, email)
	if err != nil {
		return err
	}

	fmt.Printf("Add %s to org %s", email, orgName)
	return nil
}

func RemoveMember(orgName string, email string) error {
	configRepository, orgName := loadOrg(orgName)

	if (orgName == "") {
		return errors.New("can not find default org")
	}

	userRepo := api.NewUserRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))

	users, err := userRepo.GetUserByEmail(email)
	if err != nil {
		return err
	}

	user := users.Items()[0]

	orgRepo := api.NewOrgRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))

	err = orgRepo.RmMember(orgName, user.Id())
	if err != nil {
		return err
	}

	fmt.Printf("Remove user %s from org %s", email, orgName)
	return nil
}

func ListApps(orgName string) error {
	configRepository, orgName := loadOrg(orgName)

	if (orgName == "") {
		return errors.New("can not find default org")
	}

	orgRepo := api.NewOrgRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))

	apps, err := orgRepo.GetApps(orgName)
	if err != nil {
		return err
	}

	fmt.Printf("=== Apps: [%d]\n", len(apps))

	for _, app := range apps {
		fmt.Printf("email: %s\n", app.Id())
	}
	return nil
}

func AddOrgApp(orgName string, appName string) error {
	configRepository, orgName := loadOrg(orgName)

	if (orgName == "") {
		return errors.New("can not find default org")
	}

	orgRepo := api.NewOrgRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))

	err := orgRepo.AddApp(orgName, appName)
	if err != nil {
		return err
	}

	fmt.Printf("Add %s to org %s\n", appName, orgName)
	return nil
}

func DestroyOrg(orgName string) error {
	configRepository, orgName := loadOrg(orgName)

	if (orgName == "") {
		return errors.New("can not find default org")
	}

	orgRepo := api.NewOrgRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))

	err := orgRepo.Delete(orgName)
	if err != nil {
		return err
	}

	fmt.Printf("Destroy org %s\n", orgName)
	if (orgName == configRepository.Org()) {
		configRepository.SetCurrentOrg("")
	}
	return nil
}
