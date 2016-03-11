package cmd
import (
	"fmt"
	deployApi "github.com/sjkyspa/stacks/launcher/api/api"
	"github.com/sjkyspa/stacks/launcher/api/net"
)


func RestartApp(appId string) error {
	configRepository, appId, err := load(appId)
	if err != nil {
		return err
	}
	deployRepo := deployApi.NewDeploymentRepository(configRepository, net.NewCloudControllerGateway(configRepository))
	deployment, err := deployRepo.GetDeploymentByAppName(appId)
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

func Scale(appId, serviceName string, params deployApi.ServiceConfigParams) (apiErr error) {
//	configRepository := config.NewConfigRepository(func(error) {})
//	deployRepo := deployApi.NewDeploymentRepository(configRepository, net.NewCloudControllerGateway(configRepository))
//	deployment, apiErr := deployRepo.GetDeploymentByAppName(appName)
//	if apiErr != nil {
//		return apiErr
//	}
	service, apiErr := GetService(appId, serviceName)
	if apiErr != nil {
		return apiErr
	}
	apiErr = service.Update(params)
	return
}

func ListDependentServices(appName string) error {
	return outputDependentServices(appName)
}
