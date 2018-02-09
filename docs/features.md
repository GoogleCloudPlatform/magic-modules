<div id='incomplete'
    style='background-color: #ffcccc; padding: 8pt;'>
  This document is still under construction.
  TODO(nelsonjr): Document this.
</div>

## Architecture

![Architecture Diagram][architecture]

## Features

### Semi-strong typed configuration files

`product/api.yaml` and `product/provider.yaml` are parsed, compiled and
verified by a custom serializer to reduce errors in the provider configuration.
The compiler enables many features seen in web templating languages, such as
partial inclusion.

### Multiple Target Support

Magic Module currently has support for two providers: Chef and Puppet. Each
provider shares the same API definitions, unit test data and as much of the
underlying logic as possible. Each provider has the ability to specify
provider-specific templates, provider-specific overrides (see below) and
provider-specific examples. Code is shared between providers whenever possible
to ensure that all modules work as expected and to keep code DRY.

### Templated Sources

Each module is built as a series of templates that has data injected into it.
This ensures that all modules look uniform and allows us to make sure that
features end up in all modules simultaneously.

### Template Overrides

 Since no Google Cloud API is the same, we need the ability to be flexible.
 The templates have the ability to override entire functions (create, delete,
 edit) or smaller values like an identifying URL.

### Flexible Output Mapping

Warning: TBD

Allows specifying the final module look-and-feel by using compile and copy rules
...

### Non-REST Compliant API Support

#### Objects Created by other means than via "POST path/to/collection"
[e.g. DNS: RecordSet vs. Changes]

Warning: TBD

#### Objects changed by special "one off" feature URL
[e.g. Compute: Network switchToCustomMode]

Warning: TBD

### Field Mapping
[e.g. DNS: "rrdatas" is too cryptic]

Warning: TBD
Map rrdatas => target.

### Fields not confirming to camelCase notation
[e.g. Compute: Network IPv4Range and gatewayIPv4]

Warning: TBD
Enables renaming the field, just use "field"...

### Support for read-only objects

Warning: TBD

### Support for exposing inner properties
[e.g. DNS: Project quota.<prop>]

Warning: TBD

### Support for asynchronous operations

Warning: TBD


### Include / Compile any file, including YAML and JSON

This flexibility allows us to have any common configuration be easily
centralized:

      operating_systems:
      <%= indent(include('provider/puppet~common~operating_systems.yaml'), 4) %>

Instead of being verbosely specified (and harder to maintain) on each file:

    operating_systems:
      - !ruby/object:Provider::Puppet::OperatingSystem
        name: RedHat
        versions:
          - '6'
          - '7'
      - !ruby/object:Provider::Puppet::OperatingSystem
        name: CentOS
        versions:
          - '6'
          - '7'
      ...

### Support for specific deliverables

+  Examples
+  Tests
+  Files
+  Custom code

### Support for "computed properties"

_Not yet ready!_
Warning:TBD

This is required to simplify some quirks in our API, such as if you want to
grant a privilege to a user, the "account" is actually a concatenation of the
word "user" + a dash + "email". Better if we can tease them apart so we can
check for typos, e.g. if you type "user**r**-nelsona@google.com" it would not
be caught, but if you have: "type => 'userr', principal => 'nelsona@google.com'"
we can do a better validation job.

### Support for custom validation functions

_Not yet ready!_
Warning: TBD

Some properties are special and should be crafted in a specific way. Allow to
provide a "validation function" that would fail the execution if it returns
false. This is a complement of the computed properties above, in case the user
still wants to use the direct form, maybe for his own backwards compatibility.

### Input only properties

Warning: TBD

### Output only properties

Warning: TBD

### Virtual resources

Warning: TBD

### Resource Refs

Warning: TBD

A ResourceRef is a type that allows the user to specify a short version of the
object, as it is already specified elsewhere in their manifest. For example,
when creating a VM you need: disk, network/subnetwork, etc. Each one of them
require you to have their URI. Sure you can try to build them (assuming they
are static and easy to do), but it looks like this:

    google_compute_instance { 'vm':
        ensure             => present,
        network_interfaces => [
          {
            network        => join([
                'projects', 'my-project',
                'global', 'networks', 'default',
              ], '/'),
            subnetwork     => join([
                'projects', 'my-project',
                'regions', 'us-central1',
                'subnetworks', 'default',
              ], '/'),
          },
        ],
        disks              => [
          {
            type              => 'PERSISTENT',
            boot              => true,
            mode              => 'READ_WRITE',
            auto_delete       => true,
            device_name       => 'vm',
            initialize_params => {
              source_image => 'my-image',
              disk_type    => join([
                  'projects', 'my-project,
                  'zones', 'us-central1-a',
                  'diskTypes', 'pd-standard',
                ], '/'),
              disk_size_gb => '10',
            }
          }
        ],

To avoid this we introduced ResourceRefs, which makes the code look like this:

    gcompute_network { 'production':
      ensure => present,
      ...
    }

    gcompute_subnetwork { 'servers':
      ensure  => present,
      network => 'production',
    }

    gcompute_disk { 'web-server-disk1':
      ensure     => present,
      diskSizeGb => 100,
      ...
    }

    gcompute_instance { 'myvm':
      disks => [
        'web-server-disk1',
      ],
      network_interfaces => [
        {
          network    => 'production',
          subnetwork => 'servers'
        },
      ],
    }

This is not much simpler to read, but we can also inspect the disk and _ensure_
it is up to the specifications the user care about.

### Composite (Nested) Objects

Warning: TBD

#### Puppet Example

    gcontainer_cluster { 'foo':
      node_config    => {
        machine_type => 'n1-standard-1',
        diskSizeGb   => 100,
      },
      instanceGroupUrls => [
        'aaaaaaa',
        'bbbbbbb',
        'ccccccc',
      ],
      anArrayOfObjects  => [
        { prop1 => 'aaaa',
          prop2 => 10 },
        { prop1 => 'bbbb',
          prop2 => 20 },
      ],
    }

### Debugging

To enable verbose or debug message set one or more of the following environment
variables when executing the compiler, tests or the provider:

| Environment Variable    | Usage                                 |
| ----------------------- | ------------------------------------- |
| `PUPPET_HTTP_VERBOSE=1` | for REST API call debugging in-depth  |
| `PUPPET_HTTP_DEBUG=1`   | for API call list                     |
| `CHEF_HTTP_*`           | for API debugging like Puppet         |
| `RSPEC_DEBUG=1`         | for debugging tests                   |

Examples:

    $ RSPEC_DEBUG=1 rspec

    $ PUPPET_HTTP_VERBOSE=1 puppet apply examples/managed_zone.pp


[architecture]: architecture.png
