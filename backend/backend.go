package backend

import "github.com/cnupp/cnup/controller/api/api"

type Runtime interface {
	Up()
	ToComposeFile(api.Stack) string
}
