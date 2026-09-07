package compute

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleComputeInstanceSerialPort() *schema.Resource {
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
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	zone, err := tpgresource.GetZone(d, config)
	if err != nil {
		return err
	}
	if err := d.Set("zone", zone); err != nil {
		return fmt.Errorf("Error setting zone: %s", err)
	}

	port := d.Get("port").(int)
	url := fmt.Sprintf("%sprojects/%s/zones/%s/instances/%s/serialPort",
		transport_tpg.BaseUrl(Product, config), project, zone, d.Get("instance").(string))
	url, err = transport_tpg.AddQueryParams(url, map[string]string{"port": strconv.Itoa(port)})
	if err != nil {
		return err
	}
	output, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return err
	}

	if err := d.Set("contents", output["contents"]); err != nil {
		return fmt.Errorf("Error setting contents: %s", err)
	}
	selfLink, _ := output["selfLink"].(string)
	d.SetId(selfLink)
	return nil
}

func init() {
	registry.Schema{
		Name:        "google_compute_instance_serial_port",
		ProductName: "compute",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceGoogleComputeInstanceSerialPort(),
	}.Register()
}
