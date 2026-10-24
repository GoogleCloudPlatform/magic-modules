package resourcemanager

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-provider-google/google/registry"
	tpgcloudbilling "github.com/hashicorp/terraform-provider-google/google/services/cloudbilling"
	tpgcompute "github.com/hashicorp/terraform-provider-google/google/services/compute"
	rmClient "github.com/hashicorp/terraform-provider-google/google/services/resourcemanager/client"
	tpgserviceusage "github.com/hashicorp/terraform-provider-google/google/services/serviceusage"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"
	"google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/serviceusage/v1"
)

type ServicesCall interface {
	Header() http.Header
	Do(opts ...googleapi.CallOption) (*serviceusage.Operation, error)
}

// ResourceGoogleProject returns a *schema.Resource that allows a customer
// to declare a Google Cloud Project resource.
func ResourceGoogleProject() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Create: resourceGoogleProjectCreate,
		Read:   resourceGoogleProjectRead,
		Update: resourceGoogleProjectUpdate,
		Delete: resourceGoogleProjectDelete,

		CustomizeDiff: customdiff.All(
			tpgresource.DefaultProviderDeletionPolicy("PREVENT"),
			tpgresource.SetLabelsDiff,
			resourceGoogleProjectSoftDeleteCustomDiff,
		),

		Importer: &schema.ResourceImporter{
			State: resourceProjectImportState,
		},

		Identity: &schema.ResourceIdentity{
			Version: 1,
			SchemaFunc: func() map[string]*schema.Schema {
				return map[string]*schema.Schema{
					"project_id": {
						Type:              schema.TypeString,
						RequiredForImport: true,
					},
				}
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		MigrateState: resourceGoogleProjectMigrateState,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidateProjectID(),
				Description:  `The project ID. Changing this forces a new project to be created.`,
			},
			"recover_soft_deleted": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: `If true, Terraform will retain projects in DELETE_REQUESTED state during refresh and undelete and reconcile soft-deleted projects on create or update instead of removing them from state or failing with HTTP 409 alreadyExists.`,
			},
			"lifecycle_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The Project lifecycle state.`,
			},
			"deletion_policy": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
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

func remainingTimeout(deadline time.Time) (time.Duration, error) {
	rem := time.Until(deadline)
	if rem <= 0 {
		return 0, fmt.Errorf("operation timed out")
	}
	return rem, nil
}

func sleepWithDeadline(deadline time.Time, d time.Duration) error {
	rem, err := remainingTimeout(deadline)
	if err != nil {
		return err
	}
	if d > rem {
		time.Sleep(rem)
		return fmt.Errorf("operation timed out while waiting %v for propagation", d)
	}
	time.Sleep(d)
	return nil
}

func validateSoftDeletedProjectPreflight(pid string, recoverSoftDeleted bool, configuredOrgId string, existing *cloudresourcemanager.Project) error {
	if existing == nil {
		return nil
	}
	if existing.LifecycleState == "DELETE_IN_PROGRESS" {
		return fmt.Errorf("project %q is in %s state and cannot be created or recovered", pid, existing.LifecycleState)
	}
	if existing.LifecycleState != "DELETE_REQUESTED" {
		return nil
	}
	if !recoverSoftDeleted {
		return fmt.Errorf("project %q already exists in DELETE_REQUESTED (soft-deleted) state; set recover_soft_deleted = true on google_project to undelete and reconcile it via Terraform, or run 'gcloud projects undelete %s' and import it into Terraform state", pid, pid)
	}
	if existing.Parent != nil && existing.Parent.Type == "organization" && configuredOrgId != "" && existing.Parent.Id != "" && existing.Parent.Id != configuredOrgId {
		return fmt.Errorf("project %q exists in DELETE_REQUESTED state under organization %q, which conflicts with configured org_id %q", pid, existing.Parent.Id, configuredOrgId)
	}
	return nil
}

func resourceGoogleProjectSoftDeleteCustomDiff(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
	recoverSoftDeleted, _ := d.Get("recover_soft_deleted").(bool)
	if d.Id() != "" && !d.HasChange("project_id") {
		if state, _ := d.Get("lifecycle_state").(string); state == "DELETE_REQUESTED" {
			// If user toggles recover_soft_deleted from true to false on an in-state project,
			// allow the update so resourceGoogleProjectUpdate can remove it from state.
			if !recoverSoftDeleted && d.HasChange("recover_soft_deleted") {
				return nil
			}
			pid, _ := d.Get("project_id").(string)
			oldOrg, newOrg := d.GetChange("org_id")
			existing := &cloudresourcemanager.Project{
				ProjectId:      pid,
				LifecycleState: state,
			}
			if oldOrgStr, _ := oldOrg.(string); oldOrgStr != "" {
				existing.Parent = &cloudresourcemanager.ResourceId{
					Type: "organization",
					Id:   oldOrgStr,
				}
			}
			configuredOrgId, _ := newOrg.(string)
			if err := validateSoftDeletedProjectPreflight(pid, recoverSoftDeleted, configuredOrgId, existing); err != nil {
				return err
			}
			if err := d.SetNew("lifecycle_state", "ACTIVE"); err != nil {
				return err
			}
		}
		return nil
	}

	// For new resources (d.Id() == "") or project_id changes:
	// Short-circuit when recover_soft_deleted is false so default configurations
	// and existing VCR replay cassettes never issue pre-create Projects.Get RPCs.
	if !recoverSoftDeleted {
		return nil
	}
	if !d.NewValueKnown("project_id") {
		return nil
	}
	pid, ok := d.Get("project_id").(string)
	if !ok || pid == "" {
		return nil
	}
	if meta == nil {
		return nil
	}
	config, ok := meta.(*transport_tpg.Config)
	if !ok || config == nil || config.Client == nil {
		return nil
	}
	client := rmClient.NewClient(config, config.UserAgent)
	if client == nil {
		return nil
	}
	p, err := client.Projects.Get(pid).Do()
	if err != nil {
		log.Printf("[DEBUG] Preflight Projects.Get(%q) during CustomizeDiff returned error (proceeding with create plan): %v", pid, err)
		return nil
	}
	orgId, _ := d.Get("org_id").(string)
	if err := validateSoftDeletedProjectPreflight(pid, true, orgId, p); err != nil {
		return err
	}
	if p != nil && p.LifecycleState == "DELETE_REQUESTED" {
		log.Printf("[INFO] Project %q detected in DELETE_REQUESTED state during plan; will undelete and reconcile on apply.", pid)
		if err := d.SetNew("lifecycle_state", "ACTIVE"); err != nil {
			return err
		}
		if p.ProjectNumber > 0 {
			if err := d.SetNew("number", strconv.FormatInt(p.ProjectNumber, 10)); err != nil {
				return err
			}
		}
	}
	return nil
}

func resourceGoogleProjectCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	if err = resourceGoogleProjectCheckPreRequisites(config, d, userAgent); err != nil {
		return fmt.Errorf("failed pre-requisites: %v", err)
	}

	var pid string
	pid = d.Get("project_id").(string)
	deadline := time.Now().Add(d.Timeout(schema.TimeoutCreate))

	if d.Get("recover_soft_deleted").(bool) {
		var p *cloudresourcemanager.Project
		rem, remErr := remainingTimeout(deadline)
		if remErr != nil {
			return remErr
		}
		getErr := transport_tpg.Retry(transport_tpg.RetryOptions{
			RetryFunc: func() (reqErr error) {
				p, reqErr = rmClient.NewClient(config, userAgent).Projects.Get(pid).Do()
				return reqErr
			},
			Timeout: rem,
		})
		if getErr == nil && p != nil {
			if err := validateSoftDeletedProjectPreflight(pid, true, d.Get("org_id").(string), p); err != nil {
				return err
			}
			if p.LifecycleState == "DELETE_REQUESTED" {
				return recoverSoftDeletedProject(d, meta, config, userAgent, pid, p, deadline)
			}
		}
	}

	log.Printf("[DEBUG]: Creating new project %q", pid)
	project := &cloudresourcemanager.Project{
		ProjectId: pid,
		Name:      d.Get("name").(string),
	}

	if err = getParentResourceId(d, project); err != nil {
		return err
	}

	if _, ok := d.GetOk("effective_labels"); ok {
		project.Labels = tpgresource.ExpandEffectiveLabels(d)
	}

	if _, ok := d.GetOk("tags"); ok {
		project.Tags = tpgresource.ExpandStringMap(d, "tags")
	}

	var op *cloudresourcemanager.Operation
	err = transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() (reqErr error) {
			op, reqErr = rmClient.NewClient(config, userAgent).Projects.Create(project).Do()
			return reqErr
		},
		Timeout: d.Timeout(schema.TimeoutCreate),
	})
	if err != nil {
		if transport_tpg.IsGoogleApiErrorWithCode(err, 409) {
			var p *cloudresourcemanager.Project
			rem, remErr := remainingTimeout(deadline)
			if remErr == nil {
				getErr := transport_tpg.Retry(transport_tpg.RetryOptions{
					RetryFunc: func() (reqErr error) {
						p, reqErr = rmClient.NewClient(config, userAgent).Projects.Get(pid).Do()
						return reqErr
					},
					Timeout: rem,
				})
				if getErr == nil && p != nil && (p.LifecycleState == "DELETE_REQUESTED" || p.LifecycleState == "DELETE_IN_PROGRESS") {
					if p.LifecycleState == "DELETE_REQUESTED" && d.Get("recover_soft_deleted").(bool) {
						return recoverSoftDeletedProject(d, meta, config, userAgent, pid, p, deadline)
					}
					return fmt.Errorf("error creating project %s (%s): %w", project.ProjectId, project.Name, validateSoftDeletedProjectPreflight(pid, d.Get("recover_soft_deleted").(bool), d.Get("org_id").(string), p))
				}
			}
		}
		return fmt.Errorf("error creating project %s (%s): %s. "+
			"If you received a 403 error, make sure you have the"+
			" `roles/resourcemanager.projectCreator` permission",
			project.ProjectId, project.Name, err)
	}

	d.SetId(fmt.Sprintf("projects/%s", pid))

	if err := tpgresource.SetResourceIdentityAttributes(d, map[string]interface{}{
		"project_id": pid,
	}); err != nil {
		return err
	}

	// Wait for the operation to complete
	opAsMap, err := tpgresource.ConvertToMap(op)
	if err != nil {
		return err
	}

	waitErr := ResourceManagerOperationWaitTime(config, opAsMap, "creating folder", userAgent, d.Timeout(schema.TimeoutCreate))
	if waitErr != nil {
		// The resource wasn't actually created
		d.SetId("")
		return waitErr
	}

	// Set the billing account
	if _, ok := d.GetOk("billing_account"); ok {
		err = updateProjectBillingAccount(d, config, userAgent)
		if err != nil {
			return err
		}
	}

	// Sleep to let the billing account settle before other resources
	// try to use this project.
	time.Sleep(15 * time.Second)

	err = resourceGoogleProjectRead(d, meta)
	if err != nil {
		return err
	}

	// There's no such thing as "don't auto-create network", only "delete the network
	// post-creation" - but that's what it's called in the UI and let's not confuse
	// people if we don't have to.  The GCP Console is doing the same thing - creating
	// a network and deleting it in the background.
	if !d.Get("auto_create_network").(bool) {
		// The compute API has to be enabled before we can delete a network.

		billingProject := project.ProjectId
		// err == nil indicates that the billing_project value was found
		if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
			billingProject = bp
		}

		if err = EnableServiceUsageProjectServices([]string{"compute.googleapis.com"}, project.ProjectId, billingProject, userAgent, config, d.Timeout(schema.TimeoutCreate)); err != nil {
			return errwrap.Wrapf("Error enabling the Compute Engine API required to delete the default network: {{err}} ", err)
		}

		err = forceDeleteComputeNetwork(d, config, project.ProjectId, "default")
		// Retry if API is not yet enabled.
		if err != nil && transport_tpg.IsGoogleApiErrorWithCode(err, 403) {
			time.Sleep(15 * time.Second)
			err = forceDeleteComputeNetwork(d, config, project.ProjectId, "default")
		}
		if err != nil {
			if transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
				log.Printf("[DEBUG] Default network not found for project %q, no need to delete it", project.ProjectId)
			} else {
				return errwrap.Wrapf(fmt.Sprintf("Error deleting default network in project %s: {{err}}", project.ProjectId), err)
			}
		}
	}
	return nil
}

func undeleteGoogleProject(config *transport_tpg.Config, userAgent, pid string, deadline time.Time) (*cloudresourcemanager.Project, error) {
	rem, err := remainingTimeout(deadline)
	if err != nil {
		return nil, fmt.Errorf("error undeleting project %s: %w", pid, err)
	}

	err = transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() error {
			_, reqErr := rmClient.NewClient(config, userAgent).Projects.Undelete(pid, &cloudresourcemanager.UndeleteProjectRequest{}).Do()
			return reqErr
		},
		Timeout:      rem,
		PollInterval: config.PollInterval,
	})
	if err != nil {
		// If Undelete returned 400 (FAILED_PRECONDITION) or 409, check if another actor already undeleted it
		// or if a previous retry attempt succeeded on the server side.
		if transport_tpg.IsGoogleApiErrorWithCode(err, 400) || transport_tpg.IsGoogleApiErrorWithCode(err, 409) {
			p, getErr := rmClient.NewClient(config, userAgent).Projects.Get(pid).Do()
			if getErr == nil && p != nil {
				if p.LifecycleState == "ACTIVE" {
					log.Printf("[DEBUG] Project %q is already ACTIVE after Undelete returned %v; proceeding", pid, err)
					return p, nil
				}
				if p.LifecycleState == "DELETE_IN_PROGRESS" {
					return nil, fmt.Errorf("error undeleting project %s: project is in DELETE_IN_PROGRESS state and cannot be recovered: %w", pid, err)
				}
			}
		}
		return nil, fmt.Errorf("error undeleting project %s: %w", pid, err)
	}

	pollTimeout, err := remainingTimeout(deadline)
	if err != nil {
		return nil, fmt.Errorf("error waiting for undeleted project %s to become ACTIVE: %w", pid, err)
	}
	// Bound the post-undelete ACTIVE state check to at most 1 minute to avoid spinning for 10 minutes
	// if another actor immediately re-deleted the project.
	if pollTimeout > time.Minute {
		pollTimeout = time.Minute
	}

	var recovered *cloudresourcemanager.Project
	err = transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() error {
			var reqErr error
			recovered, reqErr = rmClient.NewClient(config, userAgent).Projects.Get(pid).Do()
			if reqErr != nil {
				return reqErr
			}
			if recovered == nil {
				return fmt.Errorf("project %q returned nil response after undelete", pid)
			}
			if recovered.LifecycleState == "ACTIVE" {
				return nil
			}
			if recovered.LifecycleState == "DELETE_REQUESTED" {
				return &googleapi.Error{
					Code:    503,
					Message: fmt.Sprintf("project %q lifecycleState is %q, waiting for ACTIVE", pid, recovered.LifecycleState),
				}
			}
			return fmt.Errorf("project %q is in unexpected lifecycleState %q while waiting for ACTIVE", pid, recovered.LifecycleState)
		},
		Timeout:      pollTimeout,
		PollInterval: config.PollInterval,
	})
	if err != nil {
		return nil, fmt.Errorf("error waiting for undeleted project %s to become ACTIVE: %w", pid, err)
	}
	return recovered, nil
}

func recoverSoftDeletedProject(d *schema.ResourceData, meta interface{}, config *transport_tpg.Config, userAgent, pid string, existing *cloudresourcemanager.Project, deadline time.Time) error {
	log.Printf("[DEBUG] Recovering soft-deleted project %q", pid)
	if err := validateSoftDeletedProjectPreflight(pid, true, d.Get("org_id").(string), existing); err != nil {
		return err
	}

	if _, ok := d.GetOk("tags"); ok {
		log.Printf("[WARN] Project %q is being recovered from DELETE_REQUESTED state; create-only 'tags' cannot be reconciled via Projects.Undelete or Projects.Update. Use google_tags_tag_binding to manage tags on existing projects.", pid)
	}

	recovered, err := undeleteGoogleProject(config, userAgent, pid, deadline)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("projects/%s", pid))
	if err := tpgresource.SetResourceIdentityAttributes(d, map[string]interface{}{
		"project_id": pid,
	}); err != nil {
		return err
	}

	desiredName := d.Get("name").(string)
	needsUpdate := false
	if recovered.Name != desiredName {
		recovered.Name = desiredName
		needsUpdate = true
	}

	desiredParent := &cloudresourcemanager.Project{}
	if err := getParentResourceId(d, desiredParent); err != nil {
		return err
	}
	if desiredParent.Parent != nil {
		if recovered.Parent == nil || recovered.Parent.Type != desiredParent.Parent.Type || recovered.Parent.Id != desiredParent.Parent.Id {
			recovered.Parent = desiredParent.Parent
			needsUpdate = true
		}
	}

	desiredLabels := tpgresource.ExpandEffectiveLabels(d)
	if !reflect.DeepEqual(recovered.Labels, desiredLabels) && (len(recovered.Labels) > 0 || len(desiredLabels) > 0) {
		recovered.Labels = desiredLabels
		needsUpdate = true
	}

	if needsUpdate {
		if _, err := updateProject(config, d, desiredName, userAgent, recovered); err != nil {
			return err
		}
	}

	_, hasConfiguredBilling := d.GetOk("billing_account")
	if hasConfiguredBilling {
		if err := updateProjectBillingAccount(d, config, userAgent); err != nil {
			return err
		}
		// Sleep 15s to let billing account propagation settle across downstream services (#3388, #12779).
		if err := sleepWithDeadline(deadline, 15*time.Second); err != nil {
			return err
		}
	}

	if err := resourceGoogleProjectRead(d, meta); err != nil {
		return err
	}

	// Invariant: project must remain ACTIVE after recovery; if another actor concurrently re-deleted it, abort.
	if d.Id() == "" || d.Get("lifecycle_state").(string) != "ACTIVE" {
		return fmt.Errorf("project %q was undeleted but is no longer ACTIVE (lifecycle_state=%q); it may have been concurrently deleted by another actor", pid, d.Get("lifecycle_state").(string))
	}

	if !hasConfiguredBilling && d.Get("billing_account").(string) != "" {
		if err := d.Set("billing_account", ""); err != nil {
			return fmt.Errorf("Error setting billing_account: %w", err)
		}
		if err := updateProjectBillingAccount(d, config, userAgent); err != nil {
			return err
		}
	}

	if !d.Get("auto_create_network").(bool) {
		billingProject := pid
		if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
			billingProject = bp
		}
		rem, err := remainingTimeout(deadline)
		if err != nil {
			return err
		}
		if err := EnableServiceUsageProjectServices([]string{"compute.googleapis.com"}, pid, billingProject, userAgent, config, rem); err != nil {
			return fmt.Errorf("Error enabling the Compute Engine API required to delete the default network: %w", err)
		}
		err = forceDeleteComputeNetwork(d, config, pid, "default")
		// Retry if API is not yet enabled (#10046, #12779).
		if err != nil && transport_tpg.IsGoogleApiErrorWithCode(err, 403) {
			if sleepErr := sleepWithDeadline(deadline, 15*time.Second); sleepErr != nil {
				return sleepErr
			}
			err = forceDeleteComputeNetwork(d, config, pid, "default")
		}
		if err != nil {
			if transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
				log.Printf("[DEBUG] Default network not found for project %q, no need to delete it", pid)
			} else {
				return fmt.Errorf("Error deleting default network in project %s: %w", pid, err)
			}
		}
	}

	return nil
}

func resourceGoogleProjectCheckPreRequisites(config *transport_tpg.Config, d *schema.ResourceData, userAgent string) error {
	ib, ok := d.GetOk("billing_account")
	if !ok {
		return nil
	}
	ba := "billingAccounts/" + ib.(string)
	const perm = "billing.resourceAssociations.create"
	req := &cloudbilling.TestIamPermissionsRequest{
		Permissions: []string{perm},
	}
	resp, err := tpgcloudbilling.NewClient(config, userAgent).BillingAccounts.TestIamPermissions(ba, req).Do()
	if err != nil {
		return fmt.Errorf("failed to check permissions on billing account %q: %v", ba, err)
	}
	if !tpgresource.StringInSlice(resp.Permissions, perm) {
		return fmt.Errorf("missing permission on %q: %v", ba, perm)
	}
	if !d.Get("auto_create_network").(bool) {
		call := tpgserviceusage.NewClient(config, userAgent).Services.Get("projects/00000000000/services/serviceusage.googleapis.com")
		if config.UserProjectOverride {
			if billingProject, err := tpgresource.GetBillingProject(d, config); err == nil {
				call.Header().Add("X-Goog-User-Project", billingProject)
			}
		}
		_, err := call.Do()
		switch {
		// We are querying a dummy project since the call is already coming from the quota project.
		// If the API is enabled we get a not found message or accessNotConfigured if API is not enabled.
		case err.Error() == "googleapi: Error 403: Project '00000000000' not found or permission denied., forbidden":
			return nil
		case strings.Contains(err.Error(), "accessNotConfigured"):
			return fmt.Errorf("API serviceusage not enabled.\nFound error: %v", err)
		}
	}
	return nil
}

func resourceGoogleProjectRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	parts := strings.Split(d.Id(), "/")
	pid := parts[len(parts)-1]

	p, err := readGoogleProject(d, config, userAgent)
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 403 && strings.Contains(gerr.Message, "caller does not have permission") {
			return fmt.Errorf("the user does not have permission to access Project %q or it may not exist", pid)
		}
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Project %q", pid))
	}

	// If the project has been deleted from outside Terraform, remove it from state file
	// unless it is in DELETE_REQUESTED state and recover_soft_deleted is enabled.
	if p.LifecycleState != "ACTIVE" {
		recoverSoftDeleted, _ := d.Get("recover_soft_deleted").(bool)
		if p.LifecycleState == "DELETE_REQUESTED" && recoverSoftDeleted {
			log.Printf("[INFO] Project '%s' is in 'DELETE_REQUESTED' state and recover_soft_deleted is true; retaining in state.", pid)
			if err := populateGoogleProjectResourceData(d, p, pid, config); err != nil {
				return err
			}
			if err := d.Set("billing_account", ""); err != nil {
				return fmt.Errorf("Error setting billing_account: %s", err)
			}
			return nil
		}
		log.Printf("[WARN] Removing project '%s' from state because its lifecycle_state is '%s' (requires 'ACTIVE'). To recover a DELETE_REQUESTED project via Terraform, set recover_soft_deleted = true on the google_project resource and run terraform apply.", pid, p.LifecycleState)
		d.SetId("")
		return nil
	}
	// Explicitly set client-side fields to default values if unset

	if err := populateGoogleProjectResourceData(d, p, pid, config); err != nil {
		return err
	}

	var ba *cloudbilling.ProjectBillingInfo
	err = transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() (reqErr error) {
			ba, reqErr = tpgcloudbilling.NewClient(config, userAgent).Projects.GetBillingInfo(PrefixedProject(pid)).Do()
			return reqErr
		},
		Timeout: d.Timeout(schema.TimeoutRead),
	})
	// Read the billing account
	if err != nil && !transport_tpg.IsApiNotEnabledError(err) {
		return fmt.Errorf("Error reading billing account for project %q: %v", PrefixedProject(pid), err)
	} else if transport_tpg.IsApiNotEnabledError(err) {
		log.Printf("[WARN] Billing info API not enabled, please enable it to read billing info about project %q: %s", pid, err.Error())
	} else if ba.BillingAccountName != "" {
		// BillingAccountName is contains the resource name of the billing account
		// associated with the project, if any. For example,
		// `billingAccounts/012345-567890-ABCDEF`. We care about the ID and not
		// the `billingAccounts/` prefix, so we need to remove that. If the
		// prefix ever changes, we'll validate to make sure it's something we
		// recognize.
		_ba := strings.TrimPrefix(ba.BillingAccountName, "billingAccounts/")
		if ba.BillingAccountName == _ba {
			return fmt.Errorf("Error parsing billing account for project %q. Expected value to begin with 'billingAccounts/' but got %s", PrefixedProject(pid), ba.BillingAccountName)
		}
		if err := d.Set("billing_account", _ba); err != nil {
			return fmt.Errorf("Error setting billing_account: %s", err)
		}
	}

	return nil
}

func populateGoogleProjectResourceData(d *schema.ResourceData, p *cloudresourcemanager.Project, pid string, config *transport_tpg.Config) error {
	if err := tpgresource.DeletionPolicyReadDefault(d, config, "PREVENT"); err != nil {
		return err
	}

	if _, ok := d.GetOkExists("recover_soft_deleted"); !ok {
		if err := d.Set("recover_soft_deleted", false); err != nil {
			return fmt.Errorf("Error setting recover_soft_deleted: %s", err)
		}
	}

	if err := d.Set("lifecycle_state", p.LifecycleState); err != nil {
		return fmt.Errorf("Error setting lifecycle_state: %s", err)
	}

	if err := d.Set("project_id", pid); err != nil {
		return fmt.Errorf("Error setting project_id: %s", err)
	}
	if err := d.Set("number", strconv.FormatInt(p.ProjectNumber, 10)); err != nil {
		return fmt.Errorf("Error setting number: %s", err)
	}
	if err := d.Set("name", p.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := tpgresource.SetLabels(p.Labels, d, "labels"); err != nil {
		return fmt.Errorf("Error setting labels: %s", err)
	}
	if err := tpgresource.SetLabels(p.Labels, d, "terraform_labels"); err != nil {
		return fmt.Errorf("Error setting terraform_labels: %s", err)
	}
	if err := d.Set("effective_labels", p.Labels); err != nil {
		return fmt.Errorf("Error setting effective_labels: %s", err)
	}

	if p.Parent != nil {
		switch p.Parent.Type {
		case "organization":
			if err := d.Set("org_id", p.Parent.Id); err != nil {
				return fmt.Errorf("Error setting org_id: %s", err)
			}
			if err := d.Set("folder_id", ""); err != nil {
				return fmt.Errorf("Error setting folder_id: %s", err)
			}
		case "folder":
			if err := d.Set("folder_id", p.Parent.Id); err != nil {
				return fmt.Errorf("Error setting folder_id: %s", err)
			}
			if err := d.Set("org_id", ""); err != nil {
				return fmt.Errorf("Error setting org_id: %s", err)
			}
		}
	}
	return tpgresource.SetResourceIdentityAttributes(d, map[string]interface{}{
		"project_id": pid,
	})
}

func PrefixedProject(pid string) string {
	return "projects/" + pid
}

func getParentResourceId(d *schema.ResourceData, p *cloudresourcemanager.Project) error {
	orgId := d.Get("org_id").(string)
	folderId := d.Get("folder_id").(string)

	if orgId != "" && folderId != "" {
		return fmt.Errorf("'org_id' and 'folder_id' cannot be both set.")
	}

	if orgId != "" {
		p.Parent = &cloudresourcemanager.ResourceId{
			Id:   orgId,
			Type: "organization",
		}
	}

	if folderId != "" {
		p.Parent = &cloudresourcemanager.ResourceId{
			Id:   ParseFolderId(folderId),
			Type: "folder",
		}
	}

	return nil
}

func ParseFolderId(v interface{}) string {
	folderId, ok := v.(string)
	if !ok {
		return ""
	}
	if strings.HasPrefix(folderId, "folders/") {
		return folderId[8:]
	}
	return folderId
}

func resourceGoogleProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	parts := strings.Split(d.Id(), "/")
	pid := parts[len(parts)-1]
	project_name := d.Get("name").(string)

	// Read the project
	// we need the project even though refresh has already been called
	// because the API doesn't support patch, so we need the actual object
	p, err := readGoogleProject(d, config, userAgent)
	if err != nil {
		if transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
			return fmt.Errorf("Project %q does not exist.", pid)
		}
		return fmt.Errorf("Error checking project %q: %s", pid, err)
	}

	d.Partial(true)

	undeleted := false
	if p.LifecycleState == "DELETE_REQUESTED" {
		if !d.Get("recover_soft_deleted").(bool) {
			d.Partial(false)
			return resourceGoogleProjectRead(d, meta)
		}
		if err := validateSoftDeletedProjectPreflight(pid, true, d.Get("org_id").(string), p); err != nil {
			return err
		}
		deadline := time.Now().Add(d.Timeout(schema.TimeoutUpdate))
		if p, err = undeleteGoogleProject(config, userAgent, pid, deadline); err != nil {
			return err
		}
		undeleted = true
	}

	// Project display name has changed
	if ok := d.HasChange("name"); ok {
		p.Name = project_name
		// Do update on project
		if p, err = updateProject(config, d, project_name, userAgent, p); err != nil {
			return err
		}
	}

	// Project parent has changed
	if d.HasChange("org_id") || d.HasChange("folder_id") {
		if err := getParentResourceId(d, p); err != nil {
			return err
		}

		// Do update on project
		if p, err = updateProject(config, d, project_name, userAgent, p); err != nil {
			return err
		}
	}

	// Billing account has changed (or needs re-attachment after undelete)
	if d.HasChange("billing_account") || (undeleted && d.Get("billing_account").(string) != "") {
		err = updateProjectBillingAccount(d, config, userAgent)
		if err != nil {
			return err
		}
	}

	// Project Labels have changed
	if ok := d.HasChange("effective_labels"); ok {
		p.Labels = tpgresource.ExpandEffectiveLabels(d)

		// Do Update on project
		if p, err = updateProject(config, d, project_name, userAgent, p); err != nil {
			return err
		}
	}

	d.Partial(false)
	if err := tpgresource.SetResourceIdentityAttributes(d, map[string]interface{}{
		"project_id": pid,
	}); err != nil {
		return err
	}
	if err := resourceGoogleProjectRead(d, meta); err != nil {
		return err
	}
	if undeleted {
		// Invariant: project must remain ACTIVE after recovery; if another actor concurrently re-deleted it, abort.
		if d.Id() == "" || d.Get("lifecycle_state").(string) != "ACTIVE" {
			return fmt.Errorf("project %q was undeleted but is no longer ACTIVE (lifecycle_state=%q); it may have been concurrently deleted by another actor", pid, d.Get("lifecycle_state").(string))
		}
		_, hasConfiguredBilling := d.GetOk("billing_account")
		if !hasConfiguredBilling && d.Get("billing_account").(string) != "" {
			if err := d.Set("billing_account", ""); err != nil {
				return fmt.Errorf("Error setting billing_account: %w", err)
			}
			if err := updateProjectBillingAccount(d, config, userAgent); err != nil {
				return err
			}
		}
	}
	return nil
}

func updateProject(config *transport_tpg.Config, d *schema.ResourceData, projectName, userAgent string, desiredProject *cloudresourcemanager.Project) (*cloudresourcemanager.Project, error) {
	var newProj *cloudresourcemanager.Project
	if err := transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() (updateErr error) {
			newProj, updateErr = rmClient.NewClient(config, userAgent).Projects.Update(desiredProject.ProjectId, desiredProject).Do()
			return updateErr
		},
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}); err != nil {
		return nil, fmt.Errorf("Error updating project %q: %s", projectName, err)
	}
	return newProj, nil
}

func resourceGoogleProjectDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	deletionPolicy := d.Get("deletion_policy").(string)

	if deletionPolicy == "PREVENT" {
		return fmt.Errorf("Cannot destroy project as deletion_policy is set to PREVENT.")
	} else if deletionPolicy == "ABANDON" {
		log.Printf("[WARN] The project has been abandoned as deletion_policy set to ABANDON.")
		d.SetId("")
		return nil
	} else {
		// Only delete projects if deletion_policy isn't PREVENT or ABANDON
		parts := strings.Split(d.Id(), "/")
		pid := parts[len(parts)-1]
		if err := transport_tpg.Retry(transport_tpg.RetryOptions{
			RetryFunc: func() error {
				_, delErr := rmClient.NewClient(config, userAgent).Projects.Delete(pid).Do()
				return delErr
			},
			Timeout: d.Timeout(schema.TimeoutDelete),
		}); err != nil {
			return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Project %s", pid))
		}
	}
	d.SetId("")
	return nil
}

func resourceProjectImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// Handle identity-based import (Terraform 1.12+ import blocks with resource identity).
	// In this case d.Id() is empty and project_id is available via the identity block.
	if d.Id() == "" {
		identity, err := d.Identity()
		if err != nil || identity == nil {
			return nil, fmt.Errorf("Error getting identity for import: %v", err)
		}
		pidVal, ok := identity.GetOk("project_id")
		if !ok {
			return nil, fmt.Errorf("project_id must be set in the identity block for import")
		}
		pid := pidVal.(string)
		d.SetId(fmt.Sprintf("projects/%s", pid))
		if err := d.Set("auto_create_network", true); err != nil {
			return nil, fmt.Errorf("Error setting auto_create_network: %s", err)
		}
		return []*schema.ResourceData{d}, nil
	}

	parts := strings.Split(d.Id(), "/")
	pid := parts[len(parts)-1]
	// Prevent importing via project number, this will cause issues later
	matched, err := regexp.MatchString("^\\d+$", pid)
	if err != nil {
		return nil, fmt.Errorf("Error matching project %q: %s", pid, err)
	}

	if matched {
		return nil, fmt.Errorf("Error importing project %q, please use project_id", pid)
	}

	// Ensure the id format includes projects/
	d.SetId(fmt.Sprintf("projects/%s", pid))

	// Explicitly set to default as a workaround for `ImportStateVerify` tests, and so that users
	// don't see a diff immediately after import.
	if err := d.Set("auto_create_network", true); err != nil {
		return nil, fmt.Errorf("Error setting auto_create_network: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}

// Delete a compute network along with the firewall rules inside it.
func forceDeleteComputeNetwork(d *schema.ResourceData, config *transport_tpg.Config, projectId, networkName string) error {
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	// Read the network from the API so we can get the correct self link format. We can't construct it from the
	// base path because it might not line up exactly (compute.googleapis.com vs www.googleapis.com)
	networkUrl := fmt.Sprintf("%sprojects/%s/global/networks/%s", transport_tpg.BaseUrl(tpgcompute.Product, config), projectId, networkName)
	net, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   projectId,
		RawURL:    networkUrl,
		UserAgent: userAgent,
	})
	if err != nil {
		return err
	}
	netSelfLink, _ := net["selfLink"].(string)

	token := ""
	for paginate := true; paginate; {
		params := map[string]string{
			"filter": fmt.Sprintf("network eq %s", netSelfLink),
		}
		if token != "" {
			params["pageToken"] = token
		}
		firewallsUrl, err := transport_tpg.AddQueryParams(
			fmt.Sprintf("%sprojects/%s/global/firewalls", transport_tpg.BaseUrl(tpgcompute.Product, config), projectId),
			params)
		if err != nil {
			return err
		}

		resp, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   projectId,
			RawURL:    firewallsUrl,
			UserAgent: userAgent,
		})
		if err != nil {
			return errwrap.Wrapf("Error listing firewall rules in proj: {{err}}", err)
		}

		items, _ := resp["items"].([]interface{})
		log.Printf("[DEBUG] Found %d firewall rules in %q network", len(items), networkName)

		for _, item := range items {
			firewall, _ := item.(map[string]interface{})
			firewallName, _ := firewall["name"].(string)

			deleteUrl := fmt.Sprintf("%sprojects/%s/global/firewalls/%s", transport_tpg.BaseUrl(tpgcompute.Product, config), projectId, firewallName)
			op, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "DELETE",
				Project:   projectId,
				RawURL:    deleteUrl,
				UserAgent: userAgent,
			})
			if err != nil {
				return errwrap.Wrapf("Error deleting firewall: {{err}}", err)
			}
			err = tpgcompute.ComputeOperationWaitTime(config, op, projectId, "Deleting Firewall", userAgent, d.Timeout(schema.TimeoutCreate))
			if err != nil {
				return err
			}
		}

		if pToken, ok := resp["nextPageToken"]; ok && pToken != nil && pToken.(string) != "" {
			token = pToken.(string)
		} else {
			token = ""
		}
		paginate = token != ""
	}

	return deleteComputeNetwork(projectId, networkName, userAgent, config)
}

func updateProjectBillingAccount(d *schema.ResourceData, config *transport_tpg.Config, userAgent string) error {
	parts := strings.Split(d.Id(), "/")
	pid := parts[len(parts)-1]
	name := d.Get("billing_account").(string)
	ba := &cloudbilling.ProjectBillingInfo{}
	// If we're unlinking an existing billing account, an empty request does that, not an empty-string billing account.
	if name != "" {
		ba.BillingAccountName = "billingAccounts/" + name
	}
	updateBillingInfoFunc := func() error {
		_, err := tpgcloudbilling.NewClient(config, userAgent).Projects.UpdateBillingInfo(PrefixedProject(pid), ba).Do()
		return err
	}
	err := transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: updateBillingInfoFunc,
		Timeout:   d.Timeout(schema.TimeoutUpdate),
	})
	if err != nil {
		if err := d.Set("billing_account", ""); err != nil {
			return fmt.Errorf("Error setting billing_account: %s", err)
		}
		if _err, ok := err.(*googleapi.Error); ok {
			return fmt.Errorf("Error setting billing account %q for project %q: %v", name, PrefixedProject(pid), _err)
		}
		return fmt.Errorf("Error setting billing account %q for project %q: %v", name, PrefixedProject(pid), err)
	}
	for retries := 0; retries < 3; retries++ {
		var ba *cloudbilling.ProjectBillingInfo
		err = transport_tpg.Retry(transport_tpg.RetryOptions{
			RetryFunc: func() (reqErr error) {
				ba, reqErr = tpgcloudbilling.NewClient(config, userAgent).Projects.GetBillingInfo(PrefixedProject(pid)).Do()
				return reqErr
			},
			Timeout: d.Timeout(schema.TimeoutRead),
		})
		if err != nil {
			return fmt.Errorf("Error getting billing info for project %q: %v", PrefixedProject(pid), err)
		}
		baName := strings.TrimPrefix(ba.BillingAccountName, "billingAccounts/")
		if baName == name {
			return nil
		}
		time.Sleep(3 * time.Second)
	}
	return fmt.Errorf("Timed out waiting for billing account to return correct value.  Waiting for %s, got %s.",
		name, strings.TrimPrefix(ba.BillingAccountName, "billingAccounts/"))
}

func deleteComputeNetwork(project, network, userAgent string, config *transport_tpg.Config) error {
	url := fmt.Sprintf("%sprojects/%s/global/networks/%s", transport_tpg.BaseUrl(tpgcompute.Product, config), project, network)
	op, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "DELETE",
		Project:   project,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return errwrap.Wrapf("Error deleting network: {{err}}", err)
	}

	err = tpgcompute.ComputeOperationWaitTime(config, op, project, "Deleting Network", userAgent, 10*time.Minute)
	if err != nil {
		return err
	}
	return nil
}

func readGoogleProject(d *schema.ResourceData, config *transport_tpg.Config, userAgent string) (*cloudresourcemanager.Project, error) {
	var p *cloudresourcemanager.Project
	// Read the project
	parts := strings.Split(d.Id(), "/")
	pid := parts[len(parts)-1]
	err := transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() (reqErr error) {
			p, reqErr = rmClient.NewClient(config, userAgent).Projects.Get(pid).Do()
			return reqErr
		},
		Timeout: d.Timeout(schema.TimeoutRead),
	})
	return p, err
}

// Enables services. WARNING: Use BatchRequestEnableServices for better batching if possible.
func EnableServiceUsageProjectServices(services []string, project, billingProject, userAgent string, config *transport_tpg.Config, timeout time.Duration) error {
	// ServiceUsage does not allow more than 20 services to be enabled per
	// batchEnable API call. See
	// https://cloud.google.com/service-usage/docs/reference/rest/v1/services/batchEnable
	for i := 0; i < len(services); i += maxServiceUsageBatchSize {
		j := i + maxServiceUsageBatchSize
		if j > len(services) {
			j = len(services)
		}
		nextBatch := services[i:j]
		if len(nextBatch) == 0 {
			// All batches finished, return.
			return nil
		}

		if err := doEnableServicesRequest(nextBatch, project, billingProject, userAgent, config, timeout); err != nil {
			return err
		}
		log.Printf("[DEBUG] Finished enabling next batch of %d project services: %+v", len(nextBatch), nextBatch)
	}

	log.Printf("[DEBUG] Verifying that all services are enabled")
	return waitForServiceUsageEnabledServices(services, project, billingProject, userAgent, config, timeout)
}

func doEnableServicesRequest(services []string, project, billingProject, userAgent string, config *transport_tpg.Config, timeout time.Duration) error {
	var op *serviceusage.Operation

	// errors can come up at multiple points, so there are a few levels of
	// retrying here.
	// logicalErr / waitErr: overall error on the logical operation (enabling services)
	// but possibly also errors when retrieving the LRO (these are rare)
	// err / reqErr: precondition errors when sending the request received instead of an LRO
	logicalErr := transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() error {
			err := transport_tpg.Retry(transport_tpg.RetryOptions{
				RetryFunc: func() error {
					var reqErr error
					var call ServicesCall
					if len(services) == 1 {
						// BatchEnable returns an error for a single item, so enable with single endpoint
						name := fmt.Sprintf("projects/%s/services/%s", project, services[0])
						req := &serviceusage.EnableServiceRequest{}
						call = tpgserviceusage.NewClient(config, userAgent).Services.Enable(name, req)
					} else {
						// Batch enable for multiple services.
						name := fmt.Sprintf("projects/%s", project)
						req := &serviceusage.BatchEnableServicesRequest{ServiceIds: services}
						call = tpgserviceusage.NewClient(config, userAgent).Services.BatchEnable(name, req)
					}

					if config.UserProjectOverride && billingProject != "" {
						call.Header().Add("X-Goog-User-Project", billingProject)
					}

					op, reqErr = call.Do()
					return handleServiceUsageRetryablePreconditionError(reqErr)
				},
				Timeout:              timeout,
				ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.ServiceUsageServiceBeingActivated},
			})
			if err != nil {
				return errwrap.Wrapf("failed on request preconditions: {{err}}", err)
			}

			waitErr := tpgserviceusage.ServiceUsageOperationWait(config, op, billingProject, fmt.Sprintf("Enable Project %q Services: %+v", project, services), userAgent, timeout)
			if waitErr != nil {
				return waitErr
			}

			return nil
		},
		Timeout:              timeout,
		ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.ServiceUsageInternalError160009},
	})

	if logicalErr != nil {
		return errwrap.Wrapf("failed to enable services: {{err}}", logicalErr)
	}

	return nil
}

// Handle errors that are retryable at call time for serviceusage
// Specifically, errors in https://cloud.google.com/service-usage/docs/reference/rest/v1/services/batchEnable#response-body
// Errors in operations are handled separately.
// TODO: This should probably be turned into a retry predicate
func handleServiceUsageRetryablePreconditionError(err error) error {
	if err == nil {
		return nil
	}
	if gerr, ok := err.(*googleapi.Error); ok {
		if (gerr.Code == 400 || gerr.Code == 412) && gerr.Message == "Precondition check failed." {
			return &googleapi.Error{
				Code:    503,
				Message: "api returned \"precondition failed\" while enabling service",
			}
		}
	}
	return err
}

// Retrieve a project's services from the API
// if a service has been renamed, this function will list both the old and new
// forms of the service. LIST responses are expected to return only the old or
// new form, but we'll always return both.
func ListCurrentlyEnabledServices(project, billingProject, userAgent string, config *transport_tpg.Config, timeout time.Duration) (map[string]struct{}, error) {
	log.Printf("[DEBUG] Listing enabled services for project %s", project)
	apiServices := make(map[string]struct{})
	err := transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() error {
			ctx := context.Background()
			call := tpgserviceusage.NewClient(config, userAgent).Services.List(fmt.Sprintf("projects/%s", project)).PageSize(200)
			if config.UserProjectOverride && billingProject != "" {
				call.Header().Add("X-Goog-User-Project", billingProject)
			}
			return call.Fields("services/name,nextPageToken").Filter("state:ENABLED").
				Pages(ctx, func(r *serviceusage.ListServicesResponse) error {
					for _, v := range r.Services {
						// services are returned as "projects/{{project}}/services/{{name}}"
						name := tpgresource.GetResourceNameFromSelfLink(v.Name)

						apiServices[name] = struct{}{}

						// if a service has been renamed, set both. We'll deal
						// with setting the right values later.
						if v, ok := renamedServicesByOldAndNewServiceNames[name]; ok {
							log.Printf("[DEBUG] Adding service alias for %s to enabled services: %s", name, v)
							apiServices[v] = struct{}{}
						}
					}
					return nil
				})
		},
		Timeout: timeout,
	})
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Failed to list enabled services for project %s: {{err}}", project), err)
	}
	return apiServices, nil
}

// waitForServiceUsageEnabledServices doesn't resend enable requests - it just
// waits for service enablement status to propagate. Essentially, it waits until
// all services show up as enabled when listing services on the project.
func waitForServiceUsageEnabledServices(services []string, project, billingProject, userAgent string, config *transport_tpg.Config, timeout time.Duration) error {
	missing := make([]string, 0, len(services))
	delay := time.Duration(0)
	interval := time.Second
	err := transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() error {
			// Get the list of services that are enabled on the project
			enabledServices, err := ListCurrentlyEnabledServices(project, billingProject, userAgent, config, timeout)
			if err != nil {
				return err
			}

			missing := make([]string, 0, len(services))
			for _, s := range services {
				if _, ok := enabledServices[s]; !ok {
					missing = append(missing, s)
				}
			}
			if len(missing) > 0 {
				log.Printf("[DEBUG] waiting %v before reading project %s services...", delay, project)
				time.Sleep(delay)
				delay += interval
				interval += delay

				// Spoof a googleapi Error so retryTime will try again
				return &googleapi.Error{
					Code:    503,
					Message: fmt.Sprintf("The service(s) %q are still being enabled for project %s. This isn't a real API error, this is just eventual consistency.", missing, project),
				}
			}
			return nil
		},
		Timeout: timeout,
	})
	if err != nil {
		return errwrap.Wrap(err, fmt.Errorf("failed to enable some service(s) %q for project %s", missing, project))
	}
	return nil
}

func init() {
	registry.Schema{
		Name:        "google_project",
		ProductName: "resourcemanager",
		Type:        registry.SchemaTypeResource,
		Schema:      ResourceGoogleProject(),
	}.Register()
}
