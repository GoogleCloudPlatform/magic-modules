package compute

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleComputeNetworkEndpointGroup() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceComputeNetworkEndpointGroup().Schema)

	// Set 'Optional' schema elements
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "name")
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "zone")
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project")
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "self_link")

	return &schema.Resource{
		Read:   dataSourceComputeNetworkEndpointGroupRead,
		Schema: dsSchema,
	}
}

func dataSourceComputeNetworkEndpointGroupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	id := ""
	if name, ok := d.GetOk("name"); ok {
		project, err := tpgresource.GetProject(d, config)
		if err != nil {
			return err
		}
		zone, err := tpgresource.GetZone(d, config)
		if err != nil {
			return err
		}
		id = fmt.Sprintf("projects/%s/zones/%s/networkEndpointGroups/%s", project, zone, name.(string))
		d.SetId(id)
	} else if selfLink, ok := d.GetOk("self_link"); ok {
		parsed, err := tpgresource.ParseNetworkEndpointGroupFieldValue(selfLink.(string), d, config)
		if err != nil {
			return err
		}
		if err := d.Set("name", parsed.Name); err != nil {
			return fmt.Errorf("Error setting name: %s", err)
		}
		if err := d.Set("zone", parsed.Zone); err != nil {
			return fmt.Errorf("Error setting zone: %s", err)
		}
		if err := d.Set("project", parsed.Project); err != nil {
			return fmt.Errorf("Error setting project: %s", err)
		}
		id = fmt.Sprintf("projects/%s/zones/%s/networkEndpointGroups/%s", parsed.Project, parsed.Zone, parsed.Name)
		d.SetId(id)
	} else {
		return errors.New("Must provide either `self_link` or `zone/name`")
	}

	err := resourceComputeNetworkEndpointGroupRead(d, meta)
	if err != nil {
		return err
	}

	if d.Id() == "" {
		return fmt.Errorf("%s not found", id)
	}

	return nil
}
