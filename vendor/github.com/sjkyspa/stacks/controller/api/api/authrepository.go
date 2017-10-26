package api

import (
	"encoding/json"
	"fmt"
	"github.com/sjkyspa/stacks/controller/api/config"
	"github.com/sjkyspa/stacks/controller/api/net"
	"strings"
)

//go:generate counterfeiter -o fakes/fake_auth_repository.go . AuthRepository
type AuthRepository interface {
	Create(params UserParams) (auth Auth, apiErr error)
	Get() (user User, apiErr error)
	Delete(id string) error
}

type DefaultAuthRepository struct {
	config  config.Reader
	gateway net.Gateway
}

func NewAuthRepository(config config.Reader, gateway net.Gateway) AuthRepository {
	return DefaultAuthRepository{config: config, gateway: gateway}
}

func (authRepository DefaultAuthRepository) Create(params UserParams) (createdAuth Auth, apiErr error) {
	data, err := json.Marshal(params)
	if err != nil {
		apiErr = fmt.Errorf("%s", "Can not serialize the data")
		return
	}

	res, err := authRepository.gateway.Request("POST", "/auths", data)
	if err != nil {
		apiErr = err
		return
	}

	location := res.Header.Get("Location")
	splits := strings.Split(location, "/")

	return AuthModel{
		UserEmailField: params.Email,
		IdField:        splits[len(splits) - 1],
	}, nil
}

func (authRepository DefaultAuthRepository) Get() (user User, apiErr error) {
	var remoteUser UserModel
	apiErr = authRepository.gateway.Get("/auths", &remoteUser)
	if apiErr != nil {
		return
	}
	user = remoteUser
	return
}

func (cc DefaultAuthRepository) Delete(id string) (apiErr error) {
	apiErr = cc.gateway.Delete(fmt.Sprintf("/auths/%s", id), nil)
	return
}
