package api

import (
	"github.com/cnupp/runtimesdk/config"
	"github.com/cnupp/runtimesdk/net"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

type UpsRepository interface {
	GetUP(id string) (Up, error)
	GetUPByName(name string) (Ups, error)
	GetUpByUri(uri string) (Up, error)
	GetUps() (Ups, error)
	CreateUp(params map[string]interface{}) (Up, error)
	RemoveUp(id string) (error)
	UpdateUp(id string, params map[string]interface{}) (error)
	CreateProcedureInstance (upId string, procedureId string, params map[string]interface{}) (ProcedureInstance, error)
	GetProcedureInstanceByUri (uri string) (ProcedureInstance, error)
	GetProcedureInstance (id string) (ProcedureInstance, error)
	PublishUp(id string) (error)
	DeprecateUp(id string) (error)
}

type DefaultUpsRepository struct {
	config  config.Reader
	gateway net.Gateway
}

func NewUpsRepository(config config.Reader, gateway net.Gateway) UpsRepository {
	return DefaultUpsRepository{config, gateway}
}

func (upsRepo DefaultUpsRepository) GetUP(id string) (Up, error){
	upModel := UpModel{}
	err := upsRepo.gateway.Get(fmt.Sprintf("/ups/%s", id), &upModel)
	if err != nil {
		return nil, err
	}
	for index := range upModel.ProceduresField {
		upModel.ProceduresField[index].UpIdField = upModel.IdField
		upModel.ProceduresField[index].UpsMapper = NewUpsRepository(upsRepo.config, upsRepo.gateway)
	}


	return upModel, nil
}


func (upsRepo DefaultUpsRepository) GetUPByName(name string) (Ups, error){
	var body []byte

	res, err := upsRepo.gateway.Request("GET", fmt.Sprintf("/ups?name=%s", name), body)
	if err != nil {
		return nil, err
	}

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	upsModel := UpsModel{}
	err = json.Unmarshal(body, &upsModel)
	if err != nil {
		return nil, err
	}
	for n := range upsModel.ItemsField {
		for index := range upsModel.ItemsField[n].ProceduresField {
			upsModel.ItemsField[n].ProceduresField[index].UpIdField = upsModel.ItemsField[n].IdField
			upsModel.ItemsField[n].ProceduresField[index].UpsMapper = NewUpsRepository(upsRepo.config, upsRepo.gateway)
		}
	}
	return upsModel, nil
}

func (upsRepo DefaultUpsRepository) GetUpByUri(uri string) (Up, error){
	upModel := UpModel{}
	err := upsRepo.gateway.Get(uri, &upModel)
	if err != nil {
		return nil, err
	}
	for index := range upModel.ProceduresField {
		upModel.ProceduresField[index].UpIdField = upModel.IdField
		upModel.ProceduresField[index].UpsMapper = NewUpsRepository(upsRepo.config, upsRepo.gateway)
	}
	return upModel, nil
}

func (upsRepo DefaultUpsRepository) GetUps() (Ups, error) {
	upsModel := UpsModel{}
	err := upsRepo.gateway.Get(fmt.Sprintf("/ups"), &upsModel)
	if err != nil {
		return nil, err
	}
	for n := range upsModel.ItemsField {
		for index := range upsModel.ItemsField[n].ProceduresField {
			upsModel.ItemsField[n].ProceduresField[index].UpIdField = upsModel.ItemsField[n].IdField
			upsModel.ItemsField[n].ProceduresField[index].UpsMapper = NewUpsRepository(upsRepo.config, upsRepo.gateway)
		}
	}
	return upsModel, nil
}

func (upsRepo DefaultUpsRepository) CreateUp(params map[string]interface{}) (Up, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("Can not serilize the data")
	}

	res, err := upsRepo.gateway.Request("POST", "/ups", data)
	if err != nil {
		return nil, err
	}

	location := res.Header.Get("Location")
	var upModel UpModel
	err = upsRepo.gateway.Get(location, &upModel)
	if err != nil {
		return nil, err
	}
	for index := range upModel.ProceduresField {
		upModel.ProceduresField[index].UpIdField = upModel.IdField
		upModel.ProceduresField[index].UpsMapper = NewUpsRepository(upsRepo.config, upsRepo.gateway)
	}
	return upModel, nil
}

func (upsRepo DefaultUpsRepository) UpdateUp(id string, params map[string]interface{}) (error) {
	err := upsRepo.gateway.PUT(fmt.Sprintf("/ups/%s", id), params)
	if err != nil {
		return err
	}

	return nil

}

func (upsRepo DefaultUpsRepository) RemoveUp(id string) (error) {
	_, err := upsRepo.gateway.Request("DELETE", fmt.Sprintf("/ups/%s", id), nil)
	if err != nil {
		return err
	}
	return nil
}

func (upsRepo DefaultUpsRepository) CreateProcedureInstance(upId string, procedureId string, params map[string]interface{}) (ProcedureInstance, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("Can not serilize the data")
	}

	res, err := upsRepo.gateway.Request("POST", "/ups/" + upId + "/procedures/" + procedureId +"/instances", data)
	if err != nil {
		return nil, err
	}

	location := res.Header.Get("Location")
	var procedureInstanceModel ProcedureInstanceModel
	err = upsRepo.gateway.Get(location, &procedureInstanceModel)
	if err != nil {
		return nil, err
	}

	return procedureInstanceModel, nil
}

func (upsRepo DefaultUpsRepository) GetProcedureInstanceByUri (uri string) (ProcedureInstance, error) {
	var procedureInstance ProcedureInstanceModel
	err := upsRepo.gateway.Get(uri, &procedureInstance)
	if err != nil {
		return nil, err
	}
	return procedureInstance, nil
}

func (upsRepo DefaultUpsRepository) GetProcedureInstance (id string) (ProcedureInstance, error) {
	var procedureInstance ProcedureInstanceModel
	err := upsRepo.gateway.Get("/procedures/" + id, &procedureInstance)
	if err != nil {
		return nil, err
	}
	return procedureInstance, nil
}

func (upsRepo DefaultUpsRepository) PublishUp(id string) (error) {
	err := upsRepo.gateway.PUT(fmt.Sprintf("/ups/%s/publish", id), nil)
	if err != nil {
		return err
	}
	return nil
}


func (upsRepo DefaultUpsRepository) DeprecateUp(id string) (error) {
	err := upsRepo.gateway.PUT(fmt.Sprintf("/ups/%s/deprecate", id), nil)
	if err != nil {
		return err
	}
	return nil
}
