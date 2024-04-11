package cmd

import (
	newProvider "google/provider/new/google/provider"
	oldProvider "google/provider/old/google/provider"

	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/detector"
	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/reader"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

const detectMissingTestsDesc = "Run the missing test detector using the given services directory"

type detectMissingTestsOptions struct {
	rootOptions *rootOptions
}

func newDetectMissingTestsCmd(rootOptions *rootOptions) *cobra.Command {
	o := &detectMissingTestsOptions{
		rootOptions: rootOptions,
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
	if len(missingTests) > 0 {
		funcs := template.FuncMap{
			"join": strings.Join,
			"backTickAll": func(ss []string) []string {
				rs := make([]string, len(ss))
				for i, s := range ss {
					rs[i] = fmt.Sprintf("`%s`", s)
				}
				return rs
			},
		}
		outputTemplate, err := template.New("missing_test_output.tmpl").Funcs(funcs).ParseFiles("missing_test_output.tmpl")
		if err != nil {
			return fmt.Errorf("Error parsing missing test template file: %s", err)
		}
		if err := outputTemplate.Execute(os.Stdout, missingTests); err != nil {
			return fmt.Errorf("Error executing missing test output template: %s", err)
		}
		for resourceName, missingTestInfo := range missingTests {
			glog.Infof("%s tests parsed: %v", resourceName, missingTestInfo.Tests)
		}
	}
	return nil
}
