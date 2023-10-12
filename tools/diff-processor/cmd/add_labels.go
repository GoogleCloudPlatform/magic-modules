package cmd

import (
	newProvider "google/provider/new/google/provider"
	oldProvider "google/provider/old/google/provider"

	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

    "github.com/davecgh/go-spew/spew"
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
	stdout            io.Writer
	prId              uint64
	getIssue          func(id uint64) (labeler.Issue, error)
	dryRun            bool
}

func newAddLabelsCmd(rootOptions *rootOptions) *cobra.Command {
	o := &addLabelsOptions{
		rootOptions: rootOptions,
		computeSchemaDiff: func() diff.SchemaDiff {
			return diff.ComputeSchemaDiff(oldProvider.ResourceMap(), newProvider.ResourceMap())
		},
		enrolledTeamsYaml: labeler.EnrolledTeamsYaml,
		stdout: os.Stdout,
		getIssue: labels.GetIssue,
	}
	cmd := &cobra.Command{
		Use:   "add-labels PR_ID [--dry-run]",
		Short: addLabelsDesc,
		Long:  addLabelsDesc,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("Missing pull request ID.")
			}

			prId, err := strconv.ParseUint(args[0], 10, 0)

			if err != nil {
				return fmt.Errorf("PR_ID must be an unsigned integer: %w", err)
			}

			o.prId = prId

			return nil

		},
		RunE: func(c *cobra.Command, args []string) error {
			return o.run()
		},
	}
	cmd.Flags().BoolVar(&o.dryRun, "dry-run", false, "Do a dry run without updating labels")
	return cmd
}
func (o *addLabelsOptions) run() error {
	issue, err := o.getIssue(o.prId)

	if err != nil {
		return fmt.Errorf("Error retrieving PR data: %w", err)
	}

	hasServiceLabels := false
	oldLabels := make([]string, len(issue.Labels))
	for _, label := range issue.Labels {
		oldLabels = append(oldLabels, label.Name)
		if strings.HasPrefix(label.Name, "service/") {
			hasServiceLabels = true
		}
	}
	if hasServiceLabels {
		return nil
	}

	schemaDiff := o.computeSchemaDiff()
	affectedResources := maps.Keys(schemaDiff)
	spew.Dump(schemaDiff)
	regexpLabels, err := labeler.BuildRegexLabels(o.enrolledTeamsYaml)
	if err != nil {
		return fmt.Errorf("Error building regex labels: %w", err)
	}

	issueUpdate := labeler.IssueUpdate{
		Number: o.prId,
		Labels: labeler.ComputeLabels(affectedResources, regexpLabels),
		OldLabels: oldLabels,
	}
	labeler.UpdateIssues([]labeler.IssueUpdate{issueUpdate}, o.dryRun)

	return nil
}
