package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/detector"
	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/spf13/cobra"
)

const detectMissingIdentityDesc = "Detect resources with ResourceIdentity that are missing SetResourceIdentityAttributes in CRUD functions or missing import identity tests"

type detectMissingIdentityOptions struct {
	rootOptions       *rootOptions
	computeSchemaDiff func() diff.SchemaDiff
	stdout            io.Writer
}

func newDetectMissingIdentityCmd(rootOptions *rootOptions) *cobra.Command {
	o := &detectMissingIdentityOptions{
		rootOptions: rootOptions,
		computeSchemaDiff: func() diff.SchemaDiff {
			return schemaDiff
		},
		stdout: os.Stdout,
	}
	return &cobra.Command{
		Use:   "detect-missing-identity SERVICES_DIR",
		Short: detectMissingIdentityDesc,
		Long:  detectMissingIdentityDesc,
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return o.run(args)
		},
	}
}

func (o *detectMissingIdentityOptions) run(args []string) error {
	servicesDir := args[0]

	// Get changed resources from schema diff
	schemaDiff := o.computeSchemaDiff()
	changedResources := make([]string, 0, len(schemaDiff))
	for resourceName := range schemaDiff {
		changedResources = append(changedResources, resourceName)
	}

	results, err := detector.DetectMissingIdentityCoverage(servicesDir, changedResources)
	if err != nil {
		return fmt.Errorf("error detecting missing identity coverage: %v", err)
	}
	if err := json.NewEncoder(o.stdout).Encode(results); err != nil {
		return fmt.Errorf("error encoding json: %w", err)
	}
	return nil
}
