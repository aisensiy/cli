package compose

import (
	"github.com/sjkyspa/stacks/controller/api/api"
	"gopkg.in/yaml.v2"
	"fmt"
	"github.com/sjkyspa/stacks/client/backend"
	"strings"
	"path/filepath"
	"os"
)

type ComposeBackend struct {

}
type Service struct {
	Image       string `json:"image" yaml:"image"`
	Entrypoint  string `json:"entrypoint" yaml:"entrypoint,omitempty"`
	Command     string `json:"command" yaml:"command,omitempty"`
	Volumes     []string `json:"volumes" yaml:"volumes,omitempty"`
	Links       []string `json:"links" yaml:"links,omitempty"`
	Ports       []string `json:"ports" yaml:"ports,omitempty"`
	Environment map[string]string `json:"environment" yaml:"environment,omitempty"`
	Expose      []int `json:"expose" yaml:"expose,omitempty"`
}

type ComposeFile struct {
	Version  string `json:"version"`
	Services map[string]Service `json:"services"`
}

func (cb ComposeBackend) Up() {

}

func toString(v api.Volume) string {
	if "" == v.HostPath {
		return fmt.Sprintf("%s", v.ContainerPath)
	}

	if "" == v.Mode {
		return fmt.Sprintf("%s:%s", v.HostPath, v.ContainerPath)
	}

	var hostpath string
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic("Please ensure the current dir can be accessed")
	}

	if !filepath.IsAbs(v.HostPath) {
		hostpath = filepath.Join(dir, ".local", v.HostPath)
	}

	return fmt.Sprintf("%s:%s:%s", hostpath, v.ContainerPath, strings.ToLower(v.Mode))
}

func (cb ComposeBackend) ToComposeFile(s api.Stack) string {
	services := s.GetServices()
	composeServices := make(map[string]Service, 0)

	for name, service := range services {
		if !service.IsBuildable() {
			volumes := make([]string, 0)
			for _, v := range service.GetVolumes() {
				volumes = append(volumes, toString(v))
			}
			env := service.GetEnv()
			links := service.GetLinks()
			for _, link := range links {
				fmt.Println(link)
				env[fmt.Sprintf("%s_HOST", strings.ToUpper(link))] = link
				env[fmt.Sprintf("%s_PORT", strings.ToUpper(link))] = fmt.Sprintf("%d", services[link].GetExpose()[0])
			}


			composeServices[name] = Service{
				Image: service.GetImage(),
				Links: links,
				Volumes: volumes,
				Environment: env,
				Expose: service.GetExpose(),
				Ports: Map(service.GetExpose(), func(port int) string {return fmt.Sprintf("%d:%d", port, port)}),
			}
		} else {
			volumes := make([]string, 0)
			for _, v := range service.GetVolumes() {
				volumes = append(volumes, toString(v))
			}

			dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
			if err != nil {
				panic("Please ensure the current dir can be accessed")
			}

			volumes = append(volumes, "/var/run/docker.sock:/var/run/docker.sock")
			volumes = append(volumes, fmt.Sprintf("%s:/codebase", dir))

			env := service.GetEnv()
			links := service.GetLinks()
			for _, link := range links {
				env[fmt.Sprintf("%s_HOST", strings.ToUpper(link))] = link
				env[fmt.Sprintf("%s_PORT", strings.ToUpper(link))] = fmt.Sprintf("%d", services[link].GetExpose()[0])
			}


			composeServices["runtime"] = Service{
				Image: service.GetBuild().Name,
				Entrypoint: "/bin/sh",
				Command: "-c 'tail -f /dev/null'",
				Volumes: volumes,
				Links: links,
				Environment: env,
				Expose: service.GetExpose(),
				Ports: Map(service.GetExpose(), func(port int) string {return fmt.Sprintf("%d:%d", port, port)}),
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


func Map(src []int, mapper func(int) string) []string {
	result := make([]string, 0)
	for _, item := range src {
		result = append(result, mapper(item))
	}
	return result
}