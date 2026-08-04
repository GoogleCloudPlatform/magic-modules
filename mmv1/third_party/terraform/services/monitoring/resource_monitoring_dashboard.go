package monitoring

import (
	"fmt"
	"maps"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func stripMonitoringDashboardComputedFields(m map[string]interface{}) map[string]interface{} {
	result := maps.Clone(m)
	delete(result, "etag")
	delete(result, "name")
	return result
}

func removeKeysAbsentFromConfig(state, config map[string]interface{}) map[string]interface{} {
	if state == nil {
		return nil
	}
	if config == nil {
		config = map[string]interface{}{}
	}

	result := make(map[string]interface{}, len(state))
	for k, stateValue := range state {
		configValue, exists := config[k]
		if !exists {
			continue
		}

		switch stateValue := stateValue.(type) {
		case map[string]interface{}:
			if configValue, ok := configValue.(map[string]interface{}); ok {
				result[k] = removeKeysAbsentFromConfig(stateValue, configValue)
			} else {
				result[k] = stateValue
			}
		case []interface{}:
			if configValue, ok := configValue.([]interface{}); ok {
				result[k] = removeKeysAbsentFromConfigSlice(stateValue, configValue)
			} else {
				result[k] = stateValue
			}
		default:
			result[k] = stateValue
		}
	}
	return result
}

func removeKeysAbsentFromConfigSlice(state, config []interface{}) []interface{} {
	result := make([]interface{}, len(state))
	for i, stateValue := range state {
		if i >= len(config) {
			result[i] = stateValue
			continue
		}
		stateMap, stateIsMap := stateValue.(map[string]interface{})
		configMap, configIsMap := config[i].(map[string]interface{})
		if stateIsMap && configIsMap {
			result[i] = removeKeysAbsentFromConfig(stateMap, configMap)
		} else {
			result[i] = stateValue
		}
	}
	return result
}

// The API omits zero-valued tile positions, while Terraform's jsonencode preserves them.
func stripMonitoringDashboardZeroTilePositions(dashboard map[string]interface{}) map[string]interface{} {
	result := maps.Clone(dashboard)
	mosaicLayout, ok := result["mosaicLayout"].(map[string]interface{})
	if !ok {
		return result
	}
	tiles, ok := mosaicLayout["tiles"].([]interface{})
	if !ok {
		return result
	}

	normalizedTiles := make([]interface{}, len(tiles))
	copy(normalizedTiles, tiles)
	for i, tile := range tiles {
		tile, ok := tile.(map[string]interface{})
		if !ok {
			continue
		}
		tile = maps.Clone(tile)
		if xPos, ok := tile["xPos"].(float64); ok && xPos == 0 {
			delete(tile, "xPos")
		}
		if yPos, ok := tile["yPos"].(float64); ok && yPos == 0 {
			delete(tile, "yPos")
		}
		normalizedTiles[i] = tile
	}

	mosaicLayout = maps.Clone(mosaicLayout)
	mosaicLayout["tiles"] = normalizedTiles
	result["mosaicLayout"] = mosaicLayout
	return result
}

func monitoringDashboardDiffSuppress(_, old, new string, _ *schema.ResourceData) bool {
	oldMap, err := structure.ExpandJsonFromString(old)
	if err != nil {
		return false
	}
	newMap, err := structure.ExpandJsonFromString(new)
	if err != nil {
		return false
	}

	oldMap = stripMonitoringDashboardComputedFields(oldMap)
	newMap = stripMonitoringDashboardComputedFields(newMap)
	oldMap = removeKeysAbsentFromConfig(oldMap, newMap)
	oldMap = stripMonitoringDashboardZeroTilePositions(oldMap)
	newMap = stripMonitoringDashboardZeroTilePositions(newMap)

	return reflect.DeepEqual(oldMap, newMap)
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
			tpgresource.DefaultProviderDeletionPolicy("DELETE"),
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
			//UDP schema start
			"deletion_policy": tpgresource.DeletionPolicySchemaEntry("DELETE"),
			//UDP schema end
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

	url, err := tpgresource.ReplaceVars(d, config, transport_tpg.BaseUrl(Product, config)+"v1/projects/{{project}}/dashboards")
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

	url := transport_tpg.BaseUrl(Product, config) + "v1/" + d.Id()

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

	str, err := structure.FlattenJsonToString(res)
	if err != nil {
		return fmt.Errorf("Error reading Dashboard: %s", err)
	}
	if err = d.Set("dashboard_json", str); err != nil {
		return fmt.Errorf("Error reading Dashboard: %s", err)
	}

	if err := tpgresource.DeletionPolicyReadDefault(d, config, "DELETE"); err != nil {
		return err
	}

	return nil
}

func resourceMonitoringDashboardUpdate(d *schema.ResourceData, meta interface{}) error {

	if tpgresource.DeletionPolicyPreUpdate(d, ResourceMonitoringDashboard) {
		return ResourceMonitoringDashboard().Read(d, meta)
	}

	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	o, n := d.GetChange("dashboard_json")
	oObj, err := structure.ExpandJsonFromString(o.(string))
	if err != nil {
		return err
	}
	nObj, err := structure.ExpandJsonFromString(n.(string))
	if err != nil {
		return err
	}

	nObj["etag"] = oObj["etag"]

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	url := transport_tpg.BaseUrl(Product, config) + "v1/" + d.Id()
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

	if ok, err := tpgresource.DeletionPolicyPreDelete(d); err != nil {
		return err
	} else if ok {
		return nil
	}

	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url := transport_tpg.BaseUrl(Product, config) + "v1/" + d.Id()

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

func init() {
	registry.Schema{
		Name:        "google_monitoring_dashboard",
		ProductName: "monitoring",
		Type:        registry.SchemaTypeResource,
		Schema:      ResourceMonitoringDashboard(),
	}.Register()
}
