// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package alloydb

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceAlloydbDatabaseCluster() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceAlloydbCluster().Schema)
	// Set custom fields
	dsScema_cluster_id := map[string]*schema.Schema{
		"project": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: `Project ID of the project.`,
		},
		"location": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: `The canonical ID for the location. For example: "us-east1".`,
		},
	}
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "cluster_id")

	// Set 'Required' schema elements

	dsSchema_m := tpgresource.MergeSchemas(dsScema_cluster_id, dsSchema)

	return &schema.Resource{
		Read:   dataSourceAlloydbDatabaseInstanceRead,
		Schema: dsSchema_m,
	}
}

func dataSourceAlloydbDatabaseClusterRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	// Get feilds for ID for setting cluster filed in resource
	cluster_id := d.Get("cluster_id").(string)

	location, err := tpgresource.GetLocation(d, config)
	if err != nil {
		return err
	}
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/clusters/{{cluster_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	// Setting cluster field
	d.Set("cluster", fmt.Sprintf("projects/%s/locations/%s/clusters/%s", project, location, cluster_id))

	err = resourceAlloydbClusterRead(d, meta)
	if err != nil {
		return err
	}

	if err := tpgresource.SetDataSourceLabels(d); err != nil {
		return err
	}

	if d.Id() == "" {
		return fmt.Errorf("%s not found", id)
	}
	return nil
}
