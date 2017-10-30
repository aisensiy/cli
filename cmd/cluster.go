package cmd

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/sjkyspa/cde-client/config"
	launcherApi "github.com/sjkyspa/stacks/launcher/api/api"
	deploymentNet "github.com/sjkyspa/stacks/launcher/api/net"
	"os"
)

func ClusterList() error {
	configRepository := config.NewConfigRepository(func(error) {})
	clusterRepository := launcherApi.NewClusterRepository(configRepository, deploymentNet.NewCloudControllerGateway(configRepository))
	clusters, err := clusterRepository.GetClusters()
	if err != nil {
		return err
	}

	fmt.Printf("=== Clusters [%d]\n", len(clusters.Items()))

	for _, cluster := range clusters.Items() {
		fmt.Printf("name: %s\t id: %d \n", cluster.Name(), cluster.Id())
	}
	return nil
}

func GetCluster(clusterId string) error {
	configRepository := config.NewConfigRepository(func(error) {})
	clusterRepository := launcherApi.NewClusterRepository(configRepository, deploymentNet.NewCloudControllerGateway(configRepository))
	cluster, err := clusterRepository.GetClusterById(clusterId)
	if err != nil {
		return err
	}
	outputClusterDescription(cluster)

	return nil
}

func outputClusterDescription(cluster launcherApi.ClusterRef) {
	fmt.Printf("--- %s Cluster\n", cluster.Name())
	data := make([][]string, 3)
	data[0] = []string{"NAME", cluster.Name()}
	data[1] = []string{"TYPE", cluster.Type()}
	data[2] = []string{"URI", cluster.Uri()}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, v := range data {
		table.Append(v)
	}

	table.Render()
}

func ClusterCreate(clusterName string, clusterType string, clusterUri string) error {
	configRepository := config.NewConfigRepository(func(error) {})
	clusterRepository := launcherApi.NewClusterRepository(configRepository, deploymentNet.NewCloudControllerGateway(configRepository))
	clusterParams := launcherApi.ClusterParams{
		Name: clusterName,
		Type: clusterType,
		Uri:  clusterUri,
	}

	createdCluster, err := clusterRepository.CreateCluster(clusterParams)
	if err != nil {
		return err
	}
	fmt.Printf("create cluster %s with id %d\n", createdCluster.Name(), createdCluster.Id())
	return nil
}

func ClusterRemove(clusterId string) error {
	configRepository := config.NewConfigRepository(func(error) {})
	clusterRepository := launcherApi.NewClusterRepository(configRepository, deploymentNet.NewCloudControllerGateway(configRepository))

	err := clusterRepository.DeleteClusterById(clusterId)
	if err != nil {
		return err
	}
	fmt.Printf("delete cluster successfully\n")
	return nil
}

func ClusterUpdate(clusterId string, clusterName string, clusterType string, clusterUri string) error {
	configRepository := config.NewConfigRepository(func(error) {})
	clusterRepository := launcherApi.NewClusterRepository(configRepository, deploymentNet.NewCloudControllerGateway(configRepository))

	cluster, err := clusterRepository.GetClusterById(clusterId)
	if err != nil {
		return err
	}

	if clusterName == "" {
		clusterName = cluster.Name()
	}
	if clusterType == "" {
		clusterType = cluster.Type()
	}
	if clusterUri == "" {
		clusterUri = cluster.Uri()
	}
	clusterParams := launcherApi.ClusterParams{
		Name: clusterName,
		Type: clusterType,
		Uri:  clusterUri,
	}

	updateErr := clusterRepository.UpdateCluster(clusterId, clusterParams)
	if updateErr != nil {
		return updateErr
	}
	fmt.Printf("update cluster successfully\n")
	return nil
}
