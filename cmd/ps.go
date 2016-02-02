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

