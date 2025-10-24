package compute

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleComputeReservationBlock() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeReservationBlockRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the reservation block.",
			},
			"reservation": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the parent reservation.",
			},
			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The zone where the reservation block resides.",
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
				Description: "Type of the resource. Always compute#reservationBlock for reservation blocks.",
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
			"block_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of resources that are allocated in this reservation block.",
			},
			"in_use_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of instances that are currently in use on this reservation block.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the reservation block.",
			},
			"reservation_sub_block_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of reservation subBlocks associated with this reservation block.",
			},
			"reservation_sub_block_in_use_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of in-use reservation subBlocks associated with this reservation block.",
			},
			"reservation_maintenance": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Maintenance information for this reservation block.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"maintenance_ongoing_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of hosts in the block that have ongoing maintenance.",
						},
						"maintenance_pending_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of hosts in the block that have pending maintenance.",
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
				Description: "The physical topology of the reservation block.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The cluster name of the reservation block.",
						},
						"block": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The hash of the capacity block within the cluster.",
						},
					},
				},
			},
			"health_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Health information for the reservation block.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"health_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The health status of the reservation block.",
						},
						"healthy_sub_block_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of subBlocks that are healthy.",
						},
						"degraded_sub_block_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of subBlocks that are degraded.",
						},
					},
				},
			},
		},
	}
}

func dataSourceGoogleComputeReservationBlockRead(d *schema.ResourceData, meta interface{}) error {
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
	reservation := d.Get("reservation").(string)

	url := fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/zones/%s/reservations/%s/reservationBlocks/%s", project, zone, reservation, name)

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return fmt.Errorf("Error reading ReservationBlock: %s", err)
	}

	if res == nil {
		return fmt.Errorf("ReservationBlock %s not found", name)
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

	if err := d.Set("block_count", res["count"]); err != nil {
		return fmt.Errorf("Error setting count: %s", err)
	}

	if err := d.Set("in_use_count", res["inUseCount"]); err != nil {
		return fmt.Errorf("Error setting in_use_count: %s", err)
	}

	if err := d.Set("status", res["status"]); err != nil {
		return fmt.Errorf("Error setting status: %s", err)
	}

	if err := d.Set("reservation_sub_block_count", res["reservationSubBlockCount"]); err != nil {
		return fmt.Errorf("Error setting reservation_sub_block_count: %s", err)
	}

	if err := d.Set("reservation_sub_block_in_use_count", res["reservationSubBlockInUseCount"]); err != nil {
		return fmt.Errorf("Error setting reservation_sub_block_in_use_count: %s", err)
	}

	if reservationMaintenance, ok := res["reservationMaintenance"].(map[string]interface{}); ok {
		maintenanceList := []map[string]interface{}{
			{
				"maintenance_ongoing_count":                reservationMaintenance["maintenanceOngoingCount"],
				"maintenance_pending_count":                reservationMaintenance["maintenancePendingCount"],
				"scheduling_type":                          reservationMaintenance["schedulingType"],
				"subblock_infra_maintenance_ongoing_count": reservationMaintenance["subblockInfraMaintenanceOngoingCount"],
				"subblock_infra_maintenance_pending_count": reservationMaintenance["subblockInfraMaintenancePendingCount"],
				"instance_maintenance_ongoing_count":       reservationMaintenance["instanceMaintenanceOngoingCount"],
				"instance_maintenance_pending_count":       reservationMaintenance["instanceMaintenancePendingCount"],
			},
		}
		if err := d.Set("reservation_maintenance", maintenanceList); err != nil {
			return fmt.Errorf("Error setting reservation_maintenance: %s", err)
		}
	}

	if physicalTopology, ok := res["physicalTopology"].(map[string]interface{}); ok {
		topologyList := []map[string]interface{}{
			{
				"cluster": physicalTopology["cluster"],
				"block":   physicalTopology["block"],
			},
		}
		if err := d.Set("physical_topology", topologyList); err != nil {
			return fmt.Errorf("Error setting physical_topology: %s", err)
		}
	}

	if healthInfo, ok := res["healthInfo"].(map[string]interface{}); ok {
		healthList := []map[string]interface{}{
			{
				"health_status":            healthInfo["healthStatus"],
				"healthy_sub_block_count":  healthInfo["healthySubBlockCount"],
				"degraded_sub_block_count": healthInfo["degradedSubBlockCount"],
			},
		}
		if err := d.Set("health_info", healthList); err != nil {
			return fmt.Errorf("Error setting health_info: %s", err)
		}
	}

	d.SetId(fmt.Sprintf("projects/%s/zones/%s/reservations/%s/reservationBlocks/%s", project, zone, reservation, name))
	return nil
}
