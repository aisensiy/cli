package cmd

import (
	"fmt"
	"github.com/sjkyspa/stacks/client/pkg"
	"github.com/sjkyspa/stacks/client/config"
	"github.com/sjkyspa/stacks/controller/api/api"
	"github.com/sjkyspa/stacks/controller/api/net"
	"github.com/kr/pty"
	"net/url"
	"io/ioutil"
	"os/exec"
	"os"
	"bytes"
	"strings"
	"path/filepath"
	"github.com/sjkyspa/stacks/client/backend/compose"
	"time"
	"strconv"
	"io"
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

	dockerComposeCreate := exec.Command("docker-compose", "-f", f, "-p", app.Id(), "up", "-d")
	fi, err := pty.Start(dockerComposeCreate)
	if err != nil {
		fmt.Println(err)
		return err
	}

	io.Copy(os.Stdout, fi)
	err = dockerComposeCreate.Wait()
	if err != nil {
		fmt.Println(err)
		return err
	}

	for i := 0; i < 100; i++ {
		_, err := Pipe(
			exec.Command("docker-compose", "-f", f, "-p", app.Id(), "ps"),
			exec.Command("grep", "-E", app.Id()),
			exec.Command("/bin/sh", "-c", "grep -v Up"))

		// All the containers are started (All containers is Up, all the lines is with Up, so the grep return error.)
		if err != nil {
			break
		} else {
			fmt.Println("starting...")
			time.Sleep(1 * time.Second)
		}
	}

	containerId := func() string {
		psOutput, err := exec.Command("docker-compose", "-f", f, "-p", app.Id(), "ps", "-q", "runtime").Output()
		if err != nil {
			panic(fmt.Sprintf("Can not find the proper container id: cause %v", err))
		}
		return strings.TrimSpace(string(psOutput))
	}

	dockerExec := exec.Command("docker", "exec", "-it", containerId(), "bash")
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
	composeBackend := compose.NewComposeBackend()
	composeContent := composeBackend.ToComposeFile(stack)

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(fmt.Sprintf("Please ensure the current dir can be accessed. %v", err))
	}
	os.MkdirAll(filepath.Join(dir, ".local"), 0766)

	err = ioutil.WriteFile(".local/dockercompose.yml", []byte(composeContent), 0644)
	if err != nil {
		return "", err
	}

	return ".local/dockercompose.yml", nil
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

	dockerComposeUp := exec.Command("docker-compose", "-f", f, "-p", app.Id(), "stop")

	var out bytes.Buffer
	var errout bytes.Buffer
	dockerComposeUp.Stdout = &out
	dockerComposeUp.Stderr = &errout
	fi, err := pty.Start(dockerComposeUp)

	io.Copy(os.Stdout, fi)
	err = dockerComposeUp.Wait()
	if err != nil {
		fmt.Println(err)
		return err
	}
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

	dockerComposeUp := exec.Command("docker-compose", "-f", f, "-p", app.Id(), "down", "-v", "--remove-orphans", "--rmi", "all")

	var out bytes.Buffer
	var errout bytes.Buffer
	dockerComposeUp.Stdout = &out
	dockerComposeUp.Stderr = &errout
	fi, err := pty.Start(dockerComposeUp)
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, fi)
	err = dockerComposeUp.Wait()
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = os.RemoveAll(".local")
	if err != nil {
		fmt.Println("Error when remove the local dir .local %v", err)
		return err
	}

	fmt.Println(errout.String())
	return nil
}

func DevEnv() error {
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

	services := stack.GetServices()

	var envs []string
	var links []string

	for _, service := range services {
		if (service.IsBuildable()) {
			links = service.GetLinks()
			break
		}
	}

	for _, link := range links {
		envs = append(envs, "export " + strings.ToUpper(link + "_HOST=") + link + ";")
		envs = append(envs, "export " + strings.ToUpper(link + "_PORT=") + strconv.Itoa(services[link].GetExpose()[0]) + ";")

		linkEnvs := services[link].GetEnv()
		for name, linkEnv := range linkEnvs {
			envs = append(envs, "export " + strings.ToUpper(link + "_" + name + "=") + linkEnv + ";")
		}
	}

	for _, env := range envs {
		fmt.Println(env)
	}

	return nil
}

func Pipe(cmds ...*exec.Cmd) ([]byte, error) {
	readFile, writeFile, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	execCmds := make([]*exec.Cmd, 0)
	execCmds = append(execCmds, cmds...)
	i := 0;
	for ; i < len(execCmds) - 1; i++ {
		stdout, _ := execCmds[i].StdoutPipe()
		execCmds[i + 1].Stdin = stdout
		err := execCmds[i].Start()
		if err != nil {
			return nil, err
		}
	}
	execCmds[i].Stdout = writeFile

	err = execCmds[i].Run()
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(execCmds) - 1; i++ {
		err := execCmds[i].Wait()
		if err != nil {
			return nil, err
		}
	}

	writeFile.Close()

	content, err := ioutil.ReadAll(readFile)
	if err != nil {
		return nil, err
	}

	readFile.Close()

	return content, nil
}