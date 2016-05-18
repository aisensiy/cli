package config
import "encoding/json"

type Data struct {
	Email              string `json:"email"`
	Endpoint           string `json:"endpoint"`
	ApiEndpoint        string `json:"api_endpoint"`
	DeploymentEndpoint string `json:"deployment_endpoint"`
	GitHost            string `json:"host"`
	Auth               string `json:"auth"`
	Id                 string `json:"id"`
	Org                string `json:"org"`
}

func NewData() (data *Data) {
	data = new(Data)
	return
}

func (d *Data) JsonMarshalV3() (output []byte, err error) {
	return json.MarshalIndent(d, "", "  ")
}

func (d *Data) JsonUnmarshalV3(input []byte) (err error) {
	err = json.Unmarshal(input, d)
	if err != nil {
		return
	}

	return
}
