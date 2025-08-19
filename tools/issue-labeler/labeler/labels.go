package labeler

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"

	_ "embed"

	"github.com/golang/glog"
	"github.com/google/go-github/v68/github"
	"gopkg.in/yaml.v2"
)

var sectionRegexp = regexp.MustCompile(`#+ (New or )?Affected Resource\(s\)[^#]+`)
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

type LabelChange struct {
	Name        string
	Color       string
	IsNew       bool
	NeedsUpdate bool
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

// EnsureLabelsWithColor applies the computed changes using the GitHub API
func EnsureLabelsWithColor(repository string, labelNames []string, color string) error {
	client := newGitHubClient()
	ctx := context.Background()
	owner, repo, err := splitRepository(repository)
	if err != nil {
		return fmt.Errorf("invalid repository format: %w", err)
	}

	// Get all existing labels first
	existingLabels, err := listLabels(repository)
	if err != nil {
		return fmt.Errorf("failed to list existing labels: %w", err)
	}

	changes := ComputeLabelChanges(existingLabels, labelNames, color)
	for _, change := range changes {
		if change.IsNew {
			_, _, err := client.Issues.CreateLabel(ctx, owner, repo, &github.Label{
				Name:  &change.Name,
				Color: &change.Color,
			})
			if err != nil {
				return fmt.Errorf("failed to create label %s: %w", change.Name, err)
			}
		} else if change.NeedsUpdate {
			_, _, err := client.Issues.EditLabel(ctx, owner, repo, change.Name, &github.Label{
				Color: &change.Color,
			})
			if err != nil {
				return fmt.Errorf("failed to update label %s color: %w", change.Name, err)
			}
		}
	}

	return nil
}

// ComputeLabelChanges determines which labels need to be created or updated
func ComputeLabelChanges(existingLabels []*github.Label, desiredLabels []string, desiredColor string) []LabelChange {
	labelMap := make(map[string]*github.Label)
	for _, label := range existingLabels {
		labelMap[label.GetName()] = label
	}

	changes := make([]LabelChange, 0, len(desiredLabels))
	desiredColor = strings.ToUpper(desiredColor)

	for _, labelName := range desiredLabels {
		change := LabelChange{
			Name:  labelName,
			Color: desiredColor,
		}

		if existingLabel, exists := labelMap[labelName]; exists {
			change.IsNew = false
			change.NeedsUpdate = strings.ToUpper(existingLabel.GetColor()) != desiredColor
		} else {
			change.IsNew = true
			change.NeedsUpdate = false
		}

		changes = append(changes, change)
	}

	return changes
}
