package cloudbuild

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/api/cloudbuild/v1"
)

func (cb cloudBuild) TriggerMMPresubmitRuns(commitSha string, substitutions map[string]string) error {
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
	c, err := cloudbuild.NewService(ctx)
	if err != nil {
		return fmt.Errorf("failed to create Cloud Build service client: %s", err)
	}

	repoSource := &cloudbuild.RepoSource{
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
