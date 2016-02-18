package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/sjkyspa/stacks/Godeps/_workspace/src/golang.org/x/crypto/ssh/terminal"
	"github.com/sjkyspa/stacks/apisdk/api"
	"github.com/sjkyspa/stacks/apisdk/net"
	"github.com/sjkyspa/stacks/client/config"
	"net/http"
	"net/url"
	"syscall"
)

func Login(controller string, email string, password string) error {
	controllerURL, err := checkController(controller)

	if email == "" {
		fmt.Print("email: ")
		fmt.Scanln(&email)
	}

	if password == "" {
		fmt.Print("password: ")
		password, err = readPassword()
		fmt.Println()

		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	configRepository := config.NewConfigRepository(func(err error) {})
	configRepository.SetApiEndpoint(controllerURL)
	configRepository.SetDeploymentEndpoint("http://192.168.50.4:31089/deployment")

	return doLogin(email, password)
}

func checkController(controller string) (string, error) {
	u, err := url.Parse(controller + "/build")

	if err != nil {
		return "", err
	}

	controllerURL, err := chooseScheme(*u)

	if err != nil {
		return "", err
	}
	fmt.Println(controllerURL.String())
	if err = CheckConnection(CreateHTTPClient(), controllerURL); err != nil {
		return "", err
	}

	return controllerURL.String(), nil
}

func Register(controller string, email string, password string) error {
	controllerURL, err := checkController(controller)

	if err != nil {
		return err
	}

	if email == "" {
		fmt.Print("email: ")
		fmt.Scanln(&email)
	}

	if password == "" {
		fmt.Print("password: ")
		password, err = readPassword()
		fmt.Printf("\npassword (confirm): ")
		passwordConfirm, err := readPassword()
		fmt.Println()

		if err != nil {
			return err
		}

		if password != passwordConfirm {
			return errors.New("Password mismatch, aborting registration.")
		}
	}

	if email == "" {
		fmt.Print("email: ")
		fmt.Scanln(&email)
	}

	configRepository := config.NewConfigRepository(func(err error) {})
	configRepository.SetApiEndpoint(controllerURL)
	configRepository.SetDeploymentEndpoint("http://192.168.50.4:31089")

	userRepository := api.NewUserRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))
	userParams := api.UserParams{
		Email:    email,
		Password: password,
	}
	_, err = userRepository.Create(userParams)
	if err != nil {
		return err
	}

	fmt.Printf("Registered %s\n", email)
	return doLogin(email, password)
}

func chooseScheme(u url.URL) (url.URL, error) {
	if u.Scheme == "" {
		u.Scheme = "http"
		u, err := url.Parse(u.String())

		if err != nil {
			return url.URL{}, err
		}

		return *u, nil
	}

	return u, nil
}

func CheckConnection(client *http.Client, apiURL url.URL) error {
	errorMessage := `%s does not appear to be a valid Cde controller.
Make sure that the Controller URI is correct and the server is running.`

	baseURL := apiURL.String()

	apiURL.Path = apiURL.Path + "/apps"

	req, err := http.NewRequest("GET", apiURL.String(), bytes.NewBuffer(nil))

	if err != nil {
		return err
	}

	res, err := client.Do(req)

	if err != nil {
		fmt.Printf(errorMessage+"\n", baseURL)
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf(errorMessage, baseURL)
	}

	return nil
}

func CreateHTTPClient() *http.Client {
	tr := &http.Transport{
		DisableKeepAlives: true,
	}
	return &http.Client{Transport: tr}
}

func doLogin(email string, password string) error {
	configRepository := config.NewConfigRepository(func(err error) {})
	configRepository.SetGitHost("192.168.50.6")
	authRepository := api.NewAuthRepository(
		configRepository,
		net.NewCloudControllerGateway(configRepository))
	userParams := api.UserParams{
		Email:    email,
		Password: password,
	}
	auth, err := authRepository.Create(userParams)

	if err != nil {
		return err
	}

	userRepo := api.NewUserRepository(configRepository, net.NewCloudControllerGateway(configRepository))
	user, err := userRepo.GetUserByEmail(email)
	if err != nil {
		return err
	}
	userId := user.Items()[0].Id()
	configRepository.SetEmail(auth.UserEmail())
	configRepository.SetId(userId)
	configRepository.SetAuth(auth.Id())
	fmt.Printf("Welcome %s\n", auth.UserEmail())
	return nil
}

func readPassword() (string, error) {
	password, err := terminal.ReadPassword(int(syscall.Stdin))

	return string(password), err
}

func Logout() error {
	configRepository := config.NewConfigRepository(func(err error) {})
	token := configRepository.Auth()
	authRepository := api.NewAuthRepository(
		configRepository,
		net.NewCloudControllerGateway(configRepository))
	err := authRepository.Delete(token)

	if err != nil {
		return err
	}
	fmt.Println("Logout successfully")
	return nil
}

func Cancel(email string, password string, force bool) error {
	return nil
}

func Regenerate() error {
	return nil
}

func Whoami() error {
	configRepository := config.NewConfigRepository(func(err error) {})

	fmt.Printf("You are %s at %s\n", configRepository.Email(), configRepository.ApiEndpoint())

	return nil
}
