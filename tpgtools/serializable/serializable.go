// Copyright 2021 Google LLC. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
