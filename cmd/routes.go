package cmd
import (
	"fmt"
	"github.com/sjkyspa/stacks/client/config"
	"github.com/sjkyspa/stacks/apisdk/net"
	"github.com/sjkyspa/stacks/apisdk/api"
)

// RouteCreate creates an route.
func RoutesCreate(domainName string, path string) error {
	configRepository := config.NewConfigRepository(func(error) {})
	routeRepository := api.NewRouteRepository(configRepository, net.NewCloudControllerGateway(configRepository))
	domainRepository := api.NewDomainRepository(configRepository, net.NewCloudControllerGateway(configRepository))
	domain, err := domainRepository.GetDomain(domainName)
	if err != nil {
		return err
	}
	domainId := domain.Id()
	routeParams := api.RouteParams{
		Domain: domainId,
		Path: path,
	}
	err = routeRepository.Create(routeParams)
	if err != nil {
		fmt.Println(err)
	}

	return err
}

func RoutesList() error {
	configRepository := config.NewConfigRepository(func(error) {})
	routeRepository := api.NewRouteRepository(configRepository, net.NewCloudControllerGateway(configRepository))
	routes, err := routeRepository.GetRoutes()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("=== Routes: [%d]\n", len(routes.Items()))

	for _, route := range routes.Items() {
		fmt.Printf("id: %s path: %s domain: %s\n", route.ID(), route.Path(), route.Domain().Name)
	}
	return err
}

func RouteBindWithApp(route, appName string) error {
	configRepository := config.NewConfigRepository(func(error) {})
	appRepo := api.NewAppRepository(configRepository, net.NewCloudControllerGateway(configRepository))
	app, err := appRepo.GetApp(appName)
	if err != nil {
		return err
	}
	routeParams := api.AppRouteParams{
		Route: route,
	}
	err = app.BindWithRoute(routeParams)
	if err != nil {
		return err
	}

	fmt.Printf("=== Bind successfully\n")

	return err
}

func UnbindRouteWithApp(routeId, appName string) error {
	configRepository := config.NewConfigRepository(func(error) {})
	appRepo := api.NewAppRepository(configRepository, net.NewCloudControllerGateway(configRepository))
	app, err := appRepo.GetApp(appName)
	if err != nil {
		return err
	}
	err = app.UnbindRoute(routeId)
	if err != nil {
		return err
	}

	fmt.Printf("=== Unbind successfully\n")

	return err
}
