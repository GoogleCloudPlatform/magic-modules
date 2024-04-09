package cmd

import (
	newProvider "google/provider/new/google/provider"
	oldProvider "google/provider/old/google/provider"

	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

const changedSchemaResourcesDesc = `Compute list of resources with changed schemas.`

type changedSchemaResourcesOptions struct {
	rootOptions       *rootOptions
	computeSchemaDiff func() diff.SchemaDiff
	stdout            io.Writer
}

func newChangedSchemaResourcesCmd(rootOptions *rootOptions) *cobra.Command {
	o := &changedSchemaResourcesOptions{
		rootOptions: rootOptions,
		computeSchemaDiff: func() diff.SchemaDiff {
			return diff.ComputeSchemaDiff(oldProvider.ResourceMap(), newProvider.ResourceMap())
		},
		stdout: os.Stdout,
	}
	cmd := &cobra.Command{
		Use:   "changed-schema-resources",
		Short: changedSchemaResourcesDesc,
		Long:  changedSchemaResourcesDesc,
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			return o.run()
		},
	}
	return cmd
}
func (o *changedSchemaResourcesOptions) run() error {
	schemaDiff := o.computeSchemaDiff()
	affectedResources := maps.Keys(schemaDiff)

	if err := json.NewEncoder(o.stdout).Encode(affectedResources); err != nil {
		return fmt.Errorf("Error encoding json: %w", err)
	}

	return nil
}
