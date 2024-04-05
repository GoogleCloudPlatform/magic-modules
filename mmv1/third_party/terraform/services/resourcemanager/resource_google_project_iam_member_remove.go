package resourcemanager

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/cloudresourcemanager/v1"
)

func ResourceGoogleProjectIamMemberRemove() *schema.Resource {
	return &schema.Resource{
		Create: resourceGoogleProjectIamMemberRemoveCreate,
		Read:   resourceGoogleProjectIamMemberRemoveRead,
		Delete: resourceGoogleProjectIamMemberRemoveDelete,

		CustomizeDiff: customdiff.All(
			tpgresource.DefaultProviderProject,
		),

		Schema: map[string]*schema.Schema{
			"role": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The project that the service account will be created in. Defaults to the provider project configuration.`,
			},
			"member": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: `The Identity of the service account in the form 'serviceAccount:{email}'. This value is often used to refer to the service account in order to grant IAM permissions.`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceGoogleProjectIamMemberRemoveCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	
	project:= d.Get("project").(string)
	role := d.Get("role").(string)
	member:= d.Get("member").(string)

	iamPolicy, err := config.NewResourceManagerClient(config.UserAgent).Projects.GetIamPolicy(project,
		&cloudresourcemanager.GetIamPolicyRequest{
			Options: &cloudresourcemanager.GetPolicyOptions{
				RequestedPolicyVersion: tpgiamresource.IamPolicyVersion,
			},
		}).Do()
	for _, bind := range iamPolicy.Bindings {
		for _, existingMember := range bind.Members {
			if member == existingMember {
				if role == bind.Role {
					bind.Role = "role/viewer"
					updateRequest := &cloudresourcemanager.SetIamPolicyRequest{
						Policy:     iamPolicy,
						UpdateMask: "bindings",
					}
					_, err = config.NewResourceManagerClient(config.UserAgent).Projects.SetIamPolicy(project, updateRequest).Do()
					if err != nil {
						return fmt.Errorf("cannot update IAM binding policy on project %s: %v", project, err)
					}
				} else {
					return fmt.Errorf("Could not find Member %s with the corresponding role %s.", member, role)
				}
			}
		}
	}
	d.SetId(fmt.Sprintf("%s/%s", project, member))
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, d.Id())
	}

	return resourceGoogleProjectIamMemberRemoveRead(d, meta)
}

func resourceGoogleProjectIamMemberRemoveRead(d *schema.ResourceData, meta interface{}) error {
	/*
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	*/
	project := d.Get("project").(string)
	role:= d.Get("role").(string)
	member:= d.Get("member").(string)

	if err := d.Set("role", role); err != nil {
		return fmt.Errorf("Error setting role: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("member", member); err != nil {
		return fmt.Errorf("Error setting member: %s", err)
	}

	return nil
}

func resourceGoogleProjectIamMemberRemoveDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}
