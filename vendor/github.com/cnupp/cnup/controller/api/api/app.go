package api

import (
	"github.com/cnupp/cnup/controller/api/util"
)

type AppParams struct {
	Stack            string `json:"stackId"`
	UnifiedProcedure string `json:"unified_procedure_id"`
	Provider         string `json:"provider_uri"`
	Owner            string `json:"owner"`
	Name             string `json:"name"`
	NeedDeploy       bool   `json:"needDeploy"`
}

type AppRouteParams struct {
	Route string `json:"route"`
}

type UpdateStackParams struct {
	Stack string `json:"stack"`
}

type CreateCollaboratorParams struct {
	Email string `json:"email"`
}

type Cluster interface {
	Links() Links
	Name() string
	Endpoint() string
	Type() string
}

type ClusterModel struct {
	EndpointField string `json:"uri"`
	NameField     string `json:"name"`
	TypeField     string `json:"type"`
	LinksArray    []Link `json:"links"`
}

func (cm ClusterModel) Name() string {
	return cm.NameField
}

func (cm ClusterModel) Type() string {
	return cm.TypeField
}

func (cm ClusterModel) Endpoint() string {
	return cm.EndpointField
}

func (cm ClusterModel) Links() Links {
	return LinksModel{
		Links: cm.LinksArray,
	}
}

type App interface {
	Id() string
	Links() Links
	GetCluster() (Cluster, error)
	GetBuilds() (Builds, error)
	GetRoutes() (AppRoutes, error)
	GetBuild(id string) (Build, error)
	GetBuildByURI(uri string) (Build, error)
	GetStack() (Stack, error)
	GetEnvs() map[string]string
	SetEnv(envs map[string]interface{}) error
	UnsetEnv(keys []string) error
	CreateBuild(buildParams BuildParams) (Build, error)
	BindWithRoute(params AppRouteParams) error
	UnbindRoute(routeId string) error
	SwitchStack(params UpdateStackParams) error
	GetLogForTests(buildId, logType string, lines int64, offset int64) (LogsModel, error)
	GetPermissions(userId string) (AppPermission, error)
	GetCollaborators() ([]UserModel, error)
	AddCollaborator(param CreateCollaboratorParams) error
	RemoveCollaborator(userId string) error
	TransferToOrg(orgName string) error
	TransferToUser(userEmail string) error
	NeedDeploy() bool
}

type AppModel struct {
	ID              string            `json:"name"`
	NeedDeployField bool              `json:"needDeploy"`
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

func (a AppModel) SetEnv(envs map[string]interface{}) (err error) {
	err = a.AppMapper.SetEnv(a, envs)
	return
}

func (a AppModel) UnsetEnv(keys []string) (err error) {
	err = a.AppMapper.UnsetEnv(a, keys)
	return
}

func (a AppModel) Id() string {
	return a.ID
}

func (a AppModel) NeedDeploy() bool {
	return a.NeedDeployField
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

func (a AppModel) SwitchStack(params UpdateStackParams) error {
	return a.AppMapper.SwitchStack(a.ID, params)
}

func (a AppModel) GetLogForTests(buildId, logType string, lines int64, offset int64) (LogsModel, error) {
	return a.AppMapper.GetLog(a.ID, buildId, logType, lines, offset)
}

func (a AppModel) GetCluster() (Cluster, error) {
	return ClusterModel{
		EndpointField: "/clusters/1",
	}, nil
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
	First() (routes AppRoutes, apiError error)
	Last() (routes AppRoutes, apiError error)
	Prev() (routes AppRoutes, apiError error)
	Next() (routes AppRoutes, apiError error)
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
	AppRepo    AppRepository
}

type LogItemsModel struct {
	MessageField string `json:"message"`
}

type LogsModel struct {
	ErrorField string          `json:"error"`
	TotalField int64           `json:"total"`
	SizeField  int64           `json:"size"`
	ItemsField []LogItemsModel `json:"items"`
}

func (appRoutes AppRoutesModel) Count() int {
	return appRoutes.CountField
}
func (appRoutes AppRoutesModel) Self() (routes AppRoutes, apiError error) {
	if "" == appRoutes.SelfField {
		return
	}

	routes, apiError = appRoutes.AppRepo.GetRoutesByURI(appRoutes.SelfField)
	return
}
func (appRoutes AppRoutesModel) First() (routes AppRoutes, apiError error) {
	if "" == appRoutes.FirstField {
		return
	}

	routes, apiError = appRoutes.AppRepo.GetRoutesByURI(appRoutes.FirstField)
	return
}
func (appRoutes AppRoutesModel) Last() (routes AppRoutes, apiError error) {
	if "" == appRoutes.LastField {
		return
	}

	routes, apiError = appRoutes.AppRepo.GetRoutesByURI(appRoutes.LastField)
	return
}
func (appRoutes AppRoutesModel) Prev() (routes AppRoutes, apiError error) {
	if "" == appRoutes.PrevField {
		return
	}

	routes, apiError = appRoutes.AppRepo.GetRoutesByURI(appRoutes.PrevField)
	return
}
func (appRoutes AppRoutesModel) Next() (routes AppRoutes, apiError error) {
	if "" == appRoutes.NextField {
		return
	}

	routes, apiError = appRoutes.AppRepo.GetRoutesByURI(appRoutes.NextField)
	return
}

func (appRoutes AppRoutesModel) Items() []AppRouteModel {
	items := make([]AppRouteModel, 0)
	for _, app := range appRoutes.ItemsField {
		items = append(items, app)
	}
	return items
}

func (app AppModel) GetPermissions(userId string) (AppPermission, error) {
	appPermission, err := app.AppMapper.GetPermission(app, userId)
	return appPermission, err
}

func (app AppModel) GetCollaborators() ([]UserModel, error) {
	users, err := app.AppMapper.GetCollaborators(app.Id())
	return users, err
}

func (app AppModel) AddCollaborator(param CreateCollaboratorParams) error {
	return app.AppMapper.AddCollaborator(app.Id(), param)
}

func (app AppModel) RemoveCollaborator(userId string) error {
	return app.AppMapper.RemoveCollaborator(app.Id(), userId)
}

func (app AppModel) TransferToOrg(orgName string) error {
	return app.AppMapper.TransferToOrg(app.Id(), orgName)
}

func (app AppModel) TransferToUser(userEmail string) error {
	return app.AppMapper.TransferToUser(app.Id(), userEmail)
}
