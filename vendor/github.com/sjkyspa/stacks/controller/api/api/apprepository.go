package api

import (
	"encoding/json"
	"fmt"
	"github.com/sjkyspa/stacks/controller/api/config"
	"github.com/sjkyspa/stacks/controller/api/net"
	"errors"
)

//go:generate counterfeiter -o fakes/fake_app_repository.go . AppRepository
type AppRepository interface {
	Create(params AppParams) (createdApp App, apiErr error)
	GetApp(id string) (App, error)
	GetApps() (Apps, error)
	Delete(id string) (apiErr error)
	BindWithRoute(app App, params AppRouteParams) error
	UnbindRoute(app App, routeId string) error
	GetRoutes(app App) (routes AppRoutes, apiErr error)
	GetRoutesByURI(uri string) (routes AppRoutes, apiErr error)
	SetEnv(app App, kvs map[string]interface{}) error
	UnsetEnv(app App, keys []string) error
	SwitchStack(id string, params UpdateStackParams) (apiErr error)
	GetLog(appId, buildId, logType string, lines int64, offset int64) (LogsModel, error)
	GetPermission(app App, userId string) (AppPermission, error)
	GetCollaborators(appId string) ([]UserModel, error)
	AddCollaborator(appId string, param CreateCollaboratorParams) error
	RemoveCollaborator(appId string, userId string) error
	TransferToUser(appId string, email string) error
	TransferToOrg(appId string, orgName string) error
}

type AppPermission struct {
	Write bool `json:"write"`
	Read  bool `json:"read"`
}

type TransferAppParams struct {
	Owner     string `json:"owner"`
	OwnerType string `json:"owner_type"`
}

type CloudControllerAppRepository struct {
	config  config.Reader
	gateway net.Gateway
}

func NewAppRepository(config config.Reader, gateway net.Gateway) AppRepository {
	return CloudControllerAppRepository{config: config, gateway: gateway}
}

func (cc CloudControllerAppRepository) SetEnv(app App, kvs map[string]interface{}) (apiErr error) {
	data, err := json.Marshal(kvs)
	if err != nil {
		apiErr = fmt.Errorf("Can not serilize the data")
		return
	}

	_, err = cc.gateway.Request("POST", fmt.Sprintf("/apps/%s/env", app.Id()), data)
	return err
}

func (cc CloudControllerAppRepository) UnsetEnv(app App, keys []string) (apiErr error) {
	env := make(map[string]interface{})
	env["envs"] = keys
	apiErr = cc.gateway.PUT(fmt.Sprintf("/apps/%s/env", app.Id()), env)
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
	if id == "" {
		return nil, errors.New("Application not found")
	}

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

func (cc CloudControllerAppRepository) SwitchStack(id string, params UpdateStackParams) (apiErr error) {
	data, err := json.Marshal(params)
	if err != nil {
		apiErr = err
		return
	}

	_, apiErr = cc.gateway.Request("PUT", fmt.Sprintf("/apps/%s/switch-stack", id), data)
	return
}

func (cc CloudControllerAppRepository) Delete(id string) (apiErr error) {
	_, apiErr = cc.gateway.Request("DELETE", fmt.Sprintf("/apps/%s", id), nil)
	return
}

func (cc CloudControllerAppRepository) GetRoutes(app App) (routes AppRoutes, apiErr error) {
	var routesModel AppRoutesModel
	apiErr = cc.gateway.Get(fmt.Sprintf("/apps/"+app.Id()+"/routes"), &routesModel)
	routesModel.AppRepo = cc
	routes = routesModel
	return
}

func (cc CloudControllerAppRepository) GetRoutesByURI(uri string) (routes AppRoutes, apiErr error) {
	var routesModel AppRoutesModel
	apiErr = cc.gateway.Get(uri, &routesModel)
	routesModel.AppRepo = cc
	routes = routesModel
	return
}

func (cc CloudControllerAppRepository) GetLog(appId, buildId, logType string, lines int64, offset int64) (logs LogsModel, err error) {
	var logsModel LogsModel
	err = cc.gateway.Get(fmt.Sprintf("/apps/%s/builds/%s/log?lines=%d&log_type=%s&offset=%d", appId, buildId, lines, logType, offset), &logsModel)
	logs = logsModel
	return
}

func (cc CloudControllerAppRepository) GetPermission(app App, userId string) (permission AppPermission, err error) {
	var appPermission AppPermission
	err = cc.gateway.Get(fmt.Sprintf("/apps/%s/permissions?user=%s", app.Id(), userId), &appPermission)
	permission = appPermission
	return
}

func (cc CloudControllerAppRepository) GetCollaborators(appId string) (users []UserModel, err error) {
	err = cc.gateway.Get(fmt.Sprintf("/apps/%s/collaborators", appId), &users)
	if err != nil {
		return nil, err
	}
	return
}

func (cc CloudControllerAppRepository) AddCollaborator(appId string, param CreateCollaboratorParams) (err error) {
	data, err := json.Marshal(param)
	if err != nil {
		err = fmt.Errorf("Can not serilize the data")
		return err
	}

	_, err = cc.gateway.Request("POST", fmt.Sprintf("/apps/%s/collaborators", appId), data)

	return err
}

func (cc CloudControllerAppRepository) RemoveCollaborator(appId string, userId string) (err error) {
	_, err = cc.gateway.Request("DELETE", fmt.Sprintf("/apps/%s/collaborators/%s", appId, userId), nil)
	return
}

func (cc CloudControllerAppRepository) TransferToUser(appId string, userEmail string) (err error) {
	params := TransferAppParams{
		Owner:     userEmail,
		OwnerType: "User",
	}
	data, err := json.Marshal(params)
	if err != nil {
		return
	}
	_, err = cc.gateway.Request("PUT", fmt.Sprintf("/apps/%s/transferred", appId), data)

	return err
}

func (cc CloudControllerAppRepository) TransferToOrg(appId string, orgName string) (err error) {
	params := TransferAppParams{
		Owner:     orgName,
		OwnerType: "Organization",
	}
	data, err := json.Marshal(params)
	if err != nil {
		return
	}
	_, err = cc.gateway.Request("PUT", fmt.Sprintf("/apps/%s/transferred", appId), data)

	return err
}
