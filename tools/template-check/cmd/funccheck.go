package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/GoogleCloudPlatform/magic-modules/tools/template-check/gotemplate"
	"github.com/spf13/cobra"
)

const funcCheckDesc = `Check template files for invalid function calls`

type funcCheckOptions struct {
	rootOptions *rootOptions
	stdout      io.Writer
	fileList    []string
}

func newFuncCheckCmd(rootOptions *rootOptions) *cobra.Command {
	o := &funcCheckOptions{
		rootOptions: rootOptions,
		stdout:      os.Stdout,
	}
	command := &cobra.Command{
		Use:   "func-check",
		Short: funcCheckDesc,
		Long:  funcCheckDesc,
		RunE: func(c *cobra.Command, args []string) error {
			return o.run()
		},
	}
	command.Flags().StringSliceVar(&o.fileList, "file-list", []string{}, "file list to check")
	return command
}

func (o *funcCheckOptions) run() error {
	if len(o.fileList) == 0 {
		return nil
	}
	foundInvalidFuncs := false
	for _, fileName := range o.fileList {
		results, err := gotemplate.CheckInvalidFuncsForFile(fileName)
		if err != nil {
			return err
		}
		if len(results) > 0 {
			fmt.Fprintf(os.Stderr, "%s:\n", fileName)
			foundInvalidFuncs = true
			for _, result := range results {
				fmt.Fprintf(os.Stderr, "  %s\n", result)
			}
		}
	}
	if foundInvalidFuncs {
		return fmt.Errorf("found invalid template function calls")
	}
	return nil
}
