package cmd

import (
	"fmt"

	"io"
	"os"

	"github.com/GoogleCloudPlatform/magic-modules/tools/template-check/gotemplate"
	"github.com/spf13/cobra"
)

const versionGuardDesc = `Check the files for version guards`

type versionGuardOptions struct {
	rootOptions *rootOptions
	stdout      io.Writer
	fileList    []string
}

func newversionGuardCmd(rootOptions *rootOptions) *cobra.Command {
	o := &versionGuardOptions{
		rootOptions: rootOptions,
		stdout:      os.Stdout,
	}
	command := &cobra.Command{
		Use:   "version-guard",
		Short: versionGuardDesc,
		Long:  versionGuardDesc,
		RunE: func(c *cobra.Command, args []string) error {
			return o.run()
		},
	}

	command.Flags().StringSliceVar(&o.fileList, "file-list", []string{}, "file list to check")
	return command

}
func (o *versionGuardOptions) run() error {
	if len(o.fileList) == 0 {
		return nil
	}
	foundInvalidGuards := false
	for _, fileName := range o.fileList {
		results, err := gotemplate.CheckVersionGuardsForFile(fileName)
		if err != nil {
			return err
		}
		if len(results) > 0 {
			fmt.Fprintf(os.Stderr, "%s:\n", fileName)
			foundInvalidGuards = true
			for _, result := range results {
				fmt.Fprintf(os.Stderr, "  %s\n", result)
			}
		}
	}
	if foundInvalidGuards {
		return fmt.Errorf("found invalid version guards")
	}
	return nil
}
