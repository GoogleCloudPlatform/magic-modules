package cmd

import (
	"encoding/json"
	"fmt"
	newProvider "google/provider/new/google/provider"
	oldProvider "google/provider/old/google/provider"

	"io"
	"os"
	"sort"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/breaking_changes"
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
	breakingChanges := breaking_changes.ComputeBreakingChanges(schemaDiff)
	sort.Slice(breakingChanges, func(i, j int) bool {
		return breakingChanges[i].Message < breakingChanges[j].Message
	})
	if err := json.NewEncoder(o.stdout).Encode(breakingChanges); err != nil {
		return fmt.Errorf("error encoding json: %w", err)
	}
	return nil
}
