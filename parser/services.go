package parser

import (
	"fmt"
	"github.com/sjkyspa/stacks/Godeps/_workspace/src/github.com/docopt/docopt-go"
	"github.com/sjkyspa/stacks/client/cmd"
	"strconv"
	deployApi "github.com/sjkyspa/stacks/deploymentsdk/api"
)

func Service(argv []string) error {
	usage := `
Valid commands for services:

services:log		view serice logs
services:info    	view service basic information
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

Usage: cde services:info <app-name> <service-name>

Arguments:
  <app-name>
    the (hosted) application name.
  <service-name>
    the service name defined in stack file.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	appName := safeGetOrDefault(args, "<app-name>", "")
	serviceName := safeGetOrDefault(args, "<service-name>", "")
	if(appName == "" || serviceName == ""){
		return fmt.Errorf("Application name and service name are both required!")
	}
	return cmd.ServiceInfo(appName, serviceName)
}

func serviceUpdate(argv []string) error {
	usage := `
Update service basic information.

Usage: cde services:update <app-name> <service-name> [options]

Arguments:
  <app-name>
    the (hosted) application name.
  <service-name>
    the service name defined in stack file.


Options:
  --mem=<mem>
  	allocated memory for this service.
  --cpu=<cpu>
  	max allocated disk size.
  --instance=<instance>
  	instance number.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	appName := safeGetOrDefault(args, "<app-name>", "")
	serviceName := safeGetOrDefault(args, "<service-name>", "")
	if(appName == "" || serviceName == ""){
		return fmt.Errorf("Application name and service name are both required!")
	}

	updateParams := make(map[string]string, 0)
	updateParams["mem"] = safeGetOrDefault(args, "--mem", "")
	updateParams["cpu"] = safeGetOrDefault(args, "--cpu", "")
	updateParams["instance"] = safeGetOrDefault(args, "--instance", "")

	newServiceParams, err := mergeWithOriginService(appName, serviceName, updateParams)
	if err != nil {
		return err
	}

	return cmd.ServiceUpdate(appName, serviceName, newServiceParams)
}

func mergeWithOriginService(appName, serviceName string, params map[string]string) (deployApi.ServiceConfigParams, apiErr error){
	memory := params["mem"]
	cpu := params["cpu"]
	instances := params["instance"]

	originService, apiErr := cmd.GetService(appName, serviceName)
	if(apiErr != nil){
		return
	}

	newServiceParams := deployApi.ServiceConfigParams{}
	if mem, err := strconv.Atoi(memory); err == nil {
		newServiceParams.Memory = mem
	}else{
		newServiceParams.Memory = originService.Memory()
	}

	if ins, err := strconv.Atoi(instances); err == nil {
		newServiceParams.Instance = ins
	}else{
		newServiceParams.Instance = originService.Instance()
	}

	if cpu, err := strconv.ParseFloat(cpu, 32); err == nil {
		newServiceParams.CPUS = float32(cpu)
	}else{
		newServiceParams.CPUS = originService.CPU()
	}
	return newServiceParams, nil
}
