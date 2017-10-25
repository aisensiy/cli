package parser

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/sjkyspa/stacks/client/cmd"
	deployApi "github.com/sjkyspa/stacks/launcher/api/api"
	"os"
	"strconv"
	cli "gopkg.in/urfave/cli.v2"
)

func ServicesCommand() *cli.Command {
	return &cli.Command{
		Name:  "services",
		Usage: "Service Commands",
		Subcommands: []*cli.Command{
			{
				Name:      "create",
				Usage:     "Create Service.",
				ArgsUsage: " ",
				Action: func(c *cli.Context) error {
					if err := cmd.ServiceCreate(); err != nil {
						return cli.Exit(fmt.Sprintf("%v", err), 1)
					}
					return nil
				},
			},
			{
				Name:      "info",
				Usage:     "View service basic information.",
				ArgsUsage: "<service-name>",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "app, a",
						Usage: "Specify app with name.",
					},
				},
				Action: func(c *cli.Context) error {
					if c.Args().Get(0) == "" {
						return cli.Exit(fmt.Sprintf("USAGE: %s %s", c.Command.HelpName, c.Command.ArgsUsage), 1)
					}

					if err := cmd.ServiceInfo(c.String("app"), c.Args().Get(0)); err != nil {
						return cli.Exit(fmt.Sprintf("%v", err), 1)
					}
					return nil
				},
			},
			{
				Name:      "update",
				Usage:     "Update service basic information.",
				ArgsUsage: "<service-name>",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "app, a",
						Usage: "Specify app with name.",
					},
					&cli.StringFlag{
						Name:  "mem",
						Usage: "Specify allocated memory for this service.",
					},
					&cli.StringFlag{
						Name:  "cpu",
						Usage: "Specify max allocated cpu size.",
					},
					&cli.StringFlag{
						Name:  "instances",
						Usage: "Specify instance number",
					},
				},
				Action: func(c *cli.Context) error {
					serviceName := c.Args().Get(0)
					if serviceName == "" {
						return cli.Exit(fmt.Sprintf("USAGE: %s %s", c.Command.HelpName, c.Command.ArgsUsage), 1)
					}

					instances := c.String("instances")
					if instances != "" {
						_, err := strconv.Atoi(instances);
						if err != nil {
							return cli.Exit(fmt.Sprintf("Error: %v\n", err), 1)
						}
					}

					updateParams := make(map[string]string, 0)
					updateParams["mem"] = c.String("mem")
					updateParams["cpu"] = c.String("cpu")
					updateParams["instances"] = instances
					newServiceParams, err := mergeWithOriginService(c.String("app"), serviceName, updateParams);
					if err != nil {
						return cli.Exit(fmt.Sprintf("%v", err), 1)
					}

					if err := cmd.ServiceUpdate(c.String("app"), serviceName, newServiceParams); err != nil {
						return cli.Exit(fmt.Sprintf("%v", err), 1)
					}
					return nil
				},
			},
			{
				Name:      "logs",
				Usage:     "Prints info about the current service.",
				ArgsUsage: "<service-name>",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "app, a",
						Usage: "Specify app with name.",
					},
					&cli.StringFlag{
						Name:  "lines, n",
						Value: "100",
						Usage: "Specify the number of lines to display.",
					},
				},
				Action: func(c *cli.Context) error {
					serviceName := c.Args().Get(0)
					if serviceName == "" {
						return cli.Exit(fmt.Sprintf("USAGE: %s %s", c.Command.HelpName, c.Command.ArgsUsage), 1)
					}

					lines := c.String("lines")
					var lineNum int
					var err error
					if lineNum, err = strconv.Atoi(lines); err != nil {
						return cli.Exit(fmt.Sprintf("Error: %v\n", err), 1)
					}

					if err := cmd.ServiceLog(c.String("app"), serviceName, lineNum); err != nil {
						return cli.Exit(fmt.Sprintf("%v", err), 1)
					}
					return nil
				},
			},
		},
	}
}

func Service(argv []string) error {
	usage := `
Valid commands for services:

services:logs       view serice logs
services:info       view service basic information
services:update     update service basic information

Use 'cde help [command]' to learn more.
`
	switch argv[0] {
	case "create":
		return serviceCreate(argv)
	case "services:info":
		return serviceInfo(argv)
	case "services:update":
		return serviceUpdate(argv)
	case "services:logs":
		return serviceLogs(argv)
	case "services":
		fmt.Print(usage)
	default:
		if printHelp(argv, usage) {
			return nil
		}
		PrintUsage()
		return nil
	}
	return nil
}

func serviceCreate(argv []string) error {
	return cmd.ServiceCreate()
}

func serviceInfo(argv []string) error {
	usage := `
View service basic information.

Usage: cde services:info <service-name> [options]

Arguments:
  <service-name>
    the service name defined in stack file.

Options:
  -a --app=<app>
    the uniquely identifiable id for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	appName := safeGetOrDefault(args, "--app", "")
	serviceName := safeGetOrDefault(args, "<service-name>", "")
	if serviceName == "" {
		return fmt.Errorf("Service name is required!")
	}
	return cmd.ServiceInfo(appName, serviceName)
}

func serviceUpdate(argv []string) error {
	usage := `
Update service basic information.

Usage: cde services:update <service-name> [options]

Arguments:
  <service-name>
    the service name defined in stack file.


Options:
  -a --app=<app>
    the uniquely identifiable id for the application.
  --mem=<mem>
  	allocated memory for this service.
  --cpu=<cpu>
  	max allocated cpu size.
  --instances=<instances>
  	instance number.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	appId := safeGetOrDefault(args, "--app", "")
	serviceName := safeGetOrDefault(args, "<service-name>", "")
	if serviceName == "" {
		return fmt.Errorf("Service name is required!")
	}

	updateParams := make(map[string]string, 0)
	updateParams["mem"] = safeGetOrDefault(args, "--mem", "")
	updateParams["cpu"] = safeGetOrDefault(args, "--cpu", "")
	updateParams["instances"] = safeGetOrDefault(args, "--instances", "")

	newServiceParams, err := mergeWithOriginService(appId, serviceName, updateParams)
	if err != nil {
		return err
	}

	return cmd.ServiceUpdate(appId, serviceName, newServiceParams)
}

func mergeWithOriginService(appName, serviceName string, params map[string]string) (newServiceParams deployApi.ServiceConfigParams, apiErr error) {
	memory := params["mem"]
	cpu := params["cpu"]
	instances := params["instances"]

	originService, apiErr := cmd.GetService(appName, serviceName)
	if apiErr != nil {
		return
	}

	newServiceParams = deployApi.ServiceConfigParams{}
	if mem, err := strconv.ParseFloat(memory, 32); err == nil {
		newServiceParams.Memory = float32(mem)
	} else {
		newServiceParams.Memory = originService.Memory()
	}

	if ins, err := strconv.Atoi(instances); err == nil {
		newServiceParams.Instance = ins
	} else {
		newServiceParams.Instance = originService.Instance()
	}

	if cpu, err := strconv.ParseFloat(cpu, 32); err == nil {
		newServiceParams.CPUS = float32(cpu)
	} else {
		newServiceParams.CPUS = originService.CPU()
	}
	return newServiceParams, nil
}

func serviceLogs(argv []string) error {
	usage := `
Prints info about the current service.

Usage: cde services:logs [options]

Options:
  -a --app=<app>
    the uniquely identifiable id for the application.
  -s --service=<service>
    the service name.
  -n --lines=<lines>
    the number of lines to display
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	appId := safeGetValue(args, "--app")
	service := safeGetOrDefault(args, "--service", "main")

	lines := safeGetOrDefault(args, "--lines", "100")
	var lineNum int
	if lineNum, err = strconv.Atoi(lines); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}

	return cmd.ServiceLog(appId, service, lineNum)
}
