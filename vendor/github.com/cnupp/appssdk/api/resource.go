package api

import (
	"github.com/cnupp/appssdk/config"
	"github.com/cnupp/appssdk/net"
	"regexp"
)

type Resource interface {
	GetResourceByURI(uri string) (interface{}, error)
}

type ResourceModel struct {
	config  config.Reader
	gateway net.Gateway
}

func (rm ResourceModel) GetResourceByURI(uri string) (interface{}, error) {
	types := make(map[string]func() (interface{}, error))

	types["/apps[?]?[^/]*?$"] = func() (interface{}, error) {
		var apps AppsModel
		err := rm.gateway.Get(uri, &apps)
		apps.AppMapper = NewAppRepository(rm.config, rm.gateway)
		return apps, err
	}

	types["/apps/[^/]*?$"] = func() (interface{}, error) {
		var app AppModel
		err := rm.gateway.Get(uri, &app)
		app.BuildMapper = NewBuildMapper(rm.config, rm.gateway)
		app.AppMapper = NewAppRepository(rm.config, rm.gateway)
		return app, err
	}

	types["/apps/[^/]*?/builds[?]?[^/]*?$"] = func() (interface{}, error) {
		var model BuildsModel
		err := rm.gateway.Get(uri, &model)
		model.BuildMapper = NewBuildMapper(rm.config, rm.gateway)
		return model, err
	}

	types["/apps/[^/]*?/builds/[^/]*?$"] = func() (interface{}, error) {
		var model BuildModel
		err := rm.gateway.Get(uri, &model)
		model.BuildMapper = NewBuildMapper(rm.config, rm.gateway)
		model.Resource = NewResource(rm.config, rm.gateway)
		return model, err
	}

	for reg, ty := range types {
		matched, err := regexp.MatchString(reg, uri)
		if matched && err == nil {
			return ty()
		}
	}
	return "", nil
}

func NewResource(cf config.Reader, gw net.Gateway) Resource {
	return ResourceModel{
		config:  cf,
		gateway: gw,
	}
}
