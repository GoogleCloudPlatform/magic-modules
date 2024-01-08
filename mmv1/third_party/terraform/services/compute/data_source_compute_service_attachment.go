package compute

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleComputeAttachment() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeAttachmentRead,

		Schema: map[string]*schema.Schema{

			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGoogleComputeAttachmentRead(d *schema.ResourceData, meta interface{}) error {
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
	region := d.Get("region").(string)
	id := fmt.Sprintf("projects/%s/regions/%s/serviceAttachments/%s", project, region, name)

	attachment, err := config.NewComputeClient(userAgent).ServiceAttachments.Get(project, region, name).Do()
	if err != nil {
		return transport_tpg.HandleDataSourceNotFoundError(err, d, fmt.Sprintf("Service Attachment Not Found : %s", name), id)
	}
	if err := d.Set("self_link", attachment.SelfLink); err != nil {
		return fmt.Errorf("Error setting self_link: %s", err)
	}
	if err := d.Set("Description", attachment.Description); err != nil {
		return fmt.Errorf("Error setting description: %s", err)
	}

	d.SetId(id)
	return nil
}
