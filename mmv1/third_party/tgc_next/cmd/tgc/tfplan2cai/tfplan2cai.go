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
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/cmd/tgc/common"
	"github.com/spf13/cobra"
)

const cmdDesc = `
Support the convertion from Terraform resources data into the Cloud Asset Inventory(CAI) assets.

Supported Terraform versions = 0.12+`

func NewCmd(rootOptions *common.RootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "tfplan2cai",
		Short:         "Convert a Terraform plan to CAI assets",
		Long:          cmdDesc,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.AddCommand(newConvertCmd(rootOptions))

	return cmd
}
