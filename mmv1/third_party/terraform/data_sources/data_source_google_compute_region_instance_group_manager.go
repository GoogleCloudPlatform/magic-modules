package google

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleComputeRegionInstanceGroupManager() *schema.Resource {

	dsSchema := datasourceSchemaFromResourceSchema(resourceComputeRegionInstanceGroupManager().Schema)
	addRequiredFieldsToSchema(dsSchema, "name")
	addOptionalFieldsToSchema(dsSchema, "project", "region")
	dsSchema["wait_for_instances"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	}

	return &schema.Resource{
		Read:   dataSourceComputeRegionInstanceGroupManagerRead,
		Schema: dsSchema,
	}
}

func dataSourceComputeRegionInstanceGroupManagerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	if name, ok := d.GetOk("name"); ok {
		region, err := getRegion(d, config)
		if err != nil {
			return err
		}
		project, err := getProject(d, config)
		if err != nil {
			return err
		}
		d.SetId(fmt.Sprintf("projects/%s/regions/%s/instanceGroupManagers/%s", project, region, name.(string)))
	} else {
		return errors.New("Must provide either `resource/name`")
	}

	err := resourceComputeRegionInstanceGroupManagerRead(d, meta)

	if err != nil {
		return err
	}
	if d.Id() == "" {
		return errors.New("Instance Manager Group not found")
	}
	return nil
}
