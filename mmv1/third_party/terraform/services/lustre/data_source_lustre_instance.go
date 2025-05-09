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

	dsScema_zone := map[string]*schema.Schema{
		"zone": {
			Type:     schema.TypeString,
			Optional: true,
			Description: `The ID of the zone in which the resource belongs. If it is not provided, the provider zone is used.,
`,
		},
	}

	// Set 'Required' schema elements from resource
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "instance_id")

	// Set 'Optional' schema elements from resource
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project")

	// Set 'Required' schema elements
	dsSchema_m := tpgresource.MergeSchemas(dsScema_zone, dsSchema)

	return &schema.Resource{
		Read:   dataSourceLustreInstanceRead,
		Schema: dsSchema_m,
	}
}

func dataSourceLustreInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	location, err := tpgresource.GetZone(d, config)
	if err != nil {
		return err
	}

	// Set the ID
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/instances/{{instance_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	// Setting location field for url_param_only field
	d.Set("location", location)

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
