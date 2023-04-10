package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGameServicesGameServerDeploymentRollout() *schema.Resource {

	dsSchema := datasourceSchemaFromResourceSchema(ResourceGameServicesGameServerDeploymentRollout().Schema)
	addRequiredFieldsToSchema(dsSchema, "deployment_id")

	return &schema.Resource{
		Read:   dataSourceGameServicesGameServerDeploymentRolloutRead,
		Schema: dsSchema,
	}
}

func dataSourceGameServicesGameServerDeploymentRolloutRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	id, err := ReplaceVars(d, config, "projects/{{project}}/locations/global/gameServerDeployments/{{deployment_id}}/rollout")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}

	d.SetId(id)

	return resourceGameServicesGameServerDeploymentRolloutRead(d, meta)
}
