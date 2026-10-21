package alloydb

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceAlloydbBackup() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceAlloydbBackup().Schema)

	// Set 'Required' schema elements
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "backup_id")

	// Set 'Optional' schema elements
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "location", "project")

	return &schema.Resource{
		Read:   dataSourceAlloydbBackupRead,
		Schema: dsSchema,
	}
}

func dataSourceAlloydbBackupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	location, err := tpgresource.GetLocation(d, config)
	if err != nil {
		return err
	}
	if err := d.Set("location", location); err != nil {
		return fmt.Errorf("Error setting location: %s", err)
	}

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/backups/{{backup_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = resourceAlloydbBackupRead(d, meta)
	if err != nil {
		return err
	}

	if err := tpgresource.SetDataSourceLabels(d); err != nil {
		return err
	}

	if err := tpgresource.SetDataSourceAnnotations(d); err != nil {
		return err
	}

	if d.Id() == "" {
		return fmt.Errorf("%s not found", id)
	}

	return nil
}

func init() {
	registry.Schema{
		Name:        "google_alloydb_backup",
		ProductName: "alloydb",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceAlloydbBackup(),
	}.Register()
}
