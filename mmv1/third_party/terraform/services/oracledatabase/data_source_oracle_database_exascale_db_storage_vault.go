package oracledatabase

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceOracleDatabaseExascaleDbStorageVault() *schema.Resource {
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceOracleDatabaseExascaleDbStorageVault().Schema)
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "location", "exascale_db_storage_vault_id")
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project")
	return &schema.Resource{
		Read:   dataSourceOracleDatabaseExascaleDbStorageVaultRead,
		Schema: dsSchema,
	}
}

func dataSourceOracleDatabaseExascaleDbStorageVaultRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/exascaleDbStorageVaults/{{exascale_db_storage_vault_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	err = resourceOracleDatabaseExascaleDbStorageVaultRead(d, meta)
	if err != nil {
		return err
	}
	d.SetId(id)

	return nil
}

func init() {
	registry.Schema{
		Name:        "google_oracle_database_exascale_db_storage_vault",
		ProductName: "oracledatabase",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceOracleDatabaseExascaleDbStorageVault(),
	}.Register()
}
