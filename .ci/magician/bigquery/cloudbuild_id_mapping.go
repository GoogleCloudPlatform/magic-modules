package bigquery

import (
	"context"

	"cloud.google.com/go/bigquery"
)

type buildMapping struct {
	ParentBuildID string `bigquery:"parent_build_id"`
	ChildBuildID  string `bigquery:"child_build_id"`
}

func InsertBuildMapping(parentID, childID string) error {
	ctx := context.Background()

	// Set up BigQuery client.
	client, err := bigquery.NewClient(ctx, "YOUR_PROJECT_ID")
	if err != nil {
		return err
	}
	defer client.Close()

	// Data to insert.
	entry := &buildMapping{
		ParentBuildID: parentID,
		ChildBuildID:  childID,
	}

	// Get table reference.
	table := client.Dataset("graphite-docker-images").Table("magic_modules.cloud_build_id_mapping")

	// Insert data into table.
	inserter := table.Inserter()
	if err := inserter.Put(ctx, entry); err != nil {
		return err
	}

	return nil
}
