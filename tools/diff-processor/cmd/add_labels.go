package cmd

import (
	newProvider "google/provider/new/google/provider"
	oldProvider "google/provider/old/google/provider"

	"io"
	"os"
	"sort"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/rules"
	issueLabeler "github.com/GoogleCloudPlatform/magic-modules/tools/issue-labeler"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

const addLabelsDesc = `Add labels to a PR based on changed resources.`

type addLabelsOptions struct {
	rootOptions       *rootOptions
	computeSchemaDiff func() diff.SchemaDiff
	enrolledTeamsYaml []byte
	stdout            io.Writer
}

func newAddLabelsCmd(rootOptions *rootOptions) *cobra.Command {
	o := &addLabelsOptions{
		rootOptions: rootOptions,
		computeSchemaDiff: func() diff.SchemaDiff {
			return diff.ComputeSchemaDiff(oldProvider.ResourceMap(), newProvider.ResourceMap())
		},
		enrolledTeamsYaml: issueLabeler.EnrolledTeamsYaml,
		stdout: os.Stdout,
	}
	cmd := &cobra.Command{
		Use:   "add-labels",
		Short: addLabelsDesc,
		Long:  addLabelsDesc,
		RunE: func(c *cobra.Command, args []string) error {
			return o.run()
		},
	}
	return cmd
}
func (o *addLabelsOptions) run() error {
	schemaDiff := o.computeSchemaDiff()
	regexpLabels := issueLabeler.BuildRegexLabels(o.enrolledTeamsYaml)
	affectedResources := 
	sort.Strings(breakingChanges)
	for _, breakingChange := range breakingChanges {
		_, err := o.stdout.Write([]byte(breakingChange + "\n"))
		if err != nil {
			return err
		}
	}
	return nil
}
