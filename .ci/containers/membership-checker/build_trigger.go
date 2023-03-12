package main

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/api/cloudbuild/v1"
)

func triggerMMPresubmitRuns(projectId, repoName, commitSha, branchName string) error {
	presubmitTriggerId, ok := os.LookupEnv("PRESUBMIT_TRIGGER")
	if !ok {
		return fmt.Errorf("Did not provide PRESUBMIT_TRIGGER environment variable")
	}

	rakeTestTriggerId, ok := os.LookupEnv("RAKE_TESTS_TRIGGER")
	if !ok {
		return fmt.Errorf("Did not provide RAKE_TESTS_TRIGGER environment variable")
	}

	err := triggerCloudBuildRun(projectId, presubmitTriggerId, repoName, commitSha, branchName)
	if err != nil {
		return err
	}

	err = triggerCloudBuildRun(projectId, rakeTestTriggerId, repoName, commitSha, branchName)
	if err != nil {
		return err
	}

	return nil
}

func triggerCloudBuildRun(projectId, triggerId, repoName, commitSha, branchName string) error {
	ctx := context.Background()
	c, err := cloudbuild.NewService(ctx)
	if err != nil {
		return fmt.Errorf("Failed to create Cloud Build service client: %s", err)
	}

	repoSource := &cloudbuild.RepoSource{
		ProjectId:     projectId,
		RepoName:      repoName,
		CommitSha:     commitSha,
		Substitutions: map[string]string{"BRANCH_NAME": branchName},
	}

	_, err = c.Projects.Triggers.Run(projectId, triggerId, repoSource).Do()
	if err != nil {
		return fmt.Errorf("Failed to create Cloud Build run: %s", err)
	}

	fmt.Println("Started Cloud Build Run: ", triggerId)
	return nil
}
