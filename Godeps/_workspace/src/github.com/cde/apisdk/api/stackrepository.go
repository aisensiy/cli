package api

import (
	"encoding/json"
	"fmt"
	"github.com/cde/client/Godeps/_workspace/src/github.com/cde/apisdk/config"
	"github.com/cde/client/Godeps/_workspace/src/github.com/cde/apisdk/net"
)

//go:generate counterfeiter -o fakes/fake_stack_repository.go . StackRepository
type StackRepository interface {
	Create(params StackParams) (createdStack Stack, apiErr error)
	GetStack(id string) (Stack, error)
	GetStackByURI(uri string) (Stack, error)
	GetStacks() (Stacks, error)
	GetStackByName(name string) (Stacks, error)
	Update(id string, params StackParams) (target Stack, apiErr error)
	Delete(id string) (apiErr error)
}

type DefaultStackRepository struct {
	config  config.Reader
	gateway net.Gateway
}

func NewStackRepository(config config.Reader, gateway net.Gateway) StackRepository {
	return DefaultStackRepository{config: config, gateway: gateway}
}

func (cc DefaultStackRepository) Create(params StackParams) (createdStack Stack, apiErr error) {
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
	if apiErr != nil {
		return
	}
	createdStack = stackModel

	return
}

func (cc DefaultStackRepository) GetStack(id string) (Stack, error) {
	var stack StackModel
	apiErr := cc.gateway.Get(fmt.Sprintf("/stacks/%s", id), &stack)
	if apiErr != nil {
		return nil, apiErr
	}
	return stack, apiErr
}

func (cc DefaultStackRepository) GetStackByURI(uri string) (Stack, error) {
	var stack StackModel
	apiErr := cc.gateway.Get(uri, &stack)
	if apiErr != nil {
		return nil, apiErr
	}
	return stack, apiErr
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

func (cc DefaultStackRepository) Update(id string, params StackParams) (target Stack, apiErr error) {
	return nil, nil
}

func (cc DefaultStackRepository) Delete(id string) (apiErr error) {
	apiErr = cc.gateway.Delete(fmt.Sprintf("/stacks/%s", id), nil)
	return
}
