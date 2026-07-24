package redis

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceRedisTokenAuthUser() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceRedisTokenAuthUser().Schema)

	// Set 'Required' schema elements
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "cluster", "user_id")

	return &schema.Resource{
		Read:   dataSourceRedisTokenAuthUserRead,
		Schema: dsSchema,
	}
}

func dataSourceRedisTokenAuthUserRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	id, err := tpgresource.ReplaceVars(d, config, "{{cluster}}/tokenAuthUsers/{{user_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = resourceRedisTokenAuthUserRead(d, meta)
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
		Name:        "google_redis_token_auth_user",
		ProductName: "redis",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceRedisTokenAuthUser(),
	}.Register()
}
