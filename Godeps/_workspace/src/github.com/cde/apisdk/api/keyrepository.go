package api

import (
	"encoding/json"
	"fmt"
	"github.com/cde/client/Godeps/_workspace/src/github.com/cde/apisdk/config"
	"github.com/cde/client/Godeps/_workspace/src/github.com/cde/apisdk/net"
)

//go:generate counterfeiter -o fakes/fake_key_repository.go . KeyRepository
type KeyRepository interface {
	Upload(user User, params KeyParams) (uploaded Key, apiErr error)
	GetKey(id string) (key Key, apiErr error)
	GetKeys() (keys Keys, apiErr error)
	GetKeysForUser(user User) (keys Keys, apiErr error)
	Delete(user User, id string) (apiErr error)
}

type DefaultKeyRepository struct {
	config  config.Reader
	gateway net.Gateway
}

func (repo DefaultKeyRepository) Upload(user User, params KeyParams) (uploadedKey Key, apiErr error) {
	uploadedKey = nil
	data, err := json.Marshal(params)
	if err != nil {
		apiErr = fmt.Errorf("Can not serilize the data")
		return
	}

	res, err := repo.gateway.Request("POST", fmt.Sprintf("/users/%s/keys", user.Id()), data)
	if err != nil {
		apiErr = err
		return
	}

	location := res.Header.Get("Location")
	fmt.Println(location)

	return
}

func (repo DefaultKeyRepository) GetKey(id string) (key Key, apiErr error) {
	return nil, nil
}

func (repo DefaultKeyRepository) GetKeys() (keys Keys, apiErr error) {
	var keyModels KeysModel
	apiErr = repo.gateway.Get(fmt.Sprintf("/keys"), &keyModels)
	keys = keyModels
	return
}

func (repo DefaultKeyRepository) GetKeysForUser(user User) (keys Keys, apiErr error) {
	var keyModels KeysModel
	apiErr = repo.gateway.Get(fmt.Sprintf("/users/%s/keys", user.Id()), &keyModels)
	keys = keyModels
	return
}

func (repo DefaultKeyRepository) Delete(user User, id string) (apiErr error) {
	apiErr = repo.gateway.Delete(fmt.Sprintf("/users/%s/keys/%s", user.Id(), id), nil)
	return
}

func NewKeyRepository(config config.Reader, gateway net.Gateway) KeyRepository {
	return DefaultKeyRepository{config: config, gateway: gateway}
}
