package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"github.com/sjkyspa/stacks/controller/api/api"
	"github.com/sjkyspa/stacks/controller/api/net"
	"github.com/sjkyspa/stacks/client/config"
	"net/http"
	"net/url"
	"syscall"
)

func Login(controller string, email string, password string) error {
	formalizedURL, err := formalizeURL(controller)

	if err = IsValidController(formalizedURL); err != nil {
		return err
	}

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
	configRepository.SetEndpoint(formalizedURL.String())

	return doLogin(email, password)
}

func formalizeURL(controller string) (url.URL, error) {
	u, err := url.Parse(controller)
	if err != nil {
		return url.URL{}, err
	}

	scheme, err := chooseScheme(*u)

	if err != nil {
		return url.URL{}, err
	}

	u.Scheme = scheme

	return *u, nil
}

func Register(controller string, email string, password string) error {
	formalizedURL, err := formalizeURL(controller)

	if err != nil {
		return err
	}

	if err = IsValidController(formalizedURL); err != nil {
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
	configRepository.SetEndpoint(formalizedURL.String())

	userRepository := api.NewUserRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))
	userParams := api.UserParams{
		Email:    email,
		Password: password,
	}
	err = userRepository.Create(userParams)
	if err != nil {
		return err
	}

	fmt.Printf("Registered %s\n", email)
	return doLogin(email, password)
}

func chooseScheme(u url.URL) (string, error) {
	if u.Scheme == "" {
		u.Scheme = "http"
		_, err := url.Parse(u.String())

		if err != nil {
			return "", err
		}

		return "http", nil
	}

	return u.Scheme, nil
}

func IsValidController(apiURL url.URL) error {
	errorMessage := `%s does not appear to be a valid Cde controller.
Make sure that the Controller URI is correct and the server is running.`

	var createHttpClient = func() *http.Client {
		tr := &http.Transport{
			DisableKeepAlives: true,
		}
		return &http.Client{Transport: tr}
	}

	baseURL := apiURL.String()

	apiURL.Path += "/"

	req, err := http.NewRequest("GET", apiURL.String(), bytes.NewBuffer(nil))

	if err != nil {
		return err
	}

	res, err := createHttpClient().Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf(errorMessage, baseURL)
	}

	return nil
}

func doLogin(email string, password string) error {
	configRepository := config.NewConfigRepository(func(err error) {})
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

	fmt.Printf("You are %s at %s\n", configRepository.Email(), configRepository.Endpoint())

	return nil
}
