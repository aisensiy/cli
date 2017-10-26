package api

import (
	"fmt"
	"github.com/sjkyspa/stacks/controller/api/api"
	"github.com/sjkyspa/stacks/launcher/api/config"
	"github.com/sjkyspa/stacks/launcher/api/net"
)

type ServiceMapper interface {
	Log(appName, serviceName string, lines int) (api.LogsModel, error)
}

func NewServiceMapper(config config.Reader, gateway net.Gateway) ServiceMapper {
	return DefaultServiceMapper{
		config:  config,
		gateway: gateway,
	}
}

type DefaultServiceMapper struct {
	config  config.Reader
	gateway net.Gateway
}

func (dsm DefaultServiceMapper) Log(appName, serviceName string, lines int) (output api.LogsModel, apiErr error) {
	var logModel api.LogsModel
	apiErr = dsm.gateway.Get(fmt.Sprintf("/deployments/%s/logs/%s?log_lines=%d", appName, serviceName, lines), &logModel)
	output = logModel
	return
}
