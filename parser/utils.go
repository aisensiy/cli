package parser
import (
	"fmt"
	"os"
)

// PrintUsage runs if no matching command is found.
func PrintUsage() {
	fmt.Fprintln(os.Stderr, "Found no matching command, try 'cde help'")
	fmt.Fprintln(os.Stderr, "Usage: cde <command> [<args>...]")
}

func printHelp(argv []string, usage string) bool {
	if len(argv) > 1 {
		if argv[1] == "--help" || argv[1] == "-h" {
			fmt.Print(usage)
			return true
		}
	}

	return false
}

func safeGetValue(args map[string]interface{}, key string) string {
	return safeGetOrDefault(args, key, "")
}

func safeGetOrDefault(args map[string]interface{}, key string, defaultVal string) string {
	if args[key] == nil {
		return defaultVal
	}
	return args[key].(string)
}