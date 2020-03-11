package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	config := meta.(*Config)
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

	output, err := config.clientCompute.Instances.GetSerialPortOutput(project, zone, d.Get("instance").(string)).Port(int64(d.Get("port").(int))).Do()
	if err != nil {
		return err
	}
	contents := output.Contents
	for output.Contents != "" {
		output, err = config.clientCompute.Instances.GetSerialPortOutput(project, zone, d.Get("instance").(string)).Port(int64(d.Get("port").(int))).Start(output.Next).Do()
		if err != nil {
			return err
		}
		contents += output.Contents
	}

	d.Set("contents", contents)
	d.SetId(output.SelfLink)
	return nil
}
