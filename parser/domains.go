package parser

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/sjkyspa/stacks/client/cmd"
)

// Domains routes domain commands to their specific function.
func Domains(argv []string) error {
	usage := `
Valid commands for domains:

domains:create        create a domain
domains:list          list domains
domains:remove        remove a domain
domains:cert          attach cert to domain

Use 'cde help [command]' to learn more.
`
	switch argv[0] {
	case "domains:create":
		return domainsAdd(argv)
	case "domains:list":
		return domainsList(argv)
	case "domains:remove":
		return domainsRemove(argv)
	case "domains:cert":
		return domainCert(argv)
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

Usage: cde domains:create <domain>

Arguments:
  <domain>
    the domain name to be bound to the application, such as 'domain.testapp.com'.
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

Usage: cde domains:list
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

Usage: cde domains:remove <domain>

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

func domainCert(argv []string) error {
	usage := `
Attach cert to the domain

Usage: cde domains:cert <domain> <crt> <private-key>

Arguments:
  <domain>
    the domain name to be removed.
  <crt>
    the certificate file of the domain.
  <private-key>
    the private key file of the domain crt
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmd.DomainsCert(safeGetValue(args, "<domain>"), safeGetValue(args, "<crt>"), safeGetValue(args, "<private-key>"),)
}