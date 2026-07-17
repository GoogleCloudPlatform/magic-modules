package redis

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceRedisAuthToken() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceRedisAuthToken().Schema)

	// Set 'Required' schema elements
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "token_auth_user", "token_id")
	// Set 'Optional' schema elements
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project", "region")

	return &schema.Resource{
		Read:   dataSourceRedisAuthTokenRead,
		Schema: dsSchema,
	}
}

func dataSourceRedisAuthTokenRead(d *schema.ResourceData, meta interface{}) error {
	tokenAuthUser := d.Get("token_auth_user").(string)
	tokenId := d.Get("token_id").(string)
	id := fmt.Sprintf("%s/authTokens/%s", tokenAuthUser, tokenId)
	d.SetId(id)

	err := resourceRedisAuthTokenRead(d, meta)
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
		Name:        "google_redis_auth_token",
		ProductName: "redis",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceRedisAuthToken(),
	}.Register()
}
