/*
* Copyright 2023 Google LLC. All Rights Reserved.
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
package cloudbuild

import (
	"context"
	"fmt"
	"os"

	cloudbuildv1 "google.golang.org/api/cloudbuild/v1"
)

func (cb *Client) TriggerMMPresubmitRuns(commitSha string, substitutions map[string]string) error {
	presubmitTriggerId, ok := os.LookupEnv("GENERATE_DIFFS_TRIGGER")
	if !ok {
		return fmt.Errorf("did not provide GENERATE_DIFFS_TRIGGER environment variable")
	}

	err := triggerCloudBuildRun(PROJECT_ID, presubmitTriggerId, REPO_NAME, commitSha, substitutions)
	if err != nil {
		return err
	}

	return nil
}

func triggerCloudBuildRun(projectId, triggerId, repoName, commitSha string, substitutions map[string]string) error {
	ctx := context.Background()
	c, err := cloudbuildv1.NewService(ctx)
	if err != nil {
		return fmt.Errorf("failed to create Cloud Build service client: %s", err)
	}

	repoSource := &cloudbuildv1.RepoSource{
		ProjectId:     projectId,
		RepoName:      repoName,
		CommitSha:     commitSha,
		Substitutions: substitutions,
	}

	_, err = c.Projects.Triggers.Run(projectId, triggerId, repoSource).Do()
	if err != nil {
		return fmt.Errorf("failed to create Cloud Build run: %s", err)
	}

	fmt.Println("Started Cloud Build Run: ", triggerId)
	return nil
}
