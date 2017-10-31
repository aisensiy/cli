package cmd

import (
	"fmt"
	"github.com/cnupp/cli/config"
	"github.com/cnupp/cnup/controller/api/api"
	"github.com/cnupp/cnup/controller/api/net"
	"io/ioutil"
	"path"
	"regexp"
	"strings"
)

// AddKey creates an key.
func AddKey(sshFilePath string) error {
	configRepository := config.NewConfigRepository(func(error) {})
	userRepo := api.NewUserRepository(configRepository, net.NewCloudControllerGateway(configRepository))

	user, err := userRepo.GetUser(configRepository.Id())
	if err != nil {
		return err
	}
	public, name, err := getKey(sshFilePath)
	if err != nil {
		return err
	}
	keyParams := api.KeyParams{
		Public: public,
		Name:   name,
	}
	_, err = user.UploadKey(keyParams)
	if err != nil {
		return fmt.Errorf("Key already exists")
	}

	fmt.Println("Upload key successfully")
	return err
}

func RemoveKey(keyId string) error {
	configRepository := config.NewConfigRepository(func(error) {})
	userRepo := api.NewUserRepository(configRepository, net.NewCloudControllerGateway(configRepository))
	user, err := userRepo.GetUser(configRepository.Id())
	if err != nil {
		return err
	}

	err = user.DeleteKey(keyId)
	if err != nil {
		return err
	}
	fmt.Println("Delete the key successfully")
	return err
}

func printKeys(keys api.Keys) {
	fmt.Printf("=== Keys [%d]\n", len(keys.Items()))

	for _, key := range keys.Items() {
		var public = key.Public()
		fmt.Printf("%s %s \"%s...%s\"\n", key.ID(), key.Name(), public[:16], public[len(public)-30:len(public)-1])
	}
}

func ListKeys() error {
	configRepository := config.NewConfigRepository(func(error) {})

	userRepo := api.NewUserRepository(configRepository, net.NewCloudControllerGateway(configRepository))
	user, err := userRepo.GetUser(configRepository.Id())
	if err != nil {
		return err
	}

	keys, err := user.Keys()
	if err != nil {
		return err
	}

	printKeys(keys)
	return err
}

func getKey(filename string) (content, name string, err error) {
	regex := regexp.MustCompile("^(ssh-...|ecdsa-[^ ]+) ([^ ]+) ?(.*)")
	contents, err := ioutil.ReadFile(filename)

	if err != nil {
		return "", "", err
	}

	if regex.Match(contents) {
		name = strings.Split(path.Base(filename), ".")[0]
		return string(contents), name, nil
	}

	return "", "", fmt.Errorf("%s is not a valid ssh key", filename)
}
