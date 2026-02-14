package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/google/go-licenses/licenses"
	"github.com/spf13/cobra"
)

const rootCmdDesc = "Utilities for license check."
const licenseConfidence = 0.9

var copyrightRegex = regexp.MustCompile(`(?i)copyright (\d{4}) google`)

type rootOptions struct {
	fileList []string
}

func newRootCmd() (*cobra.Command, *rootOptions, error) {
	o := &rootOptions{}
	command := &cobra.Command{
		Use:           "license-check",
		Short:         rootCmdDesc,
		Long:          rootCmdDesc,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(c *cobra.Command, args []string) error {
			return o.run()
		},
	}
	command.Flags().StringSliceVar(&o.fileList, "file-list", []string{}, "file list to check")
	return command, o, nil
}

func (o *rootOptions) run() error {
	if len(o.fileList) == 0 {
		return nil
	}
	foundErr := false
	for _, file := range o.fileList {
		ext := filepath.Ext(file)
		if ext != ".tmpl" && ext != ".go" && ext != ".yaml" && ext != ".yml" {
			continue
		}

		if err := checkLicenseType(file); err != nil {
			fmt.Fprintf(os.Stderr, "File %s failed: %s.\n", file, err)
			foundErr = true
		}
		if err := checkCopyright(file, time.Now().Year()); err != nil {
			fmt.Fprintf(os.Stderr, "File %s failed: %s.\n", file, err)
			foundErr = true
		}
	}
	if foundErr {
		return fmt.Errorf("found file failing license check")
	}
	return nil
}

// Execute is the entry-point for all commands.
// This lets us keep all new command functions private.
func Execute() {
	rootCmd, _, err := newRootCmd()
	if err != nil {
		fmt.Printf("Error creating root logger: %s", err)
		os.Exit(1)
	}
	err = rootCmd.Execute()
	if err == nil {
		os.Exit(0)
	} else {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func checkLicenseType(filePath string) error {
	classifier, err := licenses.NewClassifier(licenseConfidence)
	if err != nil {
		return fmt.Errorf("failed to create license classifier: %w", err)
	}
	licenseName, _, err := classifier.Identify(filePath)
	if err != nil {
		return fmt.Errorf("failed to identify license for %s: %w", filePath, err)
	}
	if !strings.Contains(licenseName, "Apache") {
		return fmt.Errorf("found license type %s, expect Apache", licenseName)
	}
	return nil
}

func checkCopyright(filePath string, year int) error {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	if !copyrightRegex.MatchString(string(b)) {
		return fmt.Errorf("expected copyright string not found")
	}

	foundYears := copyrightRegex.FindStringSubmatch(string(b))
	if foundYears[1] != fmt.Sprintf("%d", year) {
		return fmt.Errorf("copyright year is not the latest year")
	}
	return nil
}
