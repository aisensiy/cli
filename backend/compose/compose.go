package compose

import (
	"github.com/sjkyspa/stacks/controller/api/api"
	"gopkg.in/yaml.v2"
	"fmt"
	"github.com/sjkyspa/stacks/client/backend"
)

type ComposeBackend struct {

}
type Service struct {
	Image      string `json:"image" yaml:"image"`
	Entrypoint string `json:"entrypoint" yaml:"entrypoint,omitempty"`
	Command    string `json:"command" yaml:"command,omitempty"`
	Volumes    []string `json:"volumes" yaml:"volumes,omitempty"`
	Links      []string `json:"links" yaml:"links,omitempty"`
	Ports      []string `json:"ports" yaml:"ports,omitempty"`
}

type ComposeFile struct {
	Version  string `json:"version"`
	Services map[string]Service `json:"services"`
}

func (cb ComposeBackend) Up() {

}

func (cb ComposeBackend) ToComposeFile(s api.Stack) string {
	services := s.GetServices()
	composeServices := make(map[string]Service, 0)

	for name, service := range services {
		if !service.IsBuildable() {
			composeServices[name] = Service{
				Image: service.GetImage(),
				Links: service.GetLinks(),
			}
		} else {
			composeServices["runtime"] = Service{
				Image: service.GetBuild().Name,
				Entrypoint: "/bin/sh",
				Command: "-c 'tail -f /dev/null'",
			}
		}
	}
	composeFile := ComposeFile{
		Version: "2",
		Services: composeServices,
	}

	out, err := yaml.Marshal(composeFile)
	if err != nil {
		panic(fmt.Sprintf("Error happend when translate to yaml %v", err))
	}

	return string(out)
}

func NewComposeBackend() backend.Runtime {
	return ComposeBackend{}
}