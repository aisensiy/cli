package api

import (
	"fmt"
	"github.com/cde/client/Godeps/_workspace/src/github.com/cde/apisdk/util"
)

type AppParams struct {
	Stack     string `json:"stackId"`
	Name      string `json:"name"`
	Mem       int    `json:"memeory"`
	Disk      int    `json:"disk"`
	Instances int    `json:"instances"`
}

type AppRouteParams struct {
	Route string `json:"route"`
}

type App interface {
	Id() string
	Mem() int
	Disk() int
	Instances() int
	Links() Links
	GetBuilds() (Builds, error)
	GetRoutes() (AppRoutes, error)
	GetBuild(id string) (Build, error)
	GetBuildByURI(uri string) (Build, error)
	GetStack() (Stack, error)
	GetEnvs() map[string]string
	SetEnv(key string, value string) error
	UnsetEnv(key string) error
	CreateBuild(buildParams BuildParams) (Build, error)
	BindWithRoute(params AppRouteParams) error
	UnbindRoute(routeId string) error
}

type AppModel struct {
	ID              string            `json:"name"`
	Memory          int               `json:"memory"`
	DiskQuota       int               `json:"disk"`
	InstancesCount  int               `json:"instances"`
	Envs            map[string]string `json:"envs"`
	LinksArray      []Link            `json:"links"`
	BuildMapper     BuildMapper
	AppMapper       AppRepository
	StackRepository StackRepository
}

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (a AppModel) GetEnvs() map[string]string {
	return a.Envs
}

func (a AppModel) SetEnv(key string, value string) (err error) {
	err = a.AppMapper.SetEnv(a, KeyValue{Key: key, Value: value})
	return
}

func (a AppModel) UnsetEnv(key string) (err error) {
	err = a.AppMapper.UnsetEnv(a, key)
	return
}

func (a AppModel) Id() string {
	return a.ID
}

func (a AppModel) Mem() int {
	return a.Memory
}

func (a AppModel) Disk() int {
	return a.DiskQuota
}

func (a AppModel) Instances() int {
	return a.InstancesCount
}
func (a AppModel) Links() Links {
	return LinksModel{
		Links: a.LinksArray,
	}
}

func (a AppModel) GetBuilds() (builds Builds, apiError error) {
	return a.BuildMapper.GetBuilds(a)
}

func (a AppModel) GetBuildByURI(uri string) (build Build, apiError error) {
	id, err := util.IDFromURI(uri)
	fmt.Printf("fdasfdas%s ", id)
	if err != nil {
		apiError = err
		return
	}
	return a.BuildMapper.GetBuild(a, id)
}

func (a AppModel) GetBuild(id string) (build Build, apiError error) {
	return a.BuildMapper.GetBuild(a, id)
}

func (a AppModel) CreateBuild(buildParams BuildParams) (build Build, apiErr error) {
	return a.BuildMapper.Create(a, buildParams)
}

func (a AppModel) GetStack() (stack Stack, apiErr error) {
	stackLink, err := a.Links().Link("stack")
	if err != nil {
		apiErr = err
		return
	}

	return a.StackRepository.GetStackByURI(stackLink.URI)
}

func (a AppModel) BindWithRoute(params AppRouteParams) error {
	return a.AppMapper.BindWithRoute(a, params)
}

func (a AppModel) UnbindRoute(routeId string) error {
	return a.AppMapper.UnbindRoute(a, routeId)
}

func (a AppModel) GetRoutes() (AppRoutes, error) {
	return a.AppMapper.GetRoutes(a)
}

type AppRef interface {
	Id() string
	Links() Links
}

type AppRefModel struct {
	IDField     string `json:"name"`
	LinksField  []Link `json:"links"`
	BuildMapper BuildMapper
}

func (arm AppRefModel) Id() string {
	return arm.IDField
}

func (arm AppRefModel) Links() Links {
	return LinksModel{
		Links: arm.LinksField,
	}
}

type Apps interface {
	Count() int
	First() Apps
	Last() Apps
	Prev() Apps
	Next() Apps
	Items() []AppRef
}

type AppsModel struct {
	CountField int           `json:"count"`
	SelfField  string        `json:"self"`
	FirstField string        `json:"first"`
	LastField  string        `json:"last"`
	PrevField  string        `json:"prev"`
	NextField  string        `json:"next"`
	ItemsField []AppRefModel `json:"items"`
	AppMapper  AppRepository
}

func (apps AppsModel) Count() int {
	return apps.CountField
}
func (apps AppsModel) Self() Apps {
	return nil
}
func (apps AppsModel) First() Apps {
	return nil
}
func (apps AppsModel) Last() Apps {
	return nil
}
func (apps AppsModel) Prev() Apps {
	return nil
}
func (apps AppsModel) Next() Apps {
	return nil
}

func (apps AppsModel) Items() []AppRef {
	items := make([]AppRef, 0)
	for _, app := range apps.ItemsField {
		items = append(items, app)
	}
	return items
}

type AppRouteModel struct {
	IDField      string       `json:"id"`
	PathField    string       `json:"path"`
	DomainField  SimpleDomain `json:"domain"`
	CreatedField string       `json:"created"`
	LinksArray   []Link       `json:"links"`
}

type AppRoutes interface {
	Count() int
	First() Apps
	Last() Apps
	Prev() Apps
	Next() Apps
	Items() []AppRouteModel
}

type AppRoutesModel struct {
	CountField int             `json:"count"`
	SelfField  string          `json:"self"`
	FirstField string          `json:"first"`
	LastField  string          `json:"last"`
	PrevField  string          `json:"prev"`
	NextField  string          `json:"next"`
	ItemsField []AppRouteModel `json:"items"`
}

func (appRoutes AppRoutesModel) Count() int {
	return appRoutes.CountField
}
func (appRoutes AppRoutesModel) Self() Apps {
	return nil
}
func (appRoutes AppRoutesModel) First() Apps {
	return nil
}
func (appRoutes AppRoutesModel) Last() Apps {
	return nil
}
func (appRoutes AppRoutesModel) Prev() Apps {
	return nil
}
func (appRoutes AppRoutesModel) Next() Apps {
	return nil
}

func (appRoutes AppRoutesModel) Items() []AppRouteModel {
	items := make([]AppRouteModel, 0)
	for _, app := range appRoutes.ItemsField {
		items = append(items, app)
	}
	return items
}
