package config

import (
	"encoding/json"
	"os"
	"errors"
	"path"
	"net/url"
	"io/ioutil"
)

type Persistor struct {
	Email    string `json:"email"`
	Endpoint string `json:"endpoint"`
	Auth     string `json:"auth"`
}

func NewPersistor() (Persistor, error) {
	filename := locateSettingsFile()

	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return Persistor{}, errors.New("Not logged in. Use 'cde login' or 'cde register' to get started.")
		}

		return Persistor{}, err
	}

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return Persistor{}, err
	}

	settings := Persistor{}
	if err = json.Unmarshal(contents, &settings); err != nil {
		return Persistor{}, err
	}

	_, err = url.Parse(settings.Endpoint)
	if err != nil {
		return Persistor{}, err
	}

	return settings, nil

}

func(p Persistor) Save() error {
	settingsContents, err := json.Marshal(p)

	if err != nil {
		return err
	}

	if err = os.MkdirAll(path.Join(FindHome(), "/.cde/"), 0775); err != nil {
		return err
	}

	return ioutil.WriteFile(locateSettingsFile(), settingsContents, 0775)
}

func SavePersistor(p Persistor) error {
	settingsContents, err := json.Marshal(p)

	if err != nil {
		return err
	}

	if err = os.MkdirAll(path.Join(FindHome(), "/.cde/"), 0775); err != nil {
		return err
	}

	return ioutil.WriteFile(locateSettingsFile(), settingsContents, 0775)
}

func FindHome() string {
	return os.Getenv("HOME")
}

func locateSettingsFile() string {
	filename := os.Getenv("CDE_PROFILE")

	if filename == "" {
		filename = "client"
	}

	return path.Join(FindHome(), ".cde", filename + ".json")
}

