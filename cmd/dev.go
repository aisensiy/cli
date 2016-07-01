package cmd

import (
	"fmt"
	"github.com/sjkyspa/stacks/client/pkg"
	"github.com/sjkyspa/stacks/client/config"
	"github.com/sjkyspa/stacks/controller/api/api"
	"github.com/sjkyspa/stacks/controller/api/net"
	"net/url"
	"io/ioutil"
	"os/exec"
	//"bytes"
	//"strings"
	"os"
	"bytes"
	"strings"
	"path/filepath"
)

func DevUp() error {

	if (!git.IsGitDirectory()) {
		return fmt.Errorf("Execute inside the app dir")
	}

	configRepository := config.NewConfigRepository(func(error) {})
	appRepository := api.NewAppRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))
	uri, err := url.Parse(configRepository.Endpoint())
	appId, err := git.DetectAppName(uri.Host)

	if err != nil || appId == "" {
		return fmt.Errorf("Please use the -remote to specfiy the app")
	}

	app, err := appRepository.GetApp(appId)
	if err != nil {
		return err
	}

	stack, err := app.GetStack()
	if err != nil {
		return err
	}

	f, err := toCompose(stack)
	if err != nil {
		return err
	}

	dockerComposeUp := exec.Command("docker-compose", "-f", f, "up", "-d")

	var out bytes.Buffer
	var errout bytes.Buffer
	dockerComposeUp.Stdin = strings.NewReader("test")
	dockerComposeUp.Stdout = &out
	dockerComposeUp.Stderr = &errout
	err = dockerComposeUp.Run()
	if err != nil {
		return err
	}

	containerId := func() string {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			panic("Please ensure the current dir can be accessed")
		}

		basename := filepath.Base(dir)
		containerNamePrefix := basename + "_" + "runtime"

		psOutput, err := exec.Command("docker-compose", "-f", f, "ps").Output()
		if err != nil {
			panic(fmt.Sprintf("Can not find the proper container id: cause %v", err))
		}
		a := string(psOutput)
		splits := strings.Split(a, "\n")
		for _, item := range splits {
			if strings.Contains(item, containerNamePrefix) {
				return strings.Split(item, " ")[0]
			}
		}
		panic("Can not find the proper container id")
	}()

	fmt.Println(fmt.Sprintf("container id %s", containerId))

	dockerExec := exec.Command("docker", "exec", "-it", containerId, "sh")
	dockerExec.Stdin = os.Stdin
	dockerExec.Stderr = os.Stderr
	dockerExec.Stdout = os.Stdout
	err = dockerExec.Run()
	if err != nil {
		return err
	}

	return nil
}

func toCompose(stack api.Stack) (string, error) {
	aa :=
	`version: '2'
services:
  runtime:
    image: hub.deepi.cn/jersey-mysql-build
    entrypoint: /bin/sh
    command: -c 'tail -f /dev/null'
    volumes:
      - /Mac/workspace/tmp/cde-stacks/jersey-mysql/template:/codee
      - /var/run/docker.sock:/var/run/docker.sock
    links:
      - mysql
  mysql:
    image: tutum/mysql
    ports:
     - 5000:5000`

	err := ioutil.WriteFile("dockercompose.yml", []byte(aa), 0600)
	if err != nil {
		return "", err
	}

	return "dockercompose.yml", nil
}

func DevDown() error {
	if (!git.IsGitDirectory()) {
		return fmt.Errorf("Execute inside the app dir")
	}

	configRepository := config.NewConfigRepository(func(error) {})
	appRepository := api.NewAppRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))
	uri, err := url.Parse(configRepository.Endpoint())
	appId, err := git.DetectAppName(uri.Host)

	if err != nil || appId == "" {
		return fmt.Errorf("Please use the -remote to specfiy the app")
	}

	app, err := appRepository.GetApp(appId)
	if err != nil {
		return err
	}

	stack, err := app.GetStack()
	if err != nil {
		return err
	}

	f, err := toCompose(stack)
	if err != nil {
		return err
	}

	dockerComposeUp := exec.Command("docker-compose", "-f", f, "stop")

	var out bytes.Buffer
	var errout bytes.Buffer
	dockerComposeUp.Stdin = strings.NewReader("test")
	dockerComposeUp.Stdout = &out
	dockerComposeUp.Stderr = &errout
	err = dockerComposeUp.Run()
	if err != nil {
		return err
	}

	fmt.Println(errout.String())
	return nil
}

func DevDestroy() error {
	if (!git.IsGitDirectory()) {
		return fmt.Errorf("Execute inside the app dir")
	}

	configRepository := config.NewConfigRepository(func(error) {})
	appRepository := api.NewAppRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))
	uri, err := url.Parse(configRepository.Endpoint())
	appId, err := git.DetectAppName(uri.Host)

	if err != nil || appId == "" {
		return fmt.Errorf("Please use the -remote to specfiy the app")
	}

	app, err := appRepository.GetApp(appId)
	if err != nil {
		return err
	}

	stack, err := app.GetStack()
	if err != nil {
		return err
	}

	f, err := toCompose(stack)
	if err != nil {
		return err
	}

	dockerComposeUp := exec.Command("docker-compose", "-f", f, "rm", "-f")

	var out bytes.Buffer
	var errout bytes.Buffer
	dockerComposeUp.Stdin = strings.NewReader("test")
	dockerComposeUp.Stdout = &out
	dockerComposeUp.Stderr = &errout
	err = dockerComposeUp.Run()
	if err != nil {
		return err
	}

	fmt.Println(errout.String())
	return nil
}