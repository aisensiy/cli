package parser

import (
	"github.com/docopt/docopt-go"
	"github.com/sjkyspa/stacks/client/cmd"
)

func Launch(argv []string) error {
	usage := `
Valid commands for launch:

launch:build  launch a build procedure

Use 'cde help [command]' to learn more.
`
	switch argv[0] {
	case "launch:build":
		return launchBuild(argv)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		PrintUsage()
		return nil
	}
	return nil
}

func launchBuild(argv []string) error {
	usage := `
Launch a build procedure.

Usage: cde launch:build (-f <filename>)

Arguments:
  <filename>
  the code base to build with
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}
	filename := safeGetValue(args, "<filename>")

	return cmd.LaunchBuild(filename)
}
