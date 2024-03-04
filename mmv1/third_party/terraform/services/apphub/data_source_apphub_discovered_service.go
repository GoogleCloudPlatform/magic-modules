package apphub

import (
    "fmt"
    "time"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-provider-google/google/tpgresource"
    transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceApphubDiscoveredService() *schema.Resource {
    return &schema.Resource{
        Read: dataSourceApphubDiscoveredServiceRead,
        Schema: map[string]*schema.Schema{
                        "location":{
                            Type: schema.TypeString,
                            Required: true,
                        },
                        "service_uri": {
                            Type: schema.TypeString,
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
                                                                "uri":{
                                                                    Type: schema.TypeString,
                                                                    Computed: true,
                                                                },
                                                                "path":{
                                                                    Type: schema.TypeString,
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

        res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
                Config:    config,
                Method:    "GET",
                Project:   billingProject,
                RawURL:    url,
                UserAgent: userAgent,
        })
        fmt.Println(res, err, "Test print")
        if err != nil {
                return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("ApphubDiscoveredService %q", d.Id()))
        }
        
        if err := d.Set("discovered_service", flattenApphubDiscoveredService(res["discoveredService"], d, config)); err != nil {
                return fmt.Errorf("Error setting discovered service: %s", err)
        }
        
        d.SetId(time.Now().UTC().String())
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
        transformed["gcp_project"] = flattenApphubDiscoveredServiceDataGcpProject(original["gcp_project"], d, config)
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


