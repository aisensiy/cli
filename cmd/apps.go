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
	fmt.Println(appParams)
	createdApp, err := appRepository.Create(appParams)
	fmt.Println(createdApp)
	return err
}

// AppsList lists apps on the Deis controller.
func AppsList(results int) error {
	return nil
}