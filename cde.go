package main

import (
	"fmt"
	docopt "github.com/docopt/docopt-go"
	"github.com/fatih/color"
	"github.com/sjkyspa/stacks/client/parser"
	"github.com/sjkyspa/stacks/version"
	"os"
	"strings"
	"github.com/sjkyspa/stacks/client/util"
	"github.com/urfave/cli"
	"github.com/sjkyspa/stacks/client/cmd"
	"time"
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
		{
			Name:    "ups",
			Usage:   "Unified Procedures Commands",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "list all Unified Procedures",
					ArgsUsage: " ",
					Action: func(c *cli.Context) error {
						return cmd.UpsList()
					},
				},
				{
					Name:  "info",
					Usage: "get info of an Unified Procedure",
					ArgsUsage: "<up-name>",
					Action: func(c *cli.Context) error {
						return cmd.UpsInfo(c.Args().First())
					},
				},
				{
					Name:  "draft",
					Usage: "create a new Unified Procedure",
					ArgsUsage: "<up-file>",
					Action: func(c *cli.Context) error {
						return cmd.UpCreate(c.Args().First())
					},
				},
			},
		},
	}

	if len(os.Args) > 1 && os.Args[1] == "ups" {
		app.Run(os.Args)
	} else {
		os.Exit(Command(os.Args[1:]))
	}
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

	if len(argv) == 1 && (argv[0] == "-h" || argv[0] == "--help" ) {
		// rearrange "cde --help" as "cde help"
		argv[0] = "help"
	}

	if len(argv) >= 2 && (argv[0] == "help" || argv[0] == "-h" || argv[0] == "--help" ) {
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
