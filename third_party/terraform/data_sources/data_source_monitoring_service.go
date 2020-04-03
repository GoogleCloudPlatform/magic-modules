package google

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	neturl "net/url"
)

type monitoringServiceTypeFlattenFunc func(map[string]interface{}, *schema.ResourceData, interface{}) error

// dataSourceMonitoringServiceType creates a Datasource resource for a type of service. It takes
// - schema for identifying the service, specific to the type (AppEngine moduleId)
// - list query filter to filter a specific service (type, ID) from the list of services for a parent
// - typeFlattenF for reading the service-specific schema (typeSchema)
func dataSourceMonitoringServiceType(
	typeSchema map[string]*schema.Schema,
	listFilter string,
	flattenF monitoringServiceTypeFlattenFunc) *schema.Resource {
	// Convert resource schema to ds schema
	dsSchema := datasourceSchemaFromResourceSchema(resourceMonitoringService().Schema)

	// Add schema specific to the service type
	dsSchema = mergeSchemas(typeSchema, dsSchema)

	return &schema.Resource{
		Read:   dataSourceMonitoringServiceTypeReadFromList(listFilter, flattenF),
		Schema: dsSchema,
	}
}

// dataSourceMonitoringServiceRead returns a ReadFunc that calls service.list with proper filters
// to identify both the type of service and underlying service resource.
// It takes the list query filter (i.e. ?filter=$listFilter) and a ReadFunc to handle reading any type-specific schema.
func dataSourceMonitoringServiceTypeReadFromList(listFilter string, flattenTypeF monitoringServiceTypeFlattenFunc) schema.ReadFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*Config)

		project, err := getProject(d, config)
		if err != nil {
			return err
		}

		filters, err := replaceVars(d, config, listFilter)
		if err != nil {
			return err
		}

		listUrlTmpl := "{{MonitoringBasePath}}projects/{{project}}/services?filter=" + neturl.QueryEscape(filters)
		url, err := replaceVars(d, config, listUrlTmpl)
		if err != nil {
			return err
		}

		resp, err := sendRequest(config, "GET", project, url, nil, isMonitoringRetryableError)
		if err != nil {
			return fmt.Errorf("unable to list Monitoring Service for data source: %v", err)
		}

		v, ok := resp["services"]
		if !ok || v == nil {
			return fmt.Errorf("no Monitoring Services found for data source")
		}
		ls, ok := v.([]interface{})
		if !ok {
			return fmt.Errorf("no Monitoring Services found for data source")
		}
		if len(ls) == 0 {
			return fmt.Errorf("no Monitoring Services found for data source")
		}
		if len(ls) > 1 {
			return fmt.Errorf("more than one Monitoring Services with given identifer found")
		}
		res := ls[0].(map[string]interface{})
		log.Printf("[DEBUG] resp: %+v", res)

		// Keep the same as resource Read
		res, err = resourceMonitoringServiceDecoder(d, meta, res)
		if err != nil {
			return err
		}

		if err := d.Set("project", project); err != nil {
			return fmt.Errorf("Error reading Service: %s", err)
		}

		if err := d.Set("display_name", flattenMonitoringServiceDisplayName(res["displayName"], d, config)); err != nil {
			return fmt.Errorf("Error reading Service: %s", err)
		}
		if err := d.Set("telemetry", flattenMonitoringServiceTelemetry(res["telemetry"], d, config)); err != nil {
			return fmt.Errorf("Error reading Service: %s", err)
		}

		if err := flattenTypeF(res, d, config); err != nil {
			return fmt.Errorf("Error reading Service: %s", err)
		}

		name := flattenMonitoringServiceName(res["name"], d, config).(string)
		d.Set("name", name)
		d.SetId(name)

		log.Printf("[DEBUG] resp: %+v", d.Get("telemetry"))
		log.Printf("[DEBUG] resp: %+v", d.Get("telemetry.0.resource_name"))
		return nil
	}
}
