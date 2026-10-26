package oracledatabase

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceOracleDatabaseOdbDnsZone() *schema.Resource {
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceOracleDatabaseOdbDnsZone().Schema)
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "location", "odbnetwork", "odb_dns_zone_id")
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project")
	return &schema.Resource{
		Read:   dataSourceOracleDatabaseOdbDnsZoneRead,
		Schema: dsSchema,
	}

}

func dataSourceOracleDatabaseOdbDnsZoneRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/odbNetworks/{{odbnetwork}}/odbDnsZones/{{odb_dns_zone_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}

	d.SetId(id)

	err = resourceOracleDatabaseOdbDnsZoneRead(d, meta)
	if err != nil {
		return err
	}

	if err := tpgresource.SetDataSourceLabels(d); err != nil {
		return err
	}

	return nil
}

func init() {
	registry.Schema{
		Name:        "google_oracle_database_odb_dns_zone",
		ProductName: "oracledatabase",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceOracleDatabaseOdbDnsZone(),
	}.Register()
}
