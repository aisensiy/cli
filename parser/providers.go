package parser

import (
	"github.com/docopt/docopt-go"
	"errors"
	"github.com/sjkyspa/stacks/client/cmd"
	"strings"
)

func Providers(argv []string) error {
	usage := `
Valid commands for providers:

providers:enroll 	enroll a new provider
providers:list		list all providers
providers:info		view info about a provider
providers:update	update provider config

Use 'cde help [command]' to learn more.
`
	switch argv[0] {
	case "providers:enroll":
		return providerEnroll(argv)
	case "providers:list":
		return providerList()
	case "providers:info":
		return providerInfo(argv)
	case "providers:update":
		return providerUpdate(argv)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		PrintUsage()
		return nil
	}
	return nil
}

func providerEnroll(argv []string) error {
	usage := `
Enroll a new provider.

Usage: cde providers:enroll <name> <type> [options] (-c <config>)...

Arguments:
  <name>
  	a provider name. No other provider can already exist with this name.
  <type>
  	provider type.
  <config>
  	provider config. Set config as kay value list. Use -c to set a key value pair.

Options:
  --for=<for>
  	specify a organization to use this provider.
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}
	name := safeGetValue(args, "<name>")
	providerType := safeGetValue(args, "<type>")
	config := safeGetValues(args, "<config>")
	consumer := safeGetValue(args, "--for")

	if name == "" || providerType == "" || len(config) <= 0 {
		return errors.New("<name> <type> <config> are essential parameters")
	}

	configMap, err := configConvert(config)

	if err != nil {
		return err
	}

	return cmd.ProviderCreate(name, providerType, consumer, configMap)
}

func configConvert(config []string) (map[string]interface{}, error) {
	configMap := map[string]interface{}{}
	for _, v := range config {
		pair := strings.Split(v, "=")
		if len(pair) != 2 {
			return nil, errors.New("invalid config format")
		}
		configMap[pair[0]] = pair[1]
	}
	return configMap, nil
}

func providerList() error {
	return cmd.ProviderList()
}

func providerInfo(argv []string) error {
	usage := `
View info about a provider.

Usage: cde providers:info <name>

Arguments:
  <name>
  	a provider name.
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}
	name := safeGetValue(args, "<name>")

	if name == "" {
		return errors.New("<name> are essential parameters")
	}

	return cmd.GetProviderByName(name)
}

func providerUpdate(argv []string) error {
	usage := `
Update provider config.

Usage: cde providers:update <name> (-c <config>)...

Arguments:
  <name>
  	a provider name.
  <config>
  	provider config. Set config as kay value list. Use -c to set a key value pair.
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}
	name := safeGetValue(args, "<name>")
	config := safeGetValues(args, "<config>")

	if name == "" || len(config) <= 0 {
		return errors.New("<name> <config> are essential parameters")
	}

	configMap, err := configConvert(config)

	if err != nil {
		return err
	}

	return cmd.ProviderUpdate(name, configMap)
}
