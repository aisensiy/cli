package api

import "fmt"

//go:generate counterfeiter -o fakes/fake_build.go . Build
type Build interface {
	Id() string
	GitSha() string
	Status() string
	Verify() Verify
	Links() Links
	GetApp() App
	Success() error
	Fail() error
	IsSuccess() bool
	IsFail() bool
	VerifySuccess() error
	IsVerifySuccess() bool
	VerifyFail() error
	IsVerifyFail() bool
}

type BuildParams struct {
	GitSha string `json:"git_sha"`
	User   string `json:"user"`
	Source string `json:"source"`
}

type Verify struct {
	StatusField string `json:"status"`
}

type BuildModel struct {
	GitShaField string      `json:"git_sha"`
	IDField     string      `json:"id"`
	StatusField string      `json:"status"`
	VerifyField Verify      `json:"verify"`
	LinksField  []Link      `json:"links"`
	AppField    App         `json:"-"`
	BuildMapper BuildMapper `json:"-"`
	Resource    Resource    `json:"-"`
}

type BuildRef struct {
	GitSha string `json:"git_sha"`
	ID     string `json:"id"`
	Status string `json:"status"`
	Links  []Link `json:"links"`
}

type Links interface {
	Self() Link
	Link(id string) (Link, error)
}

type Link struct {
	Relation string `json:"rel"`
	URI      string `json:"uri"`
}

type LinksModel struct {
	Links []Link
}

func (lm LinksModel) filter(array []Link, filterBy func(l Link) bool) []Link {
	results := make([]Link, 0)
	for _, link := range array {
		if filterBy(link) {
			results = append(results, link)
		}
	}
	return results
}

func (lm LinksModel) Self() (link Link) {
	results := lm.filter(lm.Links, func(l Link) bool {
		return l.Relation == "self"
	})

	if len(results) == 1 {
		link = results[0]
		return
	}
	return
}

func (lm LinksModel) Link(rel string) (link Link, apiErr error) {
	results := lm.filter(lm.Links, func(l Link) bool {
		return l.Relation == rel
	})
	if len(results) == 1 {
		link = results[0]
		return
	}
	apiErr = fmt.Errorf("Found %d Link in links", len(results))
	return
}

func (bm BuildModel) Id() string {
	return bm.IDField
}

func (bm BuildModel) GitSha() string {
	return bm.GitShaField
}

func (bm BuildModel) Status() string {
	return bm.StatusField
}

func (bm BuildModel) Verify() Verify {
	return bm.VerifyField
}

func (bm BuildModel) Links() Links {
	return LinksModel{
		Links: bm.LinksField,
	}
}

func (bm BuildModel) Success() (error) {
	return bm.BuildMapper.Success(bm)
}

func (bm BuildModel) GetApp() App {
	link, err := bm.Links().Link("app")
	if err != nil {
		return AppModel{}
	}
	model, err := bm.Resource.GetResourceByURI(link.URI)
	if err != nil {
		return AppModel{}
	}
	return model.(App)
}

func (bm BuildModel) Fail() error {
	return bm.BuildMapper.Fail(bm)
}

func (bm BuildModel) IsFail() bool {
	return bm.StatusField == "FAIL"
}

func (bm BuildModel) IsSuccess() bool {
	return bm.StatusField == "SUCCESS"
}

func (bm BuildModel) VerifySuccess() error {
	return bm.BuildMapper.VerifySuccess(bm)
}

func (bm BuildModel) IsVerifySuccess() bool {
	return bm.VerifyField.StatusField == "SUCCESS"
}

func (bm BuildModel) VerifyFail() error {
	return bm.BuildMapper.VerifyFail(bm)
}

func (bm BuildModel) IsVerifyFail() bool {
	return bm.VerifyField.StatusField == "FAIL"
}

type Builds interface {
	Count() int
	Self() Builds
	First() Builds
	Last() Builds
	Prev() Builds
	Next() Builds
	Items() []BuildRef
}

type BuildsModel struct {
	CountField  int        `json:"count"`
	SelfField   string     `json:"self"`
	FirstField  string     `json:"first"`
	LastField   string     `json:"last"`
	PrevField   string     `json:"prev"`
	NextField   string     `json:"next"`
	ItemsField  []BuildRef `json:"items"`
	BuildMapper BuildMapper
}

func (bsm BuildsModel) Count() int {
	return bsm.CountField
}
func (bsm BuildsModel) Self() Builds {
	return nil
}
func (bsm BuildsModel) First() Builds {
	return nil
}
func (bsm BuildsModel) Last() Builds {
	return nil
}
func (bsm BuildsModel) Prev() Builds {
	return nil
}
func (bsm BuildsModel) Next() Builds {
	return nil
}

func (bsm BuildsModel) Items() []BuildRef {
	return bsm.ItemsField
}
