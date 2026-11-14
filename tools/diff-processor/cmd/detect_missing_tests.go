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
	rootOptions             *rootOptions
	computeSchemaDiff       func() diff.SchemaDiff
	stdout                  io.Writer
	exhaustive              bool
	exhaustiveChangedFields bool
	format                  string
}

func newDetectMissingTestsCmd(rootOptions *rootOptions) *cobra.Command {
	o := &detectMissingTestsOptions{
		rootOptions: rootOptions,
		computeSchemaDiff: func() diff.SchemaDiff {
			return schemaDiff
		},
		stdout: os.Stdout,
		format: "json",
	}
	cmd := &cobra.Command{
		Use:   "detect-missing-tests SERVICES_DIR",
		Short: detectMissingTestsDesc,
		Long:  detectMissingTestsDesc,
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return o.run(args)
		},
	}
	cmd.Flags().BoolVar(&o.exhaustive, "exhaustive", false, "Check that new fields are exhaustively tested (SET, MODIFY, UNSET at every nesting level)")
	cmd.Flags().BoolVar(&o.exhaustiveChangedFields, "exhaustive-changed-fields", false, "When --exhaustive is enabled, also check modified existing fields in addition to newly added fields")
	cmd.Flags().StringVar(&o.format, "format", "json", "Output format when --exhaustive is enabled (json or markdown)")
	return cmd
}

func (o *detectMissingTestsOptions) run(args []string) error {
	allTests, errs := reader.ReadAllTests(args[0])
	for path, err := range errs {
		log.Printf("error reading path: %s, err: %v", path, err)
	}

	sd := o.computeSchemaDiff()

	if o.exhaustive {
		if o.format != "" && o.format != "json" && o.format != "markdown" {
			return fmt.Errorf("unsupported --format %q (must be \"json\" or \"markdown\")", o.format)
		}
		exhaustiveMissing, err := detector.DetectExhaustiveMissingTests(sd, allTests, o.exhaustiveChangedFields)
		if err != nil {
			return fmt.Errorf("error detecting exhaustive missing tests: %v", err)
		}
		if o.format == "markdown" {
			_, err := io.WriteString(o.stdout, detector.FormatExhaustiveMarkdown(exhaustiveMissing))
			return err
		}
		if err := json.NewEncoder(o.stdout).Encode(exhaustiveMissing); err != nil {
			return fmt.Errorf("error encoding json: %w", err)
		}
		return nil
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
