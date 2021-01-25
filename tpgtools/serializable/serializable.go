// Package serializable contains a function that returns the list of resources that tpgtools currently supports.
package serializable

import (
	"path/filepath"
	"regexp"
)

// Service contains the name of a GCP service and the resources it contains in tpgtools.
type Service struct {
	Name      string
	Resources []string
}

// ListOfResources returns a list of resources that tpgtools currently supports.
func ListOfResources(pathPrefix string) ([]*Service, error) {
	pathglob := "/api/*/*.yaml"
	pathregex := `/api/(?P<service>[a-z_-]*)/(?P<resource>[a-z_-]*).yaml`

	path := pathPrefix + pathglob
	matches, err := filepath.Glob(path)
	if err != nil {
		return nil, err
	}
	var services []*Service
	r := regexp.MustCompile(pathPrefix + pathregex)
	for _, match := range matches {
		result := r.FindAllStringSubmatch(match, -1)
		if len(result) > 0 && len(result[0]) > 1 {
			service := findServiceInList(result[0][1], services)
			if service == nil {
				services = append(services, &Service{
					Name:      result[0][1],
					Resources: []string{result[0][2]},
				})
			} else {
				service.Resources = append(service.Resources, result[0][2])
			}
		}
	}
	return services, nil
}
func findServiceInList(name string, services []*Service) *Service {
	for _, s := range services {
		if s.Name == name {
			return s
		}
	}
	return nil
}
