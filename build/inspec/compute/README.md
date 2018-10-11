# Google Compute Engine Chef Cookbook

This cookbook provides the built-in types and services for Chef to manage
Google Cloud Compute resources, as native Chef types.

## Requirements

### Platforms

#### Supported Operating Systems

This cookbook was tested on the following operating systems:

* RedHat 6, 7
* CentOS 6, 7
* Debian 7, 8
* Ubuntu 12.04, 14.04, 16.04, 16.10
* SLES 11-sp4, 12-sp2
* openSUSE 13
* Windows Server 2008 R2, 2012 R2, 2012 R2 Core, 2016 R2, 2016 R2 Core

## Example

```ruby
gauth_credential 'mycred' do
  action :serviceaccount
  path ENV['CRED_PATH'] # e.g. '/path/to/my_account.json'
  scopes [
    'https://www.googleapis.com/auth/compute'
  ]
end

gcompute_disk 'instance-test-os-1' do
  action :create
  source_image 'projects/ubuntu-os-cloud/global/images/family/ubuntu-1604-lts'
  zone 'us-west1-a'
  project ENV['PROJECT'] # ex: 'my-test-project'
  credential 'mycred'
end

gcompute_network 'mynetwork-test' do
  action :create
  project ENV['PROJECT'] # ex: 'my-test-project'
  credential 'mycred'
end

gcompute_address 'instance-test-ip' do
  action :create
  region 'us-west1'
  project ENV['PROJECT'] # ex: 'my-test-project'
  credential 'mycred'
end

gcompute_instance 'instance-test' do
  action :create
  machine_type 'n1-standard-1'
  disks [
    {
      boot: true,
      auto_delete: true,
      source: 'instance-test-os-1'
    }
  ]
  network_interfaces [
    {
      network: 'mynetwork-test',
      access_configs: [
        {
          name: 'External NAT',
          nat_ip: 'instance-test-ip',
          type: 'ONE_TO_ONE_NAT'
        }
      ]
    }
  ]
  zone 'us-west1-a'
  project ENV['PROJECT'] # ex: 'my-test-project'
  credential 'mycred'
end
```

## Credentials

All Google Cloud Platform cookbooks use an unified authentication mechanism,
provided by the `google-gauth` cookbook. Don't worry, it is automatically
installed when you install this module.

### Example

```ruby
gauth_credential 'mycred' do
  action :serviceaccount
  path ENV['CRED_PATH'] # e.g. '/path/to/my_account.json'
  scopes [
    'https://www.googleapis.com/auth/compute'
  ]
end

```

For complete details of the authentication cookbook, visit the
[google-gauth][] cookbook documentation.

## Resources

* [`gcompute_address`](#gcompute_address) -
    Represents an Address resource.
    Each virtual machine instance has an ephemeral internal IP address and,
    optionally, an external IP address. To communicate between instances on
    the same network, you can use an instance's internal IP address. To
    communicate with the Internet and instances outside of the same network,
    you must specify the instance's external IP address.
    Internal IP addresses are ephemeral and only belong to an instance for
    the lifetime of the instance; if the instance is deleted and recreated,
    the instance is assigned a new internal IP address, either by Compute
    Engine or by you. External IP addresses can be either ephemeral or
    static.
* [`gcompute_autoscaler`](#gcompute_autoscaler) -
    Represents an Autoscaler resource.
    Autoscalers allow you to automatically scale virtual machine instances in
    managed instance groups according to an autoscaling policy that you
    define.
* [`gcompute_backend_bucket`](#gcompute_backend_bucket) -
    Backend buckets allow you to use Google Cloud Storage buckets with HTTP(S)
    load balancing.
    An HTTP(S) load balancer can direct traffic to specified URLs to a
    backend bucket rather than a backend service. It can send requests for
    static content to a Cloud Storage bucket and requests for dynamic content
    a virtual machine instance.
* [`gcompute_backend_service`](#gcompute_backend_service) -
    Creates a BackendService resource in the specified project using the data
    included in the request.
* [`gcompute_disk`](#gcompute_disk) -
    Persistent disks are durable storage devices that function similarly to
    the physical disks in a desktop or a server. Compute Engine manages the
    hardware behind these devices to ensure data redundancy and optimize
    performance for you. Persistent disks are available as either standard
    hard disk drives (HDD) or solid-state drives (SSD).
    Persistent disks are located independently from your virtual machine
    instances, so you can detach or move persistent disks to keep your data
    even after you delete your instances. Persistent disk performance scales
    automatically with size, so you can resize your existing persistent disks
    or add more persistent disks to an instance to meet your performance and
    storage space requirements.
    Add a persistent disk to your instance when you need reliable and
    affordable storage with consistent performance characteristics.
* [`gcompute_firewall`](#gcompute_firewall) -
    Each network has its own firewall controlling access to and from the
    instances.
    All traffic to instances, even from other instances, is blocked by the
    firewall unless firewall rules are created to allow it.
    The default network has automatically created firewall rules that are
    shown in default firewall rules. No manually created network has
    automatically created firewall rules except for a default "allow" rule for
    outgoing traffic and a default "deny" for incoming traffic. For all
    networks except the default network, you must create any firewall rules
    you need.
* [`gcompute_forwarding_rule`](#gcompute_forwarding_rule) -
    A ForwardingRule resource. A ForwardingRule resource specifies which pool
    of target virtual machines to forward a packet to if it matches the given
    [IPAddress, IPProtocol, portRange] tuple.
* [`gcompute_global_address`](#gcompute_global_address) -
    Represents a Global Address resource. Global addresses are used for
    HTTP(S) load balancing.
* [`gcompute_global_forwarding_rule`](#gcompute_global_forwarding_rule) -
    Represents a GlobalForwardingRule resource. Global forwarding rules are
    used to forward traffic to the correct load balancer for HTTP load
    balancing. Global forwarding rules can only be used for HTTP load
    balancing.
    For more information, see
    https://cloud.google.com/compute/docs/load-balancing/http/
* [`gcompute_http_health_check`](#gcompute_http_health_check) -
    An HttpHealthCheck resource. This resource defines a template for how
    individual VMs should be checked for health, via HTTP.
* [`gcompute_https_health_check`](#gcompute_https_health_check) -
    An HttpsHealthCheck resource. This resource defines a template for how
    individual VMs should be checked for health, via HTTPS.
* [`gcompute_health_check`](#gcompute_health_check) -
    Health Checks determine whether instances are responsive and able to do work.
    They are an important part of a comprehensive load balancing configuration,
    as they enable monitoring instances behind load balancers.
    Health Checks poll instances at a specified interval. Instances that
    do not respond successfully to some number of probes in a row are marked
    as unhealthy. No new connections are sent to unhealthy instances,
    though existing connections will continue. The health check will
    continue to poll unhealthy instances. If an instance later responds
    successfully to some number of consecutive probes, it is marked
    healthy again and can receive new connections.
* [`gcompute_instance_template`](#gcompute_instance_template) -
    Defines an Instance Template resource that provides configuration settings
    for your virtual machine instances. Instance templates are not tied to the
    lifetime of an instance and can be used and reused as to deploy virtual
    machines. You can also use different templates to create different virtual
    machine configurations. Instance templates are required when you create a
    managed instance group.
    Tip: Disks should be set to autoDelete=true
    so that leftover disks are not left behind on machine deletion.
* [`gcompute_image`](#gcompute_image) -
    Represents an Image resource.
    Google Compute Engine uses operating system images to create the root
    persistent disks for your instances. You specify an image when you create
    an instance. Images contain a boot loader, an operating system, and a
    root file system. Linux operating system images are also capable of
    running containers on Compute Engine.
    Images can be either public or custom.
    Public images are provided and maintained by Google, open-source
    communities, and third-party vendors. By default, all projects have
    access to these images and can use them to create instances.  Custom
    images are available only to your project. You can create a custom image
    from root persistent disks and other images. Then, use the custom image
    to create an instance.
* [`gcompute_instance`](#gcompute_instance) -
    An instance is a virtual machine (VM) hosted on Google's infrastructure.
* [`gcompute_instance_group`](#gcompute_instance_group) -
    Represents an Instance Group resource. Instance groups are self-managed
    and can contain identical or different instances. Instance groups do not
    use an instance template. Unlike managed instance groups, you must create
    and add instances to an instance group manually.
* [`gcompute_instance_group_manager`](#gcompute_instance_group_manager) -
    Creates a managed instance group using the information that you specify in
    the request. After the group is created, it schedules an action to create
    instances in the group using the specified instance template. This
    operation is marked as DONE when the group is created even if the
    instances in the group have not yet been created. You must separately
    verify the status of the individual instances.
    A managed instance group can have up to 1000 VM instances per group.
* [`gcompute_network`](#gcompute_network) -
    Represents a Network resource.
    Your Cloud Platform Console project can contain multiple networks, and
    each network can have multiple instances attached to it. A network allows
    you to define a gateway IP and the network range for the instances
    attached to that network. Every project is provided with a default network
    with preset configurations and firewall rules. You can choose to customize
    the default network by adding or removing rules, or you can create new
    networks in that project. Generally, most users only need one network,
    although you can have up to five networks per project by default.
    A network belongs to only one project, and each instance can only belong
    to one network. All Compute Engine networks use the IPv4 protocol. Compute
    Engine currently does not support IPv6. However, Google is a major
    advocate of IPv6 and it is an important future direction.
* [`gcompute_region_autoscaler`](#gcompute_region_autoscaler) -
    Represents an Autoscaler resource.
    Autoscalers allow you to automatically scale virtual machine instances in
    managed instance groups according to an autoscaling policy that you
    define.
* [`gcompute_route`](#gcompute_route) -
    Represents a Route resource.
    A route is a rule that specifies how certain packets should be handled by
    the virtual network. Routes are associated with virtual machines by tag,
    and the set of routes for a particular virtual machine is called its
    routing table. For each packet leaving a virtual machine, the system
    searches that virtual machine's routing table for a single best matching
    route.
    Routes match packets by destination IP address, preferring smaller or more
    specific ranges over larger ones. If there is a tie, the system selects
    the route with the smallest priority value. If there is still a tie, it
    uses the layer three and four packet headers to select just one of the
    remaining matching routes. The packet is then forwarded as specified by
    the next_hop field of the winning route -- either to another virtual
    machine destination, a virtual machine gateway or a Compute
    Engine-operated gateway. Packets that do not match any route in the
    sending virtual machine's routing table will be dropped.
    A Route resource must have exactly one specification of either
    nextHopGateway, nextHopInstance, nextHopIp, or nextHopVpnTunnel.
* [`gcompute_router`](#gcompute_router) -
    Represents a Router resource.
* [`gcompute_snapshot`](#gcompute_snapshot) -
    Represents a Persistent Disk Snapshot resource.
    Use snapshots to back up data from your persistent disks. Snapshots are
    different from public images and custom images, which are used primarily
    to create instances or configure instance templates. Snapshots are useful
    for periodic backup of the data on your persistent disks. You can create
    snapshots from persistent disks even while they are attached to running
    instances.
    Snapshots are incremental, so you can create regular snapshots on a
    persistent disk faster and at a much lower cost than if you regularly
    created a full image of the disk.
* [`gcompute_ssl_certificate`](#gcompute_ssl_certificate) -
    An SslCertificate resource, used for HTTPS load balancing. This resource
    provides a mechanism to upload an SSL key and certificate to
    the load balancer to serve secure connections from the user.
* [`gcompute_ssl_policy`](#gcompute_ssl_policy) -
    Represents a SSL policy. SSL policies give you the ability to control the
    features of SSL that your SSL proxy or HTTPS load balancer negotiates.
* [`gcompute_subnetwork`](#gcompute_subnetwork) -
    A VPC network is a virtual version of the traditional physical networks
    that exist within and between physical data centers. A VPC network
    provides connectivity for your Compute Engine virtual machine (VM)
    instances, Container Engine containers, App Engine Flex services, and
    other network-related resources.
    Each GCP project contains one or more VPC networks. Each VPC network is a
    global entity spanning all GCP regions. This global VPC network allows VM
    instances and other resources to communicate with each other via internal,
    private IP addresses.
    Each VPC network is subdivided into subnets, and each subnet is contained
    within a single region. You can have more than one subnet in a region for
    a given VPC network. Each subnet has a contiguous private RFC1918 IP
    space. You create instances, containers, and the like in these subnets.
    When you create an instance, you must create it in a subnet, and the
    instance draws its internal IP address from that subnet.
    Virtual machine (VM) instances in a VPC network can communicate with
    instances in all other subnets of the same VPC network, regardless of
    region, using their RFC1918 private IP addresses. You can isolate portions
    of the network, even entire subnets, using firewall rules.
* [`gcompute_target_http_proxy`](#gcompute_target_http_proxy) -
    Represents a TargetHttpProxy resource, which is used by one or more global
    forwarding rule to route incoming HTTP requests to a URL map.
* [`gcompute_target_https_proxy`](#gcompute_target_https_proxy) -
    Represents a TargetHttpsProxy resource, which is used by one or more
    global forwarding rule to route incoming HTTPS requests to a URL map.
* [`gcompute_target_pool`](#gcompute_target_pool) -
    Represents a TargetPool resource, used for Load Balancing.
* [`gcompute_target_ssl_proxy`](#gcompute_target_ssl_proxy) -
    Represents a TargetSslProxy resource, which is used by one or more
    global forwarding rule to route incoming SSL requests to a backend
    service.
* [`gcompute_target_tcp_proxy`](#gcompute_target_tcp_proxy) -
    Represents a TargetTcpProxy resource, which is used by one or more
    global forwarding rule to route incoming TCP requests to a Backend
    service.
* [`gcompute_target_vpn_gateway`](#gcompute_target_vpn_gateway) -
    Represents a VPN gateway running in GCP. This virtual device is managed
    by Google, but used only by you.
* [`gcompute_url_map`](#gcompute_url_map) -
    UrlMaps are used to route requests to a backend service based on rules
    that you define for the host and path of an incoming URL.
* [`gcompute_vpn_tunnel`](#gcompute_vpn_tunnel) -
    VPN tunnel resource.


### gcompute_address
Represents an Address resource.

Each virtual machine instance has an ephemeral internal IP address and,
optionally, an external IP address. To communicate between instances on
the same network, you can use an instance's internal IP address. To
communicate with the Internet and instances outside of the same network,
you must specify the instance's external IP address.

Internal IP addresses are ephemeral and only belong to an instance for
the lifetime of the instance; if the instance is deleted and recreated,
the instance is assigned a new internal IP address, either by Compute
Engine or by you. External IP addresses can be either ephemeral or
static.

#### Reference Guides
* [API Reference](https://cloud.google.com/compute/docs/reference/beta/addresses)
* [Reserving a Static External IP Address](https://cloud.google.com/compute/docs/instances-and-network)
* [Reserving a Static Internal IP Address](https://cloud.google.com/compute/docs/ip-addresses/reserve-static-internal-ip-address)

#### Example

#TODO

#### Reference

```ruby
gcompute_address 'id-for-resource' do
  address            string
  address_type       'INTERNAL' or 'EXTERNAL'
  creation_timestamp time
  description        string
  id                 integer
  name               string
  network_tier       'PREMIUM' or 'STANDARD'
  region             reference to gcompute_region
  subnetwork         reference to gcompute_subnetwork
  users              [
    string,
    ...
  ]
  project            string
  credential         reference to gauth_credential
end
```

#### Actions

* `create` -
  Converges the `gcompute_address` resource into the final
  state described within the block. If the resource does not exist, Chef will
  attempt to create it.
* `delete` -
  Ensures the `gcompute_address` resource is not present.
  If the resource already exists Chef will attempt to delete it.

#### Properties

* `address` -
  The static external IP address represented by this resource. Only
  IPv4 is supported. An address may only be specified for INTERNAL
  address types. The IP address must be inside the specified subnetwork,
  if any.

* `address_type` -
  The type of address to reserve, either INTERNAL or EXTERNAL.
  If unspecified, defaults to EXTERNAL.

* `creation_timestamp` -
  Output only. Creation timestamp in RFC3339 text format.

* `description` -
  An optional description of this resource.

* `id` -
  Output only. The unique identifier for the resource.

* `name` -
  Required. Name of the resource. The name must be 1-63 characters long, and
  comply with RFC1035. Specifically, the name must be 1-63 characters
  long and match the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?`
  which means the first character must be a lowercase letter, and all
  following characters must be a dash, lowercase letter, or digit,
  except the last character, which cannot be a dash.

* `network_tier` -
  The networking tier used for configuring this address. This field can
  take the following values: PREMIUM or STANDARD. If this field is not
  specified, it is assumed to be PREMIUM.

* `subnetwork` -
  The URL of the subnetwork in which to reserve the address. If an IP
  address is specified, it must be within the subnetwork's IP range.
  This field can only be used with INTERNAL type with
  GCE_ENDPOINT/DNS_RESOLVER purposes.

* `users` -
  Output only. The URLs of the resources that are using this address.

* `region` -
  Required. URL of the region where the regional address resides.
  This field is not applicable to global addresses.

#### Label
Set the `a_label` property when attempting to set primary key
of this object. The primary key will always be referred to by the initials of
the resource followed by "_label"

### gcompute_autoscaler
Represents an Autoscaler resource.

Autoscalers allow you to automatically scale virtual machine instances in
managed instance groups according to an autoscaling policy that you
define.

#### Reference Guides
* [API Reference](https://cloud.google.com/compute/docs/reference/rest/v1/autoscalers)
* [Autoscaling Groups of Instances](https://cloud.google.com/compute/docs/autoscaler/)

#### Example

#TODO

#### Reference

```ruby
gcompute_autoscaler 'id-for-resource' do
  autoscaling_policy {
    cool_down_period_sec       integer,
    cpu_utilization            {
      utilization_target double,
    },
    custom_metric_utilizations [
      {
        metric                  string,
        utilization_target      double,
        utilization_target_type 'GAUGE', 'DELTA_PER_SECOND' or 'DELTA_PER_MINUTE',
      },
      ...
    ],
    load_balancing_utilization {
      utilization_target double,
    },
    max_num_replicas           integer,
    min_num_replicas           integer,
  }
  creation_timestamp time
  description        string
  id                 integer
  name               string
  target             reference to gcompute_instance_group_manager
  zone               reference to gcompute_zone
  project            string
  credential         reference to gauth_credential
end
```

#### Actions

* `create` -
  Converges the `gcompute_autoscaler` resource into the final
  state described within the block. If the resource does not exist, Chef will
  attempt to create it.
* `delete` -
  Ensures the `gcompute_autoscaler` resource is not present.
  If the resource already exists Chef will attempt to delete it.

#### Properties

* `id` -
  Output only. Unique identifier for the resource.

* `creation_timestamp` -
  Output only. Creation timestamp in RFC3339 text format.

* `name` -
  Required. Name of the resource. The name must be 1-63 characters long and match
  the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?` which means the
  first character must be a lowercase letter, and all following
  characters must be a dash, lowercase letter, or digit, except the last
  character, which cannot be a dash.

* `description` -
  An optional description of this resource.

* `autoscaling_policy` -
  Required. The configuration parameters for the autoscaling algorithm. You can
  define one or more of the policies for an autoscaler: cpuUtilization,
  customMetricUtilizations, and loadBalancingUtilization.
  If none of these are specified, the default will be to autoscale based
  on cpuUtilization to 0.6 or 60%.

* `autoscaling_policy/min_num_replicas`
  The minimum number of replicas that the autoscaler can scale down
  to. This cannot be less than 0. If not provided, autoscaler will
  choose a default value depending on maximum number of instances
  allowed.

* `autoscaling_policy/max_num_replicas`
  Required. The maximum number of instances that the autoscaler can scale up
  to. This is required when creating or updating an autoscaler. The
  maximum number of replicas should not be lower than minimal number
  of replicas.

* `autoscaling_policy/cool_down_period_sec`
  The number of seconds that the autoscaler should wait before it
  starts collecting information from a new instance. This prevents
  the autoscaler from collecting information when the instance is
  initializing, during which the collected usage would not be
  reliable. The default time autoscaler waits is 60 seconds.
  Virtual machine initialization times might vary because of
  numerous factors. We recommend that you test how long an
  instance may take to initialize. To do this, create an instance
  and time the startup process.

* `autoscaling_policy/cpu_utilization`
  Defines the CPU utilization policy that allows the autoscaler to
  scale based on the average CPU utilization of a managed instance
  group.

* `autoscaling_policy/cpu_utilization/utilization_target`
  The target CPU utilization that the autoscaler should maintain.
  Must be a float value in the range (0, 1]. If not specified, the
  default is 0.6.
  If the CPU level is below the target utilization, the autoscaler
  scales down the number of instances until it reaches the minimum
  number of instances you specified or until the average CPU of
  your instances reaches the target utilization.
  If the average CPU is above the target utilization, the autoscaler
  scales up until it reaches the maximum number of instances you
  specified or until the average utilization reaches the target
  utilization.

* `autoscaling_policy/custom_metric_utilizations`
  Defines the CPU utilization policy that allows the autoscaler to
  scale based on the average CPU utilization of a managed instance
  group.

* `autoscaling_policy/custom_metric_utilizations[]/metric`
  Required. The identifier (type) of the Stackdriver Monitoring metric.
  The metric cannot have negative values.
  The metric must have a value type of INT64 or DOUBLE.

* `autoscaling_policy/custom_metric_utilizations[]/utilization_target`
  Required. The target value of the metric that autoscaler should
  maintain. This must be a positive value. A utilization
  metric scales number of virtual machines handling requests
  to increase or decrease proportionally to the metric.
  For example, a good metric to use as a utilizationTarget is
  www.googleapis.com/compute/instance/network/received_bytes_count.
  The autoscaler will work to keep this value constant for each
  of the instances.

* `autoscaling_policy/custom_metric_utilizations[]/utilization_target_type`
  Required. Defines how target utilization value is expressed for a
  Stackdriver Monitoring metric. Either GAUGE, DELTA_PER_SECOND,
  or DELTA_PER_MINUTE.

* `autoscaling_policy/load_balancing_utilization`
  Configuration parameters of autoscaling based on a load balancer.

* `autoscaling_policy/load_balancing_utilization/utilization_target`
  Fraction of backend capacity utilization (set in HTTP(s) load
  balancing configuration) that autoscaler should maintain. Must
  be a positive float value. If not defined, the default is 0.8.

* `target` -
  Required. URL of the managed instance group that this autoscaler will scale.

* `zone` -
  Required. URL of the zone where the instance group resides.

#### Label
Set the `a_label` property when attempting to set primary key
of this object. The primary key will always be referred to by the initials of
the resource followed by "_label"

### gcompute_backend_bucket
Backend buckets allow you to use Google Cloud Storage buckets with HTTP(S)
load balancing.

An HTTP(S) load balancer can direct traffic to specified URLs to a
backend bucket rather than a backend service. It can send requests for
static content to a Cloud Storage bucket and requests for dynamic content
a virtual machine instance.

#### Reference Guides
* [API Reference](https://cloud.google.com/compute/docs/reference/latest/backendBuckets)
* [Using a Cloud Storage bucket as a load balancer backend](https://cloud.google.com/compute/docs/load-balancing/http/backend-bucket)

#### Example

#TODO

#### Reference

```ruby
gcompute_backend_bucket 'id-for-resource' do
  bucket_name        string
  creation_timestamp time
  description        string
  enable_cdn         boolean
  id                 integer
  name               string
  project            string
  credential         reference to gauth_credential
end
```

#### Actions

* `create` -
  Converges the `gcompute_backend_bucket` resource into the final
  state described within the block. If the resource does not exist, Chef will
  attempt to create it.
* `delete` -
  Ensures the `gcompute_backend_bucket` resource is not present.
  If the resource already exists Chef will attempt to delete it.

#### Properties

* `bucket_name` -
  Required. Cloud Storage bucket name.

* `creation_timestamp` -
  Output only. Creation timestamp in RFC3339 text format.

* `description` -
  An optional textual description of the resource; provided by the
  client when the resource is created.

* `enable_cdn` -
  If true, enable Cloud CDN for this BackendBucket.

* `id` -
  Output only. Unique identifier for the resource.

* `name` -
  Required. Name of the resource. Provided by the client when the resource is
  created. The name must be 1-63 characters long, and comply with
  RFC1035.  Specifically, the name must be 1-63 characters long and
  match the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?` which means
  the first character must be a lowercase letter, and all following
  characters must be a dash, lowercase letter, or digit, except the
  last character, which cannot be a dash.

#### Label
Set the `bb_label` property when attempting to set primary key
of this object. The primary key will always be referred to by the initials of
the resource followed by "_label"

### gcompute_backend_service
Creates a BackendService resource in the specified project using the data
included in the request.


#### Example

#TODO

#### Reference

```ruby
gcompute_backend_service 'id-for-resource' do
  affinity_cookie_ttl_sec integer
  backends                [
    {
      balancing_mode               'UTILIZATION', 'RATE' or 'CONNECTION',
      capacity_scaler              double,
      description                  string,
      group                        reference to gcompute_instance_group,
      max_connections              integer,
      max_connections_per_instance integer,
      max_rate                     integer,
      max_rate_per_instance        double,
      max_utilization              double,
    },
    ...
  ]
  cdn_policy              {
    cache_key_policy {
      include_host           boolean,
      include_protocol       boolean,
      include_query_string   boolean,
      query_string_blacklist [
        string,
        ...
      ],
      query_string_whitelist [
        string,
        ...
      ],
    },
  }
  connection_draining     {
    draining_timeout_sec integer,
  }
  creation_timestamp      time
  description             string
  enable_cdn              boolean
  health_checks           [
    string,
    ...
  ]
  iap                     {
    enabled                     boolean,
    oauth2_client_id            string,
    oauth2_client_secret        string,
    oauth2_client_secret_sha256 string,
  }
  id                      integer
  load_balancing_scheme   'INTERNAL' or 'EXTERNAL'
  name                    string
  port_name               string
  protocol                'HTTP', 'HTTPS', 'TCP' or 'SSL'
  region                  reference to gcompute_region
  session_affinity        'NONE', 'CLIENT_IP', 'GENERATED_COOKIE', 'CLIENT_IP_PROTO' or 'CLIENT_IP_PORT_PROTO'
  timeout_sec             integer
  project                 string
  credential              reference to gauth_credential
end
```

#### Actions

* `create` -
  Converges the `gcompute_backend_service` resource into the final
  state described within the block. If the resource does not exist, Chef will
  attempt to create it.
* `delete` -
  Ensures the `gcompute_backend_service` resource is not present.
  If the resource already exists Chef will attempt to delete it.

#### Properties

* `affinity_cookie_ttl_sec` -
  Lifetime of cookies in seconds if session_affinity is
  GENERATED_COOKIE. If set to 0, the cookie is non-persistent and lasts
  only until the end of the browser session (or equivalent). The
  maximum allowed value for TTL is one day.
  When the load balancing scheme is INTERNAL, this field is not used.

* `backends` -
  The list of backends that serve this BackendService.

* `backends[]/balancing_mode`
  Specifies the balancing mode for this backend.
  For global HTTP(S) or TCP/SSL load balancing, the default is
  UTILIZATION. Valid values are UTILIZATION, RATE (for HTTP(S))
  and CONNECTION (for TCP/SSL).
  This cannot be used for internal load balancing.

* `backends[]/capacity_scaler`
  A multiplier applied to the group's maximum servicing capacity
  (based on UTILIZATION, RATE or CONNECTION).
  Default value is 1, which means the group will serve up to 100%
  of its configured capacity (depending on balancingMode). A
  setting of 0 means the group is completely drained, offering
  0% of its available Capacity. Valid range is [0.0,1.0].
  This cannot be used for internal load balancing.

* `backends[]/description`
  An optional description of this resource.
  Provide this property when you create the resource.

* `backends[]/group`
  This instance group defines the list of instances that serve
  traffic. Member virtual machine instances from each instance
  group must live in the same zone as the instance group itself.
  No two backends in a backend service are allowed to use same
  Instance Group resource.
  When the BackendService has load balancing scheme INTERNAL, the
  instance group must be in a zone within the same region as the
  BackendService.

* `backends[]/max_connections`
  The max number of simultaneous connections for the group. Can
  be used with either CONNECTION or UTILIZATION balancing modes.
  For CONNECTION mode, either maxConnections or
  maxConnectionsPerInstance must be set.
  This cannot be used for internal load balancing.

* `backends[]/max_connections_per_instance`
  The max number of simultaneous connections that a single
  backend instance can handle. This is used to calculate the
  capacity of the group. Can be used in either CONNECTION or
  UTILIZATION balancing modes.
  For CONNECTION mode, either maxConnections or
  maxConnectionsPerInstance must be set.
  This cannot be used for internal load balancing.

* `backends[]/max_rate`
  The max requests per second (RPS) of the group.
  Can be used with either RATE or UTILIZATION balancing modes,
  but required if RATE mode. For RATE mode, either maxRate or
  maxRatePerInstance must be set.
  This cannot be used for internal load balancing.

* `backends[]/max_rate_per_instance`
  The max requests per second (RPS) that a single backend
  instance can handle. This is used to calculate the capacity of
  the group. Can be used in either balancing mode. For RATE mode,
  either maxRate or maxRatePerInstance must be set.
  This cannot be used for internal load balancing.

* `backends[]/max_utilization`
  Used when balancingMode is UTILIZATION. This ratio defines the
  CPU utilization target for the group. The default is 0.8. Valid
  range is [0.0, 1.0].
  This cannot be used for internal load balancing.

* `cdn_policy` -
  Cloud CDN configuration for this BackendService.

* `cdn_policy/cache_key_policy`
  The CacheKeyPolicy for this CdnPolicy.

* `cdn_policy/cache_key_policy/include_host`
  If true requests to different hosts will be cached separately.

* `cdn_policy/cache_key_policy/include_protocol`
  If true, http and https requests will be cached separately.

* `cdn_policy/cache_key_policy/include_query_string`
  If true, include query string parameters in the cache key
  according to query_string_whitelist and
  query_string_blacklist. If neither is set, the entire query
  string will be included.
  If false, the query string will be excluded from the cache
  key entirely.

* `cdn_policy/cache_key_policy/query_string_blacklist`
  Names of query string parameters to exclude in cache keys.
  All other parameters will be included. Either specify
  query_string_whitelist or query_string_blacklist, not both.
  '&' and '=' will be percent encoded and not treated as
  delimiters.

* `cdn_policy/cache_key_policy/query_string_whitelist`
  Names of query string parameters to include in cache keys.
  All other parameters will be excluded. Either specify
  query_string_whitelist or query_string_blacklist, not both.
  '&' and '=' will be percent encoded and not treated as
  delimiters.

* `connection_draining` -
  Settings for connection draining

* `connection_draining/draining_timeout_sec`
  Time for which instance will be drained (not accept new
  connections, but still work to finish started).

* `creation_timestamp` -
  Output only. Creation timestamp in RFC3339 text format.

* `description` -
  An optional description of this resource.

* `enable_cdn` -
  If true, enable Cloud CDN for this BackendService.
  When the load balancing scheme is INTERNAL, this field is not used.

* `health_checks` -
  The list of URLs to the HttpHealthCheck or HttpsHealthCheck resource
  for health checking this BackendService. Currently at most one health
  check can be specified, and a health check is required.
  For internal load balancing, a URL to a HealthCheck resource must be
  specified instead.

* `id` -
  Output only. The unique identifier for the resource.

* `iap` -
  Settings for enabling Cloud Identity Aware Proxy

* `iap/enabled`
  Enables IAP.

* `iap/oauth2_client_id`
  OAuth2 Client ID for IAP

* `iap/oauth2_client_secret`
  OAuth2 Client Secret for IAP

* `iap/oauth2_client_secret_sha256`
  Output only. OAuth2 Client Secret SHA-256 for IAP

* `load_balancing_scheme` -
  Indicates whether the backend service will be used with internal or
  external load balancing. A backend service created for one type of
  load balancing cannot be used with the other.

* `name` -
  Name of the resource. Provided by the client when the resource is
  created. The name must be 1-63 characters long, and comply with
  RFC1035. Specifically, the name must be 1-63 characters long and match
  the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?` which means the
  first character must be a lowercase letter, and all following
  characters must be a dash, lowercase letter, or digit, except the last
  character, which cannot be a dash.

* `port_name` -
  Name of backend port. The same name should appear in the instance
  groups referenced by this service. Required when the load balancing
  scheme is EXTERNAL.
  When the load balancing scheme is INTERNAL, this field is not used.

* `protocol` -
  The protocol this BackendService uses to communicate with backends.
  Possible values are HTTP, HTTPS, TCP, and SSL. The default is HTTP.
  For internal load balancing, the possible values are TCP and UDP, and
  the default is TCP.

* `region` -
  The region where the regional backend service resides.
  This field is not applicable to global backend services.

* `session_affinity` -
  Type of session affinity to use. The default is NONE.
  When the load balancing scheme is EXTERNAL, can be NONE, CLIENT_IP, or
  GENERATED_COOKIE.
  When the load balancing scheme is INTERNAL, can be NONE, CLIENT_IP,
  CLIENT_IP_PROTO, or CLIENT_IP_PORT_PROTO.
  When the protocol is UDP, this field is not used.

* `timeout_sec` -
  How many seconds to wait for the backend before considering it a
  failed request. Default is 30 seconds. Valid range is [1, 86400].

#### Label
Set the `bs_label` property when attempting to set primary key
of this object. The primary key will always be referred to by the initials of
the resource followed by "_label"

### gcompute_disk
Persistent disks are durable storage devices that function similarly to
the physical disks in a desktop or a server. Compute Engine manages the
hardware behind these devices to ensure data redundancy and optimize
performance for you. Persistent disks are available as either standard
hard disk drives (HDD) or solid-state drives (SSD).

Persistent disks are located independently from your virtual machine
instances, so you can detach or move persistent disks to keep your data
even after you delete your instances. Persistent disk performance scales
automatically with size, so you can resize your existing persistent disks
or add more persistent disks to an instance to meet your performance and
storage space requirements.

Add a persistent disk to your instance when you need reliable and
affordable storage with consistent performance characteristics.

#### Reference Guides
* [API Reference](https://cloud.google.com/compute/docs/reference/latest/disks)
* [Adding a persistent disk](https://cloud.google.com/compute/docs/disks/add-persistent-disk)

#### Example

#TODO

#### Reference

```ruby
gcompute_disk 'id-for-resource' do
  creation_timestamp             time
  description                    string
  disk_encryption_key            {
    raw_key string,
    sha256  string,
  }
  id                             integer
  label_fingerprint              fingerprint
  labels                         namevalues
  last_attach_timestamp          time
  last_detach_timestamp          time
  licenses                       [
    string,
    ...
  ]
  name                           string
  size_gb                        integer
  source_image                   string
  source_image_encryption_key    {
    raw_key string,
    sha256  string,
  }
  source_image_id                string
  source_snapshot                reference to gcompute_snapshot
  source_snapshot_encryption_key {
    raw_key string,
    sha256  string,
  }
  source_snapshot_id             string
  type                           reference to gcompute_disk_type
  users                          [
    reference to a gcompute_instance,
    ...
  ]
  zone                           reference to gcompute_zone
  project                        string
  credential                     reference to gauth_credential
end
```

#### Actions

* `create` -
  Converges the `gcompute_disk` resource into the final
  state described within the block. If the resource does not exist, Chef will
  attempt to create it.
* `delete` -
  Ensures the `gcompute_disk` resource is not present.
  If the resource already exists Chef will attempt to delete it.

#### Properties

* `label_fingerprint` -
  Output only. The fingerprint used for optimistic locking of this resource.  Used
  internally during updates.

* `creation_timestamp` -
  Output only. Creation timestamp in RFC3339 text format.

* `description` -
  An optional description of this resource. Provide this property when
  you create the resource.

* `id` -
  Output only. The unique identifier for the resource.

* `last_attach_timestamp` -
  Output only. Last attach timestamp in RFC3339 text format.

* `last_detach_timestamp` -
  Output only. Last dettach timestamp in RFC3339 text format.

* `labels` -
  Labels to apply to this disk.  A list of key->value pairs.

* `licenses` -
  Any applicable publicly visible licenses.

* `name` -
  Required. Name of the resource. Provided by the client when the resource is
  created. The name must be 1-63 characters long, and comply with
  RFC1035. Specifically, the name must be 1-63 characters long and match
  the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?` which means the
  first character must be a lowercase letter, and all following
  characters must be a dash, lowercase letter, or digit, except the last
  character, which cannot be a dash.

* `size_gb` -
  Size of the persistent disk, specified in GB. You can specify this
  field when creating a persistent disk using the sourceImage or
  sourceSnapshot parameter, or specify it alone to create an empty
  persistent disk.
  If you specify this field along with sourceImage or sourceSnapshot,
  the value of sizeGb must not be less than the size of the sourceImage
  or the size of the snapshot.

* `users` -
  Output only. Links to the users of the disk (attached instances) in form:
  project/zones/zone/instances/instance

* `type` -
  URL of the disk type resource describing which disk type to use to
  create the disk. Provide this when creating the disk.

* `source_image` -
  The source image used to create this disk. If the source image is
  deleted, this field will not be set.
  To create a disk with one of the public operating system images,
  specify the image by its family name. For example, specify
  family/debian-8 to use the latest Debian 8 image:
  projects/debian-cloud/global/images/family/debian-8
  Alternatively, use a specific version of a public operating system
  image:
  projects/debian-cloud/global/images/debian-8-jessie-vYYYYMMDD
  To create a disk with a private image that you created, specify the
  image name in the following format:
  global/images/my-private-image
  You can also specify a private image by its image family, which
  returns the latest version of the image in that family. Replace the
  image name with family/family-name:
  global/images/family/my-private-family

* `zone` -
  Required. A reference to the zone where the disk resides.

* `source_image_encryption_key` -
  The customer-supplied encryption key of the source image. Required if
  the source image is protected by a customer-supplied encryption key.

* `source_image_encryption_key/raw_key`
  Specifies a 256-bit customer-supplied encryption key, encoded in
  RFC 4648 base64 to either encrypt or decrypt this resource.

* `source_image_encryption_key/sha256`
  Output only. The RFC 4648 base64 encoded SHA-256 hash of the customer-supplied
  encryption key that protects this resource.

* `source_image_id` -
  Output only. The ID value of the image used to create this disk. This value
  identifies the exact image that was used to create this persistent
  disk. For example, if you created the persistent disk from an image
  that was later deleted and recreated under the same name, the source
  image ID would identify the exact version of the image that was used.

* `disk_encryption_key` -
  Encrypts the disk using a customer-supplied encryption key.
  After you encrypt a disk with a customer-supplied key, you must
  provide the same key if you use the disk later (e.g. to create a disk
  snapshot or an image, or to attach the disk to a virtual machine).
  Customer-supplied encryption keys do not protect access to metadata of
  the disk.
  If you do not provide an encryption key when creating the disk, then
  the disk will be encrypted using an automatically generated key and
  you do not need to provide a key to use the disk later.

* `disk_encryption_key/raw_key`
  Specifies a 256-bit customer-supplied encryption key, encoded in
  RFC 4648 base64 to either encrypt or decrypt this resource.

* `disk_encryption_key/sha256`
  Output only. The RFC 4648 base64 encoded SHA-256 hash of the customer-supplied
  encryption key that protects this resource.

* `source_snapshot` -
  The source snapshot used to create this disk. You can provide this as
  a partial or full URL to the resource. For example, the following are
  valid values:
  * `https://www.googleapis.com/compute/v1/projects/project/global/snapshots/snapshot`
  * `projects/project/global/snapshots/snapshot`
  * `global/snapshots/snapshot`

* `source_snapshot_encryption_key` -
  The customer-supplied encryption key of the source snapshot. Required
  if the source snapshot is protected by a customer-supplied encryption
  key.

* `source_snapshot_encryption_key/raw_key`
  Specifies a 256-bit customer-supplied encryption key, encoded in
  RFC 4648 base64 to either encrypt or decrypt this resource.

* `source_snapshot_encryption_key/sha256`
  Output only. The RFC 4648 base64 encoded SHA-256 hash of the customer-supplied
  encryption key that protects this resource.

* `source_snapshot_id` -
  Output only. The unique ID of the snapshot used to create this disk. This value
  identifies the exact snapshot that was used to create this persistent
  disk. For example, if you created the persistent disk from a snapshot
  that was later deleted and recreated under the same name, the source
  snapshot ID would identify the exact version of the snapshot that was
  used.

#### Label
Set the `d_label` property when attempting to set primary key
of this object. The primary key will always be referred to by the initials of
the resource followed by "_label"

### gcompute_firewall
Each network has its own firewall controlling access to and from the
instances.

All traffic to instances, even from other instances, is blocked by the
firewall unless firewall rules are created to allow it.

The default network has automatically created firewall rules that are
shown in default firewall rules. No manually created network has
automatically created firewall rules except for a default "allow" rule for
outgoing traffic and a default "deny" for incoming traffic. For all
networks except the default network, you must create any firewall rules
you need.

#### Reference Guides
* [API Reference](https://cloud.google.com/compute/docs/reference/latest/firewalls)
* [Official Documentation](https://cloud.google.com/vpc/docs/firewalls)

#### Example

#TODO

#### Reference

```ruby
gcompute_firewall 'id-for-resource' do
  allowed                 [
    {
      ip_protocol string,
      ports       [
        string,
        ...
      ],
    },
    ...
  ]
  creation_timestamp      time
  denied                  [
    {
      ip_protocol string,
      ports       [
        string,
        ...
      ],
    },
    ...
  ]
  description             string
  destination_ranges      [
    string,
    ...
  ]
  direction               'INGRESS' or 'EGRESS'
  disabled                boolean
  id                      integer
  name                    string
  network                 reference to gcompute_network
  priority                integer
  source_ranges           [
    string,
    ...
  ]
  source_service_accounts [
    string,
    ...
  ]
  source_tags             [
    string,
    ...
  ]
  target_service_accounts [
    string,
    ...
  ]
  target_tags             [
    string,
    ...
  ]
  project                 string
  credential              reference to gauth_credential
end
```

#### Actions

* `create` -
  Converges the `gcompute_firewall` resource into the final
  state described within the block. If the resource does not exist, Chef will
  attempt to create it.
* `delete` -
  Ensures the `gcompute_firewall` resource is not present.
  If the resource already exists Chef will attempt to delete it.

#### Properties

* `allowed` -
  The list of ALLOW rules specified by this firewall. Each rule
  specifies a protocol and port-range tuple that describes a permitted
  connection.

* `allowed[]/ip_protocol`
  Required. The IP protocol to which this rule applies. The protocol type is
  required when creating a firewall rule. This value can either be
  one of the following well known protocol strings (tcp, udp,
  icmp, esp, ah, sctp), or the IP protocol number.

* `allowed[]/ports`
  An optional list of ports to which this rule applies. This field
  is only applicable for UDP or TCP protocol. Each entry must be
  either an integer or a range. If not specified, this rule
  applies to connections through any port.
  Example inputs include: ["22"], ["80","443"], and
  ["12345-12349"].

* `creation_timestamp` -
  Output only. Creation timestamp in RFC3339 text format.

* `denied` -
  The list of DENY rules specified by this firewall. Each rule specifies
  a protocol and port-range tuple that describes a denied connection.

* `denied[]/ip_protocol`
  Required. The IP protocol to which this rule applies. The protocol type is
  required when creating a firewall rule. This value can either be
  one of the following well known protocol strings (tcp, udp,
  icmp, esp, ah, sctp), or the IP protocol number.

* `denied[]/ports`
  An optional list of ports to which this rule applies. This field
  is only applicable for UDP or TCP protocol. Each entry must be
  either an integer or a range. If not specified, this rule
  applies to connections through any port.
  Example inputs include: ["22"], ["80","443"], and
  ["12345-12349"].

* `description` -
  An optional description of this resource. Provide this property when
  you create the resource.

* `destination_ranges` -
  If destination ranges are specified, the firewall will apply only to
  traffic that has destination IP address in these ranges. These ranges
  must be expressed in CIDR format. Only IPv4 is supported.

* `direction` -
  Direction of traffic to which this firewall applies; default is
  INGRESS. Note: For INGRESS traffic, it is NOT supported to specify
  destinationRanges; For EGRESS traffic, it is NOT supported to specify
  sourceRanges OR sourceTags.

* `disabled` -
  Denotes whether the firewall rule is disabled, i.e not applied to the
  network it is associated with. When set to true, the firewall rule is
  not enforced and the network behaves as if it did not exist. If this
  is unspecified, the firewall rule will be enabled.

* `id` -
  Output only. The unique identifier for the resource.

* `name` -
  Required. Name of the resource. Provided by the client when the resource is
  created. The name must be 1-63 characters long, and comply with
  RFC1035. Specifically, the name must be 1-63 characters long and match
  the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?` which means the
  first character must be a lowercase letter, and all following
  characters must be a dash, lowercase letter, or digit, except the last
  character, which cannot be a dash.

* `network` -
  Required. URL of the network resource for this firewall rule. If not specified
  when creating a firewall rule, the default network is used:
  global/networks/default
  If you choose to specify this property, you can specify the network as
  a full or partial URL. For example, the following are all valid URLs:
  https://www.googleapis.com/compute/v1/projects/myproject/global/
  networks/my-network
  projects/myproject/global/networks/my-network
  global/networks/default

* `priority` -
  Priority for this rule. This is an integer between 0 and 65535, both
  inclusive. When not specified, the value assumed is 1000. Relative
  priorities determine precedence of conflicting rules. Lower value of
  priority implies higher precedence (eg, a rule with priority 0 has
  higher precedence than a rule with priority 1). DENY rules take
  precedence over ALLOW rules having equal priority.

* `source_ranges` -
  If source ranges are specified, the firewall will apply only to
  traffic that has source IP address in these ranges. These ranges must
  be expressed in CIDR format. One or both of sourceRanges and
  sourceTags may be set. If both properties are set, the firewall will
  apply to traffic that has source IP address within sourceRanges OR the
  source IP that belongs to a tag listed in the sourceTags property. The
  connection does not need to match both properties for the firewall to
  apply. Only IPv4 is supported.

* `source_service_accounts` -
  If source service accounts are specified, the firewall will apply only
  to traffic originating from an instance with a service account in this
  list. Source service accounts cannot be used to control traffic to an
  instance's external IP address because service accounts are associated
  with an instance, not an IP address. sourceRanges can be set at the
  same time as sourceServiceAccounts. If both are set, the firewall will
  apply to traffic that has source IP address within sourceRanges OR the
  source IP belongs to an instance with service account listed in
  sourceServiceAccount. The connection does not need to match both
  properties for the firewall to apply. sourceServiceAccounts cannot be
  used at the same time as sourceTags or targetTags.

* `source_tags` -
  If source tags are specified, the firewall will apply only to traffic
  with source IP that belongs to a tag listed in source tags. Source
  tags cannot be used to control traffic to an instance's external IP
  address. Because tags are associated with an instance, not an IP
  address. One or both of sourceRanges and sourceTags may be set. If
  both properties are set, the firewall will apply to traffic that has
  source IP address within sourceRanges OR the source IP that belongs to
  a tag listed in the sourceTags property. The connection does not need
  to match both properties for the firewall to apply.

* `target_service_accounts` -
  A list of service accounts indicating sets of instances located in the
  network that may make network connections as specified in allowed[].
  targetServiceAccounts cannot be used at the same time as targetTags or
  sourceTags. If neither targetServiceAccounts nor targetTags are
  specified, the firewall rule applies to all instances on the specified
  network.

* `target_tags` -
  A list of instance tags indicating sets of instances located in the
  network that may make network connections as specified in allowed[].
  If no targetTags are specified, the firewall rule applies to all
  instances on the specified network.

#### Label
Set the `f_label` property when attempting to set primary key
of this object. The primary key will always be referred to by the initials of
the resource followed by "_label"

### gcompute_forwarding_rule
A ForwardingRule resource. A ForwardingRule resource specifies which pool
of target virtual machines to forward a packet to if it matches the given
[IPAddress, IPProtocol, portRange] tuple.

#### Reference Guides
* [API Reference](https://cloud.google.com/compute/docs/reference/latest/forwardingRule)
* [Official Documentation](https://cloud.google.com/compute/docs/load-balancing/network/forwarding-rules)

#### Example

#TODO

#### Reference

```ruby
gcompute_forwarding_rule 'id-for-resource' do
  backend_service       reference to gcompute_backend_service
  creation_timestamp    time
  description           string
  id                    integer
  ip_address            string
  ip_protocol           'TCP', 'UDP', 'ESP', 'AH', 'SCTP' or 'ICMP'
  ip_version            'IPV4' or 'IPV6'
  label_fingerprint     fingerprint
  load_balancing_scheme 'INTERNAL' or 'EXTERNAL'
  name                  string
  network               reference to gcompute_network
  network_tier          'PREMIUM' or 'STANDARD'
  port_range            string
  ports                 [
    string,
    ...
  ]
  region                reference to gcompute_region
  subnetwork            reference to gcompute_subnetwork
  target                reference to gcompute_target_pool
  project               string
  credential            reference to gauth_credential
end
```

#### Actions

* `create` -
  Converges the `gcompute_forwarding_rule` resource into the final
  state described within the block. If the resource does not exist, Chef will
  attempt to create it.
* `delete` -
  Ensures the `gcompute_forwarding_rule` resource is not present.
  If the resource already exists Chef will attempt to delete it.

#### Properties

* `creation_timestamp` -
  Output only. Creation timestamp in RFC3339 text format.

* `description` -
  An optional description of this resource. Provide this property when
  you create the resource.

* `id` -
  Output only. The unique identifier for the resource.

* `ip_address` -
  The IP address that this forwarding rule is serving on behalf of.
  Addresses are restricted based on the forwarding rule's load balancing
  scheme (EXTERNAL or INTERNAL) and scope (global or regional).
  When the load balancing scheme is EXTERNAL, for global forwarding
  rules, the address must be a global IP, and for regional forwarding
  rules, the address must live in the same region as the forwarding
  rule. If this field is empty, an ephemeral IPv4 address from the same
  scope (global or regional) will be assigned. A regional forwarding
  rule supports IPv4 only. A global forwarding rule supports either IPv4
  or IPv6.
  When the load balancing scheme is INTERNAL, this can only be an RFC
  1918 IP address belonging to the network/subnet configured for the
  forwarding rule. By default, if this field is empty, an ephemeral
  internal IP address will be automatically allocated from the IP range
  of the subnet or network configured for this forwarding rule.
  An address can be specified either by a literal IP address or a URL
  reference to an existing Address resource. The following examples are
  all valid:
  * 100.1.2.3
  * https://www.googleapis.com/compute/v1/projects/project/regions/
  region/addresses/address
  * projects/project/regions/region/addresses/address
  * regions/region/addresses/address
  * global/addresses/address
  * address

* `ip_protocol` -
  The IP protocol to which this rule applies. Valid options are TCP,
  UDP, ESP, AH, SCTP or ICMP.
  When the load balancing scheme is INTERNAL, only TCP and UDP are
  valid.

* `backend_service` -
  A reference to a BackendService to receive the matched traffic.
  This is used for internal load balancing.
  (not used for external load balancing)

* `ip_version` -
  The IP Version that will be used by this forwarding rule. Valid
  options are IPV4 or IPV6. This can only be specified for a global
  forwarding rule.

* `load_balancing_scheme` -
  This signifies what the ForwardingRule will be used for and can only
  take the following values: INTERNAL, EXTERNAL The value of INTERNAL
  means that this will be used for Internal Network Load Balancing (TCP,
  UDP). The value of EXTERNAL means that this will be used for External
  Load Balancing (HTTP(S) LB, External TCP/UDP LB, SSL Proxy)

* `name` -
  Required. Name of the resource; provided by the client when the resource is
  created. The name must be 1-63 characters long, and comply with
  RFC1035. Specifically, the name must be 1-63 characters long and match
  the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?` which means the
  first character must be a lowercase letter, and all following
  characters must be a dash, lowercase letter, or digit, except the last
  character, which cannot be a dash.

* `network` -
  For internal load balancing, this field identifies the network that
  the load balanced IP should belong to for this Forwarding Rule. If
  this field is not specified, the default network will be used.
  This field is not used for external load balancing.

* `port_range` -
  This field is used along with the target field for TargetHttpProxy,
  TargetHttpsProxy, TargetSslProxy, TargetTcpProxy, TargetVpnGateway,
  TargetPool, TargetInstance.
  Applicable only when IPProtocol is TCP, UDP, or SCTP, only packets
  addressed to ports in the specified range will be forwarded to target.
  Forwarding rules with the same [IPAddress, IPProtocol] pair must have
  disjoint port ranges.
  Some types of forwarding target have constraints on the acceptable
  ports:
  * TargetHttpProxy: 80, 8080
  * TargetHttpsProxy: 443
  * TargetTcpProxy: 25, 43, 110, 143, 195, 443, 465, 587, 700, 993, 995,
  1883, 5222
  * TargetSslProxy: 25, 43, 110, 143, 195, 443, 465, 587, 700, 993, 995,
  1883, 5222
  * TargetVpnGateway: 500, 4500

* `ports` -
  This field is used along with the backend_service field for internal
  load balancing.
  When the load balancing scheme is INTERNAL, a single port or a comma
  separated list of ports can be configured. Only packets addressed to
  these ports will be forwarded to the backends configured with this
  forwarding rule.
  You may specify a maximum of up to 5 ports.

* `subnetwork` -
  A reference to a subnetwork.
  For internal load balancing, this field identifies the subnetwork that
  the load balanced IP should belong to for this Forwarding Rule.
  If the network specified is in auto subnet mode, this field is
  optional. However, if the network is in custom subnet mode, a
  subnetwork must be specified.
  This field is not used for external load balancing.

* `target` -
  A reference to a TargetPool resource to receive the matched traffic.
  For regional forwarding rules, this target must live in the same
  region as the forwarding rule. For global forwarding rules, this
  target must be a global load balancing resource. The forwarded traffic
  must be of a type appropriate to the target object.
  This field is not used for internal load balancing.

* `label_fingerprint` -
  Output only. The fingerprint used for optimistic locking of this resource.  Used
  internally during updates.

* `network_tier` -
  The networking tier used for configuring this address. This field can
  take the following values: PREMIUM or STANDARD. If this field is not
  specified, it is assumed to be PREMIUM.

* `region` -
  Required. A reference to the region where the regional forwarding rule resides.
  This field is not applicable to global forwarding rules.

#### Label
Set the `fr_label` property when attempting to set primary key
of this object. The primary key will always be referred to by the initials of
the resource followed by "_label"

### gcompute_global_address
Represents a Global Address resource. Global addresses are used for
HTTP(S) load balancing.

#### Reference Guides
* [API Reference](https://cloud.google.com/compute/docs/reference/latest/globalAddresses)
* [Reserving a Static External IP Address](https://cloud.google.com/compute/docs/ip-addresses/reserve-static-external-ip-address)

#### Example

#TODO

#### Reference

```ruby
gcompute_global_address 'id-for-resource' do
  address            string
  address_type       'EXTERNAL' or 'INTERNAL'
  creation_timestamp time
  description        string
  id                 integer
  ip_version         'IPV4' or 'IPV6'
  label_fingerprint  fingerprint
  name               string
  region             reference to gcompute_region
  project            string
  credential         reference to gauth_credential
end
```

#### Actions

* `create` -
  Converges the `gcompute_global_address` resource into the final
  state described within the block. If the resource does not exist, Chef will
  attempt to create it.
* `delete` -
  Ensures the `gcompute_global_address` resource is not present.
  If the resource already exists Chef will attempt to delete it.

#### Properties

* `address` -
  Output only. The static external IP address represented by this resource.

* `creation_timestamp` -
  Output only. Creation timestamp in RFC3339 text format.

* `description` -
  An optional description of this resource.
  Provide this property when you create the resource.

* `id` -
  Output only. The unique identifier for the resource. This identifier is defined by
  the server.

* `name` -
  Required. Name of the resource. Provided by the client when the resource is
  created. The name must be 1-63 characters long, and comply with
  RFC1035.  Specifically, the name must be 1-63 characters long and
  match the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?` which means
  the first character must be a lowercase letter, and all following
  characters must be a dash, lowercase letter, or digit, except the last
  character, which cannot be a dash.

* `label_fingerprint` -
  Output only. The fingerprint used for optimistic locking of this resource.  Used
  internally during updates.

* `ip_version` -
  The IP Version that will be used by this address. Valid options are
  IPV4 or IPV6. The default value is IPV4.

* `region` -
  Output only. A reference to the region where the regional address resides.

* `address_type` -
  The type of the address to reserve, default is EXTERNAL.
  * EXTERNAL indicates public/external single IP address.
  * INTERNAL indicates internal IP ranges belonging to some network.

#### Label
Set the `ga_label` property when attempting to set primary key
of this object. The primary key will always be referred to by the initials of
the resource followed by "_label"

### gcompute_global_forwarding_rule
Represents a GlobalForwardingRule resource. Global forwarding rules are
used to forward traffic to the correct load balancer for HTTP load
balancing. Global forwarding rules can only be used for HTTP load
balancing.

For more information, see
https://cloud.google.com/compute/docs/load-balancing/http/


#### Example

#TODO

#### Reference

```ruby
gcompute_global_forwarding_rule 'id-for-resource' do
  backend_service       reference to gcompute_backend_service
  creation_timestamp    time
  description           string
  id                    integer
  ip_address            string
  ip_protocol           'TCP', 'UDP', 'ESP', 'AH', 'SCTP' or 'ICMP'
  ip_version            'IPV4' or 'IPV6'
  load_balancing_scheme 'INTERNAL' or 'EXTERNAL'
  name                  string
  network               reference to gcompute_network
  port_range            string
  ports                 [
    string,
    ...
  ]
  region                reference to gcompute_region
  subnetwork            reference to gcompute_subnetwork
  target                string
  project               string
  credential            reference to gauth_credential
end
```

#### Actions

* `create` -
  Converges the `gcompute_global_forwarding_rule` resource into the final
  state described within the block. If the resource does not exist, Chef will
  attempt to create it.
* `delete` -
  Ensures the `gcompute_global_forwarding_rule` resource is not present.
  If the resource already exists Chef will attempt to delete it.

#### Properties

* `creation_timestamp` -
  Output only. Creation timestamp in RFC3339 text format.

* `description` -
  An optional description of this resource. Provide this property when
  you create the resource.

* `id` -
  Output only. The unique identifier for the resource.

* `ip_address` -
  The IP address that this forwarding rule is serving on behalf of.
  Addresses are restricted based on the forwarding rule's load balancing
  scheme (EXTERNAL or INTERNAL) and scope (global or regional).
  When the load balancing scheme is EXTERNAL, for global forwarding
  rules, the address must be a global IP, and for regional forwarding
  rules, the address must live in the same region as the forwarding
  rule. If this field is empty, an ephemeral IPv4 address from the same
  scope (global or regional) will be assigned. A regional forwarding
  rule supports IPv4 only. A global forwarding rule supports either IPv4
  or IPv6.
  When the load balancing scheme is INTERNAL, this can only be an RFC
  1918 IP address belonging to the network/subnet configured for the
  forwarding rule. By default, if this field is empty, an ephemeral
  internal IP address will be automatically allocated from the IP range
  of the subnet or network configured for this forwarding rule.
  An address can be specified either by a literal IP address or a URL
  reference to an existing Address resource. The following examples are
  all valid:
  * 100.1.2.3
  * https://www.googleapis.com/compute/v1/projects/project/regions/
  region/addresses/address
  * projects/project/regions/region/addresses/address
  * regions/region/addresses/address
  * global/addresses/address
  * address

* `ip_protocol` -
  The IP protocol to which this rule applies. Valid options are TCP,
  UDP, ESP, AH, SCTP or ICMP.
  When the load balancing scheme is INTERNAL, only TCP and UDP are
  valid.

* `backend_service` -
  A reference to a BackendService to receive the matched traffic.
  This is used for internal load balancing.
  (not used for external load balancing)

* `ip_version` -
  The IP Version that will be used by this forwarding rule. Valid
  options are IPV4 or IPV6. This can only be specified for a global
  forwarding rule.

* `load_balancing_scheme` -
  This signifies what the ForwardingRule will be used for and can only
  take the following values: INTERNAL, EXTERNAL The value of INTERNAL
  means that this will be used for Internal Network Load Balancing (TCP,
  UDP). The value of EXTERNAL means that this will be used for External
  Load Balancing (HTTP(S) LB, External TCP/UDP LB, SSL Proxy)

* `name` -
  Required. Name of the resource; provided by the client when the resource is
  created. The name must be 1-63 characters long, and comply with
  RFC1035. Specifically, the name must be 1-63 characters long and match
  the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?` which means the
  first character must be a lowercase letter, and all following
  characters must be a dash, lowercase letter, or digit, except the last
  character, which cannot be a dash.

* `network` -
  For internal load balancing, this field identifies the network that
  the load balanced IP should belong to for this Forwarding Rule. If
  this field is not specified, the default network will be used.
  This field is not used for external load balancing.

* `port_range` -
  This field is used along with the target field for TargetHttpProxy,
  TargetHttpsProxy, TargetSslProxy, TargetTcpProxy, TargetVpnGateway,
  TargetPool, TargetInstance.
  Applicable only when IPProtocol is TCP, UDP, or SCTP, only packets
  addressed to ports in the specified range will be forwarded to target.
  Forwarding rules with the same [IPAddress, IPProtocol] pair must have
  disjoint port ranges.
  Some types of forwarding target have constraints on the acceptable
  ports:
  * TargetHttpProxy: 80, 8080
  * TargetHttpsProxy: 443
  * TargetTcpProxy: 25, 43, 110, 143, 195, 443, 465, 587, 700, 993, 995,
  1883, 5222
  * TargetSslProxy: 25, 43, 110, 143, 195, 443, 465, 587, 700, 993, 995,
  1883, 5222
  * TargetVpnGateway: 500, 4500

* `ports` -
  This field is used along with the backend_service field for internal
  load balancing.
  When the load balancing scheme is INTERNAL, a single port or a comma
  separated list of ports can be configured. Only packets addressed to
  these ports will be forwarded to the backends configured with this
  forwarding rule.
  You may specify a maximum of up to 5 ports.

* `subnetwork` -
  A reference to a subnetwork.
  For internal load balancing, this field identifies the subnetwork that
  the load balanced IP should belong to for this Forwarding Rule.
  If the network specified is in auto subnet mode, this field is
  optional. However, if the network is in custom subnet mode, a
  subnetwork must be specified.
  This field is not used for external load balancing.

* `region` -
  Output only. A reference to the region where the regional forwarding rule resides.
  This field is not applicable to global forwarding rules.

* `target` -
  This target must be a global load balancing resource. The forwarded
  traffic must be of a type appropriate to the target object.
  Valid types: HTTP_PROXY, HTTPS_PROXY, SSL_PROXY, TCP_PROXY

#### Label
Set the `gfr_label` property when attempting to set primary key
of this object. The primary key will always be referred to by the initials of
the resource followed by "_label"

### gcompute_http_health_check
An HttpHealthCheck resource. This resource defines a template for how
individual VMs should be checked for health, via HTTP.

#### Reference Guides
* [API Reference](https://cloud.google.com/compute/docs/reference/latest/httpHealthChecks)
* [Adding Health Checks](https://cloud.google.com/compute/docs/load-balancing/health-checks#legacy_health_checks)

#### Example

#TODO

#### Reference

```ruby
gcompute_http_health_check 'id-for-resource' do
  check_interval_sec  integer
  creation_timestamp  time
  description         string
  healthy_threshold   integer
  host                string
  id                  integer
  name                string
  port                integer
  request_path        string
  timeout_sec         integer
  unhealthy_threshold integer
  project             string
  credential          reference to gauth_credential
end
```

#### Actions

* `create` -
  Converges the `gcompute_http_health_check` resource into the final
  state described within the block. If the resource does not exist, Chef will
  attempt to create it.
* `delete` -
  Ensures the `gcompute_http_health_check` resource is not present.
  If the resource already exists Chef will attempt to delete it.

#### Properties

* `check_interval_sec` -
  How often (in seconds) to send a health check. The default value is 5
  seconds.

* `creation_timestamp` -
  Output only. Creation timestamp in RFC3339 text format.

* `description` -
  An optional description of this resource. Provide this property when
  you create the resource.

* `healthy_threshold` -
  A so-far unhealthy instance will be marked healthy after this many
  consecutive successes. The default value is 2.

* `host` -
  The value of the host header in the HTTP health check request. If
  left empty (default value), the public IP on behalf of which this
  health check is performed will be used.

* `id` -
  Output only. The unique identifier for the resource. This identifier is defined by
  the server.

* `name` -
  Required. Name of the resource. Provided by the client when the resource is
  created. The name must be 1-63 characters long, and comply with
  RFC1035.  Specifically, the name must be 1-63 characters long and
  match the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?` which means
  the first character must be a lowercase letter, and all following
  characters must be a dash, lowercase letter, or digit, except the
  last character, which cannot be a dash.

* `port` -
  The TCP port number for the HTTP health check request.
  The default value is 80.

* `request_path` -
  The request path of the HTTP health check request.
  The default value is /.

* `timeout_sec` -
  How long (in seconds) to wait before claiming failure.
  The default value is 5 seconds.  It is invalid for timeoutSec to have
  greater value than checkIntervalSec.

* `unhealthy_threshold` -
  A so-far healthy instance will be marked unhealthy after this many
  consecutive failures. The default value is 2.

#### Label
Set the `hhc_label` property when attempting to set primary key
of this object. The primary key will always be referred to by the initials of
the resource followed by "_label"

### gcompute_https_health_check
An HttpsHealthCheck resource. This resource defines a template for how
individual VMs should be checked for health, via HTTPS.

#### Reference Guides
* [API Reference](https://cloud.google.com/compute/docs/reference/latest/httpsHealthChecks)
* [Adding Health Checks](https://cloud.google.com/compute/docs/load-balancing/health-checks#legacy_health_checks)

#### Example

#TODO

#### Reference

```ruby
gcompute_https_health_check 'id-for-resource' do
  check_interval_sec  integer
  creation_timestamp  time
  description         string
  healthy_threshold   integer
  host                string
  id                  integer
  name                string
  port                integer
  request_path        string
  timeout_sec         integer
  unhealthy_threshold integer
  project             string
  credential          reference to gauth_credential
end
```

#### Actions

* `create` -
  Converges the `gcompute_https_health_check` resource into the final
  state described within the block. If the resource does not exist, Chef will
  attempt to create it.
* `delete` -
  Ensures the `gcompute_https_health_check` resource is not present.
  If the resource already exists Chef will attempt to delete it.

#### Properties

* `check_interval_sec` -
  How often (in seconds) to send a health check. The default value is 5
  seconds.

* `creation_timestamp` -
  Output only. Creation timestamp in RFC3339 text format.

* `description` -
  An optional description of this resource. Provide this property when
  you create the resource.

* `healthy_threshold` -
  A so-far unhealthy instance will be marked healthy after this many
  consecutive successes. The default value is 2.

* `host` -
  The value of the host header in the HTTPS health check request. If
  left empty (default value), the public IP on behalf of which this
  health check is performed will be used.

* `id` -
  Output only. The unique identifier for the resource. This identifier is defined by
  the server.

* `name` -
  Required. Name of the resource. Provided by the client when the resource is
  created. The name must be 1-63 characters long, and comply with
  RFC1035.  Specifically, the name must be 1-63 characters long and
  match the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?` which means
  the first character must be a lowercase letter, and all following
  characters must be a dash, lowercase letter, or digit, except the
  last character, which cannot be a dash.

* `port` -
  The TCP port number for the HTTPS health check request.
  The default value is 80.

* `request_path` -
  The request path of the HTTPS health check request.
  The default value is /.

* `timeout_sec` -
  How long (in seconds) to wait before claiming failure.
  The default value is 5 seconds.  It is invalid for timeoutSec to have
  greater value than checkIntervalSec.

* `unhealthy_threshold` -
  A so-far healthy instance will be marked unhealthy after this many
  consecutive failures. The default value is 2.

#### Label
Set the `hhc_label` property when attempting to set primary key
of this object. The primary key will always be referred to by the initials of
the resource followed by "_label"

### gcompute_health_check
Health Checks determine whether instances are responsive and able to do work.
They are an important part of a comprehensive load balancing configuration,
as they enable monitoring instances behind load balancers.

Health Checks poll instances at a specified interval. Instances that
do not respond successfully to some number of probes in a row are marked
as unhealthy. No new connections are sent to unhealthy instances,
though existing connections will continue. The health check will
continue to poll unhealthy instances. If an instance later responds
successfully to some number of consecutive probes, it is marked
healthy again and can receive new connections.

#### Reference Guides
* [API Reference](https://cloud.google.com/compute/docs/reference/rest/latest/healthChecks)
* [Official Documentation](https://cloud.google.com/load-balancing/docs/health-checks)

#### Example

#TODO

#### Reference

```ruby
gcompute_health_check 'id-for-resource' do
  check_interval_sec  integer
  creation_timestamp  time
  description         string
  healthy_threshold   integer
  http_health_check   {
    host         string,
    port         integer,
    port_name    string,
    proxy_header 'NONE' or 'PROXY_V1',
    request_path string,
  }
  https_health_check  {
    host         string,
    port         integer,
    port_name    string,
    proxy_header 'NONE' or 'PROXY_V1',
    request_path string,
  }
  id                  integer
  name                string
  ssl_health_check    {
    port         integer,
    port_name    string,
    proxy_header 'NONE' or 'PROXY_V1',
    request      string,
    response     string,
  }
  tcp_health_check    {
    port         integer,
    port_name    string,
    proxy_header 'NONE' or 'PROXY_V1',
    request      string,
    response     string,
  }
  timeout_sec         integer
  type                'TCP', 'SSL', 'HTTP' or 'HTTPS'
  unhealthy_threshold integer
  project             string
  credential          reference to gauth_credential
end
```

#### Actions

* `create` -
  Converges the `gcompute_health_check` resource into the final
  state described within the block. If the resource does not exist, Chef will
  attempt to create it.
* `delete` -
  Ensures the `gcompute_health_check` resource is not present.
  If the resource already exists Chef will attempt to delete it.

#### Properties

* `check_interval_sec` -
  How often (in seconds) to send a health check. The default value is 5
  seconds.

* `creation_timestamp` -
  Output only. Creation timestamp in RFC3339 text format.

* `description` -
  An optional description of this resource. Provide this property when
  you create the resource.

* `healthy_threshold` -
  A so-far unhealthy instance will be marked healthy after this many
  consecutive successes. The default value is 2.

* `id` -
  Output only. The unique identifier for the resource. This identifier is defined by
  the server.

* `name` -
  Required. Name of the resource. Provided by the client when the resource is
  created. The name must be 1-63 characters long, and comply with
  RFC1035.  Specifically, the name must be 1-63 characters long and
  match the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?` which means
  the first character must be a lowercase letter, and all following
  characters must be a dash, lowercase letter, or digit, except the
  last character, which cannot be a dash.

* `timeout_sec` -
  How long (in seconds) to wait before claiming failure.
  The default value is 5 seconds.  It is invalid for timeoutSec to have
  greater value than checkIntervalSec.

* `unhealthy_threshold` -
  A so-far healthy instance will be marked unhealthy after this many
  consecutive failures. The default value is 2.

* `type` -
  Specifies the type of the healthCheck, either TCP, SSL, HTTP or
  HTTPS. If not specified, the default is TCP. Exactly one of the
  protocol-specific health check field must be specified, which must
  match type field.

* `http_health_check` -
  A nested object resource

* `http_health_check/host`
  The value of the host header in the HTTP health check request.
  If left empty (default value), the public IP on behalf of which this health
  check is performed will be used.

* `http_health_check/request_path`
  The request path of the HTTP health check request.
  The default value is /.

* `http_health_check/port`
  The TCP port number for the HTTP health check request.
  The default value is 80.

* `http_health_check/port_name`
  Port name as defined in InstanceGroup#NamedPort#name. If both port and
  port_name are defined, port takes precedence.

* `http_health_check/proxy_header`
  Specifies the type of proxy header to append before sending data to the
  backend, either NONE or PROXY_V1. The default is NONE.

* `https_health_check` -
  A nested object resource

* `https_health_check/host`
  The value of the host header in the HTTPS health check request.
  If left empty (default value), the public IP on behalf of which this health
  check is performed will be used.

* `https_health_check/request_path`
  The request path of the HTTPS health check request.
  The default value is /.

* `https_health_check/port`
  The TCP port number for the HTTPS health check request.
  The default value is 443.

* `https_health_check/port_name`
  Port name as defined in InstanceGroup#NamedPort#name. If both port and
  port_name are defined, port takes precedence.

* `https_health_check/proxy_header`
  Specifies the type of proxy header to append before sending data to the
  backend, either NONE or PROXY_V1. The default is NONE.

* `tcp_health_check` -
  A nested object resource

* `tcp_health_check/request`
  The application data to send once the TCP connection has been
  established (default value is empty). If both request and response are
  empty, the connection establishment alone will indicate health. The request
  data can only be ASCII.

* `tcp_health_check/response`
  The bytes to match against the beginning of the response data. If left empty
  (the default value), any response will indicate health. The response data
  can only be ASCII.

* `tcp_health_check/port`
  The TCP port number for the TCP health check request.
  The default value is 443.

* `tcp_health_check/port_name`
  Port name as defined in InstanceGroup#NamedPort#name. If both port and
  port_name are defined, port takes precedence.

* `tcp_health_check/proxy_header`
  Specifies the type of proxy header to append before sending data to the
  backend, either NONE or PROXY_V1. The default is NONE.

* `ssl_health_check` -
  A nested object resource

* `ssl_health_check/request`
  The application data to send once the SSL connection has been
  established (default value is empty). If both request and response are
  empty, the connection establishment alone will indicate health. The request
  data can only be ASCII.

* `ssl_health_check/response`
  The bytes to match against the beginning of the response data. If left empty
  (the default value), any response will indicate health. The response data
  can only be ASCII.

* `ssl_health_check/port`
  The TCP port number for the SSL health check request.
  The default value is 443.

* `ssl_health_check/port_name`
  Port name as defined in InstanceGroup#NamedPort#name. If both port and
  port_name are defined, port takes precedence.

* `ssl_health_check/proxy_header`
  Specifies the type of proxy header to append before sending data to the
  backend, either NONE or PROXY_V1. The default is NONE.

#### Label
Set the `hc_label` property when attempting to set primary key
of this object. The primary key will always be referred to by the initials of
the resource followed by "_label"

### gcompute_instance_template
Defines an Instance Template resource that provides configuration settings
for your virtual machine instances. Instance templates are not tied to the
lifetime of an instance and can be used and reused as to deploy virtual
machines. You can also use different templates to create different virtual
machine configurations. Instance templates are required when you create a
managed instance group.

Tip: Disks should be set to autoDelete=true
so that leftover disks are not left behind on machine deletion.


#### Example

#TODO

#### Reference

```ruby
gcompute_instance_template 'id-for-resource' do
  creation_timestamp time
  description        string
  id                 integer
  name               string
  properties         {
    can_ip_forward     boolean,
    description        string,
    disks              [
      {
        auto_delete         boolean,
        boot                boolean,
        device_name         string,
        disk_encryption_key {
          raw_key           string,
          rsa_encrypted_key string,
          sha256            string,
        },
        index               integer,
        initialize_params   {
          disk_name                   string,
          disk_size_gb                integer,
          disk_type                   reference to gcompute_disk_type,
          source_image                string,
          source_image_encryption_key {
            raw_key string,
            sha256  string,
          },
        },
        interface           'SCSI' or 'NVME',
        mode                'READ_WRITE' or 'READ_ONLY',
        source              reference to gcompute_disk,
        type                'SCRATCH' or 'PERSISTENT',
      },
      ...
    ],
    guest_accelerators [
      {
        accelerator_count integer,
        accelerator_type  string,
      },
      ...
    ],
    machine_type       reference to gcompute_machine_type,
    metadata           namevalues,
    network_interfaces [
      {
        access_configs  [
          {
            name   string,
            nat_ip reference to gcompute_address,
            type   ONE_TO_ONE_NAT,
          },
          ...
        ],
        alias_ip_ranges [
          {
            ip_cidr_range         string,
            subnetwork_range_name string,
          },
          ...
        ],
        name            string,
        network         reference to gcompute_network,
        network_ip      string,
        subnetwork      reference to gcompute_subnetwork,
      },
      ...
    ],
    scheduling         {
      automatic_restart   boolean,
      on_host_maintenance string,
      preemptible         boolean,
    },
    service_accounts   [
      {
        email  string,
        scopes [
          string,
          ...
        ],
      },
      ...
    ],
    tags               {
      fingerprint string,
      items       [
        string,
        ...
      ],
    },
  }
  project            string
  credential         reference to gauth_credential
end
```

#### Actions

* `create` -
  Converges the `gcompute_instance_template` resource into the final
  state described within the block. If the resource does not exist, Chef will
  attempt to create it.
* `delete` -
  Ensures the `gcompute_instance_template` resource is not present.
  If the resource already exists Chef will attempt to delete it.

#### Properties

* `creation_timestamp` -
  Output only. Creation timestamp in RFC3339 text format.

* `description` -
  An optional description of this resource. Provide this property when
  you create the resource.

* `id` -
  Output only. The unique identifier for the resource. This identifier
  is defined by the server.

* `name` -
  Required. Name of the resource. The name is 1-63 characters long
  and complies with RFC1035.

* `properties` -
  The instance properties for this instance template.

* `properties/can_ip_forward`
  Enables instances created based on this template to send packets
  with source IP addresses other than their own and receive packets
  with destination IP addresses other than their own. If these
  instances will be used as an IP gateway or it will be set as the
  next-hop in a Route resource, specify true. If unsure, leave this
  set to false.

* `properties/description`
  An optional text description for the instances that are created
  from this instance template.

* `properties/disks`
  An array of disks that are associated with the instances that are
  created from this template.

* `properties/disks[]/auto_delete`
  Specifies whether the disk will be auto-deleted when the
  instance is deleted (but not when the disk is detached from
  the instance).
  Tip: Disks should be set to autoDelete=true
  so that leftover disks are not left behind on machine
  deletion.

* `properties/disks[]/boot`
  Indicates that this is a boot disk. The virtual machine will
  use the first partition of the disk for its root filesystem.

* `properties/disks[]/device_name`
  Specifies a unique device name of your choice that is
  reflected into the /dev/disk/by-id/google-* tree of a Linux
  operating system running within the instance. This name can
  be used to reference the device for mounting, resizing, and
  so on, from within the instance.

* `properties/disks[]/disk_encryption_key`
  Encrypts or decrypts a disk using a customer-supplied
  encryption key.

* `properties/disks[]/disk_encryption_key/raw_key`
  Specifies a 256-bit customer-supplied encryption key,
  encoded in RFC 4648 base64 to either encrypt or decrypt
  this resource.

* `properties/disks[]/disk_encryption_key/rsa_encrypted_key`
  Specifies an RFC 4648 base64 encoded, RSA-wrapped
  2048-bit customer-supplied encryption key to either
  encrypt or decrypt this resource.

* `properties/disks[]/disk_encryption_key/sha256`
  Output only. The RFC 4648 base64 encoded SHA-256 hash of the
  customer-supplied encryption key that protects this
  resource.

* `properties/disks[]/index`
  Assigns a zero-based index to this disk, where 0 is
  reserved for the boot disk. For example, if you have many
  disks attached to an instance, each disk would have a
  unique index number. If not specified, the server will
  choose an appropriate value.

* `properties/disks[]/initialize_params`
  Specifies the parameters for a new disk that will be
  created alongside the new instance. Use initialization
  parameters to create boot disks or local SSDs attached to
  the new instance.

* `properties/disks[]/initialize_params/disk_name`
  Specifies the disk name. If not specified, the default
  is to use the name of the instance.

* `properties/disks[]/initialize_params/disk_size_gb`
  Specifies the size of the disk in base-2 GB.

* `properties/disks[]/initialize_params/disk_type`
  Reference to a gcompute_disk_type resource.
  Specifies the disk type to use to create the instance.
  If not specified, the default is pd-standard.

* `properties/disks[]/initialize_params/source_image`
  The source image to create this disk. When creating a
  new instance, one of initializeParams.sourceImage or
  disks.source is required.  To create a disk with one of
  the public operating system images, specify the image
  by its family name.

* `properties/disks[]/initialize_params/source_image_encryption_key`
  The customer-supplied encryption key of the source
  image. Required if the source image is protected by a
  customer-supplied encryption key.
  Instance templates do not store customer-supplied
  encryption keys, so you cannot create disks for
  instances in a managed instance group if the source
  images are encrypted with your own keys.

* `properties/disks[]/initialize_params/source_image_encryption_key/raw_key`
  Specifies a 256-bit customer-supplied encryption
  key, encoded in RFC 4648 base64 to either encrypt
  or decrypt this resource.

* `properties/disks[]/initialize_params/source_image_encryption_key/sha256`
  Output only. The RFC 4648 base64 encoded SHA-256 hash of the
  customer-supplied encryption key that protects this
  resource.

* `properties/disks[]/interface`
  Specifies the disk interface to use for attaching this
  disk, which is either SCSI or NVME. The default is SCSI.
  Persistent disks must always use SCSI and the request will
  fail if you attempt to attach a persistent disk in any
  other format than SCSI.

* `properties/disks[]/mode`
  The mode in which to attach this disk, either READ_WRITE or
  READ_ONLY. If not specified, the default is to attach the
  disk in READ_WRITE mode.

* `properties/disks[]/source`
  Reference to a gcompute_disk resource. When creating a new instance,
  one of initializeParams.sourceImage or disks.source is required.
  If desired, you can also attach existing non-root
  persistent disks using this property. This field is only
  applicable for persistent disks.
  Note that for InstanceTemplate, specify the disk name, not
  the URL for the disk.

* `properties/disks[]/type`
  Specifies the type of the disk, either SCRATCH or
  PERSISTENT. If not specified, the default is PERSISTENT.

* `properties/machine_type`
  Required. Reference to a gcompute_machine_type resource.

* `properties/metadata`
  The metadata key/value pairs to assign to instances that are
  created from this template. These pairs can consist of custom
  metadata or predefined keys.

* `properties/guest_accelerators`
  List of the type and count of accelerator cards attached to the
  instance

* `properties/guest_accelerators[]/accelerator_count`
  The number of the guest accelerator cards exposed to this
  instance.

* `properties/guest_accelerators[]/accelerator_type`
  Full or partial URL of the accelerator type resource to expose
  to this instance.

* `properties/network_interfaces`
  An array of configurations for this interface. This specifies
  how this interface is configured to interact with other
  network services, such as connecting to the internet. Only
  one network interface is supported per instance.

* `properties/network_interfaces[]/access_configs`
  An array of configurations for this interface. Currently, only
  one access config, ONE_TO_ONE_NAT, is supported. If there are no
  accessConfigs specified, then this instance will have no
  external internet access.

* `properties/network_interfaces[]/access_configs[]/name`
  Required. The name of this access configuration. The
  default and recommended name is External NAT but you can
  use any arbitrary string you would like. For example, My
  external IP or Network Access.

* `properties/network_interfaces[]/access_configs[]/nat_ip`
  Specifies the title of a gcompute_address.
  An external IP address associated with this instance.
  Specify an unused static external IP address available to
  the project or leave this field undefined to use an IP
  from a shared ephemeral IP address pool. If you specify a
  static external IP address, it must live in the same
  region as the zone of the instance.

* `properties/network_interfaces[]/access_configs[]/type`
  Required. The type of configuration. The default and only option is
  ONE_TO_ONE_NAT.

* `properties/network_interfaces[]/alias_ip_ranges`
  An array of alias IP ranges for this network interface. Can
  only be specified for network interfaces on subnet-mode
  networks.

* `properties/network_interfaces[]/alias_ip_ranges[]/ip_cidr_range`
  The IP CIDR range represented by this alias IP range.
  This IP CIDR range must belong to the specified
  subnetwork and cannot contain IP addresses reserved by
  system or used by other network interfaces. This range
  may be a single IP address (e.g. 10.2.3.4), a netmask
  (e.g. /24) or a CIDR format string (e.g. 10.1.2.0/24).

* `properties/network_interfaces[]/alias_ip_ranges[]/subnetwork_range_name`
  Optional subnetwork secondary range name specifying
  the secondary range from which to allocate the IP
  CIDR range for this alias IP range. If left
  unspecified, the primary range of the subnetwork will
  be used.

* `properties/network_interfaces[]/name`
  Output only. The name of the network interface, generated by the
  server. For network devices, these are eth0, eth1, etc

* `properties/network_interfaces[]/network`
  Specifies the title of an existing gcompute_network.  When creating
  an instance, if neither the network nor the subnetwork is specified,
  the default network global/networks/default is used; if the network
  is not specified but the subnetwork is specified, the network is
  inferred.

* `properties/network_interfaces[]/network_ip`
  An IPv4 internal network address to assign to the
  instance for this network interface. If not specified
  by the user, an unused internal IP is assigned by the
  system.

* `properties/network_interfaces[]/subnetwork`
  Reference to a gcompute_subnetwork resource.
  If the network resource is in legacy mode, do not
  provide this property.  If the network is in auto
  subnet mode, providing the subnetwork is optional. If
  the network is in custom subnet mode, then this field
  should be specified.

* `properties/scheduling`
  Sets the scheduling options for this instance.

* `properties/scheduling/automatic_restart`
  Specifies whether the instance should be automatically restarted
  if it is terminated by Compute Engine (not terminated by a user).
  You can only set the automatic restart option for standard
  instances. Preemptible instances cannot be automatically
  restarted.

* `properties/scheduling/on_host_maintenance`
  Defines the maintenance behavior for this instance. For standard
  instances, the default behavior is MIGRATE. For preemptible
  instances, the default and only possible behavior is TERMINATE.
  For more information, see Setting Instance Scheduling Options.

* `properties/scheduling/preemptible`
  Defines whether the instance is preemptible. This can only be set
  during instance creation, it cannot be set or changed after the
  instance has been created.

* `properties/service_accounts`
  A list of service accounts, with their specified scopes, authorized
  for this instance. Only one service account per VM instance is
  supported.

* `properties/service_accounts[]/email`
  Email address of the service account.

* `properties/service_accounts[]/scopes`
  The list of scopes to be made available for this service
  account.

* `properties/tags`
  A list of tags to apply to this instance. Tags are used to identify
  valid sources or targets for network firewalls and are specified by
  the client during instance creation. The tags can be later modified
  by the setTags method. Each tag within the list must comply with
  RFC1035.

* `properties/tags/fingerprint`
  Specifies a fingerprint for this request, which is essentially a
  hash of the metadata's contents and used for optimistic locking.
  The fingerprint is initially generated by Compute Engine and
  changes after every request to modify or update metadata. You
  must always provide an up-to-date fingerprint hash in order to
  update or change metadata.

* `properties/tags/items`
  An array of tags. Each tag must be 1-63 characters long, and
  comply with RFC1035.

#### Label
Set the `it_label` property when attempting to set primary key
of this object. The primary key will always be referred to by the initials of
the resource followed by "_label"

### gcompute_image
Represents an Image resource.

Google Compute Engine uses operating system images to create the root
persistent disks for your instances. You specify an image when you create
an instance. Images contain a boot loader, an operating system, and a
root file system. Linux operating system images are also capable of
running containers on Compute Engine.

Images can be either public or custom.

Public images are provided and maintained by Google, open-source
communities, and third-party vendors. By default, all projects have
access to these images and can use them to create instances.  Custom
images are available only to your project. You can create a custom image
from root persistent disks and other images. Then, use the custom image
to create an instance.


#### Example

#TODO

#### Reference

```ruby
gcompute_image 'id-for-resource' do
  archive_size_bytes         integer
  creation_timestamp         time
  deprecated                 {
    deleted     time,
    deprecated  time,
    obsolete    time,
    replacement string,
    state       'DEPRECATED', 'OBSOLETE' or 'DELETED',
  }
  description                string
  disk_size_gb               integer
  family                     string
  guest_os_features          [
    {
      type VIRTIO_SCSI_MULTIQUEUE,
    },
    ...
  ]
  id                         integer
  image_encryption_key       {
    raw_key string,
    sha256  string,
  }
  licenses                   [
    string,
    ...
  ]
  name                       string
  raw_disk                   {
    container_type TAR,
    sha1_checksum  string,
    source         string,
  }
  source_disk                reference to gcompute_disk
  source_disk_encryption_key {
    raw_key string,
    sha256  string,
  }
  source_disk_id             string
  source_type                RAW
  project                    string
  credential                 reference to gauth_credential
end
```

#### Actions

* `create` -
  Converges the `gcompute_image` resource into the final
  state described within the block. If the resource does not exist, Chef will
  attempt to create it.
* `delete` -
  Ensures the `gcompute_image` resource is not present.
  If the resource already exists Chef will attempt to delete it.

#### Properties

* `archive_size_bytes` -
  Output only. Size of the image tar.gz archive stored in Google Cloud Storage (in
  bytes).

* `creation_timestamp` -
  Output only. Creation timestamp in RFC3339 text format.

* `deprecated` -
  Output only. The deprecation status associated with this image.

* `deprecated/deleted`
  An optional RFC3339 timestamp on or after which the state of this
  resource is intended to change to DELETED. This is only
  informational and the status will not change unless the client
  explicitly changes it.

* `deprecated/deprecated`
  An optional RFC3339 timestamp on or after which the state of this
  resource is intended to change to DEPRECATED. This is only
  informational and the status will not change unless the client
  explicitly changes it.

* `deprecated/obsolete`
  An optional RFC3339 timestamp on or after which the state of this
  resource is intended to change to OBSOLETE. This is only
  informational and the status will not change unless the client
  explicitly changes it.

* `deprecated/replacement`
  The URL of the suggested replacement for a deprecated resource.
  The suggested replacement resource must be the same kind of
  resource as the deprecated resource.

* `deprecated/state`
  The deprecation state of this resource. This can be DEPRECATED,
  OBSOLETE, or DELETED. Operations which create a new resource
  using a DEPRECATED resource will return successfully, but with a
  warning indicating the deprecated resource and recommending its
  replacement. Operations which use OBSOLETE or DELETED resources
  will be rejected and result in an error.

* `description` -
  An optional description of this resource. Provide this property when
  you create the resource.

* `disk_size_gb` -
  Size of the image when restored onto a persistent disk (in GB).

* `family` -
  The name of the image family to which this image belongs. You can
  create disks by specifying an image family instead of a specific
  image name. The image family always returns its latest image that is
  not deprecated. The name of the image family must comply with
  RFC1035.

* `guest_os_features` -
  A list of features to enable on the guest OS. Applicable for
  bootable images only. Currently, only one feature can be enabled,
  VIRTIO_SCSI_MULTIQUEUE, which allows each virtual CPU to have its
  own queue. For Windows images, you can only enable
  VIRTIO_SCSI_MULTIQUEUE on images with driver version 1.2.0.1621 or
  higher. Linux images with kernel versions 3.17 and higher will
  support VIRTIO_SCSI_MULTIQUEUE.
  For new Windows images, the server might also populate this field
  with the value WINDOWS, to indicate that this is a Windows image.
  This value is purely informational and does not enable or disable
  any features.

* `guest_os_features[]/type`
  The type of supported feature. Currenty only
  VIRTIO_SCSI_MULTIQUEUE is supported. For newer Windows images,
  the server might also populate this property with the value
  WINDOWS to indicate that this is a Windows image. This value is
  purely informational and does not enable or disable any
  features.

* `id` -
  Output only. The unique identifier for the resource. This identifier
  is defined by the server.

* `image_encryption_key` -
  Encrypts the image using a customer-supplied encryption key.
  After you encrypt an image with a customer-supplied key, you must
  provide the same key if you use the image later (e.g. to create a
  disk from the image)

* `image_encryption_key/raw_key`
  Specifies a 256-bit customer-supplied encryption key, encoded in
  RFC 4648 base64 to either encrypt or decrypt this resource.

* `image_encryption_key/sha256`
  Output only. The RFC 4648 base64 encoded SHA-256 hash of the
  customer-supplied encryption key that protects this resource.

* `licenses` -
  Any applicable license URI.

* `name` -
  Required. Name of the resource; provided by the client when the resource is
  created. The name must be 1-63 characters long, and comply with
  RFC1035. Specifically, the name must be 1-63 characters long and
  match the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?` which means
  the first character must be a lowercase letter, and all following
  characters must be a dash, lowercase letter, or digit, except the
  last character, which cannot be a dash.

* `raw_disk` -
  The parameters of the raw disk image.

* `raw_disk/container_type`
  The format used to encode and transmit the block device, which
  should be TAR. This is just a container and transmission format
  and not a runtime format. Provided by the client when the disk
  image is created.

* `raw_disk/sha1_checksum`
  An optional SHA1 checksum of the disk image before unpackaging.
  This is provided by the client when the disk image is created.

* `raw_disk/source`
  The full Google Cloud Storage URL where disk storage is stored
  You must provide either this property or the sourceDisk property
  but not both.

* `source_disk` -
  Refers to a gcompute_disk object
  You must provide either this property or the
  rawDisk.source property but not both to create an image.

* `source_disk_encryption_key` -
  The customer-supplied encryption key of the source disk. Required if
  the source disk is protected by a customer-supplied encryption key.

* `source_disk_encryption_key/raw_key`
  Specifies a 256-bit customer-supplied encryption key, encoded in
  RFC 4648 base64 to either encrypt or decrypt this resource.

* `source_disk_encryption_key/sha256`
  Output only. The RFC 4648 base64 encoded SHA-256 hash of the
  customer-supplied encryption key that protects this resource.

* `source_disk_id` -
  The ID value of the disk used to create this image. This value may
  be used to determine whether the image was taken from the current
  or a previous instance of a given disk name.

* `source_type` -
  The type of the image used to create this disk. The default and
  only value is RAW

#### Label
Set the `i_label` property when attempting to set primary key
of this object. The primary key will always be referred to by the initials of
the resource followed by "_label"

### gcompute_instance
An instance is a virtual machine (VM) hosted on Google's infrastructure.


#### Example

#TODO

#### Reference

```ruby
gcompute_instance 'id-for-resource' do
  can_ip_forward     boolean
  cpu_platform       string
  creation_timestamp string
  disks              [
    {
      auto_delete         boolean,
      boot                boolean,
      device_name         string,
      disk_encryption_key {
        raw_key           string,
        rsa_encrypted_key string,
        sha256            string,
      },
      index               integer,
      initialize_params   {
        disk_name                   string,
        disk_size_gb                integer,
        disk_type                   reference to gcompute_disk_type,
        source_image                string,
        source_image_encryption_key {
          raw_key string,
          sha256  string,
        },
      },
      interface           'SCSI' or 'NVME',
      mode                'READ_WRITE' or 'READ_ONLY',
      source              reference to gcompute_disk,
      type                'SCRATCH' or 'PERSISTENT',
    },
    ...
  ]
  guest_accelerators [
    {
      accelerator_count integer,
      accelerator_type  string,
    },
    ...
  ]
  id                 integer
  label_fingerprint  string
  machine_type       reference to gcompute_machine_type
  metadata           namevalues
  min_cpu_platform   string
  name               string
  network_interfaces [
    {
      access_configs  [
        {
          name   string,
          nat_ip reference to gcompute_address,
          type   ONE_TO_ONE_NAT,
        },
        ...
      ],
      alias_ip_ranges [
        {
          ip_cidr_range         string,
          subnetwork_range_name string,
        },
        ...
      ],
      name            string,
      network         reference to gcompute_network,
      network_ip      string,
      subnetwork      reference to gcompute_subnetwork,
    },
    ...
  ]
  scheduling         {
    automatic_restart   boolean,
    on_host_maintenance string,
    preemptible         boolean,
  }
  service_accounts   [
    {
      email  string,
      scopes [
        string,
        ...
      ],
    },
    ...
  ]
  status             string
  status_message     string
  tags               {
    fingerprint string,
    items       [
      string,
      ...
    ],
  }
  zone               reference to gcompute_zone
  project            string
  credential         reference to gauth_credential
end
```

#### Actions

* `create` -
  Converges the `gcompute_instance` resource into the final
  state described within the block. If the resource does not exist, Chef will
  attempt to create it.
* `delete` -
  Ensures the `gcompute_instance` resource is not present.
  If the resource already exists Chef will attempt to delete it.

#### Properties

* `can_ip_forward` -
  Allows this instance to send and receive packets with non-matching
  destination or source IPs. This is required if you plan to use this
  instance to forward routes.

* `cpu_platform` -
  Output only. The CPU platform used by this instance.

* `creation_timestamp` -
  Output only. Creation timestamp in RFC3339 text format.

* `disks` -
  An array of disks that are associated with the instances that are
  created from this template.

* `disks[]/auto_delete`
  Specifies whether the disk will be auto-deleted when the
  instance is deleted (but not when the disk is detached from
  the instance).
  Tip: Disks should be set to autoDelete=true
  so that leftover disks are not left behind on machine
  deletion.

* `disks[]/boot`
  Indicates that this is a boot disk. The virtual machine will
  use the first partition of the disk for its root filesystem.

* `disks[]/device_name`
  Specifies a unique device name of your choice that is
  reflected into the /dev/disk/by-id/google-* tree of a Linux
  operating system running within the instance. This name can
  be used to reference the device for mounting, resizing, and
  so on, from within the instance.

* `disks[]/disk_encryption_key`
  Encrypts or decrypts a disk using a customer-supplied
  encryption key.

* `disks[]/disk_encryption_key/raw_key`
  Specifies a 256-bit customer-supplied encryption key,
  encoded in RFC 4648 base64 to either encrypt or decrypt
  this resource.

* `disks[]/disk_encryption_key/rsa_encrypted_key`
  Specifies an RFC 4648 base64 encoded, RSA-wrapped
  2048-bit customer-supplied encryption key to either
  encrypt or decrypt this resource.

* `disks[]/disk_encryption_key/sha256`
  Output only. The RFC 4648 base64 encoded SHA-256 hash of the
  customer-supplied encryption key that protects this
  resource.

* `disks[]/index`
  Assigns a zero-based index to this disk, where 0 is
  reserved for the boot disk. For example, if you have many
  disks attached to an instance, each disk would have a
  unique index number. If not specified, the server will
  choose an appropriate value.

* `disks[]/initialize_params`
  Specifies the parameters for a new disk that will be
  created alongside the new instance. Use initialization
  parameters to create boot disks or local SSDs attached to
  the new instance.

* `disks[]/initialize_params/disk_name`
  Specifies the disk name. If not specified, the default
  is to use the name of the instance.

* `disks[]/initialize_params/disk_size_gb`
  Specifies the size of the disk in base-2 GB.

* `disks[]/initialize_params/disk_type`
  Reference to a gcompute_disk_type resource.
  Specifies the disk type to use to create the instance.
  If not specified, the default is pd-standard.

* `disks[]/initialize_params/source_image`
  The source image to create this disk. When creating a
  new instance, one of initializeParams.sourceImage or
  disks.source is required.  To create a disk with one of
  the public operating system images, specify the image
  by its family name.

* `disks[]/initialize_params/source_image_encryption_key`
  The customer-supplied encryption key of the source
  image. Required if the source image is protected by a
  customer-supplied encryption key.
  Instance templates do not store customer-supplied
  encryption keys, so you cannot create disks for
  instances in a managed instance group if the source
  images are encrypted with your own keys.

* `disks[]/initialize_params/source_image_encryption_key/raw_key`
  Specifies a 256-bit customer-supplied encryption
  key, encoded in RFC 4648 base64 to either encrypt
  or decrypt this resource.

* `disks[]/initialize_params/source_image_encryption_key/sha256`
  Output only. The RFC 4648 base64 encoded SHA-256 hash of the
  customer-supplied encryption key that protects this
  resource.

* `disks[]/interface`
  Specifies the disk interface to use for attaching this
  disk, which is either SCSI or NVME. The default is SCSI.
  Persistent disks must always use SCSI and the request will
  fail if you attempt to attach a persistent disk in any
  other format than SCSI.

* `disks[]/mode`
  The mode in which to attach this disk, either READ_WRITE or
  READ_ONLY. If not specified, the default is to attach the
  disk in READ_WRITE mode.

* `disks[]/source`
  Reference to a gcompute_disk resource. When creating a new instance,
  one of initializeParams.sourceImage or disks.source is required.
  If desired, you can also attach existing non-root
  persistent disks using this property. This field is only
  applicable for persistent disks.

* `disks[]/type`
  Specifies the type of the disk, either SCRATCH or
  PERSISTENT. If not specified, the default is PERSISTENT.

* `guest_accelerators` -
  List of the type and count of accelerator cards attached to the
  instance

* `guest_accelerators[]/accelerator_count`
  The number of the guest accelerator cards exposed to this
  instance.

* `guest_accelerators[]/accelerator_type`
  Full or partial URL of the accelerator type resource to expose
  to this instance.

* `id` -
  Output only. The unique identifier for the resource. This identifier is defined by
  the server.

* `label_fingerprint` -
  A fingerprint for this request, which is essentially a hash of the
  metadata's contents and used for optimistic locking. The fingerprint
  is initially generated by Compute Engine and changes after every
  request to modify or update metadata. You must always provide an
  up-to-date fingerprint hash in order to update or change metadata.

* `metadata` -
  The metadata key/value pairs to assign to instances that are
  created from this template. These pairs can consist of custom
  metadata or predefined keys.

* `machine_type` -
  A reference to a machine type which defines VM kind.

* `min_cpu_platform` -
  Specifies a minimum CPU platform for the VM instance. Applicable
  values are the friendly names of CPU platforms

* `name` -
  The name of the resource, provided by the client when initially
  creating the resource. The resource name must be 1-63 characters long,
  and comply with RFC1035. Specifically, the name must be 1-63
  characters long and match the regular expression
  `[a-z]([-a-z0-9]*[a-z0-9])?` which means the first character must be a
  lowercase letter, and all following characters must be a dash,
  lowercase letter, or digit, except the last character, which cannot
  be a dash.

* `network_interfaces` -
  An array of configurations for this interface. This specifies
  how this interface is configured to interact with other
  network services, such as connecting to the internet. Only
  one network interface is supported per instance.

* `network_interfaces[]/access_configs`
  An array of configurations for this interface. Currently, only
  one access config, ONE_TO_ONE_NAT, is supported. If there are no
  accessConfigs specified, then this instance will have no
  external internet access.

* `network_interfaces[]/access_configs[]/name`
  Required. The name of this access configuration. The
  default and recommended name is External NAT but you can
  use any arbitrary string you would like. For example, My
  external IP or Network Access.

* `network_interfaces[]/access_configs[]/nat_ip`
  Specifies the title of a gcompute_address.
  An external IP address associated with this instance.
  Specify an unused static external IP address available to
  the project or leave this field undefined to use an IP
  from a shared ephemeral IP address pool. If you specify a
  static external IP address, it must live in the same
  region as the zone of the instance.

* `network_interfaces[]/access_configs[]/type`
  Required. The type of configuration. The default and only option is
  ONE_TO_ONE_NAT.

* `network_interfaces[]/alias_ip_ranges`
  An array of alias IP ranges for this network interface. Can
  only be specified for network interfaces on subnet-mode
  networks.

* `network_interfaces[]/alias_ip_ranges[]/ip_cidr_range`
  The IP CIDR range represented by this alias IP range.
  This IP CIDR range must belong to the specified
  subnetwork and cannot contain IP addresses reserved by
  system or used by other network interfaces. This range
  may be a single IP address (e.g. 10.2.3.4), a netmask
  (e.g. /24) or a CIDR format string (e.g. 10.1.2.0/24).

* `network_interfaces[]/alias_ip_ranges[]/subnetwork_range_name`
  Optional subnetwork secondary range name specifying
  the secondary range from which to allocate the IP
  CIDR range for this alias IP range. If left
  unspecified, the primary range of the subnetwork will
  be used.

* `network_interfaces[]/name`
  Output only. The name of the network interface, generated by the
  server. For network devices, these are eth0, eth1, etc

* `network_interfaces[]/network`
  Specifies the title of an existing gcompute_network.  When creating
  an instance, if neither the network nor the subnetwork is specified,
  the default network global/networks/default is used; if the network
  is not specified but the subnetwork is specified, the network is
  inferred.

* `network_interfaces[]/network_ip`
  An IPv4 internal network address to assign to the
  instance for this network interface. If not specified
  by the user, an unused internal IP is assigned by the
  system.

* `network_interfaces[]/subnetwork`
  Reference to a gcompute_subnetwork resource.
  If the network resource is in legacy mode, do not
  provide this property.  If the network is in auto
  subnet mode, providing the subnetwork is optional. If
  the network is in custom subnet mode, then this field
  should be specified.

* `scheduling` -
  Sets the scheduling options for this instance.

* `scheduling/automatic_restart`
  Specifies whether the instance should be automatically restarted
  if it is terminated by Compute Engine (not terminated by a user).
  You can only set the automatic restart option for standard
  instances. Preemptible instances cannot be automatically
  restarted.

* `scheduling/on_host_maintenance`
  Defines the maintenance behavior for this instance. For standard
  instances, the default behavior is MIGRATE. For preemptible
  instances, the default and only possible behavior is TERMINATE.
  For more information, see Setting Instance Scheduling Options.

* `scheduling/preemptible`
  Defines whether the instance is preemptible. This can only be set
  during instance creation, it cannot be set or changed after the
  instance has been created.

* `service_accounts` -
  A list of service accounts, with their specified scopes, authorized
  for this instance. Only one service account per VM instance is
  supported.

* `service_accounts[]/email`
  Email address of the service account.

* `service_accounts[]/scopes`
  The list of scopes to be made available for this service
  account.

* `status` -
  Output only. The status of the instance. One of the following values:
  PROVISIONING, STAGING, RUNNING, STOPPING, SUSPENDING, SUSPENDED,
  and TERMINATED.

* `status_message` -
  Output only. An optional, human-readable explanation of the status.

* `tags` -
  A list of tags to apply to this instance. Tags are used to identify
  valid sources or targets for network firewalls and are specified by
  the client during instance creation. The tags can be later modified
  by the setTags method. Each tag within the list must comply with
  RFC1035.

* `tags/fingerprint`
  Specifies a fingerprint for this request, which is essentially a
  hash of the metadata's contents and used for optimistic locking.
  The fingerprint is initially generated by Compute Engine and
  changes after every request to modify or update metadata. You
  must always provide an up-to-date fingerprint hash in order to
  update or change metadata.

* `tags/items`
  An array of tags. Each tag must be 1-63 characters long, and
  comply with RFC1035.

* `zone` -
  Required. A reference to the zone where the machine resides.

#### Label
Set the `i_label` property when attempting to set primary key
of this object. The primary key will always be referred to by the initials of
the resource followed by "_label"

### gcompute_instance_group
Represents an Instance Group resource. Instance groups are self-managed
and can contain identical or different instances. Instance groups do not
use an instance template. Unlike managed instance groups, you must create
and add instances to an instance group manually.


#### Example

#TODO

#### Reference

```ruby
gcompute_instance_group 'id-for-resource' do
  creation_timestamp time
  description        string
  id                 integer
  name               string
  named_ports        [
    {
      name string,
      port integer,
    },
    ...
  ]
  network            reference to gcompute_network
  region             reference to gcompute_region
  subnetwork         reference to gcompute_subnetwork
  zone               reference to gcompute_zone
  project            string
  credential         reference to gauth_credential
end
```

#### Actions

* `create` -
  Converges the `gcompute_instance_group` resource into the final
  state described within the block. If the resource does not exist, Chef will
  attempt to create it.
* `delete` -
  Ensures the `gcompute_instance_group` resource is not present.
  If the resource already exists Chef will attempt to delete it.

#### Properties

* `creation_timestamp` -
  Output only. Creation timestamp in RFC3339 text format.

* `description` -
  An optional description of this resource. Provide this property when
  you create the resource.

* `id` -
  Output only. A unique identifier for this instance group.

* `name` -
  The name of the instance group.
  The name must be 1-63 characters long, and comply with RFC1035.

* `named_ports` -
  Assigns a name to a port number.
  For example: {name: "http", port: 80}.
  This allows the system to reference ports by the assigned name
  instead of a port number. Named ports can also contain multiple
  ports.
  For example: [{name: "http", port: 80},{name: "http", port: 8080}]
  Named ports apply to all instances in this instance group.

* `named_ports[]/name`
  The name for this named port.
  The name must be 1-63 characters long, and comply with RFC1035.

* `named_ports[]/port`
  The port number, which can be a value between 1 and 65535.

* `network` -
  The network to which all instances in the instance group belong.

* `region` -
  The region where the instance group is located
  (for regional resources).

* `subnetwork` -
  The subnetwork to which all instances in the instance group belong.

* `zone` -
  Required. A reference to the zone where the instance group resides.

#### Label
Set the `ig_label` property when attempting to set primary key
of this object. The primary key will always be referred to by the initials of
the resource followed by "_label"

### gcompute_instance_group_manager
Creates a managed instance group using the information that you specify in
the request. After the group is created, it schedules an action to create
instances in the group using the specified instance template. This
operation is marked as DONE when the group is created even if the
instances in the group have not yet been created. You must separately
verify the status of the individual instances.

A managed instance group can have up to 1000 VM instances per group.


#### Example

#TODO

#### Reference

```ruby
gcompute_instance_group_manager 'id-for-resource' do
  base_instance_name string
  creation_timestamp time
  current_actions    {
    abandoning               integer,
    creating                 integer,
    creating_without_retries integer,
    deleting                 integer,
    none                     integer,
    recreating               integer,
    refreshing               integer,
    restarting               integer,
  }
  description        string
  id                 integer
  instance_group     reference to gcompute_instance_group
  instance_template  reference to gcompute_instance_template
  name               string
  named_ports        [
    {
      name string,
      port integer,
    },
    ...
  ]
  region             reference to gcompute_region
  target_pools       [
    reference to a gcompute_target_pool,
    ...
  ]
  target_size        integer
  zone               reference to gcompute_zone
  project            string
  credential         reference to gauth_credential
end
```

#### Actions

* `create` -
  Converges the `gcompute_instance_group_manager` resource into the final
  state described within the block. If the resource does not exist, Chef will
  attempt to create it.
* `delete` -
  Ensures the `gcompute_instance_group_manager` resource is not present.
  If the resource already exists Chef will attempt to delete it.

#### Properties

* `base_instance_name` -
  Required. The base instance name to use for instances in this group. The value
  must be 1-58 characters long. Instances are named by appending a
  hyphen and a random four-character string to the base instance name.
  The base instance name must comply with RFC1035.

* `creation_timestamp` -
  Output only. The creation timestamp for this managed instance group in RFC3339
  text format.

* `current_actions` -
  Output only. The list of instance actions and the number of instances in this
  managed instance group that are scheduled for each of those actions.

* `current_actions/abandoning`
  Output only. The total number of instances in the managed instance group that
  are scheduled to be abandoned. Abandoning an instance removes it
  from the managed instance group without deleting it.

* `current_actions/creating`
  Output only. The number of instances in the managed instance group that are
  scheduled to be created or are currently being created. If the
  group fails to create any of these instances, it tries again until
  it creates the instance successfully.
  If you have disabled creation retries, this field will not be
  populated; instead, the creatingWithoutRetries field will be
  populated.

* `current_actions/creating_without_retries`
  Output only. The number of instances that the managed instance group will
  attempt to create. The group attempts to create each instance only
  once. If the group fails to create any of these instances, it
  decreases the group's targetSize value accordingly.

* `current_actions/deleting`
  Output only. The number of instances in the managed instance group that are
  scheduled to be deleted or are currently being deleted.

* `current_actions/none`
  Output only. The number of instances in the managed instance group that are
  running and have no scheduled actions.

* `current_actions/recreating`
  Output only. The number of instances in the managed instance group that are
  scheduled to be recreated or are currently being being recreated.
  Recreating an instance deletes the existing root persistent disk
  and creates a new disk from the image that is defined in the
  instance template.

* `current_actions/refreshing`
  Output only. The number of instances in the managed instance group that are
  being reconfigured with properties that do not require a restart
  or a recreate action. For example, setting or removing target
  pools for the instance.

* `current_actions/restarting`
  Output only. The number of instances in the managed instance group that are
  scheduled to be restarted or are currently being restarted.

* `description` -
  An optional description of this resource. Provide this property when
  you create the resource.

* `id` -
  Output only. A unique identifier for this resource

* `instance_group` -
  Output only. The instance group being managed

* `instance_template` -
  Required. The instance template that is specified for this managed instance
  group. The group uses this template to create all new instances in the
  managed instance group.

* `name` -
  Required. The name of the managed instance group. The name must be 1-63
  characters long, and comply with RFC1035.

* `named_ports` -
  Named ports configured for the Instance Groups complementary to this Instance Group Manager.

* `named_ports[]/name`
  The name for this named port. The name must be 1-63 characters
  long, and comply with RFC1035.

* `named_ports[]/port`
  The port number, which can be a value between 1 and 65535.

* `region` -
  Output only. The region this managed instance group resides
  (for regional resources).

* `target_pools` -
  TargetPool resources to which instances in the instanceGroup field are
  added. The target pools automatically apply to all of the instances in
  the managed instance group.

* `target_size` -
  The target number of running instances for this managed instance
  group. Deleting or abandoning instances reduces this number. Resizing
  the group changes this number.

* `zone` -
  Required. The zone the managed instance group resides.

#### Label
Set the `igm_label` property when attempting to set primary key
of this object. The primary key will always be referred to by the initials of
the resource followed by "_label"

### gcompute_network
Represents a Network resource.

Your Cloud Platform Console project can contain multiple networks, and
each network can have multiple instances attached to it. A network allows
you to define a gateway IP and the network range for the instances
attached to that network. Every project is provided with a default network
with preset configurations and firewall rules. You can choose to customize
the default network by adding or removing rules, or you can create new
networks in that project. Generally, most users only need one network,
although you can have up to five networks per project by default.

A network belongs to only one project, and each instance can only belong
to one network. All Compute Engine networks use the IPv4 protocol. Compute
Engine currently does not support IPv6. However, Google is a major
advocate of IPv6 and it is an important future direction.

#### Reference Guides
* [API Reference](https://cloud.google.com/compute/docs/reference/rest/v1/networks)
* [Official Documentation](https://cloud.google.com/vpc/docs/vpc)

#### Example

#TODO

#### Reference

```ruby
gcompute_network 'id-for-resource' do
  auto_create_subnetworks boolean
  creation_timestamp      time
  description             string
  gateway_ipv4            string
  id                      integer
  ipv4_range              string
  name                    string
  routing_config          [
    {
      routing_mode 'REGIONAL' or 'GLOBAL',
    },
    ...
  ]
  subnetworks             [
    string,
    ...
  ]
  project                 string
  credential              reference to gauth_credential
end
```

#### Actions

* `create` -
  Converges the `gcompute_network` resource into the final
  state described within the block. If the resource does not exist, Chef will
  attempt to create it.
* `delete` -
  Ensures the `gcompute_network` resource is not present.
  If the resource already exists Chef will attempt to delete it.

#### Properties

* `description` -
  An optional description of this resource. Provide this property when
  you create the resource.

* `gateway_ipv4` -
  Output only. A gateway address for default routing to other networks. This value is
  read only and is selected by the Google Compute Engine, typically as
  the first usable address in the IPv4Range.

* `id` -
  Output only. The unique identifier for the resource.

* `ipv4_range` -
  The range of internal addresses that are legal on this network. This
  range is a CIDR specification, for example: 192.168.0.0/16. Provided
  by the client when the network is created.

* `name` -
  Required. Name of the resource. Provided by the client when the resource is
  created. The name must be 1-63 characters long, and comply with
  RFC1035. Specifically, the name must be 1-63 characters long and match
  the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?` which means the
  first character must be a lowercase letter, and all following
  characters must be a dash, lowercase letter, or digit, except the last
  character, which cannot be a dash.

* `subnetworks` -
  Output only. Server-defined fully-qualified URLs for all subnetworks in this
  network.

* `auto_create_subnetworks` -
  When set to true, the network is created in "auto subnet mode". When
  set to false, the network is in "custom subnet mode".
  In "auto subnet mode", a newly created network is assigned the default
  CIDR of 10.128.0.0/9 and it automatically creates one subnetwork per
  region.

* `creation_timestamp` -
  Output only. Creation timestamp in RFC3339 text format.

* `routing_config` -
  The network-level routing configuration for this network. Used by Cloud
  Router to determine what type of network-wide routing behavior to
  enforce.

* `routing_config[]/routing_mode`
  Required. The network-wide routing mode to use. If set to REGIONAL, this
  network's cloud routers will only advertise routes with subnetworks
  of this network in the same region as the router. If set to GLOBAL,
  this network's cloud routers will advertise routes with all
  subnetworks of this network, across regions.

#### Label
Set the `n_label` property when attempting to set primary key
of this object. The primary key will always be referred to by the initials of
the resource followed by "_label"

### gcompute_region_autoscaler
Represents an Autoscaler resource.

Autoscalers allow you to automatically scale virtual machine instances in
managed instance groups according to an autoscaling policy that you
define.

#### Reference Guides
* [API Reference](https://cloud.google.com/compute/docs/reference/rest/v1/regionAutoscalers)
* [Autoscaling Groups of Instances](https://cloud.google.com/compute/docs/autoscaler/)

#### Example

#TODO

#### Reference

```ruby
gcompute_region_autoscaler 'id-for-resource' do
  autoscaling_policy {
    cool_down_period_sec       integer,
    cpu_utilization            {
      utilization_target double,
    },
    custom_metric_utilizations [
      {
        metric                  string,
        utilization_target      double,
        utilization_target_type 'GAUGE', 'DELTA_PER_SECOND' or 'DELTA_PER_MINUTE',
      },
      ...
    ],
    load_balancing_utilization {
      utilization_target double,
    },
    max_num_replicas           integer,
    min_num_replicas           integer,
  }
  creation_timestamp time
  description        string
  id                 integer
  name               string
  region             reference to gcompute_region
  target             string
  project            string
  credential         reference to gauth_credential
end
```

#### Actions

* `create` -
  Converges the `gcompute_region_autoscaler` resource into the final
  state described within the block. If the resource does not exist, Chef will
  attempt to create it.
* `delete` -
  Ensures the `gcompute_region_autoscaler` resource is not present.
  If the resource already exists Chef will attempt to delete it.

#### Properties

* `id` -
  Output only. Unique identifier for the resource.

* `creation_timestamp` -
  Output only. Creation timestamp in RFC3339 text format.

* `name` -
  Required. Name of the resource. The name must be 1-63 characters long and match
  the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?` which means the
  first character must be a lowercase letter, and all following
  characters must be a dash, lowercase letter, or digit, except the last
  character, which cannot be a dash.

* `description` -
  An optional description of this resource.

* `autoscaling_policy` -
  Required. The configuration parameters for the autoscaling algorithm. You can
  define one or more of the policies for an autoscaler: cpuUtilization,
  customMetricUtilizations, and loadBalancingUtilization.
  If none of these are specified, the default will be to autoscale based
  on cpuUtilization to 0.6 or 60%.

* `autoscaling_policy/min_num_replicas`
  The minimum number of replicas that the autoscaler can scale down
  to. This cannot be less than 0. If not provided, autoscaler will
  choose a default value depending on maximum number of instances
  allowed.

* `autoscaling_policy/max_num_replicas`
  Required. The maximum number of instances that the autoscaler can scale up
  to. This is required when creating or updating an autoscaler. The
  maximum number of replicas should not be lower than minimal number
  of replicas.

* `autoscaling_policy/cool_down_period_sec`
  The number of seconds that the autoscaler should wait before it
  starts collecting information from a new instance. This prevents
  the autoscaler from collecting information when the instance is
  initializing, during which the collected usage would not be
  reliable. The default time autoscaler waits is 60 seconds.
  Virtual machine initialization times might vary because of
  numerous factors. We recommend that you test how long an
  instance may take to initialize. To do this, create an instance
  and time the startup process.

* `autoscaling_policy/cpu_utilization`
  Defines the CPU utilization policy that allows the autoscaler to
  scale based on the average CPU utilization of a managed instance
  group.

* `autoscaling_policy/cpu_utilization/utilization_target`
  The target CPU utilization that the autoscaler should maintain.
  Must be a float value in the range (0, 1]. If not specified, the
  default is 0.6.
  If the CPU level is below the target utilization, the autoscaler
  scales down the number of instances until it reaches the minimum
  number of instances you specified or until the average CPU of
  your instances reaches the target utilization.
  If the average CPU is above the target utilization, the autoscaler
  scales up until it reaches the maximum number of instances you
  specified or until the average utilization reaches the target
  utilization.

* `autoscaling_policy/custom_metric_utilizations`
  Defines the CPU utilization policy that allows the autoscaler to
  scale based on the average CPU utilization of a managed instance
  group.

* `autoscaling_policy/custom_metric_utilizations[]/metric`
  Required. The identifier (type) of the Stackdriver Monitoring metric.
  The metric cannot have negative values.
  The metric must have a value type of INT64 or DOUBLE.

* `autoscaling_policy/custom_metric_utilizations[]/utilization_target`
  Required. The target value of the metric that autoscaler should
  maintain. This must be a positive value. A utilization
  metric scales number of virtual machines handling requests
  to increase or decrease proportionally to the metric.
  For example, a good metric to use as a utilizationTarget is
  www.googleapis.com/compute/instance/network/received_bytes_count.
  The autoscaler will work to keep this value constant for each
  of the instances.

* `autoscaling_policy/custom_metric_utilizations[]/utilization_target_type`
  Required. Defines how target utilization value is expressed for a
  Stackdriver Monitoring metric. Either GAUGE, DELTA_PER_SECOND,
  or DELTA_PER_MINUTE.

* `autoscaling_policy/load_balancing_utilization`
  Configuration parameters of autoscaling based on a load balancer.

* `autoscaling_policy/load_balancing_utilization/utilization_target`
  Fraction of backend capacity utilization (set in HTTP(s) load
  balancing configuration) that autoscaler should maintain. Must
  be a positive float value. If not defined, the default is 0.8.

* `target` -
  Required. URL of the managed instance group that this autoscaler will scale.

* `region` -
  Required. URL of the region where the instance group resides.

#### Label
Set the `ra_label` property when attempting to set primary key
of this object. The primary key will always be referred to by the initials of
the resource followed by "_label"

### gcompute_route
Represents a Route resource.

A route is a rule that specifies how certain packets should be handled by
the virtual network. Routes are associated with virtual machines by tag,
and the set of routes for a particular virtual machine is called its
routing table. For each packet leaving a virtual machine, the system
searches that virtual machine's routing table for a single best matching
route.

Routes match packets by destination IP address, preferring smaller or more
specific ranges over larger ones. If there is a tie, the system selects
the route with the smallest priority value. If there is still a tie, it
uses the layer three and four packet headers to select just one of the
remaining matching routes. The packet is then forwarded as specified by
the next_hop field of the winning route -- either to another virtual
machine destination, a virtual machine gateway or a Compute
Engine-operated gateway. Packets that do not match any route in the
sending virtual machine's routing table will be dropped.

A Route resource must have exactly one specification of either
nextHopGateway, nextHopInstance, nextHopIp, or nextHopVpnTunnel.

#### Reference Guides
* [API Reference](https://cloud.google.com/compute/docs/reference/rest/v1/routes)
* [Using Routes](https://cloud.google.com/vpc/docs/using-routes)

#### Example

#TODO

#### Reference

```ruby
gcompute_route 'id-for-resource' do
  description         string
  dest_range          string
  name                string
  network             reference to gcompute_network
  next_hop_gateway    string
  next_hop_instance   string
  next_hop_ip         string
  next_hop_network    string
  next_hop_vpn_tunnel string
  priority            integer
  tags                [
    string,
    ...
  ]
  project             string
  credential          reference to gauth_credential
end
```

#### Actions

* `create` -
  Converges the `gcompute_route` resource into the final
  state described within the block. If the resource does not exist, Chef will
  attempt to create it.
* `delete` -
  Ensures the `gcompute_route` resource is not present.
  If the resource already exists Chef will attempt to delete it.

#### Properties

* `dest_range` -
  Required. The destination range of outgoing packets that this route applies to.
  Only IPv4 is supported.

* `description` -
  An optional description of this resource. Provide this property
  when you create the resource.

* `name` -
  Required. Name of the resource. Provided by the client when the resource is
  created. The name must be 1-63 characters long, and comply with
  RFC1035.  Specifically, the name must be 1-63 characters long and
  match the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?` which means
  the first character must be a lowercase letter, and all following
  characters must be a dash, lowercase letter, or digit, except the
  last character, which cannot be a dash.

* `network` -
  Required. The network that this route applies to.

* `priority` -
  The priority of this route. Priority is used to break ties in cases
  where there is more than one matching route of equal prefix length.
  In the case of two routes with equal prefix length, the one with the
  lowest-numbered priority value wins.
  Default value is 1000. Valid range is 0 through 65535.

* `tags` -
  A list of instance tags to which this route applies.

* `next_hop_gateway` -
  URL to a gateway that should handle matching packets.
  Currently, you can only specify the internet gateway, using a full or
  partial valid URL:
  * https://www.googleapis.com/compute/v1/projects/project/
  global/gateways/default-internet-gateway
  * projects/project/global/gateways/default-internet-gateway
  * global/gateways/default-internet-gateway

* `next_hop_instance` -
  URL to an instance that should handle matching packets.
  You can specify this as a full or partial URL. For example:
  * https://www.googleapis.com/compute/v1/projects/project/zones/zone/
  instances/instance
  * projects/project/zones/zone/instances/instance
  * zones/zone/instances/instance

* `next_hop_ip` -
  Network IP address of an instance that should handle matching packets.

* `next_hop_vpn_tunnel` -
  URL to a VpnTunnel that should handle matching packets.

* `next_hop_network` -
  Output only. URL to a Network that should handle matching packets.

#### Label
Set the `r_label` property when attempting to set primary key
of this object. The primary key will always be referred to by the initials of
the resource followed by "_label"

### gcompute_router
Represents a Router resource.

#### Reference Guides
* [API Reference](https://cloud.google.com/compute/docs/reference/rest/v1/routers)
* [Google Cloud Router](https://cloud.google.com/router/docs/)

#### Example

#TODO

#### Reference

```ruby
gcompute_router 'id-for-resource' do
  bgp                {
    advertise_mode       'DEFAULT' or 'CUSTOM',
    advertised_groups    [
      string,
      ...
    ],
    advertised_ip_ranges [
      {
        description string,
        range       string,
      },
      ...
    ],
    asn                  integer,
  }
  creation_timestamp time
  description        string
  id                 integer
  name               string
  network            reference to gcompute_network
  region             reference to gcompute_region
  project            string
  credential         reference to gauth_credential
end
```

#### Actions

* `create` -
  Converges the `gcompute_router` resource into the final
  state described within the block. If the resource does not exist, Chef will
  attempt to create it.
* `delete` -
  Ensures the `gcompute_router` resource is not present.
  If the resource already exists Chef will attempt to delete it.

#### Properties

* `id` -
  Output only. The unique identifier for the resource.

* `creation_timestamp` -
  Output only. Creation timestamp in RFC3339 text format.

* `name` -
  Required. Name of the resource. The name must be 1-63 characters long, and
  comply with RFC1035. Specifically, the name must be 1-63 characters
  long and match the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?`
  which means the first character must be a lowercase letter, and all
  following characters must be a dash, lowercase letter, or digit,
  except the last character, which cannot be a dash.

* `description` -
  An optional description of this resource.

* `network` -
  Required. A reference to the network to which this router belongs.

* `bgp` -
  BGP information specific to this router.

* `bgp/asn`
  Required. Local BGP Autonomous System Number (ASN). Must be an RFC6996
  private ASN, either 16-bit or 32-bit. The value will be fixed for
  this router resource. All VPN tunnels that link to this router
  will have the same local ASN.

* `bgp/advertise_mode`
  User-specified flag to indicate which mode to use for advertisement.
  Valid values of this enum field are: DEFAULT, CUSTOM

* `bgp/advertised_groups`
  User-specified list of prefix groups to advertise in custom mode.
  This field can only be populated if advertiseMode is CUSTOM and
  is advertised to all peers of the router. These groups will be
  advertised in addition to any specified prefixes. Leave this field
  blank to advertise no custom groups.
  This enum field has the one valid value: ALL_SUBNETS

* `bgp/advertised_ip_ranges`
  User-specified list of individual IP ranges to advertise in
  custom mode. This field can only be populated if advertiseMode
  is CUSTOM and is advertised to all peers of the router. These IP
  ranges will be advertised in addition to any specified groups.
  Leave this field blank to advertise no custom IP ranges.

* `bgp/advertised_ip_ranges[]/range`
  The IP range to advertise. The value must be a
  CIDR-formatted string.

* `bgp/advertised_ip_ranges[]/description`
  User-specified description for the IP range.

* `region` -
  Required. Region where the router resides.

#### Label
Set the `r_label` property when attempting to set primary key
of this object. The primary key will always be referred to by the initials of
the resource followed by "_label"

### gcompute_snapshot
Represents a Persistent Disk Snapshot resource.

Use snapshots to back up data from your persistent disks. Snapshots are
different from public images and custom images, which are used primarily
to create instances or configure instance templates. Snapshots are useful
for periodic backup of the data on your persistent disks. You can create
snapshots from persistent disks even while they are attached to running
instances.

Snapshots are incremental, so you can create regular snapshots on a
persistent disk faster and at a much lower cost than if you regularly
created a full image of the disk.


#### Example

#TODO

#### Reference

```ruby
gcompute_snapshot 'id-for-resource' do
  creation_timestamp         time
  description                string
  disk_size_gb               integer
  id                         integer
  labels                     [
    string,
    ...
  ]
  licenses                   [
    reference to a gcompute_license,
    ...
  ]
  name                       string
  snapshot_encryption_key    {
    raw_key string,
    sha256  string,
  }
  source                     reference to gcompute_disk
  source_disk_encryption_key {
    raw_key string,
    sha256  string,
  }
  storage_bytes              integer
  zone                       reference to gcompute_zone
  project                    string
  credential                 reference to gauth_credential
end
```

#### Actions

* `create` -
  Converges the `gcompute_snapshot` resource into the final
  state described within the block. If the resource does not exist, Chef will
  attempt to create it.
* `delete` -
  Ensures the `gcompute_snapshot` resource is not present.
  If the resource already exists Chef will attempt to delete it.

#### Properties

* `creation_timestamp` -
  Output only. Creation timestamp in RFC3339 text format.

* `id` -
  Output only. The unique identifier for the resource.

* `disk_size_gb` -
  Output only. Size of the snapshot, specified in GB.

* `name` -
  Required. Name of the resource; provided by the client when the resource is
  created. The name must be 1-63 characters long, and comply with
  RFC1035. Specifically, the name must be 1-63 characters long and match
  the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?` which means the
  first character must be a lowercase letter, and all following
  characters must be a dash, lowercase letter, or digit, except the last
  character, which cannot be a dash.

* `description` -
  An optional description of this resource.

* `storage_bytes` -
  Output only. A size of the the storage used by the snapshot. As snapshots share
  storage, this number is expected to change with snapshot
  creation/deletion.

* `licenses` -
  A list of public visible licenses that apply to this snapshot. This
  can be because the original image had licenses attached (such as a
  Windows image).  snapshotEncryptionKey nested object Encrypts the
  snapshot using a customer-supplied encryption key.

* `labels` -
  Labels to apply to this snapshot.

* `source` -
  A reference to the disk used to create this snapshot.

* `zone` -
  A reference to the zone where the disk is hosted.

* `snapshot_encryption_key` -
  The customer-supplied encryption key of the snapshot. Required if the
  source snapshot is protected by a customer-supplied encryption key.

* `snapshot_encryption_key/raw_key`
  Specifies a 256-bit customer-supplied encryption key, encoded in
  RFC 4648 base64 to either encrypt or decrypt this resource.

* `snapshot_encryption_key/sha256`
  Output only. The RFC 4648 base64 encoded SHA-256 hash of the customer-supplied
  encryption key that protects this resource.

* `source_disk_encryption_key` -
  The customer-supplied encryption key of the source snapshot. Required
  if the source snapshot is protected by a customer-supplied encryption
  key.

* `source_disk_encryption_key/raw_key`
  Specifies a 256-bit customer-supplied encryption key, encoded in
  RFC 4648 base64 to either encrypt or decrypt this resource.

* `source_disk_encryption_key/sha256`
  Output only. The RFC 4648 base64 encoded SHA-256 hash of the customer-supplied
  encryption key that protects this resource.

#### Label
Set the `s_label` property when attempting to set primary key
of this object. The primary key will always be referred to by the initials of
the resource followed by "_label"

### gcompute_ssl_certificate
An SslCertificate resource, used for HTTPS load balancing. This resource
provides a mechanism to upload an SSL key and certificate to
the load balancer to serve secure connections from the user.

#### Reference Guides
* [API Reference](https://cloud.google.com/compute/docs/reference/rest/v1/sslCertificates)
* [Official Documentation](https://cloud.google.com/load-balancing/docs/ssl-certificates)

#### Example

#TODO

#### Reference

```ruby
gcompute_ssl_certificate 'id-for-resource' do
  certificate        string
  creation_timestamp time
  description        string
  id                 integer
  name               string
  private_key        string
  project            string
  credential         reference to gauth_credential
end
```

#### Actions

* `create` -
  Converges the `gcompute_ssl_certificate` resource into the final
  state described within the block. If the resource does not exist, Chef will
  attempt to create it.
* `delete` -
  Ensures the `gcompute_ssl_certificate` resource is not present.
  If the resource already exists Chef will attempt to delete it.

#### Properties

* `certificate` -
  Required. The certificate in PEM format.
  The certificate chain must be no greater than 5 certs long.
  The chain must include at least one intermediate cert.

* `creation_timestamp` -
  Output only. Creation timestamp in RFC3339 text format.

* `description` -
  An optional description of this resource.

* `id` -
  Output only. The unique identifier for the resource.

* `name` -
  Name of the resource. Provided by the client when the resource is
  created. The name must be 1-63 characters long, and comply with
  RFC1035. Specifically, the name must be 1-63 characters long and match
  the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?` which means the
  first character must be a lowercase letter, and all following
  characters must be a dash, lowercase letter, or digit, except the last
  character, which cannot be a dash.

* `private_key` -
  Required. The write-only private key in PEM format.

#### Label
Set the `sc_label` property when attempting to set primary key
of this object. The primary key will always be referred to by the initials of
the resource followed by "_label"

### gcompute_ssl_policy
Represents a SSL policy. SSL policies give you the ability to control the
features of SSL that your SSL proxy or HTTPS load balancer negotiates.

#### Reference Guides
* [API Reference](https://cloud.google.com/compute/docs/reference/rest/v1/sslPolicies)
* [Using SSL Policies](https://cloud.google.com/compute/docs/load-balancing/ssl-policies)

#### Example

#TODO

#### Reference

```ruby
gcompute_ssl_policy 'id-for-resource' do
  creation_timestamp time
  custom_features    [
    string,
    ...
  ]
  description        string
  enabled_features   [
    string,
    ...
  ]
  fingerprint        string
  id                 integer
  min_tls_version    'TLS_1_0', 'TLS_1_1' or 'TLS_1_2'
  name               string
  profile            'COMPATIBLE', 'MODERN', 'RESTRICTED' or 'CUSTOM'
  warnings           [
    {
      code    string,
      message string,
    },
    ...
  ]
  project            string
  credential         reference to gauth_credential
end
```

#### Actions

* `create` -
  Converges the `gcompute_ssl_policy` resource into the final
  state described within the block. If the resource does not exist, Chef will
  attempt to create it.
* `delete` -
  Ensures the `gcompute_ssl_policy` resource is not present.
  If the resource already exists Chef will attempt to delete it.

#### Properties

* `creation_timestamp` -
  Output only. Creation timestamp in RFC3339 text format.

* `description` -
  An optional description of this resource.

* `id` -
  Output only. The unique identifier for the resource.

* `name` -
  Required. Name of the resource. Provided by the client when the resource is
  created. The name must be 1-63 characters long, and comply with
  RFC1035. Specifically, the name must be 1-63 characters long and match
  the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?` which means the
  first character must be a lowercase letter, and all following
  characters must be a dash, lowercase letter, or digit, except the last
  character, which cannot be a dash.

* `profile` -
  Profile specifies the set of SSL features that can be used by the
  load balancer when negotiating SSL with clients. This can be one of
  `COMPATIBLE`, `MODERN`, `RESTRICTED`, or `CUSTOM`. If using `CUSTOM`,
  the set of SSL features to enable must be specified in the
  `customFeatures` field.

* `min_tls_version` -
  The minimum version of SSL protocol that can be used by the clients
  to establish a connection with the load balancer. This can be one of
  `TLS_1_0`, `TLS_1_1`, `TLS_1_2`.

* `enabled_features` -
  Output only. The list of features enabled in the SSL policy.

* `custom_features` -
  A list of features enabled when the selected profile is CUSTOM. The
  method returns the set of features that can be specified in this
  list. This field must be empty if the profile is not CUSTOM.

* `fingerprint` -
  Output only. Fingerprint of this resource. A hash of the contents stored in this
  object. This field is used in optimistic locking.

* `warnings` -
  Output only. If potential misconfigurations are detected for this SSL policy, this
  field will be populated with warning messages.

* `warnings[]/code`
  Output only. A warning code, if applicable.

* `warnings[]/message`
  Output only. A human-readable description of the warning code.

#### Label
Set the `sp_label` property when attempting to set primary key
of this object. The primary key will always be referred to by the initials of
the resource followed by "_label"

### gcompute_subnetwork
A VPC network is a virtual version of the traditional physical networks
that exist within and between physical data centers. A VPC network
provides connectivity for your Compute Engine virtual machine (VM)
instances, Container Engine containers, App Engine Flex services, and
other network-related resources.

Each GCP project contains one or more VPC networks. Each VPC network is a
global entity spanning all GCP regions. This global VPC network allows VM
instances and other resources to communicate with each other via internal,
private IP addresses.

Each VPC network is subdivided into subnets, and each subnet is contained
within a single region. You can have more than one subnet in a region for
a given VPC network. Each subnet has a contiguous private RFC1918 IP
space. You create instances, containers, and the like in these subnets.
When you create an instance, you must create it in a subnet, and the
instance draws its internal IP address from that subnet.

Virtual machine (VM) instances in a VPC network can communicate with
instances in all other subnets of the same VPC network, regardless of
region, using their RFC1918 private IP addresses. You can isolate portions
of the network, even entire subnets, using firewall rules.

#### Reference Guides
* [API Reference](https://cloud.google.com/compute/docs/reference/rest/beta/subnetworks)
* [Private Google Access](https://cloud.google.com/vpc/docs/configure-private-google-access)
* [Cloud Networking](https://cloud.google.com/vpc/docs/using-vpc)

#### Example

#TODO

#### Reference

```ruby
gcompute_subnetwork 'id-for-resource' do
  creation_timestamp       time
  description              string
  enable_flow_logs         boolean
  fingerprint              fingerprint
  gateway_address          string
  id                       integer
  ip_cidr_range            string
  name                     string
  network                  reference to gcompute_network
  private_ip_google_access boolean
  region                   reference to gcompute_region
  secondary_ip_ranges      [
    {
      ip_cidr_range string,
      range_name    string,
    },
    ...
  ]
  project                  string
  credential               reference to gauth_credential
end
```

#### Actions

* `create` -
  Converges the `gcompute_subnetwork` resource into the final
  state described within the block. If the resource does not exist, Chef will
  attempt to create it.
* `delete` -
  Ensures the `gcompute_subnetwork` resource is not present.
  If the resource already exists Chef will attempt to delete it.

#### Properties

* `creation_timestamp` -
  Output only. Creation timestamp in RFC3339 text format.

* `description` -
  An optional description of this resource. Provide this property when
  you create the resource. This field can be set only at resource
  creation time.

* `gateway_address` -
  Output only. The gateway address for default routes to reach destination addresses
  outside this subnetwork.

* `id` -
  Output only. The unique identifier for the resource.

* `ip_cidr_range` -
  Required. The range of internal addresses that are owned by this subnetwork.
  Provide this property when you create the subnetwork. For example,
  10.0.0.0/8 or 192.168.0.0/16. Ranges must be unique and
  non-overlapping within a network. Only IPv4 is supported.

* `name` -
  Required. The name of the resource, provided by the client when initially
  creating the resource. The name must be 1-63 characters long, and
  comply with RFC1035. Specifically, the name must be 1-63 characters
  long and match the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?` which
  means the first character must be a lowercase letter, and all
  following characters must be a dash, lowercase letter, or digit,
  except the last character, which cannot be a dash.

* `network` -
  Required. The network this subnet belongs to.
  Only networks that are in the distributed mode can have subnetworks.

* `enable_flow_logs` -
  Whether to enable flow logging for this subnetwork.

* `fingerprint` -
  Output only. Fingerprint of this resource. This field is used internally during
  updates of this resource.

* `secondary_ip_ranges` -
  An array of configurations for secondary IP ranges for VM instances
  contained in this subnetwork. The primary IP of such VM must belong
  to the primary ipCidrRange of the subnetwork. The alias IPs may belong
  to either primary or secondary ranges.

* `secondary_ip_ranges[]/range_name`
  Required. The name associated with this subnetwork secondary range, used
  when adding an alias IP range to a VM instance. The name must
  be 1-63 characters long, and comply with RFC1035. The name
  must be unique within the subnetwork.

* `secondary_ip_ranges[]/ip_cidr_range`
  Required. The range of IP addresses belonging to this subnetwork secondary
  range. Provide this property when you create the subnetwork.
  Ranges must be unique and non-overlapping with all primary and
  secondary IP ranges within a network. Only IPv4 is supported.

* `private_ip_google_access` -
  Whether the VMs in this subnet can access Google services without
  assigned external IP addresses.

* `region` -
  Required. URL of the GCP region for this subnetwork.

#### Label
Set the `s_label` property when attempting to set primary key
of this object. The primary key will always be referred to by the initials of
the resource followed by "_label"

### gcompute_target_http_proxy
Represents a TargetHttpProxy resource, which is used by one or more global
forwarding rule to route incoming HTTP requests to a URL map.

#### Reference Guides
* [API Reference](https://cloud.google.com/compute/docs/reference/latest/targetHttpProxies)
* [Official Documentation](https://cloud.google.com/compute/docs/load-balancing/http/target-proxies)

#### Example

#TODO

#### Reference

```ruby
gcompute_target_http_proxy 'id-for-resource' do
  creation_timestamp time
  description        string
  id                 integer
  name               string
  url_map            reference to gcompute_url_map
  project            string
  credential         reference to gauth_credential
end
```

#### Actions

* `create` -
  Converges the `gcompute_target_http_proxy` resource into the final
  state described within the block. If the resource does not exist, Chef will
  attempt to create it.
* `delete` -
  Ensures the `gcompute_target_http_proxy` resource is not present.
  If the resource already exists Chef will attempt to delete it.

#### Properties

* `creation_timestamp` -
  Output only. Creation timestamp in RFC3339 text format.

* `description` -
  An optional description of this resource.

* `id` -
  Output only. The unique identifier for the resource.

* `name` -
  Required. Name of the resource. Provided by the client when the resource is
  created. The name must be 1-63 characters long, and comply with
  RFC1035. Specifically, the name must be 1-63 characters long and match
  the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?` which means the
  first character must be a lowercase letter, and all following
  characters must be a dash, lowercase letter, or digit, except the last
  character, which cannot be a dash.

* `url_map` -
  Required. A reference to the UrlMap resource that defines the mapping from URL
  to the BackendService.

#### Label
Set the `thp_label` property when attempting to set primary key
of this object. The primary key will always be referred to by the initials of
the resource followed by "_label"

### gcompute_target_https_proxy
Represents a TargetHttpsProxy resource, which is used by one or more
global forwarding rule to route incoming HTTPS requests to a URL map.

#### Reference Guides
* [API Reference](https://cloud.google.com/compute/docs/reference/latest/targetHttpsProxies)
* [Official Documentation](https://cloud.google.com/compute/docs/load-balancing/http/target-proxies)

#### Example

#TODO

#### Reference

```ruby
gcompute_target_https_proxy 'id-for-resource' do
  creation_timestamp time
  description        string
  id                 integer
  name               string
  quic_override      'NONE', 'ENABLE' or 'DISABLE'
  ssl_certificates   [
    reference to a gcompute_ssl_certificate,
    ...
  ]
  url_map            reference to gcompute_url_map
  project            string
  credential         reference to gauth_credential
end
```

#### Actions

* `create` -
  Converges the `gcompute_target_https_proxy` resource into the final
  state described within the block. If the resource does not exist, Chef will
  attempt to create it.
* `delete` -
  Ensures the `gcompute_target_https_proxy` resource is not present.
  If the resource already exists Chef will attempt to delete it.

#### Properties

* `creation_timestamp` -
  Output only. Creation timestamp in RFC3339 text format.

* `description` -
  An optional description of this resource.

* `id` -
  Output only. The unique identifier for the resource.

* `name` -
  Required. Name of the resource. Provided by the client when the resource is
  created. The name must be 1-63 characters long, and comply with
  RFC1035. Specifically, the name must be 1-63 characters long and match
  the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?` which means the
  first character must be a lowercase letter, and all following
  characters must be a dash, lowercase letter, or digit, except the last
  character, which cannot be a dash.

* `quic_override` -
  Specifies the QUIC override policy for this resource. This determines
  whether the load balancer will attempt to negotiate QUIC with clients
  or not. Can specify one of NONE, ENABLE, or DISABLE. If NONE is
  specified, uses the QUIC policy with no user overrides, which is
  equivalent to DISABLE. Not specifying this field is equivalent to
  specifying NONE.

* `ssl_certificates` -
  Required. A list of SslCertificate resources that are used to authenticate
  connections between users and the load balancer. Currently, exactly
  one SSL certificate must be specified.

* `url_map` -
  Required. A reference to the UrlMap resource that defines the mapping from URL
  to the BackendService.

#### Label
Set the `thp_label` property when attempting to set primary key
of this object. The primary key will always be referred to by the initials of
the resource followed by "_label"

### gcompute_target_pool
Represents a TargetPool resource, used for Load Balancing.
#### Reference Guides
* [API Reference](https://cloud.google.com/compute/docs/reference/rest/v1/targetPools)
* [Official Documentation](https://cloud.google.com/compute/docs/load-balancing/network/target-pools)

#### Example

#TODO

#### Reference

```ruby
gcompute_target_pool 'id-for-resource' do
  backup_pool        reference to gcompute_target_pool
  creation_timestamp time
  description        string
  failover_ratio     double
  health_check       reference to gcompute_http_health_check
  id                 integer
  instances          [
    reference to a gcompute_instance,
    ...
  ]
  name               string
  region             reference to gcompute_region
  session_affinity   'NONE', 'CLIENT_IP' or 'CLIENT_IP_PROTO'
  project            string
  credential         reference to gauth_credential
end
```

#### Actions

* `create` -
  Converges the `gcompute_target_pool` resource into the final
  state described within the block. If the resource does not exist, Chef will
  attempt to create it.
* `delete` -
  Ensures the `gcompute_target_pool` resource is not present.
  If the resource already exists Chef will attempt to delete it.

#### Properties

* `backup_pool` -
  This field is applicable only when the containing target pool is
  serving a forwarding rule as the primary pool, and its failoverRatio
  field is properly set to a value between [0, 1].
  backupPool and failoverRatio together define the fallback behavior of
  the primary target pool: if the ratio of the healthy instances in the
  primary pool is at or below failoverRatio, traffic arriving at the
  load-balanced IP will be directed to the backup pool.
  In case where failoverRatio and backupPool are not set, or all the
  instances in the backup pool are unhealthy, the traffic will be
  directed back to the primary pool in the "force" mode, where traffic
  will be spread to the healthy instances with the best effort, or to
  all instances when no instance is healthy.

* `creation_timestamp` -
  Output only. Creation timestamp in RFC3339 text format.

* `description` -
  An optional description of this resource.

* `failover_ratio` -
  This field is applicable only when the containing target pool is
  serving a forwarding rule as the primary pool (i.e., not as a backup
  pool to some other target pool). The value of the field must be in
  [0, 1].
  If set, backupPool must also be set. They together define the fallback
  behavior of the primary target pool: if the ratio of the healthy
  instances in the primary pool is at or below this number, traffic
  arriving at the load-balanced IP will be directed to the backup pool.
  In case where failoverRatio is not set or all the instances in the
  backup pool are unhealthy, the traffic will be directed back to the
  primary pool in the "force" mode, where traffic will be spread to the
  healthy instances with the best effort, or to all instances when no
  instance is healthy.

* `health_check` -
  A reference to a HttpHealthCheck resource.
  A member instance in this pool is considered healthy if and only if
  the health checks pass. If not specified it means all member instances
  will be considered healthy at all times.

* `id` -
  Output only. The unique identifier for the resource.

* `instances` -
  A list of virtual machine instances serving this pool.
  They must live in zones contained in the same region as this pool.

* `name` -
  Required. Name of the resource. Provided by the client when the resource is
  created. The name must be 1-63 characters long, and comply with
  RFC1035. Specifically, the name must be 1-63 characters long and match
  the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?` which means the
  first character must be a lowercase letter, and all following
  characters must be a dash, lowercase letter, or digit, except the last
  character, which cannot be a dash.

* `session_affinity` -
  Session affinity option. Must be one of these values:
  - NONE: Connections from the same client IP may go to any instance in
  the pool.
  - CLIENT_IP: Connections from the same client IP will go to the same
  instance in the pool while that instance remains healthy.
  - CLIENT_IP_PROTO: Connections from the same client IP with the same
  IP protocol will go to the same instance in the pool while that
  instance remains healthy.

* `region` -
  Required. The region where the target pool resides.

#### Label
Set the `tp_label` property when attempting to set primary key
of this object. The primary key will always be referred to by the initials of
the resource followed by "_label"

### gcompute_target_ssl_proxy
Represents a TargetSslProxy resource, which is used by one or more
global forwarding rule to route incoming SSL requests to a backend
service.

#### Reference Guides
* [API Reference](https://cloud.google.com/compute/docs/reference/latest/targetSslProxies)
* [Setting Up SSL proxy for Google Cloud Load Balancing](https://cloud.google.com/compute/docs/load-balancing/tcp-ssl/)

#### Example

#TODO

#### Reference

```ruby
gcompute_target_ssl_proxy 'id-for-resource' do
  creation_timestamp time
  description        string
  id                 integer
  name               string
  proxy_header       'NONE' or 'PROXY_V1'
  service            reference to gcompute_backend_service
  ssl_certificates   [
    reference to a gcompute_ssl_certificate,
    ...
  ]
  project            string
  credential         reference to gauth_credential
end
```

#### Actions

* `create` -
  Converges the `gcompute_target_ssl_proxy` resource into the final
  state described within the block. If the resource does not exist, Chef will
  attempt to create it.
* `delete` -
  Ensures the `gcompute_target_ssl_proxy` resource is not present.
  If the resource already exists Chef will attempt to delete it.

#### Properties

* `creation_timestamp` -
  Output only. Creation timestamp in RFC3339 text format.

* `description` -
  An optional description of this resource.

* `id` -
  Output only. The unique identifier for the resource.

* `name` -
  Required. Name of the resource. Provided by the client when the resource is
  created. The name must be 1-63 characters long, and comply with
  RFC1035. Specifically, the name must be 1-63 characters long and match
  the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?` which means the
  first character must be a lowercase letter, and all following
  characters must be a dash, lowercase letter, or digit, except the last
  character, which cannot be a dash.

* `proxy_header` -
  Specifies the type of proxy header to append before sending data to
  the backend, either NONE or PROXY_V1. The default is NONE.

* `service` -
  Required. A reference to the BackendService resource.

* `ssl_certificates` -
  Required. A list of SslCertificate resources that are used to authenticate
  connections between users and the load balancer. Currently, exactly
  one SSL certificate must be specified.

#### Label
Set the `tsp_label` property when attempting to set primary key
of this object. The primary key will always be referred to by the initials of
the resource followed by "_label"

### gcompute_target_tcp_proxy
Represents a TargetTcpProxy resource, which is used by one or more
global forwarding rule to route incoming TCP requests to a Backend
service.

#### Reference Guides
* [API Reference](https://cloud.google.com/compute/docs/reference/latest/targetTcpProxies)
* [Setting Up TCP proxy for Google Cloud Load Balancing](https://cloud.google.com/compute/docs/load-balancing/tcp-ssl/tcp-proxy)

#### Example

#TODO

#### Reference

```ruby
gcompute_target_tcp_proxy 'id-for-resource' do
  creation_timestamp time
  description        string
  id                 integer
  name               string
  proxy_header       'NONE' or 'PROXY_V1'
  service            reference to gcompute_backend_service
  project            string
  credential         reference to gauth_credential
end
```

#### Actions

* `create` -
  Converges the `gcompute_target_tcp_proxy` resource into the final
  state described within the block. If the resource does not exist, Chef will
  attempt to create it.
* `delete` -
  Ensures the `gcompute_target_tcp_proxy` resource is not present.
  If the resource already exists Chef will attempt to delete it.

#### Properties

* `creation_timestamp` -
  Output only. Creation timestamp in RFC3339 text format.

* `description` -
  An optional description of this resource.

* `id` -
  Output only. The unique identifier for the resource.

* `name` -
  Required. Name of the resource. Provided by the client when the resource is
  created. The name must be 1-63 characters long, and comply with
  RFC1035. Specifically, the name must be 1-63 characters long and match
  the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?` which means the
  first character must be a lowercase letter, and all following
  characters must be a dash, lowercase letter, or digit, except the last
  character, which cannot be a dash.

* `proxy_header` -
  Specifies the type of proxy header to append before sending data to
  the backend, either NONE or PROXY_V1. The default is NONE.

* `service` -
  Required. A reference to the BackendService resource.

#### Label
Set the `ttp_label` property when attempting to set primary key
of this object. The primary key will always be referred to by the initials of
the resource followed by "_label"

### gcompute_target_vpn_gateway
Represents a VPN gateway running in GCP. This virtual device is managed
by Google, but used only by you.

#### Reference Guides
* [API Reference](https://cloud.google.com/compute/docs/reference/rest/v1/targetVpnGateways)

#### Example

#TODO

#### Reference

```ruby
gcompute_target_vpn_gateway 'id-for-resource' do
  creation_timestamp time
  description        string
  forwarding_rules   [
    reference to a gcompute_forwarding_rule,
    ...
  ]
  id                 integer
  name               string
  network            reference to gcompute_network
  region             reference to gcompute_region
  tunnels            [
    string,
    ...
  ]
  project            string
  credential         reference to gauth_credential
end
```

#### Actions

* `create` -
  Converges the `gcompute_target_vpn_gateway` resource into the final
  state described within the block. If the resource does not exist, Chef will
  attempt to create it.
* `delete` -
  Ensures the `gcompute_target_vpn_gateway` resource is not present.
  If the resource already exists Chef will attempt to delete it.

#### Properties

* `creation_timestamp` -
  Output only. Creation timestamp in RFC3339 text format.

* `description` -
  An optional description of this resource.

* `name` -
  Required. Name of the resource. Provided by the client when the resource is
  created. The name must be 1-63 characters long, and comply with
  RFC1035.  Specifically, the name must be 1-63 characters long and
  match the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?` which means
  the first character must be a lowercase letter, and all following
  characters must be a dash, lowercase letter, or digit, except the last
  character, which cannot be a dash.

* `id` -
  Output only. The unique identifier for the resource.

* `network` -
  Required. The network this VPN gateway is accepting traffic for.

* `tunnels` -
  Output only. A list of references to VpnTunnel resources associated to this VPN gateway.

* `forwarding_rules` -
  Output only. A list of references to the ForwardingRule resources associated to this VPN
  gateway.

* `region` -
  Required. The region this gateway should sit in.

#### Label
Set the `tvg_label` property when attempting to set primary key
of this object. The primary key will always be referred to by the initials of
the resource followed by "_label"

### gcompute_url_map
UrlMaps are used to route requests to a backend service based on rules
that you define for the host and path of an incoming URL.


#### Example

#TODO

#### Reference

```ruby
gcompute_url_map 'id-for-resource' do
  creation_timestamp time
  default_service    reference to gcompute_backend_service
  description        string
  fingerprint        fingerprint
  host_rules         [
    {
      description  string,
      hosts        [
        string,
        ...
      ],
      path_matcher string,
    },
    ...
  ]
  id                 integer
  name               string
  path_matchers      [
    {
      default_service reference to gcompute_backend_service,
      description     string,
      name            string,
      path_rules      [
        {
          paths   [
            string,
            ...
          ],
          service reference to gcompute_backend_service,
        },
        ...
      ],
    },
    ...
  ]
  tests              [
    {
      description string,
      host        string,
      path        string,
      service     reference to gcompute_backend_service,
    },
    ...
  ]
  project            string
  credential         reference to gauth_credential
end
```

#### Actions

* `create` -
  Converges the `gcompute_url_map` resource into the final
  state described within the block. If the resource does not exist, Chef will
  attempt to create it.
* `delete` -
  Ensures the `gcompute_url_map` resource is not present.
  If the resource already exists Chef will attempt to delete it.

#### Properties

* `creation_timestamp` -
  Output only. Creation timestamp in RFC3339 text format.

* `default_service` -
  Required. A reference to BackendService resource if none of the hostRules match.

* `description` -
  An optional description of this resource. Provide this property when
  you create the resource.

* `host_rules` -
  The list of HostRules to use against the URL.

* `host_rules[]/description`
  An optional description of this resource. Provide this property
  when you create the resource.

* `host_rules[]/hosts`
  Required. The list of host patterns to match. They must be valid
  hostnames, except * will match any string of ([a-z0-9-.]*). In
  that case, * must be the first character and must be followed in
  the pattern by either - or ..

* `host_rules[]/path_matcher`
  Required. The name of the PathMatcher to use to match the path portion of
  the URL if the hostRule matches the URL's host portion.

* `id` -
  Output only. The unique identifier for the resource.

* `fingerprint` -
  Output only. Fingerprint of this resource. This field is used internally during
  updates of this resource.

* `name` -
  Required. Name of the resource. Provided by the client when the resource is
  created. The name must be 1-63 characters long, and comply with
  RFC1035. Specifically, the name must be 1-63 characters long and match
  the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?` which means the
  first character must be a lowercase letter, and all following
  characters must be a dash, lowercase letter, or digit, except the last
  character, which cannot be a dash.

* `path_matchers` -
  The list of named PathMatchers to use against the URL.

* `path_matchers[]/default_service`
  Required. A reference to a BackendService resource. This will be used if
  none of the pathRules defined by this PathMatcher is matched by
  the URL's path portion.

* `path_matchers[]/description`
  An optional description of this resource.

* `path_matchers[]/name`
  Required. The name to which this PathMatcher is referred by the HostRule.

* `path_matchers[]/path_rules`
  The list of path rules.

* `path_matchers[]/path_rules[]/paths`
  The list of path patterns to match. Each must start with /
  and the only place a * is allowed is at the end following
  a /. The string fed to the path matcher does not include
  any text after the first ? or #, and those chars are not
  allowed here.

* `path_matchers[]/path_rules[]/service`
  Required. A reference to the BackendService resource if this rule is
  matched.

* `tests` -
  The list of expected URL mappings. Request to update this UrlMap will
  succeed only if all of the test cases pass.

* `tests[]/description`
  Description of this test case.

* `tests[]/host`
  Required. Host portion of the URL.

* `tests[]/path`
  Required. Path portion of the URL.

* `tests[]/service`
  Required. A reference to expected BackendService resource the given URL should be mapped to.

#### Label
Set the `um_label` property when attempting to set primary key
of this object. The primary key will always be referred to by the initials of
the resource followed by "_label"

### gcompute_vpn_tunnel
VPN tunnel resource.
#### Reference Guides
* [API Reference](https://cloud.google.com/compute/docs/reference/rest/v1/vpnTunnels)
* [Cloud VPN Overview](https://cloud.google.com/vpn/docs/concepts/overview)
* [Networks and Tunnel Routing](https://cloud.google.com/vpn/docs/concepts/choosing-networks-routing)

#### Example

#TODO

#### Reference

```ruby
gcompute_vpn_tunnel 'id-for-resource' do
  creation_timestamp      time
  description             string
  ike_version             integer
  label_fingerprint       fingerprint
  labels                  namevalues
  local_traffic_selector  [
    string,
    ...
  ]
  name                    string
  peer_ip                 string
  region                  reference to gcompute_region
  remote_traffic_selector [
    string,
    ...
  ]
  router                  reference to gcompute_router
  shared_secret           string
  shared_secret_hash      string
  target_vpn_gateway      reference to gcompute_target_vpn_gateway
  project                 string
  credential              reference to gauth_credential
end
```

#### Actions

* `create` -
  Converges the `gcompute_vpn_tunnel` resource into the final
  state described within the block. If the resource does not exist, Chef will
  attempt to create it.
* `delete` -
  Ensures the `gcompute_vpn_tunnel` resource is not present.
  If the resource already exists Chef will attempt to delete it.

#### Properties

* `creation_timestamp` -
  Output only. Creation timestamp in RFC3339 text format.

* `name` -
  Required. Name of the resource. The name must be 1-63 characters long, and
  comply with RFC1035. Specifically, the name must be 1-63
  characters long and match the regular expression
  `[a-z]([-a-z0-9]*[a-z0-9])?` which means the first character
  must be a lowercase letter, and all following characters must
  be a dash, lowercase letter, or digit,
  except the last character, which cannot be a dash.

* `description` -
  An optional description of this resource.

* `target_vpn_gateway` -
  Required. URL of the Target VPN gateway with which this VPN tunnel is
  associated.

* `router` -
  URL of router resource to be used for dynamic routing.

* `peer_ip` -
  Required. IP address of the peer VPN gateway. Only IPv4 is supported.

* `shared_secret` -
  Required. Shared secret used to set the secure session between the Cloud VPN
  gateway and the peer VPN gateway.

* `shared_secret_hash` -
  Output only. Hash of the shared secret.

* `ike_version` -
  IKE protocol version to use when establishing the VPN tunnel with
  peer VPN gateway.
  Acceptable IKE versions are 1 or 2. Default version is 2.

* `local_traffic_selector` -
  Local traffic selector to use when establishing the VPN tunnel with
  peer VPN gateway. The value should be a CIDR formatted string,
  for example `192.168.0.0/16`. The ranges should be disjoint.
  Only IPv4 is supported.

* `remote_traffic_selector` -
  Remote traffic selector to use when establishing the VPN tunnel with
  peer VPN gateway. The value should be a CIDR formatted string,
  for example `192.168.0.0/16`. The ranges should be disjoint.
  Only IPv4 is supported.

* `labels` -
  Labels to apply to this VpnTunnel.

* `label_fingerprint` -
  Output only. The fingerprint used for optimistic locking of this resource.  Used
  internally during updates.

* `region` -
  Required. The region where the tunnel is located.

#### Label
Set the `vt_label` property when attempting to set primary key
of this object. The primary key will always be referred to by the initials of
the resource followed by "_label"

[google-gauth]: https://supermarket.chef.io/cookbooks/google-gauth
