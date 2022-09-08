package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type cloudFunction2Id struct {
	Project  string
	Location string
	Name     string
}

func (s *cloudFunction2Id) cloudFunction2Id() string {
	return fmt.Sprintf("projects/%s/locations/%s/functions/%s", s.Project, s.Location, s.Name)
}

func dataSourceGoogleCloudFunctions2Function() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := datasourceSchemaFromResourceSchema(resourceCloudfunctions2function().Schema)

	// Set 'Required' schema elements
	addRequiredFieldsToSchema(dsSchema, "name", "location")

	// Set 'Optional' schema elements
	addOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceGoogleCloudFunctions2FunctionRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleCloudFunctions2FunctionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	cloudFuncId := &cloudFunction2Id{
		Project:  project,
		Location: d.Get("location").(string),
		Name:     d.Get("name").(string),
	}

	d.SetId(cloudFuncId.cloudFunction2Id())

	err = resourceCloudfunctions2functionRead(d, meta)
	if err != nil {
		return err
	}

	return nil
}
