package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/fatih/color"
	"github.com/sjkyspa/stacks/client/parser"
	"github.com/sjkyspa/stacks/version"
	"os"
	"strings"
	"github.com/urfave/cli"
	"github.com/sjkyspa/stacks/client/cmd"
	"time"
	"errors"
)

func main() {
	app := cli.NewApp()
	app.Name = "CDE"
	app.Version = "0.1.4"
	app.Usage = "Cloud Development Environment"
	app.Description = "CDE command line tool"
	app.Compiled = time.Now()
	app.Author = "ThoughtWorks"
	app.EnableBashCompletion = true

	app.Commands = []cli.Command{
		upsCommand(),
		stacksCommand(),
		providersCommand(),
	}

	commandList := os.Args

	if len(commandList) > 1 &&
		!strings.Contains(commandList[1], "ups") &&
		!strings.Contains(commandList[1], "stacks") {
		os.Exit(Command(commandList[1:]))
	} else {
		commandList = preProcessCommand(commandList)
		app.Run(commandList)
	}
}

func preProcessCommand(args []string) (processedArgs []string) {
	if len(args) == 1 {
		return args
	}

	processedArgs = append([]string{args[0]}, strings.Split(args[1], ":")...)
	processedArgs = append(processedArgs, args[2:]...)
	return
}

func upsCommand() cli.Command {
	return cli.Command{
		Name:  "ups",
		Usage: "Unified Procedures Commands",
		Subcommands: []cli.Command{
			{
				Name:      "list",
				Usage:     "list all Unified Procedures",
				ArgsUsage: " ",
				Action: func(c *cli.Context) error {
					return cmd.UpsList()
				},
			},
			{
				Name:      "info",
				Usage:     "get info of an Unified Procedure",
				ArgsUsage: "<up-name>",
				Action: func(c *cli.Context) error {
					return cmd.UpsInfo(c.Args().First())
				},
			},
			{
				Name:      "draft",
				Usage:     "create a new Unified Procedure",
				ArgsUsage: "<up-file>",
				Action: func(c *cli.Context) error {
					return cmd.UpCreate(c.Args().First())
				},
			},
			{
				Name:      "remove",
				Usage:     "delete an existing Unified Procedure",
				ArgsUsage: "<up-id>",
				Action: func(c *cli.Context) error {
					return cmd.UpRemove(c.Args().First())
				},
			},
			{
				Name:      "update",
				Usage:     "update an existing Unified Procedure",
				ArgsUsage: "<up-id> <up-file>",
				Action: func(c *cli.Context) error {
					return cmd.UpUpdate(c.Args().First(), c.Args().Get(1))
				},
			},
		},
	}
}

func stacksCommand() cli.Command {
	return cli.Command{
		Name:  "stacks",
		Usage: "Stacks Commands",
		Subcommands: []cli.Command{
			{
				Name:      "list",
				Usage:     "list all Stacks",
				ArgsUsage: " ",
				Action: func(c *cli.Context) error {
					return cmd.StacksList()
				},
			},
			{
				Name:      "info",
				Usage:     "get info of a Stack",
				ArgsUsage: "<stack-name>",
				Action: func(c *cli.Context) error {
					return cmd.GetStack(c.Args().First())
				},
			},
			{
				Name:      "create",
				Usage:     "create a new Stack",
				ArgsUsage: "<stack-file>",
				Action: func(c *cli.Context) error {
					return cmd.StackCreate(c.Args().First())
				},
			},
			{
				Name:      "update",
				Usage:     "update an existing Stack",
				ArgsUsage: "<stack-id> <stack-file>",
				Action: func(c *cli.Context) error {
					return cmd.StackUpdate(c.Args().First(), c.Args().Get(1))
				},
			},
			{
				Name:      "remove",
				Usage:     "remove a Stack",
				ArgsUsage: "<stack-name>",
				Action: func(c *cli.Context) error {
					return cmd.StackRemove(c.Args().First())
				},
			},
			{
				Name:      "publish",
				Usage:     "publish a Stack",
				ArgsUsage: "<stack-id>",
				Action: func(c *cli.Context) error {
					return cmd.StackPublish(c.Args().First())
				},
			},
			{
				Name:      "unpublish",
				Usage:     "unpublish a Stack",
				ArgsUsage: "<stack-id>",
				Action: func(c *cli.Context) error {
					return cmd.StackUnPublish(c.Args().First())
				},
			},
		},
	}
}

func providersCommand() cli.Command {
	return cli.Command{
		Name:  "providers",
		Usage: "Providers Commands",
		Subcommands: []cli.Command{
			{
				Name:      "enroll",
				Usage:     "Enroll a new Provider",
				ArgsUsage: "<name> <type> [-c <config>]",
				Flags: []cli.Flag{
					cli.StringSliceFlag{
						Name: "config",
					},
				},
				Action: func(c *cli.Context) error {
					configMap, _ := configConvert(c.StringSlice("config"))
					return cmd.ProviderCreate(c.Args().Get(0), c.Args().Get(1), "", configMap)
				},
			},
		},
	}
}

func configConvert(config []string) (map[string]interface{}, error) {
	configMap := map[string]interface{}{}
	for _, v := range config {
		pair := strings.Split(v, "=")
		if len(pair) != 2 {
			return nil, errors.New("invalid config format")
		}
		configMap[pair[0]] = pair[1]
	}
	return configMap, nil
}

func Command(argv []string) int {
	usage := `
The CDE command-line
Usage: cde <command> [<args>...]
Use 'git push cde master' to deploy to an application.

Auth commands:

  register      register a new user with a controller
  login         login to a controller
  logout        logout from the current controller
  whoami        display the current user

Subcommands, use 'cde help [subcommand]' to learn more::

  apps          manage applications used to provide services
  clusters      manage clusters used to provide services
  orgs          manage organizations
  scaffold      create scaffold project quickly
  stacks        manage stacks
  domains       manage domains
  services      manage services instances in marathon
  routes      	manage routes
  keys      	manage keys
  git           manage git for applications
  config        manage environment variables that define app config
  ps            manage process status
  providers 	manage providers
`
	command, argv := parseArgs(argv)

	_, err := docopt.Parse(usage, []string{command}, false, version.Version(), true, false)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if len(argv) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: cde <command> [<args>...]")
		return 1
	}

	// Dispatch the command, passing the argv through so subcommands can
	// re-parse it according to their usage strings.
	switch command {
	case "apps":
		err = parser.Apps(argv)
	case "orgs":
		err = parser.Orgs(argv)
	case "scaffold":
		err = parser.Scaffold(argv)
	case "auth":
		err = parser.Auth(argv)
	case "domains":
		err = parser.Domains(argv)
	case "services":
		err = parser.Service(argv)
	case "stacks":
		err = parser.Stacks(argv)
	case "routes":
		err = parser.Routes(argv)
	case "keys":
		err = parser.Keys(argv)
	case "ps":
		err = parser.Ps(argv)
	case "git":
		err = parser.Git(argv)
	case "config":
		err = parser.Config(argv)
	case "dev":
		err = parser.Dev(argv)
	case "clusters":
		err = parser.Clusters(argv)
	case "providers":
		err = parser.Providers(argv)
	case "launch":
		err = parser.Launch(argv)
	case "help":
		fmt.Print(usage)
		return 0
	case "--version":
		return 0
	default:
		fmt.Fprintln(os.Stderr, "Usage: cde <command> [<args>...]")
	}

	if err != nil {
		color.Set(color.FgRed)
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		color.Unset()
		return 1
	}
	return 0
}

// parseArgs returns the provided args with "--help" as the last arg if need be,
// expands shortcuts and formats commands to be properly routed.
func parseArgs(argv []string) (string, []string) {

	if len(argv) == 0 {
		return "", argv
	}

	if len(argv) == 1 && (argv[0] == "-h" || argv[0] == "--help") {
		// rearrange "cde --help" as "cde help"
		argv[0] = "help"
	}

	if len(argv) >= 2 && (argv[0] == "help" || argv[0] == "-h" || argv[0] == "--help") {
		argv = append(argv[1:], "--help")
	}

	argv[0] = replaceShortcut(argv[0])
	return pickMainCommand(argv[0]), argv
}

func pickMainCommand(command string) string {
	return strings.Split(command, ":")[0]
}

func replaceShortcut(command string) string {
	shortcuts := map[string]string{
		"create":   "apps:create",
		"info":     "apps:info",
		"login":    "auth:login",
		"logout":   "auth:logout",
		"register": "auth:register",
		"whoami":   "auth:whoami",
	}

	if expandedCommand, ok := shortcuts[command]; ok {
		return expandedCommand
	}
	return command
}
