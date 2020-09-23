package google

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleContainerRepo() *schema.Resource {
	return &schema.Resource{
		Read: containerRegistryRepoRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"repository_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func containerRegistryRepoRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	region, ok := d.GetOk("region")
	escapedProject := strings.Replace(project, ":", "/", -1)
	if ok && region != nil && region != "" {
		if err := d.Set("repository_url", fmt.Sprintf("%s.gcr.io/%s", region, escapedProject)); err != nil {
			return fmt.Errorf("Error setting repository_url: %s", err)
		}
	} else {
		if err := d.Set("repository_url", fmt.Sprintf("gcr.io/%s", escapedProject)); err != nil {
			return fmt.Errorf("Error setting repository_url: %s", err)
		}
	}
	d.SetId(d.Get("repository_url").(string))
	return nil
}
