package api

import (
	"encoding/json"
	"fmt"
	"github.com/cde/client/Godeps/_workspace/src/github.com/cde/apisdk/config"
	"github.com/cde/client/Godeps/_workspace/src/github.com/cde/apisdk/net"
)

//go:generate counterfeiter -o fakes/fake_build_mapper.go . BuildMapper
type BuildMapper interface {
	Create(app App, params BuildParams) (build Build, apiErr error)
	GetBuilds(app App) (builds Builds, apiErr error)
	GetBuild(app App, id string) (build Build, apiErr error)
	Update(id string, params BuildParams) (updatedBuild Build, apiErr error)
	Success(build Build) (apiErr error)
	Fail(build Build) (apiErr error)
}

type DefaultBuildMapper struct {
	config  config.Reader
	gateway net.Gateway
}

func NewBuildMapper(reader config.Reader, gateway net.Gateway) BuildMapper {
	return DefaultBuildMapper{
		config:  reader,
		gateway: gateway,
	}
}

func (bm DefaultBuildMapper) Create(app App, params BuildParams) (build Build, apiErr error) {
	data, err := json.Marshal(params)
	if err != nil {
		apiErr = fmt.Errorf("Can not serilize the data")
		return
	}

	res, err := bm.gateway.Request("POST", fmt.Sprintf("/apps/%s/builds", app.Id()), data)
	if err != nil {
		apiErr = err
		return
	}

	location := res.Header.Get("Location")

	var createdBuild BuildModel
	apiErr = bm.gateway.Get(location, &createdBuild)
	if apiErr != nil {
		return
	}

	createdBuild.AppField = app
	createdBuild.BuildMapper = NewBuildMapper(bm.config, bm.gateway)
	build = createdBuild

	return
}

func (bm DefaultBuildMapper) GetBuilds(app App) (builds Builds, apiErr error) {
	var buildsModel BuildsModel
	apiErr = bm.gateway.Get(fmt.Sprintf("/apps/%s/builds", app.Id()), &buildsModel)
	if apiErr != nil {
		return
	}
	buildsModel.BuildMapper = NewBuildMapper(bm.config, bm.gateway)
	builds = buildsModel
	return
}

func (bm DefaultBuildMapper) GetBuild(app App, id string) (build Build, apiErr error) {
	var buildModel BuildModel
	apiErr = bm.gateway.Get(fmt.Sprintf("/apps/%s/builds/%s", app.Id(), id), &buildModel)
	if apiErr != nil {
		return
	}
	buildModel.AppField = app
	buildModel.BuildMapper = NewBuildMapper(bm.config, bm.gateway)
	build = buildModel
	return
}

func (bm DefaultBuildMapper) Success(build Build) (apiErr error) {
	return bm.gateway.PUT(fmt.Sprintf("/apps/%s/builds/%s/success", build.GetApp().Id(), build.Id()), nil)
}

func (bm DefaultBuildMapper) Fail(build Build) (apiErr error) {
	return bm.gateway.PUT(fmt.Sprintf("/apps/%s/builds/%s/fail", build.GetApp().Id(), build.Id()), nil)
}

func (bm DefaultBuildMapper) Update(id string, params BuildParams) (updatedApp Build, apiErr error) {
	return
}
