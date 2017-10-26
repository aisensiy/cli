package version

import "os"

func Version() string {
	cdeVersion := os.Getenv("CDE_VERSION")
	if cdeVersion == "" {
		return "0.1.4"
	} else {
		return cdeVersion
	}

}
