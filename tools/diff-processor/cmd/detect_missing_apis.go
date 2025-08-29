package cmd

import (
	"encoding/json"
	"io"

	"fmt"
	"os"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/detector"
	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/spf13/cobra"

	newTransport "google/provider/new/google/transport"
	oldTransport "google/provider/old/google/transport"
)

const detectMissingAPIsDesc = "Run the missing API detector using the given services directory"

type detectMissingAPIsOptions struct {
	rootOptions       *rootOptions
	computeSchemaDiff func() diff.SchemaDiff
	stdout            io.Writer
}

func newDetectMissingAPIsCmd(rootOptions *rootOptions) *cobra.Command {
	o := &detectMissingAPIsOptions{
		rootOptions: rootOptions,
		computeSchemaDiff: func() diff.SchemaDiff {
			return schemaDiff
		},
		stdout: os.Stdout,
	}
	return &cobra.Command{
		Use:   "detect-missing-apis test-infra-tf-file",
		Short: detectMissingAPIsDesc,
		Long:  detectMissingAPIsDesc,
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return o.run(args)
		},
	}
}

func (o *detectMissingAPIsOptions) run(args []string) error {
	fileName := args[0]
	missingAPIs, err := detector.DetectMissingAPIs(o.computeSchemaDiff(), oldTransport.DefaultBasePaths, newTransport.DefaultBasePaths, fileName)
	if err != nil {
		return fmt.Errorf("error detecting missing APIs: %v", err)
	}
	if err := json.NewEncoder(o.stdout).Encode(missingAPIs); err != nil {
		return fmt.Errorf("error encoding json: %w", err)
	}
	return nil
}
