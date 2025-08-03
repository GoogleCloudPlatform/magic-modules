package redis

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleRedisInstance() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceRedisInstance().Schema)

	// Set 'Required' schema elements
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "name")

	// Set 'Optional' schema elements
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project", "region")

	return &schema.Resource{
		Read:   dataSourceGoogleRedisInstanceRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleRedisInstanceRead(d *schema.ResourceData, meta interface{}) error {
	id, err := tpgresource.ReplaceVars(d, meta.(*transport_tpg.Config), "projects/{{project}}/locations/{{region}}/instances/{{name}}")
	if err != nil {
		return err
	}
	d.SetId(id)

	err = resourceRedisInstanceRead(d, meta)
	if err != nil {
		return err
	}

	if err := tpgresource.SetDataSourceLabels(d); err != nil {
		return err
	}
	// added to resolve a null value for reserved_ip_range. This was not getting populated due to the addtion of ignore_read
	if err := SetDataSourceReservedIpRange(d); err != nil {
		return err
	}

	if d.Id() == "" {
		return fmt.Errorf("%s not found", id)
	}
	return nil
}

func SetDataSourceReservedIpRange(d *schema.ResourceData) error {
	effectiveReservedIpRange := d.Get("effective_reserved_ip_range")
	if effectiveReservedIpRange == nil {
		return nil
	}

	if err := d.Set("reserved_ip_range", effectiveReservedIpRange); err != nil {
		return fmt.Errorf("Error setting reserved_ip_range in data source: %s", err)
	}

	return nil
}
