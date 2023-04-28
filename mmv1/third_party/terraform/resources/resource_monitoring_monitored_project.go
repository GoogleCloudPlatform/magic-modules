package google

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceMonitoringMonitoredProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceMonitoringMonitoredProjectCreate,
		Read:   resourceMonitoringMonitoredProjectRead,
		Delete: resourceMonitoringMonitoredProjectDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMonitoringMonitoredProjectImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"metrics_scope": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Required. The resource name of the existing Metrics Scope that will monitor this project. Example: locations/global/metricsScopes/{SCOPING_PROJECT_ID_OR_NUMBER}`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Immutable. The resource name of the 'MonitoredProject'. On input, the resource name includes the scoping project ID and monitored project ID. On output, it contains the equivalent project numbers. Example: 'locations/global/metricsScopes/{SCOPING_PROJECT_ID_OR_NUMBER}/projects/{MONITORED_PROJECT_ID_OR_NUMBER}'`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Output only. The time when this 'MonitoredProject' was created.`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceMonitoringMonitoredProjectCreate(d *schema.ResourceData, meta any) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]any)
	nameProp, err := expandMonitoringMonitoredProjectName(d.Get("name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("name"); !isEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}

	obj, err = resourceMonitoringMonitoredProjectEncoder(d, meta, obj)
	if err != nil {
		return err
	}

	url, err := ReplaceVars(d, config, "{{MonitoringBasePath}}v1/locations/global/metricsScopes/{{metrics_scope}}/projects")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new MonitoredProject: %#v", obj)
	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := SendRequestWithTimeout(config, "POST", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating MonitoredProject: %s", err)
	}

	// Store the ID now
	id, err := ReplaceVars(d, config, "v1/locations/global/metricsScopes/{{metrics_scope}}/projects/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating MonitoredProject %q: %#v", d.Id(), res)

	return resourceMonitoringMonitoredProjectRead(d, meta)
}

func resourceMonitoringMonitoredProjectRead(d *schema.ResourceData, meta any) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := ReplaceVars(d, config, "{{MonitoringBasePath}}v1/locations/global/metricsScopes/{{metrics_scope}}")
	if err != nil {
		return err
	}

	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := SendRequest(config, "GET", billingProject, url, userAgent, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("MonitoringMonitoredProject %q", d.Id()))
	}

	name := d.Get("name").(string)
	name = GetResourceNameFromSelfLink(name)
	if name != "" {
		project, err := config.NewResourceManagerClient(userAgent).Projects.Get(name).Do()
		if err != nil {
			return err
		}
		name = strconv.FormatInt(project.ProjectNumber, 10)
	}
	if monitoredProjects, ok := res["monitoredProjects"].([]map[string]any); ok {
		for _, monitoredProject := range monitoredProjects {
			if strings.HasSuffix(monitoredProject["name"].(string), name) {
				if err := d.Set("create_time", flattenMonitoringMonitoredProjectCreateTime(monitoredProject["createTime"], d, config)); err != nil {
					return fmt.Errorf("Error reading MonitoredProject: %s", err)
				}
			}
		}
	}

	return nil
}

func resourceMonitoringMonitoredProjectDelete(d *schema.ResourceData, meta any) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	url, err := ReplaceVars(d, config, "{{MonitoringBasePath}}v1/locations/global/metricsScopes/{{metrics_scope}}/projects/{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]any
	log.Printf("[DEBUG] Deleting MonitoredProject %q", d.Id())

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := SendRequestWithTimeout(config, "DELETE", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return handleNotFoundError(err, d, "MonitoredProject")
	}

	log.Printf("[DEBUG] Finished deleting MonitoredProject %q: %#v", d.Id(), res)
	return nil
}

func resourceMonitoringMonitoredProjectImport(d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := ParseImportId([]string{
		"v1/locations/global/metricsScopes/(?P<metrics_scope>[^/]+)/projects/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := ReplaceVars(d, config, "v1/locations/global/metricsScopes/{{metrics_scope}}/projects/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenMonitoringMonitoredProjectName(v any, d *schema.ResourceData, config *transport_tpg.Config) any {
	return v
}

func flattenMonitoringMonitoredProjectCreateTime(v any, d *schema.ResourceData, config *transport_tpg.Config) any {
	return v
}

func expandMonitoringMonitoredProjectName(v any, d TerraformResourceData, config *transport_tpg.Config) (any, error) {
	return v, nil
}

func resourceMonitoringMonitoredProjectEncoder(d *schema.ResourceData, meta any, obj map[string]any) (map[string]any, error) {
	name := d.Get("name").(string)
	name = GetResourceNameFromSelfLink(name)
	metricsScope := d.Get("metrics_scope").(string)
	metricsScope = GetResourceNameFromSelfLink(metricsScope)
	obj["name"] = fmt.Sprintf("locations/global/metricsScopes/%s/projects/%s", metricsScope, name)
	return obj, nil
}
