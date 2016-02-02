package cmd
import (
	"fmt"
	deployApi "github.com/sjkyspa/stacks/deploymentsdk/api"
	"github.com/sjkyspa/stacks/deploymentsdk/net"
	"github.com/sjkyspa/stacks/client/config"
)


func RestartApp(appName string) error {
	configRepository := config.NewConfigRepository(func(error) {})
	deployRepo := deployApi.NewDeploymentRepository(configRepository, net.NewCloudControllerGateway(configRepository))
	deployment, err := deployRepo.GetDeploymentByAppName(appName)
	if err != nil {
		return err
	}

	err = deployment.Restart()
	if err != nil {
		return err
	}
	fmt.Println("Restart the application successfully")
	return nil
}

func Scale(appName, serviceName string, params deployApi.ServiceConfigParams) (apiErr error) {
	configRepository := config.NewConfigRepository(func(error) {})
	deployRepo := deployApi.NewDeploymentRepository(configRepository, net.NewCloudControllerGateway(configRepository))
	deployment, apiErr := deployRepo.GetDeploymentByAppName(appName)
	if apiErr != nil {
		return apiErr
	}
	service, apiErr := deployment.GetService(serviceName)
	if apiErr != nil {
		return apiErr
	}
	apiErr = service.Scale(params)
	return
}

