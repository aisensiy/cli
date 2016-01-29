package api

import (
	"encoding/json"
	"fmt"
	"github.com/cde/client/Godeps/_workspace/src/github.com/cde/apisdk/config"
	"github.com/cde/client/Godeps/_workspace/src/github.com/cde/apisdk/net"
)

//go:generate counterfeiter -o fakes/fake_app_repository.go . AppRepository
type AppRepository interface {
	Create(params AppParams) (createdApp App, apiErr error)
	GetApp(id string) (App, error)
	GetApps() (Apps, error)
	Update(id string, params AppParams) (updatedApp App, apiErr error)
	Delete(id string) (apiErr error)
	BindWithRoute(app App, params AppRouteParams) error
	UnbindRoute(app App, routeId string) error
	GetRoutes(app App) (routes AppRoutes, apiErr error)
	SetEnv(app App, kv KeyValue) error
	UnsetEnv(app App, key string) error
}

type CloudControllerAppRepository struct {
	config  config.Reader
	gateway net.Gateway
}

func NewAppRepository(config config.Reader, gateway net.Gateway) AppRepository {
	return CloudControllerAppRepository{config: config, gateway: gateway}
}

func (cc CloudControllerAppRepository) SetEnv(app App, kv KeyValue) (apiErr error) {
	data, err := json.Marshal(kv)
	if err != nil {
		apiErr = fmt.Errorf("Can not serilize the data")
		return
	}

	_, err = cc.gateway.Request("POST", fmt.Sprintf("/apps/%s/env", app.Id()), data)
	return err
}

func (cc CloudControllerAppRepository) UnsetEnv(app App, key string) (apiErr error) {
	apiErr = cc.gateway.Delete(fmt.Sprintf("/apps/%s/env/%s", app.Id(), key), nil)
	return
}

func (cc CloudControllerAppRepository) Create(params AppParams) (createdApp App, apiErr error) {
	data, err := json.Marshal(params)
	if err != nil {
		apiErr = fmt.Errorf("Can not serilize the data")
		return
	}

	res, err := cc.gateway.Request("POST", "/apps", data)
	if err != nil {
		apiErr = err
		return
	}

	location := res.Header.Get("Location")

	var appModel AppModel
	apiErr = cc.gateway.Get(location, &appModel)
	if apiErr != nil {
		return
	}
	appModel.BuildMapper = NewBuildMapper(cc.config, cc.gateway)
	appModel.AppMapper = NewAppRepository(cc.config, cc.gateway)
	appModel.StackRepository = NewStackRepository(cc.config, cc.gateway)
	createdApp = appModel

	return
}

func (cc CloudControllerAppRepository) GetApp(id string) (app App, apiErr error) {
	var remoteApp AppModel
	apiErr = cc.gateway.Get(fmt.Sprintf("/apps/%s", id), &remoteApp)
	if apiErr != nil {
		return
	}
	remoteApp.BuildMapper = NewBuildMapper(cc.config, cc.gateway)
	remoteApp.AppMapper = NewAppRepository(cc.config, cc.gateway)
	remoteApp.StackRepository = NewStackRepository(cc.config, cc.gateway)
	app = remoteApp
	return
}

func (cc CloudControllerAppRepository) GetApps() (apps Apps, apiErr error) {
	var remoteApps AppsModel
	apiErr = cc.gateway.Get(fmt.Sprintf("/apps"), &remoteApps)
	if apiErr != nil {
		return
	}
	remoteApps.AppMapper = cc
	apps = remoteApps
	return
}

func (cc CloudControllerAppRepository) BindWithRoute(app App, params AppRouteParams) error {
	data, err := json.Marshal(params)
	if err != nil {
		err = fmt.Errorf("Can not serilize the data")
		return err
	}

	_, err = cc.gateway.Request("POST", fmt.Sprintf("/apps/%s/routes", app.Id()), data)

	return err
}

func (cc CloudControllerAppRepository) UnbindRoute(app App, routeId string) error {
	_, err := cc.gateway.Request("DELETE", fmt.Sprintf("/apps/%s/routes/%s", app.Id(), routeId), nil)
	return err
}

func (cc CloudControllerAppRepository) Update(id string, params AppParams) (updatedApp App, apiErr error) {
	return
}

func (cc CloudControllerAppRepository) Delete(id string) (apiErr error) {
	return
}

func (cc CloudControllerAppRepository) GetRoutes(app App) (routes AppRoutes, apiErr error) {
	var routesModel AppRoutesModel
	apiErr = cc.gateway.Get(fmt.Sprintf("/apps/"+app.Id()+"/routes"), &routesModel)
	routes = routesModel
	return
}
