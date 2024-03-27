package cmd

import (
	newProvider "google/provider/new/google/provider"
	oldProvider "google/provider/old/google/provider"

	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/GoogleCloudPlatform/magic-modules/tools/issue-labeler/labeler"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

const changedSchemaLabelsDesc = `Compute service labels based on schema changes.`

type changedSchemaLabelsOptions struct {
	rootOptions       *rootOptions
	computeSchemaDiff func() diff.SchemaDiff
	enrolledTeamsYaml []byte
	stdout            io.Writer
}

func newChangedSchemaLabelsCmd(rootOptions *rootOptions) *cobra.Command {
	o := &changedSchemaLabelsOptions{
		rootOptions: rootOptions,
		computeSchemaDiff: func() diff.SchemaDiff {
			return diff.ComputeSchemaDiff(oldProvider.ResourceMap(), newProvider.ResourceMap())
		},
		enrolledTeamsYaml: labeler.EnrolledTeamsYaml,
		stdout:            os.Stdout,
	}
	cmd := &cobra.Command{
		Use:   "changed-schema-labels",
		Short: changedSchemaLabelsDesc,
		Long:  changedSchemaLabelsDesc,
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			return o.run()
		},
	}
	return cmd
}
func (o *changedSchemaLabelsOptions) run() error {
	schemaDiff := o.computeSchemaDiff()
	affectedResources := maps.Keys(schemaDiff)
	regexpLabels, err := labeler.BuildRegexLabels(o.enrolledTeamsYaml)
	if err != nil {
		return fmt.Errorf("Error building regex labels: %w", err)
	}

	labels := labeler.ComputeLabels(affectedResources, regexpLabels)

	if err = json.NewEncoder(o.stdout).Encode(labels); err != nil {
		return fmt.Errorf("Error encoding json: %w", err)
	}

	return nil
}
