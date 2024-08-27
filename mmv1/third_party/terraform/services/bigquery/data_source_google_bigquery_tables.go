// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package bigquery

import (
  "context"
	"fmt"

  "google.golang.org/api/iterator"
	bq "google.golang.org/api/bigquery/v2"


	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)


func DataSourceGoogleBigQueryTables() *schema.Resource {

  dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceBigQueryTable().Schema)
  
  tpgresource.AddRequiredFieldsToSchema(dsSchema, "dataset_id")
  tpgresource.AddOptionalFieldsToSchema(dsSchema, "project_id")

  dsSchema["tables"] = &schema.Schema{
      Type:     schema.TypeList,
      Computed: true,
      Elem:     &schema.Schema{Type: schema.TypeString},
      Description: "A list of table names in the dataset.",
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

  tablesService := bq.NewTablesService(bigquery_service)
  listCall := tablesService.List(project, dataset_id)
  tables, err := listCall.Do()

  if err != nil { 
    return fmt.Errorf("Error listing tables: %s", err)
  }

  // Prepare the list of table names
  var table_names []interface{}
  for _, table := range tables.Tables { 
    table_names = append(table_names, table.TableReference.TableId) 
  }

  if err := d.Set("tables", table_names); err != nil {
      return fmt.Errorf("Error setting tables in state: %w", err)
  }

  return nil
}
