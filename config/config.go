package config

import (
	"sync"
	"github.com/cde/apisdk/config"
)

type Data struct {

}

type DefaultConfigRepository struct {
	data      *Data
	mutex     *sync.RWMutex
	initOnce  *sync.Once
	persistor Persistor
	onError   func(error)
}

func NewConfigRepository(errorHandler func(error)) config.ConfigRepository {
	p, _ := NewPersistor()
	return DefaultConfigRepository{
		persistor: p,
		mutex:     new(sync.RWMutex),
	}
}

func (c *DefaultConfigRepository) read(cb func()) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	cb()
}

func (c *DefaultConfigRepository) write(cb func()) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	cb()

	err := c.persistor.Save()
	if err != nil {
		c.onError(err)
	}
}


func (c DefaultConfigRepository) ApiEndpoint() (endpoint string) {
	c.read(func() {
		endpoint = c.persistor.Endpoint
	})
	return
}

func (c DefaultConfigRepository) SetApiEndpoint(endpoint string) {
	c.write(func() {
		c.persistor.Endpoint = endpoint
	})
}

func (c DefaultConfigRepository) SetEmail(email string) {
	c.write(func() {
		c.persistor.Email = email
	})
}

func (c DefaultConfigRepository) SetAuth(auth string) {
	c.write(func() {
		c.persistor.Auth = auth
	})
}

func (c DefaultConfigRepository) SetId(id string) {
	c.write(func() {
		c.persistor.Id = id
	})
}

func (c DefaultConfigRepository) Email() (email string) {
	c.read(func() {
		email = c.persistor.Email
	})
	return
}

func (c DefaultConfigRepository) Auth() (auth string) {
	c.read(func() {
		auth = c.persistor.Auth
	})
	return
}

func (c DefaultConfigRepository) Id() (id string) {
	c.read(func() {
		id = c.persistor.Id
	})
	return
}

func (c DefaultConfigRepository) Close() {

}
