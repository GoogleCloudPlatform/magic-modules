package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceMonitoringServiceCloudEndpoints() *schema.Resource {
	endptSchema := map[string]*schema.Schema{
		"service": {
			Type:        schema.TypeString,
			Required:    true,
			Description: `Cloud Endpoints service. Learn more at https://cloud.google.com/endpoints.`,
		},
	}
	filter := `cloud_endpoints.service="{{service}}"`

	return dataSourceMonitoringServiceType(endptSchema, filter, dataSourceMonitoringServiceCloudEndpointsRead)
}

func dataSourceMonitoringServiceCloudEndpointsRead(res map[string]interface{}, d *schema.ResourceData, meta interface{}) error {
	var cloudEndpoints map[string]interface{}
	if v, ok := res["cloud_endpoints"]; ok {
		cloudEndpoints = v.(map[string]interface{})
	}
	if len(cloudEndpoints) == 0 {
		return nil
	}

	if err := d.Set("service", cloudEndpoints["service"]); err != nil {
		return err
	}
	return nil
}
