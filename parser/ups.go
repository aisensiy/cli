package parser

import (
	"fmt"
	"github.com/sjkyspa/stacks/client/cmd"
)

func Ups(argv []string) error {
	usage := `Valid commands for ups:

cde ups:draft unified_procedure.yml
cde ups:publish <ups_name>
cde ups:list
cde ups:info <unified_procedure_name>
cde ups:deprecate <ups_name>
cde ups:update <ups_name> unified_procedure.yml
cde ups:remove <ups_name>

Use 'cde help [command]' to learn more.`

	switch argv[0] {
	case "ups:draft":
		return upsCreate(argv)
	case "ups:publish":
		return upsPublish(argv)
	case "ups:list":
		return upsList()
	case "ups:info":
		return upsInfo(argv)
	case "ups:deprecate":
		return upsDeprecate(argv)
	case "ups:update":
		return upsUpdate(argv)
	case "ups:remove":
		return upsRemove(argv)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "ups" {
			return stackList()
		}

		PrintUsage()
		return nil
	}
}

func upsCreate(argv []string) error {
	fmt.Println("TODO ups:draft command")
	return nil
}

func upsPublish(argv []string) error {
	fmt.Println("TODO ups:publish command")
	return nil
}

func upsList() error {
	return cmd.UpsList()
}

func upsInfo(argv []string) error {
	fmt.Println("TODO ups:info command")
	return nil
}

func upsDeprecate(argv []string) error {
	fmt.Println("TODO ups:deprecate command")
	return nil
}

func upsUpdate(argv []string) error {
	fmt.Println("TODO ups:update command")
	return nil
}

func upsRemove(argv []string) error {
	fmt.Println("TODO ups:remove command")
	return nil
}
