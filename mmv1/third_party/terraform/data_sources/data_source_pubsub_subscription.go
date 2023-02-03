package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGooglePubsubSubscription() *schema.Resource {

	dsSchema := DatasourceSchemaFromResourceSchema(resourcePubsubSubscription().Schema)
	AddRequiredFieldsToSchema(dsSchema, "name")
	AddOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceGooglePubsubSubscriptionRead,
		Schema: dsSchema,
	}
}

func dataSourceGooglePubsubSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	id, err := ReplaceVars(d, config, "projects/{{project}}/subscriptions/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)
	return resourcePubsubSubscriptionRead(d, meta)
}
