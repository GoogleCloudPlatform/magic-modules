package cmd

import (
	newProvider "google/provider/new/google/provider"
	oldProvider "google/provider/old/google/provider"

	"io"
	"os"
	"sort"

	"github.com/GoogleCloudPlatform/magic-modules/.ci/diff-processor/diff"
	"github.com/GoogleCloudPlatform/magic-modules/.ci/diff-processor/rules"
	"github.com/spf13/cobra"
)

const breakingChangesDesc = `Check for breaking changes between the new / old Terraform provider versions.`

type breakingChangesOptions struct {
	rootOptions       *rootOptions
	computeSchemaDiff func() diff.SchemaDiff
	stdout            io.Writer
}

func newBreakingChangesCmd(rootOptions *rootOptions) *cobra.Command {
	o := &breakingChangesOptions{
		rootOptions: rootOptions,
		computeSchemaDiff: func() diff.SchemaDiff {
			return diff.ComputeSchemaDiff(oldProvider.ResourceMap(), newProvider.ResourceMap())
		},
		stdout: os.Stdout,
	}
	cmd := &cobra.Command{
		Use:   "breaking-changes",
		Short: breakingChangesDesc,
		Long:  breakingChangesDesc,
		RunE: func(c *cobra.Command, args []string) error {
			return o.run()
		},
	}
	return cmd
}
func (o *breakingChangesOptions) run() error {
	schemaDiff := o.computeSchemaDiff()
	breakingChanges := rules.ComputeBreakingChanges(schemaDiff)
	sort.Strings(breakingChanges)
	for _, breakingChange := range breakingChanges {
		_, err := o.stdout.Write([]byte(breakingChange + "\n"))
		if err != nil {
			return err
		}
	}
	return nil
}
