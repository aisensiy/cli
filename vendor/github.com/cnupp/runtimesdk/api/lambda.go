package api

import (
	"encoding/json"
	. "github.com/cnupp/appssdk/api"
)

type Lambda interface {
	ID() string
	GitSha() string
	Status() string
	Links() Links
	Destroy() error
	IsRunning() bool
	GetService(string) LauncherService
	ToJson() string
	Progress() Progress
}

type LambdaParams struct {
	App     string `json:"app"`
	Build   string `json:"build"`
	Cluster string `json:"cluster"`
}

type Endpoint struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type Progress struct {
	Status string `json:"status"`
}

type LambdaModel struct {
	IDField     string                  `json:"id"`
	GitShaField string                  `json:"releaseVersion"`
	StatusField string                  `json:"status"`
	ProgressField Progress                  `json:"progress"`
	Services    map[string]ServiceModel `json:"services"`
	LinksArray  []Link                  `json:"links"`
	LambdaRepo  LambdaRepository        `json:"-"`
}

func (dm LambdaModel) ID() string {
	return dm.IDField
}

func (p Progress) IsSucceed() bool {
	return p.Status == "SUCCEED"
}

func (dm LambdaModel) Status() string {
	return dm.StatusField
}

func (dm LambdaModel) Links() Links {
	return LinksModel{
		Links: dm.LinksArray,
	}
}

func (dm LambdaModel) Destroy() error {
	return dm.LambdaRepo.Destroy(dm.ID())
}

func (dm LambdaModel) GitSha() string {
	return dm.GitShaField
}

func (dm LambdaModel) IsRunning() bool {
	return dm.StatusField == "RUNNING"
}

func (dm LambdaModel) GetService(name string) LauncherService {
	if val, ok := dm.Services[name]; ok {
		return val
	}
	return nil
}

func (dm LambdaModel) ToJson() string {
	content, err := json.Marshal(dm)
	if err != nil {
		return "{}"
	}

	return string(content)
}

func (dm LambdaModel) Progress() Progress {
	return dm.ProgressField
}
