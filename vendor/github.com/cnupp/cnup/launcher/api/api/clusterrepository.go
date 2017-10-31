package api

import (
	"github.com/cnupp/cnup/launcher/api/config"
	"github.com/cnupp/cnup/launcher/api/net"
	"fmt"
	"encoding/json"
)

type ClusterRepository interface {
	GetClusters() (Clusters, error)
	GetClusterById(id string) (ClusterRef, error)
	CreateCluster(params ClusterParams)(ClusterRef, error)
	DeleteClusterById(id string) (error)
	UpdateCluster(id string, params ClusterParams) (error)
}

type DefaultClusterRepository struct {
	config  config.Reader
	gateway net.Gateway
}

func (dcr DefaultClusterRepository) GetClusters() (clusters Clusters, clusterError error)  {
	var remoteClusters ClustersModel
	clusterError = dcr.gateway.Get(fmt.Sprintf("/clusters"), &remoteClusters)
	if clusterError != nil {
		return
	}
	remoteClusters.ClusterMapper = dcr
	clusters = remoteClusters
	return
}

func (dcr DefaultClusterRepository) GetClusterById(id string) (ClusterRef, error){
	var cluster ClusterRefModel
	clusterError := dcr.gateway.Get(fmt.Sprintf("/clusters/%s", id), &cluster)
	if clusterError != nil {
		return nil, clusterError
	}
	return cluster, clusterError
}

func (dcr DefaultClusterRepository)CreateCluster(params ClusterParams)(ClusterRef, error){
	data, err := json.Marshal(params)
	if err != nil {
		clusterErr := fmt.Errorf("Can not serilize the data")
		return nil, clusterErr
	}

	res,err := dcr.gateway.Request("POST", "/clusters", data)
	if err != nil {
		return nil, err
	}

	location := res.Header.Get("Location")

	var clusterModel ClusterRefModel
	clusterError := dcr.gateway.Get(location, &clusterModel)
	if clusterError != nil {
		return nil, clusterError
	}
	return clusterModel, clusterError
}

func (dcr DefaultClusterRepository)DeleteClusterById(id string)(error){
	err := dcr.gateway.Delete(fmt.Sprintf("/clusters/%s", id), nil)
	return err
}

func (dcr DefaultClusterRepository)UpdateCluster(id string, params ClusterParams) (error){
	err := dcr.gateway.PUT(fmt.Sprintf("/clusters/%s", id), params)
	return err
}

func NewClusterRepository (config config.Reader, gateway net.Gateway) ClusterRepository {
	return DefaultClusterRepository{config: config, gateway: gateway}
}
