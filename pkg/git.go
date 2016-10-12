package git

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

// CreateRemote adds a git remote in the current directory.
func CreateRemote(host, remote, appID string) error {
	cmd := exec.Command("git", "remote", "add", remote, RemoteURL(host, appID))
	stderr, err := cmd.StderrPipe()

	if err != nil {
		return err
	}

	if err = cmd.Start(); err != nil {
		return err
	}

	output, _ := ioutil.ReadAll(stderr)
	fmt.Print(string(output))

	if err := cmd.Wait(); err != nil {
		return err
	}

	fmt.Printf("Git remote %s added\n", remote)

	return nil
}

// DeleteRemote removes a git remote in the current directory.
func DeleteRemote(appID string) error {
	name, err := remoteNameFromAppID(appID)

	if err != nil {
		return err
	}

	if _, err = exec.Command("git", "remote", "remove", name).Output(); err != nil {
		return err
	}

	fmt.Printf("Git remote %s removed\n", name)

	return nil
}

func HasRemoteNameForApp(remoteName, appId string) (bool){
	name, err := remoteNameFromAppID(appId)
	if(err != nil) {
		return false
	}
	return name == remoteName
}

func remoteNameFromAppID(appID string) (string, error) {
	out, err := exec.Command("git", "remote", "-v").Output()

	if err != nil {
		return "", err
	}

	cmd := string(out)

	for _, line := range strings.Split(cmd, "\n") {
		if strings.Contains(line, appID) {
			return strings.Split(line, "\t")[0], nil
		}
	}

	return "", errors.New("Could not find remote matching app in 'git remote -v'")
}

// DetectAppName detects if there is cde remote in git.
func DetectAppName(host string) (string, error) {
	remote, err := findRemote(host)

	if err != nil {
		return "", errors.New("Cannot detect the app name.\n" +
			"You may not be in a project OR no application has been created for this project")
	}

	ss := strings.Split(remote, "/")
	return strings.Split(ss[len(ss)-1], ".")[0], nil
}

func findRemote(host string) (string, error) {
	out, err := exec.Command("git", "remote", "-v").Output()

	if err != nil {
		return "", err
	}

	cmd := string(out)

	for _, line := range strings.Split(cmd, "\n") {
		for _, remote := range strings.Split(line, " ") {
			if strings.Contains(remote, host) {
				return strings.Split(remote, "\t")[1], nil
			}
		}
	}

	return "", errors.New("Could not find cde remote in 'git remote -v'")
}

// RemoteURL returns the git URL of app.
func RemoteURL(host, appID string) string {
	return fmt.Sprintf("ssh://git@%s:2222/%s.git", host, appID)
}

func DeleteCdeRemote() error {
	if _, err := exec.Command("git", "remote", "remove", "cde").Output(); err != nil {
		return err
	}
	return nil
}

func IsGitDirectory() bool {
	_, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Output()
	return err==nil
}

