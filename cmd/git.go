package cmd

import (
	"github.com/sjkyspa/stacks/client/pkg"
)

func GitRemote(appID, remote string) error {
	host := "192.168.50.6"
	appID, err := load(appID, host)

	if err != nil {
		return err
	}

	return git.CreateRemote(host, remote, appID)
}

func load(appID string, host string) (string, error) {
	appID, err := git.DetectAppName(host)
	if err != nil {
		return "", err
	}
	return appID, nil
}
