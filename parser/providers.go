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

Use 'cde help [command]' to learn more.
`
	switch argv[0] {
	case "providers:enroll":
		return providerEnroll(argv)
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

Usage: cde providers:enroll <name> <type> (-c <config>)...

Arguments:
  <name>
  	a provider name. No other provider can already exist with this name.
  <type>
  	provider type.
  <config>
  	provider config. Set config as kay value list. Use -c to set a key value pair.
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}
	name := safeGetValue(args, "<name>")
	providerType := safeGetValue(args, "<type>")
	config := safeGetValues(args, "<config>")

	if name == "" || providerType == "" || len(config) <= 0 {
		return errors.New("<name> <type> <config> are essential parameters")
	}

	configMap, err := configConvert(config)

	if err != nil {
		return err
	}

	return cmd.ProviderCreate(name, strings.ToUpper(providerType), configMap)
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
