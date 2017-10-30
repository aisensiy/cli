package cmd

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/sjkyspa/cde/config"
	"github.com/sjkyspa/stacks/launcher/api/api"
	"github.com/sjkyspa/stacks/launcher/api/net"
	"os"

	"encoding/json"
	"github.com/ghodss/yaml"
	"io/ioutil"
)

func createUpsRepoository() (upsRepository api.UpsRepository) {
	configRepository := config.NewConfigRepository(func(error) {})
	upsRepository = api.NewUpsRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))
	return
}

func UpsList() error {
	upsRepository := createUpsRepoository()
	ups, err := upsRepository.GetUps()
	if err != nil || ups.Count() == 0 {
		err = fmt.Errorf("up not found")
		return err
	}

	fmt.Printf("=== Unified Procedures: [%d]\n", ups.Count())
	for _, up := range ups.Items() {
		fmt.Printf("name: %s; id: %s\n", up.Name(), up.Id())
	}
	return nil
}

func UpsInfo(upName string) error {
	upsRepository := createUpsRepoository()
	ups, err := upsRepository.GetUPByName(upName)

	if err != nil || ups.Count() == 0 {
		err = fmt.Errorf("up not found")
		return err
	}

	upId := ups.Items()[0].Id()
	up, err := upsRepository.GetUP(upId)

	outputUpDescription(up)
	outputUpBuildProcedure(up)
	return nil
}

func UpCreate(filename string) error {
	upsRepository := createUpsRepoository()

	content, err := getUpFileContent(filename)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	content, err = yaml.YAMLToJSON(content)
	upParams := make(map[string]interface{})
	err = json.Unmarshal(content, &upParams)
	if err != nil {
		return err
	}
	upModel, err := upsRepository.CreateUp(upParams)
	if err != nil {
		return err
	}

	fmt.Printf("created new UP [%s] with id %s", upModel.Name(), upModel.Id())
	return nil
}

func UpRemove(idOrName string) error {
	fmt.Println(idOrName)
	upsRepository := createUpsRepoository()

	err := upsRepository.RemoveUp(idOrName)
	if err != nil {
		fmt.Println(err)
		fmt.Printf("failed")
		return err
	}
	fmt.Printf("success")
	return nil
}

func UpUpdate(idOrName string, fileName string) error {
	upsRepository := createUpsRepoository()
	content, err := getUpFileContent(fileName)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	content, err = yaml.YAMLToJSON(content)
	upParams := make(map[string]interface{})
	err = json.Unmarshal(content, &upParams)
	if err != nil {
		return err
	}

	err = upsRepository.UpdateUp(idOrName, upParams)
	if err != nil {
		fmt.Println(err)
		fmt.Printf("failed")
		return err
	}
	fmt.Printf("success")
	return nil
}

func getUpFileContent(filename string) (content []byte, err error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return []byte{}, err
	}
	return contents, err
}

func outputUpDescription(up api.Up) {
	fmt.Println("--- Unified Procedures Detail\n")

	data := make([][]string, 2)
	data[0] = []string{"id", up.Id()}
	data[1] = []string{"name", up.Name()}

	table := tablewriter.NewWriter(os.Stdout)
	table.AppendBulk(data)
	table.Render()
}

func outputUpBuildProcedure(up api.Up) {
	fmt.Println("--- Build Procedure Detail\n")

	build, _ := up.GetProcedureByType("BUILD")

	data := make([][]string, 2)
	data[0] = []string{"id", build.Id()}
	data[1] = []string{"type", build.Type()}

	table := tablewriter.NewWriter(os.Stdout)
	table.AppendBulk(data)
	table.Render()
}
