package memorystore

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func DataSourceMemorystoreAuthToken() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceMemorystoreAuthToken().Schema)

	// Set 'Required' schema elements
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "token_auth_user", "token_id")

	return &schema.Resource{
		Read:   dataSourceMemorystoreAuthTokenRead,
		Schema: dsSchema,
	}
}

func dataSourceMemorystoreAuthTokenRead(d *schema.ResourceData, meta interface{}) error {
	tokenAuthUser := d.Get("token_auth_user").(string)
	tokenId := d.Get("token_id").(string)
	id := fmt.Sprintf("%s/authTokens/%s", tokenAuthUser, tokenId)
	d.SetId(id)

	err := resourceMemorystoreAuthTokenRead(d, meta)
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
		Name:        "google_memorystore_auth_token",
		ProductName: "memorystore",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceMemorystoreAuthToken(),
	}.Register()
}
