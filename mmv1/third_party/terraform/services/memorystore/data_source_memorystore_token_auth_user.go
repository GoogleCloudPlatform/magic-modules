package memorystore

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceMemorystoreTokenAuthUser() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceMemorystoreTokenAuthUser().Schema)

	// Set 'Required' schema elements
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "instance", "user_id")
	// Set 'Optional' schema elements
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project", "location")

	return &schema.Resource{
		Read:   dataSourceMemorystoreTokenAuthUserRead,
		Schema: dsSchema,
	}
}

func dataSourceMemorystoreTokenAuthUserRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/instances/{{instance}}/tokenAuthUsers/{{user_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = resourceMemorystoreTokenAuthUserRead(d, meta)
	if err != nil {
		return err
	}

	if d.Id() == "" {
		return fmt.Errorf("%s not found", id)
	}
	return nil
}

func init() {
	registry.Schema{
		Name:        "google_memorystore_token_auth_user",
		ProductName: "memorystore",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceMemorystoreTokenAuthUser(),
	}.Register()
}
