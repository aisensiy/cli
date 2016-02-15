package cmd
import (
	"github.com/sjkyspa/stacks/client/pkg"
	"github.com/sjkyspa/stacks/client/config"
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
