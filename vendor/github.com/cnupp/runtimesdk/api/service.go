package api

import (
	"encoding/json"
	"github.com/cnupp/appssdk/api"
)

type Task interface {
	Port() int
	Host() string
}

type TaskModel struct {
	PortField int    `json:"port"`
	HostField string `json:"host"`
}

func (t TaskModel) Port() int {
	return t.PortField
}

func (t TaskModel) Host() string {
	return t.HostField
}

type LauncherService interface {
	Id() string
	Env() string
	Instance() int
	CPU() float32
	Memory() float32
	Tasks() []Task
	Update(params ServiceConfigParams) error
	Name() string
	Log(lines int) (api.LogsModel, error)
	InternalEndpoint() (Endpoint, error)
}

type Endpoints struct {
	Internal Endpoint `json:"internal"`
}

type ServiceModel struct {
	IdField           string                 `json:"id"`
	NameField         string                 `json:"name"`
	EnvField          map[string]string      `json:"env"`
	InstanceField     int                    `json:"instances"`
	CPUField          float32                `json:"cpus"`
	MemoryField       float32                `json:"mem"`
	DiskField         int                    `json:"disk"`
	PortsField        []int                  `json:"ports"`
	ContainerField    map[string]interface{} `json:"container"`
	DependenciesField []string               `json:"dependencies"`
	TasksField        []TaskModel            `json:"tasks"`
	AppName           string                 `json:"-"`
	Endpoints         Endpoints              `json:"endpoint"`
	repo              DeploymentRepository   `json:"-"`
	serviceMapper     ServiceMapper          `json:"-"`
}

func (s ServiceModel) Id() string {
	return s.IdField
}
func (s ServiceModel) Env() string {
	content, _ := json.Marshal(s.EnvField)
	return string(content)
}
func (s ServiceModel) Instance() int {
	return s.InstanceField
}
func (s ServiceModel) CPU() float32 {
	return s.CPUField
}
func (s ServiceModel) Memory() float32 {
	return s.MemoryField
}
func (s ServiceModel) Name() string {
	return s.NameField
}
func (s ServiceModel) Tasks() []Task {
	tasks := make([]Task, 0)
	for _, task := range s.TasksField {
		tasks = append(tasks, task)
	}
	return tasks
}
func (s ServiceModel) Log(lines int) (api.LogsModel, error) {
	return s.serviceMapper.Log(s.AppName, s.Name(), lines)
}

func (s ServiceModel) Update(params ServiceConfigParams) error {
	return s.repo.ScaleServiceForApp(s.AppName, s.Name(), params)
}

func (s ServiceModel) InternalEndpoint() (Endpoint, error) {
	return s.Endpoints.Internal, nil
}

type ServiceConfigParams struct {
	Instance int     `json:"instances"`
	Memory   float32 `json:"mem"`
	CPUS     float32 `json:"cpus"`
}
