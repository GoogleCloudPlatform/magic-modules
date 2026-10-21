package networkconnectivity

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceNetworkConnectivitySpoke() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceNetworkConnectivitySpoke().Schema)

	// Set 'Required' schema elements
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "name")

	// Set 'Optional' schema elements
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project", "location")

	return &schema.Resource{
		Read:   dataSourceNetworkConnectivitySpokeRead,
		Schema: dsSchema,
	}
}

func dataSourceNetworkConnectivitySpokeRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	location, err := tpgresource.GetLocation(d, config)
	if err != nil {
		return err
	}
	if err := d.Set("location", location); err != nil {
		return fmt.Errorf("Error setting location: %s", err)
	}

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/spokes/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = resourceNetworkConnectivitySpokeRead(d, meta)
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

func init() {
	registry.Schema{
		Name:        "google_network_connectivity_spoke",
		ProductName: "networkconnectivity",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceNetworkConnectivitySpoke(),
	}.Register()
}
