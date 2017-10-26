package api

import (
	"encoding/json"
	"fmt"
	"github.com/sjkyspa/stacks/controller/api/config"
	"github.com/sjkyspa/stacks/controller/api/net"
)

//go:generate counterfeiter -o fakes/fake_app_repository.go . AppRepository
type OrgRepository interface {
	Create(params OrgParams) (org Org, apiErr error)
	GetOrg(name string) (Org, error)
	GetOrgMembers(name string) (users []UserModel, apiErr error)
	AddMember(orgName string, userEmail string) (apiErr error)
	RmMember(orgName string, userId string) (apiErr error)
	GetApps(orgName string) (apps []AppModel, apiErr error)
	AddApp(orgName string, appName string) (apiErr error)
	Delete(orgName string) (apiErr error)
}

type CloudControllerOrgRepository struct {
	config  config.Reader
	gateway net.Gateway
}

type AddMemberParams struct {
	Email string `json:"email"`
}

type AddAppParams struct {
	Name string `json:"name"`
}

func NewOrgRepository(config config.Reader, gateway net.Gateway) OrgRepository {
	return CloudControllerOrgRepository{config: config, gateway: gateway}
}

func (cc CloudControllerOrgRepository) Create(params OrgParams) (org Org, apiErr error) {
	data, err := json.Marshal(params)
	if err != nil {
		apiErr = fmt.Errorf("Can not serilize the data")
		return
	}

	res, err := cc.gateway.Request("POST", "/orgs", data)
	if err != nil {
		apiErr = err
		return
	}

	location := res.Header.Get("Location")

	var orgModel OrgModel
	apiErr = cc.gateway.Get(location, &orgModel)
	if apiErr != nil {
		return
	}
	orgModel.OrgMapper = NewOrgRepository(cc.config, cc.gateway)
	org = orgModel
	return
}

func (cc CloudControllerOrgRepository) GetOrg(orgName string) (org Org, apiErr error) {
	var orgModel OrgModel
	apiErr = cc.gateway.Get(fmt.Sprintf("/orgs/%s", orgName), &orgModel)
	if apiErr != nil {
		return
	}
	orgModel.OrgMapper = NewOrgRepository(cc.config, cc.gateway)
	org = orgModel
	return
}

func (cc CloudControllerOrgRepository) GetOrgMembers(orgName string) (users []UserModel, apiErr error) {
	apiErr = cc.gateway.Get(fmt.Sprintf("/orgs/%s/members", orgName), &users)
	return
}

func (cc CloudControllerOrgRepository) AddMember(orgName string, userEmail string) (apiErr error) {
	params := AddMemberParams{
		Email: userEmail,
	}
	data, err := json.Marshal(params)
	if err != nil {
		apiErr = fmt.Errorf("Can not serilize the data")
		return
	}
	_, apiErr = cc.gateway.Request("POST", fmt.Sprintf("/orgs/%s/members", orgName), data)
	return
}

func (cc CloudControllerOrgRepository) RmMember(orgName string, userId string) (apiErr error) {
	apiErr = cc.gateway.Delete(fmt.Sprintf("/orgs/%s/members/%s", orgName, userId), nil)
	return
}

func (cc CloudControllerOrgRepository) GetApps(orgName string) (apps []AppModel, apiErr error) {
	apiErr = cc.gateway.Get(fmt.Sprintf("/orgs/%s/apps", orgName), &apps)
	return
}

func (cc CloudControllerOrgRepository) AddApp(orgName string, appName string) (apiErr error) {
	params := AddAppParams{
		Name: appName,
	}
	data, err := json.Marshal(params)
	if err != nil {
		apiErr = fmt.Errorf("Can not serilize the data")
		return
	}
	_, apiErr = cc.gateway.Request("POST", fmt.Sprintf("/orgs/%s/apps", orgName), data)
	return
}

func (cc CloudControllerOrgRepository) Delete(orgName string) (apiErr error) {
	return cc.gateway.Delete("/orgs/"+orgName, nil)
}
