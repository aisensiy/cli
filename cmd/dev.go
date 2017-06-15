package cmd

import (
	"bytes"
	"fmt"
	"github.com/kr/pty"
	"github.com/sjkyspa/stacks/client/backend/compose"
	"github.com/sjkyspa/stacks/client/config"
	"github.com/sjkyspa/stacks/client/pkg"
	"github.com/sjkyspa/stacks/controller/api/api"
	"github.com/sjkyspa/stacks/controller/api/net"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"regexp"
	"errors"
)

func DevUp() error {

	if !git.IsGitDirectory() {
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

	dockerExec := exec.Command("docker", "exec", "-it", containerId(), "bash", "-c", "cd /codebase; exec ${SHELL:-bash}")
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
	if !git.IsGitDirectory() {
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
	if !git.IsGitDirectory() {
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
	if !git.IsGitDirectory() {
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
		if service.IsBuildable() {
			links = service.GetLinks()
			break
		}
	}

	for _, link := range links {
		envs = append(envs, "export "+strings.ToUpper(link+"_HOST=")+"localhost;")
		mappingPort, err := getServiceMappingPort(appId, link, services[link].GetExpose()[0])
		if err != nil {
			return err
		}
		envs = append(envs, "export "+strings.ToUpper(link+"_PORT=")+strconv.Itoa(mappingPort)+";")

		linkEnvs := services[link].GetEnv()
		for name, linkEnv := range linkEnvs {
			envs = append(envs, "export "+strings.ToUpper(link+"_"+name+"=")+linkEnv+";")
		}
	}

	for _, env := range envs {
		fmt.Println(env)
	}

	return nil
}

func getServiceMappingPort(appName string, serviceName string, port int) (int, error) {
	out, err := exec.Command("docker", "ps").Output()

	if err != nil {
		return 0, err
	}

	containerInfo := string(out)

	for _, line := range strings.Split(containerInfo, "\n") {
		if strings.Contains(line, strings.Replace(appName, "-", "", -1)+"_"+serviceName) {
			containerId := strings.Split(line, " ")[0]
			output, err := exec.Command("docker", "port", containerId).Output()

			if err != nil {
				return 0, err
			}
			mappingInfo := string(output)

			for _, infoLine := range strings.Split(mappingInfo, "\n") {
				re, _ := regexp.Compile(strconv.Itoa(port) + `/.*`)
				res := re.FindSubmatch([]byte(line))

				if len(res) > 0 {
					return strconv.Atoi(strings.Split(infoLine, ":")[1])
				}
			}
			return 0, errors.New(fmt.Sprintf("Cannot find mapping port for service %s port %d", serviceName, strconv.Itoa(port)))
		}
	}

	return 0, errors.New("Cannot find container for service " + serviceName)
}

func Pipe(cmds ...*exec.Cmd) ([]byte, error) {
	readFile, writeFile, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	execCmds := make([]*exec.Cmd, 0)
	execCmds = append(execCmds, cmds...)
	i := 0
	for ; i < len(execCmds)-1; i++ {
		stdout, _ := execCmds[i].StdoutPipe()
		execCmds[i+1].Stdin = stdout
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

	for i := 0; i < len(execCmds)-1; i++ {
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
