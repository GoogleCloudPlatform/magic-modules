package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	resourceManagerV2 "google.golang.org/api/cloudresourcemanager/v2"
)

func dataSourceGoogleActiveFolder() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleActiveFolderRead,

		Schema: map[string]*schema.Schema{
			"parent": {
				Type:     schema.TypeString,
				Required: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGoogleActiveFolderRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	var folderMatch *resourceManagerV2.Folder
	parent := d.Get("parent").(string)
	displayName := d.Get("display_name").(string)
	token := ""

	for paginate := true; paginate; {
		resp, err := config.NewResourceManagerV2Client(userAgent).Folders.List().Parent(parent).PageToken(token).Do()
		if err != nil {
			return fmt.Errorf("error reading folder list: %s", err)
		}

		for _, folder := range resp.Folders {
			if folder.DisplayName == displayName && folder.LifecycleState == "ACTIVE" {
				if folderMatch != nil {
					return fmt.Errorf("more than one matching folder found")
				}
				folderMatch = folder
			}
		}
		token = resp.NextPageToken
		paginate = token != ""
	}

	if folderMatch == nil {
		return fmt.Errorf("folder not found: %s", displayName)
	}

	d.SetId(folderMatch.Name)
	if err := d.Set("name", folderMatch.Name); err != nil {
		return fmt.Errorf("Error setting folder name: %s", err)
	}

	return nil
}
