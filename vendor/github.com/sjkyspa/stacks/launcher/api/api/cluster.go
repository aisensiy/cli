package api

type ClusterParams struct {
	Name	string `json:"name"`
	Type	string `json:"type"`
	Uri	string `json:"uri"`
}

type ClusterRef interface {
	Id() int
	Name() string
	Type() string
	Uri() string
}

type ClusterRefModel struct {
	IDField     int `json:"id"`
	NAMEField     string `json:"name"`
	TYPEField     string `json:"type"`
	URIField      string`json:"uri"`
}

func (crm ClusterRefModel) Id() int {
	return crm.IDField
}

func (crm ClusterRefModel) Name() string {
	return crm.NAMEField
}

func (crm ClusterRefModel) Type() string {
	return crm.TYPEField
}

func (crm ClusterRefModel) Uri() string {
	return crm.URIField
}

type Clusters interface {
	Count() int
	First() Clusters
	Last() Clusters
	Prev() Clusters
	Next() Clusters
	Items() []ClusterRef
}

type ClustersModel struct {
	CountField int           `json:"count"`
	SelfField  string        `json:"self"`
	FirstField string        `json:"first"`
	LastField  string        `json:"last"`
	PrevField  string        `json:"prev"`
	NextField  string        `json:"next"`
	ItemsField []ClusterRefModel `json:"items"`
	ClusterMapper  ClusterRepository
}

func (clusters ClustersModel) Count() int {
	return clusters.CountField
}
func (clusters ClustersModel) Self() Clusters {
	return nil
}
func (clusters ClustersModel) First() Clusters {
	return nil
}
func (clusters ClustersModel) Last() Clusters {
	return nil
}
func (clusters ClustersModel) Prev() Clusters {
	return nil
}
func (clusters ClustersModel) Next() Clusters {
	return nil
}

func (clusters ClustersModel) Items() []ClusterRef {
	items := make([]ClusterRef, 0)
	for _, cluster := range clusters.ItemsField {
		items = append(items, cluster)
	}
	return items
}
