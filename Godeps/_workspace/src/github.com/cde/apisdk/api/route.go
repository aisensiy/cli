package api

type RouteParams struct {
	Domain string `json:"domain"`
	Path   string `json:"path"`
}

type SimpleDomain struct {
	Name string `json:"name"`
}

type Route interface {
	ID() string
	Path() string
	Domain() SimpleDomain
	Created() string
	Links() Links
}

type RouteModel struct {
	IDField      string     `json:"id"`
	PathField    string        `json:"path"`
	DomainField  SimpleDomain        `json:"domain"`
	CreatedField string        `json:"created"`
	LinksArray   []Link     `json:"links"`
}

func (route RouteModel) ID() string {
	return route.IDField
}

func (route RouteModel) Path() string {
	return route.PathField
}

func (route RouteModel) Domain() SimpleDomain {
	return route.DomainField
}

func (route RouteModel) Created() string {
	return route.CreatedField
}

func (route RouteModel) Links() Links {
	return LinksModel{
		Links: route.LinksArray,
	}
}

type Routes interface {
	Count() int
	First() Routes
	Last() Routes
	Prev() Routes
	Next() Routes
	Items() []Route
}

type RoutesModel struct {
	CountField int            `json:"count"`
	SelfField  string         `json:"self"`
	FirstField string         `json:"first"`
	LastField  string         `json:"last"`
	PrevField  string         `json:"prev"`
	NextField  string         `json:"next"`
	ItemsField []RouteModel  `json:"items"`
}

func (routes RoutesModel) Count() int {
	return routes.CountField
}
func (routes RoutesModel) Self() Routes {
	return nil
}
func (routes RoutesModel) First() Routes {
	return nil
}
func (routes RoutesModel) Last() Routes {
	return nil
}
func (routes RoutesModel) Prev() Routes {
	return nil
}
func (routes RoutesModel) Next() Routes {
	return nil
}

func (routes RoutesModel) Items() []Route {
	items := make([]Route, 0)
	for _, route := range routes.ItemsField {
		items = append(items, route)
	}
	return items
}
