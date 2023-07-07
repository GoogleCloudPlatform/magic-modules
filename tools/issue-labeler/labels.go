package main

import (
	"regexp"
	"sort"
	"strings"

	"github.com/golang/glog"
)

func labels(issueBody string, enrolledTeams map[string][]string) string {
	sectionRegexp := regexp.MustCompile(`### (New or )?Affected Resource\(s\)[^#]+`)
	affectedResources := sectionRegexp.FindString(issueBody)
	var labels []string
	for label, resources := range enrolledTeams {
		for _, resource := range resources {
			if strings.Contains(affectedResources, resource) {
				glog.Infof("found resource %q, applying label %q", resource, label)
				labels = append(labels, "\""+label+"\"")
				break
			}
		}
	}

	if len(labels) > 0 {
		sort.Strings(labels)
		return "[" + strings.Join(labels, ", ") + "]"
	}
	return ""
}
