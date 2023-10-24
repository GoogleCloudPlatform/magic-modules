package cloudbuild

import (
	"context"
	"encoding/json"
	"fmt"
	"magician/bigquery"
	"os"
	"time"

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

	fmt.Println("Starting Cloud Build Run: ", triggerId)

	operation, err := c.Projects.Triggers.Run(projectId, triggerId, repoSource).Do()
	if err != nil {
		return fmt.Errorf("failed to create Cloud Build run: %s", err)
	}

	currentBuildId := os.Getenv("BUILD_ID")
	triggeredBuildId, err := getTriggeredCloudBuildId(c, operation)
	if err != nil {
		fmt.Println("Failed to get build id from triggered cloud build operation: ", err)
	} else if currentBuildId != "" && triggeredBuildId != "" {
		bigquery.InsertBuildMapping(currentBuildId, triggeredBuildId)
	}

	fmt.Println("Started Cloud Build Run: ", triggerId)
	return nil
}

func getTriggeredCloudBuildId(c *cloudbuild.Service, operation *cloudbuild.Operation) (string, error) {
	maxRetries := 30 // this gives a maximum waiting time of 30 * 5 seconds = 150 seconds
	retryCount := 0
	var err error

	for {
		if operation.Done || retryCount >= maxRetries {
			break
		}
		time.Sleep(time.Second * 5)
		operation, err = c.Operations.Get(operation.Name).Do()
		if err != nil {
			return "", fmt.Errorf("failed to wait for operation to finish: %s", err)
		}
		retryCount++
	}

	if operation.Done && operation.Response != nil {
		fmt.Println("got response from API \n", operation.Response)
		responseBytes, err := operation.Response.MarshalJSON()
		if err != nil {
			return "", fmt.Errorf("failed to deserialize operation resonse: %s", err)
		}

		var buildData map[string]interface{}
		err = json.Unmarshal(responseBytes, &buildData)
		if err != nil {
			return "", fmt.Errorf("failed to deserialize operation resonse: %s", err)
		}

		if buildId, exists := buildData["id"]; exists {
			fmt.Println("BUILD_ID: ", buildId)
			return buildId.(string), nil
		}
	}

	return "", fmt.Errorf("failed to extract buildId: %s", err)
}
