package api
import "fmt"


type StackParams struct {
	Name string `json:"name"`
	Content string `json:"content"`
}

type Image struct {
	Name string `json:"name"`
}

type Stack interface {
	Id() string
	Name() string
	Links() Links
	GetBuildImage() Image
}

type StackModel struct {
	IDField    string `json:"id"`
	NameField  string `json:"name"`
	LinksField []Link `json:"links"`
}

func (a StackModel) Id() string {
	return a.IDField
}

func (a StackModel) Name() string {
	return a.NameField
}

func (a StackModel) Links() Links {
	return LinksModel{
		Links: a.LinksField,
	}
}

func (a StackModel) GetBuildImage() Image {
	return Image{
		Name: fmt.Sprintf("hub.deepi.cn/%s-%s", a.NameField, "build"),
	}
}

type Stacks interface {
	Count() int
	First() Stacks
	Last() Stacks
	Prev() Stacks
	Next() Stacks
	Items() []Stack
}

type StacksModel struct {
	CountField int            `json:"count"`
	SelfField  string         `json:"self"`
	FirstField string         `json:"first"`
	LastField  string         `json:"last"`
	PrevField  string         `json:"prev"`
	NextField  string         `json:"next"`
	ItemsField []StackModel   `json:"items"`
}

func (stacks StacksModel) Count() int {
	return stacks.CountField
}

func (stacks StacksModel) Self() Stacks {
	return nil
}
func (stacks StacksModel) First() Stacks {
	return nil
}
func (stacks StacksModel) Last() Stacks {
	return nil
}
func (stacks StacksModel) Prev() Stacks {
	return nil
}
func (stacks StacksModel) Next() Stacks {
	return nil
}

func (stacks StacksModel) Items() []Stack {
	items := make([]Stack, 0)
	for _, stack := range stacks.ItemsField {
		items = append(items, stack)
	}
	return items
}
