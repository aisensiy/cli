package api

type Auth interface {
	Id() string
	UserEmail() string
}

type AuthModel struct {
	IdField string
	UserEmailField string
}

func (a AuthModel) Id() string {
	return a.IdField
}

func (a AuthModel) UserEmail() string {
	return a.UserEmailField
}
