package parser

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/sjkyspa/stacks/client/cmd"
	deployApi "github.com/sjkyspa/stacks/launcher/api/api"
	"os"
	"strconv"
)

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
