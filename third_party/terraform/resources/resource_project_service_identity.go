package google

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceProjectServiceIdentity() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectServiceIdentityCreate,
		Read:   resourceProjectServiceIdentityRead,
		Delete: resourceProjectServiceIdentityDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Read:   schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"service": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceProjectServiceIdentityCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	url, err := replaceVars(d, config, "{{ServiceUsageBasePath}}projects/{{project}}/services/{{service}}:generateServiceIdentity")
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	billingProject := project

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequestWithTimeout(config, "POST", billingProject, url, nil, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating Service Identity: %s", err)
	}

	err = serviceUsageOperationWaitTime(
		config, res, project, "Creating Service Identity",
		d.Timeout(schema.TimeoutCreate))

	if err != nil {
		return err
	}

	id, err := replaceVars(d, config, "projects/{{project}}/services/{{service}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating Service Identity %q: %#v", d.Id(), res)
	return nil
}

// There is no read endpoint for this API.
func resourceProjectServiceIdentityRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

// There is no delete endpoint for this API.
func resourceProjectServiceIdentityDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
