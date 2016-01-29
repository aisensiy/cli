package api

type KeyParams struct {
	Public string `json:"public"`
	Name   string `json:"name"`
}

type Key interface {
	ID() string
	Fingerprint() string
	Name() string
	Public() string
	Owner() string
	Links() []Link
}

type KeyModel struct {
	IDField          string `json:"id"`
	PublicField      string `json:"public"`
	FingerprintField string `json:"fingerprint"`
	NameField        string `json:"name"`
	OwnerField       string `json:"owner"`
	LinksField       []Link `json:"links"`
}

func (key KeyModel) ID() string {
	return key.IDField
}

func (key KeyModel) Fingerprint() string {
	return key.FingerprintField
}

func (key KeyModel) Name() string {
	return key.NameField
}

func (key KeyModel) Public() string {
	return key.PublicField
}

func (key KeyModel) Owner() string {
	return key.OwnerField
}

func (key KeyModel) Links() []Link {
	//	return LinksModel{
	//		Links: key.LinksField,
	//	}
	return key.LinksField
}

type Keys interface {
	Count() int
	First() Keys
	Last() Keys
	Prev() Keys
	Next() Keys
	Items() []Key
}

type KeysModel struct {
	CountField int            `json:"count"`
	SelfField  string         `json:"self"`
	FirstField string         `json:"first"`
	LastField  string         `json:"last"`
	PrevField  string         `json:"prev"`
	NextField  string         `json:"next"`
	ItemsField []KeyModel  `json:"items"`

}

func (keys KeysModel) Count() int {
	return keys.CountField
}

func (keys KeysModel) Self() Keys {
	return nil
}

func (keys KeysModel) First() Keys {
	return nil
}

func (keys KeysModel) Last() Keys {
	return nil
}

func (keys KeysModel) Prev() Keys {
	return nil
}

func (keys KeysModel) Next() Keys {
	return nil
}

func (keys KeysModel) Items() []Key {
	items := make([]Key, 0)
	for _, key := range keys.ItemsField {
		items = append(items, key)
	}
	return items
}
