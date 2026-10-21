package securesourcemanager

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceSecureSourceManagerRepository() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceSecureSourceManagerRepository().Schema)

	// Set 'Required' schema elements
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "repository_id")

	// Set 'Optional' schema elements
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "location", "project")

	return &schema.Resource{
		Read:   dataSourceSecureSourceManagerRepositoryRead,
		Schema: dsSchema,
	}
}

func dataSourceSecureSourceManagerRepositoryRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	location, err := tpgresource.GetLocation(d, config)
	if err != nil {
		return err
	}
	if err := d.Set("location", location); err != nil {
		return fmt.Errorf("Error setting location: %s", err)
	}

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/repositories/{{repository_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = resourceSecureSourceManagerRepositoryRead(d, meta)
	if err != nil {
		return err
	}

	if err := d.Set("deletion_policy", nil); err != nil {
		return fmt.Errorf("Error setting deletion_policy: %s", err)
	}

	if d.Id() == "" {
		return fmt.Errorf("%s not found", id)
	}
	return nil
}

func init() {
	registry.Schema{
		Name:        "google_secure_source_manager_repository",
		ProductName: "securesourcemanager",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceSecureSourceManagerRepository(),
	}.Register()
}
