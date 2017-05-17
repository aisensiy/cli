package cmd

import (
	"github.com/sjkyspa/stacks/client/config"
	"github.com/sjkyspa/stacks/client/pkg"
)

func load(appID string) (config.ConfigRepository, string, error) {
	configRepository := config.NewConfigRepository(func(error) {})
	if appID == "" {
		var err error
		appID, err = git.DetectAppName(configRepository.GitHost())

		if err != nil {
			return configRepository, "", err
		}
	}

	return configRepository, appID, nil
}

func loadOrg(orgName string) (config.ConfigRepository, string) {
	configRepository := config.NewConfigRepository(func(error) {})
	if orgName == "" {
		orgName = configRepository.Org()
	}

	return configRepository, orgName
}
