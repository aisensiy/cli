package parser

import (
	"github.com/docopt/docopt-go"
	"github.com/sjkyspa/stacks/client/cmd"
	"github.com/urfave/cli"
	"fmt"
)

func RoutesCommand() cli.Command{
	return cli.Command {
		Name: "routes",
		Usage: "Routes Commands",
		Subcommands: []cli.Command {
			{
				Name: "create",
				Usage: "Create a new routes",
				ArgsUsage: "<domain> <path>",
				Action: func(c *cli.Context) error {
					if c.Args().Get(0) == "" || c.Args().Get(1) == "" {
						return cli.NewExitError(fmt.Sprintf("USAGE: %s %s", c.Command.HelpName, c.Command.ArgsUsage), 1)
					}
					if err := cmd.RoutesCreate(c.Args().Get(0), c.Args().Get(1)); err != nil {
						return cli.NewExitError(err, 1)
					}
					return nil
				},
			},
			{
				Name: "list",
				Usage: "List accessible routes",
				ArgsUsage: " ",
				Action: func(c *cli.Context) error {
					if err := cmd.RoutesList(); err != nil {
						return cli.NewExitError(err, 1)
					}
					return nil
				},
			},
			{
				Name: "bind",
				Usage: "Bind a route with an app",
				ArgsUsage: "<route> <app>",
				Action: func(c *cli.Context) error {
					if c.Args().Get(0) == "" || c.Args().Get(1) == "" {
						return cli.NewExitError(fmt.Sprintf("USAGE: %s %s", c.Command.HelpName, c.Command.ArgsUsage), 1)
					}
					if err := cmd.RouteBindWithApp(c.Args().Get(0), c.Args().Get(1)); err != nil {
						return cli.NewExitError(err, 1)
					}
					return nil
				},
			},
			{
				Name: "unbind",
				Usage: "Unbind a route with an app",
				ArgsUsage: "<route> <app>",
				Action: func(c *cli.Context) error {
					if c.Args().Get(0) == "" || c.Args().Get(1) == "" {
						return cli.NewExitError(fmt.Sprintf("USAGE: %s %s", c.Command.HelpName, c.Command.ArgsUsage), 1)
					}
					if err := cmd.UnbindRouteWithApp(c.Args().Get(0), c.Args().Get(1)); err != nil {
						return cli.NewExitError(err, 1)
					}
					return nil
				},
			},
		},
	}
}

func Routes(argv []string) error {
	usage := `
Valid commands for routes:

routes:create        create a new routes
routes:list          list accessible routes
routes:bind          bind a route with an app
routes:unbind        unbind a route with an app
Use 'cde help [command]' to learn more.
`
	switch argv[0] {
	case "routes:create":
		return routeCreate(argv)
	case "routes:list":
		return routesList()
	case "routes:bind":
		return bindRouteWithApp(argv)
	case "routes:unbind":
		return unbindRouteWithApp(argv)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "routes" {
			argv[0] = "routes:list"
			return routesList()
		}

		PrintUsage()
		return nil
	}
	return nil
}

func routesList() error {
	return cmd.RoutesList()
}

func routeCreate(argv []string) error {
	usage := `
Creates a new route.

Usage: cde routes:create <domain> <path>

Arguments:
  <domain>
    the domain name.
  <path>
  	the route path
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	domain := safeGetValue(args, "<domain>")
	path := safeGetValue(args, "<path>")
	return cmd.RoutesCreate(domain, path)
}

func bindRouteWithApp(argv []string) error {
	usage := `
bind a route with an app.

Usage: cde routes:bind <route> <app>

Arguments:
  <route>
    the domain/path
  <app>
  	the app name
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	route := safeGetValue(args, "<route>")
	app := safeGetValue(args, "<app>")
	return cmd.RouteBindWithApp(route, app)
}

func unbindRouteWithApp(argv []string) error {
	usage := `
unbind a route with an app.

Usage: cde routes:unbind <route> <app>

Arguments:
  <route>
    the domain/path
  <app>
  	the app name
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	route := safeGetValue(args, "<route>")
	app := safeGetValue(args, "<app>")
	return cmd.UnbindRouteWithApp(route, app)
}
