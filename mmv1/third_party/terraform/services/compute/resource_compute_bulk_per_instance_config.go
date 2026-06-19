package compute

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"

	"google.golang.org/api/googleapi"
)

var (
	_ = bytes.Clone
	_ = context.WithCancel
	_ = base64.NewDecoder
	_ = json.Marshal
	_ = fmt.Sprintf
	_ = log.Print
	_ = http.Get
	_ = reflect.ValueOf
	_ = regexp.Match
	_ = slices.Min([]int{1})
	_ = sort.IntSlice{}
	_ = strconv.Atoi
	_ = strings.Trim
	_ = time.Now
	_ = errwrap.Wrap
	_ = cty.BoolVal
	_ = diag.Diagnostic{}
	_ = customdiff.All
	_ = id.UniqueId
	_ = logging.LogLevel
	_ = retry.Retry
	_ = schema.Noop
	_ = validation.All
	_ = structure.ExpandJsonFromString
	_ = terraform.State{}
	_ = tpgresource.SetLabels
	_ = transport_tpg.Config{}
	_ = verify.ValidateEnum
	_ = googleapi.Error{}
)

func init() {
	registry.Schema{
		Name:        "google_compute_bulk_per_instance_config",
		ProductName: "compute",
		Type:        registry.SchemaTypeResource,
		Schema:      ResourceComputeBulkPerInstanceConfig(),
	}.Register()
}

func ResourceComputeBulkPerInstanceConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeBulkPerInstanceConfigCreate,
		Read:   resourceComputeBulkPerInstanceConfigRead,
		Delete: resourceComputeBulkPerInstanceConfigDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			tpgresource.DefaultProviderProject,
			tpgresource.DefaultProviderZone,
			tpgresource.DefaultProviderDeletionPolicy("DELETE"),
		),

		Identity: &schema.ResourceIdentity{
			Version: 1,
			SchemaFunc: func() map[string]*schema.Schema {
				return map[string]*schema.Schema{
					"name": {
						Type:              schema.TypeString,
						OptionalForImport: true,
					},
					"zone": {
						Type:              schema.TypeString,
						OptionalForImport: true,
					},
					"instance_group_manager": {
						Type:              schema.TypeString,
						RequiredForImport: true,
					},
					"project": {
						Type:              schema.TypeString,
						OptionalForImport: true,
					},
				}
			},
		},
		ResourceBehavior: schema.ResourceBehavior{
			MutableIdentity: false,
		},

		Schema: map[string]*schema.Schema{
			"instance_group_manager": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      `The instance group manager this bulk per instance config is part of.`,
			},
			"per_instance_configs": {
				Type:        schema.TypeSet,
				Required:    true,
				ForceNew:    true,
				Description: `The list of per-instance configs.`,
				Elem:        computeBulkPerInstanceConfigPerInstanceConfigsSchema(),
				// Default schema.HashSchema is used.
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name for this bulk per-instance config.`,
			},
			"zone": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      `Zone where the containing instance group manager is located`,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"deletion_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: `Whether Terraform will be prevented from destroying the bulk per instance config. Defaults to "DELETE".
When a 'terraform destroy' or 'terraform apply' would delete the bulk per instance config,
the command will fail if this field is set to "PREVENT" in Terraform state.
When set to "ABANDON", the command will remove the resource from Terraform
management without updating or deleting the resource in the API.
When set to "DELETE", deleting the resource is allowed.
`,
			},
		},
		UseJSONNumber: true,
	}
}

func computeBulkPerInstanceConfigPerInstanceConfigsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name for this per-instance config and its corresponding instance.`,
			},
		},
	}
}

func resourceComputeBulkPerInstanceConfigCreate(d *schema.ResourceData, meta any) error {
	config := meta.(*transport_tpg.Config)

	lockName, err := tpgresource.ReplaceVars(d, config, "instanceGroupManager/{{project}}/{{zone}}/{{instance_group_manager}}/{{name}}")
	if err != nil {
		return err
	}
	transport_tpg.MutexStore.Lock(lockName)
	defer transport_tpg.MutexStore.Unlock(lockName)

	log.Printf("[DEBUG] Creating new BulkPerInstanceConfig: %s", d.Get("name"))

	perInstanceConfigs, err := expandRawPerInstanceConfigs(d.Get("per_instance_configs"))
	if err != nil {
		return err
	}
	err = callCreateInstances(d, meta, perInstanceConfigs)

	if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating BulkPerInstanceConfig: %s", err)
	}

	id, err := tpgresource.ReplaceVars(d, config, "{{project}}/{{zone}}/{{instance_group_manager}}/{{name}}")
	if err != nil {

		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating BulkPerInstanceConfig %q", d.Id())

	identity, err := d.Identity()
	if err == nil && identity != nil {
		if nameValue, ok := d.GetOk("name"); ok && nameValue.(string) != "" {
			if err = identity.Set("name", nameValue.(string)); err != nil {
				return fmt.Errorf("Error setting name: %s", err)
			}
		}
		if zoneValue, ok := d.GetOk("zone"); ok && zoneValue.(string) != "" {
			if err = identity.Set("zone", zoneValue.(string)); err != nil {
				return fmt.Errorf("Error setting zone: %s", err)
			}
		}
		if instanceGroupManagerValue, ok := d.GetOk("instance_group_manager"); ok && instanceGroupManagerValue.(string) != "" {
			if err = identity.Set("instance_group_manager", instanceGroupManagerValue.(string)); err != nil {
				return fmt.Errorf("Error setting instance_group_manager: %s", err)
			}
		}
		if projectValue, ok := d.GetOk("project"); ok && projectValue.(string) != "" {
			if err = identity.Set("project", projectValue.(string)); err != nil {
				return fmt.Errorf("Error setting project: %s", err)
			}
		}
	} else {
		log.Printf("[DEBUG] (Create) identity not set: %s", err)
	}

	return resourceComputeBulkPerInstanceConfigRead(d, meta)
}

func resourceComputeBulkPerInstanceConfigRead(d *schema.ResourceData, meta any) error {
	res, err := getListManagedInstancesResponse(d, meta)
	if err != nil {
		return err
	}

	if res == nil {
		// Object isn't there any more - remove it from the state.
		log.Printf("[DEBUG] Removing ComputeBulkPerInstanceConfig because it couldn't be matched.")
		d.SetId("")
		return nil
	}

	// ListManagedInstances returns all managed instances for a given instance group manager.
	// BulkPerInstanceConfig manages only a selection of those instances.
	perInstanceConfigs, err := expandRawPerInstanceConfigs(d.Get("per_instance_configs"))
	if err != nil {
		return nil
	}
	thisPerInstanceConfigsInstanceNames := getBulkPerInstanceConfigInstanceNames(perInstanceConfigs)
	var filteredPerInstanceConfigs []any
	for _, perInstanceConfig := range res["perInstanceConfigs"].([]any) {
		if slices.Contains(thisPerInstanceConfigsInstanceNames, perInstanceConfig.(map[string]any)["name"].(string)) {
			filteredPerInstanceConfigs = append(filteredPerInstanceConfigs, perInstanceConfig.(map[string]any))
		}
	}
	res["perInstanceConfigs"] = filteredPerInstanceConfigs

	config := meta.(*transport_tpg.Config)
	// Explicitly set virtual fields to default values if unset
	if _, ok := d.GetOkExists("deletion_policy"); !ok {
		//prioritize config's value if present
		if config.DeletionPolicy != "" {
			if err := d.Set("deletion_policy", config.DeletionPolicy); err != nil {
				return fmt.Errorf("Error setting deletion_policy: %s", err)
			}
		} else {
			if err := d.Set("deletion_policy", "DELETE"); err != nil {
				return fmt.Errorf("Error setting deletion_policy: %s", err)
			}
		}
	}
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for BulkPerInstanceConfig: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading BulkPerInstanceConfig: %s", err)
	}

	zone, err := tpgresource.GetZone(d, config)
	if err != nil {
		return err
	}
	if err := d.Set("zone", zone); err != nil {
		return fmt.Errorf("Error reading BulkPerInstanceConfig: %s", err)
	}
	err = ResourceComputeBulkPerInstanceConfigFlatten(d, res)
	if err != nil {
		return err
	}

	identity, err := d.Identity()
	if err == nil && identity != nil {
		if v, ok := identity.GetOk("name"); !ok && v == "" {
			err = identity.Set("name", d.Get("name").(string))
			if err != nil {
				return fmt.Errorf("Error setting name: %s", err)
			}
		}
		if v, ok := identity.GetOk("zone"); !ok && v == "" {
			err = identity.Set("zone", d.Get("zone").(string))
			if err != nil {
				return fmt.Errorf("Error setting zone: %s", err)
			}
		}
		if v, ok := identity.GetOk("instance_group_manager"); !ok && v == "" {
			err = identity.Set("instance_group_manager", d.Get("instance_group_manager").(string))
			if err != nil {
				return fmt.Errorf("Error setting instance_group_manager: %s", err)
			}
		}
		if v, ok := identity.GetOk("project"); !ok && v == "" {
			err = identity.Set("project", d.Get("project").(string))
			if err != nil {
				return fmt.Errorf("Error setting project: %s", err)
			}
		}
	} else {
		log.Printf("[DEBUG] (Read) identity not set: %s", err)
	}

	return nil
}

func resourceComputeBulkPerInstanceConfigDelete(d *schema.ResourceData, meta any) error {
	if d.Get("deletion_policy").(string) == "PREVENT" {
		return fmt.Errorf("cannot destroy ComputeBulkPerInstanceConfig without setting deletion_policy=\"DELETE\" and running `terraform apply`")
	}
	if d.Get("deletion_policy").(string) == "ABANDON" {
		log.Printf("[DEBUG] deletion_policy set to \"ABANDON\", removing BulkPerInstanceConfig %q from Terraform state without deletion", d.Id())
		return nil
	}
	config := meta.(*transport_tpg.Config)

	lockName, err := tpgresource.ReplaceVars(d, config, "instanceGroupManager/{{project}}/{{zone}}/{{instance_group_manager}}/{{name}}")
	if err != nil {
		return err
	}
	transport_tpg.MutexStore.Lock(lockName)
	defer transport_tpg.MutexStore.Unlock(lockName)

	perInstanceConfigs, err := expandRawPerInstanceConfigs(d.Get("per_instance_configs"))
	if err != nil {
		return nil
	}
	instanceNamesWithZones := getBulkPerInstanceConfigInstanceNamesWithZones(perInstanceConfigs, d, meta)

	err = callDeleteInstances(d, meta, instanceNamesWithZones)

	if err != nil {
		return err
	}

	instanceNames := getBulkPerInstanceConfigInstanceNames(perInstanceConfigs)
	err = waitForInstancesToBeDeleted(d, meta, instanceNames)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting BulkPerInstanceConfig %s", d.Id())
	return nil
}

func getListManagedInstancesResponse(d *schema.ResourceData, meta any) (map[string]any, error) {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return nil, err
	}

	url, err := tpgresource.ReplaceVars(d, config, transport_tpg.BaseUrl(Product, config)+"projects/{{project}}/zones/{{zone}}/instanceGroupManagers/{{instance_group_manager}}/listManagedInstances")
	if err != nil {
		return nil, err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return nil, fmt.Errorf("Error fetching project for BulkPerInstanceConfig: %s", err)
	}

	headers := make(http.Header)
	listManagedInstancesResponse, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "POST",
		Project:   project,
		RawURL:    url,
		UserAgent: userAgent,
		Headers:   headers,
	})
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] Finished reading ComputeBulkPerInstanceConfig %q", d.Id())

	return mapListManagedInstancesResponseToBulkPerInstanceConfig(d, listManagedInstancesResponse)
}

func callCreateInstances(d *schema.ResourceData, meta any, perInstanceConfigs []map[string]any) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	createInstancesRequest := map[string]any{
		"instances": perInstanceConfigs,
	}
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, transport_tpg.BaseUrl(Product, config)+"projects/{{project}}/zones/{{zone}}/instanceGroupManagers/{{instance_group_manager}}/createInstances")
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for BulkPerInstanceConfig: %s", err)
	}

	headers := make(http.Header)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "POST",
		Project:   project,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      createInstancesRequest,
		Timeout:   d.Timeout(schema.TimeoutCreate),
		Headers:   headers,
	})
	if err != nil {
		return fmt.Errorf("Error creating instances: %s", err)
	}

	err = ComputeOperationWaitTime(config, res, project, "Creating instances", userAgent, d.Timeout(schema.TimeoutCreate))

	log.Printf("[DEBUG] Finished creating instances %#v", res)

	return nil
}

func callDeleteInstances(d *schema.ResourceData, meta any, instanceNamesWithZones []string) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	url, err := tpgresource.ReplaceVars(d, config, transport_tpg.BaseUrl(Product, config)+"projects/{{project}}/zones/{{zone}}/instanceGroupManagers/{{instance_group_manager}}/deleteInstances")
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for BulkPerInstanceConfig: %s", err)
	}

	obj := map[string]any{
		"instances": instanceNamesWithZones,
	}

	headers := make(http.Header)

	log.Printf("[DEBUG] Deleting instances %#v", instanceNamesWithZones)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "POST",
		Project:   project,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutDelete),
		Headers:   headers,
	})

	if err != nil {
		return fmt.Errorf("Error deleting instances %#v: %s", instanceNamesWithZones, err)
	}
	err = ComputeOperationWaitTime(
		config, res, project, "Deleting instances", userAgent,
		d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return fmt.Errorf("Error deleting instances %#v: %s", instanceNamesWithZones, err)
	}
	return nil
}

func waitForInstancesToBeDeleted(d *schema.ResourceData, meta any, instanceNames []string) error {
	retryConf := retry.StateChangeConf{
		Pending:      []string{"deleting"},
		Target:       []string{"deleted"},
		Refresh:      checkIfInstancesDeleted(d, meta, instanceNames),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		PollInterval: time.Duration(10) * time.Second,
	}
	_, err := retryConf.WaitForState()
	if err != nil {
		return err
	}
	return nil
}

func checkIfInstancesDeleted(d *schema.ResourceData, meta any, instanceNames []string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		listManagedInstancesResponse, err := getListManagedInstancesResponse(d, meta)
		if err != nil {
			log.Printf("[WARNING] Error in fetching managed instances: %s\n", err)
			return nil, "error", err
		}

		// No managed instances in instance group manager
		if listManagedInstancesResponse == nil {
			return true, "deleted", nil
		}
		var allInstanceNames []string
		for _, perInstanceConfig := range listManagedInstancesResponse["perInstanceConfigs"].([]any) {
			name := perInstanceConfig.(map[string]any)["name"].(string)
			allInstanceNames = append(allInstanceNames, name)
		}

		var instancesNotDeleted []string
		for _, name := range allInstanceNames {
			if slices.Contains(instanceNames, name) {
				instancesNotDeleted = append(instancesNotDeleted, name)
			}
		}
		if len(instancesNotDeleted) == 0 {
			return true, "deleted", nil
		}
		return nil, "deleting", nil
	}
}

func getBulkPerInstanceConfigInstanceNamesWithZones(perInstanceConfigs []map[string]any, d *schema.ResourceData, meta any) []string {
	config := meta.(*transport_tpg.Config)
	var instanceNames []string
	for _, perInstanceConfig := range perInstanceConfigs {
		if instanceName, ok := perInstanceConfig["name"].(string); ok && instanceName != "" {
			fullName, _ := tpgresource.ReplaceVars(d, config, "zones/{{zone}}/instances/"+instanceName)
			instanceNames = append(instanceNames, fullName)
		}
	}
	return instanceNames
}

func getBulkPerInstanceConfigInstanceNames(perInstanceConfigs []map[string]any) []string {
	var instanceNames []string
	for _, perInstanceConfig := range perInstanceConfigs {
		if instanceName, ok := perInstanceConfig["name"].(string); ok && instanceName != "" {
			instanceNames = append(instanceNames, instanceName)
		}
	}
	return instanceNames
}

func expandRawPerInstanceConfigs(rawPerInstanceConfigs any) ([]map[string]any, error) {
	rawPerInstanceConfigs = rawPerInstanceConfigs.(*schema.Set).List()
	if rawPerInstanceConfigs == nil {
		return nil, nil
	}
	l := rawPerInstanceConfigs.([]any)
	req := make([]map[string]any, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]any)
		transformed := make(map[string]any)

		transformedName := original["name"]
		if val := reflect.ValueOf(transformedName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["name"] = transformedName
		}

		req = append(req, transformed)
	}
	return req, nil
}

func flattenNestedComputeBulkPerInstanceConfigPerInstanceConfigs(v any) any {
	if v == nil {
		return v
	}
	l := v.([]any)
	transformed := schema.NewSet(schema.HashResource(computeBulkPerInstanceConfigPerInstanceConfigsSchema()), []any{})
	for _, raw := range l {
		original := raw.(map[string]any)
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed.Add(map[string]any{
			"name": original["name"],
		})
	}
	return transformed
}

func mapListManagedInstancesResponseToBulkPerInstanceConfig(d *schema.ResourceData, res map[string]any) (map[string]any, error) {
	managedInstances, ok := res["managedInstances"]
	if !ok {
		return nil, nil
	}

	var perInstanceConfigsList []any
	for _, managedInstanceRaw := range managedInstances.([]any) {
		managedInstance := managedInstanceRaw.(map[string]any)
		instanceName, ok := managedInstance["name"].(string)
		if !ok {
			continue
		}

		config := map[string]any{"name": instanceName}
		perInstanceConfigsList = append(perInstanceConfigsList, config)
	}

	result := map[string]any{
		"perInstanceConfigs": perInstanceConfigsList,
		"name":               d.Get("name"),
	}
	return result, nil
}

func ResourceComputeBulkPerInstanceConfigFlatten(d *schema.ResourceData, perInstanceConfig map[string]any) error {
	var err error

	if err = d.Set("name", perInstanceConfig["name"]); err != nil {
		return fmt.Errorf("Error reading BulkPerInstanceConfig: %s", err)
	}
	if err = d.Set("per_instance_configs", flattenNestedComputeBulkPerInstanceConfigPerInstanceConfigs(perInstanceConfig["perInstanceConfigs"])); err != nil {
		return fmt.Errorf("Error reading BulkPerInstanceConfig: %s", err)
	}

	return nil
}
