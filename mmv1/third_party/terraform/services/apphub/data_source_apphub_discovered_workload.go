package apphub

import (
    "fmt"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-provider-google/google/tpgresource"
    transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceApphubDiscoveredWorkload() *schema.Resource {
    return &schema.Resource{
        Read: dataSourceApphubDiscoveredWorkloadRead,
        Schema: map[string]*schema.Schema{
        		"project":{
        		    Type: schema.TypeString,
        		    Optional: true,
        		 },
                        "location":{
                            Type: schema.TypeString,
                            Required: true,
                        },
                        "workload_uri": {
                            Type: schema.TypeString,
                            Required: true,
                        },
                        "discovered_workload": {
                                Type:     schema.TypeList,
                                Computed: true,
                                Elem: &schema.Resource{
                                        Schema: map[string]*schema.Schema{
                                                "name": {
                                                        Type:     schema.TypeString,
                                                        Computed: true,
                                                },
                                                "workload_reference": {
                                                        Type:     schema.TypeList,
                                                        Computed: true,
                                                        Elem: &schema.Resource{
                                                            Schema: map[string]*schema.Schema{
                                                                "uri":{
                                                                    Type: schema.TypeString,
                                                                    Computed: true,
                                                                },
                                                            },
                                                        },
                                                },
                                                "workload_properties": {
                                                        Type:     schema.TypeList,
                                                        Computed: true,
                                                        Elem: &schema.Resource{
                                                            Schema: map[string]*schema.Schema{
                                                                "gcp_project":{
                                                                    Type: schema.TypeString,
                                                                    Computed: true,
                                                                },
                                                                "location":{
                                                                    Type: schema.TypeString,
                                                                    Computed: true,
                                                                },
                                                                "zone":{
                                                                    Type: schema.TypeString,
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

func dataSourceApphubDiscoveredWorkloadRead(d *schema.ResourceData, meta interface{}) error {
        config := meta.(*transport_tpg.Config)
        userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
        if err != nil {
                return err
        }

        url, err := tpgresource.ReplaceVars(d, config, "{{ApphubBasePath}}projects/{{project}}/locations/{{location}}/discoveredWorkloads:lookup?uri={{workload_uri}}" )
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
        })

        if err != nil {
            return transport_tpg.HandleDataSourceNotFoundError(err, d, fmt.Sprintf("ApphubDiscoveredWorkload %q", d.Id()), url)
        }
        
        if err := d.Set("discovered_workload", flattenApphubDiscoveredWorkload(res["discoveredWorkload"], d, config)); err != nil {
                return fmt.Errorf("Error setting discovered workload: %s", err)
        }
        
        id, err := tpgresource.ReplaceVars(d, config, "{{workload_uri}}")
        if err != nil {
        	return err
        }
        d.SetId(id)
        
        return nil

}

func flattenApphubDiscoveredWorkload(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
        if v == nil {
                return nil
        }
        original := v.(map[string]interface{})
        if len(original) == 0 {
                return nil
        }
        
        transformed := make(map[string]interface{})
        transformed["name"] = flattenApphubDiscoveredWorkloadDataName(original["name"], d, config)
        transformed["workload_reference"] =  flattenApphubWorkloadReference(original["workloadReference"], d, config)
        transformed["workload_properties"] = flattenApphubWorkloadProperties(original["workloadProperties"], d, config)
        return []interface{}{transformed}
}


func flattenApphubWorkloadReference(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
        if v == nil {
                return nil
        }
        original := v.(map[string]interface{})
        if len(original) == 0 {
                return nil
        }
        transformed := make(map[string]interface{})
        transformed["uri"] = flattenApphubDiscoveredWorkloadDataUri(original["uri"], d, config)
        return []interface{}{transformed}
}

func flattenApphubWorkloadProperties(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
        if v == nil {
                return nil
        }
        original := v.(map[string]interface{})
        if len(original) == 0 {
                return nil
        }
        transformed := make(map[string]interface{})
        transformed["gcp_project"] = flattenApphubDiscoveredWorkloadDataGcpProject(original["gcpProject"], d, config)
        transformed["location"] = flattenApphubDiscoveredWorkloadDataLocation(original["location"], d, config)
        transformed["zone"] = flattenApphubDiscoveredWorkloadDataZone(original["zone"], d, config)
        return []interface{}{transformed}
}

func flattenApphubDiscoveredWorkloadDataName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
        return v
}

func flattenApphubDiscoveredWorkloadDataUri(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
        return v
}

func flattenApphubDiscoveredWorkloadDataGcpProject(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
        return v
}

func flattenApphubDiscoveredWorkloadDataLocation(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
        return v
}

func flattenApphubDiscoveredWorkloadDataZone(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
        return v
}

