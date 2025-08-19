package alloydb

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceAlloydbLocations() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceAlloydbLocationsRead,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Project ID of the project.`,
			},
			"locations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: `Resource name for the location, which may vary between implementations. For example: "projects/example-project/locations/us-east1`,
						},
						"location_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: `The canonical id for this location. For example: "us-east1".`,
						},
						"display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: `The friendly name for this location, typically a nearby city name. For example, "Tokyo".`,
						},
						"labels": {
							Type:        schema.TypeMap,
							Computed:    true,
							Optional:    true,
							Description: `Cross-service attributes for the location. For example {"cloud.googleapis.com/region": "us-east1"}`,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"metadata": {
							Type:        schema.TypeMap,
							Computed:    true,
							Optional:    true,
							Description: `Service-specific metadata. For example the available capacity at the given location.`,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlloydbLocationsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error ππfetching project: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{AlloydbBasePath}}projects/{{project}}/locations")
	if err != nil {
		return fmt.Errorf("Error setting api endpoint")
	}

	opts := transport_tpg.GetPaginatedItemsOptions{
		ResourceData:   d,
		Config:         config,
		BillingProject: &billingProject,
		UserAgent:      userAgent,
		URL:            url,
		ResourceToList: "locations",
		Params:         map[string]string{"filter": "name:projects/{{project}}/locations/*"},
	}
	listedLocations, err := transport_tpg.GetPaginatedItems(opts)
	if err != nil {
		return err
	}
	var locations []map[string]interface{}
	for _, loc := range listedLocations {
		locationDetails := make(map[string]interface{})
		if loc["name"] != nil {
			locationDetails["name"] = loc["name"].(string)
		}
		if loc["locationId"] != nil {
			locationDetails["location_id"] = loc["locationId"].(string)
		}
		if loc["displayName"] != nil {
			locationDetails["display_id"] = loc["displayName"].(string)
		}
		if loc["labels"] != nil {
			labels := make(map[string]string)
			for k, v := range loc["labels"].(map[string]interface{}) {
				labels[k] = v.(string)
			}
			locationDetails["labels"] = labels
		}
		if loc["metadata"] != nil {
			metadata := make(map[string]string)
			for k, v := range loc["metadata"].(map[interface{}]interface{}) {
				metadata[k.(string)] = v.(string)
			}
			locationDetails["metadata"] = metadata
		}
		locations = append(locations, locationDetails)
	}

	if err := d.Set("locations", locations); err != nil {
		return fmt.Errorf("Error setting locations: %s", err)
	}
	d.SetId(fmt.Sprintf("projects/%s/locations", project))
	return nil
}
