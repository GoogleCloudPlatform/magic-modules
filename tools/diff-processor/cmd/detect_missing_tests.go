package cmd

import (
	"encoding/json"
	newProvider "google/provider/new/google/provider"
	oldProvider "google/provider/old/google/provider"
	"io"

	"fmt"
	"os"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/detector"
	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/GoogleCloudPlatform/magic-modules/tools/test-reader/reader"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

const detectMissingTestsDesc = "Run the missing test detector using the given services directory"

type detectMissingTestsOptions struct {
	rootOptions *rootOptions
	stdout      io.Writer
}

func newDetectMissingTestsCmd(rootOptions *rootOptions) *cobra.Command {
	o := &detectMissingTestsOptions{
		rootOptions: rootOptions,
		stdout:      os.Stdout,
	}
	return &cobra.Command{
		Use:   "detect-missing-tests SERVICES_DIR",
		Short: detectMissingTestsDesc,
		Long:  detectMissingTestsDesc,
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return o.run(args)
		},
	}
}

func (o *detectMissingTestsOptions) run(args []string) error {
	allTests, errs := reader.ReadAllTests(args[0])
	for path, err := range errs {
		glog.Infof("error reading path: %s, err: %v", path, err)
	}

	schemaDiff := diff.ComputeSchemaDiff(oldProvider.ResourceMap(), newProvider.ResourceMap())

	missingTests, err := detector.DetectMissingTests(schemaDiff, allTests)
	if err != nil {
		return fmt.Errorf("error detecting missing tests: %v", err)
	}
	if err := json.NewEncoder(o.stdout).Encode(missingTests); err != nil {
		return fmt.Errorf("error encoding json: %w", err)
	}
	return nil
}
