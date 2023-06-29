package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/go-yaml/yaml"
	"github.com/golang/glog"
)

func main() {
	file, err := os.ReadFile("enrolled_teams.yaml")
	if err != nil {
		glog.Exitf("Error reading enrolled teams yaml: %v", err)
	}
	enrolledTeams := make(map[string][]string)
	err = yaml.Unmarshal(file, &enrolledTeams)
	if err != nil {
		glog.Exitf("Error unmarshalling enrolled teams yaml: %v", err)
	}
	issueBody := os.Getenv("ISSUE_BODY")
	sectionRegexp := regexp.MustCompile(`### (New or )?Affected Resource\(s\)[^#]+`)
	affectedResources := sectionRegexp.FindString(issueBody)
	var labels []string
	for label, resources := range enrolledTeams {
		for _, resource := range resources {
			if strings.Contains(affectedResources, resource) {
				labels = append(labels, "\""+label+"\"")
				break
			}
		}
	}

	if len(labels) > 0 {
		fmt.Println("[" + strings.Join(labels, ", ") + "]")
	}
}
