package parser
import (
	"os"
	"fmt"
	"strconv"
	"github.com/cde/client/cmd"
	docopt "github.com/docopt/docopt-go"
	"errors"
)

func Apps(argv []string) error {
	usage := `
Valid commands for apps:

apps:create        create a new application
apps:list          list accessible applications
apps:info          view info about an application

Use 'cde help [command]' to learn more.
`
	switch argv[0] {
	case "apps:create":
		return appCreate(argv)
	case "apps:list":
		return appList(argv)
	case "apps:info":
		return appInfo(argv)
	case "apps":
		fmt.Print(usage)
		return nil
	default:
		if printHelp(argv, usage) {
			return nil
		}
		PrintUsage()
		return nil
	}

}

func appCreate(argv []string) error {
	usage := `
Creates a new application.

Usage: cde apps:create <name> <stack> <owner> [options]

Arguments:
  <name>
  	a uniquely identifiable name for the application. No other app can already
    exist with this name.
  <stack>

Options:
  --mem
  	allocated memory for this app. [default: 512]
  --disk
  	max allocated disk size. [default: 20]
  --instances
  	default started instance number. [default: 1]
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	name := safeGetValue(args, "<name>")
	stack := safeGetValue(args, "<stack>")
	if stack == "" || name == "" {
		return errors.New("<name> <stack> are essential parameters")
	}
	fmt.Println(args)
	memory := safeGetOrDefault(args, "--mem", "512")
	disk := safeGetOrDefault(args, "--disk", "20")
	instances := safeGetOrDefault(args, "--instances", "1")

	var mem, ins, diskSize int

	if mem, err = strconv.Atoi(memory); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}

	if ins, err = strconv.Atoi(instances); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}

	if diskSize, err = strconv.Atoi(disk); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}
	return cmd.AppCreate(name, stack, owner, mem, diskSize, ins)
}

func appList(argv []string) error {
	usage := `
Lists applications visible to the current user.

Usage: cde apps:list [options]

Options:
  -l --limit=<num>
    the maximum number of results to display, defaults to config setting
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	results, err:= strconv.Atoi(safeGetOrDefault(args, "--limit", "20"))

	if err != nil {
		return err
	}

	return cmd.AppsList(results)
}

func appInfo(argv []string) error {
	return nil
}
