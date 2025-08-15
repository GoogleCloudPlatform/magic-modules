package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const rootCmdDesc = "Utilities for template checks."

type rootOptions struct {
}

func newRootCmd() (*cobra.Command, *rootOptions, error) {
	o := &rootOptions{}
	cmd := &cobra.Command{
		Use:           "template-check",
		Short:         rootCmdDesc,
		Long:          rootCmdDesc,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.AddCommand(newversionGuardCmd(o))
	return cmd, o, nil
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
