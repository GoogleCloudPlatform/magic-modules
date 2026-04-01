package compute

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceComputeSharedVpcHostProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeSharedVpcHostProjectCreate,
		Read:   resourceComputeSharedVpcHostProjectRead,
		Update:	resourceComputeSharedVpcHostProjectUpdate,
		Delete: resourceComputeSharedVpcHostProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(4 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
            tpgresource.DefaultProviderDeletionPolicy("DELETE"),
        ),

		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the project that will serve as a Shared VPC host project`,
			},
			//UDP schema start
			"deletion_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: `Whether Terraform will be prevented from destroying the instance. Defaults to "DELETE".
When a 'terraform destroy' or 'terraform apply' would delete the instance,
the command will fail if this field is set to "PREVENT" in Terraform state.
When set to "ABANDON", the command will remove the resource from Terraform
management without updating or deleting the resource in the API.
When set to "DELETE", deleting the resource is allowed.
`,
			},
			//UDP schema end
		},
		UseJSONNumber: true,
	}
}

func resourceComputeSharedVpcHostProjectCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	hostProject := d.Get("project").(string)
	op, err := config.NewComputeClient(userAgent).Projects.EnableXpnHost(hostProject).Do()
	if err != nil {
		return fmt.Errorf("Error enabling Shared VPC Host %q: %s", hostProject, err)
	}

	d.SetId(hostProject)

	err = ComputeOperationWaitTime(config, op, hostProject, "Enabling Shared VPC Host", userAgent, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		d.SetId("")
		return err
	}

	return nil
}

func resourceComputeSharedVpcHostProjectRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	hostProject := d.Id()

	project, err := config.NewComputeClient(userAgent).Projects.Get(hostProject).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Project data for project %q", hostProject))
	}

	if project.XpnProjectStatus != "HOST" {
		log.Printf("[WARN] Removing Shared VPC host resource %q because it's not enabled server-side", hostProject)
		d.SetId("")
	}

	if err := d.Set("project", hostProject); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}

    //UDP default read start
    // Explicitly set virtual fields to default values if unset
    if _, ok := d.GetOkExists("deletion_policy"); !ok {
        //prioritize config's value if present
        if config.DeletionPolicy != ""{
            if err := d.Set("deletion_policy", config.DeletionPolicy); err != nil {
                return fmt.Errorf("Error setting deletion_policy: %s", err)
            }
        }else{
            if err := d.Set("deletion_policy", "DELETE"); err != nil {
                return fmt.Errorf("Error setting deletion_policy: %s", err)
            }
        }
    }
    //UDP default read end
	return nil
}

//UDP update start
func resourceComputeSharedVpcHostProjectUpdate(d *schema.ResourceData, meta interface{}) error {
    // Only the root field "deletion_policy", "labels", "terraform_labels", and virtual fields are mutable
    return resourceComputeSharedVpcHostProjectRead(d, meta)
}
//UDP update end

func resourceComputeSharedVpcHostProjectDelete(d *schema.ResourceData, meta interface{}) error {
    //UDP pre-delete start
    if d.Get("deletion_policy").(string) == "PREVENT" {
        return fmt.Errorf("cannot destroy Shared VPC Host without setting deletion_policy=\"DELETE\" and running `terraform apply`")
    }
    if d.Get("deletion_policy").(string) == "ABANDON" {
        log.Printf("[DEBUG] deletion_policy set to \"ABANDON\", removing Shared VPC Host %q from Terraform state without deletion", d.Id())
        return nil
    }
    //UDP pre-delete end	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	hostProject := d.Get("project").(string)

	op, err := config.NewComputeClient(userAgent).Projects.DisableXpnHost(hostProject).Do()
	if err != nil {
		return fmt.Errorf("Error disabling Shared VPC Host %q: %s", hostProject, err)
	}

	err = ComputeOperationWaitTime(config, op, hostProject, "Disabling Shared VPC Host", userAgent, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
