package monitoring

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// removeComputedKeys removes keys from the old configuration that don't exist in the new configuration.
// This prevents spurious diffs when the API adds computed fields that weren't in the original user config.
func removeComputedKeys(old map[string]interface{}, new map[string]interface{}) map[string]interface{} {
	if old == nil {
		return old
	}
	if new == nil {
		new = make(map[string]interface{})
	}

	for k, oldVal := range old {
		newVal, exists := new[k]

		if !exists {
			delete(old, k)
			continue
		}

		if oldMap, okOld := oldVal.(map[string]interface{}); okOld {
			if newMap, okNew := newVal.(map[string]interface{}); okNew {
				old[k] = removeComputedKeys(oldMap, newMap)
			}
			continue
		}

		if oldSlice, okOld := oldVal.([]interface{}); okOld {
			if newSlice, okNew := newVal.([]interface{}); okNew {
				for i := range oldSlice {
					if i < len(newSlice) {
						if oldElem, okOldElem := oldSlice[i].(map[string]interface{}); okOldElem {
							if newElem, okNewElem := newSlice[i].(map[string]interface{}); okNewElem {
								oldSlice[i] = removeComputedKeys(oldElem, newElem)
							}
						}
					}
				}
			}
			continue
		}
	}

	return old
}

// apiDefaultFields are fields that the Monitoring API adds as defaults when not specified.
// These fields are normalized away during diff suppression to prevent spurious diffs.
// Only these specific fields are normalized; other fields are compared as-is.
var apiDefaultFields = map[string]bool{
	// Empty string fields from dashboardFilters and other locations
	"labelKey":       true,
	"stringValue":    true,
	"legendTemplate": true,
	"label":          true,
	"unitOverride":   true,
	// Boolean fields that API sets to false
	"showLegend":         true,
	"outputFullDuration": true,
	// Position fields that API sets to 0
	"xPos": true,
	"yPos": true,
}

// normalizeDefaults recursively removes API default values from dashboard JSON.
// Only fields in apiDefaultFields are normalized, preventing overly broad removal.
// This allows comparison between user configs and API responses without spurious diffs.
func normalizeDefaults(obj interface{}) interface{} {
	switch v := obj.(type) {
	case map[string]interface{}:
		normalized := make(map[string]interface{})
		for k, val := range v {
			// Skip default values for known API fields
			if apiDefaultFields[k] {
				switch tv := val.(type) {
				case string:
					if tv == "" {
						continue
					}
				case bool:
					if !tv {
						continue
					}
				case float64:
					if tv == 0 {
						continue
					}
				case []interface{}:
					if len(tv) == 0 {
						continue
					}
				}
			}

			// Recursively process nested structures
			switch tv := val.(type) {
			case map[string]interface{}:
				normalized[k] = normalizeDefaults(tv)
			case []interface{}:
				normalizedArray := make([]interface{}, len(tv))
				for i, elem := range tv {
					normalizedArray[i] = normalizeDefaults(elem)
				}
				normalized[k] = normalizedArray
			default:
				normalized[k] = val
			}
		}
		return normalized
	case []interface{}:
		normalizedArray := make([]interface{}, len(v))
		for i, elem := range v {
			normalizedArray[i] = normalizeDefaults(elem)
		}
		return normalizedArray
	default:
		return obj
	}
}

func monitoringDashboardDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	oldMap, err := structure.ExpandJsonFromString(old)
	if err != nil {
		return false
	}
	newMap, err := structure.ExpandJsonFromString(new)
	if err != nil {
		return false
	}

	oldMap = removeComputedKeys(oldMap, newMap)
	oldNormalized := normalizeDefaults(oldMap)
	newNormalized := normalizeDefaults(newMap)

	// Compare as JSON strings after normalization to suppress spurious diffs
	oldJSON, _ := json.Marshal(oldNormalized)
	newJSON, _ := json.Marshal(newNormalized)
	return string(oldJSON) == string(newJSON)
}

func ResourceMonitoringDashboard() *schema.Resource {
	return &schema.Resource{
		Create: resourceMonitoringDashboardCreate,
		Read:   resourceMonitoringDashboardRead,
		Update: resourceMonitoringDashboardUpdate,
		Delete: resourceMonitoringDashboardDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMonitoringDashboardImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Update: schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(4 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			tpgresource.DefaultProviderProject,
		),

		Schema: map[string]*schema.Schema{
			"dashboard_json": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: monitoringDashboardDiffSuppress,
				StateFunc: func(v interface{}) string {
					json, _ := structure.NormalizeJsonString(v)
					return json
				},
				Description: `The JSON representation of a dashboard, following the format at https://cloud.google.com/monitoring/api/ref_v3/rest/v1/projects.dashboards.`,
			},
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The ID of the project in which the resource belongs. If it is not provided, the provider project is used.`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceMonitoringDashboardCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj, err := structure.ExpandJsonFromString(d.Get("dashboard_json").(string))
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{MonitoringBasePath}}v1/projects/{{project}}/dashboards")
	if err != nil {
		return err
	}
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:               config,
		Method:               "POST",
		Project:              project,
		RawURL:               url,
		UserAgent:            userAgent,
		Body:                 obj,
		Timeout:              d.Timeout(schema.TimeoutCreate),
		ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.IsMonitoringConcurrentEditError},
	})
	if err != nil {
		return fmt.Errorf("Error creating Dashboard: %s", err)
	}

	name, ok := res["name"]
	if !ok {
		return fmt.Errorf("Create response didn't contain critical fields. Create may not have succeeded.")
	}
	d.SetId(name.(string))

	return resourceMonitoringDashboardRead(d, config)
}

func resourceMonitoringDashboardRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url := config.MonitoringBasePath + "v1/" + d.Id()

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:               config,
		Method:               "GET",
		Project:              project,
		RawURL:               url,
		UserAgent:            userAgent,
		ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.IsMonitoringConcurrentEditError},
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("MonitoringDashboard %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting Dashboard: %s", err)
	}

	// Remove system-managed fields that change on every update or are derived from the ID
	if res != nil {
		delete(res, "etag")
		delete(res, "name")
	}

	str, err := structure.FlattenJsonToString(res)
	if err != nil {
		return fmt.Errorf("Error reading Dashboard: %s", err)
	}
	if err = d.Set("dashboard_json", str); err != nil {
		return fmt.Errorf("Error reading Dashboard: %s", err)
	}

	return nil
}

func resourceMonitoringDashboardUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	_, n := d.GetChange("dashboard_json")
	nObj, err := structure.ExpandJsonFromString(n.(string))
	if err != nil {
		return err
	}

	// Fetch current dashboard to get the latest etag
	url := config.MonitoringBasePath + "v1/" + d.Id()
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	currentDashboard, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:               config,
		Method:               "GET",
		Project:              project,
		RawURL:               url,
		UserAgent:            userAgent,
		ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.IsMonitoringConcurrentEditError},
	})
	if err != nil {
		return fmt.Errorf("Error fetching Dashboard for update: %s", err)
	}

	// Preserve etag from current API state for update request
	if etag, ok := currentDashboard["etag"]; ok {
		nObj["etag"] = etag
	}

	_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:               config,
		Method:               "PATCH",
		Project:              project,
		RawURL:               url,
		UserAgent:            userAgent,
		Body:                 nObj,
		Timeout:              d.Timeout(schema.TimeoutUpdate),
		ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.IsMonitoringConcurrentEditError},
	})
	if err != nil {
		return fmt.Errorf("Error updating Dashboard %q: %s", d.Id(), err)
	}

	return resourceMonitoringDashboardRead(d, config)
}

func resourceMonitoringDashboardDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url := config.MonitoringBasePath + "v1/" + d.Id()

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:               config,
		Method:               "DELETE",
		Project:              project,
		RawURL:               url,
		UserAgent:            userAgent,
		Timeout:              d.Timeout(schema.TimeoutDelete),
		ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.IsMonitoringConcurrentEditError},
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("MonitoringDashboard %q", d.Id()))
	}

	return nil
}

func resourceMonitoringDashboardImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)

	// current import_formats can't import fields with forward slashes in their value
	parts, err := tpgresource.GetImportIdQualifiers([]string{"projects/(?P<project>[^/]+)/dashboards/(?P<id>[^/]+)", "(?P<id>[^/]+)"}, d, config, d.Id())
	if err != nil {
		return nil, err
	}

	if err := d.Set("project", parts["project"]); err != nil {
		return nil, fmt.Errorf("Error setting project: %s", err)
	}
	d.SetId(fmt.Sprintf("projects/%s/dashboards/%s", parts["project"], parts["id"]))

	return []*schema.ResourceData{d}, nil
}
