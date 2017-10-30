package cmd

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/sjkyspa/cde/config"
	"github.com/sjkyspa/stacks/controller/api/api"
	"github.com/sjkyspa/stacks/controller/api/net"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func LaunchBuild(filename, appName string) error {
	file, err := read(filename)

	if err != nil {
		return err
	}

	configRepository := config.NewConfigRepository(func(err error) {

	})
	launcherEntrypoint := configRepository.DeploymentEndpoint()

	request, errChannel, err := toRequest(file, launcherEntrypoint)
	if err != nil {
		return err
	}

	client := http.Client{}
	res, err := client.Do(request)
	if err != nil {
		if errc := <-errChannel; errc != nil {
			return errors.New(fmt.Sprintf("multiple errors happend: %s %s", errc, err))
		} else {
			return err
		}
	}

	location := res.Header.Get("Location")
	gateway := net.NewCloudControllerGateway(configRepository)
	apps := api.NewAppRepository(configRepository, gateway)

	if appName == "" {
		_, appId, err := load("")
		if err != nil {
			return err
		}
		appName = appId
	}

	app, err := apps.GetApp(appName)
	if err != nil {
		return err
	}

	build, err := app.CreateBuild(api.BuildParams{
		GitSha: "mockedsh",
		User:   "should_be_replaced_to_the_build_owner",
		Source: location,
	})

	if err != nil {
		fmt.Println("create build", err)
		return err
	}
	cerr := make(chan error)
	succeed := make(chan bool)
	cleaning := make(chan bool)

	go func() {
		for {
			select {
			case <-cleaning:
				// todo: should fetch log again
				return
			default:
				// todo: fetch log
				fmt.Println("fetch log")
				time.Sleep(time.Second * 5)
			}
		}
	}()

	go func() {
		for {
			if build.IsSuccess() {
				succeed <- true
				cleaning <- true
				return
			}
			if build.IsFail() {
				succeed <- false
				cleaning <- true
				cerr <- errors.New("Build is failed")
			}
			time.Sleep(time.Second * 5)
			build, err = app.GetBuild(build.Id())
			if err != nil {
				cerr <- err
				return
			}
		}
	}()

	for {
		select {
		case succ := <-succeed:
			if !succ {
				return errors.New("Build fail")
			} else {
				color.Green("Build Success")
				return nil
			}
		case err := <-cerr:
			return err
		default:
			time.Sleep(time.Second * 5)
		}
	}
}

func toRequest(file *os.File, entrypoint string) (*http.Request, chan error, error) {
	reader, writer := io.Pipe()
	newWriter := multipart.NewWriter(writer)
	errChannel := make(chan error, 1)
	go func() {
		defer file.Close()
		defer writer.Close()

		part, err := newWriter.CreateFormFile("file", filepath.Base(file.Name()))
		if err != nil {
			errChannel <- errors.New("unable to create multipart")
			return
		}

		if _, err := io.Copy(part, file); err != nil {
			errChannel <- err
			return
		}

		errChannel <- newWriter.Close()
	}()

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/files", entrypoint), reader)
	request.Header.Set("Content-Type", newWriter.FormDataContentType())
	return request, errChannel, err
}

func abs(filename string) (string, error) {
	if filepath.IsAbs(filename) {
		return filename, nil
	} else {
		if pwd, err := os.Getwd(); err == nil {
			return filepath.Join(pwd, filename), nil
		} else {
			return "", err
		}
	}
}

func read(filename string) (*os.File, error) {
	abs, err := abs(filename)

	if _, err := os.Stat(abs); os.IsNotExist(err) {
		return nil, errors.New(fmt.Sprintf("File %s not found", filename))
	}

	file, err := os.Open(abs)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Can not open file %s", abs))
	}
	return file, nil
}
