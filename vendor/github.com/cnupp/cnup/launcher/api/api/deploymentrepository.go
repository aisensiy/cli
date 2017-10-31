package api

import (
	"encoding/json"
	"fmt"
	"github.com/cnupp/cnup/controller/api/api"
	"github.com/cnupp/cnup/launcher/api/config"
	"github.com/cnupp/cnup/launcher/api/net"
)

type DeploymentRepository interface {
	Create(appName string, params DeploymentParams) (Deployment, error)
	GetDeploymentByAppName(appName string) (deployment Deployment, apiErr error)
	GetDependentServicesForApp(appName string) (services []LauncherService, apiErr error)
	GetDependentServiceForApp(appName, serviceId string) (service LauncherService, apiErr error)
	Destroy(appName string) error
	Restart(appName string) error
	ScaleServiceForApp(appName, serviceName string, params ServiceConfigParams) error
	Log(appName string, lines int) (api.LogsModel, error)
}

func NewDeploymentRepository(config config.Reader, gateway net.Gateway) DeploymentRepository {
	return DefaultDeploymentRepository{
		config:  config,
		gateway: gateway,
	}
}

type DefaultDeploymentRepository struct {
	config  config.Reader
	gateway net.Gateway
}

func (ddr DefaultDeploymentRepository) Create(appName string, params DeploymentParams) (deployment Deployment, apiErr error) {
	data, err := json.Marshal(params)
	if err != nil {
		apiErr = err
		return
	}

	res, err := ddr.gateway.Request("PUT", "/deployments/"+appName, data)
	if err != nil {
		apiErr = err
		return
	}

	location, err := res.Location()
	if err != nil {
		apiErr = err
		return
	}

	var deploymentModel DeploymentModel
	err = ddr.gateway.Get(location.String(), &deploymentModel)
	if err != nil {
		apiErr = err
		return
	}

	deployment = deploymentModel
	return
}

func (ddr DefaultDeploymentRepository) GetDeploymentByAppName(appName string) (deployment Deployment, apiErr error) {
	var deploymentModel DeploymentModel
	apiErr = ddr.gateway.Get(fmt.Sprintf("/deployments/%s", appName), &deploymentModel)
	if apiErr != nil {
		return
	}
	deploymentModel.repo = ddr
	deploymentModel.appName = appName
	deployment = deploymentModel
	return
}

func (ddr DefaultDeploymentRepository) GetDependentServicesForApp(appName string) ([]LauncherService, error) {
	var serviceModels []ServiceModel
	apiErr := ddr.gateway.Get(fmt.Sprintf("/deployments/%s/services", appName), &serviceModels)
	if apiErr != nil {
		return nil, apiErr
	}
	services := make([]LauncherService, 0)
	for _, service := range serviceModels {
		service.AppName = appName
		service.repo = ddr
		services = append(services, service)
	}
	return services, nil
}

func (ddr DefaultDeploymentRepository) GetDependentServiceForApp(appName, serviceId string) (service LauncherService, apiErr error) {
	var serviceModel ServiceModel
	apiErr = ddr.gateway.Get(fmt.Sprintf("/deployments/%s/services/%s", appName, serviceId), &serviceModel)
	serviceModel.AppName = appName
	serviceModel.repo = ddr
	serviceModel.serviceMapper = NewServiceMapper(ddr.config, ddr.gateway)
	service = serviceModel
	return
}

func (ddr DefaultDeploymentRepository) Destroy(appName string) (apiErr error) {
	apiErr = ddr.gateway.Delete(fmt.Sprintf("/deployments/%s", appName), nil)
	return
}

func (ddr DefaultDeploymentRepository) Restart(appName string) (apiErr error) {
	apiErr = ddr.gateway.PUT(fmt.Sprintf("/deployments/%s/restart", appName), nil)
	return
}

func (ddr DefaultDeploymentRepository) ScaleServiceForApp(appName, serviceName string, params ServiceConfigParams) (apiErr error) {
	apiErr = ddr.gateway.PUT(fmt.Sprintf("/deployments/%s/services/%s", appName, serviceName), params)
	return
}

func (ddr DefaultDeploymentRepository) Log(appName string, lines int) (log api.LogsModel, apiErr error) {
	var logModel api.LogsModel
	apiErr = ddr.gateway.Get(fmt.Sprintf("/deployments/%s/logs/main?log_lines=%d", appName, lines), &logModel)
	log = logModel
	return
}
