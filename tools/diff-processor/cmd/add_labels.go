package cmd

import (
	newProvider "google/provider/new/google/provider"
	oldProvider "google/provider/old/google/provider"

	"fmt"
	"strconv"
	"strings"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/labels"
	"github.com/GoogleCloudPlatform/magic-modules/tools/issue-labeler/labeler"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

const addLabelsDesc = `Add labels to a PR based on changed resources.`

type addLabelsOptions struct {
	rootOptions       *rootOptions
	computeSchemaDiff func() diff.SchemaDiff
	enrolledTeamsYaml []byte
	getIssue          func(repository string, id uint64) (labeler.Issue, error)
	updateIssues      func(repository string, issueUpdates []labeler.IssueUpdate, dryRun bool)
	dryRun            bool
}

func newAddLabelsCmd(rootOptions *rootOptions) *cobra.Command {
	o := &addLabelsOptions{
		rootOptions: rootOptions,
		computeSchemaDiff: func() diff.SchemaDiff {
			return diff.ComputeSchemaDiff(oldProvider.ResourceMap(), newProvider.ResourceMap())
		},
		enrolledTeamsYaml: labeler.EnrolledTeamsYaml,
		getIssue:          labels.GetIssue,
		updateIssues:      labeler.UpdateIssues,
	}
	cmd := &cobra.Command{
		Use:   "add-labels PR_ID [--dry-run]",
		Short: addLabelsDesc,
		Long:  addLabelsDesc,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return o.run(args)
		},
	}
	cmd.Flags().BoolVar(&o.dryRun, "dry-run", false, "Do a dry run without updating labels")
	return cmd
}
func (o *addLabelsOptions) run(args []string) error {
	prId, err := strconv.ParseUint(args[0], 10, 0)
	if err != nil {
		return fmt.Errorf("PR_ID must be an unsigned integer: %w", err)
	}

	repository := "GoogleCloudPlatform/magic-modules"
	issue, err := o.getIssue(repository, prId)

	if err != nil {
		return fmt.Errorf("Error retrieving PR data: %w", err)
	}

	hasServiceLabels := false
	oldLabels := make(map[string]struct{}, len(issue.Labels))
	for _, label := range issue.Labels {
		oldLabels[label.Name] = struct{}{}
		if strings.HasPrefix(label.Name, "service/") {
			hasServiceLabels = true
		}
	}
	if hasServiceLabels {
		return nil
	}

	schemaDiff := o.computeSchemaDiff()
	affectedResources := maps.Keys(schemaDiff)
	regexpLabels, err := labeler.BuildRegexLabels(o.enrolledTeamsYaml)
	if err != nil {
		return fmt.Errorf("Error building regex labels: %w", err)
	}

	newLabels := make(map[string]struct{}, len(oldLabels))
	for label, _ := range oldLabels {
		newLabels[label] = struct{}{}
	}
	for _, label := range labeler.ComputeLabels(affectedResources, regexpLabels) {
		newLabels[label] = struct{}{}
	}

	// Only update the issue if new labels should be added
	if len(newLabels) != len(oldLabels) {
		issueUpdate := labeler.IssueUpdate{
			Number:    prId,
			Labels:    maps.Keys(newLabels),
			OldLabels: maps.Keys(oldLabels),
		}

		o.updateIssues(repository, []labeler.IssueUpdate{issueUpdate}, o.dryRun)
	}

	return nil
}
