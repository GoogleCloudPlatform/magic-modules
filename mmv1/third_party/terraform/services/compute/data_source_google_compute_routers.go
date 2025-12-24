package compute

import (
    "fmt"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-provider-google/google/tpgresource"
    transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleComputeRouters() *schema.Resource {
    return &schema.Resource{
        Read: dataSourceGoogleComputeRoutersRead,
        Schema: map[string]*schema.Schema{
            "project": {
                Type:     schema.TypeString,
                Optional: true,
                Computed: true,
            },
            "region": {
                Type:     schema.TypeString,
                Optional: true,
                Computed: true,
            },
            "routers": {
                Type:     schema.TypeList,
                Computed: true,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "name": {
                            Type:     schema.TypeString,
                            Computed: true,
                        },
                        "network": {
                            Type:     schema.TypeString,
                            Computed: true,
                        },
                        "description": {
                            Type:     schema.TypeString,
                            Computed: true,
                        },
                        "self_link": {
                            Type:     schema.TypeString,
                            Computed: true,
                        },
                        "bgp": {
                            Type:     schema.TypeList,
                            Computed: true,
                            Elem: &schema.Resource{
                                Schema: map[string]*schema.Schema{
                                    "asn": {
                                        Type:     schema.TypeInt,
                                        Computed: true,
                                    },
                                    "advertise_mode": {
                                        Type:     schema.TypeString,
                                        Computed: true,
                                    },
                                    "advertised_groups": {
                                        Type:     schema.TypeList,
                                        Computed: true,
                                        Elem:     &schema.Schema{Type: schema.TypeString},
                                    },
                                    "advertised_ip_ranges": {
                                        Type:     schema.TypeList,
                                        Computed: true,
                                        Elem: &schema.Resource{
                                            Schema: map[string]*schema.Schema{
                                                "range": {
                                                    Type:     schema.TypeString,
                                                    Computed: true,
                                                },
                                                "description": {
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
                },
            },
        },
    }
}

func dataSourceGoogleComputeRoutersRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*transport_tpg.Config)
    userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
    if err != nil {
        return err
    }

    project, err := tpgresource.GetProject(d, config)
    if err != nil {
        return err
    }

    region, err := tpgresource.GetRegion(d, config)
    if err != nil {
        return err
    }

    d.SetId(fmt.Sprintf("projects/%s/regions/%s", project, region))

    list, err := config.NewComputeClient(userAgent).Routers.List(project, region).Do()
    if err != nil {
        return fmt.Errorf("Error retrieving list of routers: %s", err)
    }

    var routers []map[string]interface{}
    for _, router := range list.Items {
        var bgpList []interface{}
        if router.Bgp != nil {
            var advertisedIpRanges []interface{}
            for _, ipRange := range router.Bgp.AdvertisedIpRanges {
                advertisedIpRanges = append(advertisedIpRanges, map[string]interface{}{
                    "range":       ipRange.Range,
                    "description": ipRange.Description,
                })
            }
            bgpList = []interface{}{
                map[string]interface{}{
                    "asn":                  router.Bgp.Asn,
                    "advertise_mode":       router.Bgp.AdvertiseMode,
                    "advertised_groups":    router.Bgp.AdvertisedGroups,
                    "advertised_ip_ranges": advertisedIpRanges,
                },
            }
        }

        routers = append(routers, map[string]interface{}{
            "name":        router.Name,
            "network":     router.Network,
            "description": router.Description,
            "self_link":   router.SelfLink,
            "bgp":         bgpList,
        })
    }

    if err := d.Set("routers", routers); err != nil {
        return fmt.Errorf("Error setting routers: %s", err)
    }

    if err := d.Set("project", project); err != nil {
        return fmt.Errorf("Error setting project: %s", err)
    }

    if err := d.Set("region", region); err != nil {
        return fmt.Errorf("Error setting region: %s", err)
    }

    return nil
}