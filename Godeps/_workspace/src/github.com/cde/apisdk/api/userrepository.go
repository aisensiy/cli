package api

import (
	"encoding/json"
	"fmt"
	"github.com/cde/client/Godeps/_workspace/src/github.com/cde/apisdk/config"
	"github.com/cde/client/Godeps/_workspace/src/github.com/cde/apisdk/net"
)

//go:generate counterfeiter -o fakes/fake_user_repository.go . UserRepository
type UserRepository interface {
	Create(params UserParams) (createdUser User, apiErr error)
	GetUser(id string) (User, error)
	GetUserByEmail(email string) (Users, error)
}

type DefaultUserRepository struct {
	config  config.Reader
	gateway net.Gateway
}

func NewUserRepository(config config.Reader, gateway net.Gateway) UserRepository {
	return DefaultUserRepository{config: config, gateway: gateway}
}

func (cc DefaultUserRepository) Create(params UserParams) (createdUser User, apiErr error) {
	data, err := json.Marshal(params)
	if err != nil {
		apiErr = fmt.Errorf("Can not serilize the data")
		return
	}

	res, err := cc.gateway.Request("POST", "/users", data)
	if err != nil {
		apiErr = err
		return
	}

	location := res.Header.Get("Location")

	var userModel UserModel
	apiErr = cc.gateway.Get(location, &userModel)
	if apiErr != nil {
		return
	}

	userModel.KeyRepo = NewKeyRepository(cc.config, cc.gateway)
	createdUser = userModel
	return
}

func (cc DefaultUserRepository) GetUser(id string) (user User, apiErr error) {
	var remoteUser UserModel
	apiErr = cc.gateway.Get(fmt.Sprintf("/users/%s", id), &remoteUser)
	if apiErr != nil {
		return
	}

	remoteUser.KeyRepo = NewKeyRepository(cc.config, cc.gateway)
	user = remoteUser
	return
}

func (cc DefaultUserRepository) GetUserByEmail(email string) (users Users, apiErr error) {
	var usersModel UsersModel
	apiErr = cc.gateway.Get(fmt.Sprintf("/users?email=%s", email), &usersModel)
	if apiErr != nil {
		return nil, apiErr
	}
	if usersModel.Count() < 1 {
		apiErr = fmt.Errorf("User not found")
		return nil, apiErr
	}
	users = usersModel
	return
}
