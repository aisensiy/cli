package parser

import (
	"errors"
	"fmt"
	docopt "github.com/docopt/docopt-go"
	"github.com/sjkyspa/cde/cmd"
	cli "gopkg.in/urfave/cli.v2"
)

func ClustersCommands() *cli.Command {
	return &cli.Command{
		Name:  "clusters",
		Usage: "Cluster Commands",
		Subcommands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List clusters.",
				ArgsUsage: " ",
				Action: func(c *cli.Context) error {
					if err := cmd.ClusterList(); err != nil {
						return cli.Exit(fmt.Sprintf("%v", err), 1)
					}
					return nil
				},
			},
			{
				Name:      "info",
				Usage:     "View info about a cluster",
				ArgsUsage: "<cluster-id>",
				Action: func(c *cli.Context) error {
					if !c.Args().Present() {
						return cli.Exit(fmt.Sprintf("USAGE: %s %s", c.Command.HelpName, c.Command.ArgsUsage), 1)
					}
					if err := cmd.GetCluster(c.Args().First()); err != nil {
						return cli.Exit(fmt.Sprintf("%v", err), 1)
					}
					return nil
				},
			},
			{
				Name:      "create",
				Usage:     "Create a new cluster",
				ArgsUsage: "<name> <type> <uri>",
				Action: func(c *cli.Context) error {
					if c.Args().Get(0) == "" || c.Args().Get(1) == "" || c.Args().Get(2) == "" {
						return cli.Exit(fmt.Sprintf("USAGE: %s %s", c.Command.HelpName, c.Command.ArgsUsage), 1)
					}
					if err := cmd.ClusterCreate(c.Args().Get(0), c.Args().Get(1), c.Args().Get(2)); err != nil {
						return cli.Exit(fmt.Sprintf("%v", err), 1)
					}
					return nil
				},
			},
			{
				Name:      "update",
				Usage:     "Update a cluster",
				ArgsUsage: "<cluster-id>",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "name, n",
						Usage: "New name for cluster.",
					},
					&cli.StringFlag{
						Name:  "type, t",
						Usage: "New type for cluster.",
					},
					&cli.StringFlag{
						Name:  "uri, u",
						Usage: "New uri for cluster.",
					},
				},
				Action: func(c *cli.Context) error {
					if !c.Args().Present() {
						return cli.Exit(fmt.Sprintf("USAGE: %s %s", c.Command.HelpName, c.Command.ArgsUsage), 1)
					}
					newName := c.String("name")
					newType := c.String("type")
					newUri := c.String("uri")
					if newName == "" && newType == "" && newUri == "" {
						return cli.Exit("name, type or uri should be given", 1)
					}

					if err := cmd.ClusterUpdate(c.Args().First(), newName, newType, newUri); err != nil {
						return cli.Exit(fmt.Sprintf("%v", err), 1)
					}
					return nil
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete a cluster",
				ArgsUsage: "<cluster-id>",
				Action: func(c *cli.Context) error {
					if !c.Args().Present() {
						return cli.Exit(fmt.Sprintf("USAGE: %s %s", c.Command.HelpName, c.Command.ArgsUsage), 1)
					}
					if err := cmd.ClusterRemove(c.Args().First()); err != nil {
						return cli.Exit(fmt.Sprintf("%v", err), 1)
					}
					return nil
				},
			},
		},
	}
}

// Config routes config commands to their specific function.
func Clusters(argv []string) error {
	usage := `
Valid commands for config:

clusters:list           list clusters
clusters:info           view info about a cluster
clusters:create         create a new cluster
clusters:delete         unset environment variables for a cluster
clusters:update         unset environment variables for a cluster

Use 'cde help [command]' to learn more.
`

	switch argv[0] {
	case "clusters:list":
		return clustersList(argv)
	case "clusters:info":
		return clusterInfo(argv)
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
Creates a new cluster.

Usage: cde clusters:create <name> <type> <uri>

Arguments:
  <name>
  	cluster name
  <type>
  	cluster type
  <uri>
  	cluster uri
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	clusterName := safeGetValue(args, "<name>")
	clusterType := safeGetValue(args, "<type>")
	clusterUri := safeGetValue(args, "<uri>")

	if clusterName == "" || clusterType == "" || clusterUri == "" {
		return errors.New("<name> <type> <uri> are essential parameters")
	}

	return cmd.ClusterCreate(clusterName, clusterType, clusterUri)
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

Usage: cde clusters:delete <cluster-id>

Arguments:
  <cluster-id>
  	a cluster Id
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	clusterId := safeGetValue(args, "<cluster-id>")

	return cmd.ClusterRemove(clusterId)
}

func clusterInfo(argv []string) error {
	usage := `
Prints info about an cluster.

Usage: cde clusters:info <cluster-id>

Arguments:
  <cluster-id>
  	a cluster Id
	`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	clusterId := safeGetValue(args, "<cluster-id>")

	return cmd.GetCluster(clusterId)
}

func clustersUpdate(argv []string) error {
	usage := `
Update cluster info.

Usage: cde clusters:update <cluster-id> [options]

Arguments:
  <cluster-id>
  	a cluster Id

Options:
  -n --name=<name>
    the new name for cluster
  -t --type=<org>
    the new type for cluster
  -u --uri=<uri>
    the new uri for cluster
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	clusterId := safeGetValue(args, "<cluster-id>")
	clusterName := safeGetValue(args, "--name")
	clusterType := safeGetValue(args, "--type")
	clusterUri := safeGetValue(args, "--uri")

	if clusterName == "" && clusterType == "" && clusterUri == "" {
		return errors.New("name, type or uri should given")
	}

	return cmd.ClusterUpdate(clusterId, clusterName, clusterType, clusterUri)
}
