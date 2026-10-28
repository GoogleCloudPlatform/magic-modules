package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/detector"
	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/GoogleCloudPlatform/magic-modules/tools/test-reader/reader"
	"github.com/spf13/cobra"
)

const detectMissingTestsDesc = "Run the missing test detector using the given services directory"

type detectMissingTestsOptions struct {
	rootOptions       *rootOptions
	computeSchemaDiff func() diff.SchemaDiff
	stdout            io.Writer
}

func newDetectMissingTestsCmd(rootOptions *rootOptions) *cobra.Command {
	o := &detectMissingTestsOptions{
		rootOptions: rootOptions,
		computeSchemaDiff: func() diff.SchemaDiff {
			return schemaDiff
		},
		stdout: os.Stdout,
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
		log.Printf("error reading path: %s, err: %v", path, err)
	}

	sd := schemaDiff
	if o.computeSchemaDiff != nil {
		sd = o.computeSchemaDiff()
	}
	missingTests, err := detector.DetectMissingTests(sd, allTests)
	if err != nil {
		return fmt.Errorf("error detecting missing tests: %v", err)
	}
	if err := json.NewEncoder(o.stdout).Encode(missingTests); err != nil {
		return fmt.Errorf("error encoding json: %w", err)
	}
	return nil
}

const detectMissingUpdateTestsDesc = "Run the missing update test detector using the given services directory"

func newDetectMissingUpdateTestsCmd(rootOptions *rootOptions) *cobra.Command {
	o := &detectMissingTestsOptions{
		rootOptions: rootOptions,
		computeSchemaDiff: func() diff.SchemaDiff {
			return schemaDiff
		},
		stdout: os.Stdout,
	}
	return &cobra.Command{
		Use:   "detect-missing-update-tests SERVICES_DIR",
		Short: detectMissingUpdateTestsDesc,
		Long:  detectMissingUpdateTestsDesc,
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return o.runUpdateTests(args)
		},
	}
}

func (o *detectMissingTestsOptions) runUpdateTests(args []string) error {
	allTests, errs := reader.ReadAllTests(args[0])
	for path, err := range errs {
		log.Printf("error reading path: %s, err: %v", path, err)
	}

	sd := schemaDiff
	if o.computeSchemaDiff != nil {
		sd = o.computeSchemaDiff()
	}
	missingUpdateTests, err := detector.DetectMissingUpdateTests(sd, allTests)
	if err != nil {
		return fmt.Errorf("error detecting missing update tests: %v", err)
	}
	if err := json.NewEncoder(o.stdout).Encode(missingUpdateTests); err != nil {
		return fmt.Errorf("error encoding json: %w", err)
	}
	return nil
}
