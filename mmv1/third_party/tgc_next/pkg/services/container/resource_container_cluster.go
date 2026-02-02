package container

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"google.golang.org/api/container/v1"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/tpgresource"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/verify"
)

// ContainerClusterAssetType is the CAI asset type name for container cluster.
const ContainerClusterAssetType string = "container.googleapis.com/Cluster"

// ContainerClusterSchemaName is the TF resource schema name for container cluster.
const ContainerClusterSchemaName string = "google_container_cluster"

// Single-digit hour is equivalent to hour with leading zero e.g. suppress diff 1:00 => 01:00.
// Assume either value could be in either format.
func Rfc3339TimeDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	if (len(old) == 4 && "0"+old == new) || (len(new) == 4 && "0"+new == old) {
		return true
	}
	return false
}

var (
	instanceGroupManagerURL = regexp.MustCompile(fmt.Sprintf("projects/(%s)/zones/([a-z0-9-]*)/instanceGroupManagers/([^/]*)", verify.ProjectRegex))

	masterAuthorizedNetworksConfig = &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cidr_blocks": {
				Type: schema.TypeSet,
				// This should be kept Optional. Expressing the
				// parent with no entries and omitting the
				// parent entirely are semantically different.
				Optional:    true,
				Elem:        cidrBlockConfig,
				Description: `External networks that can access the Kubernetes cluster master through HTTPS.`,
			},
			"gcp_public_cidrs_access_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether Kubernetes master is accessible via Google Compute Engine Public IPs.`,
			},
			"private_endpoint_enforcement_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether authorized networks is enforced on the private endpoint or not. Defaults to false.`,
			},
		},
	}
	cidrBlockConfig = &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cidr_block": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsCIDRNetwork(0, 32),
				Description:  `External network that can access Kubernetes master through HTTPS. Must be specified in CIDR notation.`,
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Field for users to identify CIDR blocks.`,
			},
		},
	}

	ipAllocationCidrBlockFields = []string{"ip_allocation_policy.0.cluster_ipv4_cidr_block", "ip_allocation_policy.0.services_ipv4_cidr_block"}
	ipAllocationRangeFields     = []string{"ip_allocation_policy.0.cluster_secondary_range_name", "ip_allocation_policy.0.services_secondary_range_name"}

	addonsConfigKeys = []string{
		"addons_config.0.http_load_balancing",
		"addons_config.0.horizontal_pod_autoscaling",
		"addons_config.0.network_policy_config",
		"addons_config.0.cloudrun_config",
		"addons_config.0.gcp_filestore_csi_driver_config",
		"addons_config.0.dns_cache_config",
		"addons_config.0.gce_persistent_disk_csi_driver_config",
		"addons_config.0.gke_backup_agent_config",
		"addons_config.0.config_connector_config",
		"addons_config.0.gcs_fuse_csi_driver_config",
		"addons_config.0.stateful_ha_config",
		"addons_config.0.ray_operator_config",
		"addons_config.0.parallelstore_csi_driver_config",
		"addons_config.0.lustre_csi_driver_config",
		"addons_config.0.istio_config",
		"addons_config.0.kalm_config",
	}

	privateClusterConfigKeys = []string{
		"private_cluster_config.0.enable_private_endpoint",
		"private_cluster_config.0.enable_private_nodes",
		"private_cluster_config.0.master_ipv4_cidr_block",
		"private_cluster_config.0.private_endpoint_subnetwork",
		"private_cluster_config.0.master_global_access_config",
	}

	suppressDiffForAutopilot = schema.SchemaDiffSuppressFunc(func(k, oldValue, newValue string, d *schema.ResourceData) bool {
		if v, _ := d.Get("enable_autopilot").(bool); v {
			if k == "dns_config.0.additive_vpc_scope_dns_domain" {
				return false
			}
			if k == "dns_config.#" {
				if avpcDomain, _ := d.Get("dns_config.0.additive_vpc_scope_dns_domain").(string); avpcDomain != "" || d.HasChange("dns_config.0.additive_vpc_scope_dns_domain") {
					return false
				}
			}
			return true
		}
		return false
	})

	suppressDiffForPreRegisteredFleet = schema.SchemaDiffSuppressFunc(func(k, oldValue, newValue string, d *schema.ResourceData) bool {
		// Suppress if the cluster has been pre registered to fleet.
		if v, _ := d.Get("fleet.0.pre_registered").(bool); v {
			log.Printf("[DEBUG] fleet suppress pre_registered: %v\n", v)
			return true
		}
		// Suppress the addition of a fleet block (count changes 0 -> 1) if the "project" field being added is null or empty.
		if k == "fleet.#" && oldValue == "0" && newValue == "1" {
			// When transitioning from 0->1 blocks, d.Get/d.GetOk effectively reads the 'new' config value.
			projectVal, projectIsSet := d.GetOk("fleet.0.project")
			if !projectIsSet || projectVal.(string) == "" {
				log.Printf("[DEBUG] Suppressing diff for 'fleet.#' (0 -> 1) because fleet.0.project is null or empty in config.\n")
				return true
			}
		}
		return false
	})

	suppressDiffForConfidentialNodes = schema.SchemaDiffSuppressFunc(func(k, oldValue, newValue string, d *schema.ResourceData) bool {
		k = strings.Replace(k, "confidential_instance_type", "enabled", 1)
		if v, _ := d.Get(k).(bool); v {
			return oldValue == "SEV" && newValue == ""
		}
		return false
	})
)

// Defines default nodel pool settings for the entire cluster. These settings are
// overridden if specified on the specific NodePool object.
func clusterSchemaNodePoolDefaults() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    true,
		Description: `The default nodel pool settings for the entire cluster.`,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"node_config_defaults": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: `Subset of NodeConfig message that has defaults.`,
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"containerd_config":                      schemaContainerdConfig(),
							"gcfs_config":                            schemaGcfsConfig(),
							"insecure_kubelet_readonly_port_enabled": schemaInsecureKubeletReadonlyPortEnabled(),
							"logging_variant":                        schemaLoggingVariant(),
						},
					},
				},
			},
		},
	}
}

func rfc5545RecurrenceDiffSuppress(k, o, n string, d *schema.ResourceData) bool {
	// This diff gets applied in the cloud console if you specify
	// "FREQ=DAILY" in your config and add a maintenance exclusion.
	if o == "FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR,SA,SU" && n == "FREQ=DAILY" {
		return true
	}
	// Writing a full diff suppress for identical recurrences would be
	// complex and error-prone - it's not a big problem if a user
	// changes the recurrence and it's textually difference but semantically
	// identical.
	return false
}

// Has the field (e.g. enable_l4_ilb_subsetting and enable_fqdn_network_policy) been enabled before?
func isBeenEnabled(_ context.Context, old, new, _ interface{}) bool {
	if old == nil || new == nil {
		return false
	}

	// if subsetting is enabled, but is not now
	if old.(bool) && !new.(bool) {
		return true
	}

	return false
}

func suppressDiffForClusterDnsScope(k, o, n string, d *schema.ResourceData) bool {
	if o == "" && n == "DNS_SCOPE_UNSPECIFIED" {
		return true
	}
	return false
}

func ResourceContainerCluster() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the cluster, unique within the project and location.`,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)

					if len(value) > 40 {
						errors = append(errors, fmt.Errorf(
							"%q cannot be longer than 40 characters", k))
					}
					if !regexp.MustCompile("^[a-z0-9-]+$").MatchString(value) {
						errors = append(errors, fmt.Errorf(
							"%q can only contain lowercase letters, numbers and hyphens", k))
					}
					if !regexp.MustCompile("^[a-z]").MatchString(value) {
						errors = append(errors, fmt.Errorf(
							"%q must start with a letter", k))
					}
					if !regexp.MustCompile("[a-z0-9]$").MatchString(value) {
						errors = append(errors, fmt.Errorf(
							"%q must end with a number or a letter", k))
					}
					return
				},
			},

			"operation": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The location (region or zone) in which the cluster master will be created, as well as the default node location. If you specify a zone (such as us-central1-a), the cluster will be a zonal cluster with a single cluster master. If you specify a region (such as us-west1), the cluster will be a regional cluster with multiple masters spread across zones in the region, and with default node locations in those zones as well.`,
			},

			"node_locations": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of zones in which the cluster's nodes are located. Nodes must be in the region of their regional cluster or in the same region as their cluster's zone for zonal clusters. If this is specified for a zonal cluster, omit the cluster's zone.`,
			},

			"deletion_protection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: `When the field is set to true or unset in Terraform state, a terraform apply or terraform destroy that would delete the cluster will fail. When the field is set to false, deleting the cluster is allowed.`,
			},

			"addons_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: `The configuration for addons supported by GKE.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"http_load_balancing": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: addonsConfigKeys,
							MaxItems:     1,
							Description:  `The status of the HTTP (L7) load balancing controller addon, which makes it easy to set up HTTP load balancers for services in a cluster. It is enabled by default; set disabled = true to disable.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
						"horizontal_pod_autoscaling": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: addonsConfigKeys,
							MaxItems:     1,
							Description:  `The status of the Horizontal Pod Autoscaling addon, which increases or decreases the number of replica pods a replication controller has based on the resource usage of the existing pods. It ensures that a Heapster pod is running in the cluster, which is also used by the Cloud Monitoring service. It is enabled by default; set disabled = true to disable.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
						"network_policy_config": {
							Type:          schema.TypeList,
							Optional:      true,
							Computed:      true,
							AtLeastOneOf:  addonsConfigKeys,
							MaxItems:      1,
							Description:   `Whether we should enable the network policy addon for the master. This must be enabled in order to enable network policy for the nodes. To enable this, you must also define a network_policy block, otherwise nothing will happen. It can only be disabled if the nodes already do not have network policies enabled. Defaults to disabled; set disabled = false to enable.`,
							ConflictsWith: []string{"enable_autopilot"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
						"gcp_filestore_csi_driver_config": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: addonsConfigKeys,
							MaxItems:     1,
							Description:  `The status of the Filestore CSI driver addon, which allows the usage of filestore instance as volumes. Defaults to disabled for Standard clusters; set enabled = true to enable. It is enabled by default for Autopilot clusters; set enabled = true to enable it explicitly.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
						"cloudrun_config": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: addonsConfigKeys,
							MaxItems:     1,
							Description:  `The status of the CloudRun addon. It is disabled by default. Set disabled = false to enable.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"load_balancer_type": {
										Type:         schema.TypeString,
										ValidateFunc: validation.StringInSlice([]string{"LOAD_BALANCER_TYPE_INTERNAL"}, false),
										Optional:     true,
									},
								},
							},
						},
						"dns_cache_config": {
							Type:          schema.TypeList,
							Optional:      true,
							Computed:      true,
							AtLeastOneOf:  addonsConfigKeys,
							MaxItems:      1,
							Description:   `The status of the NodeLocal DNSCache addon. It is disabled by default. Set enabled = true to enable.`,
							ConflictsWith: []string{"enable_autopilot"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
						"gce_persistent_disk_csi_driver_config": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: addonsConfigKeys,
							MaxItems:     1,
							Description:  `Whether this cluster should enable the Google Compute Engine Persistent Disk Container Storage Interface (CSI) Driver. Set enabled = true to enable. The Compute Engine persistent disk CSI Driver is enabled by default on newly created clusters for the following versions: Linux clusters: GKE version 1.18.10-gke.2100 or later, or 1.19.3-gke.2100 or later.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
						"gke_backup_agent_config": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: addonsConfigKeys,
							MaxItems:     1,
							Description:  `The status of the Backup for GKE Agent addon. It is disabled by default. Set enabled = true to enable.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
						"gcs_fuse_csi_driver_config": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: addonsConfigKeys,
							MaxItems:     1,
							Description:  `The status of the GCS Fuse CSI driver addon, which allows the usage of gcs bucket as volumes. Defaults to disabled; set enabled = true to enable.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
						"parallelstore_csi_driver_config": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: addonsConfigKeys,
							MaxItems:     1,
							Description:  `The status of the Parallelstore CSI driver addon, which allows the usage of Parallelstore instances as volumes. Defaults to disabled; set enabled = true to enable.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
						"lustre_csi_driver_config": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: addonsConfigKeys,
							MaxItems:     1,
							Description:  `Configuration for the Lustre CSI driver. Defaults to disabled; set enabled = true to enable.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: `Whether the Lustre CSI driver is enabled for this cluster.`,
									},
									"enable_legacy_lustre_port": {
										Type:     schema.TypeBool,
										Optional: true,
										Description: `If set to true, the Lustre CSI driver will initialize LNet (the virtual network layer for Lustre kernel module) using port 6988.
										This flag is required to workaround a port conflict with the gke-metadata-server on GKE nodes.`,
									},
								},
							},
						},
						"istio_config": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: addonsConfigKeys,
							MaxItems:     1,
							Description:  `The status of the Istio addon.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: `The status of the Istio addon, which makes it easy to set up Istio for services in a cluster. It is disabled by default. Set disabled = false to enable.`,
									},
									"auth": {
										Type:     schema.TypeString,
										Optional: true,
										// We can't use a Terraform-level default because it won't be true when the block is disabled: true
										DiffSuppressFunc: tpgresource.EmptyOrDefaultStringSuppress("AUTH_NONE"),
										ValidateFunc:     validation.StringInSlice([]string{"AUTH_NONE", "AUTH_MUTUAL_TLS"}, false),
										Description:      `The authentication type between services in Istio. Available options include AUTH_MUTUAL_TLS.`,
									},
								},
							},
						},
						"kalm_config": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: addonsConfigKeys,
							MaxItems:     1,
							Description:  `Configuration for the KALM addon, which manages the lifecycle of k8s. It is disabled by default; Set enabled = true to enable.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
						"config_connector_config": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: addonsConfigKeys,
							MaxItems:     1,
							Description:  `The of the Config Connector addon.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
						"stateful_ha_config": {
							Type:          schema.TypeList,
							Optional:      true,
							Computed:      true,
							AtLeastOneOf:  addonsConfigKeys,
							MaxItems:      1,
							Description:   `The status of the Stateful HA addon, which provides automatic configurable failover for stateful applications. Defaults to disabled; set enabled = true to enable.`,
							ConflictsWith: []string{"enable_autopilot"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
						"ray_operator_config": {
							Type:         schema.TypeList,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: addonsConfigKeys,
							MaxItems:     3,
							Description:  `The status of the Ray Operator addon, which enabled management of Ray AI/ML jobs on GKE. Defaults to disabled; set enabled = true to enable.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"ray_cluster_logging_config": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										MaxItems:    1,
										Description: `The status of Ray Logging, which scrapes Ray cluster logs to Cloud Logging. Defaults to disabled; set enabled = true to enable.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Type:     schema.TypeBool,
													Required: true,
												},
											},
										},
									},
									"ray_cluster_monitoring_config": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										MaxItems:    1,
										Description: `The status of Ray Cluster monitoring, which shows Ray cluster metrics in Cloud Console. Defaults to disabled; set enabled = true to enable.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Type:     schema.TypeBool,
													Required: true,
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

			"cluster_autoscaling": {
				Type:     schema.TypeList,
				MaxItems: 1,
				// This field is Optional + Computed because we automatically set the
				// enabled value to false if the block is not returned in API responses.
				Optional:    true,
				Computed:    true,
				Description: `Per-cluster configuration of Node Auto-Provisioning with Cluster Autoscaler to automatically adjust the size of the cluster and create/delete node pools based on the current needs of the cluster's workload. See the guide to using Node Auto-Provisioning for more details.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:          schema.TypeBool,
							Optional:      true,
							Computed:      true,
							ConflictsWith: []string{"enable_autopilot"},
							Description:   `Whether node auto-provisioning is enabled. Resource limits for cpu and memory must be defined to enable node auto-provisioning.`,
						},
						"resource_limits": {
							Type:             schema.TypeList,
							Optional:         true,
							ConflictsWith:    []string{"enable_autopilot"},
							DiffSuppressFunc: suppressDiffForAutopilot,
							Description:      `Global constraints for machine resources in the cluster. Configuring the cpu and memory types is required if node auto-provisioning is enabled. These limits will apply to node pool autoscaling in addition to node auto-provisioning.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The type of the resource. For example, cpu and memory. See the guide to using Node Auto-Provisioning for a list of types.`,
									},
									"minimum": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: `Minimum amount of the resource in the cluster.`,
									},
									"maximum": {
										Type:         schema.TypeInt,
										Description:  `Maximum amount of the resource in the cluster.`,
										Required:     true,
										ValidateFunc: validation.IntAtLeast(1),
									},
								},
							},
						},
						"auto_provisioning_defaults": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: `Contains defaults for a node pool created by NAP.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"oauth_scopes": {
										Type:             schema.TypeList,
										Optional:         true,
										Computed:         true,
										Elem:             &schema.Schema{Type: schema.TypeString},
										DiffSuppressFunc: containerClusterAddedScopesSuppress,
										Description:      `Scopes that are used by NAP when creating node pools.`,
									},
									"service_account": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "default",
										Description: `The Google Cloud Platform Service Account to be used by the node VMs.`,
									},
									"disk_size": {
										Type:             schema.TypeInt,
										Optional:         true,
										Default:          100,
										Description:      `Size of the disk attached to each node, specified in GB. The smallest allowed disk size is 10GB.`,
										DiffSuppressFunc: suppressDiffForAutopilot,
										ValidateFunc:     validation.IntAtLeast(10),
									},
									"disk_type": {
										Type:             schema.TypeString,
										Optional:         true,
										Default:          "pd-standard",
										Description:      `Type of the disk attached to each node.`,
										DiffSuppressFunc: suppressDiffForAutopilot,
										ValidateFunc:     validation.StringInSlice([]string{"pd-standard", "pd-ssd", "pd-balanced"}, false),
									},
									"image_type": {
										Type:             schema.TypeString,
										Optional:         true,
										Default:          "COS_CONTAINERD",
										Description:      `The default image type used by NAP once a new node pool is being created.`,
										DiffSuppressFunc: suppressDiffForAutopilot,
										ValidateFunc:     validation.StringInSlice([]string{"COS_CONTAINERD", "COS", "UBUNTU_CONTAINERD", "UBUNTU"}, false),
									},
									"min_cpu_platform": {
										Type:             schema.TypeString,
										Optional:         true,
										DiffSuppressFunc: tpgresource.EmptyOrDefaultStringSuppress("automatic"),
										Description:      `Minimum CPU platform to be used by this instance. The instance may be scheduled on the specified or newer CPU platform. Applicable values are the friendly names of CPU platforms, such as Intel Haswell.`,
									},
									"boot_disk_kms_key": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: `The Customer Managed Encryption Key used to encrypt the boot disk attached to each node in the node pool.`,
									},
									"shielded_instance_config": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: `Shielded Instance options.`,
										MaxItems:    1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enable_secure_boot": {
													Type:        schema.TypeBool,
													Optional:    true,
													Default:     false,
													Description: `Defines whether the instance has Secure Boot enabled.`,
													AtLeastOneOf: []string{
														"cluster_autoscaling.0.auto_provisioning_defaults.0.shielded_instance_config.0.enable_secure_boot",
														"cluster_autoscaling.0.auto_provisioning_defaults.0.shielded_instance_config.0.enable_integrity_monitoring",
													},
												},
												"enable_integrity_monitoring": {
													Type:        schema.TypeBool,
													Optional:    true,
													Default:     true,
													Description: `Defines whether the instance has integrity monitoring enabled.`,
													AtLeastOneOf: []string{
														"cluster_autoscaling.0.auto_provisioning_defaults.0.shielded_instance_config.0.enable_secure_boot",
														"cluster_autoscaling.0.auto_provisioning_defaults.0.shielded_instance_config.0.enable_integrity_monitoring",
													},
												},
											},
										},
									},
									"management": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										MaxItems:    1,
										Description: `NodeManagement configuration for this NodePool.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"auto_upgrade": {
													Type:        schema.TypeBool,
													Optional:    true,
													Computed:    true,
													Description: `Specifies whether node auto-upgrade is enabled for the node pool. If enabled, node auto-upgrade helps keep the nodes in your node pool up to date with the latest release version of Kubernetes.`,
												},
												"auto_repair": {
													Type:        schema.TypeBool,
													Optional:    true,
													Computed:    true,
													Description: `Specifies whether the node auto-repair is enabled for the node pool. If enabled, the nodes in this node pool will be monitored and, if they fail health checks too many times, an automatic repair action will be triggered.`,
												},
												"upgrade_options": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: `Specifies the Auto Upgrade knobs for the node pool.`,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"auto_upgrade_start_time": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: `This field is set when upgrades are about to commence with the approximate start time for the upgrades, in RFC3339 text format.`,
															},
															"description": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: `This field is set when upgrades are about to commence with the description of the upgrade.`,
															},
														},
													},
												},
											},
										},
									},
									"upgrade_settings": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: `Specifies the upgrade settings for NAP created node pools`,
										Computed:    true,
										MaxItems:    1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"max_surge": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: `The maximum number of nodes that can be created beyond the current size of the node pool during the upgrade process.`,
												},
												"max_unavailable": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: `The maximum number of nodes that can be simultaneously unavailable during the upgrade process.`,
												},
												"strategy": {
													Type:         schema.TypeString,
													Optional:     true,
													Computed:     true,
													Description:  `Update strategy of the node pool.`,
													ValidateFunc: validation.StringInSlice([]string{"NODE_POOL_UPDATE_STRATEGY_UNSPECIFIED", "BLUE_GREEN", "SURGE"}, false),
												},
												"blue_green_settings": {
													Type:        schema.TypeList,
													Optional:    true,
													Computed:    true,
													MaxItems:    1,
													Description: `Settings for blue-green upgrade strategy.`,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"node_pool_soak_duration": {
																Type:     schema.TypeString,
																Optional: true,
																Computed: true,
																Description: `Time needed after draining entire blue pool. After this period, blue pool will be cleaned up.

																A duration in seconds with up to nine fractional digits, ending with 's'. Example: "3.5s".`,
															},
															"standard_rollout_policy": {
																Type:        schema.TypeList,
																Optional:    true,
																Computed:    true,
																MaxItems:    1,
																Description: `Standard policy for the blue-green upgrade.`,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"batch_percentage": {
																			Type:         schema.TypeFloat,
																			Optional:     true,
																			Computed:     true,
																			ValidateFunc: validation.FloatBetween(0.0, 1.0),
																			ExactlyOneOf: []string{
																				"cluster_autoscaling.0.auto_provisioning_defaults.0.upgrade_settings.0.blue_green_settings.0.standard_rollout_policy.0.batch_percentage",
																				"cluster_autoscaling.0.auto_provisioning_defaults.0.upgrade_settings.0.blue_green_settings.0.standard_rollout_policy.0.batch_node_count",
																			},
																			Description: `Percentage of the bool pool nodes to drain in a batch. The range of this field should be (0.0, 1.0].`,
																		},
																		"batch_node_count": {
																			Type:     schema.TypeInt,
																			Optional: true,
																			Computed: true,
																			ExactlyOneOf: []string{
																				"cluster_autoscaling.0.auto_provisioning_defaults.0.upgrade_settings.0.blue_green_settings.0.standard_rollout_policy.0.batch_percentage",
																				"cluster_autoscaling.0.auto_provisioning_defaults.0.upgrade_settings.0.blue_green_settings.0.standard_rollout_policy.0.batch_node_count",
																			},
																			Description: `Number of blue nodes to drain in a batch.`,
																		},
																		"batch_soak_duration": {
																			Type:     schema.TypeString,
																			Optional: true,
																			Default:  "0s",
																			Description: `Soak time after each batch gets drained.

																			A duration in seconds with up to nine fractional digits, ending with 's'. Example: "3.5s".`,
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
							},
						},
						"auto_provisioning_locations": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The list of Google Compute Engine zones in which the NodePool's nodes can be created by NAP.`,
						},
						"autoscaling_profile": {
							Type:             schema.TypeString,
							Default:          "BALANCED",
							Optional:         true,
							DiffSuppressFunc: suppressDiffForAutopilot,
							ValidateFunc:     validation.StringInSlice([]string{"BALANCED", "OPTIMIZE_UTILIZATION"}, false),
							Description:      `Configuration options for the Autoscaling profile feature, which lets you choose whether the cluster autoscaler should optimize for resource utilization or resource availability when deciding to remove nodes from a cluster. Can be BALANCED or OPTIMIZE_UTILIZATION. Defaults to BALANCED.`,
						},
						"default_compute_class_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Specifies whether default compute class behaviour is enabled. If enabled, cluster autoscaler will use Compute Class with name default for all the workloads, if not overriden.`,
						},
					},
				},
			},

			"cluster_ipv4_cidr": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  verify.OrEmpty(verify.ValidateRFC1918Network(8, 32)),
				ConflictsWith: []string{"ip_allocation_policy"},
				Description:   `The IP address range of the Kubernetes pods in this cluster in CIDR notation (e.g. 10.96.0.0/14). Leave blank to have one automatically chosen or specify a /14 block in 10.0.0.0/8. This field will only work for routes-based clusters, where ip_allocation_policy is not defined.`,
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: ` Description of the cluster.`,
			},

			"binary_authorization": {
				Type:             schema.TypeList,
				Optional:         true,
				DiffSuppressFunc: BinaryAuthorizationDiffSuppress,
				MaxItems:         1,
				Description:      "Configuration options for the Binary Authorization feature.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:          schema.TypeBool,
							Optional:      true,
							Deprecated:    "Deprecated in favor of evaluation_mode.",
							Description:   "Enable Binary Authorization for this cluster.",
							ConflictsWith: []string{"enable_autopilot", "binary_authorization.0.evaluation_mode"},
						},
						"evaluation_mode": {
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							ValidateFunc:  validation.StringInSlice([]string{"DISABLED", "PROJECT_SINGLETON_POLICY_ENFORCE"}, false),
							Description:   "Mode of operation for Binary Authorization policy evaluation.",
							ConflictsWith: []string{"binary_authorization.0.enabled"},
						},
					},
				},
			},

			"enable_kubernetes_alpha": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: `Whether to enable Kubernetes Alpha features for this cluster. Note that when this option is enabled, the cluster cannot be upgraded and will be automatically deleted after 30 days.`,
			},

			"enable_k8s_beta_apis": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: `Configuration for Kubernetes Beta APIs.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled_apis": {
							Type:        schema.TypeSet,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Enabled Kubernetes Beta APIs.`,
						},
					},
				},
			},

			"enable_tpu": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: `Whether to enable Cloud TPU resources in this cluster.`,
			},

			"enable_legacy_abac": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: `Whether the ABAC authorizer is enabled for this cluster. When enabled, identities in the system, including service accounts, nodes, and controllers, will have statically granted permissions beyond those provided by the RBAC configuration or IAM. Defaults to false.`,
			},

			"enable_shielded_nodes": {
				Type:          schema.TypeBool,
				Optional:      true,
				Default:       true,
				Description:   `Enable Shielded Nodes features on all nodes in this cluster. Defaults to true.`,
				ConflictsWith: []string{"enable_autopilot"},
			},

			"enable_autopilot": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: `Enable Autopilot for this cluster.`,
				// ConflictsWith: many fields, see https://cloud.google.com/kubernetes-engine/docs/concepts/autopilot-overview#comparison. The conflict is only set one-way, on other fields w/ this field.
			},

			"allow_net_admin": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Enable NET_ADMIN for this cluster.`,
			},

			"authenticator_groups_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: `Configuration for the Google Groups for GKE feature.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"security_group": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The name of the RBAC security group for use with Google security groups in Kubernetes RBAC. Group name must be in format gke-security-groups@yourdomain.com.`,
						},
					},
				},
			},

			"initial_node_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: `The number of nodes to create in this cluster's default node pool. In regional or multi-zonal clusters, this is the number of nodes per zone. Must be set if node_pool is not set. If you're using google_container_node_pool objects with no default node pool, you'll need to set this to a value of at least 1, alongside setting remove_default_node_pool to true.`,
			},

			"logging_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: `Logging configuration for the cluster.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_components": {
							Type:        schema.TypeList,
							Required:    true,
							Description: `GKE components exposing logs. Valid values include SYSTEM_COMPONENTS, APISERVER, CONTROLLER_MANAGER, KCP_CONNECTION, KCP_SSHD, KCP_HPA, SCHEDULER, and WORKLOADS.`,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{"SYSTEM_COMPONENTS", "APISERVER", "CONTROLLER_MANAGER", "KCP_CONNECTION", "KCP_SSHD", "KCP_HPA", "SCHEDULER", "WORKLOADS"}, false),
							},
						},
					},
				},
			},

			"logging_service": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"cluster_telemetry"},
				ValidateFunc:  validation.StringInSlice([]string{"logging.googleapis.com", "logging.googleapis.com/kubernetes", "none"}, false),
				Description:   `The logging service that the cluster should write logs to. Available options include logging.googleapis.com(Legacy Stackdriver), logging.googleapis.com/kubernetes(Stackdriver Kubernetes Engine Logging), and none. Defaults to logging.googleapis.com/kubernetes.`,
			},

			"maintenance_policy": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: `The maintenance policy to use for the cluster.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"daily_maintenance_window": {
							Type:     schema.TypeList,
							Optional: true,
							ExactlyOneOf: []string{
								"maintenance_policy.0.daily_maintenance_window",
								"maintenance_policy.0.recurring_window",
							},
							MaxItems:    1,
							Description: `Time window specified for daily maintenance operations. Specify start_time in RFC3339 format "HH:MM”, where HH : [00-23] and MM : [00-59] GMT.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_time": {
										Type:             schema.TypeString,
										Required:         true,
										ValidateFunc:     verify.ValidateRFC3339Time,
										DiffSuppressFunc: Rfc3339TimeDiffSuppress,
									},
									"duration": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"recurring_window": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							ExactlyOneOf: []string{
								"maintenance_policy.0.daily_maintenance_window",
								"maintenance_policy.0.recurring_window",
							},
							Description: `Time window for recurring maintenance operations.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_time": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: verify.ValidateRFC3339Date,
									},
									"end_time": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: verify.ValidateRFC3339Date,
									},
									"recurrence": {
										Type:             schema.TypeString,
										Required:         true,
										DiffSuppressFunc: rfc5545RecurrenceDiffSuppress,
									},
								},
							},
						},
						"maintenance_exclusion": {
							Type:        schema.TypeSet,
							Optional:    true,
							MaxItems:    20,
							Description: `Exceptions to maintenance window. Non-emergency maintenance should not occur in these windows.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"exclusion_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"start_time": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: verify.ValidateRFC3339Date,
									},
									"end_time": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: verify.ValidateRFC3339Date,
									},
									"exclusion_options": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: `Maintenance exclusion related options.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"end_time_behavior": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringInSlice([]string{"UNTIL_END_OF_SUPPORT"}, false),
													Description:  `The behavior of the exclusion end time.`,
												},
												"scope": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validation.StringInSlice([]string{"NO_UPGRADES", "NO_MINOR_UPGRADES", "NO_MINOR_OR_NODE_UPGRADES"}, false),
													Description:  `The scope of automatic upgrades to restrict in the exclusion window.`,
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

			"protect_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: `Enable/Disable Protect API features for the cluster.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"workload_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: `WorkloadConfig defines which actions are enabled for a cluster's workload configurations.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"audit_mode": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Sets which mode of auditing should be used for the cluster's workloads. Accepted values are DISABLED, BASIC.`,
									},
								},
							},
							AtLeastOneOf: []string{
								"protect_config.0.workload_config",
								"protect_config.0.workload_vulnerability_mode",
							},
						},
						"workload_vulnerability_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `Sets which mode to use for Protect workload vulnerability scanning feature. Accepted values are DISABLED, BASIC.`,
							AtLeastOneOf: []string{
								"protect_config.0.workload_config",
								"protect_config.0.workload_vulnerability_mode",
							},
						},
					},
				},
			},

			"security_posture_config": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Computed:    true,
				Description: `Defines the config needed to enable/disable features for the Security Posture API`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ValidateFunc:     validation.StringInSlice([]string{"DISABLED", "BASIC", "ENTERPRISE", "MODE_UNSPECIFIED"}, false),
							Description:      `Sets the mode of the Kubernetes security posture API's off-cluster features. Available options include DISABLED, BASIC, and ENTERPRISE.`,
							DiffSuppressFunc: tpgresource.EmptyOrDefaultStringSuppress("MODE_UNSPECIFIED"),
						},
						"vulnerability_mode": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ValidateFunc:     validation.StringInSlice([]string{"VULNERABILITY_DISABLED", "VULNERABILITY_BASIC", "VULNERABILITY_ENTERPRISE", "VULNERABILITY_MODE_UNSPECIFIED"}, false),
							Description:      `Sets the mode of the Kubernetes security posture API's workload vulnerability scanning. Available options include VULNERABILITY_DISABLED, VULNERABILITY_BASIC and VULNERABILITY_ENTERPRISE.`,
							DiffSuppressFunc: tpgresource.EmptyOrDefaultStringSuppress("VULNERABILITY_MODE_UNSPECIFIED"),
						},
					},
				},
			},
			"monitoring_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: `Monitoring configuration for the cluster.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_components": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: `GKE components exposing metrics. Valid values include SYSTEM_COMPONENTS, APISERVER, SCHEDULER, CONTROLLER_MANAGER, STORAGE, HPA, POD, DAEMONSET, DEPLOYMENT, STATEFULSET, KUBELET, CADVISOR, DCGM and JOBSET.`,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"managed_prometheus": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: `Configuration for Google Cloud Managed Services for Prometheus.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: `Whether or not the managed collection is enabled.`,
									},
									"auto_monitoring_config": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										MaxItems:    1,
										Description: `Configuration for GKE Workload Auto-Monitoring.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"scope": {
													Type:         schema.TypeString,
													Required:     true,
													Description:  `The scope of auto-monitoring.`,
													ValidateFunc: validation.StringInSlice([]string{"ALL", "NONE"}, false),
												},
											},
										},
									},
								},
							},
						},
						"advanced_datapath_observability_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: `Configuration of Advanced Datapath Observability features.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable_metrics": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: `Whether or not the advanced datapath metrics are enabled.`,
									},
									"enable_relay": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: `Whether or not Relay is enabled.`,
									},
								},
							},
						},
					},
				},
			},

			"notification_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: `The notification config for sending cluster upgrade notifications`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pubsub": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Description: `Notification config for Cloud Pub/Sub`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: `Whether or not the notification config is enabled`,
									},
									"topic": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The pubsub topic to push upgrade notifications to. Must be in the same project as the cluster. Must be in the format: projects/{project}/topics/{topic}.`,
									},
									"filter": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: `Allows filtering to one or more specific event types. If event types are present, those and only those event types will be transmitted to the cluster. Other types will be skipped. If no filter is specified, or no event types are present, all event types will be sent`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"event_type": {
													Type:        schema.TypeList,
													Required:    true,
													Description: `Can be used to filter what notifications are sent. Valid values include include UPGRADE_AVAILABLE_EVENT, UPGRADE_EVENT, SECURITY_BULLETIN_EVENT, and UPGRADE_INFO_EVENT`,
													Elem: &schema.Schema{
														Type:         schema.TypeString,
														ValidateFunc: validation.StringInSlice([]string{"UPGRADE_AVAILABLE_EVENT", "UPGRADE_EVENT", "SECURITY_BULLETIN_EVENT", "UPGRADE_INFO_EVENT"}, false),
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

			"confidential_nodes": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: `Configuration for the confidential nodes feature, which makes nodes run on confidential VMs. Warning: This configuration can't be changed (or added/removed) after cluster creation without deleting and recreating the entire cluster.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							ForceNew:    true,
							Description: `Whether Confidential Nodes feature is enabled for all nodes in this cluster.`,
						},
						"confidential_instance_type": {
							Type:             schema.TypeString,
							Optional:         true,
							ForceNew:         true,
							DiffSuppressFunc: suppressDiffForConfidentialNodes,
							Description:      `Defines the type of technology used by the confidential node.`,
							ValidateFunc:     validation.StringInSlice([]string{"SEV", "SEV_SNP", "TDX"}, false),
						},
					},
				},
			},

			"master_auth": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Computed:    true,
				Description: `The authentication information for accessing the Kubernetes master. Some values in this block are only returned by the API if your service account has permission to get credentials for your GKE cluster. If you see an unexpected diff unsetting your client cert, ensure you have the container.clusters.getCredentials permission.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_certificate_config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							ForceNew:    true,
							Description: `Whether client certificate authorization is enabled for this cluster.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"issue_client_certificate": {
										Type:        schema.TypeBool,
										Required:    true,
										ForceNew:    true,
										Description: `Whether client certificate authorization is enabled for this cluster.`,
									},
								},
							},
						},

						"client_certificate": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Base64 encoded public certificate used by clients to authenticate to the cluster endpoint.`,
						},

						"client_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: `Base64 encoded private key used by clients to authenticate to the cluster endpoint.`,
						},

						"cluster_ca_certificate": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Base64 encoded public certificate that is the root of trust for the cluster.`,
						},
					},
				},
			},

			"master_authorized_networks_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        masterAuthorizedNetworksConfig,
				Description: `The desired configuration options for master authorized networks. Omit the nested cidr_blocks attribute to disallow external access (except the cluster node IPs, which GKE automatically whitelists).`,
			},

			"min_master_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The minimum version of the master. GKE will auto-update the master to new versions, so this does not guarantee the current master version--use the read-only master_version field to obtain that. If unset, the cluster's version will be set by GKE to the version of the most recent official release (which is not necessarily the latest version).`,
			},

			"monitoring_service": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"cluster_telemetry"},
				ValidateFunc:  validation.StringInSlice([]string{"monitoring.googleapis.com", "monitoring.googleapis.com/kubernetes", "none"}, false),
				Description:   `The monitoring service that the cluster should write metrics to. Automatically send metrics from pods in the cluster to the Google Cloud Monitoring API. VM metrics will be collected by Google Compute Engine regardless of this setting Available options include monitoring.googleapis.com(Legacy Stackdriver), monitoring.googleapis.com/kubernetes(Stackdriver Kubernetes Engine Monitoring), and none. Defaults to monitoring.googleapis.com/kubernetes.`,
			},

			"network": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "default",
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      `The name or self_link of the Google Compute Engine network to which the cluster is connected. For Shared VPC, set this to the self link of the shared network.`,
			},

			"network_policy": {
				Type:             schema.TypeList,
				Optional:         true,
				MaxItems:         1,
				Description:      `Configuration options for the NetworkPolicy feature.`,
				ConflictsWith:    []string{"enable_autopilot"},
				DiffSuppressFunc: containerClusterNetworkPolicyDiffSuppress,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: `Whether network policy is enabled on the cluster.`,
						},
						"provider": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateFunc:     validation.StringInSlice([]string{"PROVIDER_UNSPECIFIED", "CALICO"}, false),
							DiffSuppressFunc: tpgresource.EmptyOrDefaultStringSuppress("PROVIDER_UNSPECIFIED"),
							Description:      `The selected network policy provider.`,
						},
					},
				},
			},

			"node_config": schemaNodeConfig(),

			"node_pool_defaults": clusterSchemaNodePoolDefaults(),

			"node_pool_auto_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: `Node pool configs that apply to all auto-provisioned node pools in autopilot clusters and node auto-provisioning enabled clusters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_kubelet_config": schemaNodePoolAutoConfigNodeKubeletConfig(),
						"linux_node_config":   schemaNodePoolAutoConfigLinuxNodeConfig(),
						"network_tags": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: `Collection of Compute Engine network tags that can be applied to a node's underlying VM instance.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tags": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: `List of network tags applied to auto-provisioned node pools.`,
									},
								},
							},
						},
						"resource_manager_tags": {
							Type:        schema.TypeMap,
							Optional:    true,
							Description: `A map of resource manager tags. Resource manager tag keys and values have the same definition as resource manager tags. Keys must be in the format tagKeys/{tag_key_id}, and values are in the format tagValues/456. The field is ignored (both PUT & PATCH) when empty.`,
						},
					},
				},
			},

			"node_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The Kubernetes version on the nodes. Must either be unset or set to the same value as min_master_version on create. Defaults to the default version set by GKE which is not necessarily the latest version. This only affects nodes in the default node pool. While a fuzzy version can be specified, it's recommended that you specify explicit versions as Terraform will see spurious diffs when fuzzy versions are used. See the google_container_engine_versions data source's version_prefix field to approximate fuzzy versions in a Terraform-compatible way. To update nodes in other node pools, use the version attribute on the node pool.`,
			},

			"pod_security_policy_config": {
				Type:             schema.TypeList,
				Optional:         true,
				Description:      `Configuration for the PodSecurityPolicy feature.`,
				MaxItems:         1,
				DiffSuppressFunc: podSecurityPolicyCfgSuppress,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: `Enable the PodSecurityPolicy controller for this cluster. If enabled, pods must be valid under a PodSecurityPolicy to be created.`,
						},
					},
				},
			},
			"pod_autoscaling": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: `PodAutoscaling is used for configuration of parameters for workload autoscaling`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hpa_profile": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"NONE", "PERFORMANCE"}, false),
							Description: `
								HPA Profile is used to configure the Horizontal Pod Autoscaler (HPA) profile for the cluster.
								Available options include:
								- NONE: Customers explicitly opt-out of HPA profiles.
								- PERFORMANCE: PERFORMANCE is used when customers opt-in to the performance HPA profile. In this profile we support a higher number of HPAs per cluster and faster metrics collection for workload autoscaling.
							`,
						},
					},
				},
			},
			"secret_manager_config": {
				Type:             schema.TypeList,
				Optional:         true,
				Description:      `Configuration for the Secret Manager feature.`,
				MaxItems:         1,
				DiffSuppressFunc: SecretManagerCfgSuppress,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: `Enable the Secret manager csi component.`,
						},
						"rotation_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: `Configuration for Secret Manager auto rotation.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: `Enable the Secret manager auto rotation.`,
									},
									"rotation_interval": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: `The interval between two consecutive rotations. Default rotation interval is 2 minutes`,
									},
								},
							},
						},
					},
				},
			},
			"secret_sync_config": {
				Type:             schema.TypeList,
				Optional:         true,
				Description:      `Configuration for the Sync as k8s secrets feature.`,
				MaxItems:         1,
				DiffSuppressFunc: SecretSyncCfgSuppress,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: `Enable the Sync as k8s secret add-on.`,
						},
						"rotation_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: `Configuration for Secret Sync auto rotation.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: `Enable the Secret sync auto rotation.`,
									},
									"rotation_interval": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: `The interval between two consecutive rotations. Default rotation interval is 2 minutes`,
									},
								},
							},
						},
					},
				},
			},
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The ID of the project in which the resource belongs. If it is not provided, the provider project is used.`,
			},

			"subnetwork": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      `The name or self_link of the Google Compute Engine subnetwork in which the cluster's instances are launched.`,
			},

			"self_link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Server-defined URL for the resource.`,
			},

			"endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The IP address of this cluster's Kubernetes master.`,
			},

			"master_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The current version of the master in the cluster. This may be different than the min_master_version set in the config if the master has been updated by GKE.`,
			},

			"services_ipv4_cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The IP address range of the Kubernetes services in this cluster, in CIDR notation (e.g. 1.2.3.4/29). Service addresses are typically put in the last /16 from the container CIDR.`,
			},

			"ip_allocation_policy": {
				Type:          schema.TypeList,
				MaxItems:      1,
				ForceNew:      true,
				Computed:      true,
				Optional:      true,
				ConflictsWith: []string{"cluster_ipv4_cidr"},
				Description:   `Configuration of cluster IP allocation for VPC-native clusters. Adding this block enables IP aliasing, making the cluster VPC-native instead of routes-based.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// GKE creates/deletes secondary ranges in VPC
						"cluster_ipv4_cidr_block": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ForceNew:         true,
							ConflictsWith:    ipAllocationRangeFields,
							DiffSuppressFunc: tpgresource.CidrOrSizeDiffSuppress,
							Description:      `The IP address range for the cluster pod IPs. Set to blank to have a range chosen with the default size. Set to /netmask (e.g. /14) to have a range chosen with a specific netmask. Set to a CIDR notation (e.g. 10.96.0.0/14) from the RFC-1918 private networks (e.g. 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16) to pick a specific range to use.`,
						},

						"services_ipv4_cidr_block": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ForceNew:         true,
							ConflictsWith:    ipAllocationRangeFields,
							DiffSuppressFunc: tpgresource.CidrOrSizeDiffSuppress,
							Description:      `The IP address range of the services IPs in this cluster. Set to blank to have a range chosen with the default size. Set to /netmask (e.g. /14) to have a range chosen with a specific netmask. Set to a CIDR notation (e.g. 10.96.0.0/14) from the RFC-1918 private networks (e.g. 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16) to pick a specific range to use.`,
						},

						// User manages secondary ranges manually
						"cluster_secondary_range_name": {
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							ForceNew:      true,
							ConflictsWith: ipAllocationCidrBlockFields,
							Description:   `The name of the existing secondary range in the cluster's subnetwork to use for pod IP addresses. Alternatively, cluster_ipv4_cidr_block can be used to automatically create a GKE-managed one.`,
						},

						"services_secondary_range_name": {
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							ForceNew:      true,
							ConflictsWith: ipAllocationCidrBlockFields,
							Description:   `The name of the existing secondary range in the cluster's subnetwork to use for service ClusterIPs. Alternatively, services_ipv4_cidr_block can be used to automatically create a GKE-managed one.`,
						},

						"stack_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "IPV4",
							ValidateFunc: validation.StringInSlice([]string{"IPV4", "IPV4_IPV6"}, false),
							Description:  `The IP Stack type of the cluster. Choose between IPV4 and IPV4_IPV6. Default type is IPV4 Only if not set`,
						},
						"pod_cidr_overprovision_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: `Configuration for cluster level pod cidr overprovision. Default is disabled=false.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
						"additional_pod_ranges_config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: `AdditionalPodRangesConfig is the configuration for additional pod secondary ranges supporting the ClusterUpdate message.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"pod_range_names": {
										Type:        schema.TypeSet,
										MinItems:    1,
										Required:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: `Name for pod secondary ipv4 range which has the actual range defined ahead.`,
									},
								},
							},
						},
						"additional_ip_ranges_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `AdditionalIPRangesConfig is the configuration for individual additional subnetworks attached to the cluster`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"subnetwork": {
										Type:             schema.TypeString,
										Required:         true,
										DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
										Description:      `Name of the subnetwork. This can be the full path of the subnetwork or just the name.`,
									},
									"pod_ipv4_range_names": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: `List of secondary ranges names within this subnetwork that can be used for pod IPs.`,
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"auto_ipam_config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: `AutoIpamConfig contains all information related to Auto IPAM.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: `The flag that enables Auto IPAM on this cluster.`,
									},
								},
							},
						},
						"network_tier_config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: `Used to determine the default network tier for external IP addresses on cluster resources, such as node pools and load balancers.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"network_tier": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Network tier configuration.`,
									},
								},
							},
						},
					},
				},
			},

			// Defaults to "VPC_NATIVE" during create only
			"networking_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"VPC_NATIVE", "ROUTES"}, false),
				Description:  `Determines whether alias IPs or routes will be used for pod IPs in the cluster. Defaults to VPC_NATIVE for new clusters.`,
			},

			"remove_default_node_pool": {
				Type:          schema.TypeBool,
				Optional:      true,
				Description:   `If true, deletes the default node pool upon cluster creation. If you're using google_container_node_pool resources with no default node pool, this should be set to true, alongside setting initial_node_count to at least 1.`,
				ConflictsWith: []string{"enable_autopilot"},
			},

			"control_plane_endpoints_config": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Computed:    true,
				Optional:    true,
				Description: `Configuration for all of the cluster's control plane endpoints. Currently supports only DNS endpoint configuration and disable IP endpoint. Other IP endpoint configurations are available in private_cluster_config.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dns_endpoint_config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: `DNS endpoint configuration.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"endpoint": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: `The cluster's DNS endpoint.`,
									},
									"allow_external_traffic": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: `Controls whether user traffic is allowed over this endpoint. Note that GCP-managed services may still use the endpoint even if this is false.`,
									},
									"enable_k8s_tokens_via_dns": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: `Controls whether the k8s token auth is allowed via dns.`,
									},
									"enable_k8s_certs_via_dns": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: `Controls whether the k8s certs auth is allowed via dns.`,
									},
								},
							},
						},
						"ip_endpoints_config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: `IP endpoint configuration.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: `Controls whether to allow direct IP access.`,
									},
								},
							},
						},
					},
				},
			},

			"private_cluster_config": {
				Type:             schema.TypeList,
				MaxItems:         1,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: containerClusterPrivateClusterConfigSuppress,
				Description:      `Configuration for private clusters, clusters with private nodes.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// enable_private_endpoint is orthogonal to private_endpoint_subnetwork.
						// User can create a private_cluster_config block without including
						// either one of those two fields. Both fields are optional.
						// At the same time, we use 'AtLeastOneOf' to prevent an empty block
						// like 'private_cluster_config{}'
						"enable_private_endpoint": {
							Type:             schema.TypeBool,
							Optional:         true,
							AtLeastOneOf:     privateClusterConfigKeys,
							DiffSuppressFunc: containerClusterPrivateClusterConfigSuppress,
							Description:      `When true, the cluster's private endpoint is used as the cluster endpoint and access through the public endpoint is disabled. When false, either endpoint can be used.`,
						},
						"enable_private_nodes": {
							Type:             schema.TypeBool,
							Optional:         true,
							ForceNew:         true,
							AtLeastOneOf:     privateClusterConfigKeys,
							DiffSuppressFunc: containerClusterPrivateClusterConfigSuppress,
							Description:      `Enables the private cluster feature, creating a private endpoint on the cluster. In a private cluster, nodes only have RFC 1918 private addresses and communicate with the master's private endpoint via private networking.`,
						},
						"master_ipv4_cidr_block": {
							Type:         schema.TypeString,
							Computed:     true,
							Optional:     true,
							ForceNew:     true,
							AtLeastOneOf: privateClusterConfigKeys,
							ValidateFunc: verify.OrEmpty(validation.IsCIDRNetwork(28, 28)),
							Description:  `The IP range in CIDR notation to use for the hosted master network. This range will be used for assigning private IP addresses to the cluster master(s) and the ILB VIP. This range must not overlap with any other ranges in use within the cluster's network, and it must be a /28 subnet. See Private Cluster Limitations for more details. This field only applies to private clusters, when enable_private_nodes is true.`,
						},
						"peering_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the peering between this cluster and the Google owned VPC.`,
						},
						"private_endpoint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The internal IP address of this cluster's master endpoint.`,
						},
						"private_endpoint_subnetwork": {
							Type:             schema.TypeString,
							Optional:         true,
							ForceNew:         true,
							AtLeastOneOf:     privateClusterConfigKeys,
							DiffSuppressFunc: containerClusterPrivateClusterConfigSuppress,
							Description:      `Subnetwork in cluster's network where master's endpoint will be provisioned.`,
						},
						"public_endpoint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The external IP address of this cluster's master endpoint.`,
						},
						"master_global_access_config": {
							Type:         schema.TypeList,
							MaxItems:     1,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: privateClusterConfigKeys,
							Description:  "Controls cluster master global access settings.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: `Whether the cluster master is accessible globally or not.`,
									},
								},
							},
						},
					},
				},
			},

			"resource_labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `The GCE resource labels (a map of key/value pairs) to be applied to the cluster.

				**Note**: This field is non-authoritative, and will only manage the labels present in your configuration.
				Please refer to the field 'effective_labels' for all of the labels present on the resource.`,
			},
			"terraform_labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: `The combination of labels configured directly on the resource and default labels configured on the provider.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"effective_labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: `All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"label_fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The fingerprint of the set of labels for this cluster.`,
			},

			"default_max_pods_per_node": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				Description:   `The default maximum number of pods per node in this cluster. This doesn't work on "routes-based" clusters, clusters that don't have IP Aliasing enabled.`,
				ConflictsWith: []string{"enable_autopilot"},
			},

			"vertical_pod_autoscaling": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: `Vertical Pod Autoscaling automatically adjusts the resources of pods controlled by it.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: `Enables vertical pod autoscaling.`,
						},
					},
				},
			},
			"workload_identity_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				// Computed is unsafe to remove- this API may return `"workloadIdentityConfig": {},` or omit the key entirely
				// and both will be valid. Note that we don't handle the case where the API returns nothing & the user has defined
				// workload_identity_config today.
				Computed:    true,
				Description: `Configuration for the use of Kubernetes Service Accounts in GCP IAM policies.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"workload_pool": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The workload pool to attach all Kubernetes service accounts to.",
						},
					},
				},
			},

			"identity_service_config": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: `Configuration for Identity Service which allows customers to use external identity providers with the K8S API.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to enable the Identity Service component.",
						},
					},
				},
			},

			"service_external_ips_config": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: `If set, and enabled=true, services with external ips field will not be blocked`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: `When enabled, services with external ips specified will be allowed.`,
						},
					},
				},
			},

			"mesh_certificates": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: `If set, and enable_certificates=true, the GKE Workload Identity Certificates controller and node agent will be deployed in the cluster.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_certificates": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: `When enabled the GKE Workload Identity Certificates controller and node agent will be deployed in the cluster.`,
						},
					},
				},
			},

			"database_encryption": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: `Application-layer Secrets Encryption settings. The object format is {state = string, key_name = string}. Valid values of state are: "ENCRYPTED"; "DECRYPTED". key_name is the name of a CloudKMS key.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"ENCRYPTED", "DECRYPTED"}, false),
							Description:  `ENCRYPTED or DECRYPTED.`,
						},
						"key_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The key to use to encrypt/decrypt secrets.`,
						},
					},
				},
			},

			"release_channel": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: `Configuration options for the Release channel feature, which provide more control over automatic upgrades of your GKE clusters. Note that removing this field from your config will not unenroll it. Instead, use the "UNSPECIFIED" channel.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"channel": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"UNSPECIFIED", "RAPID", "REGULAR", "STABLE", "EXTENDED"}, false),
							Description: `The selected release channel. Accepted values are:
* UNSPECIFIED: Not set.
* RAPID: Weekly upgrade cadence; Early testers and developers who requires new features.
* REGULAR: Multiple per month upgrade cadence; Production users who need features not yet offered in the Stable channel.
* STABLE: Every few months upgrade cadence; Production users who need stability above all else, and for whom frequent upgrades are too risky.
* EXTENDED: GKE provides extended support for Kubernetes minor versions through the Extended channel. With this channel, you can stay on a minor version for up to 24 months.`,
						},
					},
				},
			},

			"gke_auto_upgrade_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: `Configuration options for the auto-upgrade patch type feature, which provide more control over the speed of automatic upgrades of your GKE clusters.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"patch_mode": {
							Type:     schema.TypeString,
							Required: true,
							Description: `The selected auto-upgrade patch type. Accepted values are:
* ACCELERATED: Upgrades to the latest available patch version in a given minor and release channel.`,
						},
					},
				},
			},

			"tpu_ipv4_cidr_block": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: `The IP address range of the Cloud TPUs in this cluster, in CIDR notation (e.g. 1.2.3.4/29).`,
			},

			"cluster_telemetry": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: `Telemetry integration for the cluster.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"DISABLED", "ENABLED", "SYSTEM_ONLY"}, false),
							Description:  `Type of the integration.`,
						},
					},
				},
			},

			"default_snat_status": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: `Whether the cluster disables default in-node sNAT rules. In-node sNAT rules will be disabled when defaultSnatStatus is disabled.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: `When disabled is set to false, default IP masquerade rules will be applied to the nodes to prevent sNAT on cluster internal traffic.`,
						},
					},
				},
			},

			"datapath_provider": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				Description:      `The desired datapath provider for this cluster. By default, uses the IPTables-based kube-proxy implementation.`,
				ValidateFunc:     validation.StringInSlice([]string{"DATAPATH_PROVIDER_UNSPECIFIED", "LEGACY_DATAPATH", "ADVANCED_DATAPATH"}, false),
				DiffSuppressFunc: tpgresource.EmptyOrDefaultStringSuppress("DATAPATH_PROVIDER_UNSPECIFIED"),
			},
			"enable_cilium_clusterwide_network_policy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether Cilium cluster-wide network policy is enabled on this cluster.`,
				Default:     false,
			},
			"enable_intranode_visibility": {
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				Description:   `Whether Intra-node visibility is enabled for this cluster. This makes same node pod to pod traffic visible for VPC network.`,
				ConflictsWith: []string{"enable_autopilot"},
			},
			"enable_l4_ilb_subsetting": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether L4ILB Subsetting is enabled for this cluster.`,
				Default:     false,
			},
			"disable_l4_lb_firewall_reconciliation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Disable L4 load balancer VPC firewalls to enable firewall policies.`,
				Default:     false,
			},
			"enable_multi_networking": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether multi-networking is enabled for this cluster.`,
				Default:     false,
			},
			"enable_fqdn_network_policy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether FQDN Network Policy is enabled on this cluster.`,
				Default:     false,
			},
			"private_ipv6_google_access": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The desired state of IPv6 connectivity to Google Services. By default, no private IPv6 access to or from Google Services (all access will be via IPv4).`,
				Computed:    true,
			},

			"cost_management_config": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: `Cost management configuration for the cluster.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: `Whether to enable GKE cost allocation. When you enable GKE cost allocation, the cluster name and namespace of your GKE workloads appear in the labels field of the billing export to BigQuery. Defaults to false.`,
						},
					},
				},
			},

			"resource_usage_export_config": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: `Configuration for the ResourceUsageExportConfig feature.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_network_egress_metering": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: `Whether to enable network egress metering for this cluster. If enabled, a daemonset will be created in the cluster to meter network egress traffic.`,
						},
						"enable_resource_consumption_metering": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: `Whether to enable resource consumption metering on this cluster. When enabled, a table will be created in the resource export BigQuery dataset to store resource consumption data. The resulting table can be joined with the resource usage table or with BigQuery billing export. Defaults to true.`,
						},
						"bigquery_destination": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							Description: `Parameters for using BigQuery as the destination of resource usage export.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dataset_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The ID of a BigQuery Dataset.`,
									},
								},
							},
						},
					},
				},
			},
			"dns_config": {
				Type:             schema.TypeList,
				Optional:         true,
				MaxItems:         1,
				DiffSuppressFunc: suppressDiffForAutopilot,
				Description:      `Configuration for Cloud DNS for Kubernetes Engine.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"additive_vpc_scope_dns_domain": {
							Type:        schema.TypeString,
							Description: `Enable additive VPC scope DNS in a GKE cluster.`,
							Optional:    true,
						},
						"cluster_dns": {
							Type:         schema.TypeString,
							Default:      "PROVIDER_UNSPECIFIED",
							ValidateFunc: validation.StringInSlice([]string{"PROVIDER_UNSPECIFIED", "PLATFORM_DEFAULT", "CLOUD_DNS", "KUBE_DNS"}, false),
							Description:  `Which in-cluster DNS provider should be used.`,
							Optional:     true,
						},
						"cluster_dns_scope": {
							Type:             schema.TypeString,
							ValidateFunc:     validation.StringInSlice([]string{"DNS_SCOPE_UNSPECIFIED", "CLUSTER_SCOPE", "VPC_SCOPE"}, false),
							Description:      `The scope of access to cluster DNS records.`,
							Optional:         true,
							DiffSuppressFunc: suppressDiffForClusterDnsScope,
						},
						"cluster_dns_domain": {
							Type:        schema.TypeString,
							Description: `The suffix used for all cluster service records.`,
							Optional:    true,
						},
					},
				},
			},
			"gateway_api_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: `Configuration for GKE Gateway API controller.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"channel": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"CHANNEL_DISABLED", "CHANNEL_EXPERIMENTAL", "CHANNEL_STANDARD"}, false),
							Description:  `The Gateway API release channel to use for Gateway API.`,
						},
					},
				},
			},
			"fleet": {
				Type:             schema.TypeList,
				Optional:         true,
				MaxItems:         1,
				Description:      `Fleet configuration of the cluster.`,
				DiffSuppressFunc: suppressDiffForPreRegisteredFleet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The Fleet host project of the cluster.`,
						},
						"membership": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Full resource name of the registered fleet membership of the cluster.`,
						},
						"pre_registered": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the cluster has been registered via the fleet API.`,
						},
						"membership_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Short name of the fleet membership, for example "member-1".`,
						},
						"membership_location": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Location of the fleet membership, for example "us-central1".`,
						},
						"membership_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"LIGHTWEIGHT"}, false),
							Description:  `The type of the cluster's fleet membership.`,
						},
					},
				},
			},
			"user_managed_keys_config": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: `The custom keys configuration of the cluster.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_ca": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: `The Certificate Authority Service caPool to use for the cluster CA in this cluster.`,
						},
						"etcd_api_ca": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: `The Certificate Authority Service caPool to use for the etcd API CA in this cluster.`,
						},
						"etcd_peer_ca": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: `The Certificate Authority Service caPool to use for the etcd peer CA in this cluster.`,
						},
						"aggregation_ca": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: `The Certificate Authority Service caPool to use for the aggreation CA in this cluster.`,
						},
						"service_account_signing_keys": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: `The Cloud KMS cryptoKeyVersions to use for signing service account JWTs issued by this cluster.`,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"service_account_verification_keys": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: `The Cloud KMS cryptoKeyVersions to use for verifying service account JWTs issued by this cluster.`,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"control_plane_disk_encryption_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The Cloud KMS cryptoKey to use for Confidential Hyperdisk on the control plane nodes.`,
						},
						"gkeops_etcd_backup_encryption_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Resource path of the Cloud KMS cryptoKey to use for encryption of internal etcd backups.`,
						},
					},
				},
			},
			"workload_alts_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: `Configuration for direct-path (via ALTS) with workload identity.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_alts": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: `Whether the alts handshaker should be enabled or not for direct-path. Requires Workload Identity (workloadPool must be non-empty).`,
						},
					},
				},
			},
			"enterprise_config": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Computed:    true,
				Description: `Defines the config needed to enable/disable GKE Enterprise`,
				Deprecated:  `GKE Enterprise features are now available without an Enterprise tier. This field is deprecated and will be removed in a future major release`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_tier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the effective cluster tier. Available options include STANDARD and ENTERPRISE.`,
							Deprecated:  `GKE Enterprise features are now available without an Enterprise tier. This field is deprecated and will be removed in a future major release`,
						},
						"desired_tier": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ValidateFunc:     validation.StringInSlice([]string{"STANDARD", "ENTERPRISE"}, false),
							Description:      `Indicates the desired cluster tier. Available options include STANDARD and ENTERPRISE.`,
							Deprecated:       `GKE Enterprise features are now available without an Enterprise tier. This field is deprecated and will be removed in a future major release`,
							DiffSuppressFunc: tpgresource.EmptyOrDefaultStringSuppress("CLUSTER_TIER_UNSPECIFIED"),
						},
					},
				},
			},
			"in_transit_encryption_config": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  `Defines the config of in-transit encryption`,
				ValidateFunc: validation.StringInSlice([]string{"IN_TRANSIT_ENCRYPTION_CONFIG_UNSPECIFIED", "IN_TRANSIT_ENCRYPTION_DISABLED", "IN_TRANSIT_ENCRYPTION_INTER_NODE_TRANSPARENT"}, false),
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
							Description: `Specifies the total network bandwidth tier for NodePools in the cluster.`,
						},
					},
				},
			},
			"anonymous_authentication_config": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Computed:    true,
				Description: `AnonymousAuthenticationConfig allows users to restrict or enable anonymous access to the cluster.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"ENABLED", "LIMITED"}, false),
							Description: `Setting this to LIMITED will restrict authentication of anonymous users to health check endpoints only.
 Accepted values are:
* ENABLED: Authentication of anonymous users is enabled for all endpoints.
* LIMITED: Anonymous access is only allowed for health check endpoints.`,
						},
					},
				},
			},
			"rbac_binding_config": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Computed:    true,
				Description: `RBACBindingConfig allows user to restrict ClusterRoleBindings an RoleBindings that can be created.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_insecure_binding_system_unauthenticated": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Setting this to true will allow any ClusterRoleBinding and RoleBinding with subjects system:anonymous or system:unauthenticated.`,
						},
						"enable_insecure_binding_system_authenticated": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Setting this to true will allow any ClusterRoleBinding and RoleBinding with subjects system:authenticated.`,
						},
					},
				},
			},
		},
	}
}

// Setting a guest accelerator block to count=0 is the equivalent to omitting the block: it won't get
// sent to the API and it won't be stored in state. This diffFunc will try to compare the old + new state
// by only comparing the blocks with a positive count and ignoring those with count=0
//
// One quirk with this approach is that configs with mixed count=0 and count>0 accelerator blocks will
// show a confusing diff if there are config changes that result in a legitimate diff as the count=0
// blocks will not be in state.
func resourceNodeConfigEmptyGuestAccelerator(_ context.Context, diff *schema.ResourceDiff, meta any) error {
	old, new := diff.GetChange("node_config.0.guest_accelerator")
	oList, ok := old.([]any)
	if !ok {
		return fmt.Errorf("type assertion failed, expected []any, got %T", old)
	}
	nList, ok := new.([]any)
	if !ok {
		return fmt.Errorf("type assertion failed, expected []any, got %T", new)
	}

	if len(nList) == len(oList) || len(nList) == 0 {
		return nil
	}
	var hasAcceleratorWithEmptyCount bool
	// the list of blocks in a desired state with count=0 accelerator blocks in it
	// will be longer than the current state.
	// this index tracks the location of positive count accelerator blocks
	index := 0
	for _, item := range nList {
		nAccel, ok := item.(map[string]any)
		if !ok {
			return fmt.Errorf("type assertion failed, expected []any, got %T", item)
		}
		if nAccel["count"].(int) == 0 {
			hasAcceleratorWithEmptyCount = true
			// Ignore any 'empty' accelerators because they aren't sent to the API
			continue
		}
		if index >= len(oList) {
			// Return early if there are more positive count accelerator blocks in the desired state
			// than the current state since a difference in 'legit' blocks is a valid diff.
			// This will prevent array index overruns
			return nil
		}
		// Delete Optional + Computed field from old and new map.
		oAccel, ok := oList[index].(map[string]any)
		if !ok {
			return fmt.Errorf("type assertion failed, expected []any, got %T", oList[index])
		}
		delete(nAccel, "gpu_driver_installation_config")
		delete(oAccel, "gpu_driver_installation_config")
		if !reflect.DeepEqual(oAccel, nAccel) {
			return nil
		}
		index += 1
	}

	if hasAcceleratorWithEmptyCount && index == len(oList) {
		// If the number of count>0 blocks match, there are count=0 blocks present and we
		// haven't already returned due to a legitimate diff
		err := diff.Clear("node_config.0.guest_accelerator")
		if err != nil {
			return err
		}
	}

	return nil
}

func containerClusterMutexKey(project, location, clusterName string) string {
	return fmt.Sprintf("google-container-cluster/%s/%s/%s", project, location, clusterName)
}

func containerClusterFullName(project, location, cluster string) string {
	return fmt.Sprintf("projects/%s/locations/%s/clusters/%s", project, location, cluster)
}

// Suppress unremovable default scope values from GCP.
// If the default service account would not otherwise have it, the `monitoring.write` scope
// is added to a GKE cluster's scopes regardless of what the user provided.
// monitoring.write is inherited from monitoring (rw) and cloud-platform, so it won't always
// be present.
// Enabling Stackdriver features through logging_service and monitoring_service may enable
// monitoring or logging.write. We've chosen not to suppress in those cases because they're
// removable by disabling those features.
func containerClusterAddedScopesSuppress(k, old, new string, d *schema.ResourceData) bool {
	o, n := d.GetChange("cluster_autoscaling.0.auto_provisioning_defaults.0.oauth_scopes")
	if o == nil || n == nil {
		return false
	}

	addedScopes := []string{
		"https://www.googleapis.com/auth/monitoring.write",
	}

	// combine what the default scopes are with what was passed
	m := tpgresource.GolangSetFromStringSlice(append(addedScopes, tpgresource.ConvertStringArr(n.([]interface{}))...))
	combined := tpgresource.StringSliceFromGolangSet(m)

	// compare if the combined new scopes and default scopes differ from the old scopes
	if len(combined) != len(tpgresource.ConvertStringArr(o.([]interface{}))) {
		return false
	}

	for _, i := range combined {
		if tpgresource.StringInSlice(tpgresource.ConvertStringArr(o.([]interface{})), i) {
			continue
		}

		return false
	}

	return true
}

// We want to suppress diffs for empty/disabled private cluster config.
func containerClusterPrivateClusterConfigSuppress(k, old, new string, d *schema.ResourceData) bool {
	o, n := d.GetChange("private_cluster_config.0.enable_private_endpoint")
	suppressEndpoint := !o.(bool) && !n.(bool)

	o, n = d.GetChange("private_cluster_config.0.enable_private_nodes")
	suppressNodes := !o.(bool) && !n.(bool)

	// Do not suppress diffs when private_endpoint_subnetwork is configured
	_, hasSubnet := d.GetOk("private_cluster_config.0.private_endpoint_subnetwork")

	// Do not suppress diffs when master_global_access_config is configured
	_, hasGlobalAccessConfig := d.GetOk("private_cluster_config.0.master_global_access_config")

	if k == "private_cluster_config.0.enable_private_endpoint" {
		return suppressEndpoint && !hasSubnet
	} else if k == "private_cluster_config.0.enable_private_nodes" {
		return suppressNodes && !hasSubnet
	} else if k == "private_cluster_config.#" {
		return suppressEndpoint && suppressNodes && !hasSubnet && !hasGlobalAccessConfig
	} else if k == "private_cluster_config.0.private_endpoint_subnetwork" {
		// Before regular compare, for the sake of private flexible cluster,
		// suppress diffs in private_endpoint_subnetwork when
		//   master_ipv4_cidr_block is set
		//   && private_endpoint_subnetwork is unset in terraform (new value == "")
		//   && private_endpoint_subnetwork is returned from resource (old value != "")
		_, hasMasterCidr := d.GetOk("private_cluster_config.0.master_ipv4_cidr_block")
		return (hasMasterCidr && new == "" && old != "") || tpgresource.CompareSelfLinkOrResourceName(k, old, new, d)
	}
	return false
}

// Autopilot clusters have preconfigured defaults: https://cloud.google.com/kubernetes-engine/docs/concepts/autopilot-overview#comparison.
// This function modifies the diff so users can see what these will be during plan time.
func containerClusterAutopilotCustomizeDiff(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
	if d.HasChange("enable_autopilot") && d.Get("enable_autopilot").(bool) {
		if err := d.SetNew("enable_intranode_visibility", true); err != nil {
			return err
		}
		if err := d.SetNew("networking_mode", "VPC_NATIVE"); err != nil {
			return err
		}
	}
	if d.Get("enable_autopilot").(bool) && d.HasChange("dns_config.0.additive_vpc_scope_dns_domain") {
		return d.ForceNew("dns_config.0.additive_vpc_scope_dns_domain")
	}
	return nil
}

// node_version only applies to the default node pool, so it should conflict with remove_default_node_pool = true
func containerClusterNodeVersionRemoveDefaultCustomizeDiff(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
	// node_version is computed, so we can only check this on initial creation
	o, _ := d.GetChange("name")
	if o != "" {
		return nil
	}
	if d.Get("node_version").(string) != "" && d.Get("remove_default_node_pool").(bool) {
		return fmt.Errorf("node_version can only be specified if remove_default_node_pool is not true")
	}
	return nil
}

func containerClusterNetworkPolicyEmptyCustomizeDiff(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
	// we want to set computed only in the case that there wasn't a previous network_policy configured
	// because we default a returned empty network policy to a configured false, this will only apply
	// on the first run, if network_policy is not configured - all other runs will store empty configurations
	// as enabled=false and provider=PROVIDER_UNSPECIFIED
	o, n := d.GetChange("network_policy")
	if o == nil && n == nil {
		return d.SetNewComputed("network_policy")
	}
	return nil
}

func podSecurityPolicyCfgSuppress(k, old, new string, r *schema.ResourceData) bool {
	if k == "pod_security_policy_config.#" && old == "1" && new == "0" {
		if v, ok := r.GetOk("pod_security_policy_config"); ok {
			cfgList := v.([]interface{})
			if len(cfgList) > 0 {
				d := cfgList[0].(map[string]interface{})
				// Suppress if old value was {enabled == false}
				return !d["enabled"].(bool)
			}
		}
	}
	return false
}

func SecretManagerCfgSuppress(k, old, new string, r *schema.ResourceData) bool {
	if k == "secret_manager_config.#" && old == "1" && new == "0" {
		if v, ok := r.GetOk("secret_manager_config"); ok {
			cfgList := v.([]interface{})
			if len(cfgList) > 0 {
				d := cfgList[0].(map[string]interface{})
				// Suppress if old value was {enabled == false}
				return !d["enabled"].(bool)
			}
		}
	}
	return false
}

func SecretSyncCfgSuppress(k, old, new string, r *schema.ResourceData) bool {
	if k == "secret_sync_config.#" && old == "1" && new == "0" {
		if v, ok := r.GetOk("secret_sync_config"); ok {
			cfgList := v.([]interface{})
			if len(cfgList) > 0 {
				d := cfgList[0].(map[string]interface{})
				// Suppress if old value was {enabled == false}
				return !d["enabled"].(bool)
			}
		}
	}
	return false
}

func containerClusterNetworkPolicyDiffSuppress(k, old, new string, r *schema.ResourceData) bool {
	// if network_policy configuration is empty, we store it as populated and enabled=false, and
	// provider=PROVIDER_UNSPECIFIED, in the case that it was previously stored with this state,
	// and the configuration removed, we want to suppress the diff
	if k == "network_policy.#" && old == "1" && new == "0" {
		o, _ := r.GetChange("network_policy.0.enabled")
		if !o.(bool) {
			return true
		}
	}

	return false
}

func BinaryAuthorizationDiffSuppress(k, old, new string, r *schema.ResourceData) bool {
	// An empty config is equivalent to a config with enabled set to false.
	if k == "binary_authorization.#" && old == "1" && new == "0" {
		o, _ := r.GetChange("binary_authorization.0.enabled")
		if !o.(bool) && !r.HasChange("binary_authorization.0.evaluation_mode") {
			return true
		}
	}

	return false
}

func validateNodePoolAutoConfig(cluster *container.Cluster) error {
	if cluster == nil || cluster.NodePoolAutoConfig == nil {
		return nil
	}
	if cluster.NodePoolAutoConfig != nil && cluster.NodePoolAutoConfig.NetworkTags != nil && len(cluster.NodePoolAutoConfig.NetworkTags.Tags) > 0 {
		if (cluster.Autopilot == nil || !cluster.Autopilot.Enabled) && (cluster.Autoscaling == nil || !cluster.Autoscaling.EnableNodeAutoprovisioning) {
			return fmt.Errorf("node_pool_auto_config network tags can only be set if enable_autopilot or cluster_autoscaling is enabled")
		}
	}

	return nil
}

func containerClusterSurgeSettingsCustomizeDiff(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
	if v, ok := d.GetOk("cluster_autoscaling.0.auto_provisioning_defaults.0.upgrade_settings.0.strategy"); ok {
		if v != "SURGE" {
			if _, maxSurgeIsPresent := d.GetOk("cluster_autoscaling.0.auto_provisioning_defaults.0.upgrade_settings.0.max_surge"); maxSurgeIsPresent {
				return fmt.Errorf("Surge upgrade settings max_surge/max_unavailable can only be used when strategy is set to SURGE")
			}
		}
		if v != "SURGE" {
			if _, maxSurgeIsPresent := d.GetOk("cluster_autoscaling.0.auto_provisioning_defaults.0.upgrade_settings.0.max_unavailable"); maxSurgeIsPresent {
				return fmt.Errorf("Surge upgrade settings max_surge/max_unavailable can only be used when strategy is set to SURGE")
			}
		}
	}

	return nil
}

func containerClusterEnableK8sBetaApisCustomizeDiff(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
	// separate func to allow unit testing
	return containerClusterEnableK8sBetaApisCustomizeDiffFunc(d)
}

func containerClusterEnableK8sBetaApisCustomizeDiffFunc(d tpgresource.TerraformResourceDiff) error {
	// The Kubernetes Beta APIs cannot be disabled once they have been enabled by users.
	// The reason why we don't allow disabling is that the controller does not have the
	// ability to clean up the Kubernetes objects created by the APIs. If the user
	// removes the already enabled Kubernetes Beta API from the list, we need to force
	// a new cluster.
	if !d.HasChange("enable_k8s_beta_apis.0.enabled_apis") {
		return nil
	}
	old, new := d.GetChange("enable_k8s_beta_apis.0.enabled_apis")
	if old != "" && new != "" {
		oldAPIsSet := old.(*schema.Set)
		newAPIsSet := new.(*schema.Set)
		for _, oldAPI := range oldAPIsSet.List() {
			if !newAPIsSet.Contains(oldAPI) {
				return d.ForceNew("enable_k8s_beta_apis.0.enabled_apis")
			}
		}
	}

	return nil
}

func containerClusterNodeVersionCustomizeDiff(_ context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	// separate func to allow unit testing
	return containerClusterNodeVersionCustomizeDiffFunc(diff)
}

func containerClusterNodeVersionCustomizeDiffFunc(diff tpgresource.TerraformResourceDiff) error {
	oldValueName, _ := diff.GetChange("name")
	if oldValueName != "" {
		return nil
	}

	_, newValueNode := diff.GetChange("node_version")
	_, newValueMaster := diff.GetChange("min_master_version")

	if newValueNode == "" || newValueMaster == "" {
		return nil
	}

	//ignore -gke.X suffix for now. If it becomes a problem later, we can fix it
	masterVersion := strings.Split(newValueMaster.(string), "-")[0]
	nodeVersion := strings.Split(newValueNode.(string), "-")[0]

	if masterVersion != nodeVersion {
		return fmt.Errorf("Resource argument node_version (value: %s) must either be unset or set to the same value as min_master_version (value: %s) on create.", newValueNode, newValueMaster)
	}

	return nil
}
