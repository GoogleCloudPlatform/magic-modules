// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package lustre

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceLustreInstance() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceLustreInstance().Schema)

	// Set custom fields
	dsScema_custom := map[string]*schema.Schema{
		"location": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: `The canonical ID for the region. For example: "us-east1".`,
		},
	}

	// Set 'Required' schema elements from resource
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "instance_id")

	// Set 'Required' schema elements from resource
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project")

	// Set 'Required' schema elements
	dsSchema_m := tpgresource.MergeSchemas(dsScema_custom, dsSchema)

	return &schema.Resource{
		Read:   dataSourceLustreInstanceRead,
		Schema: dsSchema_m,
	}
}

func dataSourceLustreInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	// Get feilds for setting cluster field in resource
	instance_id := d.Get("instance_id").(string)

	location, err := tpgresource.GetRegion(d, config)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	// Set the ID
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/instances/{{instance_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	// Setting instance field as this is set as a required field in instance resource
	d.Set("instance_id", fmt.Sprintf("projects/%s/locations/%s/instances/%s", project, location, instance_id))

	err = resourceLustreInstanceRead(d, meta)
	if err != nil {
		return err
	}

	if err := tpgresource.SetDataSourceLabels(d); err != nil {
		return err
	}

	if d.Id() == "" {
		return fmt.Errorf("%s not found", d.Id())
	}

	return nil
}
