// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const rootCmdDesc = `
Validate that a terraform plan conforms to a Constraint Framework
policy library written to expect Google CAI (Cloud Asset Inventory) data.

Supported Terraform versions = 0.12+`

var allowedVerbosity = map[string]struct{}{
	"debug":    {},
	"info":     {},
	"warning":  {},
	"error":    {},
	"critical": {},
	"none":     {},
}

type rootOptions struct {
	verbosity            string
	errorLogger          *zap.Logger
	outputLogger         *zap.Logger
	useStructuredLogging bool
}

func newRootCmd() (*cobra.Command, *rootOptions, error) {
	o := &rootOptions{
		useStructuredLogging: os.Getenv("USE_STRUCTURED_LOGGING") == "true",
	}

	cmd := &cobra.Command{
		Use:           "tfplan2cai",
		Short:         "Convert a terraform plan to CAI object",
		Long:          rootCmdDesc,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(c *cobra.Command, args []string) error {
			if _, ok := allowedVerbosity[o.verbosity]; !ok {
				return errors.New("verbosity must be one of: debug, info, warning, error, critical, none")
			}

			// set this up in PersistentPreRun because we need to wait for flags to be parsed
			// to have accurate verbosity.
			errorLogger := newErrorLogger(o.verbosity, o.useStructuredLogging, zapcore.Lock(os.Stderr))
			defer errorLogger.Sync()
			zap.RedirectStdLog(errorLogger)
			o.errorLogger = errorLogger

			outputLogger := newOutputLogger(zapcore.Lock(os.Stdout))
			defer outputLogger.Sync()
			o.outputLogger = outputLogger

			return nil
		},
	}

	cmd.PersistentFlags().StringVar(&o.verbosity, "verbosity", "info", "Set verbosity level. One of: debug, info, warning, error, critical, none.")

	cmd.AddCommand(newConvertCmd(o))
	cmd.AddCommand(newListSupportedResourcesCmd())

	return cmd, o, nil
}

// Execute is the entry-point for all commands.
// This lets us keep all new command functions private.
func Execute() {
	rootCmd, rootOptions, err := newRootCmd()

	if err != nil {
		fmt.Printf("Error creating root logger: %s", err)
		os.Exit(1)
	}

	err = rootCmd.Execute()

	if err == nil {
		os.Exit(0)
	} else {
		if rootOptions.errorLogger == nil {
			fmt.Println(err.Error())
		} else {
			rootOptions.errorLogger.Error(err.Error())
		}
		os.Exit(1)
	}
}
