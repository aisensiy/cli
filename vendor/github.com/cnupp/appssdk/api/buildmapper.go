package api

import (
	"encoding/json"
	"fmt"

	"github.com/cnupp/appssdk/config"
	"github.com/cnupp/appssdk/net"
)

//go:generate counterfeiter -o fakes/fake_build_mapper.go . BuildMapper
type BuildMapper interface {
	Create(app App, params BuildParams) (build Build, apiErr error)
	GetBuilds(app App) (builds Builds, apiErr error)
	GetBuild(app App, id string) (build Build, apiErr error)
	Update(id string, params BuildParams) (updatedBuild Build, apiErr error)
	Success(build Build) (apiErr error)
	Fail(build Build) (apiErr error)
	VerifySuccess(build Build) (apiErr error)
	VerifyFail(build Build) (apiErr error)
	CreateVerify(build Build, params VerifyParams) (verify Verify, apiErr error)
	GetVerify(app App, build Build, id string) (verify Verify, apiErr error)
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

	res, err := bm.gateway.Request("POST", fmt.Sprintf("/apps/%s/builds", app.Name()), data)
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
	apiErr = bm.gateway.Get(fmt.Sprintf("/apps/%s/builds", app.Name()), &buildsModel)
	if apiErr != nil {
		return
	}
	buildsModel.BuildMapper = NewBuildMapper(bm.config, bm.gateway)
	builds = buildsModel
	return
}

func (bm DefaultBuildMapper) GetBuild(app App, id string) (build Build, apiErr error) {
	var buildModel BuildModel
	apiErr = bm.gateway.Get(fmt.Sprintf("/apps/%s/builds/%s", app.Name(), id), &buildModel)
	if apiErr != nil {
		return
	}
	buildModel.AppField = app
	buildModel.BuildMapper = NewBuildMapper(bm.config, bm.gateway)
	buildModel.Resource = NewResource(bm.config, bm.gateway)
	build = buildModel
	return
}

func (bm DefaultBuildMapper) GetVerify(app App, build Build, id string) (verify Verify, apiErr error) {
	var verifyModel VerifyModel
	apiErr = bm.gateway.Get(fmt.Sprintf("/apps/%s/builds/%s/verifies/%s", app.Name(), build.Id(), id), &verifyModel)
	if apiErr != nil {
		return
	}
	verifyModel.BuildMapper = NewBuildMapper(bm.config, bm.gateway)
	verifyModel.Resource = NewResource(bm.config, bm.gateway)
	verifyModel.BuildField = build
	verify = verifyModel
	return
}

func (bm DefaultBuildMapper) Success(build Build) error {
	return bm.gateway.PUT(fmt.Sprintf("/apps/%s/builds/%s/success", build.GetApp().Name(), build.Id()), nil)
}

func (bm DefaultBuildMapper) Fail(build Build) (apiErr error) {
	return bm.gateway.PUT(fmt.Sprintf("/apps/%s/builds/%s/fail", build.GetApp().Name(), build.Id()), nil)
}

func (bm DefaultBuildMapper) Update(id string, params BuildParams) (updatedApp Build, apiErr error) {
	return
}

func (bm DefaultBuildMapper) VerifySuccess(build Build) (apiErr error) {
	return bm.gateway.PUT(fmt.Sprintf("/apps/%s/builds/%s/verify/success", build.GetApp().Name(), build.Id()), nil)
}

func (bm DefaultBuildMapper) VerifyFail(build Build) (apiErr error) {
	return bm.gateway.PUT(fmt.Sprintf("/apps/%s/builds/%s/verify/fail", build.GetApp().Name(), build.Id()), nil)
}

func (bm DefaultBuildMapper) CreateVerify(build Build, params VerifyParams) (verify Verify, apiErr error) {
	data, err := json.Marshal(params)
	if err != nil {
		apiErr = fmt.Errorf("Can not serilize the data")
		return
	}
	url := fmt.Sprintf("/apps/%s/builds/%s/verifies", build.GetApp().Name(), build.Id())
	res, err := bm.gateway.Request("POST", url, data)
	if err != nil {
		fmt.Println("error hanppend when request ", url, err)
		apiErr = err
		return
	}

	location := res.Header.Get("Location")

	var createdVerify VerifyModel
	apiErr = bm.gateway.Get(location, &createdVerify)
	if apiErr != nil {
		return
	}

	createdVerify.BuildField = build
	createdVerify.BuildMapper = NewBuildMapper(bm.config, bm.gateway)
	createdVerify.Resource = NewResource(bm.config, bm.gateway)
	verify = createdVerify

	return
}
