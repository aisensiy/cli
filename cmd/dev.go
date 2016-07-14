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
	"os"
	"bytes"
	"strings"
	"path/filepath"
	"github.com/sjkyspa/stacks/client/backend/compose"
	"time"
	"bufio"
	"sync"
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

	dockerComposeCreate := exec.Command("docker-compose", "-f", f, "create")
	stdoutReadFile, stdoutWriteFile, err := os.Pipe()
	stderrReadFile, stderrWriteFile, err := os.Pipe()
	dockerComposeCreate.Stdout = stdoutWriteFile
	dockerComposeCreate.Stderr = stderrWriteFile
	err = dockerComposeCreate.Start()
	if err != nil {
		return err
	}
	createDone := make(chan bool, 1)
	go func() {
		dockerComposeCreate.Wait()
		createDone <- true
	}()

	var wg sync.WaitGroup
	go (func() {
		wg.Add(1)
		defer wg.Done()
		stdoutStop := make(chan int, 1)
		stderrStop := make(chan int, 1)
		go func() {
			for {
				select {
				case <-stdoutStop:
					return
				default:
					stdout := bufio.NewScanner(stdoutReadFile)
					if stdout.Scan() {
						fmt.Println(stdout.Text())
					}
				}
			}
		}()
		go func() {
			for {
				select {
				case <-stderrStop:
					return
				default:
					stderr := bufio.NewScanner(stderrReadFile)
					if stderr.Scan() {
						fmt.Fprintln(os.Stderr, string(stderr.Bytes()))
					}
				}
			}
		}()
		for {
			select {
			case <-createDone:
				stdoutStop <- 0
				stderrStop <- 0
				return
			default:
				fmt.Println("Creating ...")
				time.Sleep(5 * time.Second)
			}
		}
	})()

	wg.Wait()
	dockerComposeUp := exec.Command("docker-compose", "-f", f, "up", "-d")

	var composeUpOut bytes.Buffer
	var composeUpErr bytes.Buffer
	dockerComposeUp.Stdout = &composeUpOut
	dockerComposeUp.Stderr = &composeUpErr
	err = dockerComposeUp.Run()
	if err != nil {
		return err
	}

	for i := 0; i < 100; i++ {
		content, err := Pipe(
			exec.Command("docker-compose", "-f", f, "ps"),
			exec.Command("grep", "-E", "runtime|db"),
			exec.Command("/bin/sh", "-c", "grep -v Up||true"),
			exec.Command("wc", "-l"))
		if err != nil {
			fmt.Fprint(os.Stderr, fmt.Sprintf("error occured: %v", err))
			return err
		}
		if "0" == strings.TrimSpace(string(content)) {
			break
		} else {
			fmt.Println("starting...")
			time.Sleep(1 * time.Second)
		}
	}

	psOutput, err := exec.Command("docker-compose", "-f", f, "logs").Output()
	if err != nil {
		panic(fmt.Sprintf("Can not find the proper container id: cause %v", err))
	}

	fmt.Println(string(psOutput))

	containerId := func() string {
		containerNamePrefix := "local" + "_" + "runtime"

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

	dockerComposeUp := exec.Command("docker-compose", "-f", f, "stop")

	var out bytes.Buffer
	var errout bytes.Buffer
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

	dockerComposeUp := exec.Command("docker-compose", "-f", f, "down", "-v", "--remove-orphans", "--rmi", "all")

	var out bytes.Buffer
	var errout bytes.Buffer
	dockerComposeUp.Stdout = &out
	dockerComposeUp.Stderr = &errout
	err = dockerComposeUp.Run()
	if err != nil {
		return err
	}

	fmt.Println(errout.String())
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