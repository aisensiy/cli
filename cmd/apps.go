package cmd
import (
	"os"
	"fmt"
	"errors"
	"strings"
	"net/url"
	"github.com/olekukonko/tablewriter"
	"github.com/sjkyspa/stacks/apisdk/api"
	"github.com/sjkyspa/stacks/apisdk/net"
	"github.com/sjkyspa/stacks/client/config"
	"github.com/sjkyspa/stacks/client/pkg"
	deploymentApi "github.com/sjkyspa/stacks/deploymentsdk/api"
	deploymentNet "github.com/sjkyspa/stacks/deploymentsdk/net"
)

// AppCreate creates an app.
func AppCreate(appId string, stackName string, memory int, disk int, instances int) error {
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
		Name: appId,
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
	host = configRepository.GitHost()
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
	configRepository, appId, err := load(appId)
	if err != nil {
		return err
	}
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

func outputDescription(app api.App) {
	fmt.Printf("--- %s Application\n", app.Id())
	data := make([][]string, 4)
	data[0] = []string{"ID", app.Id()}
	data[1] = []string{"instances", fmt.Sprintf("%d", app.Instances())}
	data[2] = []string{"memory", fmt.Sprintf("%d", app.Mem())}
	data[3] = []string{"disk", fmt.Sprintf("%d", app.Disk())}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, v := range data {
		table.Append(v)
	}
	table.Render() // Send output
}

func outputRoutes(app api.App) {
	boundRoutes, err := app.GetRoutes()
	fmt.Println("--- Access routes:\n")

	if (err != nil) {
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
	configRepository, appId, err := load(appId)
	if err != nil {
		return err
	}
	repo := deploymentApi.NewDeploymentRepository(configRepository, deploymentNet.NewCloudControllerGateway(configRepository))
	servicesModel, err := repo.GetDependentServicesForApp(appId)
	fmt.Print("--- Dependent services:\n")
	if (err != nil) {
		fmt.Print(err)
		return err
	}
	servicesArray := servicesModel.Apps()
	for index, service := range servicesArray {
		fmt.Printf("-----> Service %d:\n", index+1)
		table := tablewriter.NewWriter(os.Stdout)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.Append([]string{"ID", service.Id()})
		table.Append([]string{"Name", service.Name()})
		table.Append([]string{"Instances", fmt.Sprintf("%d", service.Instance())})
		table.Append([]string{"Memory", fmt.Sprintf("%v", service.Memory())})
		table.Append([]string{"Env", service.Env()})
		table.Render() // Send output
	}
	return nil
}

func DestroyApp(appId string) error {
	configRepository, currentApp, err := load("")
	if err != nil {
		return err
	}
	if appId != "" && appId != currentApp {
		return errors.New(fmt.Sprintf("current dir's app %s != %s\n", currentApp, appId))
	}
	deployRepo := deploymentApi.NewDeploymentRepository(configRepository, deploymentNet.NewCloudControllerGateway(configRepository))
	err = deployRepo.Destroy(currentApp)
	if err != nil {
		return err
	}

	appRepository := api.NewAppRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))
	err = appRepository.Delete(currentApp)
	if err != nil {
		return err
	}

	fmt.Printf("destroy %s successfully!\n", currentApp)

	if (git.HasRemoteNameForApp("cde", currentApp)) {
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
	configRepository, appName, err := load(appName)
	if err != nil {
		return err
	}
	appRepository := api.NewAppRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))
	params := api.UpdateStackParams{
		Stack: stackName,
	}
	err = appRepository.SwitchStack(appName, params)
	if err != nil {
		return err
	}

	return nil
}

func AppLog(appId string, lines int) error {
	configRepository, appId, err := load(appId)
	if err != nil {
		return err
	}
	deploymentRepository := deploymentApi.NewDeploymentRepository(configRepository,
		deploymentNet.NewCloudControllerGateway(configRepository))
	deployment, err := deploymentRepository.GetDeploymentByAppName(appId)
	if err != nil {
		return err
	}
	output, err := deployment.Log(lines)
	if err != nil {
		return err
	}
	fmt.Println(output)
	return nil
}

func ServiceLog(appId, serviceName string, lines int) error {
	configRepository, appId, err := load(appId)
	if err != nil {
		return err
	}
	deploymentRepository := deploymentApi.NewDeploymentRepository(configRepository,
		deploymentNet.NewCloudControllerGateway(configRepository))
	deployment, err := deploymentRepository.GetDeploymentByAppName(appId)
	if err != nil {
		return err
	}

	service, err := deployment.GetService(serviceName)
	if err != nil {
		return err
	}
	output, err := service.Log(lines)
	if err != nil {
		return err
	}
	fmt.Println(output)
	return nil
}
