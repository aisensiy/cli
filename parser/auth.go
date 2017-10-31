package parser

import (
	"fmt"

	docopt "github.com/docopt/docopt-go"
	"github.com/cnupp/cli/cmd"
	cli "gopkg.in/urfave/cli.v2"
)

func AuthCommand() *cli.Command {
	return &cli.Command{
		Name:  "auth",
		Usage: "Auth Commands",
		Subcommands: []*cli.Command{
			{
				Name:      "register",
				Usage:     "Register a new user on a specific controller",
				ArgsUsage: "<controller>",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "email, e",
						Usage: "Provide email for the new user",
					},
					&cli.StringFlag{
						Name:  "password, p",
						Usage: "Provide password for the new user",
					},
				},
				Action: func(c *cli.Context) error {
					if c.Args().Get(0) == "" {
						return cli.Exit(fmt.Sprintf("USAGE: %s %s", c.Command.HelpName, c.Command.ArgsUsage), 1)
					}
					err := cmd.Register(c.Args().First(), c.String("email"), c.String("password"))
					if err != nil {
						return cli.Exit(fmt.Sprintf("%v", err), 1)
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
				ArgsUsage: "<controller>",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "email, e",
						Usage: "Provide user email",
					},
					&cli.StringFlag{
						Name:  "password, p",
						Usage: "Provide user password",
					},
					&cli.StringFlag{
						Name:  "ssl-verify, s",
						Usage: "Disable SSL certificate verification for API requests",
					},
				},
				Action: func(c *cli.Context) error {
					if c.Args().Get(0) == "" {
						return cli.Exit(fmt.Sprintf("USAGE: %s %s", c.Command.HelpName, c.Command.ArgsUsage), 1)
					}
					err := cmd.Login(c.Args().First(), c.String("email"), c.String("password"))
					if err != nil {
						return cli.Exit(fmt.Sprintf("%v", err), 1)
					}
					return nil
				},
			},
			{
				Name:      "logout",
				Usage:     "Log out from a controoler",
				ArgsUsage: " ",
				Action: func(c *cli.Context) error {
					err := cmd.Logout()
					if err != nil {
						return cli.Exit(fmt.Sprintf("%v", err), 1)
					}
					return nil
				},
			},
		},
	}
}

// Auth routes auth commands to the specific function.
func Auth(argv []string) error {
	usage := `
Valid commands for auth:

auth:register          register a new user
auth:login             authenticate against a controller
auth:logout            clear the current user session
auth:whoami            display the current user

Use 'cde help [command]' to learn more.
`

	switch argv[0] {
	case "auth:register":
		return authRegister(argv)
	case "auth:login":
		return authLogin(argv)
	case "auth:logout":
		return authLogout(argv)
		//	case "auth:passwd":
		//		return authPasswd(argv)
	case "auth:whoami":
		return authWhoami(argv)
		//	case "auth:cancel":
		//		return authCancel(argv)
		//	case "auth:regenerate":
		//		return authRegenerate(argv)
	case "auth":
		fmt.Print(usage)
		return nil
	default:
		PrintUsage()
		return nil
	}
}

func authRegister(argv []string) error {
	usage := `
Registers a new user with a Cde controller.

Usage: cde auth:register <controller> [options]

Arguments:
  <controller>
    fully-qualified controller URI, e.g. 'http://cde.com/'

Options:
  --email=<email>
    provide an email address.
  --password=<password>
    provide a password for the new account.
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	controller := safeGetValue(args, "<controller>")
	password := safeGetValue(args, "--password")
	email := safeGetValue(args, "--email")

	return cmd.Register(controller, email, password)
}

func authLogin(argv []string) error {
	usage := `
Logs in by authenticating against a controller.

Usage: cde auth:login <controller> [options]

Arguments:
  <controller>
    a fully-qualified controller URI, e.g. "http://cde.local3.cdeapp.com/".

Options:
  --email=<email>
    provide a email for the account.
  --password=<password>
    provide a password for the account.
  --ssl-verify=false
    disables SSL certificate verification for API requests
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	controller := safeGetValue(args, "<controller>")
	email := safeGetValue(args, "--email")
	password := safeGetValue(args, "--password")

	return cmd.Login(controller, email, password)
}

func authLogout(argv []string) error {
	usage := `
Logs out from a controller and clears the user session.

Usage: cde auth:logout
`

	if _, err := docopt.Parse(usage, argv, true, "", false, true); err != nil {
		return err
	}

	return cmd.Logout()
}

//func authPasswd(argv []string) error {
//	usage := `
//Changes the password for the current user.
//
//Usage: cde auth:passwd [options]
//
//Options:
//  --password=<password>
//    the current password for the account.
//  --new-password=<new-password>
//    the new password for the account.
//  --username=<username>
//    the account's username.
//`
//
//	args, err := docopt.Parse(usage, argv, true, "", false, true)
//
//	if err != nil {
//		return err
//	}
//
//	username := safeGetValue(args, "--username")
//	password := safeGetValue(args, "--password")
//	newPassword := safeGetValue(args, "--new-password")
//
//	return cmd.Passwd(username, password, newPassword)
//}

func authWhoami(argv []string) error {
	usage := `
Displays the currently logged in user.

Usage: cde auth:whoami
`

	if _, err := docopt.Parse(usage, argv, true, "", false, true); err != nil {
		return err
	}

	return cmd.Whoami()
}

func authCancel(argv []string) error {
	usage := `
Cancels and removes the current account.

Usage: cde auth:cancel [options]

Options:
  --username=<username>
    provide a username for the account.
  --password=<password>
    provide a password for the account.
  --yes
    force "yes" when prompted.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	username := safeGetValue(args, "--username")
	password := safeGetValue(args, "--password")
	yes := args["--yes"].(bool)

	return cmd.Cancel(username, password, yes)
}

func authRegenerate(argv []string) error {
	usage := `
Regenerates auth token, defaults to regenerating token for the current user.

Usage: cde auth:regenerate
`

	_, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmd.Regenerate()
}
