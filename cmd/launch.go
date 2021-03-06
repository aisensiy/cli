package cmd

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/cnupp/cli/config"
	"github.com/cnupp/appssdk/api"
	"github.com/cnupp/appssdk/net"
	runtimeApi "github.com/cnupp/runtimesdk/api"
	runtimeNet "github.com/cnupp/runtimesdk/net"
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

func LaunchVerify(buildId, appName string) error {
	configRepository := config.NewConfigRepository(func(err error) {

	})
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
	build, err := app.GetBuild(buildId)
	if err != nil {
		return err
	}

	verify, err := build.CreateVerify(api.VerifyParams{})
	if err != nil {
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
			if verify.IsSuccess() {
				succeed <- true
				cleaning <- true
				return
			}
			if verify.IsFail() {
				succeed <- false
				cleaning <- true
				cerr <- errors.New("Verify is failed")
			}
			time.Sleep(time.Second * 5)
			verify, err = build.GetVerify(verify.Id())
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
				return errors.New("Verify fail")
			} else {
				color.Green("Verify Success")
				return nil
			}
		case err := <-cerr:
			return err
		default:
			time.Sleep(time.Second * 5)
		}
	}
}

func LaunchDeployment(releaseId, appName string, providerName string) error {
	configRepository := config.NewConfigRepository(func(err error) {

	})
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
	release, err := app.GetRelease(releaseId)
	if err != nil {
		return err
	}

	runtimeGateway := runtimeNet.NewCloudControllerGateway(configRepository)
	upsRepository := runtimeApi.NewUpsRepository(configRepository, runtimeGateway)
	upLink, err := app.Links().Link("unified_procedure")
	if err != nil {
		return err
	}
	up, err := upsRepository.GetUpByUri(upLink.URI)
	if err != nil {
		return err
	}

	var procedure runtimeApi.Procedure
	procedure, err = up.GetProcedureByType("RUN")
	if err != nil {
		return err
	}

	providerRepository := runtimeApi.NewProviderRepository(configRepository, runtimeGateway)
	var provider runtimeApi.Provider
	if providerName == "" {
		providerLink, err := app.Links().Link("provider")
		if err != nil {
			return err
		}
		provider, err = providerRepository.GetProviderByUri(providerLink.URI)
	} else {
		provider, err = providerRepository.GetProviderByName(providerName)
	}
	if err != nil {
		return err
	}

	providerParam := make(map[string]interface{})
	providerParam["id"] = provider.ID()

	procedureAppParam := make(map[string]interface{})
	procedureAppParam["image"] = release.ImageName() + ":" + release.Version()

	procedureRuntimeParam := make(map[string]interface{})

	procedureParam := make(map[string]interface{})
	procedureParam["app"] = procedureAppParam
	procedureParam["runtime"] = procedureRuntimeParam

	ownerParam := make(map[string]interface{})
	ownerParam["id"] = app.Id()
	ownerParam["name"] = app.Name()

	params := make(map[string]interface{})
	params["provider"] = providerParam
	params["procedure"] = procedureParam
	params["owner"] = ownerParam

	instance, err := procedure.CreateInstance(params)
	if err != nil {
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
			if instance.Status() == "SUCCEED" {
				succeed <- true
				cleaning <- true
				return
			}
			if instance.Status() == "FAILED" {
				succeed <- false
				cleaning <- true
				cerr <- errors.New("Deployment is failed")
			}
			time.Sleep(time.Second * 5)
			instance, err = upsRepository.GetProcedureInstance(instance.Id())
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
				return errors.New("Deployment Fail")
			} else {
				color.Green("Deployment Success")
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
