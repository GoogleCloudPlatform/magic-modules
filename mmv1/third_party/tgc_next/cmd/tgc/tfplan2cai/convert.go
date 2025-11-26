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

package tfplan2cai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/cmd/tgc/common"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/caiasset"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/tfplan2cai"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const convertDesc = `
This command will convert a Terraform plan json file into Cloud Asset Inventory(CAI) assets
and output them as a JSON array.

Note:
  Only supported resources will be converted. Non supported resources are
  omitted from results.
  Run "tgc tfplan2cai list-supported-resources" to see all supported
  resources.

Example:
tgc tfplan2cai convert ./example/terraform.tfplan --project my-project \
    --ancestry organization/my-org/folder/my-folder
`

func multiEnvSearch(ks []string) string {
	for _, k := range ks {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return ""
}

type convertOptions struct {
	project     string
	ancestry    string
	offline     bool
	rootOptions *common.RootOptions
	outputPath  string
	dryRun      bool
}

var origConvertFunc = func(ctx context.Context, path, project, zone, region string, ancestry map[string]string, offline bool, errorLogger *zap.Logger, userAgent string) ([]caiasset.Asset, error) {
	jsonPlan, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %s", path, err)
	}

	return tfplan2cai.Convert(ctx, jsonPlan, &tfplan2cai.Options{
		ErrorLogger:    errorLogger,
		Offline:        offline,
		DefaultProject: project,
		DefaultRegion:  region,
		DefaultZone:    zone,
		UserAgent:      userAgent,
		AncestryCache:  ancestry,
	})
}

var convertFunc = origConvertFunc

func newConvertCmd(rootOptions *common.RootOptions) *cobra.Command {
	o := &convertOptions{
		rootOptions: rootOptions,
	}

	cmd := &cobra.Command{
		Use:   "convert TFPLAN_JSON",
		Short: "convert a Terraform plan to Google CAI assets",
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

	cmd.Flags().StringVar(&o.project, "project", "", "Provider project override (override the default project configuration assigned to the google terraform provider when converting resources)")
	cmd.Flags().StringVar(&o.ancestry, "ancestry", "", "Override the ancestry location of the project when validating resources")
	cmd.Flags().BoolVar(&o.offline, "offline", false, "Do not make network requests")
	cmd.Flags().StringVar(&o.outputPath, "output-path", "", "If specified, write the convert result into the specified output file")
	cmd.Flags().BoolVar(&o.dryRun, "dry-run", false, "Only parse & validate args")
	cmd.Flags().MarkHidden("dry-run")

	return cmd
}

func (o *convertOptions) validateArgs(args []string) error {
	if len(args) != 1 {
		return errors.New("missing required argument TFPLAN_JSON")
	}
	if o.offline && o.ancestry == "" {
		return errors.New("please set ancestry via --ancestry in offline mode")
	}
	return nil
}

func (o *convertOptions) run(plan string) error {
	ctx := context.Background()
	ancestryCache := map[string]string{}
	if o.project != "" {
		ancestryCache[o.project] = o.ancestry
	}
	zone := multiEnvSearch([]string{
		"GOOGLE_ZONE",
		"GCLOUD_ZONE",
		"CLOUDSDK_COMPUTE_ZONE",
	})

	region := multiEnvSearch([]string{
		"GOOGLE_REGION",
		"GCLOUD_REGION",
		"CLOUDSDK_COMPUTE_REGION",
	})
	userAgent := "tfplan2cai"
	assets, err := convertFunc(ctx, plan, o.project, zone, region, ancestryCache, o.offline, o.rootOptions.ErrorLogger, userAgent)
	if err != nil {
		return err
	}

	if len(o.outputPath) > 0 {
		f, err := os.OpenFile(o.outputPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
		if err != nil {
			return err
		}
		defer f.Close()

		if err := json.NewEncoder(f).Encode(assets); err != nil {
			return fmt.Errorf("encoding json: %w", err)
		}
		return nil
	}

	if o.rootOptions.UseStructuredLogging {
		o.rootOptions.OutputLogger.Info(
			"converted resources",
			zap.Any("resource_body", assets),
		)
		return nil
	}

	// Legacy behavior
	if err := json.NewEncoder(os.Stdout).Encode(assets); err != nil {
		return fmt.Errorf("encoding json: %w", err)
	}
	return nil
}
