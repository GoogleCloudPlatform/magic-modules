package compute

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleComputeReservation() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceComputeReservation().Schema)
	dsSchema_block_name := map[string]*schema.Schema{
		"block_names": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "List of all reservation block names in the parent reservation.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}

	// Set 'Required' schema elements
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "name")
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "zone")
	// Set `Optional` schema elements
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project")
	// Merge schemas
	dsSchema_m := tpgresource.MergeSchemas(dsSchema_block_name, dsSchema)

	return &schema.Resource{
		Read:   dataSourceGoogleComputeReservationRead,
		Schema: dsSchema_m,
	}
}

func dataSourceGoogleComputeReservationRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	err := resourceComputeReservationRead(d, meta)
	if err != nil {
		return err
	}
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	zone, err := tpgresource.GetZone(d, config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)

	// Fetch the list of all reservation blocks from this reservation
	listUrl := fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/zones/%s/reservations/%s/reservationBlocks?alt=json&maxResults=500", project, zone, name)

	listRes, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    listUrl,
		UserAgent: userAgent,
	})
	if err != nil {
		return fmt.Errorf("Error listing ReservationBlocks: %s", err)
	}

	blockNames := []string{}
	if listRes != nil {
		if items, ok := listRes["items"].([]interface{}); ok {
			for _, item := range items {
				if block, ok := item.(map[string]interface{}); ok {
					if blockName, ok := block["name"].(string); ok {
						blockNames = append(blockNames, blockName)
					}
				}
			}
		}
	}

	if err := d.Set("block_names", blockNames); err != nil {
		return fmt.Errorf("Error setting block_names: %s", err)
	}

	d.SetId(fmt.Sprintf("projects/%s/zones/%s/reservations/%s", project, zone, name))
	return nil
}
