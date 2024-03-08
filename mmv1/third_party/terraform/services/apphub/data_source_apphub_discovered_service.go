package apphub

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/googleapi"
)

func DataSourceApphubDiscoveredService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApphubDiscoveredServiceRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"location": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service_uri": {
				Type:     schema.TypeString,
				Required: true,
			},
			"discovered_service": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_reference": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"uri": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"path": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"service_properties": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"gcp_project": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"location": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"zone": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceApphubDiscoveredServiceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ApphubBasePath}}projects/{{project}}/locations/{{location}}/discoveredServices:lookup?uri={{service_uri}}")
	if err != nil {
		return err
	}

	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	var res map[string]interface{}

	err = transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() (rerr error) {
			res, rerr = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: userAgent,
			})
			return rerr
		},
		Timeout: d.Timeout(schema.TimeoutRead),
		ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{
			func(err error) (bool, string) {
				gerr, ok := err.(*googleapi.Error)
				if !ok {
					return false, ""
				}

				if gerr.Code == 404 {
					log.Printf("[DEBUG] Dismissed an error as retryable based on error code: %s", err)
					return true, fmt.Sprintf("Retryable error code %d", gerr.Code)
				}
				return false, ""
			},
		},
	})

	if err != nil {
		return transport_tpg.HandleDataSourceNotFoundError(err, d, fmt.Sprintf("ApphubDiscoveredService %q", d.Id()), url)
	}

	if err := d.Set("discovered_service", flattenApphubDiscoveredService(res["discoveredService"], d, config)); err != nil {
		return fmt.Errorf("Error setting discovered service: %s", err)
	}

	d.SetId(res["discoveredService"].(map[string]interface{})["name"].(string))

	return nil

}

func flattenApphubDiscoveredService(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["name"] = flattenApphubDiscoveredServiceDataName(original["name"], d, config)
	transformed["service_reference"] = flattenApphubServiceReference(original["serviceReference"], d, config)
	transformed["service_properties"] = flattenApphubServiceProperties(original["serviceProperties"], d, config)
	return []interface{}{transformed}
}

func flattenApphubServiceReference(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["uri"] = flattenApphubDiscoveredServiceDataUri(original["uri"], d, config)
	transformed["path"] = flattenApphubDiscoveredServiceDataPath(original["path"], d, config)
	return []interface{}{transformed}
}

func flattenApphubServiceProperties(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["gcp_project"] = flattenApphubDiscoveredServiceDataGcpProject(original["gcpProject"], d, config)
	transformed["location"] = flattenApphubDiscoveredServiceDataLocation(original["location"], d, config)
	transformed["zone"] = flattenApphubDiscoveredServiceDataZone(original["zone"], d, config)
	return []interface{}{transformed}
}

func flattenApphubDiscoveredServiceDataName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApphubDiscoveredServiceDataUri(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApphubDiscoveredServiceDataPath(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApphubDiscoveredServiceDataGcpProject(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApphubDiscoveredServiceDataLocation(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApphubDiscoveredServiceDataZone(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}
