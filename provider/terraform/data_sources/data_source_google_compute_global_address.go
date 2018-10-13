package google

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGoogleComputeGlobalAddress() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeGlobalAddressRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"self_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func dataSourceGoogleComputeGlobalAddressRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	address, err := config.clientCompute.GlobalAddresses.Get(project, name).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Global Address Not Found : %s", name))
	}

	d.Set("address", address.Address)
	d.Set("status", address.Status)
	d.Set("self_link", address.SelfLink)
	d.Set("project", project)

	d.SetId(strconv.FormatUint(uint64(address.Id), 10))
	return nil
}
