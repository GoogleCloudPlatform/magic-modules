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

	"github.com/spf13/cobra"

	"github.com/GoogleCloudPlatform/magic-modules/tools/issue-labeler/labeler"
)

var (
	// used for flags
	backfillSince  string
	backfillDryRun bool
)

var backfillIssueLabels = &cobra.Command{
	Use:   "backfill-issue-labels [--dry-run] [--since=1973-01-01]",
	Short: "Backfills labels on old issues",
	Long:  "Backfills labels on old issues",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// For now actual usage is handled inside UpdateIssues. This is just a new quick check.
		_, ok := os.LookupEnv("GITHUB_TOKEN")
		if !ok {
			return fmt.Errorf("did not provide GITHUB_TOKEN environment variable")
		}
		return execBackfillIssueLabels()
	},
}

func execBackfillIssueLabels() error {
	regexpLabels, err := labeler.BuildRegexLabels(labeler.EnrolledTeamsYaml)
	if err != nil {
		return fmt.Errorf("building regex labels: %w", err)
	}
	repository := "hashicorp/terraform-provider-google"
	issues, err := labeler.GetIssues(repository, backfillSince)
	if err != nil {
		return fmt.Errorf("getting github issues: %w", err)
	}
	issueUpdates := labeler.ComputeIssueUpdates(issues, regexpLabels)
	err = labeler.UpdateIssues(repository, issueUpdates, backfillDryRun)
	if err != nil {
		return fmt.Errorf("updating github issues: %w", err)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(backfillIssueLabels)
	backfillIssueLabels.Flags().BoolVar(&backfillDryRun, "dry-run", false, "Only log write actions instead of updating issues")
	backfillIssueLabels.Flags().StringVar(&backfillSince, "since", "1973-01-01", "Only apply labels to issues filed after given date")
}
