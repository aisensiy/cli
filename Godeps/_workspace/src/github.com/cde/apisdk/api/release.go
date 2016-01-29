package api

type Release interface {
	Id() string
	Version() int
	Envs() map[string]string
	Status() string
	Links() Links
	GetApp() App
	Success() (error)
	Fail() (error)
	IsSuccess() (bool)
	IsFail() (bool)
}

type ReleaseParams struct {
	BuildId string `json:"build_id"`
}

type ReleaseModel struct {
	IDField       string    `json:"id"`
	VersionField  int    `json:"version"`
	EnvsField     map[string]string    `json:"envs"`
	StatusField   string    `json:"status"`
	AppField      App          `json:"-"`
	LinksField    []Link     `json:"links"`
	ReleaseMapper ReleaseMapper `json:"-"`
}

func (bm ReleaseModel) Id() string {
	return bm.IDField
}

func (bm ReleaseModel) Status() string {
	return bm.StatusField
}

func (bm ReleaseModel) Envs() map[string]string {
	return bm.EnvsField
}

func (bm ReleaseModel) Version() int {
	return bm.VersionField
}

func (bm ReleaseModel) Links() Links {
	return LinksModel{
		Links: bm.LinksField,
	}
}

func (bm ReleaseModel) Success() (apiErr error) {
	return bm.ReleaseMapper.Success(bm)
}

func (bm ReleaseModel) GetApp() App {
	return bm.AppField.(App)
}

func (bm ReleaseModel) Fail() (error) {
	return bm.ReleaseMapper.Fail(bm)
}

func (bm ReleaseModel) IsFail() (bool) {
	return bm.StatusField == "FAIL"
}

func (bm ReleaseModel) IsSuccess() (bool) {
	return bm.StatusField == "SUCCESS"
}

type Releases interface {
	Count() int
	Self() Releases
	First() Releases
	Last() Releases
	Prev() Releases
	Next() Releases
	Items() []ReleaseModel
}

type ReleasesModel struct {
	CountField    int            `json:"count"`
	SelfField     string         `json:"self"`
	FirstField    string         `json:"first"`
	LastField     string         `json:"last"`
	PrevField     string         `json:"prev"`
	NextField     string         `json:"next"`
	ItemsField    []ReleaseModel  `json:"items"`
	ReleaseMapper ReleaseMapper
}

func (bsm ReleasesModel) Count() int {
	return bsm.CountField
}
func (bsm ReleasesModel) Self() Releases {
	return nil
}
func (bsm ReleasesModel) First() Releases {
	return nil
}
func (bsm ReleasesModel) Last() Releases {
	return nil
}
func (bsm ReleasesModel) Prev() Releases {
	return nil
}
func (bsm ReleasesModel) Next() Releases {
	return nil
}

func (bsm ReleasesModel) Items() []ReleaseModel {
	return bsm.ItemsField
}
