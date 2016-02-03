package parser

import (
	"fmt"
	"github.com/sjkyspa/stacks/Godeps/_workspace/src/github.com/docopt/docopt-go"
	"github.com/sjkyspa/stacks/client/cmd"
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
