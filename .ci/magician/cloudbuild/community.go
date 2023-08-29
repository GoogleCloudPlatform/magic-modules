package cloudbuild

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/api/cloudbuild/v1"
)

func (cb cloudBuild) ApproveCommunityChecker(prNumber, commitSha string) error {
	buildId, err := getPendingBuildId(PROJECT_ID, commitSha)
	if err != nil {
		return err
	}

	if buildId == "" {
		return fmt.Errorf("Failed to find pending build for PR %s", prNumber)
	}

	err = approveBuild(PROJECT_ID, buildId)
	if err != nil {
		return err
	}

	return nil
}

func (cb cloudBuild) GetAwaitingApprovalBuildLink(prNumber, commitSha string) (string, error) {
	buildId, err := getPendingBuildId(PROJECT_ID, commitSha)
	if err != nil {
		return "", err
	}

	if buildId == "" {
		return "", fmt.Errorf("failed to find pending build for PR %s", prNumber)
	}

	targetUrl := fmt.Sprintf("https://console.cloud.google.com/cloud-build/builds;region=global/%s?project=%s", buildId, PROJECT_ID)

	return targetUrl, nil
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
