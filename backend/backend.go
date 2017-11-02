package backend

import "github.com/cnupp/appssdk/api"

type Runtime interface {
	Up()
	ToComposeFile(api.Stack) string
}
