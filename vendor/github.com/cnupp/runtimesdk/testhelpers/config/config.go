package config

import (
	"github.com/cnupp/runtimesdk/config"
)

var conf_endpoint string
var conf_email string
var conf_auth string
var conf_id string

type fakeConfigRepository struct {
}

func (c fakeConfigRepository) DeploymentEndpoint() string {
	return conf_endpoint
}

func (c fakeConfigRepository) SetDeploymentEndpoint(endpoint string) {
	conf_endpoint = endpoint
}

func (c fakeConfigRepository) Email() string {
	return conf_email
}

func (c fakeConfigRepository) Id() string {
	return conf_id
}

func (c fakeConfigRepository) SetEmail(email string) {
	conf_email = email
}

func (c fakeConfigRepository) Auth() string {
	return conf_auth
}

func (c fakeConfigRepository) SetAuth(auth string) {
	conf_auth = auth
}

func (c fakeConfigRepository) SetId(id string) {
	conf_id = id
}

func (c fakeConfigRepository) Close() {

}

func NewConfigRepository() config.ConfigRepository {
	return fakeConfigRepository{}
}

func NewRepositoryWithDefaults() config.ConfigRepository {
	configRepo := NewConfigRepository()

	return configRepo
}
