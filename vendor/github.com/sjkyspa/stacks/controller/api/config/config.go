package config

type ConfigRepository interface {
	Reader
	Writer
	Close()
}

type Reader interface {
	ApiEndpoint() string
	Email() string
	Auth() string
	Id() string
}

type Writer interface {
	SetApiEndpoint(string)
	SetEmail(string)
	SetAuth(string)
	SetId(string)
}
