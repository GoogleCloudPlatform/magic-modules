package eventarc

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceEventarcTrigger() *schema.Resource {
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceEventarcTrigger().Schema)

	// Set 'Required' schema elements
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "name")

	// Set 'Optional' schema elements
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "location", "project")

	return &schema.Resource{
		Read:   dataSourceEventarcTriggerRead,
		Schema: dsSchema,
	}
}

func dataSourceEventarcTriggerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	location, err := tpgresource.GetLocation(d, config)
	if err != nil {
		return err
	}
	if err := d.Set("location", location); err != nil {
		return fmt.Errorf("Error setting location: %s", err)
	}

	if v, ok := d.GetOk("name"); ok {
		if err := d.Set("name", tpgresource.GetResourceNameFromSelfLink(v.(string))); err != nil {
			return fmt.Errorf("Error setting name: %s", err)
		}
	}

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/triggers/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = resourceEventarcTriggerRead(d, meta)
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
		Name:        "google_eventarc_trigger",
		ProductName: "eventarc",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceEventarcTrigger(),
	}.Register()
}
