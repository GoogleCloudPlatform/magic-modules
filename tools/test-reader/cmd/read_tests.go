package cmd

import (
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/magic-modules/tools/test-reader/reader"
	"github.com/spf13/cobra"
)

const readTestsDesc = "Run the missing test detector using the given services directory"

type readTestsOptions struct {
	rootOptions *rootOptions
	testPrefix  string
}

func newReadTestsCmd(rootOptions *rootOptions) *cobra.Command {
	o := &readTestsOptions{
		rootOptions: rootOptions,
	}
	cmd := &cobra.Command{
		Use:   "read-tests SERVICES_DIR",
		Short: readTestsDesc,
		Long:  readTestsDesc,
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return o.run(args)
		},
	}
	cmd.Flags().StringVar(&o.testPrefix, "test-prefix", "", "Only display results for matching tests")
	return cmd
}

func (o *readTestsOptions) run(args []string) error {
	allTests, errs := reader.ReadAllTests(args[0])
	for path, err := range errs {
		fmt.Printf("error reading path: %s, err: %v\n", path, err)
	}

	total := 0
	for _, test := range allTests {
		if !strings.HasPrefix(test.Name, o.testPrefix) {
			continue
		}
		fmt.Printf("%s:\n", test.Name)
		for index, step := range test.Steps {
			fmt.Printf("  Step %d:\n", index)
			for resourceType, resources := range step {
				for _, resource := range resources {
					fmt.Printf("    %s:\n", resourceType)
					for field, value := range resource {
						fmt.Printf("      %s: %v\n", field, value)
					}
				}
			}
		}
		fmt.Println("")
		total += 1
	}
	fmt.Printf("Found %d tests\n", total)
	return nil
}
