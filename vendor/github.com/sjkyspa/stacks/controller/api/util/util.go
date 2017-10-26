package util

import (
	"strings"
)

func IDFromURI(uri string) (string, error) {
	parts := strings.Split(uri, "/")
	return parts[len(parts)-1], nil
}
