package api

import (
	"fmt"
	"github.com/cnupp/appssdk/config"
	"github.com/cnupp/appssdk/net"
)

//go:generate counterfeiter -o fakes/fake_event_repository.go . EventRepository
type EventRepository interface {
	GetEvents(eventType string) (Events, error)
	GetEventsByURI(uri string) (Events, error)
}

type DefaultEventRepository struct {
	config  config.Reader
	gateway net.Gateway
}

func NewEventRepository(config config.Reader, gateway net.Gateway) EventRepository {
	return DefaultEventRepository{
		config:  config,
		gateway: gateway,
	}
}

func (der DefaultEventRepository) GetEvents(eventType string) (events Events, apiError error) {
	var eventsModel EventsModel
	err := der.gateway.Get(fmt.Sprintf("/events?type=%s", eventType), &eventsModel)
	if err != nil {
		apiError = err
		return
	}
	eventsModel.EventRepository = der
	events = eventsModel

	return
}

func (der DefaultEventRepository) GetEventsByURI(uri string) (events Events, apiError error) {
	var eventsModel EventsModel
	err := der.gateway.Get(uri, &eventsModel)
	if err != nil {
		apiError = err
		return
	}

	eventsModel.EventRepository = der
	events = eventsModel
	return
}
