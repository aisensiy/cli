package cmd
import (
	"fmt"
	"github.com/cde/apisdk/api"
	"github.com/cde/apisdk/config"
	"github.com/cde/apisdk/net"
)

// AppCreate creates an app.
func AppCreate(name string, stack string, owner string, memory int, disk int, instances int) error {
	configRepository := config.NewRepositoryFromFilepath("", func(err error) {})
	appRepository := api.NewAppRepository(config.NewRepositoryFromFilepath("", func(err error) {}),
		net.NewCloudControllerGateway(configRepository))
	appParams := api.AppParams{Name: name, Stack: stack,
		Owner:owner, Mem: memory,
		Disk:disk, Instances:instances}
	fmt.Println(appParams)
	createdApp, err := appRepository.Create(appParams)
	fmt.Println(createdApp)
	return err
}

// AppsList lists apps on the Deis controller.
func AppsList(results int) error {
	return nil
}