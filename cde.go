package main

import (
	"fmt"
	docopt "github.com/docopt/docopt-go"
	"github.com/fatih/color"
	"github.com/sjkyspa/stacks/client/parser"
	"github.com/sjkyspa/stacks/version"
	"os"
	"strings"
)

func main() {
	os.Exit(Command(os.Args[1:]))
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
`
	command, argv := parseArgs(argv)

	_, err := docopt.Parse(usage, []string{command}, false, version.Version, true, false)
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
	if len(argv) == 1 {
		// rearrange "cde --help" as "cde help"
		if argv[0] == "--help" || argv[0] == "-h" {
			argv[0] = "help"
		}
	}

	if len(argv) >= 2 {
		// Rearrange "cde help <command>" to "cde <command> --help".
		if argv[0] == "help" || argv[0] == "--help" || argv[0] == "-h" {
			argv = append(argv[1:], "--help")
		}
	}

	if len(argv) > 0 {
		argv[0] = replaceShortcut(argv[0])
		index := strings.Index(argv[0], ":")

		if index != -1 {
			command := argv[0]
			return command[:index], argv
		}

		return argv[0], argv
	}

	return "", argv
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

	expandedCommand := shortcuts[command]
	if expandedCommand == "" {
		return command
	}

	return expandedCommand
}
