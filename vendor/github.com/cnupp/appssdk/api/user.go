package api

type UserParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type User interface {
	Id() string
	Email() string
	UploadKey(keyParam KeyParams) (Key, error)
	DeleteKey(keyId string) error
	Keys() (Keys, error)
	Links() Links
}

type UserModel struct {
	IdField    string `json:"id"`
	EmailField string `json:"email"`
	LinksArray []Link `json:"links"`
	KeyRepo    KeyRepository
}

func (u UserModel) Id() string {
	return u.IdField
}

func (u UserModel) Email() string {
	return u.EmailField
}

func (u UserModel) Links() Links {
	return LinksModel{
		Links: u.LinksArray,
	}
}

func (u UserModel) UploadKey(keyParam KeyParams) (Key, error) {
	return u.KeyRepo.Upload(u, keyParam)
}

func (u UserModel) DeleteKey(keyId string) error {
	return u.KeyRepo.Delete(u, keyId)
}

func (u UserModel) Keys() (Keys, error) {
	return u.KeyRepo.GetKeysForUser(u)
}

type Users interface {
	Count() int
	First() Users
	Last() Users
	Prev() Users
	Next() Users
	Items() []User
}

type UsersModel struct {
	CountField int         `json:"count"`
	SelfField  string      `json:"self"`
	FirstField string      `json:"first"`
	LastField  string      `json:"last"`
	PrevField  string      `json:"prev"`
	NextField  string      `json:"next"`
	ItemsField []UserModel `json:"items"`
}

func (users UsersModel) Count() int {
	return users.CountField
}

func (users UsersModel) Self() Users {
	return nil
}
func (users UsersModel) First() Users {
	return nil
}
func (users UsersModel) Last() Users {
	return nil
}
func (users UsersModel) Prev() Users {
	return nil
}
func (users UsersModel) Next() Users {
	return nil
}

func (users UsersModel) Items() []User {
	items := make([]User, 0)
	for _, user := range users.ItemsField {
		items = append(items, user)
	}
	return items
}
