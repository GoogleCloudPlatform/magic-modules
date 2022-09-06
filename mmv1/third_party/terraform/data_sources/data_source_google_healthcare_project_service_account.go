package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleHealthcareProjectServiceAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleHealthcareProjectServiceAccountRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
			"user_project": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"email_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGoogleHealthcareProjectServiceAccountRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	rmProject, err := config.NewResourceManagerClient(userAgent).Projects.Get(project).Do()
	if err != nil {
		return handleNotFoundError(err, d, "Project not found")
	}
	projectNumber := rmProject.ProjectNumber

	serviceAccountEmail := fmt.Sprintf("service-%v@gcp-sa-healthcare.iam.gserviceaccount.com", projectNumber)

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("email_address", serviceAccountEmail); err != nil {
		return fmt.Errorf("Error setting email_address: %s", err)
	}

	d.SetId(serviceAccountEmail)

	return nil
}
