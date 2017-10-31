package api

type EventRef interface {
	ID() string
	Type() string
	Links() Links
	Entity() map[string]interface{}
}

type EventRefModel struct {
	IDField     string                 `json:"id"`
	TypeField   string                 `json:"type"`
	LinksField  LinksModel             `json:"links"`
	EntityField map[string]interface{} `json:"content"`
}

func (erm EventRefModel) ID() string {
	return erm.IDField
}

func (erm EventRefModel) Type() string {
	return erm.TypeField
}

func (erm EventRefModel) Links() Links {
	return erm.LinksField
}

func (erm EventRefModel) Entity() map[string]interface{} {
	return erm.EntityField
}

type Events interface {
	Count() int
	Next() (Events, error)
	Prev() (Events, error)
	Items() []EventRef
}

type EventsModel struct {
	CountField      int             `json:"count"`
	NextField       string          `json:"next"`
	PrevField       string          `json:"prev"`
	ItemsField      []EventRefModel `json:"items"`
	EventRepository EventRepository `json:"-"`
}

func (em EventsModel) Count() int {
	return em.CountField
}

func (em EventsModel) Next() (e Events, apiError error) {
	if "" == em.NextField {
		e = nil
		apiError = nil
		return
	}

	events, err := em.EventRepository.GetEventsByURI(em.NextField)
	if err != nil {
		apiError = err
		return
	}

	e = events

	return
}

func (em EventsModel) Prev() (e Events, apiError error) {
	if "" == em.PrevField {
		e = EventsModel{
			CountField: 0,
			ItemsField: []EventRefModel{},
			NextField:  "",
			PrevField:  "",
		}
		apiError = nil
		return
	}

	events, err := em.EventRepository.GetEventsByURI(em.PrevField)
	if err != nil {
		apiError = err
		return
	}

	e = events

	return
}

func (em EventsModel) Items() []EventRef {
	items := make([]EventRef, 0)
	for _, event := range em.ItemsField {
		items = append(items, event)
	}
	return items
}
