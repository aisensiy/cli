package api

import (
	"encoding/json"
	"fmt"
	"github.com/sjkyspa/stacks/controller/api/config"
	"github.com/sjkyspa/stacks/controller/api/net"
)

//go:generate counterfeiter -o fakes/fake_route_repository.go . RouteRepository
type RouteRepository interface {
	Create(params RouteParams) (apiErr error)
	GetRoutes() (routes Routes, apiErr error)
	GetAppsForRoute(routeId string) (apps Apps, apiErr error)
}

type DefaultRouteRepository struct {
	config  config.Reader
	gateway net.Gateway
}

func NewRouteRepository(config config.Reader, gateway net.Gateway) RouteRepository {
	return DefaultRouteRepository{config: config, gateway: gateway}
}

func (repo DefaultRouteRepository) Create(params RouteParams) (apiErr error) {
	data, err := json.Marshal(params)
	if err != nil {
		apiErr = fmt.Errorf("Can not serilize the data")
		return
	}

	res, apiErr := repo.gateway.Request("POST", "/routes", data)
	if apiErr != nil {
		return
	}
	location := res.Header.Get("Location")
	fmt.Println(location)
	return
}

func (repo DefaultRouteRepository) GetRoutes() (routes Routes, apiErr error) {
	var routeModels RoutesModel
	apiErr = repo.gateway.Get(fmt.Sprintf("/routes"), &routeModels)
	routes = routeModels
	return
}

func (repo DefaultRouteRepository) GetAppsForRoute(routeId string) (apps Apps, apiErr error) {
	var appModels AppsModel
	apiErr = repo.gateway.Get(fmt.Sprintf("/routes/"+routeId+"/apps"), &appModels)
	//	appModels.AppMapper = NewAppRepository(repo.config, repo.gateway)
	apps = appModels
	return
}
