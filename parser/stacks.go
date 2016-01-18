package parser

import (
	"github.com/cde/client/cmd"
	"github.com/docopt/docopt-go"
)

func Stacks(argv []string) error {
	usage := `
Valid commands for apps:

stacks:create        create a new stack
stacks:list          list accessible stacks

Use 'cde help [command]' to learn more.
`
	switch argv[0] {
	case "stacks:create":
		return stackCreate(argv)
	case "stacks:list":
		return stackList()
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "stacks" {
			argv[0] = "stacks:list"
			return stackList()
		}

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

func stackList() error {
	return cmd.StacksList()

}