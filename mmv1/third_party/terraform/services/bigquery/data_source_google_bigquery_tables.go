// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package bigquery

import (
  "context"
	"fmt"
  "log"
  "time"

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
            Type:        schema.TypeList,
            Computed:    true,
            Elem:        &schema.Schema{Type: schema.TypeString},
            Description: "A list of table names in the dataset.",
        },
    }
  
	return &schema.Resource{
		Read: DataSourceGoogleBigQueryTablesRead,
		Schema: dsSchema, 
  } 
}

func DataSourceGoogleBigQueryTablesRead(d *schema.ResourceData, meta interface{}) error {

  ctx := context.Background()
  config := meta.(*transport_tpg.Config)

  dataset_id := d.Get("dataset_id").(string)

  project, err := tpgresource.GetProject(d, config)

  if err != nil {
      return fmt.Errorf("Error fetching project: %s", err)
  }

  bigquery_service, err := bq.NewService(ctx)

  if err != nil {
      return fmt.Errorf("Error creating BigQuery service: %s", err)
  }

  tables, err := bq.NewTablesService(bigquery_service).List(project, dataset_id).Do()

  if err != nil { 
    return fmt.Errorf("Error listing tables: %s", err)
  }

  var table_names []interface{}
  for _, table := range tables.Tables { 
    log.Printf("[INFO] Found BigQuery table: %s", table.TableReference.TableId)
    table_name := table.TableReference.TableId
    table_names = append(table_names, table_name)
  }

  if err := d.Set("tables", table_names); err != nil {
      log.Printf("[ERROR] Failed to set 'tables' attribute: %s", err)
      return fmt.Errorf("error setting 'tables' attribute: %w", err)
  }

  d.SetId(time.Now().UTC().String())

  return nil
}
