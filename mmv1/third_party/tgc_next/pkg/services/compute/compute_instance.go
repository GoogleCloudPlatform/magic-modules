package compute

import (
	"strings"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/cai2hcl/converters/utils"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/tpgresource"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/verify"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	compute "google.golang.org/api/compute/v0.beta"
)

// ComputeInstanceAssetType is the CAI asset type name for compute instance.
const ComputeInstanceAssetType string = "compute.googleapis.com/Instance"

// ComputeInstanceSchemaName is the TF resource schema name for compute instance.
const ComputeInstanceSchemaName string = "google_compute_instance"

var (
	advancedMachineFeaturesKeys = []string{
		"advanced_machine_features.0.enable_nested_virtualization",
		"advanced_machine_features.0.threads_per_core",
		"advanced_machine_features.0.turbo_mode",
		"advanced_machine_features.0.visible_core_count",
		"advanced_machine_features.0.performance_monitoring_unit",
		"advanced_machine_features.0.enable_uefi_networking",
	}

	bootDiskKeys = []string{
		"boot_disk.0.guest_os_features",
		"boot_disk.0.auto_delete",
		"boot_disk.0.device_name",
		"boot_disk.0.disk_encryption_key_raw",
		"boot_disk.0.kms_key_self_link",
		"boot_disk.0.disk_encryption_key_rsa",
		"boot_disk.0.disk_encryption_service_account",
		"boot_disk.0.initialize_params",
		"boot_disk.0.mode",
		"boot_disk.0.source",
	}

	initializeParamsKeys = []string{
		"boot_disk.0.initialize_params.0.size",
		"boot_disk.0.initialize_params.0.type",
		"boot_disk.0.initialize_params.0.image",
		"boot_disk.0.initialize_params.0.labels",
		"boot_disk.0.initialize_params.0.resource_manager_tags",
		"boot_disk.0.initialize_params.0.provisioned_iops",
		"boot_disk.0.initialize_params.0.provisioned_throughput",
		"boot_disk.0.initialize_params.0.enable_confidential_compute",
		"boot_disk.0.initialize_params.0.source_image_encryption_key",
		"boot_disk.0.initialize_params.0.snapshot",
		"boot_disk.0.initialize_params.0.source_snapshot_encryption_key",
		"boot_disk.0.initialize_params.0.storage_pool",
		"boot_disk.0.initialize_params.0.resource_policies",
		"boot_disk.0.initialize_params.0.architecture",
	}

	schedulingKeys = []string{
		"scheduling.0.on_host_maintenance",
		"scheduling.0.automatic_restart",
		"scheduling.0.preemptible",
		"scheduling.0.node_affinities",
		"scheduling.0.min_node_cpus",
		"scheduling.0.provisioning_model",
		"scheduling.0.instance_termination_action",
		"scheduling.0.termination_time",
		"scheduling.0.availability_domain",
		"scheduling.0.max_run_duration",
		"scheduling.0.on_instance_stop_action",
		"scheduling.0.maintenance_interval",
		"scheduling.0.host_error_timeout_seconds",
		"scheduling.0.graceful_shutdown",
		"scheduling.0.local_ssd_recovery_timeout",
	}

	shieldedInstanceConfigKeys = []string{
		"shielded_instance_config.0.enable_secure_boot",
		"shielded_instance_config.0.enable_vtpm",
		"shielded_instance_config.0.enable_integrity_monitoring",
	}
)

func ResourceComputeInstance() *schema.Resource {
	return &schema.Resource{
		// A compute instance is more or less a superset of a compute instance
		// template. Please attempt to maintain consistency with the
		// resource_compute_instance_template schema when updating this one.
		Schema: map[string]*schema.Schema{
			"boot_disk": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: `The boot disk for the instance.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_delete": {
							Type:         schema.TypeBool,
							Optional:     true,
							AtLeastOneOf: bootDiskKeys,
							Default:      true,
							Description:  `Whether the disk will be auto-deleted when the instance is deleted.`,
						},

						"device_name": {
							Type:         schema.TypeString,
							Optional:     true,
							AtLeastOneOf: bootDiskKeys,
							Computed:     true,
							ForceNew:     true,
							Description:  `Name with which attached disk will be accessible under /dev/disk/by-id/`,
						},

						"disk_encryption_key_raw": {
							Type:          schema.TypeString,
							Optional:      true,
							AtLeastOneOf:  bootDiskKeys,
							ForceNew:      true,
							ConflictsWith: []string{"boot_disk.0.kms_key_self_link", "boot_disk.0.disk_encryption_key_rsa"},
							Sensitive:     true,
							Description:   `A 256-bit customer-supplied encryption key, encoded in RFC 4648 base64 to encrypt this disk. Only one of kms_key_self_link, disk_encryption_key_raw and disk_encryption_key_rsa may be set.`,
						},

						"disk_encryption_key_rsa": {
							Type:          schema.TypeString,
							Optional:      true,
							AtLeastOneOf:  bootDiskKeys,
							ForceNew:      true,
							ConflictsWith: []string{"boot_disk.0.kms_key_self_link", "boot_disk.0.disk_encryption_key_raw"},
							Sensitive:     true,
							Description:   `Specifies an RFC 4648 base64 encoded, RSA-wrapped 2048-bit customer-supplied encryption key to either encrypt or decrypt this resource. Only one of kms_key_self_link, disk_encryption_key_raw and disk_encryption_key_rsa may be set.`,
						},

						"disk_encryption_key_sha256": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The RFC 4648 base64 encoded SHA-256 hash of the customer-supplied encryption key that protects this resource.`,
						},

						"disk_encryption_service_account": {
							Type:         schema.TypeString,
							Optional:     true,
							AtLeastOneOf: bootDiskKeys,
							ForceNew:     true,
							Description:  `The service account being used for the encryption request for the given KMS key. If absent, the Compute Engine default service account is used`,
						},

						"interface": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"SCSI", "NVME"}, false),
							Description:  `The disk interface used for attaching this disk. One of SCSI or NVME. (This field is shared with attached_disk and only used for specific cases, please don't specify this field without advice from Google.)`,
						},

						"kms_key_self_link": {
							Type:             schema.TypeString,
							Optional:         true,
							AtLeastOneOf:     bootDiskKeys,
							ForceNew:         true,
							ConflictsWith:    []string{"boot_disk.0.disk_encryption_key_raw", "boot_disk.0.disk_encryption_key_rsa"},
							DiffSuppressFunc: tpgresource.CompareSelfLinkRelativePaths,
							Computed:         true,
							Description:      `The self_link of the encryption key that is stored in Google Cloud KMS to encrypt this disk. Only one of kms_key_self_link, disk_encryption_key_raw and disk_encryption_key_rsa may be set.`,
						},

						"guest_os_features": {
							Type:         schema.TypeList,
							Optional:     true,
							AtLeastOneOf: bootDiskKeys,
							ForceNew:     true,
							Computed:     true,
							Description:  `A list of features to enable on the guest operating system. Applicable only for bootable images.`,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						"initialize_params": {
							Type:         schema.TypeList,
							Optional:     true,
							AtLeastOneOf: bootDiskKeys,
							Computed:     true,
							ForceNew:     true,
							MaxItems:     1,
							Description:  `Parameters with which a disk was created alongside the instance.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"size": {
										Type:         schema.TypeInt,
										Optional:     true,
										AtLeastOneOf: initializeParamsKeys,
										Computed:     true,
										ForceNew:     true,
										ValidateFunc: validation.IntAtLeast(1),
										Description:  `The size of the image in gigabytes.`,
									},

									"type": {
										Type:         schema.TypeString,
										Optional:     true,
										AtLeastOneOf: initializeParamsKeys,
										Computed:     true,
										ForceNew:     true,
										Description:  `The Google Compute Engine disk type. Such as pd-standard, pd-ssd or pd-balanced.`,
									},

									"image": {
										Type:         schema.TypeString,
										Optional:     true,
										AtLeastOneOf: initializeParamsKeys,
										Computed:     true,
										ForceNew:     true,
										Description:  `The image from which this disk was initialised.`,
									},

									"source_image_encryption_key": {
										Type:         schema.TypeList,
										Optional:     true,
										AtLeastOneOf: initializeParamsKeys,
										MaxItems:     1,
										Description:  `The encryption key used to decrypt the source image.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"raw_key": {
													Type:        schema.TypeString,
													Optional:    true,
													ForceNew:    true,
													Sensitive:   true,
													Description: `Specifies a 256-bit customer-supplied encryption key, encoded in RFC 4648 base64 to either encrypt or decrypt this resource. Only one of kms_key_self_link, rsa_encrypted_key and raw_key may be set.`,
												},

												"rsa_encrypted_key": {
													Type:        schema.TypeString,
													Optional:    true,
													ForceNew:    true,
													Sensitive:   true,
													Description: `Specifies an RFC 4648 base64 encoded, RSA-wrapped 2048-bit customer-supplied encryption key to either encrypt or decrypt this resource. Only one of kms_key_self_link, rsa_encrypted_key and raw_key may be set.`,
												},

												"kms_key_self_link": {
													Type:             schema.TypeString,
													Optional:         true,
													ForceNew:         true,
													Computed:         true,
													DiffSuppressFunc: tpgresource.CompareCryptoKeyVersions,
													Description:      `The self link of the encryption key that is stored in Google Cloud KMS. Only one of kms_key_self_link, rsa_encrypted_key and raw_key may be set.`,
												},

												"kms_key_service_account": {
													Type:        schema.TypeString,
													Optional:    true,
													ForceNew:    true,
													Description: `The service account being used for the encryption request for the given KMS key. If absent, the Compute Engine default service account is used.`,
												},

												"sha256": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: `The SHA256 hash of the encryption key used to encrypt this disk.`,
												},
											},
										},
									},

									"snapshot": {
										Type:             schema.TypeString,
										Optional:         true,
										AtLeastOneOf:     initializeParamsKeys,
										Computed:         true,
										ForceNew:         true,
										DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
										Description:      `The snapshot from which this disk was initialised.`,
									},

									"source_snapshot_encryption_key": {
										Type:         schema.TypeList,
										Optional:     true,
										AtLeastOneOf: initializeParamsKeys,
										MaxItems:     1,
										Description:  `The encryption key used to decrypt the source snapshot.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"raw_key": {
													Type:        schema.TypeString,
													Optional:    true,
													ForceNew:    true,
													Sensitive:   true,
													Description: `Specifies a 256-bit customer-supplied encryption key, encoded in RFC 4648 base64 to either encrypt or decrypt this resource. Only one of kms_key_self_link, rsa_encrypted_key and raw_key may be set.`,
												},

												"rsa_encrypted_key": {
													Type:        schema.TypeString,
													Optional:    true,
													ForceNew:    true,
													Sensitive:   true,
													Description: `Specifies an RFC 4648 base64 encoded, RSA-wrapped 2048-bit customer-supplied encryption key to either encrypt or decrypt this resource. Only one of kms_key_self_link, rsa_encrypted_key and raw_key may be set.`,
												},

												"kms_key_self_link": {
													Type:             schema.TypeString,
													Optional:         true,
													ForceNew:         true,
													Computed:         true,
													DiffSuppressFunc: tpgresource.CompareCryptoKeyVersions,
													Description:      `The self link of the encryption key that is stored in Google Cloud KMS. Only one of kms_key_self_link, rsa_encrypted_key and raw_key may be set.`,
												},

												"kms_key_service_account": {
													Type:        schema.TypeString,
													Optional:    true,
													ForceNew:    true,
													Description: `The service account being used for the encryption request for the given KMS key. If absent, the Compute Engine default service account is used.`,
												},

												"sha256": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: `The SHA256 hash of the encryption key used to encrypt this disk.`,
												},
											},
										},
									},

									"labels": {
										Type:         schema.TypeMap,
										Optional:     true,
										AtLeastOneOf: initializeParamsKeys,
										Computed:     true,
										ForceNew:     true,
										Description:  `A set of key/value label pairs assigned to the disk.`,
									},

									"resource_manager_tags": {
										Type:         schema.TypeMap,
										Optional:     true,
										ForceNew:     true,
										AtLeastOneOf: initializeParamsKeys,
										Description:  `A map of resource manager tags. Resource manager tag keys and values have the same definition as resource manager tags. Keys must be in the format tagKeys/{tag_key_id}, and values are in the format tagValues/456. The field is ignored (both PUT & PATCH) when empty.`,
									},

									"resource_policies": {
										Type:             schema.TypeList,
										Elem:             &schema.Schema{Type: schema.TypeString},
										Optional:         true,
										ForceNew:         true,
										Computed:         true,
										AtLeastOneOf:     initializeParamsKeys,
										DiffSuppressFunc: tpgresource.CompareSelfLinkRelativePaths,
										MaxItems:         1,
										Description:      `A list of self_links of resource policies to attach to the instance's boot disk. Modifying this list will cause the instance to recreate. Currently a max of 1 resource policy is supported.`,
									},

									"provisioned_iops": {
										Type:         schema.TypeInt,
										Optional:     true,
										AtLeastOneOf: initializeParamsKeys,
										Computed:     true,
										ForceNew:     true,
										Description:  `Indicates how many IOPS to provision for the disk. This sets the number of I/O operations per second that the disk can handle.`,
									},

									"provisioned_throughput": {
										Type:         schema.TypeInt,
										Optional:     true,
										AtLeastOneOf: initializeParamsKeys,
										Computed:     true,
										ForceNew:     true,
										Description:  `Indicates how much throughput to provision for the disk. This sets the number of throughput mb per second that the disk can handle.`,
									},

									"enable_confidential_compute": {
										Type:         schema.TypeBool,
										Optional:     true,
										AtLeastOneOf: initializeParamsKeys,
										ForceNew:     true,
										Description:  `A flag to enable confidential compute mode on boot disk`,
									},

									"storage_pool": {
										Type:             schema.TypeString,
										Optional:         true,
										AtLeastOneOf:     initializeParamsKeys,
										ForceNew:         true,
										DiffSuppressFunc: tpgresource.CompareResourceNames,
										Description:      `The URL of the storage pool in which the new disk is created`,
									},

									"architecture": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ForceNew:     true,
										AtLeastOneOf: initializeParamsKeys,
										ValidateFunc: validation.StringInSlice([]string{"X86_64", "ARM64"}, false),
										Description:  `The architecture of the disk. One of "X86_64" or "ARM64".`,
									},
								},
							},
						},

						"mode": {
							Type:         schema.TypeString,
							Optional:     true,
							AtLeastOneOf: bootDiskKeys,
							ForceNew:     true,
							Default:      "READ_WRITE",
							ValidateFunc: validation.StringInSlice([]string{"READ_WRITE", "READ_ONLY"}, false),
							Description:  `Read/write mode for the disk. One of "READ_ONLY" or "READ_WRITE".`,
						},

						"source": {
							Type:             schema.TypeString,
							Optional:         true,
							AtLeastOneOf:     bootDiskKeys,
							Computed:         true,
							ForceNew:         true,
							ConflictsWith:    []string{"boot_disk.initialize_params"},
							DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
							Description:      `The name or self_link of the disk attached to this instance.`,
						},
					},
				},
			},

			"machine_type": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      `The machine type to create.`,
				DiffSuppressFunc: tpgresource.CompareResourceNames,
			},

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidateRFC1035Name(1, 63),
				Description:  `The name of the instance. One of name or self_link must be provided.`,
			},

			"network_interface": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: `The networks attached to the instance.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
							Description:      `The name or self_link of the network attached to this interface.`,
						},

						"subnetwork": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
							Description:      `The name or self_link of the subnetwork attached to this interface.`,
						},

						"network_attachment": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ForceNew:         true,
							DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
							Description:      `The URL of the network attachment that this interface should connect to in the following format: projects/{projectNumber}/regions/{region_name}/networkAttachments/{network_attachment_name}.`,
						},

						"subnetwork_project": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The project in which the subnetwork belongs.`,
						},

						"network_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The private IP address assigned to the instance.`,
						},

						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the interface`,
						},
						"nic_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"GVNIC", "VIRTIO_NET", "IDPF", "MRDMA", "IRDMA"}, false),
							Description:  `The type of vNIC to be used on this interface. Possible values:GVNIC, VIRTIO_NET, IDPF, MRDMA, and IRDMA`,
						},
						"access_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `Access configurations, i.e. IPs via which this instance can be accessed via the Internet.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"nat_ip": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: `The IP address that is be 1:1 mapped to the instance's network ip.`,
									},

									"network_tier": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: `The networking tier used for configuring this instance. One of PREMIUM or STANDARD.`,
									},

									"public_ptr_domain_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The DNS domain name for the public PTR record.`,
									},
									"security_policy": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `A full or partial URL to a security policy to add to this instance. If this field is set to an empty string it will remove the associated security policy.`,
									},
								},
							},
						},

						"alias_ip_range": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `An array of alias IP ranges for this network interface.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_cidr_range": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The IP CIDR range represented by this alias IP range.`,
									},
									"subnetwork_range_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The subnetwork secondary range name specifying the secondary range from which to allocate the IP CIDR range for this alias IP range.`,
									},
								},
							},
						},

						"stack_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"IPV4_ONLY", "IPV4_IPV6", "IPV6_ONLY", ""}, false),
							Description:  `The stack type for this network interface to identify whether the IPv6 feature is enabled or not. If not specified, IPV4_ONLY will be used.`,
						},

						"ipv6_access_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `One of EXTERNAL, INTERNAL to indicate whether the IP can be accessed from the Internet. This field is always inherited from its subnetwork.`,
						},

						"ipv6_access_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `An array of IPv6 access configurations for this interface. Currently, only one IPv6 access config, DIRECT_IPV6, is supported. If there is no ipv6AccessConfig specified, then this instance will have no external IPv6 Internet access.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"network_tier": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The service-level to be provided for IPv6 traffic when the subnet has an external subnet. Only PREMIUM tier is valid for IPv6`,
									},
									"public_ptr_domain_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The domain name to be used when creating DNSv6 records for the external IPv6 ranges.`,
									},
									"external_ipv6": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
										Description: `The first IPv6 address of the external IPv6 range associated with this instance, prefix length is stored in externalIpv6PrefixLength in ipv6AccessConfig. To use a static external IP address, it must be unused and in the same region as the instance's zone. If not specified, Google Cloud will automatically assign an external IPv6 address from the instance's subnetwork.`,
									},
									"external_ipv6_prefix_length": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
										Description: `The prefix length of the external IPv6 range.`,
									},
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
										Description: `The name of this access configuration. In ipv6AccessConfigs, the recommended name is External IPv6.`,
									},
									"security_policy": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `A full or partial URL to a security policy to add to this instance. If this field is set to an empty string it will remove the associated security policy.`,
									},
								},
							},
						},

						"internal_ipv6_prefix_length": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: `The prefix length of the primary internal IPv6 range.`,
						},

						"ipv6_address": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `An IPv6 internal network address for this network interface. If not specified, Google Cloud will automatically assign an internal IPv6 address from the instance's subnetwork.`,
						},

						"queue_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: `The networking queue count that's specified by users for the network interface. Both Rx and Tx queues will be set to this number. It will be empty if not specified.`,
						},

						"security_policy": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `A full or partial URL to a security policy to add to this instance. If this field is set to an empty string it will remove the associated security policy.`,
						},
					},
				},
			},
			"network_performance_config": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: `Configures network performance settings for the instance. If not specified, the instance will be created with its default network performance configuration.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_egress_bandwidth_tier": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"TIER_1", "DEFAULT"}, false),
							Description:  `The egress bandwidth tier to enable. Possible values:TIER_1, DEFAULT`,
						},
					},
				},
			},
			"allow_stopping_for_update": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `If true, allows Terraform to stop the instance to update its properties. If you try to update a property that requires stopping the instance without setting this field, the update will fail.`,
			},

			"attached_disk": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `List of disks attached to the instance`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
							Description:      `The name or self_link of the disk attached to this instance.`,
						},

						"device_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `Name with which the attached disk is accessible under /dev/disk/by-id/`,
						},

						"mode": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "READ_WRITE",
							ValidateFunc: validation.StringInSlice([]string{"READ_WRITE", "READ_ONLY"}, false),
							Description:  `Read/write mode for the disk. One of "READ_ONLY" or "READ_WRITE".`,
						},

						"disk_encryption_key_raw": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: `A 256-bit customer-supplied encryption key, encoded in RFC 4648 base64 to encrypt this disk. Only one of kms_key_self_link, disk_encryption_key_rsa and disk_encryption_key_raw may be set.`,
						},

						"disk_encryption_key_rsa": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: `Specifies an RFC 4648 base64 encoded, RSA-wrapped 2048-bit customer-supplied encryption key to either encrypt or decrypt this resource. Only one of kms_key_self_link, disk_encryption_key_rsa and disk_encryption_key_raw may be set.`,
						},

						"kms_key_self_link": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: tpgresource.CompareSelfLinkRelativePaths,
							Computed:         true,
							Description:      `The self_link of the encryption key that is stored in Google Cloud KMS to encrypt this disk. Only one of kms_key_self_link, disk_encryption_key_rsa and disk_encryption_key_raw may be set.`,
						},

						"disk_encryption_service_account": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The service account being used for the encryption request for the given KMS key. If absent, the Compute Engine default service account is used`,
						},

						"disk_encryption_key_sha256": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The RFC 4648 base64 encoded SHA-256 hash of the customer-supplied encryption key that protects this resource.`,
						},
					},
				},
			},

			"can_ip_forward": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: `Whether sending and receiving of packets with non-matching source or destination IPs is allowed.`,
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `A brief description of the resource.`,
			},

			"deletion_protection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: `Whether deletion protection is enabled on this instance.`,
			},

			"enable_display": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether the instance has virtual displays enabled.`,
			},

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
							Description: `The number of the guest accelerator cards exposed to this instance.`,
						},
						"type": {
							Type:             schema.TypeString,
							Required:         true,
							ForceNew:         true,
							DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
							Description:      `The accelerator type resource exposed to this instance. E.g. nvidia-tesla-k80.`,
						},
					},
				},
			},

			"params": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: `Stores additional params passed with the request, but not persisted as part of resource payload.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_manager_tags": {
							Type:        schema.TypeMap,
							Optional:    true,
							Description: `A map of resource manager tags. Resource manager tag keys and values have the same definition as resource manager tags. Keys must be in the format tagKeys/{tag_key_id}, and values are in the format tagValues/456. The field is ignored (both PUT & PATCH) when empty.`,
						},
					},
				},
			},

			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `A set of key/value label pairs assigned to the instance.

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

			"metadata": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Metadata key/value pairs made available within the instance.`,
			},

			"partner_metadata": {
				Type:                  schema.TypeMap,
				Optional:              true,
				DiffSuppressFunc:      ComparePartnerMetadataDiff,
				DiffSuppressOnRefresh: true,
				Elem:                  &schema.Schema{Type: schema.TypeString},
				Description:           `Partner Metadata Map made available within the instance.`,
			},

			"metadata_startup_script": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Metadata startup scripts made available within the instance.`,
			},

			"min_cpu_platform": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The minimum CPU platform specified for the VM instance.`,
			},

			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The ID of the project in which the resource belongs. If self_link is provided, this value is ignored. If neither self_link nor project are provided, the provider project is used.`,
			},

			"scheduling": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: `The scheduling strategy being used by the instance.`,
				Elem: &schema.Resource{
					// !!! IMPORTANT !!!
					// We have a custom diff function for the scheduling block due to issues with Terraform's
					// diff on schema.Set. If changes are made to this block, they must be reflected in that
					// method. See schedulingHasChangeWithoutReboot in compute_instance_helpers.go
					Schema: map[string]*schema.Schema{
						"on_host_maintenance": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: schedulingKeys,
							Description:  `Describes maintenance behavior for the instance. One of MIGRATE or TERMINATE,`,
						},

						"automatic_restart": {
							Type:         schema.TypeBool,
							Optional:     true,
							AtLeastOneOf: schedulingKeys,
							Default:      true,
							Description:  `Specifies if the instance should be restarted if it was terminated by Compute Engine (not a user).`,
						},

						"preemptible": {
							Type:         schema.TypeBool,
							Optional:     true,
							Default:      false,
							AtLeastOneOf: schedulingKeys,
							ForceNew:     true,
							Description:  `Whether the instance is preemptible.`,
						},

						"node_affinities": {
							Type:             schema.TypeSet,
							Optional:         true,
							AtLeastOneOf:     schedulingKeys,
							Elem:             instanceSchedulingNodeAffinitiesElemSchema(),
							DiffSuppressFunc: tpgresource.EmptyOrDefaultStringSuppress(""),
							Description:      `Specifies node affinities or anti-affinities to determine which sole-tenant nodes your instances and managed instance groups will use as host systems.`,
						},

						"min_node_cpus": {
							Type:         schema.TypeInt,
							Optional:     true,
							AtLeastOneOf: schedulingKeys,
						},

						"provisioning_model": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							AtLeastOneOf: schedulingKeys,
							Description:  `Whether the instance is spot. If this is set as SPOT.`,
						},

						"instance_termination_action": {
							Type:         schema.TypeString,
							Optional:     true,
							AtLeastOneOf: schedulingKeys,
							Description:  `Specifies the action GCE should take when SPOT VM is preempted.`,
						},
						"termination_time": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							AtLeastOneOf: schedulingKeys,
							Description: `Specifies the timestamp, when the instance will be terminated,
in RFC3339 text format. If specified, the instance termination action
will be performed at the termination time.`,
						},
						"availability_domain": {
							Type:         schema.TypeInt,
							Optional:     true,
							AtLeastOneOf: schedulingKeys,
							Description:  `Specifies the availability domain, which this instance should be scheduled on.`,
						},
						"max_run_duration": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `The timeout for new network connections to hosts.`,
							MaxItems:    1,
							ForceNew:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"seconds": {
										Type:     schema.TypeInt,
										Required: true,
										ForceNew: true,
										Description: `Span of time at a resolution of a second.
Must be from 0 to 315,576,000,000 inclusive.`,
									},
									"nanos": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
										Description: `Span of time that's a fraction of a second at nanosecond
resolution. Durations less than one second are represented
with a 0 seconds field and a positive nanos field. Must
be from 0 to 999,999,999 inclusive.`,
									},
								},
							},
						},
						"on_instance_stop_action": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							ForceNew:    true,
							Description: `Defines the behaviour for instances with the instance_termination_action.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"discard_local_ssd": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: `If true, the contents of any attached Local SSD disks will be discarded.`,
										Default:     false,
										ForceNew:    true,
									},
								},
							},
						},
						"host_error_timeout_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: `Specify the time in seconds for host error detection, the value must be within the range of [90, 330] with the increment of 30, if unset, the default behavior of host error recovery will be used.`,
						},

						"maintenance_interval": {
							Type:         schema.TypeString,
							Optional:     true,
							AtLeastOneOf: schedulingKeys,
							Description:  `Specifies the frequency of planned maintenance events. The accepted values are: PERIODIC`,
						},
						"local_ssd_recovery_timeout": {
							Type:     schema.TypeList,
							Optional: true,
							Description: `Specifies the maximum amount of time a Local Ssd Vm should wait while
  recovery of the Local Ssd state is attempted. Its value should be in
  between 0 and 168 hours with hour granularity and the default value being 1
  hour.`,
							MaxItems: 1,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"seconds": {
										Type:     schema.TypeInt,
										Required: true,
										ForceNew: true,
										Description: `Span of time at a resolution of a second.
Must be from 0 to 315,576,000,000 inclusive.`,
									},
									"nanos": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
										Description: `Span of time that's a fraction of a second at nanosecond
resolution. Durations less than one second are represented
with a 0 seconds field and a positive nanos field. Must
be from 0 to 999,999,999 inclusive.`,
									},
								},
							},
						},
						"graceful_shutdown": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `Settings for the instance to perform a graceful shutdown.`,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: `Opts-in for graceful shutdown.`,
									},
									"max_duration": {
										Type:     schema.TypeList,
										Optional: true,
										Description: `The time allotted for the instance to gracefully shut down.
										If the graceful shutdown isn't complete after this time, then the instance
										transitions to the STOPPING state.`,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"seconds": {
													Type:     schema.TypeInt,
													Required: true,
													Description: `Span of time at a resolution of a second.
													The value must be between 1 and 3600, which is 3,600 seconds (one hour).`,
												},
												"nanos": {
													Type:     schema.TypeInt,
													Optional: true,
													Description: `Span of time that's a fraction of a second at nanosecond
													resolution. Durations less than one second are represented
													with a 0 seconds field and a positive nanos field. Must
													be from 0 to 999,999,999 inclusive.`,
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

			"scratch_disk": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: `The scratch disks attached to the instance.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"device_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `Name with which the attached disk is accessible under /dev/disk/by-id/`,
						},
						"interface": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"SCSI", "NVME"}, false),
							Description:  `The disk interface used for attaching this disk. One of SCSI or NVME.`,
						},
						"size": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.IntAtLeast(375),
							Default:      375,
							Description:  `The size of the disk in gigabytes. One of 375 or 3000.`,
						},
					},
				},
			},

			"service_account": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: `The service account to attach to the instance.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The service account e-mail address.`,
						},

						"scopes": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: `A list of service scopes.`,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								StateFunc: func(v interface{}) string {
									return tpgresource.CanonicalizeServiceScope(v.(string))
								},
							},
							Set: tpgresource.StringScopeHashcode,
						},
					},
				},
			},

			"shielded_instance_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				// Since this block is used by the API based on which
				// image being used, the field needs to be marked as Computed.
				Computed:         true,
				DiffSuppressFunc: tpgresource.EmptyOrDefaultStringSuppress(""),
				Description:      `The shielded vm config being used by the instance.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_secure_boot": {
							Type:         schema.TypeBool,
							Optional:     true,
							AtLeastOneOf: shieldedInstanceConfigKeys,
							Default:      false,
							Description:  `Whether secure boot is enabled for the instance.`,
						},

						"enable_vtpm": {
							Type:         schema.TypeBool,
							Optional:     true,
							AtLeastOneOf: shieldedInstanceConfigKeys,
							Default:      true,
							Description:  `Whether the instance uses vTPM.`,
						},

						"enable_integrity_monitoring": {
							Type:         schema.TypeBool,
							Optional:     true,
							AtLeastOneOf: shieldedInstanceConfigKeys,
							Default:      true,
							Description:  `Whether integrity monitoring is enabled for the instance.`,
						},
					},
				},
			},
			"advanced_machine_features": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: `Controls for advanced machine-related behavior features.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_nested_virtualization": {
							Type:         schema.TypeBool,
							Optional:     true,
							AtLeastOneOf: advancedMachineFeaturesKeys,
							Description:  `Whether to enable nested virtualization or not.`,
						},
						"threads_per_core": {
							Type:         schema.TypeInt,
							Optional:     true,
							AtLeastOneOf: advancedMachineFeaturesKeys,
							Description:  `The number of threads per physical core. To disable simultaneous multithreading (SMT) set this to 1. If unset, the maximum number of threads supported per core by the underlying processor is assumed.`,
						},
						"turbo_mode": {
							Type:         schema.TypeString,
							Optional:     true,
							AtLeastOneOf: advancedMachineFeaturesKeys,
							Description:  `Turbo frequency mode to use for the instance. Currently supported modes is "ALL_CORE_MAX".`,
							ValidateFunc: validation.StringInSlice([]string{"ALL_CORE_MAX"}, false),
						},
						"visible_core_count": {
							Type:         schema.TypeInt,
							Optional:     true,
							AtLeastOneOf: advancedMachineFeaturesKeys,
							Description:  `The number of physical cores to expose to an instance. Multiply by the number of threads per core to compute the total number of virtual CPUs to expose to the instance. If unset, the number of cores is inferred from the instance\'s nominal CPU count and the underlying platform\'s SMT width.`,
						},
						"performance_monitoring_unit": {
							Type:         schema.TypeString,
							Optional:     true,
							AtLeastOneOf: advancedMachineFeaturesKeys,
							ValidateFunc: validation.StringInSlice([]string{"STANDARD", "ENHANCED", "ARCHITECTURAL"}, false),
							Description:  `The PMU is a hardware component within the CPU core that monitors how the processor runs code. Valid values for the level of PMU are "STANDARD", "ENHANCED", and "ARCHITECTURAL".`,
						},
						"enable_uefi_networking": {
							Type:         schema.TypeBool,
							Optional:     true,
							ForceNew:     true,
							AtLeastOneOf: advancedMachineFeaturesKeys,
							Description:  `Whether to enable UEFI networking for the instance.`,
						},
					},
				},
			},
			"confidential_instance_config": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `The Confidential VM config being used by the instance.  on_host_maintenance has to be set to TERMINATE or this will fail to create.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_confidential_compute": {
							Type:         schema.TypeBool,
							Optional:     true,
							Description:  `Defines whether the instance should have confidential compute enabled. Field will be deprecated in a future release`,
							AtLeastOneOf: []string{"confidential_instance_config.0.enable_confidential_compute", "confidential_instance_config.0.confidential_instance_type"},
						},
						"confidential_instance_type": {
							Type:     schema.TypeString,
							Optional: true,
							Description: `
								The confidential computing technology the instance uses.
								SEV is an AMD feature. TDX is an Intel feature. One of the following
								values is required: SEV, SEV_SNP, TDX. If SEV_SNP, min_cpu_platform =
								"AMD Milan" is currently required.`,
							AtLeastOneOf: []string{"confidential_instance_config.0.enable_confidential_compute", "confidential_instance_config.0.confidential_instance_type"},
						},
					},
				},
			},
			"desired_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"RUNNING", "TERMINATED", "SUSPENDED"}, false),
				Description:  `Desired status of the instance. Either "RUNNING", "SUSPENDED" or "TERMINATED".`,
			},
			"current_status": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `
					Current status of the instance.
					This could be one of the following values: PROVISIONING, STAGING, RUNNING, STOPPING, SUSPENDING, SUSPENDED, REPAIRING, and TERMINATED.
					For more information about the status of the instance, see [Instance life cycle](https://cloud.google.com/compute/docs/instances/instance-life-cycle).`,
			},
			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: `The list of tags attached to the instance.`,
			},

			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The zone of the instance. If self_link is provided, this value is ignored. If neither self_link nor zone are provided, the provider zone is used.`,
			},

			"cpu_platform": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The CPU platform used by this instance.`,
			},

			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The server-assigned unique identifier of this instance.`,
			},

			"creation_timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Creation timestamp in RFC3339 text format.`,
			},

			"label_fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The unique fingerprint of the labels.`,
			},

			"metadata_fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The unique fingerprint of the metadata.`,
			},

			"self_link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The URI of the created resource.`,
			},

			"tags_fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The unique fingerprint of the tags.`,
			},

			"hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `A custom hostname for the instance. Must be a fully qualified DNS name and RFC-1035-valid. Valid format is a series of labels 1-63 characters long matching the regular expression [a-z]([-a-z0-9]*[a-z0-9]), concatenated with periods. The entire hostname must not exceed 253 characters. Changing this forces a new resource to be created.`,
			},

			"resource_policies": {
				Type:             schema.TypeList,
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: tpgresource.CompareSelfLinkRelativePaths,
				Optional:         true,
				MaxItems:         1,
				Description:      `A list of self_links of resource policies to attach to the instance. Currently a max of 1 resource policy is supported.`,
			},

			"reservation_affinity": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the reservations that this instance can consume from.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"ANY_RESERVATION", "SPECIFIC_RESERVATION", "NO_RESERVATION"}, false),
							Description:  `The type of reservation from which this instance can consume resources.`,
						},

						"specific_reservation": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							ForceNew:    true,
							Description: `Specifies the label selector for the reservation to use.`,

							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: `Corresponds to the label key of a reservation resource. To target a SPECIFIC_RESERVATION by name, specify compute.googleapis.com/reservation-name as the key and specify the name of your reservation as the only value.`,
									},
									"values": {
										Type:        schema.TypeList,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Required:    true,
										ForceNew:    true,
										Description: `Corresponds to the label values of a reservation resource.`,
									},
								},
							},
						},
					},
				},
			},

			"key_revocation_action_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"STOP", "NONE", ""}, false),
				Description:  `Action to be taken when a customer's encryption key is revoked. Supports "STOP" and "NONE", with "NONE" being the default.`,
			},

			"instance_encryption_key": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: `Encryption key used to provide data encryption on the given instance.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kms_key_self_link": {
							Type:             schema.TypeString,
							Optional:         true,
							ForceNew:         true,
							DiffSuppressFunc: tpgresource.CompareCryptoKeyVersions,
							Computed:         true,
							Description:      `The self link of the encryption key that is stored in Google Cloud KMS.`,
						},

						"kms_key_service_account": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: `The service account being used for the encryption request for the given KMS key. If absent, the Compute Engine default service account is used.`,
						},

						"sha256": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The SHA256 hash of the customer's encryption key.`,
						},
					},
				},
			},
		},
		UseJSONNumber: true,
	}
}

func flattenAliasIpRangeTgc(ranges []*compute.AliasIpRange) []map[string]interface{} {
	rangesSchema := make([]map[string]interface{}, 0, len(ranges))
	for _, ipRange := range ranges {
		rangesSchema = append(rangesSchema, map[string]interface{}{
			"ip_cidr_range":         ipRange.IpCidrRange,
			"subnetwork_range_name": ipRange.SubnetworkRangeName,
		})
	}
	return rangesSchema
}

func flattenSchedulingTgc(resp *compute.Scheduling) []map[string]interface{} {
	schedulingMap := make(map[string]interface{}, 0)

	// gracefulShutdown is not in the cai asset, so graceful_shutdown is skipped.

	if resp.InstanceTerminationAction != "" {
		schedulingMap["instance_termination_action"] = resp.InstanceTerminationAction
	}

	if resp.MinNodeCpus != 0 {
		schedulingMap["min_node_cpus"] = resp.MinNodeCpus
	}

	schedulingMap["on_host_maintenance"] = resp.OnHostMaintenance

	if resp.AutomaticRestart != nil && !*resp.AutomaticRestart {
		schedulingMap["automatic_restart"] = *resp.AutomaticRestart
	}

	if resp.Preemptible {
		schedulingMap["preemptible"] = resp.Preemptible
	}

	if resp.NodeAffinities != nil && len(resp.NodeAffinities) > 0 {
		nodeAffinities := []map[string]interface{}{}
		for _, na := range resp.NodeAffinities {
			nodeAffinities = append(nodeAffinities, map[string]interface{}{
				"key":      na.Key,
				"operator": na.Operator,
				"values":   tpgresource.ConvertStringArrToInterface(na.Values),
			})
		}
		schedulingMap["node_affinities"] = nodeAffinities
	}

	schedulingMap["provisioning_model"] = resp.ProvisioningModel

	if resp.AvailabilityDomain != 0 {
		schedulingMap["availability_domain"] = resp.AvailabilityDomain
	}

	if resp.MaxRunDuration != nil {
		schedulingMap["max_run_duration"] = flattenComputeMaxRunDuration(resp.MaxRunDuration)
	}

	if resp.OnInstanceStopAction != nil {
		schedulingMap["on_instance_stop_action"] = flattenOnInstanceStopAction(resp.OnInstanceStopAction)
	}

	if resp.HostErrorTimeoutSeconds != 0 {
		schedulingMap["host_error_timeout_seconds"] = resp.HostErrorTimeoutSeconds
	}

	if resp.MaintenanceInterval != "" {
		schedulingMap["maintenance_interval"] = resp.MaintenanceInterval
	}

	if resp.LocalSsdRecoveryTimeout != nil {
		schedulingMap["local_ssd_recovery_timeout"] = flattenComputeLocalSsdRecoveryTimeout(resp.LocalSsdRecoveryTimeout)
	}

	if len(schedulingMap) == 0 {
		return nil
	}

	return []map[string]interface{}{schedulingMap}
}

func flattenNetworkInterfacesTgc(networkInterfaces []*compute.NetworkInterface, project string) ([]map[string]interface{}, string, string, error) {
	flattened := make([]map[string]interface{}, len(networkInterfaces))
	var internalIP, externalIP string

	for i, iface := range networkInterfaces {
		var ac []map[string]interface{}
		ac, externalIP = flattenAccessConfigs(iface.AccessConfigs)

		flattened[i] = map[string]interface{}{
			"network_ip":                  iface.NetworkIP,
			"access_config":               ac,
			"alias_ip_range":              flattenAliasIpRangeTgc(iface.AliasIpRanges),
			"nic_type":                    iface.NicType,
			"stack_type":                  iface.StackType,
			"ipv6_access_config":          flattenIpv6AccessConfigs(iface.Ipv6AccessConfigs),
			"ipv6_address":                iface.Ipv6Address,
			"network":                     tpgresource.ConvertSelfLinkToV1(iface.Network),
			"subnetwork":                  tpgresource.ConvertSelfLinkToV1(iface.Subnetwork),
			"internal_ipv6_prefix_length": iface.InternalIpv6PrefixLength,
		}

		subnetProject := utils.ParseFieldValue(iface.Subnetwork, "projects")
		if subnetProject != project {
			flattened[i]["subnetwork_project"] = subnetProject
		}

		// The field name is computed, no it is not converted.

		if iface.StackType != "IPV4_ONLY" {
			flattened[i]["stack_type"] = iface.StackType
		}

		if iface.QueueCount != 0 {
			flattened[i]["queue_count"] = iface.QueueCount
		}

		if internalIP == "" {
			internalIP = iface.NetworkIP
		}

		if iface.NetworkAttachment != "" {
			networkAttachment, err := tpgresource.GetRelativePath(iface.NetworkAttachment)
			if err != nil {
				return nil, "", "", err
			}
			flattened[i]["network_attachment"] = networkAttachment
		}

		// the security_policy for a network_interface is found in one of its accessConfigs.
		if len(iface.AccessConfigs) > 0 && iface.AccessConfigs[0].SecurityPolicy != "" {
			flattened[i]["security_policy"] = iface.AccessConfigs[0].SecurityPolicy
		} else if len(iface.Ipv6AccessConfigs) > 0 && iface.Ipv6AccessConfigs[0].SecurityPolicy != "" {
			flattened[i]["security_policy"] = iface.Ipv6AccessConfigs[0].SecurityPolicy
		}
	}
	return flattened, internalIP, externalIP, nil
}

func flattenServiceAccountsTgc(serviceAccounts []*compute.ServiceAccount) []map[string]interface{} {
	result := make([]map[string]interface{}, len(serviceAccounts))
	for i, serviceAccount := range serviceAccounts {
		scopes := serviceAccount.Scopes
		if len(scopes) == 0 {
			scopes = []string{}
		}
		result[i] = map[string]interface{}{
			"email":  serviceAccount.Email,
			"scopes": scopes,
		}
	}
	return result
}

func flattenGuestAcceleratorsTgc(accelerators []*compute.AcceleratorConfig) []map[string]interface{} {
	acceleratorsSchema := make([]map[string]interface{}, len(accelerators))
	for i, accelerator := range accelerators {
		acceleratorsSchema[i] = map[string]interface{}{
			"count": accelerator.AcceleratorCount,
			"type":  tpgresource.GetResourceNameFromSelfLink(accelerator.AcceleratorType),
		}
	}
	return acceleratorsSchema
}

func flattenReservationAffinityTgc(affinity *compute.ReservationAffinity) []map[string]interface{} {
	if affinity == nil {
		return nil
	}

	// The values of ConsumeReservationType in cai assets are NO_ALLOCATION, SPECIFIC_ALLOCATION, ANY_ALLOCATION
	crt := strings.ReplaceAll(affinity.ConsumeReservationType, "_ALLOCATION", "_RESERVATION")
	flattened := map[string]interface{}{
		"type": crt,
	}

	if crt == "SPECIFIC_RESERVATION" {
		flattened["specific_reservation"] = []map[string]interface{}{{
			"key":    affinity.Key,
			"values": affinity.Values,
		}}
	}

	return []map[string]interface{}{flattened}
}
