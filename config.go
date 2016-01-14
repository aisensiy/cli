package main

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
	persistor config.Persistor
	onError   func(error)
}

func NewRepositoryFromFilepath(filepath string, errorHandler func(error)) config.ConfigRepository {
	if errorHandler == nil {
		return nil
	}
	return NewConfigRepositoryFromPersistor(config.NewPersistor(filepath), errorHandler)
}

func NewConfigRepositoryFromPersistor(persistor config.Persistor, errorHandler func(error)) config.ConfigRepository {
	return DefaultConfigRepository{}
}

type Reader interface {
	ApiEndpoint() string
}

type Writer interface {
	SetApiEndpoint(string)
}

func(c DefaultConfigRepository) ApiEndpoint() string {
	return "http://www.tw.com"
}

func(c DefaultConfigRepository) SetApiEndpoint(endpoint string) {

}

func(c DefaultConfigRepository) Close() {

}
