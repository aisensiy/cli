package api

import (
	"encoding/json"
	"fmt"
	"github.com/cnupp/appssdk/config"
	"github.com/cnupp/appssdk/net"
	"io/ioutil"
)

//go:generate counterfeiter -o fakes/fake_stack_repository.go . StackRepository
type StackRepository interface {
	Create(params map[string]interface{}) (createdStack Stack, apiErr error)
	GetStack(id string) (Stack, error)
	GetStackByURI(uri string) (Stack, error)
	GetStacks() (Stacks, error)
	GetStackByName(name string) (Stacks, error)
	Update(id string, params map[string]interface{}) (apiErr error)
	Delete(id string) (apiErr error)
	Publish(id string) (apiErr error)
	UnPublish(id string) (apiErr error)
}

type DefaultStackRepository struct {
	config  config.Reader
	gateway net.Gateway
}

func NewStackRepository(config config.Reader, gateway net.Gateway) StackRepository {
	return DefaultStackRepository{config: config, gateway: gateway}
}

func (cc DefaultStackRepository) Create(params map[string]interface{}) (createdStack Stack, apiErr error) {
	data, err := json.Marshal(params)
	if err != nil {
		apiErr = fmt.Errorf("Can not serilize the data")
		return
	}

	res, err := cc.gateway.Request("POST", "/stacks", data)
	if err != nil {
		apiErr = err
		return
	}

	location := res.Header.Get("Location")
	var stackModel StackModel
	apiErr = cc.gateway.Get(location, &stackModel)
	stackModel.StackMapper = cc
	if apiErr != nil {
		return
	}
	createdStack = stackModel

	return
}

func (cc DefaultStackRepository) GetStack(id string) (Stack, error) {
	var data []byte

	res, err := cc.gateway.Request("GET", fmt.Sprintf("/stacks/%s", id), nil)
	if err != nil {
		return nil, err
	}
	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	stackModel := StackModel{}
	json.Unmarshal(data, &stackModel)
	stackModel.StackMapper = cc
	return stackModel, nil
}

func (cc DefaultStackRepository) GetStackByURI(uri string) (Stack, error) {
	var data []byte
	res, err := cc.gateway.Request("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	stackModel := StackModel{}
	json.Unmarshal(data, &stackModel)
	return stackModel, nil
}

func (cc DefaultStackRepository) GetStacks() (Stacks, error) {
	var stacks StacksModel
	apiErr := cc.gateway.Get(fmt.Sprintf("/stacks"), &stacks)
	if apiErr != nil {
		return nil, apiErr
	}
	return stacks, apiErr
}

func (cc DefaultStackRepository) GetStackByName(name string) (Stacks, error) {
	var stacks StacksModel
	apiErr := cc.gateway.Get(fmt.Sprintf("/stacks?name=%s", name), &stacks)
	if apiErr != nil {
		return nil, apiErr
	}
	if stacks.Count() < 1 {
		apiErr = fmt.Errorf("Stack not found")
		return nil, apiErr
	}
	return stacks, apiErr
}

func (cc DefaultStackRepository) Update(id string, params map[string]interface{}) (apiErr error) {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	_, apiErr = cc.gateway.Request("PUT", fmt.Sprintf("/stacks/%s", id), data)
	return
}

func (cc DefaultStackRepository) Delete(id string) (apiErr error) {
	apiErr = cc.gateway.Delete(fmt.Sprintf("/stacks/%s", id), nil)
	return
}

func (cc DefaultStackRepository) Publish(id string) (apiErr error) {
	apiErr = cc.gateway.PUT(fmt.Sprintf("/stacks/%s/published", id), nil)
	return
}

func (cc DefaultStackRepository) UnPublish(id string) (apiErr error) {
	apiErr = cc.gateway.PUT(fmt.Sprintf("/stacks/%s/unpublished", id), nil)
	return
}
