package api

import "errors"

type TemplateModel struct {
	Type string
	URI  string
}

type Up interface {
	Id() string
	Template() TemplateModel
	Name() string
	Status() string
	Procedures() []ProcedureModel
	GetProcedureByType(typeName string) (ProcedureModel, error)
}

type UpModel struct {
	IdField         string           `json:"id"`
	NameField       string           `json:"name"`
	StatusField       string           `json:"status"`
	TemplateFiled   TemplateModel         `json:"template"`
	ProceduresField []ProcedureModel `json:"procedures"`
}

func (u UpModel) Id() string {
	return u.IdField
}

func (u UpModel) Name() string {
	return u.NameField
}

func (u UpModel) Status() string {
	return u.StatusField
}

func (u UpModel) Template() TemplateModel {
	return u.TemplateFiled
}

func (u UpModel) Procedures() []ProcedureModel {
	return u.ProceduresField
}

func (u UpModel) GetProcedureByType(typeName string) (ProcedureModel, error) {
	for _, procedure := range u.ProceduresField {
		if procedure.Type() == typeName {
			return procedure, nil
		}
	}
	return ProcedureModel{}, errors.New("No procedure found")
}

type Procedure interface {
	Id() string
	Type() string
	Links() []LinkModel
	App() ProcedureAppModel
	Runtime() ProcedureRuntimeModel
}

type ProcedureModel struct {
	IdField      string                `json:"id"`
	TypeField    string                `json:"type"`
	LinksField   []LinkModel           `json:"links"`
	AppField     ProcedureAppModel     `json:"app"`
	RuntimeField ProcedureRuntimeModel `json:"runtime"`
}

func (p ProcedureModel) Id() string {
	return p.IdField
}

func (p ProcedureModel) Type() string {
	return p.TypeField
}

func (p ProcedureModel) Links() []LinkModel {
	return p.LinksField
}

func (p ProcedureModel) App() ProcedureAppModel {
	return p.AppField
}

func (p ProcedureModel) Runtime() ProcedureRuntimeModel {
	return p.RuntimeField
}

type ProcedureAppModel struct {
	NameField      string            `json:"name"`
	CpuField       float64           `json:"cpu"`
	MemField       float64           `json:"mem"`
	DiskField      float64           `json:"disk"`
	InstancesField int               `json:"instances"`
	ImageField     string            `json:"image"`
	EnvField       map[string]string `json:"environment"`
	ExposesField   []int             `json:"exposes"`
	VolumesField   []VolumeModel     `json:"volumes"`
	LinksField     []string          `json:"links"`
	HealthsField   []HealthModel     `json:"healths"`
}

type VolumeModel struct {
	ModeField   string `json:"mode"`
	SourceField string `json:"source"`
	TargetField string `json:"target"`
	ScopeField  string `json:"scope"`
}

type HealthModel struct {
	ProtocolField    string `json:"protocol"`
	PortField        int    `json:"port"`
	MappedField      int    `json:"mapped"`
	IgnoreField      int    `json:"ignore"`
	IntervalField    int    `json:"interval"`
	TimeoutField     int    `json:"timeout"`
	ConsecutiveField int    `json:"consecutive"`
}

type LinkModel struct {
	RelField string `json:"rel"`
	UriField string `json:"uri"`
}

type ProcedureRuntimeModel struct {
	IdField        string              `json:"id"`
	LinksField     []LinkModel         `json:"links"`
	ProcessesField []ProcedureAppModel `json:"processes"`
}

type Ups interface {
	First() string
	Last() string
	Prev() string
	Next() string
	Count() int
	Self() string
	Items() []UpModel
}

type UpsModel struct {
	FirstField string    `json:"first"`
	LastField  string    `json:"last"`
	PrevField  string    `json:"prev"`
	NextField  string    `json:"next"`
	CountField int       `json:"count"`
	SelfField  string    `json:"self"`
	ItemsField []UpModel `json:"items"`
}

func (u UpsModel) First() string {
	return u.FirstField
}

func (u UpsModel) Last() string {
	return u.LastField
}

func (u UpsModel) Prev() string {
	return u.PrevField
}

func (u UpsModel) Next() string {
	return u.NextField
}

func (u UpsModel) Count() int {
	return u.CountField
}

func (u UpsModel) Self() string {
	return u.SelfField
}

func (u UpsModel) Items() []UpModel {
	return u.ItemsField
}
