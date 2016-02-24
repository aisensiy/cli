package parser

import (
	"github.com/sjkyspa/stacks/Godeps/_workspace/src/github.com/docopt/docopt-go"
	"github.com/sjkyspa/stacks/client/cmd"
)

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
