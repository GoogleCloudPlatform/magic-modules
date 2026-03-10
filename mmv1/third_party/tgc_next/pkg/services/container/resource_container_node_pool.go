package container

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/tpgresource"
)

// ContainerNodePoolAssetType is the CAI asset type name for container node pool.
const ContainerNodePoolAssetType string = "container.googleapis.com/NodePool"

// ContainerNodePoolSchemaName is the TF resource schema name for container node pool.
const ContainerNodePoolSchemaName string = "google_container_node_pool"

var schemaBlueGreenSettings = &schema.Schema{
	Type:     schema.TypeList,
	Optional: true,
	Computed: true,
	MaxItems: 1,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"standard_rollout_policy": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: `Standard rollout policy is the default policy for blue-green.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"batch_percentage": {
							Type:         schema.TypeFloat,
							Optional:     true,
							Computed:     true,
							Description:  `Percentage of the blue pool nodes to drain in a batch.`,
							ValidateFunc: validation.FloatBetween(0.0, 1.0),
						},
						"batch_node_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: `Number of blue nodes to drain in a batch.`,
						},
						"batch_soak_duration": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `Soak time after each batch gets drained.`,
						},
					},
				},
			},
			"node_pool_soak_duration": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Time needed after draining entire blue pool. After this period, blue pool will be cleaned up.`,
			},
		},
	},
	Description: `Settings for BlueGreen node pool upgrade.`,
}

var schemaNodePool = map[string]*schema.Schema{
	"autoscaling": {
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: `Configuration required by cluster autoscaler to adjust the size of the node pool to the current cluster usage.`,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"min_node_count": {
					Type:         schema.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntAtLeast(0),
					Description:  `Minimum number of nodes per zone in the node pool. Must be >=0 and <= max_node_count. Cannot be used with total limits.`,
				},

				"max_node_count": {
					Type:         schema.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntAtLeast(0),
					Description:  `Maximum number of nodes per zone in the node pool. Must be >= min_node_count. Cannot be used with total limits.`,
				},

				"total_min_node_count": {
					Type:         schema.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntAtLeast(0),
					Description:  `Minimum number of all nodes in the node pool. Must be >=0 and <= total_max_node_count. Cannot be used with per zone limits.`,
				},

				"total_max_node_count": {
					Type:         schema.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntAtLeast(0),
					Description:  `Maximum number of all nodes in the node pool. Must be >= total_min_node_count. Cannot be used with per zone limits.`,
				},

				"location_policy": {
					Type:         schema.TypeString,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.StringInSlice([]string{"BALANCED", "ANY"}, false),
					Description:  `Location policy specifies the algorithm used when scaling-up the node pool. "BALANCED" - Is a best effort policy that aims to balance the sizes of available zones. "ANY" - Instructs the cluster autoscaler to prioritize utilization of unused reservations, and reduces preemption risk for Spot VMs.`,
				},
			},
		},
	},

	"placement_policy": {
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: `Specifies the node placement policy`,
		ForceNew:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:        schema.TypeString,
					Required:    true,
					Description: `Type defines the type of placement policy`,
				},
				"policy_name": {
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
					Description: `If set, refers to the name of a custom resource policy supplied by the user. The resource policy must be in the same project and region as the node pool. If not found, InvalidArgument error is returned.`,
				},
				"tpu_topology": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: `The TPU topology like "2x4" or "2x2x2". https://cloud.google.com/kubernetes-engine/docs/concepts/plan-tpus#topology`,
				},
			},
		},
	},

	"queued_provisioning": {
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: `Specifies the configuration of queued provisioning`,
		ForceNew:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"enabled": {
					Type:        schema.TypeBool,
					Required:    true,
					ForceNew:    true,
					Description: `Whether nodes in this node pool are obtainable solely through the ProvisioningRequest API`,
				},
			},
		},
	},

	"max_pods_per_node": {
		Type:        schema.TypeInt,
		Optional:    true,
		ForceNew:    true,
		Computed:    true,
		Description: `The maximum number of pods per node in this node pool. Note that this does not work on node pools which are "route-based" - that is, node pools belonging to clusters that do not have IP Aliasing enabled.`,
	},

	"node_locations": {
		Type:        schema.TypeSet,
		Optional:    true,
		Computed:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: `The list of zones in which the node pool's nodes should be located. Nodes must be in the region of their regional cluster or in the same region as their cluster's zone for zonal clusters. If unspecified, the cluster-level node_locations will be used.`,
	},

	"upgrade_settings": {
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    true,
		MaxItems:    1,
		Description: `Specify node upgrade settings to change how many nodes GKE attempts to upgrade at once. The number of nodes upgraded simultaneously is the sum of max_surge and max_unavailable. The maximum number of nodes upgraded simultaneously is limited to 20.`,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"max_surge": {
					Type:         schema.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntAtLeast(0),
					Description:  `The number of additional nodes that can be added to the node pool during an upgrade. Increasing max_surge raises the number of nodes that can be upgraded simultaneously. Can be set to 0 or greater.`,
				},

				"max_unavailable": {
					Type:         schema.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntAtLeast(0),
					Description:  `The number of nodes that can be simultaneously unavailable during an upgrade. Increasing max_unavailable raises the number of nodes that can be upgraded in parallel. Can be set to 0 or greater.`,
				},

				"strategy": {
					Type:         schema.TypeString,
					Optional:     true,
					Default:      "SURGE",
					ValidateFunc: validation.StringInSlice([]string{"SURGE", "BLUE_GREEN"}, false),
					Description:  `Update strategy for the given nodepool.`,
				},

				"blue_green_settings": schemaBlueGreenSettings,
			},
		},
	},

	"initial_node_count": {
		Type:        schema.TypeInt,
		Optional:    true,
		ForceNew:    true,
		Computed:    true,
		Description: `The initial number of nodes for the pool. In regional or multi-zonal clusters, this is the number of nodes per zone. Changing this will force recreation of the resource.`,
	},

	"instance_group_urls": {
		Type:        schema.TypeList,
		Computed:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: `The resource URLs of the managed instance groups associated with this node pool.`,
	},

	"managed_instance_group_urls": {
		Type:        schema.TypeList,
		Computed:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: `List of instance group URLs which have been assigned to this node pool.`,
	},

	"management": {
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    true,
		MaxItems:    1,
		Description: `Node management configuration, wherein auto-repair and auto-upgrade is configured.`,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"auto_repair": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     true,
					Description: `Whether the nodes will be automatically repaired. Enabled by default.`,
				},

				"auto_upgrade": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     true,
					Description: `Whether the nodes will be automatically upgraded. Enabled by default.`,
				},
			},
		},
	},

	"name": {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		ForceNew:    true,
		Description: `The name of the node pool. If left blank, Terraform will auto-generate a unique name.`,
	},

	"name_prefix": {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		ForceNew:    true,
		Description: `Creates a unique name for the node pool beginning with the specified prefix. Conflicts with name.`,
	},

	"node_config": schemaNodeConfig(),

	"node_count": {
		Type:         schema.TypeInt,
		Optional:     true,
		Computed:     true,
		ValidateFunc: validation.IntAtLeast(0),
		Description:  `The number of nodes per instance group. This field can be used to update the number of nodes per instance group but should not be used alongside autoscaling.`,
	},

	"version": {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: `The Kubernetes version for the nodes in this pool. Note that if this field and auto_upgrade are both specified, they will fight each other for what the node version should be, so setting both is highly discouraged. While a fuzzy version can be specified, it's recommended that you specify explicit versions as Terraform will see spurious diffs when fuzzy versions are used. See the google_container_engine_versions data source's version_prefix field to approximate fuzzy versions in a Terraform-compatible way.`,
	},

	"network_config": {
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    true,
		MaxItems:    1,
		Description: `Networking configuration for this NodePool. If specified, it overrides the cluster-level defaults.`,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{

				"create_pod_range": {
					Type:        schema.TypeBool,
					Optional:    true,
					ForceNew:    true,
					Description: `Whether to create a new range for pod IPs in this node pool. Defaults are provided for pod_range and pod_ipv4_cidr_block if they are not specified.`,
				},
				"enable_private_nodes": {
					Type:        schema.TypeBool,
					Optional:    true,
					Computed:    true,
					Description: `Whether nodes have internal IP addresses only.`,
				},
				"pod_range": {
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
					Computed:    true,
					Description: `The ID of the secondary range for pod IPs. If create_pod_range is true, this ID is used for the new range. If create_pod_range is false, uses an existing secondary range with this ID.`,
				},
				"pod_ipv4_cidr_block": {
					Type:             schema.TypeString,
					Optional:         true,
					ForceNew:         true,
					Computed:         true,
					DiffSuppressFunc: tpgresource.CidrOrSizeDiffSuppress,
					Description:      `The IP address range for pod IPs in this node pool. Only applicable if create_pod_range is true. Set to blank to have a range chosen with the default size. Set to /netmask (e.g. /14) to have a range chosen with a specific netmask. Set to a CIDR notation (e.g. 10.96.0.0/14) to pick a specific range to use.`,
				},
				"additional_node_network_configs": {
					Type:        schema.TypeList,
					Optional:    true,
					Computed:    true,
					ForceNew:    true,
					Description: `We specify the additional node networks for this node pool using this list. Each node network corresponds to an additional interface`,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"network": {
								Type:        schema.TypeString,
								Optional:    true,
								Computed:    true,
								ForceNew:    true,
								Description: `Name of the VPC where the additional interface belongs.`,
							},
							"subnetwork": {
								Type:        schema.TypeString,
								Optional:    true,
								Computed:    true,
								ForceNew:    true,
								Description: `Name of the subnetwork where the additional interface belongs.`,
							},
						},
					},
				},
				"additional_pod_network_configs": {
					Type:        schema.TypeList,
					Optional:    true,
					ForceNew:    true,
					Description: `We specify the additional pod networks for this node pool using this list. Each pod network corresponds to an additional alias IP range for the node`,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"subnetwork": {
								Type:        schema.TypeString,
								Optional:    true,
								ForceNew:    true,
								Description: `Name of the subnetwork where the additional pod network belongs.`,
							},
							"secondary_pod_range": {
								Type:        schema.TypeString,
								Optional:    true,
								ForceNew:    true,
								Description: `The name of the secondary range on the subnet which provides IP address for this pod range.`,
							},
							"max_pods_per_node": {
								Type:        schema.TypeInt,
								Optional:    true,
								ForceNew:    true,
								Computed:    true,
								Description: `The maximum number of pods per node which use this pod network.`,
							},
						},
					},
				},
				"pod_cidr_overprovision_config": {
					Type:        schema.TypeList,
					Optional:    true,
					Computed:    true,
					ForceNew:    true,
					MaxItems:    1,
					Description: `Configuration for node-pool level pod cidr overprovision. If not set, the cluster level setting will be inherited`,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"disabled": {
								Type:     schema.TypeBool,
								Required: true,
							},
						},
					},
				},
				"network_performance_config": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: `Network bandwidth tier configuration.`,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"total_egress_bandwidth_tier": {
								Type:        schema.TypeString,
								Required:    true,
								Description: `Specifies the total network bandwidth tier for the NodePool. [Valid values](https://cloud.google.com/kubernetes-engine/docs/reference/rest/v1/projects.locations.clusters.nodePools#NodePool.Tier) include: "TIER_1" and "TIER_UNSPECIFIED".`,
							},
						},
					},
				},
				"subnetwork": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: `The subnetwork path for the node pool. Format: projects/{project}/regions/{region}/subnetworks/{subnetwork} . If the cluster is associated with multiple subnetworks, the subnetwork for the node pool is picked based on the IP utilization during node pool creation and is immutable.`,
				},
			},
		},
	},

	"node_drain_config": {
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    true,
		Description: `Node drain configuration for this NodePool.`,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"respect_pdb_during_node_pool_deletion": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: `Whether to respect PodDisruptionBudget policy during node pool deletion.`,
				},
			},
		},
	},
}

func ResourceContainerNodePool() *schema.Resource {
	return &schema.Resource{
		Schema: tpgresource.MergeSchemas(
			schemaNodePool,
			map[string]*schema.Schema{
				"project": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					ForceNew:    true,
					Description: `The ID of the project in which to create the node pool. If blank, the provider-configured project will be used.`,
				},
				"cluster": {
					Type:        schema.TypeString,
					Required:    true,
					ForceNew:    true,
					Description: `The cluster to create the node pool for. Cluster must be present in location provided for zonal clusters.`,
				},
				"location": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					ForceNew:    true,
					Description: `The location (region or zone) of the cluster.`,
				},
				"operation": {
					Type:     schema.TypeString,
					Computed: true,
				},
			}),
	}
}
