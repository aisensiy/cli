package parser

import (
	"github.com/cde/client/cmd"
	docopt "github.com/docopt/docopt-go"
)

// Domains routes domain commands to their specific function.
func Domains(argv []string) error {
	usage := `
Valid commands for domains:

domains:add           create a domain
domains:list          list domains
domains:remove        remove a domain

Use 'deis help [command]' to learn more.
`
	switch argv[0] {
	case "domains:add":
		return domainsAdd(argv)
	case "domains:list":
		return domainsList(argv)
	case "domains:remove":
		return domainsRemove(argv)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "domains" {
			argv[0] = "domains:list"
			return domainsList(argv)
		}

		PrintUsage()
		return nil
	}
}

func domainsAdd(argv []string) error {
	usage := `
Binds a domain to an application.

Usage: deis domains:add <domain>

Arguments:
  <domain>
    the domain name to be bound to the application, such as 'domain.deisapp.com'.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmd.DomainsAdd(safeGetValue(args, "<domain>"))
}

func domainsList(argv []string) error {
	usage := `
Lists domains bound to an application.

Usage: deis domains:list
`

	_, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}


	return cmd.DomainsList()
}

func domainsRemove(argv []string) error {
	usage := `
Unbinds a domain for an application.

Usage: deis domains:remove <domain>

Arguments:
  <domain>
    the domain name to be removed.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmd.DomainsRemove(safeGetValue(args, "<domain>"))
}
