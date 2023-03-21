package google

import (
	"fmt"

	alloydb "cloud.google.com/go/alloydb/apiv1"
	gax "github.com/googleapis/gax-go/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/iterator"
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
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: `Service-specific metadata. For example the available capacity at the given location.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlloydbLocationsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project: %s", err)
	}
	billingProject := project
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	var listLocationsIterator *alloydb.LocationIterator
	locsReq := new(locationpb.ListLocationsRequest)
	err = nil
	err = retryTime(func() error {
		locationsClient := config.NewAlloydbClient()
		//set call options on locationsClient
		listLocationsIterator = locationsClient.ListLocations(config.context, locsReq, gax.WithPath(fmt.Sprintf("v1/projects/%s/locations", billingProject)))
		return nil
	}, 5)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Locations %q", d.Id()))
	}

	var locations []map[string]interface{}
	for {
		loc, err := listLocationsIterator.Next()
		if err == iterator.Done {
			break
		}
		locationDetails := make(map[string]interface{})
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
		if loc.Metadata != nil {
			locationDetails["metadata"] = fmt.Sprintf("%v", loc.Metadata)
		}
		locations = append(locations, locationDetails)
	}

	if err := d.Set("locations", locations); err != nil {
		return fmt.Errorf("Error setting locations: %s", err)
	}
	d.SetId(fmt.Sprintf("projects/%s/locations", project))
	return nil
}
