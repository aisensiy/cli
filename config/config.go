package config

import (
	"sync"
	"path/filepath"
	"os"
)

type ConfigRepository interface {
	Reader
	Writer
	Close()
}

type Reader interface {
	ApiEndpoint() string
	DeploymentEndpoint() string
	Email() string
	Auth() string
	Id() string
}

type Writer interface {
	SetApiEndpoint(string)
	SetDeploymentEndpoint(string)
	SetEmail(string)
	SetAuth(string)
	SetId(string)
}

type DefaultConfigRepository struct {
	data      *Data
	mutex     *sync.RWMutex
	initOnce  *sync.Once
	persistor Persistor
	onError   func(error)
}

func NewConfigRepository(errorHandler func(error)) ConfigRepository {
	if errorHandler == nil {
		return nil
	}
	path := DefaultFilePath()
	return NewRepositoryFromPersistor(NewDiskPersistor(path), errorHandler)
}

func NewRepositoryFromFilepath(filepath string, errorHandler func(error)) ConfigRepository {
	if errorHandler == nil {
		return nil
	}
	return NewRepositoryFromPersistor(NewDiskPersistor(filepath), errorHandler)
}

func NewRepositoryFromPersistor(persistor Persistor, errorHandler func(error)) ConfigRepository {
	data := NewData()

	return &DefaultConfigRepository{
		data:      data,
		mutex:     new(sync.RWMutex),
		initOnce:  new(sync.Once),
		persistor: persistor,
		onError:   errorHandler,
	}
}

func (c *DefaultConfigRepository) init() {
	c.initOnce.Do(func() {
		err := c.persistor.Load(c.data)
		if err != nil {
			c.onError(err)
		}
	})
}


func (c *DefaultConfigRepository) read(cb func()) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	c.init()

	cb()
}

func (c *DefaultConfigRepository) write(cb func()) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	cb()

	err := c.persistor.Save(c.data)
	if err != nil {
		c.onError(err)
	}
}


func (c DefaultConfigRepository) ApiEndpoint() (endpoint string) {
	c.read(func() {
		endpoint = c.data.Endpoint
	})
	return
}

func (c DefaultConfigRepository) SetApiEndpoint(endpoint string) {
	c.write(func() {
		c.data.Endpoint = endpoint
	})
}

func (c DefaultConfigRepository) DeploymentEndpoint() (endpoint string) {
	c.read(func() {
		endpoint = c.data.DeploymentEndpoint
	})
	return
}

func (c DefaultConfigRepository) SetDeploymentEndpoint(endpoint string) {
	c.write(func() {
		c.data.DeploymentEndpoint = endpoint
	})
}

func (c DefaultConfigRepository) SetEmail(email string) {
	c.write(func() {
		c.data.Email = email
	})
}

func (c DefaultConfigRepository) SetAuth(auth string) {
	c.write(func() {
		c.data.Auth = auth
	})
}

func (c DefaultConfigRepository) SetId(id string) {
	c.write(func() {
		c.data.Id = id
	})
}

func (c DefaultConfigRepository) Email() (email string) {
	c.read(func() {
		email = c.data.Email
	})
	return
}

func (c DefaultConfigRepository) Auth() (auth string) {
	c.read(func() {
		auth = c.data.Auth
	})
	return
}

func (c DefaultConfigRepository) Id() (id string) {
	c.read(func() {
		id = c.data.Id
	})
	return
}

func (c DefaultConfigRepository) Close() {
	c.read(func() {
		// perform a read to ensure write lock has been cleared
	})
}

func DefaultFilePath() string {
	var configDir string

	if os.Getenv("CDE_HOME") != "" {
		cfHome := os.Getenv("CDE_HOME")
		configDir = filepath.Join(cfHome, ".cde")
	} else {
		configDir = filepath.Join(userHomeDir(), ".cde")
	}

	return filepath.Join(configDir, "config.json")
}

var userHomeDir = func() string {
	return os.Getenv("HOME")
}
