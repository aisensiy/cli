package api

type DomainParams struct {
	Name string `json:"name"`
}

type Domain interface {
	Id() string
	Name() string
}

type DomainModel struct {
	IdField   string     `json:"id"`
	NameField string     `json:"name"`
}

func (d DomainModel) Id() string {
	return d.IdField
}

func (d DomainModel) Name() string {
	return d.NameField
}


type DomainRef interface {
	Id() string
	Name() string
	Links() Links
}

type DomainRefModel struct {
	IdField    string `json:"id"`
	NameField  string `json:"name"`
	LinksField []Link `json:"links"`
}

func (drm DomainRefModel) Id() string {
	return drm.IdField
}

func (drm DomainRefModel) Links() Links {
	return LinksModel{
		Links: drm.LinksField,
	}
}

func (drm DomainRefModel) Name() string {
	return drm.NameField
}

type Domains interface {
	Count() int
	First() Domains
	Last() Domains
	Prev() Domains
	Next() Domains
	Items() []DomainRef
}

type DomainsModel struct {
	CountField   int            `json:"count"`
	SelfField    string         `json:"self"`
	FirstField   string         `json:"first"`
	LastField    string         `json:"last"`
	PrevField    string         `json:"prev"`
	NextField    string         `json:"next"`
	ItemsField   []DomainRefModel  `json:"items"`
	DomainMapper DomainRepository
}

func (domains DomainsModel) Count() int {
	return domains.CountField
}
func (domains DomainsModel) Self() Domains {
	return nil
}
func (domains DomainsModel) First() Domains {
	return nil
}
func (domains DomainsModel) Last() Domains {
	return nil
}
func (domains DomainsModel) Prev() Domains {
	return nil
}
func (domains DomainsModel) Next() Domains {
	return nil
}

func (domains DomainsModel) Items() []DomainRef {
	items := make([]DomainRef, 0)
	for _, app := range domains.ItemsField {
		items = append(items, app)
	}
	return items
}
