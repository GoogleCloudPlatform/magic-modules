package labeler

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"

	_ "embed"

	"github.com/golang/glog"
	"github.com/google/go-github/v61/github"
	"gopkg.in/yaml.v2"
)

var sectionRegexp = regexp.MustCompile(`### (New or )?Affected Resource\(s\)[^#]+`)
var commentRegexp = regexp.MustCompile(`<!--.*?-->`)
var resourceRegexp = regexp.MustCompile(`google_[\w*.]+`)

var (
	//go:embed enrolled_teams.yml
	EnrolledTeamsYaml []byte
)

type LabelData struct {
	Team      string   `yaml:"team,omitempty"`
	Resources []string `yaml:"resources"`
}

type RegexpLabel struct {
	Regexp *regexp.Regexp
	Label  string
}

func BuildRegexLabels(teamsYaml []byte) ([]RegexpLabel, error) {
	enrolledTeams := make(map[string]LabelData)
	regexpLabels := []RegexpLabel{}
	if err := yaml.Unmarshal(teamsYaml, &enrolledTeams); err != nil {
		return regexpLabels, fmt.Errorf("unmarshalling enrolled teams yaml: %w", err)
	}

	for label, data := range enrolledTeams {
		for _, resource := range data.Resources {
			exactResource := fmt.Sprintf("^%s$", resource)
			regexpLabels = append(regexpLabels, RegexpLabel{
				Regexp: regexp.MustCompile(exactResource),
				Label:  label,
			})
		}
	}

	sort.Slice(regexpLabels, func(i, j int) bool {
		return regexpLabels[i].Label < regexpLabels[j].Label
	})

	return regexpLabels, nil
}

func ExtractAffectedResources(body string) []string {
	section := sectionRegexp.FindString(body)
	section = commentRegexp.ReplaceAllString(section, "")
	if section != "" {
		return resourceRegexp.FindAllString(section, -1)
	}

	return []string{}
}

func ComputeLabels(resources []string, regexpLabels []RegexpLabel) []string {
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

	labels := []string{}
	for label := range labelSet {
		labels = append(labels, label)
	}
	sort.Strings(labels)

	return labels
}

// EnsureLabelsWithColor ensures service labels exist with the correct color
func EnsureLabelsWithColor(repository string, labelNames []string, color string) error {
	glog.Infof("Modifying labels for %s", repository)
	client := newGitHubClient()
	owner, repo, err := splitRepository(repository)
	if err != nil {
		return fmt.Errorf("invalid repository format: %w", err)
	}

	// Get all existing labels first
	existingLabels, err := listLabels(repository)
	if err != nil {
		return fmt.Errorf("failed to list existing labels: %w", err)
	}

	// Create a map for quick lookup
	labelMap := make(map[string]*github.Label)
	for _, label := range existingLabels {
		labelMap[label.GetName()] = label
	}

	ctx := context.Background()
	desiredColor := strings.ToUpper(color)

	// Process each desired label
	for _, labelName := range labelNames {
		if existingLabel, exists := labelMap[labelName]; exists {
			// Update if color doesn't match
			if strings.ToUpper(existingLabel.GetColor()) != desiredColor {
				existingLabel.Color = &desiredColor
				_, _, err = client.Issues.EditLabel(ctx, owner, repo, labelName, existingLabel)
				if err != nil {
					return fmt.Errorf("failed to update label %s color: %w", labelName, err)
				}
				glog.Infof("Updated label %q color from %q to %q", labelName, existingLabel.GetColor(), color)
			} else {
				glog.Infof("Label %q already exists with correct color", labelName)
			}
		} else {
			// Create new label
			_, _, err = client.Issues.CreateLabel(ctx, owner, repo, &github.Label{
				Name:  &labelName,
				Color: &color,
			})
			if err != nil {
				return fmt.Errorf("failed to create label %s: %w", labelName, err)
			}
			glog.Infof("Created new label %q with color %q", labelName, color)
		}
	}

	return nil
}
