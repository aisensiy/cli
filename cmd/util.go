package cmd

import (
	"github.com/cnupp/cli/config"
	"github.com/cnupp/cli/pkg"
)

func load(appID string) (config.ConfigRepository, string, error) {
	configRepository := config.NewConfigRepository(func(error) {})
	if appID == "" {
		var err error
		appID, err = git.DetectAppName(configRepository.GitHost())

		if err != nil {
			return nil, "", err
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
