package cmd
import (
	"fmt"
	"github.com/sjkyspa/stacks/apisdk/api"
	"github.com/sjkyspa/stacks/apisdk/net"
	"github.com/sjkyspa/stacks/client/config"
	"github.com/sjkyspa/stacks/client/pkg"
	"net/url"
	"strings"
	deploymentApi "github.com/sjkyspa/stacks/deploymentsdk/api"
	deploymentNet "github.com/sjkyspa/stacks/deploymentsdk/net"
)

// AppCreate creates an app.
func AppCreate(name string, stackName string, memory int, disk int, instances int) error {
	configRepository := config.NewConfigRepository(func(error) {})
	appRepository := api.NewAppRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))

	stackRepo := api.NewStackRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))

	stacks, err := stackRepo.GetStackByName(stackName)
	if err != nil {
		return err
	}
	stackId := stacks.Items()[0].Id()

	appParams := api.AppParams{
		Name: name,
		Stack: stackId,
		Mem: memory,
		Disk:disk,
		Instances:instances}
	createdApp, err := appRepository.Create(appParams)
	if err != nil {
		return err
	}
	u, err := url.Parse(configRepository.ApiEndpoint())
	if err != nil {
		return err
	}
	host := u.Host
	if strings.Index(host, ":") != -1 {
		splits := strings.Split(host, ":")
		host = splits[0]
	}
	git.DeleteCdeRemote()
	host = "192.168.50.6"
	if err = git.CreateRemote(host, "cde", createdApp.Id()); err != nil {
		if err.Error() == "exit status 128" {
			fmt.Println("To replace the existing git remote entry, run:")
			fmt.Printf("  git remote rename cde cde.old && cde git:remote -a %s\n", createdApp.Id())
		}
		return err
	}

	fmt.Println("remote available at", git.RemoteURL(host, createdApp.Id()))
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
	outputDescription(app)
	outputRoutes(app)
	outputDependentServices(appId)
	
	return nil
}

func outputDescription(app api.App){
	fmt.Printf("=== %s Application\n", app.Id())
	fmt.Println("id:        ", app.Id())
	fmt.Println("instances: ", app.Instances())
	fmt.Println("memory:    ", app.Mem())
	fmt.Println("disk:      ", app.Disk())
}

func outputRoutes(app api.App) {
	boundRoutes, err := app.GetRoutes()
	fmt.Println("--- Access routes:\n")

	if(err != nil) {
		fmt.Print(err)
		return
	}
	for boundRoutes != nil {
		routes := boundRoutes.Items()
		for _, route := range routes {
			fmt.Println(route.DomainField.Name + "/" + route.PathField + " \n")
		}
		boundRoutes, _ = boundRoutes.Next()
	}

}

func outputDependentServices(appId string) error {
	configRepository := config.NewConfigRepository(func(error) {})
	repo := deploymentApi.NewDeploymentRepository(configRepository, deploymentNet.NewCloudControllerGateway(configRepository))
	servicesModel, err := repo.GetDependentServicesForApp(appId)
	fmt.Print("--- Dependent services:\n")
	if(err != nil) {
		fmt.Print(err)
		return err
	}
	servicesArray := servicesModel.Apps()
	for _, service := range servicesArray  {
		fmt.Println("Id:      ", service.Id())
		fmt.Println("Instance(s):      ", service.Instance())
		fmt.Println("Memory:      ", service.Memory())
		fmt.Println("Env:      ", service.Env())
	}
	return nil
}

func DestroyApp(appId string) error {
	configRepository := config.NewConfigRepository(func(error) {})
	deployRepo := deploymentApi.NewDeploymentRepository(configRepository, deploymentNet.NewCloudControllerGateway(configRepository))
	err := deployRepo.Destroy(appId)
	if err != nil {
		return err
	}

	appRepository := api.NewAppRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))
	err = appRepository.Delete(appId)
	if err != nil {
		return err
	}

	if (git.HasRemoteNameForApp("cde", appId)) {
		err = git.DeleteRemote("cde")
		if (err != nil) {
			fmt.Print("Remove 'cde' remote failed. \n Please execute git cmd in the app directory: `git remote remove cde`")
		}
	}else {
		fmt.Print("Please execute git cmd in the app directory: `git remote remove cde`")
	}

	return nil
}

func SwitchStack(appName string, stackName string) error {

	configRepository := config.NewConfigRepository(func(error) {})
	appRepository := api.NewAppRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))
	params := api.UpdateStackParams{
		Stack: stackName,
	}
	err := appRepository.SwitchStack(appName, params)
	if err != nil {
		return err
	}

	return nil
}
