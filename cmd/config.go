package cmd

import (
	"fmt"
	"regexp"

	"github.com/sjkyspa/cde-client/pkg/prettyprint"
	"github.com/sjkyspa/stacks/controller/api/api"
	"github.com/sjkyspa/stacks/controller/api/net"
)

// ConfigList lists an app's config.
func ConfigList(appId string, oneLine bool) error {
	configRepository, appId, err := load(appId)
	if err != nil {
		return err
	}
	appRepository := api.NewAppRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))

	app, err := appRepository.GetApp(appId)

	if err != nil {
		return err
	}

	envs := app.GetEnvs()

	if oneLine {
		for key, value := range envs {
			fmt.Printf("%s=%s ", key, value)
		}
		fmt.Println()
	} else {
		fmt.Printf("=== %s Config\n", appId)

		configMap := make(map[string]string)

		// config.Values is type interface, so it needs to be converted to a string
		for key, value := range envs {
			configMap[key] = fmt.Sprintf("%v", value)
		}

		fmt.Print(prettyprint.PrettyTabs(configMap, 6))
	}

	return nil
}

// ConfigSet sets an app's config variables.
func ConfigSet(appId string, configVars []string) error {
	configRepository, appId, err := load(appId)
	if err != nil {
		return err
	}
	appRepository := api.NewAppRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))

	app, err := appRepository.GetApp(appId)

	if err != nil {
		return err
	}

	fmt.Print("Creating config... ")

	configMap := parseConfig(configVars)
	err = app.SetEnv(configMap)

	if err != nil {
		return err
	}

	return ConfigList(appId, false)
}

// ConfigUnset removes a config variable from an app.
func ConfigUnset(appId string, keys []string) error {
	configRepository, appId, err := load(appId)
	if err != nil {
		return err
	}
	appRepository := api.NewAppRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))

	app, err := appRepository.GetApp(appId)

	if err != nil {
		return err
	}

	fmt.Print("Removing config... ")

	err = app.UnsetEnv(keys)

	if err != nil {
		return err
	}

	fmt.Print("done\n\n")

	return ConfigList(appId, false)
}

func parseConfig(configVars []string) map[string]interface{} {
	configMap := make(map[string]interface{})

	regex := regexp.MustCompile(`^([A-z_]+[A-z0-9_]*)=([\s\S]+)$`)
	for _, config := range configVars {
		if regex.MatchString(config) {
			captures := regex.FindStringSubmatch(config)
			configMap[captures[1]] = captures[2]
		} else {
			fmt.Printf("'%s' does not match the pattern 'key=var', ex: MODE=test\n", config)
		}
	}

	return configMap
}

func formatConfig(configVars map[string]interface{}) string {
	var formattedConfig string

	for key, value := range configVars {
		formattedConfig += fmt.Sprintf("%s=%s\n", key, value)
	}

	return formattedConfig
}
