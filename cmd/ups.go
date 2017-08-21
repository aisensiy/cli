package cmd

import (
	"fmt"
	"github.com/sjkyspa/stacks/client/config"
	"github.com/sjkyspa/stacks/launcher/api/api"
	"github.com/sjkyspa/stacks/launcher/api/net"
	"github.com/olekukonko/tablewriter"
	"os"
)

func UpsList() error {
	configRepository := config.NewConfigRepository(func(error) {})
	upsRepository := api.NewUpsRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))
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
	configRepository := config.NewConfigRepository(func(error) {})
	upsRepository := api.NewUpsRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))
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
