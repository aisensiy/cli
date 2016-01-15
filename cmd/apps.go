package cmd
import (
	"fmt"
	"github.com/cde/apisdk/api"
	"github.com/cde/apisdk/net"
	"github.com/cde/client/config"
)

// AppCreate creates an app.
func AppCreate(name string, stack string, memory int, disk int, instances int) error {
	configRepository := config.NewConfigRepository(func(error) {})
	appRepository := api.NewAppRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))
	appParams := api.AppParams{
		Name: name,
		Stack: stack,
		Mem: memory,
		Disk:disk,
		Instances:instances}
	createdApp, err := appRepository.Create(appParams)
	if createdApp != nil {
		fmt.Println(createdApp)
	}
	return err
}

func AppsList() error {
	configRepository := config.NewConfigRepository(func(error) {})
	appRepository := api.NewAppRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))
	apps, err := appRepository.GetApps()
	if err != nil {
		return err
	}
	fmt.Printf("=== Apps [%d]\n", len(apps.Items()))

	for _, app := range apps.Items() {
		fmt.Printf("id: %s\n", app.Id())
	}
	return nil
}

func GetApp(appId string) error {
	configRepository := config.NewConfigRepository(func(error) {})
	appRepository := api.NewAppRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))
	app, err := appRepository.GetApp(appId)
	if err != nil {
		return err
	}

	fmt.Printf("=== %s Application\n", app.Id())
	fmt.Println("id:        ", app.Id())
	fmt.Println("instances: ", app.Instances())
	fmt.Println("memory:    ", app.Mem())
	fmt.Println("disk:      ", app.Disk())

	fmt.Println()

	return nil
}