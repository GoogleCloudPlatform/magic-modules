package resourcemanager

import (
	"strings"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/verify"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ProjectAssetType is the CAI asset type name for project.
const ProjectAssetType string = "cloudresourcemanager.googleapis.com/Project"

// ProjectSchemaName is the TF resource schema name for resourcemanager project.
const ProjectSchemaName string = "google_project"

func ParseFolderId(v interface{}) string {
	folderId := v.(string)
	if strings.HasPrefix(folderId, "folders/") {
		return folderId[8:]
	}
	return folderId
}

// ResourceGoogleProject returns a *schema.Resource that allows a customer
// to declare a Google Cloud Project resource.
func ResourceGoogleProject() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidateProjectID(),
				Description:  `The project ID. Changing this forces a new project to be created.`,
			},
			"deletion_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "PREVENT",
				Description: `The deletion policy for the Project. Setting PREVENT will protect the project against any destroy actions caused by a terraform apply or terraform destroy. Setting ABANDON allows the resource
				to be abandoned rather than deleted. Possible values are: "PREVENT", "ABANDON", "DELETE"`,
				ValidateFunc: validation.StringInSlice([]string{"PREVENT", "ABANDON", "DELETE"}, false),
			},
			"auto_create_network": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: `Create the 'default' network automatically.  Default true. If set to false, the default network will be deleted.  Note that, for quota purposes, you will still need to have 1 network slot available to create the project successfully, even if you set auto_create_network to false, since the network will exist momentarily.`,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: verify.ValidateProjectName(),
				Description:  `The display name of the project.`,
			},
			"org_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"folder_id"},
				Description:   `The numeric ID of the organization this project belongs to. Changing this forces a new project to be created.  Only one of org_id or folder_id may be specified. If the org_id is specified then the project is created at the top level. Changing this forces the project to be migrated to the newly specified organization.`,
			},
			"folder_id": {
				Type:          schema.TypeString,
				Optional:      true,
				StateFunc:     ParseFolderId,
				ConflictsWith: []string{"org_id"},
				Description:   `The numeric ID of the folder this project should be created under. Only one of org_id or folder_id may be specified. If the folder_id is specified, then the project is created under the specified folder. Changing this forces the project to be migrated to the newly specified folder.`,
			},
			"number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The numeric identifier of the project.`,
			},
			"billing_account": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The alphanumeric ID of the billing account this project belongs to. The user or service account performing this operation with Terraform must have Billing Account Administrator privileges (roles/billing.admin) in the organization. See Google Cloud Billing API Access Control for more details.`,
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `A set of key/value label pairs to assign to the project.
				
				**Note**: This field is non-authoritative, and will only manage the labels present in your configuration.
				Please refer to the field 'effective_labels' for all of the labels present on the resource.`,
			},

			"terraform_labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: `(ReadOnly) The combination of labels configured directly on the resource and default labels configured on the provider.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"effective_labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: `All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `A map of resource manager tags. Resource manager tag keys and values have the same definition as resource manager tags. Keys must be in the format tagKeys/{tag_key_id}, and values are in the format tagValues/456. The field is ignored when empty. This field is only set at create time and modifying this field after creation will trigger recreation. To apply tags to an existing resource, see the google_tags_tag_value resource.`,
			},
		},
		UseJSONNumber: true,
	}
}
