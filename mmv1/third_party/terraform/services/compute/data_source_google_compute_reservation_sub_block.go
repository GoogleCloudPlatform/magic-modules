package compute

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleComputeReservationSubBlock() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeReservationSubBlockRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the reservation sub-block.",
			},
			"reservation_block": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the parent reservation block.",
			},
			"reservation": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the parent reservation.",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The zone where the reservation sub-block resides.",
			},
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The project in which the resource belongs.",
			},
			"kind": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the resource. Always compute#reservationSubBlock for reservation sub-blocks.",
			},
			"resource_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for the resource.",
			},
			"creation_timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation timestamp in RFC3339 text format.",
			},
			"self_link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server-defined fully-qualified URL for this resource.",
			},
			"self_link_with_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server-defined URL for this resource with the resource id.",
			},
			"sub_block_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of hosts that are allocated in this reservation sub-block.",
			},
			"in_use_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of instances that are currently in use on this reservation sub-block.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the reservation sub-block.",
			},
			"reservation_sub_block_maintenance": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Maintenance information for this reservation sub-block.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"maintenance_ongoing_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of hosts in the sub-block that have ongoing maintenance.",
						},
						"maintenance_pending_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of hosts in the sub-block that have pending maintenance.",
						},
						"scheduling_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of maintenance for the reservation.",
						},
						"subblock_infra_maintenance_ongoing_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of subblock Infrastructure that has ongoing maintenance.",
						},
						"subblock_infra_maintenance_pending_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of subblock Infrastructure that has pending maintenance.",
						},
						"instance_maintenance_ongoing_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of instances that have ongoing maintenance.",
						},
						"instance_maintenance_pending_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of instances that have pending maintenance.",
						},
					},
				},
			},
			"physical_topology": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The physical topology of the reservation sub-block.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The cluster name of the reservation sub-block.",
						},
						"block": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The hash of the capacity block within the cluster.",
						},
						"sub_block": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The hash of the capacity sub-block within the capacity block.",
						},
					},
				},
			},
			"health_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Health information for the reservation sub-block.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"health_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The health status of the reservation sub-block.",
						},
						"healthy_host_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of healthy hosts in the reservation sub-block.",
						},
						"degraded_host_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of degraded hosts in the reservation sub-block.",
						},
						"healthy_infra_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of healthy infrastructure in the reservation sub-block.",
						},
						"degraded_infra_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of degraded infrastructure in the reservation sub-block.",
						},
					},
				},
			},
		},
	}
}

func dataSourceGoogleComputeReservationSubBlockRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	zone, err := tpgresource.GetZone(d, config)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	reservationBlock := d.Get("reservation_block").(string)
	reservation := d.Get("reservation").(string)

	url := fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/zones/%s/reservations%%2F%s%%2FreservationBlocks%%2F%s/reservationSubBlocks/%s", project, zone, reservation, reservationBlock, name)

	log.Printf("[DEBUG] URL  %s ", url)

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return fmt.Errorf("Error reading ReservationSubBlock: %s", err)
	}

	if res == nil {
		return fmt.Errorf("ReservationSubBlock %s not found", name)
	}

	// Flatten the resource field if it exists
	if resource, ok := res["resource"]; ok {
		if resourceMap, ok := resource.(map[string]interface{}); ok {
			res = resourceMap
		}
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}

	if err := d.Set("kind", res["kind"]); err != nil {
		return fmt.Errorf("Error setting kind: %s", err)
	}

	if err := d.Set("resource_id", res["id"]); err != nil {
		return fmt.Errorf("Error setting resource_id: %s", err)
	}

	if err := d.Set("creation_timestamp", res["creationTimestamp"]); err != nil {
		return fmt.Errorf("Error setting creation_timestamp: %s", err)
	}

	if err := d.Set("self_link", res["selfLink"]); err != nil {
		return fmt.Errorf("Error setting self_link: %s", err)
	}

	if err := d.Set("self_link_with_id", res["selfLinkWithId"]); err != nil {
		return fmt.Errorf("Error setting self_link_with_id: %s", err)
	}

	if err := d.Set("sub_block_count", res["count"]); err != nil {
		return fmt.Errorf("Error setting count: %s", err)
	}

	if err := d.Set("in_use_count", res["inUseCount"]); err != nil {
		return fmt.Errorf("Error setting in_use_count: %s", err)
	}

	if err := d.Set("status", res["status"]); err != nil {
		return fmt.Errorf("Error setting status: %s", err)
	}

	if reservationSubBlockMaintenance, ok := res["reservationSubBlockMaintenance"].(map[string]interface{}); ok {
		maintenanceList := []map[string]interface{}{
			{
				"maintenance_ongoing_count":                reservationSubBlockMaintenance["maintenanceOngoingCount"],
				"maintenance_pending_count":                reservationSubBlockMaintenance["maintenancePendingCount"],
				"scheduling_type":                          reservationSubBlockMaintenance["schedulingType"],
				"subblock_infra_maintenance_ongoing_count": reservationSubBlockMaintenance["subblockInfraMaintenanceOngoingCount"],
				"subblock_infra_maintenance_pending_count": reservationSubBlockMaintenance["subblockInfraMaintenancePendingCount"],
				"instance_maintenance_ongoing_count":       reservationSubBlockMaintenance["instanceMaintenanceOngoingCount"],
				"instance_maintenance_pending_count":       reservationSubBlockMaintenance["instanceMaintenancePendingCount"],
			},
		}
		if err := d.Set("reservation_sub_block_maintenance", maintenanceList); err != nil {
			return fmt.Errorf("Error setting reservation_sub_block_maintenance: %s", err)
		}
	}

	if physicalTopology, ok := res["physicalTopology"].(map[string]interface{}); ok {
		topologyList := []map[string]interface{}{
			{
				"cluster":   physicalTopology["cluster"],
				"block":     physicalTopology["block"],
				"sub_block": physicalTopology["subBlock"],
			},
		}
		if err := d.Set("physical_topology", topologyList); err != nil {
			return fmt.Errorf("Error setting physical_topology: %s", err)
		}
	}

	if healthInfo, ok := res["healthInfo"].(map[string]interface{}); ok {
		healthList := []map[string]interface{}{
			{
				"health_status":        healthInfo["healthStatus"],
				"healthy_host_count":   healthInfo["healthyHostCount"],
				"degraded_host_count":  healthInfo["degradedHostCount"],
				"healthy_infra_count":  healthInfo["healthyInfraCount"],
				"degraded_infra_count": healthInfo["degradedInfraCount"],
			},
		}
		if err := d.Set("health_info", healthList); err != nil {
			return fmt.Errorf("Error setting health_info: %s", err)
		}
	}

	d.SetId(fmt.Sprintf("projects/%s/zones/%s/reservations/%s/reservationBlocks/%s/reservationSubBlocks/%s", project, zone, reservation, reservationBlock, name))
	return nil
}
