package config

type ConfigRepository interface {
	Reader
	Writer
	Close()
}

type Reader interface {
	DeploymentEndpoint() string
	Auth() string
}

type Writer interface {
	SetDeploymentEndpoint(string)
}
