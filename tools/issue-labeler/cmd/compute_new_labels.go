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
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/GoogleCloudPlatform/magic-modules/tools/issue-labeler/labeler"
)

var computeNewLabels = &cobra.Command{
	Use:   "compute-new-labels",
	Short: "Computes labels that should be added to an issue based on its body",
	Long:  "Computes labels that should be added to an issue based on its body",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return execComputeNewLabels()
	},
}

func execComputeNewLabels() error {
	regexpLabels, err := labeler.BuildRegexLabels(labeler.EnrolledTeamsYaml)
	if err != nil {
		return fmt.Errorf("building regex labels: %w", err)
	}
	issueBody := os.Getenv("ISSUE_BODY")
	affectedResources := labeler.ExtractAffectedResources(issueBody)
	labels := labeler.ComputeLabels(affectedResources, regexpLabels)

	// If there are more than 3 service labels, treat this as a cross-provider issue.
	// Note that labeler.ComputeLabels() currently only returns service labels, but
	// the logic here remains defensive in case that changes.
	var serviceLabels []string
	var nonServiceLabels []string
	for _, l := range labels {
		if strings.HasPrefix(l, "service/") {
			serviceLabels = append(serviceLabels, l)
		} else {
			nonServiceLabels = append(nonServiceLabels, l)
		}
	}
	if len(serviceLabels) > 3 {
		serviceLabels = []string{"service/terraform"}
	}
	labels = append(nonServiceLabels, serviceLabels...)

	if len(labels) > 0 {
		labels = append(labels, "forward/review")
		sort.Strings(labels)
		fmt.Println(`["` + strings.Join(labels, `", "`) + `"]`)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(computeNewLabels)
}
