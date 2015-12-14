package parser

import (
	"github.com/cde/client/cmd"
)

func Stacks(argv []string) error {
	switch argv[0] {
	case "init":
		return stackCreate(argv)
	default:
		PrintUsage()
		return nil
	}

}

func stackCreate(argv []string) error {
	return cmd.StackCreate()
}