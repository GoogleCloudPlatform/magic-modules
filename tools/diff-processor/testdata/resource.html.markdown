---
subcategory: "Compute Engine"
description: |-
  Manages abcdefg.
---

# google_test_resource

This resource combines some sections in google_compute_instance, google_container_attached_cluster, network_services_endpoint_policy and irrelvant parts are trimmed. 

## Example Usage

Lorem ipsum

## Example usage - Confidential Computing

Lorem ipsum


## Argument Reference

The following arguments are supported:

* `boot_disk` - (Required) The boot disk for the instance.
    Structure is [documented below](#nested_boot_disk).

* `name` - (Required) A unique name for the resource, required by GCE.
    Changing this forces a new resource to be created.

* `zone` - (Optional) The zone that the machine should be created in. If it is not provided, the provider zone is used.

* `network_interface` - (Required) Networks to attach to the instance. This can
    be specified multiple times. Structure is [documented below](#nested_network_interface).

* `params` - (Optional) Additional instance parameters.

---

<a name="nested_boot_disk"></a>The `boot_disk` block supports:

* `auto_delete` - (Optional) Whether the disk will be auto-deleted when the instance
    is deleted. Defaults to true.

* `device_name` - (Optional) Name with which attached disk will be accessible.
    On the instance, this device will be `/dev/disk/by-id/google-{{device_name}}`.

* `mode` - (Optional) The mode in which to attach this disk, either `READ_WRITE`
  or `READ_ONLY`. If not specified, the default is to attach the disk in `READ_WRITE` mode.

* `disk_encryption_key_raw` - (Optional) A 256-bit [customer-supplied encryption key]
    (https://cloud.google.com/compute/docs/disks/customer-supplied-encryption),
    encoded in [RFC 4648 base64](https://tools.ietf.org/html/rfc4648#section-4)
    to encrypt this disk. Only one of `kms_key_self_link` and `disk_encryption_key_raw`
    may be set.

* `kms_key_self_link` - (Optional) The self_link of the encryption key that is
    stored in Google Cloud KMS to encrypt this disk. Only one of `kms_key_self_link`
    and `disk_encryption_key_raw` may be set.

* `initialize_params` - (Optional) Parameters for a new disk that will be created
    alongside the new instance. Either `initialize_params` or `source` must be set.
    Structure is [documented below](#nested_initialize_params).

* `source` - (Optional) The name or self_link of the existing disk (such as those managed by
    `google_compute_disk`) or disk image. To create an instance from a snapshot, first create a
    `google_compute_disk` from a snapshot and reference it here.

<a name="nested_initialize_params"></a>The `initialize_params` block supports:

* `size` - (Optional) The size of the image in gigabytes. If not specified, it
    will inherit the size of its base image.

* `type` - (Optional) The GCE disk type. Such as pd-standard, pd-balanced or pd-ssd.

* `image` - (Optional) The image from which to initialize this disk. This can be
    one of: the image's `self_link`, `projects/{project}/global/images/{image}`,
    `projects/{project}/global/images/family/{family}`, `global/images/{image}`,
    `global/images/family/{family}`, `family/{family}`, `{project}/{family}`,
    `{project}/{image}`, `{family}`, or `{image}`. If referred by family, the
    images names must include the family name. If they don't, use the
    [google_compute_image data source](/docs/providers/google/d/compute_image.html).
    For instance, the image `centos-6-v20180104` includes its family name `centos-6`.
    These images can be referred by family name here.

* `labels` - (Optional) A set of key/value label pairs assigned to the disk. This
    field is only applicable for persistent disks.

* `resource_manager_tags` - (Optional) A tag is a key-value pair that can be attached to a Google Cloud resource. You can use tags to conditionally allow or deny policies based on whether a resource has a specific tag. This value is not returned by the API. In Terraform, this value cannot be updated and changing it will recreate the resource.

* `provisioned_iops` - (Optional) Indicates how many IOPS to provision for the disk.
    This sets the number of I/O operations per second that the disk can handle.
    For more details,see the [Hyperdisk documentation](https://cloud.google.com/compute/docs/disks/hyperdisks).
    Note: Updating currently is only supported for hyperdisk skus via disk update
    api/gcloud without the need to delete and recreate the disk, hyperdisk allows
    for an update of IOPS every 4 hours. To update your hyperdisk more frequently,
    you'll need to manually delete and recreate it.

* `provisioned_throughput` - (Optional) Indicates how much throughput to provision for the disk.
    This sets the number of throughput mb per second that the disk can handle.
    For more details,see the [Hyperdisk documentation](https://cloud.google.com/compute/docs/disks/hyperdisks).
    Note: Updating currently is only supported for hyperdisk skus via disk update
    api/gcloud without the need to delete and recreate the disk, hyperdisk allows
    for an update of throughput every 4 hours. To update your hyperdisk more
    frequently, you'll need to manually delete and recreate it.

* `enable_confidential_compute` - (Optional) Whether this disk is using confidential compute mode.
    Note: Only supported on hyperdisk skus, disk_encryption_key is required when setting to true.

* `storage_pool` - (Optional) The URL of the storage pool in which the new disk is created.
    For example:
    * https://www.googleapis.com/compute/v1/projects/{project}/zones/{zone}/storagePools/{storagePool}
    * /projects/{project}/zones/{zone}/storagePools/{storagePool}


<a name="nested_network_interface"></a>The `network_interface` block supports:

* `network` - (Optional) The name or self_link of the network to attach this interface to.
    Either `network` or `subnetwork` must be provided. If network isn't provided it will
    be inferred from the subnetwork.

*  `subnetwork` - (Optional) The name or self_link of the subnetwork to attach this
    interface to. Either `network` or `subnetwork` must be provided. If network isn't provided
    it will be inferred from the subnetwork. The subnetwork must exist in the same region this
    instance will be created in. If the network resource is in
    [legacy](https://cloud.google.com/vpc/docs/legacy) mode, do not specify this field. If the
    network is in auto subnet mode, specifying the subnetwork is optional. If the network is
    in custom subnet mode, specifying the subnetwork is required.


*  `subnetwork_project` - (Optional) The project in which the subnetwork belongs.
   If the `subnetwork` is a self_link, this field is ignored in favor of the project
   defined in the subnetwork self_link. If the `subnetwork` is a name and this
   field is not provided, the provider project is used.

* `network_ip` - (Optional) The private IP address to assign to the instance. If
    empty, the address will be automatically assigned.

* `access_config` - (Optional) Access configurations, i.e. IPs via which this
    instance can be accessed via the Internet. Omit to ensure that the instance
    is not accessible from the Internet. If omitted, ssh provisioners will not
    work unless Terraform can send traffic to the instance's network (e.g. via
    tunnel or because it is running on another cloud instance on that network).
    This block can be repeated multiple times. Structure [documented below](#nested_access_config).

* `alias_ip_range` - (Optional) An
    array of alias IP ranges for this network interface. Can only be specified for network
    interfaces on subnet-mode networks. Structure [documented below](#nested_alias_ip_range).

* `nic_type` - (Optional) The type of vNIC to be used on this interface. Possible values: GVNIC, VIRTIO_NET.

* `network_attachment` - (Optional) [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html) The URL of the network attachment that this interface should connect to in the following format: `projects/{projectNumber}/regions/{region_name}/networkAttachments/{network_attachment_name}`.

* `stack_type` - (Optional) The stack type for this network interface to identify whether the IPv6 feature is enabled or not. Values are IPV4_IPV6 or IPV4_ONLY. If not specified, IPV4_ONLY will be used.

* `ipv6_access_config` - (Optional) An array of IPv6 access configurations for this interface.
Currently, only one IPv6 access config, DIRECT_IPV6, is supported. If there is no ipv6AccessConfig
specified, then this instance will have no external IPv6 Internet access. Structure [documented below](#nested_ipv6_access_config).

* `queue_count` - (Optional) The networking queue count that's specified by users for the network interface. Both Rx and Tx queues will be set to this number. It will be empty if not specified.

* `security_policy` - (Optional) [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html) A full or partial URL to a security policy to add to this instance. If this field is set to an empty string it will remove the associated security policy.

<a name="nested_access_config"></a>The `access_config` block supports:

* `nat_ip` - (Optional) The IP address that will be 1:1 mapped to the instance's
    network ip. If not given, one will be generated.

* `public_ptr_domain_name` - (Optional) The DNS domain name for the public PTR record.
    To set this field on an instance, you must be verified as the owner of the domain.
    See [the docs](https://cloud.google.com/compute/docs/instances/create-ptr-record) for how
    to become verified as a domain owner.

* `network_tier` - (Optional) The [networking tier](https://cloud.google.com/network-tiers/docs/overview) used for configuring this instance.
    This field can take the following values: PREMIUM, FIXED_STANDARD or STANDARD. If this field is
    not specified, it is assumed to be PREMIUM.

<a name="nested_ipv6_access_config"></a>The `ipv6_access_config` block supports:

* `external_ipv6` - (Optional) The first IPv6 address of the external IPv6 range associated
    with this instance, prefix length is stored in externalIpv6PrefixLength in ipv6AccessConfig.
    To use a static external IP address, it must be unused and in the same region as the instance's zone.
    If not specified, Google Cloud will automatically assign an external IPv6 address from the instance's subnetwork.

* `external_ipv6_prefix_length` - (Optional) The prefix length of the external IPv6 range.

* `name` - (Optional) The name of this access configuration. In ipv6AccessConfigs, the recommended name
    is "External IPv6".

* `network_tier` - (Optional) The service-level to be provided for IPv6 traffic when the
    subnet has an external subnet. Only PREMIUM or STANDARD tier is valid for IPv6.

* `public_ptr_domain_name` - (Optional) The domain name to be used when creating DNSv6
    records for the external IPv6 ranges..

<a name="nested_alias_ip_range"></a>The `alias_ip_range` block supports:

* `ip_cidr_range` - The IP CIDR range represented by this alias IP range. This IP CIDR range
    must belong to the specified subnetwork and cannot contain IP addresses reserved by
    system or used by other network interfaces. This range may be a single IP address
    (e.g. 10.2.3.4), a netmask (e.g. /24) or a CIDR format string (e.g. 10.1.2.0/24).

* `subnetwork_range_name` - (Optional) The subnetwork secondary range name specifying
    the secondary range from which to allocate the IP CIDR range for this alias IP
    range. If left unspecified, the primary range of the subnetwork will be used.

<a name="nested_params"></a>The `params` block supports:

* `resource_manager_tags` (Optional) - A tag is a key-value pair that can be attached to a Google Cloud resource. You can use tags to conditionally allow or deny policies based on whether a resource has a specific tag. This value is not returned by the API. In Terraform, this value cannot be updated and changing it will recreate the resource.

- - -


- `labels` -
  (Optional)
  Set of label tags associated with the TcpRoute resource.
  **Note**: This field is non-authoritative, and will only manage the labels present in your configuration.
  Please refer to the field `effective_labels` for all of the labels present on the resource.

* `description` -
  (Optional)
  A free-text description of the resource. Max length 1024 characters.

* `traffic_port_selector` -
  (Optional)
  Port selector for the (matched) endpoints. If no port selector is provided, the matched config is applied to all ports.
  Structure is [documented below](#nested_traffic_port_selector).

* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.


<a name="nested_traffic_port_selector"></a>The `traffic_port_selector` block supports:

* `ports` -
  (Required)
  List of ports. Can be port numbers or port range (example, [80-90] specifies all ports from 80 to 90, including 80 and 90) or named ports or * to specify all ports. If the list is empty, all ports are selected.


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `id` - an identifier for the resource with format `projects/{{project}}/zones/{{zone}}/instances/{{name}}`

* `network_interface.0.access_config.0.nat_ip` - If the instance has an access config, either the given external ip (in the `nat_ip` field) or the ephemeral (generated) ip (if you didn't provide one).

* `workload_identity_config` -
  Workload Identity settings.
  Structure is [documented below](#nested_workload_identity_config).

* `errors` -
  A set of errors found in the cluster.
  Structure is [documented below](#nested_errors).
 

<a name="nested_workload_identity_config"></a>The `workload_identity_config` block contains:

* `identity_provider` -
  (Optional)
  The ID of the OIDC Identity Provider (IdP) associated to
  the Workload Identity Pool.

* `issuer_uri` -
  (Optional)
  The OIDC issuer URL for this cluster.

* `workload_pool` -
  (Optional)
  The Workload Identity Pool associated to the cluster.

<a name="nested_errors"></a>The `errors` block contains:

* `message` -
  (Optional)
  Human-friendly description of the error.

## Ephemeral Attributes Reference

The following write-only attributes are supported:

* `shared_secret_wo` -
  (Optional)
  Shared secret used to set the secure session between the Cloud VPN
  gateway and the peer VPN gateway.
  Note: This property is write-only and will not be read from the API. For more info see [updating write-only attributes](/docs/providers/google/guides/using_write_only_attributes.html#updating-write-only-attributes)
  **Note**: This property is write-only and will not be read from the API.

* `sensitive_params` -
  (Optional)
  Different parameters are configured primarily using the the `params` field on this
  resource. This block contains the parameters which contain secrets or passwords so that they can be marked
  sensitive and hidden from plan output. The name of the field, eg: secret_access_key, will be the key
  in the `params` map in the api request.
  Credentials may not be specified in both locations and will cause an error. Changing from one location
  to a different credential configuration in the config will require an apply to update state.
  Structure is [documented below](#nested_sensitive_params).

<a name="nested_sensitive_params"></a>The `sensitive_params` block supports:

* `secret_access_key_wo` -
  (Optional)
  The Secret Access Key of the AWS account transferring data from.
  **Note**: This property is write-only and will not be read from the API.

## Timeouts

Lorem ipsum

## Import

Lorem ipsum

