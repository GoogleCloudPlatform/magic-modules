package container

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/tpgresource"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/verify"
	"google.golang.org/api/container/v1"
)

// Matches gke-default scope from https://cloud.google.com/sdk/gcloud/reference/container/clusters/create
var defaultOauthScopes = []string{
	"https://www.googleapis.com/auth/devstorage.read_only",
	"https://www.googleapis.com/auth/logging.write",
	"https://www.googleapis.com/auth/monitoring",
	"https://www.googleapis.com/auth/service.management.readonly",
	"https://www.googleapis.com/auth/servicecontrol",
	"https://www.googleapis.com/auth/trace.append",
}

func schemaContainerdConfig() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    true,
		Description: "Parameters for containerd configuration.",
		MaxItems:    1,
		Elem: &schema.Resource{Schema: map[string]*schema.Schema{
			"private_registry_access_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Parameters for private container registries configuration.",
				MaxItems:    1,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"enabled": {
						Type:        schema.TypeBool,
						Required:    true,
						Description: "Whether or not private registries are configured.",
					},
					"certificate_authority_domain_config": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "Parameters for configuring CA certificate and domains.",
						Elem: &schema.Resource{Schema: map[string]*schema.Schema{
							"fqdns": {
								Type:        schema.TypeList,
								Required:    true,
								Description: "List of fully-qualified-domain-names. IPv4s and port specification are supported.",
								Elem:        &schema.Schema{Type: schema.TypeString},
							},
							"gcp_secret_manager_certificate_config": {
								Type:        schema.TypeList,
								Required:    true,
								Description: "Parameters for configuring a certificate hosted in GCP SecretManager.",
								MaxItems:    1,
								Elem: &schema.Resource{Schema: map[string]*schema.Schema{
									"secret_uri": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "URI for the secret that hosts a certificate. Must be in the format 'projects/PROJECT_NUM/secrets/SECRET_NAME/versions/VERSION_OR_LATEST'.",
									},
								}},
							},
						}},
					},
				}},
			},
			"writable_cgroups": {
				Type:        schema.TypeList,
				Description: `Parameters for writable cgroups configuration.`,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: `Whether writable cgroups are enabled.`,
						},
					},
				},
			},
			"registry_hosts": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Configures containerd registry host configuration. Each registry_hosts entry represents a hosts.toml file.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Defines the host name of the registry server.",
						},
						"hosts": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Configures a list of host-specific configurations for the server.",
							Elem: &schema.Resource{Schema: map[string]*schema.Schema{
								"host": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "Configures the registry host/mirror.",
								},
								"capabilities": {
									Type:        schema.TypeList,
									Optional:    true,
									Description: "Represent the capabilities of the registry host, specifying what operations a host is capable of performing.",
									Elem:        &schema.Schema{Type: schema.TypeString},
								},
								"override_path": {
									Type:        schema.TypeBool,
									Optional:    true,
									Description: "Indicate the host's API root endpoint is defined in the URL path rather than by the API specification.",
								},
								"dial_timeout": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Specifies the maximum duration allowed for a connection attempt to complete.",
								},
								"header": {
									Type:        schema.TypeList,
									Optional:    true,
									Description: "Configures the registry host headers.",
									Elem: &schema.Resource{Schema: map[string]*schema.Schema{
										"key": {
											Type:        schema.TypeString,
											Required:    true,
											Description: "Configures the header key.",
										},
										"value": {
											Type:        schema.TypeList,
											Required:    true,
											Description: "Configures the header value.",
											Elem:        &schema.Schema{Type: schema.TypeString},
										},
									}},
								},
								"ca": {
									Type:        schema.TypeList,
									Optional:    true,
									Description: "Configures the registry host certificate.",
									Elem: &schema.Resource{Schema: map[string]*schema.Schema{
										"gcp_secret_manager_secret_uri": {
											Type:        schema.TypeString,
											Optional:    true,
											Description: "URI for the Secret Manager secret that hosts the certificate.",
										},
									}},
								},
								"client": {
									Type:        schema.TypeList,
									Optional:    true,
									Description: "Configures the registry host client certificate and key.",
									Elem: &schema.Resource{Schema: map[string]*schema.Schema{
										"cert": {
											Type:        schema.TypeList,
											Required:    true,
											MaxItems:    1,
											Description: "Configures the client certificate.",
											Elem: &schema.Resource{Schema: map[string]*schema.Schema{
												"gcp_secret_manager_secret_uri": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "URI for the Secret Manager secret that hosts the client certificate.",
												},
											}},
										},
										"key": {
											Type:        schema.TypeList,
											Optional:    true,
											MaxItems:    1,
											Description: "Configures the client private key.",
											Elem: &schema.Resource{Schema: map[string]*schema.Schema{
												"gcp_secret_manager_secret_uri": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "URI for the Secret Manager secret that hosts the private key.",
												},
											}},
										},
									}},
								},
							},
							},
						},
					}},
			},
		}},
	}
}

// Note: this is a bool internally, but implementing as an enum internally to
// make it easier to accept API level defaults.
func schemaInsecureKubeletReadonlyPortEnabled() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Computed:     true,
		Description:  "Controls whether the kubelet read-only port is enabled. It is strongly recommended to set this to `FALSE`. Possible values: `TRUE`, `FALSE`.",
		ValidateFunc: validation.StringInSlice([]string{"FALSE", "TRUE"}, false),
	}
}

func schemaLoggingVariant() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Computed:     true,
		Description:  `Type of logging agent that is used as the default value for node pools in the cluster. Valid values include DEFAULT and MAX_THROUGHPUT.`,
		ValidateFunc: validation.StringInSlice([]string{"DEFAULT", "MAX_THROUGHPUT"}, false),
	}
}

func schemaGcfsConfig() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    true,
		MaxItems:    1,
		Description: `GCFS configuration for this node.`,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"enabled": {
					Type:        schema.TypeBool,
					Required:    true,
					Description: `Whether or not GCFS is enabled`,
				},
			},
		},
	}
}

func schemaNodeConfig() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    true,
		ForceNew:    true,
		Description: `The configuration of the nodepool`,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"containerd_config": schemaContainerdConfig(),
				"disk_size_gb": {
					Type:         schema.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntAtLeast(10),
					Description:  `Size of the disk attached to each node, specified in GB. The smallest allowed disk size is 10GB.`,
				},

				"disk_type": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					Description: `Type of the disk attached to each node. Such as pd-standard, pd-balanced or pd-ssd`,
				},

				"boot_disk": schemaBootDiskConfig(),

				"guest_accelerator": {
					Type:        schema.TypeList,
					Optional:    true,
					Computed:    true,
					ForceNew:    true,
					Description: `List of the type and count of accelerator cards attached to the instance.`,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"count": {
								Type:        schema.TypeInt,
								Required:    true,
								ForceNew:    true,
								Description: `The number of the accelerator cards exposed to an instance.`,
							},
							"type": {
								Type:             schema.TypeString,
								Required:         true,
								ForceNew:         true,
								DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
								Description:      `The accelerator type resource name.`,
							},
							"gpu_driver_installation_config": {
								Type:        schema.TypeList,
								MaxItems:    1,
								Optional:    true,
								Computed:    true,
								ForceNew:    true,
								Description: `Configuration for auto installation of GPU driver.`,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"gpu_driver_version": {
											Type:         schema.TypeString,
											Required:     true,
											ForceNew:     true,
											Description:  `Mode for how the GPU driver is installed.`,
											ValidateFunc: validation.StringInSlice([]string{"GPU_DRIVER_VERSION_UNSPECIFIED", "INSTALLATION_DISABLED", "DEFAULT", "LATEST"}, false),
										},
									},
								},
							},
							"gpu_partition_size": {
								Type:        schema.TypeString,
								Optional:    true,
								ForceNew:    true,
								Description: `Size of partitions to create on the GPU. Valid values are described in the NVIDIA mig user guide (https://docs.nvidia.com/datacenter/tesla/mig-user-guide/#partitioning)`,
							},
							"gpu_sharing_config": {
								Type:        schema.TypeList,
								MaxItems:    1,
								Optional:    true,
								ForceNew:    true,
								Description: `Configuration for GPU sharing.`,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"gpu_sharing_strategy": {
											Type:        schema.TypeString,
											Required:    true,
											ForceNew:    true,
											Description: `The type of GPU sharing strategy to enable on the GPU node. Possible values are described in the API package (https://pkg.go.dev/google.golang.org/api/container/v1#GPUSharingConfig)`,
										},
										"max_shared_clients_per_gpu": {
											Type:        schema.TypeInt,
											Required:    true,
											ForceNew:    true,
											Description: `The maximum number of containers that can share a GPU.`,
										},
									},
								},
							},
						},
					},
				},

				"image_type": {
					Type:             schema.TypeString,
					Optional:         true,
					Computed:         true,
					DiffSuppressFunc: tpgresource.CaseDiffSuppress,
					Description:      `The image type to use for this node. Note that for a given image type, the latest version of it will be used.`,
				},

				"labels": {
					Type:     schema.TypeMap,
					Optional: true,
					// Computed=true because GKE Sandbox will automatically add labels to nodes that can/cannot run sandboxed pods.
					Computed:         true,
					Elem:             &schema.Schema{Type: schema.TypeString},
					Description:      `The map of Kubernetes labels (key/value pairs) to be applied to each node. These will added in addition to any default label(s) that Kubernetes may apply to the node.`,
					DiffSuppressFunc: containerNodePoolLabelsSuppress,
				},

				"resource_labels": {
					Type:             schema.TypeMap,
					Optional:         true,
					Elem:             &schema.Schema{Type: schema.TypeString},
					DiffSuppressFunc: containerNodePoolResourceLabelsDiffSuppress,
					Description:      `The GCE resource labels (a map of key/value pairs) to be applied to the node pool.`,
				},

				"local_ssd_count": {
					Type:         schema.TypeInt,
					Optional:     true,
					Computed:     true,
					ForceNew:     true,
					ValidateFunc: validation.IntAtLeast(0),
					Description:  `The number of local SSD disks to be attached to the node.`,
				},

				"logging_variant": schemaLoggingVariant(),
				"ephemeral_storage_config": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: `Parameters for the ephemeral storage filesystem. If unspecified, ephemeral storage is backed by the boot disk.`,
					ForceNew:    true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"local_ssd_count": {
								Type:         schema.TypeInt,
								Required:     true,
								ForceNew:     true,
								ValidateFunc: validation.IntAtLeast(0),
								Description:  `Number of local SSDs to use to back ephemeral storage. Uses NVMe interfaces. Each local SSD must be 375 or 3000 GB in size, and all local SSDs must share the same size.`,
							},
						},
					},
				},
				"ephemeral_storage_local_ssd_config": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: `Parameters for the ephemeral storage filesystem. If unspecified, ephemeral storage is backed by the boot disk.`,
					ForceNew:    true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"local_ssd_count": {
								Type:         schema.TypeInt,
								Required:     true,
								ForceNew:     true,
								ValidateFunc: validation.IntAtLeast(0),
								Description:  `Number of local SSDs to use to back ephemeral storage. Uses NVMe interfaces. Each local SSD must be 375 or 3000 GB in size, and all local SSDs must share the same size.`,
							},
							"data_cache_count": {
								Type:         schema.TypeInt,
								Optional:     true,
								ForceNew:     true,
								ValidateFunc: validation.IntAtLeast(0),
								Description:  `Number of local SSDs to be utilized for GKE Data Cache. Uses NVMe interfaces.`,
							},
						},
					},
				},

				"local_nvme_ssd_block_config": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: `Parameters for raw-block local NVMe SSDs.`,
					ForceNew:    true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"local_ssd_count": {
								Type:         schema.TypeInt,
								Required:     true,
								ForceNew:     true,
								ValidateFunc: validation.IntAtLeast(0),
								Description:  `Number of raw-block local NVMe SSD disks to be attached to the node. Each local SSD is 375 GB in size.`,
							},
						},
					},
				},

				"secondary_boot_disks": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    127,
					Description: `Secondary boot disks for preloading data or container images.`,
					ForceNew:    true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"disk_image": {
								Type:        schema.TypeString,
								Required:    true,
								ForceNew:    true,
								Description: `Disk image to create the secondary boot disk from`,
							},
							"mode": {
								Type:        schema.TypeString,
								Optional:    true,
								ForceNew:    true,
								Description: `Mode for how the secondary boot disk is used.`,
							},
						},
					},
				},

				"gcfs_config": schemaGcfsConfig(),

				"gvnic": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: `Enable or disable gvnic in the node pool.`,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:        schema.TypeBool,
								Required:    true,
								Description: `Whether or not gvnic is enabled`,
							},
						},
					},
				},

				"machine_type": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					Description: `The name of a Google Compute Engine machine type.`,
				},

				"metadata": {
					Type:        schema.TypeMap,
					Optional:    true,
					Computed:    true,
					ForceNew:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Description: `The metadata key/value pairs assigned to instances in the cluster.`,
				},

				"min_cpu_platform": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					ForceNew:    true,
					Description: `Minimum CPU platform to be used by this instance. The instance may be scheduled on the specified or newer CPU platform.`,
				},

				"oauth_scopes": {
					Type:        schema.TypeSet,
					Optional:    true,
					Computed:    true,
					ForceNew:    true,
					Description: `The set of Google API scopes to be made available on all of the node VMs.`,
					Elem: &schema.Schema{
						Type: schema.TypeString,
						StateFunc: func(v interface{}) string {
							return tpgresource.CanonicalizeServiceScope(v.(string))
						},
					},
					DiffSuppressFunc: containerClusterAddedScopesSuppress,
					Set:              tpgresource.StringScopeHashcode,
				},

				"preemptible": {
					Type:        schema.TypeBool,
					Optional:    true,
					ForceNew:    true,
					Default:     false,
					Description: `Whether the nodes are created as preemptible VM instances.`,
				},
				"reservation_affinity": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: `The reservation affinity configuration for the node pool.`,
					ForceNew:    true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"consume_reservation_type": {
								Type:         schema.TypeString,
								Required:     true,
								ForceNew:     true,
								Description:  `Corresponds to the type of reservation consumption.`,
								ValidateFunc: validation.StringInSlice([]string{"UNSPECIFIED", "NO_RESERVATION", "ANY_RESERVATION", "SPECIFIC_RESERVATION"}, false),
							},
							"key": {
								Type:        schema.TypeString,
								Optional:    true,
								ForceNew:    true,
								Description: `The label key of a reservation resource.`,
							},
							"values": {
								Type:        schema.TypeSet,
								Description: "The label values of the reservation resource.",
								ForceNew:    true,
								Optional:    true,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
						},
					},
				},
				"spot": {
					Type:        schema.TypeBool,
					Optional:    true,
					ForceNew:    true,
					Default:     false,
					Description: `Whether the nodes are created as spot VM instances.`,
				},

				"service_account": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					ForceNew:    true,
					Description: `The Google Cloud Platform Service Account to be used by the node VMs.`,
				},

				"tags": {
					Type:        schema.TypeList,
					Optional:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Description: `The list of instance tags applied to all nodes.`,
				},

				"storage_pools": {
					Type:        schema.TypeList,
					Optional:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Description: `The list of Storage Pools where boot disks are provisioned.`,
				},

				"shielded_instance_config": {
					Type:        schema.TypeList,
					Optional:    true,
					Computed:    true,
					ForceNew:    true,
					Description: `Shielded Instance options.`,
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enable_secure_boot": {
								Type:        schema.TypeBool,
								Optional:    true,
								ForceNew:    true,
								Default:     false,
								Description: `Defines whether the instance has Secure Boot enabled.`,
							},
							"enable_integrity_monitoring": {
								Type:        schema.TypeBool,
								Optional:    true,
								ForceNew:    true,
								Default:     true,
								Description: `Defines whether the instance has integrity monitoring enabled.`,
							},
						},
					},
				},

				"taint": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: `List of Kubernetes taints to be applied to each node.`,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"key": {
								Type:        schema.TypeString,
								Required:    true,
								Description: `Key for taint.`,
							},
							"value": {
								Type:        schema.TypeString,
								Required:    true,
								Description: `Value for taint.`,
							},
							"effect": {
								Type:         schema.TypeString,
								Required:     true,
								ValidateFunc: validation.StringInSlice([]string{"NO_SCHEDULE", "PREFER_NO_SCHEDULE", "NO_EXECUTE"}, false),
								Description:  `Effect for taint.`,
							},
						},
					},
				},

				"effective_taints": {
					Type:        schema.TypeList,
					Computed:    true,
					Description: `List of kubernetes taints applied to each node.`,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"key": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: `Key for taint.`,
							},
							"value": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: `Value for taint.`,
							},
							"effect": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: `Effect for taint.`,
							},
						},
					},
				},

				"workload_metadata_config": {
					Computed:    true,
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: `The workload metadata configuration for this node.`,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"mode": {
								Type:         schema.TypeString,
								Required:     true,
								ValidateFunc: validation.StringInSlice([]string{"MODE_UNSPECIFIED", "GCE_METADATA", "GKE_METADATA"}, false),
								Description:  `Mode is the configuration for how to expose metadata to workloads running on the node.`,
							},
						},
					},
				},
				"sandbox_config": {
					Type:        schema.TypeList,
					Optional:    true,
					ForceNew:    true,
					MaxItems:    1,
					Description: `Sandbox configuration for this node.`,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"sandbox_type": {
								Type:         schema.TypeString,
								Required:     true,
								Description:  `Type of the sandbox to use for the node (e.g. 'gvisor')`,
								ValidateFunc: validation.StringInSlice([]string{"gvisor"}, false),
							},
						},
					},
				},
				"boot_disk_kms_key": {
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
					Description: `The Customer Managed Encryption Key used to encrypt the boot disk attached to each node in the node pool.`,
				},
				// Note that AtLeastOneOf can't be set because this schema is reused by
				// two different resources.
				"kubelet_config": {
					Type:        schema.TypeList,
					Optional:    true,
					Computed:    true,
					MaxItems:    1,
					Description: `Node kubelet configs.`,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"cpu_manager_policy": {
								Type:         schema.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringInSlice([]string{"static", "none", ""}, false),
								Description:  `Control the CPU management policy on the node.`,
							},
							"memory_manager": {
								Type:        schema.TypeList,
								Optional:    true,
								MaxItems:    1,
								Description: `Configuration for the Memory Manager on the node. The memory manager optimizes memory and hugepages allocation for pods, especially those in the Guaranteed QoS class, by influencing NUMA affinity.`,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"policy": {
											Type:         schema.TypeString,
											Optional:     true,
											Computed:     true,
											Description:  `The Memory Manager policy to use. This policy guides how memory and hugepages are allocated and managed for pods on the node, influencing NUMA affinity.`,
											ValidateFunc: validation.StringInSlice([]string{"None", "Static", ""}, false),
										},
									},
								},
							},
							"topology_manager": {
								Type:        schema.TypeList,
								Optional:    true,
								MaxItems:    1,
								Description: `Configuration for the Topology Manager on the node. The Topology Manager aligns CPU, memory, and device resources on a node to optimize performance, especially for NUMA-aware workloads, by ensuring resource co-location.`,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"policy": {
											Type:         schema.TypeString,
											Optional:     true,
											Computed:     true,
											Description:  `The Topology Manager policy to use. This policy dictates how resource alignment is handled on the node.`,
											ValidateFunc: validation.StringInSlice([]string{"none", "restricted", "single-numa-node", "best-effort", ""}, false),
										},
										"scope": {
											Type:         schema.TypeString,
											Optional:     true,
											Computed:     true,
											Description:  `The Topology Manager scope, defining the granularity at which policy decisions are applied. Valid values are "container" (resources are aligned per container within a pod) or "pod" (resources are aligned for the entire pod).`,
											ValidateFunc: validation.StringInSlice([]string{"container", "pod", ""}, false),
										},
									},
								},
							},
							"cpu_cfs_quota": {
								Type:        schema.TypeBool,
								Computed:    true,
								Optional:    true,
								Description: `Enable CPU CFS quota enforcement for containers that specify CPU limits.`,
							},
							"cpu_cfs_quota_period": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: `Set the CPU CFS quota period value 'cpu.cfs_period_us'.`,
							},
							"insecure_kubelet_readonly_port_enabled": schemaInsecureKubeletReadonlyPortEnabled(),
							"pod_pids_limit": {
								Type:        schema.TypeInt,
								Optional:    true,
								Description: `Controls the maximum number of processes allowed to run in a pod.`,
							},
							"max_parallel_image_pulls": {
								Type:        schema.TypeInt,
								Optional:    true,
								Computed:    true,
								Description: `Set the maximum number of image pulls in parallel.`,
							},
							"container_log_max_size": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: `Defines the maximum size of the container log file before it is rotated.`,
							},
							"container_log_max_files": {
								Type:        schema.TypeInt,
								Optional:    true,
								Description: `Defines the maximum number of container log files that can be present for a container.`,
							},
							"image_gc_low_threshold_percent": {
								Type:        schema.TypeInt,
								Optional:    true,
								Description: `Defines the percent of disk usage before which image garbage collection is never run. Lowest disk usage to garbage collect to.`,
							},
							"image_gc_high_threshold_percent": {
								Type:        schema.TypeInt,
								Optional:    true,
								Description: `Defines the percent of disk usage after which image garbage collection is always run.`,
							},
							"image_minimum_gc_age": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: `Defines the minimum age for an unused image before it is garbage collected.`,
							},
							"image_maximum_gc_age": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: `Defines the maximum age an image can be unused before it is garbage collected.`,
							},
							"allowed_unsafe_sysctls": {
								Type:        schema.TypeList,
								Optional:    true,
								Description: `Defines a comma-separated allowlist of unsafe sysctls or sysctl patterns which can be set on the Pods.`,
								Elem:        &schema.Schema{Type: schema.TypeString},
							},
							"single_process_oom_kill": {
								Type:        schema.TypeBool,
								Optional:    true,
								Description: `Defines whether to enable single process OOM killer.`,
							},
							"eviction_max_pod_grace_period_seconds": {
								Type:        schema.TypeInt,
								Optional:    true,
								Description: `Defines the maximum allowed grace period (in seconds) to use when terminating pods in response to a soft eviction threshold being met.`,
							},
							"eviction_soft": {
								Type:        schema.TypeList,
								Optional:    true,
								MaxItems:    1,
								Description: `Defines a map of signal names to quantities or percentage that defines soft eviction thresholds.`,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"memory_available": {
											Type:        schema.TypeString,
											Optional:    true,
											Description: `Defines quantity of soft eviction threshold for memory.available.`,
										},
										"nodefs_available": {
											Type:        schema.TypeString,
											Optional:    true,
											Description: `Defines percentage of soft eviction threshold for nodefs.available.`,
										},
										"nodefs_inodes_free": {
											Type:        schema.TypeString,
											Optional:    true,
											Description: `Defines percentage of soft eviction threshold for nodefs.inodesFree.`,
										},
										"imagefs_available": {
											Type:        schema.TypeString,
											Optional:    true,
											Description: `Defines percentage of soft eviction threshold for imagefs.available.`,
										},
										"imagefs_inodes_free": {
											Type:        schema.TypeString,
											Optional:    true,
											Description: `Defines percentage of soft eviction threshold for imagefs.inodesFree.`,
										},
										"pid_available": {
											Type:        schema.TypeString,
											Optional:    true,
											Description: `Defines percentage of soft eviction threshold for pid.available.`,
										},
									},
								},
							},
							"eviction_soft_grace_period": {
								Type:        schema.TypeList,
								Optional:    true,
								MaxItems:    1,
								Description: `Defines a map of signal names to durations that defines grace periods for soft eviction thresholds. Each soft eviction threshold must have a corresponding grace period.`,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"memory_available": {
											Type:        schema.TypeString,
											Optional:    true,
											Description: `Defines grace period for the memory.available soft eviction threshold.`,
										},
										"nodefs_available": {
											Type:        schema.TypeString,
											Optional:    true,
											Description: `Defines grace period for the nodefs.available soft eviction threshold.`,
										},
										"nodefs_inodes_free": {
											Type:        schema.TypeString,
											Optional:    true,
											Description: `Defines grace period for the nodefs.inodesFree soft eviction threshold.`,
										},
										"imagefs_available": {
											Type:        schema.TypeString,
											Optional:    true,
											Description: `Defines grace period for the imagefs.available soft eviction threshold`,
										},
										"imagefs_inodes_free": {
											Type:        schema.TypeString,
											Optional:    true,
											Description: `Defines grace period for the imagefs.inodesFree soft eviction threshold.`,
										},
										"pid_available": {
											Type:        schema.TypeString,
											Optional:    true,
											Description: `Defines grace period for the pid.available soft eviction threshold.`,
										},
									},
								},
							},
							"eviction_minimum_reclaim": {
								Type:        schema.TypeList,
								Optional:    true,
								MaxItems:    1,
								Description: `Defines a map of signal names to percentage that defines minimum reclaims. It describes the minimum amount of a given resource the kubelet will reclaim when performing a pod eviction.`,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"memory_available": {
											Type:        schema.TypeString,
											Optional:    true,
											Description: `Defines percentage of minimum reclaim for memory.available.`,
										},
										"nodefs_available": {
											Type:        schema.TypeString,
											Optional:    true,
											Description: `Defines percentage of minimum reclaim for nodefs.available.`,
										},
										"nodefs_inodes_free": {
											Type:        schema.TypeString,
											Optional:    true,
											Description: `Defines percentage of minimum reclaim for nodefs.inodesFree.`,
										},
										"imagefs_available": {
											Type:        schema.TypeString,
											Optional:    true,
											Description: `Defines percentage of minimum reclaim for imagefs.available.`,
										},
										"imagefs_inodes_free": {
											Type:        schema.TypeString,
											Optional:    true,
											Description: `Defines percentage of minimum reclaim for imagefs.inodesFree.`,
										},
										"pid_available": {
											Type:        schema.TypeString,
											Optional:    true,
											Description: `Defines percentage of minimum reclaim for pid.available.`,
										},
									},
								},
							},
						},
					},
				},
				"linux_node_config": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Computed:    true,
					Description: `Parameters that can be configured on Linux nodes.`,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"sysctls": {
								Type:        schema.TypeMap,
								Optional:    true,
								Elem:        &schema.Schema{Type: schema.TypeString},
								Description: `The Linux kernel parameters to be applied to the nodes and all pods running on the nodes.`,
							},
							"cgroup_mode": {
								Type:             schema.TypeString,
								Optional:         true,
								Computed:         true,
								ValidateFunc:     validation.StringInSlice([]string{"CGROUP_MODE_UNSPECIFIED", "CGROUP_MODE_V1", "CGROUP_MODE_V2"}, false),
								Description:      `cgroupMode specifies the cgroup mode to be used on the node.`,
								DiffSuppressFunc: tpgresource.EmptyOrDefaultStringSuppress("CGROUP_MODE_UNSPECIFIED"),
							},
							"node_kernel_module_loading": {
								Type:        schema.TypeList,
								Optional:    true,
								MaxItems:    1,
								Description: `The settings for kernel module loading.`,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"policy": {
											Type:             schema.TypeString,
											Optional:         true,
											ValidateFunc:     validation.StringInSlice([]string{"POLICY_UNSPECIFIED", "ENFORCE_SIGNED_MODULES", "DO_NOT_ENFORCE_SIGNED_MODULES"}, false),
											Description:      `The policy for kernel module loading.`,
											DiffSuppressFunc: tpgresource.EmptyOrDefaultStringSuppress("POLICY_UNSPECIFIED"),
										},
									},
								},
							},
							"transparent_hugepage_enabled": {
								Type:             schema.TypeString,
								Optional:         true,
								Computed:         true,
								ValidateFunc:     validation.StringInSlice([]string{"TRANSPARENT_HUGEPAGE_ENABLED_ALWAYS", "TRANSPARENT_HUGEPAGE_ENABLED_MADVISE", "TRANSPARENT_HUGEPAGE_ENABLED_NEVER", "TRANSPARENT_HUGEPAGE_ENABLED_UNSPECIFIED"}, false),
								Description:      `The Linux kernel transparent hugepage setting.`,
								DiffSuppressFunc: tpgresource.EmptyOrDefaultStringSuppress("TRANSPARENT_HUGEPAGE_ENABLED_UNSPECIFIED"),
							},
							"transparent_hugepage_defrag": {
								Type:             schema.TypeString,
								Optional:         true,
								ValidateFunc:     validation.StringInSlice([]string{"TRANSPARENT_HUGEPAGE_DEFRAG_ALWAYS", "TRANSPARENT_HUGEPAGE_DEFRAG_DEFER", "TRANSPARENT_HUGEPAGE_DEFRAG_DEFER_WITH_MADVISE", "TRANSPARENT_HUGEPAGE_DEFRAG_MADVISE", "TRANSPARENT_HUGEPAGE_DEFRAG_NEVER", "TRANSPARENT_HUGEPAGE_DEFRAG_UNSPECIFIED"}, false),
								Description:      `The Linux kernel transparent hugepage defrag setting.`,
								DiffSuppressFunc: tpgresource.EmptyOrDefaultStringSuppress("TRANSPARENT_HUGEPAGE_DEFRAG_UNSPECIFIED"),
							},
							"hugepages_config": {
								Type:        schema.TypeList,
								Optional:    true,
								MaxItems:    1,
								Description: `Amounts for 2M and 1G hugepages.`,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"hugepage_size_2m": {
											Type:        schema.TypeInt,
											Optional:    true,
											Description: `Amount of 2M hugepages.`,
										},
										"hugepage_size_1g": {
											Type:        schema.TypeInt,
											Optional:    true,
											Description: `Amount of 1G hugepages.`,
										},
									},
								},
							},
						},
					},
				},
				"windows_node_config": {
					Type:        schema.TypeList,
					Optional:    true,
					Computed:    true,
					MaxItems:    1,
					Description: `Parameters that can be configured on Windows nodes.`,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"osversion": {
								Type:         schema.TypeString,
								Optional:     true,
								Default:      "OS_VERSION_UNSPECIFIED",
								Description:  `The OS Version of the windows nodepool.Values are OS_VERSION_UNSPECIFIED,OS_VERSION_LTSC2019 and OS_VERSION_LTSC2022`,
								ValidateFunc: validation.StringInSlice([]string{"OS_VERSION_UNSPECIFIED", "OS_VERSION_LTSC2019", "OS_VERSION_LTSC2022"}, false),
							},
						},
					},
				},
				"node_group": {
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
					Description: `Setting this field will assign instances of this pool to run on the specified node group. This is useful for running workloads on sole tenant nodes.`,
				},

				"advanced_machine_features": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: `Specifies options for controlling advanced machine features.`,
					ForceNew:    true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"threads_per_core": {
								Type:        schema.TypeInt,
								Required:    true,
								ForceNew:    true,
								Description: `The number of threads per physical core. To disable simultaneous multithreading (SMT) set this to 1. If unset, the maximum number of threads supported per core by the underlying processor is assumed.`,
							},
							"enable_nested_virtualization": {
								Type:        schema.TypeBool,
								Optional:    true,
								ForceNew:    true,
								Description: `Whether the node should have nested virtualization enabled.`,
							},
							"performance_monitoring_unit": {
								Type:         schema.TypeString,
								Optional:     true,
								ValidateFunc: verify.ValidateEnum([]string{"ARCHITECTURAL", "STANDARD", "ENHANCED"}),
								Description:  `Level of Performance Monitoring Unit (PMU) requested. If unset, no access to the PMU is assumed.`,
							},
						},
					},
				},
				"sole_tenant_config": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: `Node affinity options for sole tenant node pools.`,
					ForceNew:    true,
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"node_affinity": {
								Type:        schema.TypeSet,
								Required:    true,
								ForceNew:    true,
								Description: `.`,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"key": {
											Type:        schema.TypeString,
											Required:    true,
											ForceNew:    true,
											Description: `.`,
										},
										"operator": {
											Type:         schema.TypeString,
											Required:     true,
											ForceNew:     true,
											Description:  `.`,
											ValidateFunc: validation.StringInSlice([]string{"IN", "NOT_IN"}, false),
										},
										"values": {
											Type:        schema.TypeList,
											Required:    true,
											ForceNew:    true,
											Description: `.`,
											Elem:        &schema.Schema{Type: schema.TypeString},
										},
									},
								},
							},
							"min_node_cpus": {
								Type:        schema.TypeInt,
								Optional:    true,
								Description: `Specifies the minimum number of vCPUs that each sole tenant node must have to use CPU overcommit. If not specified, the CPU overcommit feature is disabled.`,
							},
						},
					},
				},
				"host_maintenance_policy": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: `The maintenance policy for the hosts on which the GKE VMs run on.`,
					ForceNew:    true,
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"maintenance_interval": {
								Type:         schema.TypeString,
								Required:     true,
								ForceNew:     true,
								Description:  `.`,
								ValidateFunc: validation.StringInSlice([]string{"MAINTENANCE_INTERVAL_UNSPECIFIED", "AS_NEEDED", "PERIODIC"}, false),
							},
						},
					},
				},
				"confidential_nodes": {
					Type:        schema.TypeList,
					Optional:    true,
					Computed:    true,
					MaxItems:    1,
					Description: `Configuration for the confidential nodes feature, which makes nodes run on confidential VMs.`,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:        schema.TypeBool,
								Required:    true,
								Description: `Whether Confidential Nodes feature is enabled for all nodes in this pool.`,
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
				"fast_socket": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: `Enable or disable NCCL Fast Socket in the node pool.`,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:        schema.TypeBool,
								Required:    true,
								Description: `Whether or not NCCL Fast Socket is enabled`,
							},
						},
					},
				},
				"resource_manager_tags": {
					Type:        schema.TypeMap,
					Optional:    true,
					Description: `A map of resource manager tags. Resource manager tag keys and values have the same definition as resource manager tags. Keys must be in the format tagKeys/{tag_key_id}, and values are in the format tagValues/456. The field is ignored (both PUT & PATCH) when empty.`,
				},
				"enable_confidential_storage": {
					Type:        schema.TypeBool,
					Optional:    true,
					ForceNew:    true,
					Description: `If enabled boot disks are configured with confidential mode.`,
				},
				"local_ssd_encryption_mode": {
					Type:         schema.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringInSlice([]string{"STANDARD_ENCRYPTION", "EPHEMERAL_KEY_ENCRYPTION"}, false),
					Description:  `LocalSsdEncryptionMode specified the method used for encrypting the local SSDs attached to the node.`,
				},
				"max_run_duration": {
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
					Description: `The runtime of each node in the node pool in seconds, terminated by 's'. Example: "3600s".`,
				},
				"flex_start": {
					Type:        schema.TypeBool,
					Optional:    true,
					ForceNew:    true,
					Description: `Enables Flex Start provisioning model for the node pool`,
				},
			},
		},
	}
}

func schemaBootDiskConfig() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    true,
		MaxItems:    1,
		Description: `Boot disk configuration for node pools nodes.`,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"disk_type": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					Description: `Type of the disk attached to each node. Such as pd-standard, pd-balanced or pd-ssd`,
				},
				"size_gb": {
					Type:         schema.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntAtLeast(10),
					Description:  `Size of the disk attached to each node, specified in GB. The smallest allowed disk size is 10GB.`,
				},
				"provisioned_iops": {
					Type:        schema.TypeInt,
					Optional:    true,
					Computed:    true,
					Description: `Configured IOPs provisioning. Only valid with disk type hyperdisk-balanced.`,
				},
				"provisioned_throughput": {
					Type:        schema.TypeInt,
					Optional:    true,
					Computed:    true,
					Description: `Configured throughput provisioning. Only valid with disk type hyperdisk-balanced.`,
				},
			},
		},
	}
}

// Separate since this currently only supports a single value -- a subset of
// the overall NodeKubeletConfig
func schemaNodePoolAutoConfigNodeKubeletConfig() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    true,
		MaxItems:    1,
		Description: `Node kubelet configs.`,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"insecure_kubelet_readonly_port_enabled": schemaInsecureKubeletReadonlyPortEnabled(),
			},
		},
	}
}

// Separate since this currently only supports a subset of the overall
// LinuxNodeConfig
func schemaNodePoolAutoConfigLinuxNodeConfig() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: `Linux node configuration options.`,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"cgroup_mode": {
					Type:             schema.TypeString,
					Optional:         true,
					Computed:         true,
					ValidateFunc:     validation.StringInSlice([]string{"CGROUP_MODE_UNSPECIFIED", "CGROUP_MODE_V1", "CGROUP_MODE_V2"}, false),
					Description:      `cgroupMode specifies the cgroup mode to be used on the node.`,
					DiffSuppressFunc: tpgresource.EmptyOrDefaultStringSuppress("CGROUP_MODE_UNSPECIFIED"),
				},
				"node_kernel_module_loading": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: `The settings for kernel module loading.`,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"policy": {
								Type:             schema.TypeString,
								Optional:         true,
								ValidateFunc:     validation.StringInSlice([]string{"POLICY_UNSPECIFIED", "ENFORCE_SIGNED_MODULES", "DO_NOT_ENFORCE_SIGNED_MODULES"}, false),
								Description:      `The policy for kernel module loading.`,
								DiffSuppressFunc: tpgresource.EmptyOrDefaultStringSuppress("POLICY_UNSPECIFIED"),
							},
						},
					},
				},
			},
		},
	}
}

func expandNodeConfigDefaults(configured interface{}) *container.NodeConfigDefaults {
	configs := configured.([]interface{})
	if len(configs) == 0 || configs[0] == nil {
		return nil
	}
	config := configs[0].(map[string]interface{})

	nodeConfigDefaults := &container.NodeConfigDefaults{}
	nodeConfigDefaults.ContainerdConfig = expandContainerdConfig(config["containerd_config"])
	if v, ok := config["insecure_kubelet_readonly_port_enabled"]; ok {
		nodeConfigDefaults.NodeKubeletConfig = &container.NodeKubeletConfig{
			InsecureKubeletReadonlyPortEnabled: expandInsecureKubeletReadonlyPortEnabled(v),
			ForceSendFields:                    []string{"InsecureKubeletReadonlyPortEnabled"},
		}
	}
	if variant, ok := config["logging_variant"]; ok {
		nodeConfigDefaults.LoggingConfig = &container.NodePoolLoggingConfig{
			VariantConfig: &container.LoggingVariantConfig{
				Variant: variant.(string),
			},
		}
	}

	if v, ok := config["gcfs_config"]; ok && len(v.([]interface{})) > 0 {
		gcfsConfig := v.([]interface{})[0].(map[string]interface{})
		nodeConfigDefaults.GcfsConfig = &container.GcfsConfig{
			Enabled: gcfsConfig["enabled"].(bool),
		}
	}

	return nodeConfigDefaults
}

func expandNodeConfig(d tpgresource.TerraformResourceData, prefix string, v interface{}) *container.NodeConfig {
	nodeConfigs := v.([]interface{})
	nc := &container.NodeConfig{
		// Defaults can't be set on a list/set in the schema, so set the default on create here.
		OauthScopes: defaultOauthScopes,
	}
	if len(nodeConfigs) == 0 {
		return nc
	}

	nodeConfig := nodeConfigs[0].(map[string]interface{})

	if v, ok := nodeConfig["containerd_config"]; ok {
		nc.ContainerdConfig = expandContainerdConfig(v)
	}

	if v, ok := nodeConfig["machine_type"]; ok {
		nc.MachineType = v.(string)
	}

	if v, ok := nodeConfig["guest_accelerator"]; ok {
		accels := v.([]interface{})
		guestAccelerators := make([]*container.AcceleratorConfig, 0, len(accels))
		for _, raw := range accels {
			data := raw.(map[string]interface{})
			if data["count"].(int) == 0 {
				continue
			}
			guestAcceleratorConfig := &container.AcceleratorConfig{
				AcceleratorCount: int64(data["count"].(int)),
				AcceleratorType:  data["type"].(string),
				GpuPartitionSize: data["gpu_partition_size"].(string),
			}

			if v, ok := data["gpu_driver_installation_config"]; ok && len(v.([]interface{})) > 0 {
				gpuDriverInstallationConfig := data["gpu_driver_installation_config"].([]interface{})[0].(map[string]interface{})
				guestAcceleratorConfig.GpuDriverInstallationConfig = &container.GPUDriverInstallationConfig{
					GpuDriverVersion: gpuDriverInstallationConfig["gpu_driver_version"].(string),
				}
			}

			if v, ok := data["gpu_sharing_config"]; ok && len(v.([]interface{})) > 0 {
				gpuSharingConfig := data["gpu_sharing_config"].([]interface{})[0].(map[string]interface{})
				guestAcceleratorConfig.GpuSharingConfig = &container.GPUSharingConfig{
					GpuSharingStrategy:     gpuSharingConfig["gpu_sharing_strategy"].(string),
					MaxSharedClientsPerGpu: int64(gpuSharingConfig["max_shared_clients_per_gpu"].(int)),
				}
			}

			guestAccelerators = append(guestAccelerators, guestAcceleratorConfig)
		}
		nc.Accelerators = guestAccelerators
	}

	if v, ok := nodeConfig["disk_size_gb"]; ok {
		nc.DiskSizeGb = int64(v.(int))
	}

	if v, ok := nodeConfig["disk_type"]; ok {
		nc.DiskType = v.(string)
	}

	if v, ok := nodeConfig["boot_disk"]; ok {
		nc.BootDisk = expandBootDiskConfig(v)
	}

	if v, ok := nodeConfig["local_ssd_count"]; ok {
		nc.LocalSsdCount = int64(v.(int))
	}

	if v, ok := nodeConfig["logging_variant"]; ok {
		nc.LoggingConfig = &container.NodePoolLoggingConfig{
			VariantConfig: &container.LoggingVariantConfig{
				Variant: v.(string),
			},
		}
	}

	if v, ok := nodeConfig["local_nvme_ssd_block_config"]; ok && len(v.([]interface{})) > 0 {
		conf := v.([]interface{})[0].(map[string]interface{})
		nc.LocalNvmeSsdBlockConfig = &container.LocalNvmeSsdBlockConfig{
			LocalSsdCount: int64(conf["local_ssd_count"].(int)),
		}
	}

	if v, ok := nodeConfig["ephemeral_storage_local_ssd_config"]; ok && len(v.([]interface{})) > 0 {
		conf := v.([]interface{})[0].(map[string]interface{})
		nc.EphemeralStorageLocalSsdConfig = &container.EphemeralStorageLocalSsdConfig{
			LocalSsdCount: int64(conf["local_ssd_count"].(int)),
		}
		dataCacheCount, ok := conf["data_cache_count"]
		if ok {
			nc.EphemeralStorageLocalSsdConfig.DataCacheCount = int64(dataCacheCount.(int))
		}

	}

	if v, ok := nodeConfig["secondary_boot_disks"]; ok && len(v.([]interface{})) > 0 {
		conf, confOK := v.([]interface{})[0].(map[string]interface{})
		if confOK {
			modeValue, modeOK := conf["mode"]
			diskImage := conf["disk_image"].(string)
			if modeOK {
				nc.SecondaryBootDisks = append(nc.SecondaryBootDisks, &container.SecondaryBootDisk{
					DiskImage: diskImage,
					Mode:      modeValue.(string),
				})
			} else {
				nc.SecondaryBootDisks = append(nc.SecondaryBootDisks, &container.SecondaryBootDisk{
					DiskImage: diskImage,
				})
			}
		} else {
			nc.SecondaryBootDisks = append(nc.SecondaryBootDisks, &container.SecondaryBootDisk{
				DiskImage: "",
			})
		}
	}

	if v, ok := nodeConfig["gcfs_config"]; ok && len(v.([]interface{})) > 0 {
		conf := v.([]interface{})[0].(map[string]interface{})
		nc.GcfsConfig = &container.GcfsConfig{
			Enabled: conf["enabled"].(bool),
		}
	}

	if v, ok := nodeConfig["gvnic"]; ok && len(v.([]interface{})) > 0 {
		conf := v.([]interface{})[0].(map[string]interface{})
		nc.Gvnic = &container.VirtualNIC{
			Enabled: conf["enabled"].(bool),
		}
	}

	if v, ok := nodeConfig["fast_socket"]; ok && len(v.([]interface{})) > 0 {
		conf := v.([]interface{})[0].(map[string]interface{})
		nc.FastSocket = &container.FastSocket{
			Enabled: conf["enabled"].(bool),
		}
	}

	if v, ok := nodeConfig["reservation_affinity"]; ok && len(v.([]interface{})) > 0 {
		conf := v.([]interface{})[0].(map[string]interface{})
		valuesSet := conf["values"].(*schema.Set)
		values := make([]string, valuesSet.Len())
		for i, value := range valuesSet.List() {
			values[i] = value.(string)
		}

		nc.ReservationAffinity = &container.ReservationAffinity{
			ConsumeReservationType: conf["consume_reservation_type"].(string),
			Key:                    conf["key"].(string),
			Values:                 values,
		}
	}

	if scopes, ok := nodeConfig["oauth_scopes"]; ok {
		scopesSet := scopes.(*schema.Set)
		scopes := make([]string, scopesSet.Len())
		for i, scope := range scopesSet.List() {
			scopes[i] = tpgresource.CanonicalizeServiceScope(scope.(string))
		}

		nc.OauthScopes = scopes
	}

	if v, ok := nodeConfig["service_account"]; ok {
		nc.ServiceAccount = v.(string)
	}

	if v, ok := nodeConfig["metadata"]; ok {
		m := make(map[string]string)
		for k, val := range v.(map[string]interface{}) {
			m[k] = val.(string)
		}
		nc.Metadata = m
	}

	if v, ok := nodeConfig["image_type"]; ok {
		nc.ImageType = v.(string)
	}

	if v, ok := nodeConfig["labels"]; ok {
		m := make(map[string]string)
		for k, val := range v.(map[string]interface{}) {
			m[k] = val.(string)
		}
		nc.Labels = m
	}

	if v, ok := nodeConfig["resource_labels"]; ok {
		m := make(map[string]string)
		for k, val := range v.(map[string]interface{}) {
			m[k] = val.(string)
		}
		nc.ResourceLabels = m
	}

	if v, ok := nodeConfig["resource_manager_tags"]; ok && len(v.(map[string]interface{})) > 0 {
		nc.ResourceManagerTags = expandResourceManagerTags(v)
	}

	if v, ok := nodeConfig["tags"]; ok {
		tagsList := v.([]interface{})
		tags := []string{}
		for _, v := range tagsList {
			if v != nil {
				tags = append(tags, v.(string))
			}
		}
		nc.Tags = tags
	}

	if v, ok := nodeConfig["storage_pools"]; ok {
		spList := v.([]interface{})
		storagePools := []string{}
		for _, v := range spList {
			if v != nil {
				storagePools = append(storagePools, v.(string))
			}
		}
		nc.StoragePools = storagePools
	}
	if v, ok := nodeConfig["shielded_instance_config"]; ok && len(v.([]interface{})) > 0 {
		conf := v.([]interface{})[0].(map[string]interface{})
		nc.ShieldedInstanceConfig = &container.ShieldedInstanceConfig{
			EnableSecureBoot:          conf["enable_secure_boot"].(bool),
			EnableIntegrityMonitoring: conf["enable_integrity_monitoring"].(bool),
		}
	}

	// Preemptible Is Optional+Default, so it always has a value
	nc.Preemptible = nodeConfig["preemptible"].(bool)

	// Spot Is Optional+Default, so it always has a value
	nc.Spot = nodeConfig["spot"].(bool)

	if v, ok := nodeConfig["min_cpu_platform"]; ok {
		nc.MinCpuPlatform = v.(string)
	}

	if v, ok := nodeConfig["taint"]; ok && len(v.([]interface{})) > 0 {
		taints := v.([]interface{})
		nodeTaints := make([]*container.NodeTaint, 0, len(taints))
		for _, raw := range taints {
			data := raw.(map[string]interface{})
			taint := &container.NodeTaint{
				Key:    data["key"].(string),
				Value:  data["value"].(string),
				Effect: data["effect"].(string),
			}

			nodeTaints = append(nodeTaints, taint)
		}

		nc.Taints = nodeTaints
	}

	if v, ok := nodeConfig["workload_metadata_config"]; ok {
		nc.WorkloadMetadataConfig = expandWorkloadMetadataConfig(v)
	}

	if v, ok := nodeConfig["boot_disk_kms_key"]; ok {
		nc.BootDiskKmsKey = v.(string)
	}

	if v, ok := nodeConfig["kubelet_config"]; ok {
		nc.KubeletConfig = expandKubeletConfig(v)
	}

	if v, ok := nodeConfig["linux_node_config"]; ok {
		nc.LinuxNodeConfig = expandLinuxNodeConfig(v)
	}

	if v, ok := nodeConfig["windows_node_config"]; ok {
		nc.WindowsNodeConfig = expandWindowsNodeConfig(v)
	}

	if v, ok := nodeConfig["node_group"]; ok {
		nc.NodeGroup = v.(string)
	}

	if v, ok := nodeConfig["advanced_machine_features"]; ok && len(v.([]interface{})) > 0 {
		advanced_machine_features := v.([]interface{})[0].(map[string]interface{})
		nc.AdvancedMachineFeatures = &container.AdvancedMachineFeatures{
			ThreadsPerCore:             int64(advanced_machine_features["threads_per_core"].(int)),
			EnableNestedVirtualization: advanced_machine_features["enable_nested_virtualization"].(bool),
			PerformanceMonitoringUnit:  advanced_machine_features["performance_monitoring_unit"].(string),
		}
	}

	if v, ok := nodeConfig["sole_tenant_config"]; ok && len(v.([]interface{})) > 0 {
		nc.SoleTenantConfig = expandSoleTenantConfig(v)
	}

	if v, ok := nodeConfig["enable_confidential_storage"]; ok {
		nc.EnableConfidentialStorage = v.(bool)
	}

	if v, ok := nodeConfig["local_ssd_encryption_mode"]; ok {
		nc.LocalSsdEncryptionMode = v.(string)
	}

	if v, ok := nodeConfig["max_run_duration"]; ok {
		nc.MaxRunDuration = v.(string)
	}

	if v, ok := nodeConfig["flex_start"]; ok {
		nc.FlexStart = v.(bool)
	}

	if v, ok := nodeConfig["confidential_nodes"]; ok {
		nc.ConfidentialNodes = expandConfidentialNodes(v)
	}

	return nc
}

func expandBootDiskConfig(v interface{}) *container.BootDisk {
	bd := &container.BootDisk{}
	if v == nil {
		return nil
	}
	ls := v.([]interface{})
	if len(ls) == 0 {
		return nil
	}
	cfg := ls[0].(map[string]interface{})

	if v, ok := cfg["disk_type"]; ok {
		bd.DiskType = v.(string)
	}

	if v, ok := cfg["size_gb"]; ok {
		bd.SizeGb = int64(v.(int))
	}

	if v, ok := cfg["provisioned_iops"]; ok {
		bd.ProvisionedIops = int64(v.(int))
	}

	if v, ok := cfg["provisioned_throughput"]; ok {
		bd.ProvisionedThroughput = int64(v.(int))
	}

	return bd
}

func expandResourceManagerTags(v interface{}) *container.ResourceManagerTags {
	if v == nil {
		return nil
	}

	rmts := make(map[string]string)

	if v != nil {
		rmts = tpgresource.ConvertStringMap(v.(map[string]interface{}))
	}

	return &container.ResourceManagerTags{
		Tags:            rmts,
		ForceSendFields: []string{"Tags"},
	}
}

func expandWorkloadMetadataConfig(v interface{}) *container.WorkloadMetadataConfig {
	if v == nil {
		return nil
	}
	ls := v.([]interface{})
	if len(ls) == 0 {
		return nil
	}
	wmc := &container.WorkloadMetadataConfig{}

	cfg := ls[0].(map[string]interface{})

	if v, ok := cfg["mode"]; ok {
		wmc.Mode = v.(string)
	}

	return wmc
}

func expandInsecureKubeletReadonlyPortEnabled(v interface{}) bool {
	if v == "TRUE" {
		return true
	}
	return false
}

func expandKubeletConfig(v interface{}) *container.NodeKubeletConfig {
	if v == nil {
		return nil
	}
	ls := v.([]interface{})
	if len(ls) == 0 {
		return nil
	}
	cfg := ls[0].(map[string]interface{})
	kConfig := &container.NodeKubeletConfig{}
	if cpuManagerPolicy, ok := cfg["cpu_manager_policy"]; ok {
		kConfig.CpuManagerPolicy = cpuManagerPolicy.(string)
	}
	if cpuCfsQuota, ok := cfg["cpu_cfs_quota"]; ok {
		kConfig.CpuCfsQuota = cpuCfsQuota.(bool)
	}
	if cpuCfsQuotaPeriod, ok := cfg["cpu_cfs_quota_period"]; ok {
		kConfig.CpuCfsQuotaPeriod = cpuCfsQuotaPeriod.(string)
	}
	if insecureKubeletReadonlyPortEnabled, ok := cfg["insecure_kubelet_readonly_port_enabled"]; ok {
		kConfig.InsecureKubeletReadonlyPortEnabled = expandInsecureKubeletReadonlyPortEnabled(insecureKubeletReadonlyPortEnabled)
		kConfig.ForceSendFields = append(kConfig.ForceSendFields, "InsecureKubeletReadonlyPortEnabled")
	}
	if podPidsLimit, ok := cfg["pod_pids_limit"]; ok {
		kConfig.PodPidsLimit = int64(podPidsLimit.(int))
	}
	if maxParallelImagePulls, ok := cfg["max_parallel_image_pulls"]; ok {
		kConfig.MaxParallelImagePulls = int64(maxParallelImagePulls.(int))
	}
	if containerLogMaxSize, ok := cfg["container_log_max_size"]; ok {
		kConfig.ContainerLogMaxSize = containerLogMaxSize.(string)
	}
	if containerLogMaxFiles, ok := cfg["container_log_max_files"]; ok {
		kConfig.ContainerLogMaxFiles = int64(containerLogMaxFiles.(int))
	}
	if imageGcLowThresholdPercent, ok := cfg["image_gc_low_threshold_percent"]; ok {
		kConfig.ImageGcLowThresholdPercent = int64(imageGcLowThresholdPercent.(int))
	}
	if imageGcHighThresholdPercent, ok := cfg["image_gc_high_threshold_percent"]; ok {
		kConfig.ImageGcHighThresholdPercent = int64(imageGcHighThresholdPercent.(int))
	}
	if imageMinimumGcAge, ok := cfg["image_minimum_gc_age"]; ok {
		kConfig.ImageMinimumGcAge = imageMinimumGcAge.(string)
	}
	if imageMaximumGcAge, ok := cfg["image_maximum_gc_age"]; ok {
		kConfig.ImageMaximumGcAge = imageMaximumGcAge.(string)
	}
	if allowedUnsafeSysctls, ok := cfg["allowed_unsafe_sysctls"]; ok {
		sysctls := allowedUnsafeSysctls.([]interface{})
		kConfig.AllowedUnsafeSysctls = make([]string, len(sysctls))
		for i, s := range sysctls {
			kConfig.AllowedUnsafeSysctls[i] = s.(string)
		}
	}
	if singleProcessOomKill, ok := cfg["single_process_oom_kill"]; ok {
		kConfig.SingleProcessOomKill = singleProcessOomKill.(bool)
	}
	if evictionMaxPodGracePeriodSeconds, ok := cfg["eviction_max_pod_grace_period_seconds"]; ok {
		kConfig.EvictionMaxPodGracePeriodSeconds = int64(evictionMaxPodGracePeriodSeconds.(int))
	}
	if v, ok := cfg["eviction_soft"]; ok && len(v.([]interface{})) > 0 {
		es := v.([]interface{})[0].(map[string]interface{})
		evictionSoft := &container.EvictionSignals{}
		if val, ok := es["memory_available"]; ok {
			evictionSoft.MemoryAvailable = val.(string)
		}
		if val, ok := es["nodefs_available"]; ok {
			evictionSoft.NodefsAvailable = val.(string)
		}
		if val, ok := es["imagefs_available"]; ok {
			evictionSoft.ImagefsAvailable = val.(string)
		}
		if val, ok := es["imagefs_inodes_free"]; ok {
			evictionSoft.ImagefsInodesFree = val.(string)
		}
		if val, ok := es["nodefs_inodes_free"]; ok {
			evictionSoft.NodefsInodesFree = val.(string)
		}
		if val, ok := es["pid_available"]; ok {
			evictionSoft.PidAvailable = val.(string)
		}
		kConfig.EvictionSoft = evictionSoft
	}

	if v, ok := cfg["memory_manager"]; ok {
		kConfig.MemoryManager = expandMemoryManager(v)
	}
	if v, ok := cfg["topology_manager"]; ok {
		kConfig.TopologyManager = expandTopologyManager(v)
	}

	if v, ok := cfg["eviction_soft_grace_period"]; ok && len(v.([]interface{})) > 0 {
		es := v.([]interface{})[0].(map[string]interface{})
		periods := &container.EvictionGracePeriod{}
		if val, ok := es["memory_available"]; ok {
			periods.MemoryAvailable = val.(string)
		}
		if val, ok := es["nodefs_available"]; ok {
			periods.NodefsAvailable = val.(string)
		}
		if val, ok := es["imagefs_available"]; ok {
			periods.ImagefsAvailable = val.(string)
		}
		if val, ok := es["imagefs_inodes_free"]; ok {
			periods.ImagefsInodesFree = val.(string)
		}
		if val, ok := es["nodefs_inodes_free"]; ok {
			periods.NodefsInodesFree = val.(string)
		}
		if val, ok := es["pid_available"]; ok {
			periods.PidAvailable = val.(string)
		}
		kConfig.EvictionSoftGracePeriod = periods
	}
	if v, ok := cfg["eviction_minimum_reclaim"]; ok && len(v.([]interface{})) > 0 {
		es := v.([]interface{})[0].(map[string]interface{})
		reclaim := &container.EvictionMinimumReclaim{}
		if val, ok := es["memory_available"]; ok {
			reclaim.MemoryAvailable = val.(string)
		}
		if val, ok := es["nodefs_available"]; ok {
			reclaim.NodefsAvailable = val.(string)
		}
		if val, ok := es["imagefs_available"]; ok {
			reclaim.ImagefsAvailable = val.(string)
		}
		if val, ok := es["imagefs_inodes_free"]; ok {
			reclaim.ImagefsInodesFree = val.(string)
		}
		if val, ok := es["nodefs_inodes_free"]; ok {
			reclaim.NodefsInodesFree = val.(string)
		}
		if val, ok := es["pid_available"]; ok {
			reclaim.PidAvailable = val.(string)
		}
		kConfig.EvictionMinimumReclaim = reclaim
	}
	return kConfig
}

func expandTopologyManager(v interface{}) *container.TopologyManager {
	if v == nil {
		return nil
	}
	ls := v.([]interface{})
	if len(ls) == 0 {
		return nil
	}
	if ls[0] == nil {
		return &container.TopologyManager{}
	}
	cfg := ls[0].(map[string]interface{})

	topologyManager := &container.TopologyManager{}

	if v, ok := cfg["policy"]; ok {
		topologyManager.Policy = v.(string)
	}

	if v, ok := cfg["scope"]; ok {
		topologyManager.Scope = v.(string)
	}

	return topologyManager
}

func expandMemoryManager(v interface{}) *container.MemoryManager {
	if v == nil {
		return nil
	}
	ls := v.([]interface{})
	if len(ls) == 0 {
		return nil
	}
	if ls[0] == nil {
		return &container.MemoryManager{}
	}
	cfg := ls[0].(map[string]interface{})

	memoryManager := &container.MemoryManager{}

	if v, ok := cfg["policy"]; ok {
		memoryManager.Policy = v.(string)
	}

	return memoryManager
}

func expandLinuxNodeConfig(v interface{}) *container.LinuxNodeConfig {
	if v == nil {
		return nil
	}
	ls := v.([]interface{})
	if len(ls) == 0 {
		return nil
	}
	if ls[0] == nil {
		return &container.LinuxNodeConfig{}
	}
	cfg := ls[0].(map[string]interface{})

	linuxNodeConfig := &container.LinuxNodeConfig{}
	sysctls := expandSysctls(cfg)
	if sysctls != nil {
		linuxNodeConfig.Sysctls = sysctls
	}
	cgroupMode := expandCgroupMode(cfg)
	if len(cgroupMode) != 0 {
		linuxNodeConfig.CgroupMode = cgroupMode
	}

	if v, ok := cfg["transparent_hugepage_enabled"]; ok {
		linuxNodeConfig.TransparentHugepageEnabled = v.(string)
	}
	if v, ok := cfg["transparent_hugepage_defrag"]; ok {
		linuxNodeConfig.TransparentHugepageDefrag = v.(string)
	}

	if v, ok := cfg["hugepages_config"]; ok {
		linuxNodeConfig.Hugepages = expandHugepagesConfig(v)
	}

	if v, ok := cfg["node_kernel_module_loading"]; ok {
		linuxNodeConfig.NodeKernelModuleLoading = expandNodeKernelModuleLoading(v)
	}

	return linuxNodeConfig
}

func expandWindowsNodeConfig(v interface{}) *container.WindowsNodeConfig {
	if v == nil {
		return nil
	}
	ls := v.([]interface{})
	if len(ls) == 0 {
		return nil
	}
	cfg := ls[0].(map[string]interface{})
	osversionRaw, ok := cfg["osversion"]
	if !ok {
		return nil
	}
	return &container.WindowsNodeConfig{
		OsVersion: osversionRaw.(string),
	}
}

func expandSysctls(cfg map[string]interface{}) map[string]string {
	sysCfgRaw, ok := cfg["sysctls"]
	if !ok {
		return nil
	}
	sysctls := make(map[string]string)
	for k, v := range sysCfgRaw.(map[string]interface{}) {
		sysctls[k] = v.(string)
	}
	return sysctls
}

func expandCgroupMode(cfg map[string]interface{}) string {
	cgroupMode, ok := cfg["cgroup_mode"]
	if !ok {
		return ""
	}

	return cgroupMode.(string)
}

func expandHugepagesConfig(v interface{}) *container.HugepagesConfig {
	if v == nil {
		return nil
	}
	ls := v.([]interface{})
	if len(ls) == 0 {
		return nil
	}
	if ls[0] == nil {
		return &container.HugepagesConfig{}
	}
	cfg := ls[0].(map[string]interface{})

	hugepagesConfig := &container.HugepagesConfig{}

	if v, ok := cfg["hugepage_size_2m"]; ok {
		hugepagesConfig.HugepageSize2m = int64(v.(int))
	}

	if v, ok := cfg["hugepage_size_1g"]; ok {
		hugepagesConfig.HugepageSize1g = int64(v.(int))
	}

	return hugepagesConfig
}

func expandNodeKernelModuleLoading(v interface{}) *container.NodeKernelModuleLoading {
	if v == nil {
		return nil
	}
	ls := v.([]interface{})
	if len(ls) == 0 {
		return nil
	}
	if ls[0] == nil {
		return &container.NodeKernelModuleLoading{}
	}
	cfg := ls[0].(map[string]interface{})

	NodeKernelModuleLoading := &container.NodeKernelModuleLoading{}

	if v, ok := cfg["policy"]; ok {
		NodeKernelModuleLoading.Policy = v.(string)
	}

	return NodeKernelModuleLoading
}

func expandContainerdConfig(v interface{}) *container.ContainerdConfig {
	if v == nil {
		return nil
	}
	ls := v.([]interface{})
	if len(ls) == 0 {
		return nil
	}
	if ls[0] == nil {
		return &container.ContainerdConfig{}
	}

	cfg := ls[0].(map[string]interface{})

	cc := &container.ContainerdConfig{}
	cc.PrivateRegistryAccessConfig = expandPrivateRegistryAccessConfig(cfg["private_registry_access_config"])
	cc.WritableCgroups = expandWritableCgroups(cfg["writable_cgroups"])
	cc.RegistryHosts = expandRegistryHosts(cfg["registry_hosts"])
	return cc
}

func expandPrivateRegistryAccessConfig(v interface{}) *container.PrivateRegistryAccessConfig {
	if v == nil {
		return nil
	}
	ls := v.([]interface{})
	if len(ls) == 0 {
		return nil
	}
	if ls[0] == nil {
		return &container.PrivateRegistryAccessConfig{}
	}
	cfg := ls[0].(map[string]interface{})

	pracc := &container.PrivateRegistryAccessConfig{}
	if enabled, ok := cfg["enabled"]; ok {
		pracc.Enabled = enabled.(bool)
	}
	if caCfgRaw, ok := cfg["certificate_authority_domain_config"]; ok {
		ls := caCfgRaw.([]interface{})
		pracc.CertificateAuthorityDomainConfig = make([]*container.CertificateAuthorityDomainConfig, len(ls))
		for i, caCfg := range ls {
			pracc.CertificateAuthorityDomainConfig[i] = expandCADomainConfig(caCfg)
		}
	}

	return pracc
}

func expandCADomainConfig(v interface{}) *container.CertificateAuthorityDomainConfig {
	if v == nil {
		return nil
	}
	cfg := v.(map[string]interface{})

	caConfig := &container.CertificateAuthorityDomainConfig{}
	if v, ok := cfg["fqdns"]; ok {
		fqdns := v.([]interface{})
		caConfig.Fqdns = make([]string, len(fqdns))
		for i, dn := range fqdns {
			caConfig.Fqdns[i] = dn.(string)
		}
	}

	caConfig.GcpSecretManagerCertificateConfig = expandGCPSecretManagerCertificateConfig(cfg["gcp_secret_manager_certificate_config"])

	return caConfig
}

func expandGCPSecretManagerCertificateConfig(v interface{}) *container.GCPSecretManagerCertificateConfig {
	if v == nil {
		return nil
	}
	ls := v.([]interface{})
	if len(ls) == 0 {
		return nil
	}
	if ls[0] == nil {
		return &container.GCPSecretManagerCertificateConfig{}
	}
	cfg := ls[0].(map[string]interface{})

	gcpSMConfig := &container.GCPSecretManagerCertificateConfig{}
	if v, ok := cfg["secret_uri"]; ok {
		gcpSMConfig.SecretUri = v.(string)
	}
	return gcpSMConfig
}

func expandWritableCgroups(v interface{}) *container.WritableCgroups {
	if v == nil {
		return nil
	}
	ls := v.([]interface{})
	if len(ls) == 0 {
		return nil
	}
	if ls[0] == nil {
		return &container.WritableCgroups{}
	}
	cfg := ls[0].(map[string]interface{})

	wcg := &container.WritableCgroups{}
	if enabled, ok := cfg["enabled"]; ok {
		wcg.Enabled = enabled.(bool)
	}
	return wcg
}

func expandRegistryHosts(v interface{}) []*container.RegistryHostConfig {
	if v == nil {
		return nil
	}
	ls := v.([]interface{})
	if len(ls) == 0 {
		return nil
	}
	registryHosts := make([]*container.RegistryHostConfig, 0, len(ls))
	for _, raw := range ls {
		data := raw.(map[string]interface{})
		rh := &container.RegistryHostConfig{
			Server: data["server"].(string),
		}
		if v, ok := data["hosts"]; ok {
			hosts := v.([]interface{})
			rh.Hosts = make([]*container.HostConfig, 0, len(hosts))
			for _, rawHost := range hosts {
				hostData := rawHost.(map[string]interface{})
				h := &container.HostConfig{
					Host: hostData["host"].(string),
				}
				if v, ok := hostData["override_path"]; ok {
					h.OverridePath = v.(bool)
				}
				if v, ok := hostData["dial_timeout"]; ok {
					h.DialTimeout = v.(string)
				}
				if v, ok := hostData["capabilities"]; ok {
					cap := v.([]interface{})
					h.Capabilities = make([]string, len(cap))
					for i, c := range cap {
						h.Capabilities[i] = c.(string)
					}
				}
				if v, ok := hostData["header"]; ok {
					headers := v.([]interface{})
					h.Header = make([]*container.RegistryHeader, len(headers))
					for i, headerRaw := range headers {
						h.Header[i] = expandRegistryHeader(headerRaw)
					}
				}
				if v, ok := hostData["ca"]; ok {
					ca := v.([]interface{})
					h.Ca = make([]*container.CertificateConfig, len(ca))
					for i, caRaw := range ca {
						h.Ca[i] = expandRegistryCertificateConfig(caRaw)
					}
				}
				if v, ok := hostData["client"]; ok {
					client := v.([]interface{})
					h.Client = make([]*container.CertificateConfigPair, len(client))
					for i, clientRaw := range client {
						h.Client[i] = expandRegistryCertificateConfigPair(clientRaw)
					}
				}
				rh.Hosts = append(rh.Hosts, h)
			}
		}
		registryHosts = append(registryHosts, rh)
	}
	return registryHosts
}

func expandRegistryHeader(v interface{}) *container.RegistryHeader {
	header := &container.RegistryHeader{}
	if v == nil {
		return header
	}
	ls := v.(map[string]interface{})
	if val, ok := ls["key"]; ok {
		header.Key = val.(string)
	}
	if val, ok := ls["value"]; ok {
		headerVal := val.([]interface{})
		header.Value = make([]string, len(headerVal))
		for i, hv := range headerVal {
			header.Value[i] = hv.(string)
		}
	}
	return header
}

func expandRegistryCertificateConfig(v interface{}) *container.CertificateConfig {
	cfg := &container.CertificateConfig{}
	if v == nil {
		return cfg
	}
	ls := v.(map[string]interface{})
	if val, ok := ls["gcp_secret_manager_secret_uri"]; ok {
		cfg.GcpSecretManagerSecretUri = val.(string)
	}
	return cfg
}

func expandRegistryCertificateConfigPair(v interface{}) *container.CertificateConfigPair {
	cfg := &container.CertificateConfigPair{}
	if v == nil {
		return cfg
	}
	ls := v.(map[string]interface{})
	if val, ok := ls["cert"]; ok {
		certRaw := val.([]interface{})
		if len(certRaw) > 0 {
			cfg.Cert = expandRegistryCertificateConfig(certRaw[0])
		}
	}
	if val, ok := ls["key"]; ok {
		keyRaw := val.([]interface{})
		if len(keyRaw) > 0 {
			cfg.Key = expandRegistryCertificateConfig(keyRaw[0])
		}
	}
	return cfg
}

func expandSoleTenantConfig(v interface{}) *container.SoleTenantConfig {
	if v == nil {
		return nil
	}
	ls := v.([]interface{})
	if len(ls) == 0 {
		return nil
	}
	stConfig := &container.SoleTenantConfig{}
	cfg := ls[0].(map[string]interface{})
	if affinitiesRaw, ok := cfg["node_affinity"]; ok {
		affinities := make([]*container.NodeAffinity, 0)
		for _, v := range affinitiesRaw.(*schema.Set).List() {
			na := v.(map[string]interface{})
			affinities = append(affinities, &container.NodeAffinity{
				Key:      na["key"].(string),
				Operator: na["operator"].(string),
				Values:   tpgresource.ConvertStringArr(na["values"].([]interface{})),
			})
		}
		stConfig.NodeAffinities = affinities
	}
	if v, ok := cfg["min_node_cpus"]; ok {
		stConfig.MinNodeCpus = int64(v.(int))
	}
	return stConfig
}

func expandConfidentialNodes(configured interface{}) *container.ConfidentialNodes {
	l := configured.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}
	config := l[0].(map[string]interface{})
	return &container.ConfidentialNodes{
		Enabled:                  config["enabled"].(bool),
		ConfidentialInstanceType: config["confidential_instance_type"].(string),
	}
}

func flattenNodeConfigDefaults(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	transformed := map[string]interface{}{}
	transformed["containerd_config"] = flattenContainerdConfig(c["containerdConfig"])
	transformed["insecure_kubelet_readonly_port_enabled"] = flattenInsecureKubeletReadonlyPortEnabled(c["nodeKubeletConfig"])
	transformed["logging_variant"] = flattenLoggingVariant(c["loggingConfig"])
	transformed["gcfs_config"] = flattenGcfsConfig(c["gcfsConfig"])

	return []map[string]interface{}{transformed}
}

func flattenNodeConfig(v interface{}, _ interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}

	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	transformed := map[string]interface{}{
		"machine_type":                       c["machineType"],
		"containerd_config":                  flattenContainerdConfig(c["containerdConfig"]),
		"disk_size_gb":                       c["diskSizeGb"],
		"disk_type":                          c["diskType"],
		"boot_disk":                          flattenBootDiskConfig(c["bootDisk"]),
		"guest_accelerator":                  flattenContainerGuestAccelerators(c["accelerators"]),
		"local_ssd_count":                    c["localSsdCount"],
		"logging_variant":                    flattenLoggingVariant(c["loggingConfig"]),
		"local_nvme_ssd_block_config":        flattenLocalNvmeSsdBlockConfig(c["localNvmeSsdBlockConfig"]),
		"gcfs_config":                        flattenGcfsConfig(c["gcfsConfig"]),
		"ephemeral_storage_local_ssd_config": flattenEphemeralStorageLocalSsdConfig(c["ephemeralStorageLocalSsdConfig"]),
		"gvnic":                              flattenGvnic(c["gvnic"]),
		"reservation_affinity":               flattenGKEReservationAffinity(c["reservationAffinity"]),
		"service_account":                    c["serviceAccount"],
		"metadata":                           c["metadata"],
		"image_type":                         c["imageType"],
		"labels":                             c["labels"],
		"resource_labels":                    c["resourceLabels"],
		"tags":                               c["tags"],
		"preemptible":                        c["preemptible"],
		"secondary_boot_disks":               flattenSecondaryBootDisks(c["secondaryBootDisks"]),
		"storage_pools":                      c["storagePools"],
		"spot":                               c["spot"],
		"min_cpu_platform":                   c["minCpuPlatform"],
		"shielded_instance_config":           flattenShieldedInstanceConfig(c["shieldedInstanceConfig"]),
		"taint":                              flattenEffectiveTaints(c["taints"]),
		"workload_metadata_config":           flattenWorkloadMetadataConfig(c["workloadMetadataConfig"]),
		"confidential_nodes":                 flattenConfidentialNodes(c["confidentialNodes"]),
		"boot_disk_kms_key":                  c["bootDiskKmsKey"],
		"kubelet_config":                     flattenKubeletConfig(c["kubeletConfig"]),
		"linux_node_config":                  flattenLinuxNodeConfig(c["linuxNodeConfig"]),
		"windows_node_config":                flattenWindowsNodeConfig(c["windowsNodeConfig"]),
		"node_group":                         c["nodeGroup"],
		"advanced_machine_features":          flattenAdvancedMachineFeaturesConfig(c["advancedMachineFeatures"]),
		"max_run_duration":                   c["maxRunDuration"],
		"flex_start":                         c["flexStart"],
		"sole_tenant_config":                 flattenSoleTenantConfig(c["soleTenantConfig"]),
		"fast_socket":                        flattenFastSocket(c["fastSocket"]),
		"resource_manager_tags":              flattenResourceManagerTags(c["resourceManagerTags"]),
		"enable_confidential_storage":        c["enableConfidentialStorage"],
		"local_ssd_encryption_mode":          c["localSsdEncryptionMode"],
	}

	if v, ok := c["oauthScopes"].([]interface{}); ok && len(v) > 0 {
		transformed["oauth_scopes"] = v
	}

	return []map[string]interface{}{transformed}
}

func flattenBootDiskConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	transformed := map[string]interface{}{
		"disk_type":              c["diskType"],
		"size_gb":                c["sizeGb"],
		"provisioned_iops":       c["provisionedIops"],
		"provisioned_throughput": c["provisionedThroughput"],
	}

	return []map[string]interface{}{transformed}
}

func flattenResourceManagerTags(v interface{}) map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	rmt := make(map[string]interface{})
	if tags, ok := c["tags"].(map[string]interface{}); ok {
		for k, val := range tags {
			rmt[k] = val
		}
	}

	return rmt
}

func flattenAdvancedMachineFeaturesConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	transformed := map[string]interface{}{
		"threads_per_core":             c["threadsPerCore"],
		"enable_nested_virtualization": c["enableNestedVirtualization"],
		"performance_monitoring_unit":  c["performanceMonitoringUnit"],
	}

	return []map[string]interface{}{transformed}
}

func flattenContainerGuestAccelerators(v interface{}) []map[string]interface{} {
	result := []map[string]interface{}{}
	if v == nil {
		return result
	}
	c, ok := v.([]interface{})
	if !ok {
		return result
	}

	for _, item := range c {
		accel, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		accelerator := map[string]interface{}{
			"count":              accel["acceleratorCount"],
			"type":               accel["acceleratorType"],
			"gpu_partition_size": accel["gpuPartitionSize"],
		}
		if v, ok := accel["gpuDriverInstallationConfig"].(map[string]interface{}); ok {
			accelerator["gpu_driver_installation_config"] = []map[string]interface{}{
				{
					"gpu_driver_version": v["gpuDriverVersion"],
				},
			}
		}
		if v, ok := accel["gpuSharingConfig"].(map[string]interface{}); ok {
			accelerator["gpu_sharing_config"] = []map[string]interface{}{
				{
					"gpu_sharing_strategy":       v["gpuSharingStrategy"],
					"max_shared_clients_per_gpu": v["maxSharedClientsPerGpu"],
				},
			}
		}
		result = append(result, accelerator)
	}
	return result
}

func flattenShieldedInstanceConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	transformed := map[string]interface{}{
		"enable_secure_boot":          c["enableSecureBoot"],
		"enable_integrity_monitoring": c["enableIntegrityMonitoring"],
	}

	return []map[string]interface{}{transformed}
}

func flattenLocalNvmeSsdBlockConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	transformed := map[string]interface{}{
		"local_ssd_count": c["localSsdCount"],
	}

	return []map[string]interface{}{transformed}
}

func flattenEphemeralStorageLocalSsdConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	transformed := map[string]interface{}{
		"local_ssd_count":  c["localSsdCount"],
		"data_cache_count": c["dataCacheCount"],
	}

	return []map[string]interface{}{transformed}
}

func flattenSecondaryBootDisks(v interface{}) []map[string]interface{} {
	result := []map[string]interface{}{}
	if v == nil {
		return result
	}
	c, ok := v.([]interface{})
	if !ok {
		return result
	}

	for _, item := range c {
		disk, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		secondaryBootDisk := map[string]interface{}{
			"disk_image": disk["diskImage"],
			"mode":       disk["mode"],
		}
		if disk["diskImage"] == "" || disk["diskImage"] == nil { // required field
			secondaryBootDisk["disk_image"] = ""
		}
		result = append(result, secondaryBootDisk)
	}
	return result
}

func flattenInsecureKubeletReadonlyPortEnabled(v interface{}) string {
	// Convert bool from the API to the enum values used internally
	if v == nil {
		return "FALSE"
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return "FALSE"
	}
	if val, ok := c["insecureKubeletReadonlyPortEnabled"].(bool); ok && val {
		return "TRUE"
	}
	return "FALSE"
}

func flattenLoggingVariant(v interface{}) string {
	variant := "DEFAULT"
	if v == nil {
		return variant
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return variant
	}
	if vc, ok := c["variantConfig"].(map[string]interface{}); ok {
		if v, ok := vc["variant"].(string); ok && v != "" {
			variant = v
		}
	}
	return variant
}

func flattenGcfsConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	enabled, ok := c["enabled"].(bool)
	if !ok {
		enabled = false
	}
	transformed := map[string]interface{}{
		"enabled": enabled,
	}

	return []map[string]interface{}{transformed}
}

func flattenGvnic(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	transformed := map[string]interface{}{
		"enabled": c["enabled"],
	}

	return []map[string]interface{}{transformed}
}

func flattenGKEReservationAffinity(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	transformed := map[string]interface{}{
		"consume_reservation_type": c["consumeReservationType"],
		"key":                      c["key"],
		"values":                   c["values"],
	}

	return []map[string]interface{}{transformed}
}

// flattenTaints records the set of taints already present in state.
func flattenTaints(c []*container.NodeTaint, oldTaints []interface{}) []map[string]interface{} {
	taintKeys := map[string]struct{}{}
	for _, raw := range oldTaints {
		data := raw.(map[string]interface{})
		taintKey := data["key"].(string)
		taintKeys[taintKey] = struct{}{}
	}

	result := []map[string]interface{}{}
	for _, taint := range c {
		if _, ok := taintKeys[taint.Key]; ok {
			result = append(result, map[string]interface{}{
				"key":    taint.Key,
				"value":  taint.Value,
				"effect": taint.Effect,
			})
		}
	}

	return result
}

// flattenEffectiveTaints records the complete set of taints returned from GKE.
func flattenEffectiveTaints(v interface{}) []map[string]interface{} {
	result := []map[string]interface{}{}
	if v == nil {
		return result
	}
	c, ok := v.([]interface{})
	if !ok {
		return result
	}
	for _, raw := range c {
		taint, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		result = append(result, map[string]interface{}{
			"key":    taint["key"],
			"value":  taint["value"],
			"effect": taint["effect"],
		})
	}

	return result
}

func flattenWorkloadMetadataConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	transformed := map[string]interface{}{
		"mode": c["mode"],
	}
	return []map[string]interface{}{transformed}
}

func containerNodePoolLabelsSuppress(k, old, new string, d *schema.ResourceData) bool {
	// Node configs are embedded into multiple resources (container cluster and
	// container node pool) so we determine the node config key dynamically.
	idx := strings.Index(k, ".labels.")
	if idx < 0 {
		return false
	}

	root := k[:idx]

	// Right now, GKE only applies its own out-of-band labels when you enable
	// Sandbox. We only need to perform diff suppression in this case;
	// otherwise, the default Terraform behavior is fine.
	o, n := d.GetChange(root + ".sandbox_config.0.sandbox_type")
	if o == nil || n == nil {
		return false
	}

	// Pull the entire changeset as a list rather than trying to deal with each
	// element individually.
	o, n = d.GetChange(root + ".labels")
	if o == nil || n == nil {
		return false
	}

	labels := n.(map[string]interface{})

	// Remove all current labels, skipping GKE-managed ones if not present in
	// the new configuration.
	for key, value := range o.(map[string]interface{}) {
		if nv, ok := labels[key]; ok && nv == value {
			delete(labels, key)
		} else if !strings.HasPrefix(key, "sandbox.gke.io/") {
			// User-provided label removed in new configuration.
			return false
		}
	}

	// If, at this point, the map still has elements, the new configuration
	// added an additional taint.
	if len(labels) > 0 {
		return false
	}

	return true
}

func containerNodePoolResourceLabelsDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	// Suppress diffs for server-specified labels prefixed with "goog-gke"
	if strings.Contains(k, "resource_labels.goog-gke") && new == "" {
		return true
	}

	// Let diff be determined by resource_labels (above)
	if strings.Contains(k, "resource_labels.%") {
		return true
	}

	// For other keys, don't suppress diff.
	return false
}

func flattenKubeletConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	transformed := map[string]interface{}{
		"cpu_cfs_quota":                          c["cpuCfsQuota"],
		"cpu_cfs_quota_period":                   c["cpuCfsQuotaPeriod"],
		"cpu_manager_policy":                     c["cpuManagerPolicy"],
		"memory_manager":                         flattenMemoryManager(c["memoryManager"]),
		"topology_manager":                       flattenTopologyManager(c["topologyManager"]),
		"insecure_kubelet_readonly_port_enabled": flattenInsecureKubeletReadonlyPortEnabled(v),
		"pod_pids_limit":                         c["podPidsLimit"],
		"container_log_max_size":                 c["containerLogMaxSize"],
		"container_log_max_files":                c["containerLogMaxFiles"],
		"image_gc_low_threshold_percent":         c["imageGcLowThresholdPercent"],
		"image_gc_high_threshold_percent":        c["imageGcHighThresholdPercent"],
		"image_minimum_gc_age":                   c["imageMinimumGcAge"],
		"image_maximum_gc_age":                   c["imageMaximumGcAge"],
		"allowed_unsafe_sysctls":                 c["allowedUnsafeSysctls"],
		"single_process_oom_kill":                c["singleProcessOomKill"],
		"max_parallel_image_pulls":               c["maxParallelImagePulls"],
		"eviction_max_pod_grace_period_seconds":  c["evictionMaxPodGracePeriodSeconds"],
		"eviction_soft":                          flattenEvictionSignals(c["evictionSoft"]),
		"eviction_soft_grace_period":             flattenEvictionGracePeriod(c["evictionSoftGracePeriod"]),
		"eviction_minimum_reclaim":               flattenEvictionMinimumReclaim(c["evictionMinimumReclaim"]),
	}

	return []map[string]interface{}{transformed}
}

func flattenTopologyManager(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	transformed := map[string]interface{}{
		"policy": c["policy"],
		"scope":  c["scope"],
	}

	return []map[string]interface{}{transformed}
}

func flattenMemoryManager(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	transformed := map[string]interface{}{
		"policy": c["policy"],
	}

	return []map[string]interface{}{transformed}
}

func flattenNodePoolAutoConfigNodeKubeletConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	transformed := map[string]interface{}{}
	if c != nil {
		transformed["insecure_kubelet_readonly_port_enabled"] = flattenInsecureKubeletReadonlyPortEnabled(c)
	}

	return []map[string]interface{}{transformed}
}

func flattenEvictionSignals(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	transformed := map[string]interface{}{
		"memory_available":    c["memoryAvailable"],
		"nodefs_available":    c["nodefsAvailable"],
		"nodefs_inodes_free":  c["nodefsInodesFree"],
		"imagefs_available":   c["imagefsAvailable"],
		"imagefs_inodes_free": c["imagefsInodesFree"],
		"pid_available":       c["pidAvailable"],
	}

	return []map[string]interface{}{transformed}
}

func flattenEvictionGracePeriod(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	transformed := map[string]interface{}{
		"memory_available":    c["memoryAvailable"],
		"nodefs_available":    c["nodefsAvailable"],
		"nodefs_inodes_free":  c["nodefsInodesFree"],
		"imagefs_available":   c["imagefsAvailable"],
		"imagefs_inodes_free": c["imagefsInodesFree"],
		"pid_available":       c["pidAvailable"],
	}

	return []map[string]interface{}{transformed}
}

func flattenEvictionMinimumReclaim(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	transformed := map[string]interface{}{
		"memory_available":    c["memoryAvailable"],
		"nodefs_available":    c["nodefsAvailable"],
		"nodefs_inodes_free":  c["nodefsInodesFree"],
		"imagefs_available":   c["imagefsAvailable"],
		"imagefs_inodes_free": c["imagefsInodesFree"],
		"pid_available":       c["pidAvailable"],
	}

	return []map[string]interface{}{transformed}
}

func flattenLinuxNodeConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	sysctls := make(map[string]interface{})
	if rawSysctls, ok := c["sysctls"].(map[string]interface{}); ok {
		for k, val := range rawSysctls {
			sysctls[k] = val
		}
	}

	transformed := map[string]interface{}{
		"sysctls":                      sysctls,
		"cgroup_mode":                  c["cgroupMode"],
		"hugepages_config":             flattenHugepagesConfig(c["hugepages"]),
		"transparent_hugepage_enabled": c["transparentHugepageEnabled"],
		"transparent_hugepage_defrag":  c["transparentHugepageDefrag"],
		"node_kernel_module_loading":   flattenNodeKernelModuleLoading(c["nodeKernelModuleLoading"]),
	}

	return []map[string]interface{}{transformed}
}

func flattenWindowsNodeConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(c) == 0 {
		return nil
	}
	transformed := map[string]interface{}{
		"osversion": c["osVersion"],
	}

	return []map[string]interface{}{transformed}
}

func flattenHugepagesConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	transformed := map[string]interface{}{
		"hugepage_size_2m": c["hugepageSize2m"],
		"hugepage_size_1g": c["hugepageSize1g"],
	}

	return []map[string]interface{}{transformed}
}

func flattenNodeKernelModuleLoading(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	r := make(map[string]interface{})
	if val, ok := c["policy"]; ok {
		r["policy"] = val
	}
	return []map[string]interface{}{r}
}

func flattenContainerdConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	r := map[string]interface{}{}
	if c["privateRegistryAccessConfig"] != nil {
		r["private_registry_access_config"] = flattenPrivateRegistryAccessConfig(c["privateRegistryAccessConfig"])
	}
	if c["writableCgroups"] != nil {
		r["writable_cgroups"] = flattenWritableCgroups(c["writableCgroups"])
	}
	if c["registryHosts"] != nil {
		r["registry_hosts"] = flattenRegistryHosts(c["registryHosts"])
	}
	return []map[string]interface{}{r}
}

func flattenRegistryHosts(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	ls, ok := v.([]interface{})
	if !ok {
		return nil
	}
	if len(ls) == 0 {
		return nil
	}

	items := []map[string]interface{}{}
	for _, raw := range ls {
		host, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		item := make(map[string]interface{})
		item["server"] = host["server"]
		if host["hosts"] != nil {
			item["hosts"] = flattenHostInRegistryHosts(host["hosts"])
		}
		items = append(items, item)
	}
	if len(items) == 0 {
		return nil
	}
	return items
}

func flattenHostInRegistryHosts(v interface{}) []map[string]interface{} {
	if v == nil {
		return []map[string]interface{}{}
	}
	ls, ok := v.([]interface{})
	if !ok {
		return []map[string]interface{}{}
	}
	items := make([]map[string]interface{}, 0, len(ls))
	for _, raw := range ls {
		h, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		item := make(map[string]interface{})
		item["host"] = h["host"]
		item["capabilities"] = h["capabilities"]
		item["override_path"] = h["overridePath"]
		item["dial_timeout"] = h["dialTimeout"]

		if h["header"] != nil {
			if headers, ok := h["header"].([]interface{}); ok {
				tmp := make([]interface{}, len(headers))
				for i, rawVal := range headers {
					if val, ok := rawVal.(map[string]interface{}); ok {
						tmp[i] = map[string]interface{}{
							"key":   val["key"],
							"value": val["value"],
						}
					}
				}
				item["header"] = tmp
			}
		}

		if h["ca"] != nil {
			if cas, ok := h["ca"].([]interface{}); ok {
				tmp := make([]interface{}, len(cas))
				for i, rawVal := range cas {
					if val, ok := rawVal.(map[string]interface{}); ok {
						if uri, ok := val["gcpSecretManagerSecretUri"].(string); ok && uri != "" {
							tmp[i] = map[string]interface{}{
								"gcp_secret_manager_secret_uri": uri,
							}
						}
					}
				}
				item["ca"] = tmp
			}
		}

		if h["client"] != nil {
			if clients, ok := h["client"].([]interface{}); ok {
				tmp := make([]interface{}, len(clients))
				for i, rawVal := range clients {
					if val, ok := rawVal.(map[string]interface{}); ok {
						currentClient := map[string]interface{}{}
						if certRaw, ok := val["cert"].(map[string]interface{}); ok && certRaw["gcpSecretManagerSecretUri"] != "" {
							currentClient["cert"] = []interface{}{
								map[string]interface{}{
									"gcp_secret_manager_secret_uri": certRaw["gcpSecretManagerSecretUri"],
								},
							}
						}

						if keyRaw, ok := val["key"].(map[string]interface{}); ok && keyRaw["gcpSecretManagerSecretUri"] != "" {
							currentClient["key"] = []interface{}{
								map[string]interface{}{
									"gcp_secret_manager_secret_uri": keyRaw["gcpSecretManagerSecretUri"],
								},
							}
						}
						tmp[i] = currentClient
					}
				}
				item["client"] = tmp
			}
		}
		items = append(items, item)
	}
	return items
}

func flattenPrivateRegistryAccessConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return []map[string]interface{}{}
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return []map[string]interface{}{}
	}
	enabled, ok := c["enabled"].(bool)
	if !ok {
		enabled = false
	}
	r := map[string]interface{}{
		"enabled": enabled,
	}
	if c["certificateAuthorityDomainConfig"] != nil {
		if caConfigs, ok := c["certificateAuthorityDomainConfig"].([]interface{}); ok {
			flattenedCaConfigs := make([]interface{}, 0, len(caConfigs))
			for _, caCfg := range caConfigs {
				flattened := flattenCADomainConfig(caCfg)
				if len(flattened) > 0 {
					flattenedCaConfigs = append(flattenedCaConfigs, flattened[0])
				}
			}
			r["certificate_authority_domain_config"] = flattenedCaConfigs
		}
	}
	return []map[string]interface{}{r}
}

func flattenCADomainConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return []map[string]interface{}{}
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return []map[string]interface{}{}
	}
	r := map[string]interface{}{
		"fqdns": c["fqdns"],
	}
	if c["gcpSecretManagerCertificateConfig"] != nil {
		r["gcp_secret_manager_certificate_config"] = flattenGCPSecretManagerCertificateConfig(c["gcpSecretManagerCertificateConfig"])
	}
	return []map[string]interface{}{r}
}

func flattenGCPSecretManagerCertificateConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return []map[string]interface{}{}
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return []map[string]interface{}{}
	}
	r := map[string]interface{}{
		"secret_uri": c["secretUri"],
	}
	return []map[string]interface{}{r}
}

func flattenWritableCgroups(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	enabled, ok := c["enabled"].(bool)
	if !ok {
		enabled = false
	}
	transformed := map[string]interface{}{
		"enabled": enabled,
	}

	return []map[string]interface{}{transformed}
}

func flattenConfidentialNodes(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	enabled, ok := c["enabled"].(bool)
	if !ok {
		enabled = false
	}
	transformed := map[string]interface{}{
		"enabled":                    enabled,
		"confidential_instance_type": c["confidentialInstanceType"],
	}

	return []map[string]interface{}{transformed}
}

func flattenSoleTenantConfig(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	affinities := []map[string]interface{}{}
	if rawAffinities, ok := c["nodeAffinities"].([]interface{}); ok {
		for _, raw := range rawAffinities {
			affinity, ok := raw.(map[string]interface{})
			if !ok {
				continue
			}
			affinities = append(affinities, map[string]interface{}{
				"key":      affinity["key"],
				"operator": affinity["operator"],
				"values":   affinity["values"],
			})
		}
	}
	transformed := map[string]interface{}{
		"node_affinity": affinities,
		"min_node_cpus": c["minNodeCpus"],
	}

	return []map[string]interface{}{transformed}
}

func flattenFastSocket(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	c, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	transformed := map[string]interface{}{
		"enabled": c["enabled"],
	}

	return []map[string]interface{}{transformed}
}
