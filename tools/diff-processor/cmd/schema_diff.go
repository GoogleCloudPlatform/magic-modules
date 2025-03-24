package cmd

import (
	newProvider "google/provider/new/google/provider"
	oldProvider "google/provider/old/google/provider"

	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/spf13/cobra"
)

const schemaDiffDesc = `Return a simple summary of the schema diff for this build.`

var schemaDiff = diff.ComputeSchemaDiff(oldProvider.ResourceMap(), newProvider.ResourceMap())

type simpleSchemaDiff struct {
	AddedResources, ModifiedResources, RemovedResources []string
}

type schemaDiffOptions struct {
	rootOptions       *rootOptions
	computeSchemaDiff func() diff.SchemaDiff
	stdout            io.Writer
}

func newSchemaDiffCmd(rootOptions *rootOptions) *cobra.Command {
	o := &schemaDiffOptions{
		rootOptions: rootOptions,
		computeSchemaDiff: func() diff.SchemaDiff {
			return schemaDiff
		},
		stdout: os.Stdout,
	}
	cmd := &cobra.Command{
		Use:   "schema-diff",
		Short: schemaDiffDesc,
		Long:  schemaDiffDesc,
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			return o.run()
		},
	}
	return cmd
}
func (o *schemaDiffOptions) run() error {
	schemaDiff := o.computeSchemaDiff()

	simple := simpleSchemaDiff{}

	for k, d := range schemaDiff {
		if d.ResourceConfig.Old == nil {
			simple.AddedResources = append(simple.AddedResources, k)
		} else if d.ResourceConfig.New == nil {
			simple.RemovedResources = append(simple.RemovedResources, k)
		} else {
			simple.ModifiedResources = append(simple.ModifiedResources, k)
		}
	}

	sort.Strings(simple.AddedResources)
	sort.Strings(simple.ModifiedResources)
	sort.Strings(simple.RemovedResources)

	if err := json.NewEncoder(o.stdout).Encode(simple); err != nil {
		return fmt.Errorf("Error encoding json: %w", err)
	}

	return nil
}
