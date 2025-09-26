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

func (cb *Client) ApproveDownstreamGenAndTest(prNumber, commitSha string) error {
	buildId, err := getPendingBuildId(PROJECT_ID, commitSha)
	if err != nil {
		return err
	}

	if buildId == "" {
		fmt.Printf("WARNING: Failed to find pending build for PR %s\nThis build may have been approved already.\n", prNumber)
		return nil
	}

	err = approveBuild(PROJECT_ID, buildId)
	if err != nil {
		return err
	}

	return nil
}

func getPendingBuildId(projectId, commitSha string) (string, error) {
	COMMUNITY_CHECKER_TRIGGER, ok := os.LookupEnv("DOWNSTREAM_GENERATION_AND_TEST_TRIGGER")
	if !ok {
		return "", fmt.Errorf("Did not provide DOWNSTREAM_GENERATION_AND_TEST_TRIGGER environment variable")
	}

	ctx := context.Background()

	c, err := cloudbuildv1.NewService(ctx)
	if err != nil {
		return "", err
	}

	filter := fmt.Sprintf("trigger_id=%s AND status=PENDING", COMMUNITY_CHECKER_TRIGGER)
	// Builds will be sorted by createTime, descending order.
	// 50 should be enough to include the one needs auto approval
	pageSize := int64(50)

	builds, err := c.Projects.Builds.List(projectId).Filter(filter).PageSize(pageSize).Do()
	if err != nil {
		return "", err
	}

	for _, build := range builds.Builds {
		if build.Substitutions["COMMIT_SHA"] == commitSha {
			return build.Id, nil
		}
	}

	return "", nil
}

func approveBuild(projectId, buildId string) error {
	ctx := context.Background()

	c, err := cloudbuildv1.NewService(ctx)
	if err != nil {
		return err
	}

	name := fmt.Sprintf("projects/%s/builds/%s", projectId, buildId)

	approveBuildRequest := &cloudbuildv1.ApproveBuildRequest{
		ApprovalResult: &cloudbuildv1.ApprovalResult{
			Decision: "APPROVED",
		},
	}

	_, err = c.Projects.Builds.Approve(name, approveBuildRequest).Do()
	if err != nil {
		return err
	}

	fmt.Println("Auto approved build ", buildId)

	return nil
}
