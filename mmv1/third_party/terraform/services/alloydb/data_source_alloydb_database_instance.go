// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package alloydb

import (
	"fmt"

	// "log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceAlloydbDatabaseInstance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlloydbDatabaseInstanceRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the alloydb cluster that the instance belongs to.'alloydb_cluster_id'`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the alloydb instance.'alloydb_instance_id'`,
			},
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
		},
	}
}

func dataSourceAlloydbDatabaseInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/clusters/{{cluster_id}}/instances/{{instance_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = resourceAlloydbInstanceRead(d, meta)
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
