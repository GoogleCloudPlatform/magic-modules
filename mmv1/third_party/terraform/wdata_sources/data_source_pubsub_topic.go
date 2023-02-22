package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGooglePubsubTopic() *schema.Resource {

	dsSchema := DatasourceSchemaFromResourceSchema(ResourcePubsubTopic().Schema)
	AddRequiredFieldsToSchema(dsSchema, "name")
	AddOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceGooglePubsubTopicRead,
		Schema: dsSchema,
	}
}

func dataSourceGooglePubsubTopicRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	id, err := ReplaceVars(d, config, "projects/{{project}}/topics/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)
	return resourcePubsubTopicRead(d, meta)
}
