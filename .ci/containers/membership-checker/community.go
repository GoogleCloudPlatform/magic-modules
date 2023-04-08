package main

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/api/cloudbuild/v1"
)

func approveCommunityChecker(prNumber, projectId, commitSha string) error {
	buildId, err := getPendingBuildId(projectId, commitSha)
	if err != nil {
		return err
	}

	if buildId == "" {
		return fmt.Errorf("Failed to find pending build for PR %s", prNumber)
	}

	err = approveBuild(projectId, buildId)
	if err != nil {
		return err
	}

	return nil
}

func postAwaitingApprovalBuildLink(prNumber, GITHUB_TOKEN, projectId, commitSha string) error {
	buildId, err := getPendingBuildId(projectId, commitSha)
	if err != nil {
		return err
	}

	if buildId == "" {
		return fmt.Errorf("Failed to find pending build for PR %s", prNumber)
	}

	targetUrl := fmt.Sprintf("https://console.cloud.google.com/cloud-build/builds;region=global/%s?project=%s", buildId, projectId)

	postBody := map[string]string{
		"context":    "Approve Build",
		"state":      "success",
		"target_url": targetUrl,
	}

	err = postBuildStatus(prNumber, GITHUB_TOKEN, commitSha, postBody)
	if err != nil {
		return err
	}

	return nil
}

func postBuildStatus(prNumber, GITHUB_TOKEN, commitSha string, body map[string]string) error {

	url := fmt.Sprintf("https://api.github.com/repos/GoogleCloudPlatform/magic-modules/statuses/%s", commitSha)

	_, err := requestCall(url, "POST", GITHUB_TOKEN, nil, body)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully posted community-checker build link to pull request %s", prNumber)

	return nil
}

func getPendingBuildId(projectId, commitSha string) (string, error) {
	COMMUNITY_CHECKER_TRIGGER, ok := os.LookupEnv("COMMUNITY_CHECKER_TRIGGER")
	if !ok {
		return "", fmt.Errorf("Did not provide COMMUNITY_CHECKER_TRIGGER environment variable")
	}

	ctx := context.Background()

	c, err := cloudbuild.NewService(ctx)
	if err != nil {
		return "", err
	}

	filter := fmt.Sprintf("trigger_id=%s AND status=PENDING", COMMUNITY_CHECKER_TRIGGER)
	// Builds will be sorted by createTime, descending order.
	// 10 should be enough to include the one needs auto approval
	pageSize := int64(10)

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

	c, err := cloudbuild.NewService(ctx)
	if err != nil {
		return err
	}

	name := fmt.Sprintf("projects/%s/builds/%s", projectId, buildId)

	approveBuildRequest := &cloudbuild.ApproveBuildRequest{
		ApprovalResult: &cloudbuild.ApprovalResult{
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

func addAwaitingApprovalLabel(prNumber, GITHUB_TOKEN string) error {
	url := fmt.Sprintf("https://api.github.com/repos/GoogleCloudPlatform/magic-modules/issues/%s/labels", prNumber)

	body := map[string][]string{
		"labels": []string{"awaiting-approval"},
	}
	_, err := requestCall(url, "POST", GITHUB_TOKEN, nil, body)

	if err != nil {
		return fmt.Errorf("Failed to add awaiting approval label: %s", err)
	}

	return nil

}

func removeAwaitingApprovalLabel(prNumber, GITHUB_TOKEN string) error {
	url := fmt.Sprintf("https://api.github.com/repos/GoogleCloudPlatform/magic-modules/issues/%s/labels/awaiting-approval", prNumber)
	_, err := requestCall(url, "DELETE", GITHUB_TOKEN, nil, nil)

	if err != nil {
		return fmt.Errorf("Failed to remove awaiting approval label: %s", err)
	}

	return nil
}
