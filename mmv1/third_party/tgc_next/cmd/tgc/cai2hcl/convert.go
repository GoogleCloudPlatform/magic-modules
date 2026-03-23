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

package cai2hcl

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/cmd/tgc/common"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/cai2hcl"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/caiasset"
)

const convertDesc = `
This command will convert a JSON file of Cloud Asset Inventory(CAI) assets
into Terraform HCL (HashiCorp Configuration Language) native syntax.

Note:
  Only supported resources will be converted. Non supported resources are
  omitted from results.

Example:
tgc cai2hcl convert ./example/caiassets.json
`

type convertOptions struct {
	rootOptions *common.RootOptions
	outputPath  string
	dryRun      bool
}

var origConvertFunc = func(ctx context.Context, path string, errorLogger *zap.Logger) ([]byte, error) {
	assetPayload, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %s", path, err)
	}

	var assets []caiasset.Asset
	if err := json.Unmarshal(assetPayload, &assets); err != nil {
		return nil, err
	}

	return cai2hcl.Convert(assets, &cai2hcl.Options{
		ErrorLogger: errorLogger,
	})
}

var convertFunc = origConvertFunc

func newConvertCmd(rootOptions *common.RootOptions) *cobra.Command {
	o := &convertOptions{
		rootOptions: rootOptions,
	}

	cmd := &cobra.Command{
		Use:   "convert CAI_ASSETS_JSON",
		Short: "convert Google CAI assets into Terraform HCL",
		Long:  convertDesc,
		PreRunE: func(c *cobra.Command, args []string) error {
			return o.validateArgs(args)
		},
		RunE: func(c *cobra.Command, args []string) error {
			if o.dryRun {
				return nil
			}

			return o.run(args[0])
		},
	}

	cmd.Flags().StringVar(&o.outputPath, "output-path", "", "If specified, write the convert result into the specified output file")
	cmd.Flags().BoolVar(&o.dryRun, "dry-run", false, "Only parse & validate args")
	cmd.Flags().MarkHidden("dry-run")

	return cmd
}

func (o *convertOptions) validateArgs(args []string) error {
	if len(args) != 1 {
		return errors.New("missing required argument CAI_ASSETS_JSON")
	}
	return nil
}

func (o *convertOptions) run(path string) error {
	ctx := context.Background()

	hclBlocks, err := convertFunc(ctx, path, o.rootOptions.ErrorLogger)
	if err != nil {
		return err
	}

	if len(o.outputPath) > 0 {
		f, err := os.OpenFile(o.outputPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
		if err != nil {
			return err
		}
		defer f.Close()

		err = ioutil.WriteFile(o.outputPath, hclBlocks, 0644)
		if err != nil {
			return err
		}

		return nil
	}

	if o.rootOptions.UseStructuredLogging {
		o.rootOptions.OutputLogger.Info(
			"converted resources",
			zap.String("resource_body", string(hclBlocks)),
		)
		return nil
	}

	os.Stdout.Write(hclBlocks)
	os.Stdout.Write([]byte("\n")) // Add a newline

	return nil
}
