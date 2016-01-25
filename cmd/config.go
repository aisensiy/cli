package cmd

import (
	"fmt"
	"regexp"

	"github.com/cde/client/pkg/prettyprint"
	"github.com/cde/apisdk/api"
	"github.com/cde/client/config"
	"github.com/cde/apisdk/net"
)

// ConfigList lists an app's config.
func ConfigList(appID string, oneLine bool) error {
	configRepository := config.NewConfigRepository(func(error) {})
	appRepository := api.NewAppRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))

	app, err := appRepository.GetApp(appID)

	if err != nil {
		return err
	}

	envs := app.GetEnvs()

	if oneLine {
		for _, key := range envs {
			fmt.Printf("%s=%s ", key, envs[key])
		}
		fmt.Println()
	} else {
		fmt.Printf("=== %s Config\n", appID)

		configMap := make(map[string]string)

		// config.Values is type interface, so it needs to be converted to a string
		for _, key := range envs {
			configMap[key] = fmt.Sprintf("%v", envs[key])
		}

		fmt.Print(prettyprint.PrettyTabs(configMap, 6))
	}

	return nil
}

// ConfigSet sets an app's config variables.
func ConfigSet(appID string, key string, value string) error {
	configRepository := config.NewConfigRepository(func(error) {})
	appRepository := api.NewAppRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))

	app, err := appRepository.GetApp(appID)

	if err != nil {
		return err
	}

	fmt.Print("Creating config... ")

	err = app.SetEnv(key, value)

	if err != nil {
		return err
	}

	return ConfigList(appID, false)
}

// ConfigUnset removes a config variable from an app.
func ConfigUnset(appID string, key string) error {
	configRepository := config.NewConfigRepository(func(error) {})
	appRepository := api.NewAppRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))

	app, err := appRepository.GetApp(appID)

	if err != nil {
		return err
	}

	fmt.Print("Removing config... ")

	err = app.UnsetEnv(key)

	if err != nil {
		return err
	}

	fmt.Print("done\n\n")

	return ConfigList(appID, false)
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