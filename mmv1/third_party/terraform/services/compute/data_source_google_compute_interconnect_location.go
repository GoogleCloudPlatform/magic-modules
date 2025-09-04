package compute

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleComputeInterconnectLocation() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeInterconnectLocationRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"peeringdb_facility_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"facility_provider": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"facility_provider_facility_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"continent": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"city": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"supports_pzs": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"available_features": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"available_link_types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceGoogleComputeInterconnectLocationRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)

	location, err := config.NewComputeClient(userAgent).InterconnectLocations.Get(project, name).Do()
	if err != nil {
		return transport_tpg.HandleDataSourceNotFoundError(err, d, fmt.Sprintf("Interconnect Location %q not found", name), "")
	}

	d.SetId(location.Name)
	d.Set("project", project)
	d.Set("description", location.Description)
	d.Set("self_link", location.SelfLink)
	d.Set("peeringdb_facility_id", location.PeeringdbFacilityId)
	d.Set("address", location.Address)
	d.Set("facility_provider", location.FacilityProvider)
	d.Set("facility_provider_facility_id", location.FacilityProviderFacilityId)
	d.Set("status", location.Status)
	d.Set("continent", location.Continent)
	d.Set("city", location.City)
	d.Set("availability_zone", location.AvailabilityZone)
	d.Set("supports_pzs", location.SupportsPzs)
	d.Set("available_features", location.AvailableFeatures)
	d.Set("available_link_types", location.AvailableLinkTypes)

	return nil
}
