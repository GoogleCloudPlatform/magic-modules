/*
* Copyright 2024 Google LLC. All Rights Reserved.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */
package cmd

import (
	"flag"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/GoogleCloudPlatform/magic-modules/tools/issue-labeler/constants"
	"github.com/GoogleCloudPlatform/magic-modules/tools/issue-labeler/labeler"
)

var setupLabels = &cobra.Command{
	Use:   "setup-labels",
	Short: "Sets up labels for the relevant services",
	Long:  "Computes labels that should be added to an issue based on its body",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := args[0]
		return execSetupLabels(repo)
	},
}

func execSetupLabels(repo string) error {
	flag.Set("logtostderr", "true")
	regexpLabels, err := labeler.BuildRegexLabels(labeler.EnrolledTeamsYaml)
	if err != nil {
		return fmt.Errorf("building regex labels: %w", err)
	}
	var serviceLabels = make([]string, 0, len(regexpLabels))
	var serviceLabelMap = map[string]bool{}
	for _, r := range regexpLabels {
		_, alreadyExists := serviceLabelMap[r.Label]
		if alreadyExists {
			continue
		}
		serviceLabels = append(serviceLabels, r.Label)
		serviceLabelMap[r.Label] = true
	}
	err = labeler.EnsureLabelsWithColor(repo, serviceLabels, constants.GITHUB_YELLOW)
	return err
}

func init() {
	rootCmd.AddCommand(setupLabels)
}
