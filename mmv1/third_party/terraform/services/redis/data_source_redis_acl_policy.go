package redis

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceRedisAclPolicy() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceRedisAclPolicy().Schema)

	// Set 'Required' schema elements
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "acl_policy_id")
	// Set 'Optional' schema elements
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project", "location")

	return &schema.Resource{
		Read:   dataSourceRedisAclPolicyRead,
		Schema: dsSchema,
	}
}

func dataSourceRedisAclPolicyRead(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*transport_tpg.Config)

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	location, err := tpgresource.GetLocation(d, config)
	if err != nil {
		return err
	}

	aclPolicyId := d.Get("acl_policy_id").(string)

	id := fmt.Sprintf("projects/%s/locations/%s/aclPolicies/%s", project, location, aclPolicyId)
	d.SetId(id)

	// Setting location field, as this is set as a required field in instance resource to build the url
	d.Set("location", location)

	err = resourceRedisAclPolicyRead(d, meta)
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
		Name:        "google_redis_acl_policy",
		ProductName: "redis",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceRedisAclPolicy(),
	}.Register()
}
