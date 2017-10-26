package api

import (
	"github.com/sjkyspa/stacks/controller/api/api"
)

type Deployment interface {
	ID() string
	Version() string
	Status() string
	MarathonApp() string
	Links() api.Links
	Destroy() error
	Restart() error
	GetService(serviceName string) (LauncherService, error)
	GetDependentServices() ([]LauncherService, error)
	Log(lines int) (api.LogsModel, error)
}

type DeploymentParams struct {
	App     string `json:"app"`
	Release string `json:"release"`
	Cluster string `json:"cluster"`
}

type DeploymentModel struct {
	IDField          string               `json:"id"`
	VersionField     string               `json:"releaseVersion"`
	MarathonAppField string               `json:"marathonApp"`
	StatusField      string               `json:"status"`
	LinksArray       []api.Link           `json:"links"`
	repo             DeploymentRepository `json:"-"`
	appName          string               `json:"-"`
}

func (dm DeploymentModel) ID() string {
	return dm.IDField
}

func (dm DeploymentModel) Version() string {
	return dm.VersionField
}

func (dm DeploymentModel) Status() string {
	return dm.StatusField
}

func (dm DeploymentModel) MarathonApp() string {
	return dm.MarathonAppField
}

func (dm DeploymentModel) Links() api.Links {
	return api.LinksModel{
		Links: dm.LinksArray,
	}
}

func (dm DeploymentModel) Destroy() error {
	return dm.repo.Destroy(dm.appName)
}

func (dm DeploymentModel) Restart() error {
	return dm.repo.Restart(dm.appName)
}

func (dm DeploymentModel) GetService(serviceName string) (LauncherService, error) {
	return dm.repo.GetDependentServiceForApp(dm.appName, serviceName)
}

func (dm DeploymentModel) GetDependentServices() ([]LauncherService, error) {
	return dm.repo.GetDependentServicesForApp(dm.appName)
}

func (dm DeploymentModel) Log(lines int) (api.LogsModel, error) {
	return dm.repo.Log(dm.appName, lines)
}
