package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleComputeRegionNetworkEndpointGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeRegionNetworkEndpointGroupRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"self_link": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func dataSourceGoogleComputeRegionNetworkEndpointGroupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	project, region, name, err := GetRegionalResourcePropertiesFromSelfLinkOrSchema(d, config)
	if err != nil {
		return err
	}

	networEndpointGroup, err := config.NewComputeClient(userAgent).RegionInstanceGroups.Get(
		project, region, name).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Region Network Endpoint Group %q", name))
	}

	d.SetId(fmt.Sprintf("projects/%s/regions/%s/networkEndpointGroups/%s", project, region, name))
	if err := d.Set("self_link", networEndpointGroup.SelfLink); err != nil {
		return fmt.Errorf("Error setting self_link: %s", err)
	}

	if err := d.Set("name", name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}

	if err := d.Set("region", region); err != nil {
		return fmt.Errorf("Error setting region: %s", err)
	}

	return nil
}
