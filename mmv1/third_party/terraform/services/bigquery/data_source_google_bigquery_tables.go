package bigquery

import (
	"context"
	"fmt"
	"log"

	bq "google.golang.org/api/bigquery/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleBigQueryTables() *schema.Resource {

	dsSchema := map[string]*schema.Schema{
		"dataset_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The ID of the dataset containing the tables.",
		},
		"project": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The ID of the project in which the dataset is located. If it is not provided, the provider project is used.",
		},
    "tables": {
			Type:         schema.TypeMap,
			Computed:     true,
			Elem: &schema.Schema{ 
        Type:     schema.TypeString,
        Optional: true,
			},
			Description: "A map of table names in the dataset.",
		},
	}

	return &schema.Resource{
		Read:   DataSourceGoogleBigQueryTablesRead,
		Schema: dsSchema,
	}
}

func DataSourceGoogleBigQueryTablesRead(d *schema.ResourceData, meta interface{}) error {

	ctx := context.Background()
	config := meta.(*transport_tpg.Config)

	datasetID := d.Get("dataset_id").(string)

	project, err := tpgresource.GetProject(d, config)

	if err != nil {
		return fmt.Errorf("Error fetching project: %s", err)
	}

  bigqueryService, err := bq.NewService(ctx)

  if err != nil {
		return fmt.Errorf("Error creating BigQuery service: %s", err)
	}

  tablesService := bq.NewTablesService(bigqueryService)

	tableMap := make(map[string]interface{})

	nextPageToken := ""
	for {
		listCall := tablesService.List(project, datasetID)
		if nextPageToken != "" {
			listCall.PageToken(nextPageToken)
		}

		tables, err := listCall.Do()
		if err != nil {
			return fmt.Errorf("Error listing tables: %s", err)
		}

		for _, table := range tables.Tables {
			tableName := table.TableReference.TableId
			log.Printf("[INFO] Found BigQuery table: %s", tableName)

			tableMap[tableName] = nil
		}

		if tables.NextPageToken == "" {
			break
		}

		nextPageToken = tables.NextPageToken
	}

	if err := d.Set("tables", tableMap); err != nil {
		return fmt.Errorf("error setting 'tables' attribute: %w", err)
	}

  id := fmt.Sprintf("projects/%s/datasets/%s/tables", project, datasetID)
  d.SetId(id)

	return nil
}
