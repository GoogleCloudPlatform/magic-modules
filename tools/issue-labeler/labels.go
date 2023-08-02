package main

import (
	"regexp"
	"sort"

	"github.com/golang/glog"
)

func serviceLabels(issueBody string, enrolledTeams map[string][]string) []string {
	sectionRegexp := regexp.MustCompile(`### (New or )?Affected Resource\(s\)[^#]+`)
	affectedResources := sectionRegexp.FindString(issueBody)
	var results []string
	for label, resources := range enrolledTeams {
		for _, resource := range resources {
			resourceRegexp := regexp.MustCompile(resource)
			if resourceRegexp.MatchString(affectedResources) {
				glog.Infof("found resource %q, applying label %q", resource, label)
				results = append(results, label)
				break
			}
		}
	}

	sort.Strings(results)

	return results
}
