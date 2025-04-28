package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const rootCmdDesc = "Utilities for license check."

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
		b, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		ext := filepath.Ext(file)
		if ext != ".tmpl" && ext != ".go" && ext != ".yaml" {
			continue
		}
		if !strings.Contains(string(b), `Licensed under the Apache License, Version 2.0 (the "License");`) {
			fmt.Fprintf(os.Stderr, "File %s does not contain Apache License.\n", file)
			foundErr = true
		}
	}
	if foundErr {
		return fmt.Errorf("found file missing license")
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
