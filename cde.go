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
		authCommand(),
	}

	commandList := os.Args

	if len(commandList) > 1 &&
		!strings.Contains(commandList[1], "ups") &&
		!strings.Contains(commandList[1], "stacks") &&
		!strings.Contains(commandList[1], "providers") &&
		!strings.Contains(commandList[1], "login") &&
		!strings.Contains(commandList[1], "create") &&
		!strings.Contains(commandList[1], "info") &&
		!strings.Contains(commandList[1], "logout") &&
		!strings.Contains(commandList[1], "whoami") &&
		!strings.Contains(commandList[1], "register") {
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

	args[1] = replaceShortcut(args[1])

	//TODO: filter command name, because only some commands have list subcommands
	if len(args) == 2 && !strings.Contains(args[1], ":") {
		args[1] = args[1] + ":list"
	}


	processedArgs = append([]string{args[0]}, strings.Split(args[1], ":")...)
	processedArgs = append(processedArgs, args[2:]...)
	return
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

func authCommand() cli.Command {
	return cli.Command{
		Name:  "auth",
		Usage: "Auth Commands",
		Subcommands: []cli.Command{
			{
				Name:      "register",
				Usage:     "Register a new user on a specific controller",
				ArgsUsage: "[controller]",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "email, e",
						Usage: "Provide email for the new user",
					},
					cli.StringFlag{
						Name:  "password, p",
						Usage: "Provide password for the new user",
					},
				},
				Action: func(c *cli.Context) error {
					if c.Args().Get(0) == "" {
						return cli.NewExitError(fmt.Sprintf("USAGE: %s %s", c.Command.HelpName, c.Command.ArgsUsage), 1)
					}
					err := cmd.Register(c.Args().First(), c.String("email"), c.String("password"))
					if err != nil {
						return cli.NewExitError(err, 1)
					}
					return nil
				},
			},
			{
				Name:      "whoami",
				Usage:     "Display current user.",
				ArgsUsage: " ",
				Action: func(c *cli.Context) error {
					return cmd.Whoami()
				},
			},
			{
				Name:      "login",
				Usage:     "Log in on a specific controller.",
				ArgsUsage: "[controller]",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "email, e",
						Usage: "Provide user email",
					},
					cli.StringFlag{
						Name:  "password, p",
						Usage: "Provide user password",
					},
					cli.StringFlag{
						Name:  "ssl-verify, s",
						Usage: "Disable SSL certificate verification for API requests",
					},
				},
				Action: func(c *cli.Context) error {
					if c.Args().Get(0) == "" {
						return cli.NewExitError(fmt.Sprintf("USAGE: %s %s", c.Command.HelpName, c.Command.ArgsUsage), 1)
					}
					err := cmd.Login(c.Args().First(), c.String("email"), c.String("password"))
					if err != nil {
						return cli.NewExitError(err, 1)
					}
					return nil
				},
			},
			{
				Name:      "logout",
				Usage:     "Log out from a controoler",
				ArgsUsage: " ",
				Action: func(c *cli.Context) error {
					err :=  cmd.Logout()
					if err != nil {
						return cli.NewExitError(err, 1)
					}
					return nil
				},
			},
		},
	}
}

func upsCommand() cli.Command {
	return cli.Command{
		Name:  "ups",
		Usage: "Unified Procedures Commands",
		Subcommands: []cli.Command{
			{
				Name:      "list",
				Usage:     "List all Unified Procedures",
				ArgsUsage: " ",
				Action: func(c *cli.Context) error {
					return cmd.UpsList()
				},
			},
			{
				Name:      "info",
				Usage:     "Get info of an Unified Procedure",
				ArgsUsage: "[up-name]",
				Action: func(c *cli.Context) error {
					return cmd.UpsInfo(c.Args().First())
				},
			},
			{
				Name:      "draft",
				Usage:     "Create a new Unified Procedure",
				ArgsUsage: "[up-file]",
				Action: func(c *cli.Context) error {
					return cmd.UpCreate(c.Args().First())
				},
			},
			{
				Name:      "update",
				Usage:     "Update an existing Unified Procedure",
				ArgsUsage: "[up-id] [up-file]",
				Action: func(c *cli.Context) error {
					return cmd.UpUpdate(c.Args().First(), c.Args().Get(1))
				},
			},
			{
				Name:      "remove",
				Usage:     "Delete an Unified Procedure",
				ArgsUsage: "[up-id]",
				Action: func(c *cli.Context) error {
					return cmd.UpRemove(c.Args().First())
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
				Usage:     "List all Stacks",
				ArgsUsage: " ",
				Action: func(c *cli.Context) error {
					err := cmd.StacksList()
					if err != nil {
						return cli.NewExitError(err, 1)
					}
					return nil
				},
			},
			{
				Name:      "info",
				Usage:     "Get info of a Stack",
				ArgsUsage: "[stack-name]",
				Action: func(c *cli.Context) error {
					err := cmd.GetStack(c.Args().First())
					if err != nil {
						return cli.NewExitError(err, 1)
					}
					return nil
				},
			},
			{
				Name:      "create",
				Usage:     "Create a new Stack",
				ArgsUsage: "[stack-file]",
				Action: func(c *cli.Context) error {
					err := cmd.StackCreate(c.Args().First())
					if err != nil {
						return cli.NewExitError(err, 1)
					}
					return nil
				},
			},
			{
				Name:      "update",
				Usage:     "Update an existing Stack",
				ArgsUsage: "[stack-id] [stack-file]",
				Action: func(c *cli.Context) error {
					err := cmd.StackUpdate(c.Args().First(), c.Args().Get(1))
					if err != nil {
						return cli.NewExitError(err, 1)
					}
					return nil
				},
			},
			{
				Name:      "remove",
				Usage:     "Delete a Stack",
				ArgsUsage: "[stack-name]",
				Action: func(c *cli.Context) error {
					err := cmd.StackRemove(c.Args().First())
					if err != nil {
						return cli.NewExitError(err, 1)
					}
					return nil
				},
			},
			{
				Name:      "publish",
				Usage:     "Publish a Stack",
				ArgsUsage: "[stack-id]",
				Action: func(c *cli.Context) error {
					err := cmd.StackPublish(c.Args().First())
					if err != nil {
						return cli.NewExitError(err, 1)
					}
					return nil
				},
			},
			{
				Name:      "unpublish",
				Usage:     "Unpublish a Stack",
				ArgsUsage: "[stack-id]",
				Action: func(c *cli.Context) error {
					err := cmd.StackUnPublish(c.Args().First())
					if err != nil {
						return cli.NewExitError(err, 1)
					}
					return nil

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
				Name:      "list",
				Usage:     "List all Providers",
				ArgsUsage: " ",
				Action: func(c *cli.Context) error {
					err := cmd.ProviderList()
					if err != nil {
						return cli.NewExitError(err, 1)
					}
					return nil
				},
			},
			{
				Name:      "info",
				Usage:     "Get info of a Provider",
				ArgsUsage: "[provider-name]",
				Action: func(c *cli.Context) error {
					if (c.Args().Get(0) == "") {
						return cli.NewExitError(fmt.Sprintf("USAGE: %s %s", c.Command.HelpName, c.Command.ArgsUsage), 1)
					}
					err := cmd.GetProviderByName(c.Args().Get(0))
					if err != nil {
						return cli.NewExitError(err, 1)
					}
					return nil
				},
			},
			{
				Name:      "enroll",
				Usage:     "Enroll a new Provider",
				ArgsUsage: "[name] [type]",
				Flags: []cli.Flag{
					cli.StringSliceFlag{
						Name:  "config, c",
						Usage: "Set provider's configuration. Key \"endpoint\" is required.",
					},
					cli.StringFlag{
						Name:  "for, f",
						Usage: "Specify an organization for the provider.",
					},
				},
				Action: func(c *cli.Context) error {
					if c.Args().Get(0) == "" || c.Args().Get(1) == "" {
						return cli.NewExitError(fmt.Sprintf("USAGE: %s %s", c.Command.HelpName, c.Command.ArgsUsage), 1)
					}
					configMap, err := enrollConfigConvert(c.StringSlice("config"))
					if err != nil {
						return cli.NewExitError(err, 1)
					}
					consumer := c.String("for")
					err = cmd.ProviderCreate(c.Args().Get(0), c.Args().Get(1), consumer, configMap)
					if err != nil {
						return cli.NewExitError(err, 1)
					}
					return nil
				},
			},
			{
				Name:      "update",
				Usage:     "Update an existing Provider",
				ArgsUsage: "[name]",
				Flags: []cli.Flag{
					cli.StringSliceFlag{
						Name:  "config, c",
						Usage: "Set provider's configuration. Key \"endpoint\" is required.",
					},
				},
				Action: func(c *cli.Context) error {
					if c.Args().Get(0) == "" {
						return cli.NewExitError(fmt.Sprintf("USAGE: %s %s", c.Command.HelpName, c.Command.ArgsUsage), 1)
					}
					configMap, err := updateConfigConvert(c.StringSlice("config"))
					if err != nil {
						return cli.NewExitError(err, 1)
					}
					err = cmd.ProviderUpdate(c.Args().Get(0), configMap)
					if err != nil {
						return cli.NewExitError(err, 1)
					}
					return nil
				},
			},
		},
	}
}

func updateConfigConvert(config []string) (map[string]interface{}, error) {
	configMap := map[string]interface{}{}
	if len(configMap) == 0 {
		return nil, errors.New("please input at least one config")
	}
	for _, v := range config {
		pair := strings.Split(v, "=")
		if len(pair) != 2 {
			pair = append(pair, "")
		}
		configMap[pair[0]] = pair[1]
	}
	return configMap, nil
}


func enrollConfigConvert(config []string) (map[string]interface{}, error) {
	configMap := map[string]interface{}{}
	if len(configMap) == 0 {
		return nil, errors.New("please input at least one config")
	}
	for _, v := range config {
		pair := strings.Split(v, "=")
		if len(pair) != 2 {
			return nil, errors.New("invalid config format")
		}
		if pair[1] == "" {
			return nil, errors.New("config value should not be empty")
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
	case "domains":
		err = parser.Domains(argv)
	case "services":
		err = parser.Service(argv)
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
