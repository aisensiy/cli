package backend

import "github.com/sjkyspa/stacks/controller/api/api"

type Runtime interface {
	Up()
	ToComposeFile(api.Stack) string
}