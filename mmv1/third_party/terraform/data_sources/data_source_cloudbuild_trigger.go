package google

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceresourceCloudBuildTrigger() *schema.Resource {
	dsSchema := datasourceSchemaFromResourceSchema(resourceCloudBuildTrigger().Schema)
	addRequiredFieldsToSchema(dsSchema, "trigger_id")
	addOptionalFieldsToSchema(dsSchema, "project")
	dsSchema["location"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
		Description: `The [Cloud Build location](https://cloud.google.com/build/docs/locations) for the trigger.
	If not specified, "global" is used.`,
		Default: "global",
	}
	return &schema.Resource{
		Read:   dataSourceGoogleCloudBuildTriggerRead,
		Schema: dsSchema,
	}

}

func dataSourceGoogleCloudBuildTriggerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	id, err := replaceVars(d, config, "projects/{{project}}/locations/{{location}}/triggers/{{trigger_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	id = strings.ReplaceAll(id, "/locations//", "/")
	id = strings.ReplaceAll(id, "/locations/global/", "/")
	d.SetId(id)
	return resourceCloudBuildTriggerRead(d, meta)
}