package parser

import (
	"github.com/cde/client/cmd"
	"fmt"
"github.com/docopt/docopt-go"
)

func Stacks(argv []string) error {
	usage := `
Valid commands for apps:

stacks:create        create a new stack
stacks:list          list accessible stacks
stacks:info          view info about an stack

Use 'cde help [command]' to learn more.
`
	switch argv[0] {
	case "stacks:create":
		return stackCreate(argv)
	case "stacks":
		fmt.Print(usage)
		return nil
	default:
		PrintUsage()
		return nil
	}

}

func stackCreate(argv []string) error {
	usage := `
Create a stack.

Usage: deis stacks:create <stack>

Arguments:
  <stack>
    the stack name.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmd.StackCreate(safeGetValue(args, "<stack>"))
}