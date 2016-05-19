package parser

import (
	"fmt"
	docopt "github.com/docopt/docopt-go"
	"github.com/sjkyspa/stacks/client/cmd"
	"regexp"
)

func Orgs(argv []string) error {
	usage := `
Valid commands for apps:

orgs:create             create a new organization
orgs:info               view info about an organization
orgs:current            set org as a default org
orgs:members            list members of the organization
orgs:add-member         add member to organization
orgs:rm-member          remove member from organization
orgs:apps               list apps of the organization
orgs:add-app            add app to organization
orgs:remove             destory an organization

Use 'cde help [command]' to learn more.
`
	switch argv[0] {
	case "orgs:create":
		return orgCreate(argv)
	case "orgs:info":
		return orgInfo(argv)
	case "orgs:current":
		return orgCurrent(argv)
	case "orgs:members":
		return orgMembers(argv)
	case "orgs:add-member":
		return orgAddMember(argv)
	case "orgs:rm-member":
		return orgRmMember(argv)
	case "orgs:apps":
		return orgApps(argv)
	case "orgs:add-app":
		return orgAddApp(argv)
	case "orgs:remove":
		return orgDestroy(argv)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		PrintUsage()
		return nil
	}
	return nil
}

func orgCreate(argv []string) error {
	usage := `
Creates a new organization.

Usage: cde orgs:create [options]

Options:
  -o --org=<name>
    tell system to deploy this app or not, 1 means need, 0 mean no, default 1
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	name := safeGetOrDefault(args, "--org", "1")

	regex := regexp.MustCompile(`^[a-z0-9\-]+$`)
	if !regex.MatchString(name) {
		return fmt.Errorf("'%s' does not match the pattern '[a-z0-9-]+'\n", name)
	}

	return cmd.OrgCreate(name)
}

func orgInfo(argv []string) error {
	usage := `
Prints info about the current organization.

Usage: cde orgs:info [options]

Options:
  -o --org=<org>
    the uniquely identifiable id for the organization.
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	orgName := safeGetValue(args, "--org")

	return cmd.GetOrg(orgName)
}

func orgCurrent(argv []string) error {
	usage := `
Set org as a default org.

Usage: cde orgs:current [options]

Options:
  -o --org=<org>
    the uniquely identifiable id for the organization.
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	orgName := safeGetValue(args, "--org")

	return cmd.SetCurrentOrg(orgName)
}

func orgMembers(argv []string) error {
	usage := `
List members of the organization

Usage: cde orgs:members [options]

Options:
  -o --org=<org>
    the uniquely identifiable id for the organization.
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	orgName := safeGetValue(args, "--org")

	return cmd.ListMembers(orgName)
}

func orgAddMember(argv []string) error {
	usage := `
Add member to organization

Usage: cde orgs:add-member [options]

Options:
  -o --org=<org>
    the uniquely identifiable id for the organization.
  -e --email=<email>
    the user email need to add to org
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	orgName := safeGetValue(args, "--org")
	email := safeGetValue(args, "--email")

	return cmd.AddMember(orgName, email)
}


func orgRmMember(argv []string) error {
	usage := `
Remove member from organization

Usage: cde orgs:rm-member [options]

Options:
  -o --org=<org>
    the uniquely identifiable id for the organization.
  -e --email=<email>
    the user email need to add to org
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	orgName := safeGetValue(args, "--org")
	email := safeGetValue(args, "--email")

	return cmd.RemoveMember(orgName, email)
}

func orgApps(argv []string) error {
	usage := `
List apps of the organization

Usage: cde orgs:apps [options]

Options:
  -o --org=<org>
    the uniquely identifiable id for the organization.
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	orgName := safeGetValue(args, "--org")

	return cmd.ListApps(orgName)
}

func orgAddApp(argv []string) error {
	usage := `
Add app to organization

Usage: cde orgs:add-app [options]

Options:
  -o --org=<org>
    the uniquely identifiable id for the organization.
  -a --app=<app>
    the app name.
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	orgName := safeGetValue(args, "--org")
	appName := safeGetValue(args, "--app")

	return cmd.AddOrgApp(orgName, appName)
}

func orgDestroy(argv []string) error {
	usage := `
Destroy organization

Usage: cde orgs:remove [options]

Options:
  -o --org=<org>
    the uniquely identifiable id for the organization.
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	orgName := safeGetValue(args, "--org")

	return cmd.DestroyOrg(orgName)
}
