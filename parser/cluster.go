package parser

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/sjkyspa/stacks/client/cmd"
)

// Config routes config commands to their specific function.
func Clusters(argv []string) error {
	usage := `
Valid commands for config:

clusters:list           list clusters
clusters:create         set environment variables for an app
clusters:delete         unset environment variables for an app
clusters:update         unset environment variables for an app

Use 'cde help [command]' to learn more.
`

	switch argv[0] {
	case "clusters:list":
		return clustersList(argv)
	case "clusters:create":
		return clustersCreate(argv)
	case "clusters:delete":
		return clustersDelete(argv)
	case "clusters:update":
		return clustersUpdate(argv)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "clusters" {
			argv[0] = "clusters:list"
			return clustersList(argv)
		}

		PrintUsage()
		return nil
	}
}



func clustersCreate(argv []string) error {
	usage := `
List the cluster for user.

Usage: cde clusters:create
`

	_, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmd.ClusterCreate()
}

func clustersList(argv []string) error {
	usage := `
List the cluster for user.

Usage: cde clusters:list
`

	_, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmd.ClusterList()
}

func clustersDelete(argv []string) error {
	usage := `
Delete the cluster.

Usage: cde clusters:delete
`

	_, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmd.ClusterList()
}


func clustersUpdate(argv []string) error {
	usage := `
Update cluster info.

Usage: cde clusters:update
`

	_, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmd.ClusterList()
}

