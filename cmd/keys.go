package cmd
import (
	"fmt"
	"path"
	"regexp"
	"strings"
	"io/ioutil"
	"github.com/cde/apisdk/api"
	"github.com/cde/apisdk/net"
	"github.com/cde/client/config"
)

// AddKey creates an key.
func AddKey(sshFilePath string) error {
	configRepository := config.NewConfigRepository(func(error) {})
	userRepo := api.NewUserRepository(configRepository, net.NewCloudControllerGateway(configRepository))

	user, err := userRepo.GetUser(configRepository.Id())
	if err != nil {
		fmt.Println(err)
		return err
	}
	public, name, err := getKey(sshFilePath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	keyParams := api.KeyParams{
		Public: public,
		Name: name,
	}
	_, err = user.UploadKey(keyParams)
	if err != nil {
		fmt.Println(err)
		return err
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
		fmt.Printf("id: %s\n", key.ID())
		fmt.Printf("fingerprint: %s\n", key.Fingerprint())
		fmt.Printf("ssh: %s\n", key.Public())
	}
}

func ListKeys(userId string) error {
	configRepository := config.NewConfigRepository(func(error) {})

	if userId == "" {
		keyRepo := api.NewKeyRepository(configRepository, net.NewCloudControllerGateway(configRepository))
		keys, err := keyRepo.GetKeys()
		printKeys(keys)
		return err
	}

	userRepo := api.NewUserRepository(configRepository, net.NewCloudControllerGateway(configRepository))
	user, err := userRepo.GetUser(userId)
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