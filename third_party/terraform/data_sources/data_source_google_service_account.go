package google

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleServiceAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleServiceAccountRead,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateRFC1035Name(6, 30),
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"unique_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGoogleServiceAccountRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceAccountName, err := serviceAccountFQN(d.Get("account_id").(string), d, config)
	if err != nil {
		return err
	}

	sa, err := config.clientIAM.Projects.ServiceAccounts.Get(serviceAccountName).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Service Account %q", serviceAccountName))
	}

	d.SetId(sa.Name)
	if err := d.Set("email", sa.Email); err != nil {
		return fmt.Errorf("Error reading email: %s", err)
	}
	if err := d.Set("unique_id", sa.UniqueId); err != nil {
		return fmt.Errorf("Error reading unique_id: %s", err)
	}
	if err := d.Set("project", sa.ProjectId); err != nil {
		return fmt.Errorf("Error reading project: %s", err)
	}
	if err := d.Set("account_id", strings.Split(sa.Email, "@")[0]); err != nil {
		return fmt.Errorf("Error reading account_id: %s", err)
	}
	if err := d.Set("name", sa.Name); err != nil {
		return fmt.Errorf("Error reading name: %s", err)
	}
	if err := d.Set("display_name", sa.DisplayName); err != nil {
		return fmt.Errorf("Error reading display_name: %s", err)
	}

	return nil
}
