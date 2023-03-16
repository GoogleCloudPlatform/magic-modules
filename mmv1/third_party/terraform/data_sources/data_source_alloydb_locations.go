package google

import (
	"fmt"

	alloydb "cloud.google.com/go/alloydb/apiv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	locationpb "google.golang.org/genproto/googleapis/cloud/location"
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
	config := meta.(*Config)
	//userAgent, err := generateUserAgentString(d, config.userAgent)
	// if err != nil {
	// 	return err
	// }
	//billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project: %s", err)
	}
	//billingProject := project

	// err == nil indicates that the billing_project value was found
	// if bp, err := getBillingProject(d, config); err == nil {
	// 	billingProject = bp
	// }

	var listLocationsIterator *alloydb.LocationIterator
	//pageToken := ""
	err = nil
	//locationsClientBasePath := removeBasePathVersion(removeBasePathVersion(config.AlloydbBasePath))
	var locsReq *locationpb.ListLocationsRequest
	locsReq.PageSize = 5
	err = retryTime(func() error {
		locationsClient := config.NewAlloydbClient()
		//set call options on locationsClient
		listLocationsIterator = locationsClient.ListLocations(config.context, locsReq)
		return nil
	}, 5)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Locations %q", d.Id()))
	}
	//var fetchedLocations []*locationpb.Location
	fmt.Println("[INFO] kanthara: %q", listLocationsIterator)

	var locations []map[string]interface{}
	for {
		//fetchedLocations := res["locations"].([]interface{})
		for i := 0; i < 4; i++ {
			loc, _ := listLocationsIterator.Next()
			locationDetails := make(map[string]interface{})
			//l := loc.(map[string]interface{})
			if loc.Name != "" {
				locationDetails["name"] = loc.Name
			}
			if loc.LocationId != "" {
				locationDetails["location_id"] = loc.LocationId
			}
			if loc.DisplayName != "" {
				locationDetails["display_id"] = loc.DisplayName
			}
			if loc.Labels != nil {
				labels := make(map[string]string)
				for k, v := range loc.Labels {
					labels[k] = v
				}
				locationDetails["labels"] = labels
			}
			// if loc.Metadata != nil {
			// 	metadata := make(map[string]string)
			// 	for k, v := range loc.Metadata.(map[interface{}]interface{}) {
			// 		metadata[k.(string)] = v.(string)
			// 	}
			// 	locationDetails["metadata"] = metadata
			// }
			locations = append(locations, locationDetails)
		}
		break
		// if listLocationsResponse.NextPageToken == nil || listLocationsResponse.NextPageToken.(string) == "" {
		// 	break
		// }
		// err = retryTime(func() error {
		// 	listLocationsResponse, err = config.NewLocationsClient(userAgent).Locations.List(project).PageToken(listLocationsResponse.NextPageToken).Do()
		// 	return err
		// }, 5)
		// if err != nil {
		// 	return handleNotFoundError(err, d, fmt.Sprintf("Locations %q", d.Id()))
		// }
	}

	if err := d.Set("locations", locations); err != nil {
		return fmt.Errorf("Error setting locations: %s", err)
	}
	d.SetId(fmt.Sprintf("projects/%s/locations", project))
	return nil
}
