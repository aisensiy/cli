package parser

import (
	"github.com/cde/client/cmd"
)

func Service(argv []string) error {
	switch argv[0] {
	case "create":
		return serviceCreate(argv)
	default:
		PrintUsage()
		return nil
	}

}

func serviceCreate(argv []string) error {
	return cmd.ServiceCreate()
}

