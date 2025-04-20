package datastream

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleDatastreamPrivateConnection() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceDatastreamPrivateConnection().Schema)
	// Set 'Required' schema elements
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "private_connection_id")

	// Set 'Optional' schema elements
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "display_name", "vpc_peering_config", "location", "project")

	return &schema.Resource{
		Read:   dataSourceGoogleDatastreamPrivateConnectionRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleDatastreamPrivateConnectionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/privateConnections/{{private_connection_id}}")
	if err != nil {
		return err
	}

	err = resourceDatastreamPrivateConnectionRead(d, meta)
	if err != nil {
		return err
	}

	if err := tpgresource.SetDataSourceLabels(d); err != nil {
		return err
	}

	// Store the ID now
	if d.Id() == "" {
		return fmt.Errorf("%s not found", id)
	}
	return nil
}
