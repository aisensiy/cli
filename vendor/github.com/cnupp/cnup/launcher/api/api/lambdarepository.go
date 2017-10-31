package api

import (
	"encoding/json"
	"fmt"
	"github.com/cnupp/cnup/launcher/api/config"
	"github.com/cnupp/cnup/launcher/api/net"
	"time"
)

type LambdaRepository interface {
	Create(params LambdaParams) (Lambda, error)
	GetLambda(id string) (lambda Lambda, apiErr error)
	GetLambdaByURI(url string) (lambda Lambda, apiErr error)
	Destroy(appName string) error
}

func NewLambdaRepository(config config.Reader, gateway net.Gateway) LambdaRepository {
	return DefaultLambdaRepository{
		config:  config,
		gateway: gateway,
	}
}

type DefaultLambdaRepository struct {
	config  config.Reader
	gateway net.Gateway
}

func (ddr DefaultLambdaRepository) GetLambdaByURI(uri string) (lambda Lambda, apiErr error) {
	var lambdaModel LambdaModel
	apiErr = ddr.gateway.Get(uri, &lambdaModel)
	lambdaModel.LambdaRepo = ddr
	lambda = lambdaModel
	return
}

func (ddr DefaultLambdaRepository) Create(params LambdaParams) (lambda Lambda, apiErr error) {
	data, err := json.Marshal(params)
	if err != nil {
		apiErr = err
		return
	}

	res, err := ddr.gateway.Request("POST", "/lambdas", data)
	if err != nil {
		apiErr = err
		return
	}

	location, err := res.Location()
	if err != nil {
		apiErr = err
		return
	}

	time.Sleep(time.Second * 5)

	var lambdaModel LambdaModel
	lambdaModel.LambdaRepo = ddr
	err = ddr.gateway.Get(location.String(), &lambdaModel)
	if err != nil {
		apiErr = err
		return
	}

	lambda = lambdaModel
	return
}

func (ddr DefaultLambdaRepository) GetLambda(id string) (lambda Lambda, apiErr error) {
	var lambdaModel LambdaModel
	apiErr = ddr.gateway.Get(fmt.Sprintf("/lambdas/%s", id), &lambdaModel)
	if apiErr != nil {
		return
	}
	lambdaModel.LambdaRepo = ddr
	lambdaModel.IDField = id
	lambda = lambdaModel
	return
}

func (ddr DefaultLambdaRepository) Destroy(appName string) (apiErr error) {
	apiErr = ddr.gateway.Delete(fmt.Sprintf("/lambdas/%s", appName), nil)
	return
}
