package compute

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// DataSourceGoogleComputeInterconnectGroups returns the schema of the google_compute_interconnect_groups data source.
func DataSourceGoogleComputeInterconnectGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeInterconnectGroupsRead,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The ID of the project whose interconnect groups are listed. If it is not provided, the provider project is used.`,
			},
			"interconnect_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The interconnect groups of the project.`,
				Elem:        &schema.Resource{Schema: dataSourceGoogleComputeInterconnectGroupsItemSchema()},
			},
		},
	}
}

// dataSourceGoogleComputeInterconnectGroupsItemSchema returns the schema of one group of the list, built from google_compute_interconnect_group.
func dataSourceGoogleComputeInterconnectGroupsItemSchema() map[string]*schema.Schema {
	groupSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceComputeInterconnectGroup().Schema)
	// These fields are left out: project repeats the argument, and deletion_policy is a setting of the resource.
	tpgresource.DeleteFieldsFromSchema(groupSchema, "project", "deletion_policy")
	// A nested block gets no id from the SDK, so the id of the resource is added by hand, as in google_kms_crypto_keys.
	groupSchema["id"] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: `An identifier for the group with format projects/{{project}}/global/interconnectGroups/{{name}}.`,
	}
	return groupSchema
}

// dataSourceGoogleComputeInterconnectGroupsRead lists every group of the project with the list call of the google_compute_interconnect_group list resource.
func dataSourceGoogleComputeInterconnectGroupsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	groupSchema := dataSourceGoogleComputeInterconnectGroupsItemSchema()
	groups := make([]map[string]interface{}, 0)
	// ListComputeInterconnectGroups follows the pages and hands over each group as the data of a google_compute_interconnect_group resource.
	err = ListComputeInterconnectGroups(config, project, func(rd *schema.ResourceData) error {
		group := make(map[string]interface{}, len(groupSchema))
		for key := range groupSchema {
			if key == "id" {
				group[key] = rd.Id()
			} else {
				group[key] = rd.Get(key)
			}
		}
		groups = append(groups, group)
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error retrieving list of interconnect groups: %w", err)
	}

	if err := d.Set("interconnect_groups", groups); err != nil {
		return fmt.Errorf("Error setting interconnect_groups: %w", err)
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %w", err)
	}

	d.SetId(fmt.Sprintf("projects/%s/global/interconnectGroups", project))

	return nil
}

func init() {
	registry.Schema{
		Name:        "google_compute_interconnect_groups",
		ProductName: "compute",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceGoogleComputeInterconnectGroups(),
	}.Register()
}
