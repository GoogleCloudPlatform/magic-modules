package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleComputeInstanceSerialPort() *schema.Resource {
	return &schema.Resource{
		Read: computeInstanceSerialPortRead,
		Schema: map[string]*schema.Schema{
			"port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"instance": {
				Type:     schema.TypeString,
				Required: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"contents": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func computeInstanceSerialPortRead(d *schema.ResourceData, meta interface{}) error {
	var m providerMeta

	err := d.GetProviderMeta(&m)
	if err != nil {
		return err
	}
	config := meta.(*Config)
	config.clientCompute.UserAgent = fmt.Sprintf("%s %s", config.clientCompute.UserAgent, m.ModuleKey)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	d.Set("project", project)
	zone, err := getZone(d, config)
	if err != nil {
		return err
	}
	d.Set("zone", zone)

	port := int64(d.Get("port").(int))
	output, err := config.clientCompute.Instances.GetSerialPortOutput(project, zone, d.Get("instance").(string)).Port(port).Do()
	if err != nil {
		return err
	}

	d.Set("contents", output.Contents)
	d.SetId(output.SelfLink)
	return nil
}
