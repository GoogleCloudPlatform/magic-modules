// Copyright 2024 Google LLC
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

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/cmd/tgc/cai2hcl"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/cmd/tgc/common"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/cmd/tgc/tfplan2cai"
)

const rootCmdDesc = `
Convert between Terraform resource data and Google Cloud Platform's native API inventory format,
Cloud Asset Inventory (CAI) assets.

Supported Terraform versions = 0.12+`

var allowedVerbosity = map[string]struct{}{
	"debug":    {},
	"info":     {},
	"warning":  {},
	"error":    {},
	"critical": {},
	"none":     {},
}

func newRootCmd() (*cobra.Command, *common.RootOptions, error) {
	o := &common.RootOptions{
		UseStructuredLogging: os.Getenv("USE_STRUCTURED_LOGGING") == "true",
	}

	cmd := &cobra.Command{
		Use:           "tgc",
		Short:         "Convert between Terraform resource data and CAI assets",
		Long:          rootCmdDesc,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(c *cobra.Command, args []string) error {
			if _, ok := allowedVerbosity[o.Verbosity]; !ok {
				return errors.New("verbosity must be one of: debug, info, warning, error, critical, none")
			}

			// set this up in PersistentPreRun because we need to wait for flags to be parsed
			// to have accurate verbosity.
			errorLogger := common.NewErrorLogger(o.Verbosity, o.UseStructuredLogging, zapcore.Lock(os.Stderr))
			defer errorLogger.Sync()
			zap.RedirectStdLog(errorLogger)
			o.ErrorLogger = errorLogger

			outputLogger := common.NewOutputLogger(zapcore.Lock(os.Stdout))
			defer outputLogger.Sync()
			o.OutputLogger = outputLogger

			return nil
		},
	}

	cmd.PersistentFlags().StringVar(&o.Verbosity, "verbosity", "info", "Set verbosity level. One of: debug, info, warning, error, critical, none.")

	cmd.AddCommand(tfplan2cai.NewCmd(o))
	cmd.AddCommand(cai2hcl.NewCmd(o))

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
		if rootOptions.ErrorLogger == nil {
			fmt.Println(err.Error())
		} else {
			rootOptions.ErrorLogger.Error(err.Error())
		}
		os.Exit(1)
	}
}
