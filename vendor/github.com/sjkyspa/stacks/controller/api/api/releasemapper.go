package api

import (
	"fmt"
	"github.com/sjkyspa/stacks/controller/api/config"
	"github.com/sjkyspa/stacks/controller/api/net"
)

type ReleaseMapper interface {
	Create(app App) (release Release, apiErr error)
	GetReleases(app App) (releases Releases, apiErr error)
	GetRelease(app App, id string) (release Release, apiErr error)
	Success(release Release) (apiErr error)
	Fail(release Release) (apiErr error)
}

type DefaultReleaseMapper struct {
	config  config.Reader
	gateway net.Gateway
}

func NewReleaseMapper(reader config.Reader, gateway net.Gateway) ReleaseMapper {
	return DefaultReleaseMapper{
		config:  reader,
		gateway: gateway,
	}
}

func (bm DefaultReleaseMapper) Create(app App) (release Release, apiErr error) {
	res, err := bm.gateway.Request("POST", fmt.Sprintf("/apps/%s/releases", app.Id()), []byte("{}"))
	if err != nil {
		apiErr = err
		return
	}

	location := res.Header.Get("Location")

	var createdRelease ReleaseModel
	apiErr = bm.gateway.Get(location, &createdRelease)
	if apiErr != nil {
		return
	}

	createdRelease.AppField = app
	createdRelease.ReleaseMapper = NewReleaseMapper(bm.config, bm.gateway)
	release = createdRelease

	return
}

func (bm DefaultReleaseMapper) GetReleases(app App) (releases Releases, apiErr error) {
	var releasesModel ReleasesModel
	apiErr = bm.gateway.Get(fmt.Sprintf("/apps/%s/releases", app.Id()), &releasesModel)
	if apiErr != nil {
		return
	}
	releasesModel.ReleaseMapper = NewReleaseMapper(bm.config, bm.gateway)
	releases = releasesModel
	return
}

func (bm DefaultReleaseMapper) GetRelease(app App, id string) (release Release, apiErr error) {
	var releaseModel ReleaseModel
	apiErr = bm.gateway.Get(fmt.Sprintf("/apps/%s/releases/%s", app.Id(), id), &releaseModel)
	if apiErr != nil {
		return
	}
	releaseModel.AppField = app
	releaseModel.ReleaseMapper = NewReleaseMapper(bm.config, bm.gateway)
	release = releaseModel
	return
}

func (bm DefaultReleaseMapper) Success(release Release) (apiErr error) {
	return bm.gateway.PUT(fmt.Sprintf("/apps/%s/releases/%s/success", release.GetApp().Id(), release.Id()), nil)
}

func (bm DefaultReleaseMapper) Fail(release Release) (apiErr error) {
	return bm.gateway.PUT(fmt.Sprintf("/apps/%s/releases/%s/fail", release.GetApp().Id(), release.Id()), nil)
}

func (bm DefaultReleaseMapper) Update(id string, params ReleaseParams) (updatedApp Release, apiErr error) {
	return
}
