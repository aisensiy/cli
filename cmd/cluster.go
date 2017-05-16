package cmd

import (
	"fmt"
	launcherApi "github.com/sjkyspa/stacks/launcher/api/api"
	deploymentNet "github.com/sjkyspa/stacks/launcher/api/net"
	"github.com/sjkyspa/stacks/client/config"
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
		fmt.Printf("name: %s\n", cluster.Name())
	}
	return nil
}

func ClusterCreate() error {

	return nil
}