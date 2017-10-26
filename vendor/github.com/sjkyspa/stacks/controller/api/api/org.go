package api

import ()

type OrgParams struct {
	Name string `json:"name"`
}

type Org interface {
	Name() string
	Links() Links
}

type OrgModel struct {
	NAME       string `json:"name"`
	LinksArray []Link `json:"links"`
	OrgMapper  OrgRepository
}

func (o OrgModel) Name() string {
	return o.NAME
}

func (o OrgModel) Links() Links {
	return LinksModel{
		Links: o.LinksArray,
	}
}
