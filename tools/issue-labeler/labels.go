package main

import (
	"fmt"
	"regexp"
	"sort"

	_ "embed"
	"github.com/golang/glog"
	"gopkg.in/yaml.v2"
)

var sectionRegexp = regexp.MustCompile(`### (New or )?Affected Resource\(s\)[^#]+`)
var commentRegexp = regexp.MustCompile(`<!--.*?-->`)
var resourceRegexp = regexp.MustCompile(`google_[\w*.]+`)

var (
	//go:embed enrolled_teams.yml
	enrolledTeamsYaml []byte
)

type labelData struct {
	Team      string   `yaml:"team,omitempty"`
	Resources []string `yaml:"resources"`
}

type regexpLabel struct {
	Regexp *regexp.Regexp
	Label  string
}

func buildRegexLabels(teamsYaml []byte) ([]regexpLabel, error) {
	enrolledTeams := make(map[string]labelData)
	regexpLabels := []regexpLabel{}
	if err := yaml.Unmarshal(teamsYaml, &enrolledTeams); err != nil {
		return regexpLabels, fmt.Errorf("Error unmarshalling enrolled teams yaml: %w", err)
	}

	for label, data := range enrolledTeams {
		for _, resource := range data.Resources {
			regexpLabels = append(regexpLabels, regexpLabel{
				Regexp: regexp.MustCompile(resource),
				Label:  label,
			})
		}
	}

	sort.Slice(regexpLabels, func(i, j int) bool {
		return regexpLabels[i].Label < regexpLabels[j].Label
	})

	return regexpLabels, nil
}

func extractAffectedResources(body string) []string {
	section := sectionRegexp.FindString(body)
	section = commentRegexp.ReplaceAllString(section, "")
	if section != "" {
		return resourceRegexp.FindAllString(section, -1)
	}

	return []string{}
}

func computeLabels(resources []string, regexpLabels []regexpLabel) []string {
	labelSet := make(map[string]struct{})
	for _, resource := range resources {
		for _, rl := range regexpLabels {
			if rl.Regexp.MatchString(resource) {
				glog.Infof("found resource %q, applying label %q", resource, rl.Label)
				labelSet[rl.Label] = struct{}{}
				break
			}
		}
	}

	var labels []string

	for l, _ := range labelSet {
		labels = append(labels, l)
	}
	sort.Strings(labels)

	return labels
}
