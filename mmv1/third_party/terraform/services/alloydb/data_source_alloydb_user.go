package alloydb

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceAlloydbUser() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceAlloydbUser().Schema)

	// Set custom fields
	dsSchemaClusterId := map[string]*schema.Schema{
		"cluster_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: `The ID of the alloydb cluster that the user belongs to.`,
		},
		"project": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: `Project ID of the project.`,
		},
		"location": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: `The canonical ID for the location. For example: "us-east1".`,
		},
	}

	dsSchemaMerged := tpgresource.MergeSchemas(dsSchemaClusterId, dsSchema)

	// Set 'Required' schema elements
	tpgresource.AddRequiredFieldsToSchema(dsSchemaMerged, "cluster_id", "user_id")

	// Set 'Optional' schema elements
	tpgresource.AddOptionalFieldsToSchema(dsSchemaMerged, "location", "project")

	return &schema.Resource{
		Read:   dataSourceAlloydbUserRead,
		Schema: dsSchemaMerged,
	}
}

func dataSourceAlloydbUserRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	clusterId := d.Get("cluster_id").(string)

	location, err := tpgresource.GetLocation(d, config)
	if err != nil {
		return err
	}
	if err := d.Set("location", location); err != nil {
		return fmt.Errorf("Error setting location: %s", err)
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}

	// Setting cluster field as this is set as a required field in user resource
	cluster := fmt.Sprintf("projects/%s/locations/%s/clusters/%s", project, location, clusterId)
	if err := d.Set("cluster", cluster); err != nil {
		return fmt.Errorf("Error setting cluster: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/clusters/{{cluster_id}}/users/{{user_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = resourceAlloydbUserRead(d, meta)
	if err != nil {
		return err
	}

	if d.Id() == "" {
		return fmt.Errorf("%s not found", id)
	}

	return nil
}

func init() {
	registry.Schema{
		Name:        "google_alloydb_user",
		ProductName: "alloydb",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceAlloydbUser(),
	}.Register()
}
