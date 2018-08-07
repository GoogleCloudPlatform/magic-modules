# products/{{product}}/chef.yaml

The `chef.yaml` file contains Chef specific settings, overrides and other
customizations to build the Chef cookbook for the product.

A `chef.yaml` file should derive from `Provider::Chef::Config` object.

> Please note that you may find Ruby code inlined in the `chef.yaml`
> throughout this doc and in the products/*/chef.yaml. The main reason for
> that are 2 fold:
>
> 1. Chef is in Ruby. The code listed will be placed into the final product
>    generated, so it *has* to be in Ruby. If your provider is targeting
>    something on another language, e.g. Go or Python, those scripts will likely
>    be written in Go or Python respectively.
> 2. We decided to inline the Ruby functions because they were usually small and
>    managing many small files. That was an arbitrary decision, but localized to
>    the Chef provider. If you write another provider you are free to do in a
>    different way if you wish.
>
> You will still have to deal with a minor amount of Ruby, mostly
> [ERB (Embedded RuBy)][erb-home]. However it will be very minimal and
> restricted to helping you build the code in your target language. For example
> you can use ERB to iterate through objects, properties and other niceties
> provided by the compiler. For example:
>
> ```
> <%
>   object.properties.each do |prop|
>     ... your code (or template) iterated with every property
>   end
> -%>
> ```

Example:

    --- !ruby/object:Provider::Chef::Config
    manifest: !ruby/object:Provider::Chef::Manifest
      version: '0.1.0'
      source: 'https://github.com/GoogleCloudPlatform/chef-google-compute'
      ...
      ...


## Skeleton

    --- !ruby/object:Provider::Chef::Config
    manifest: !ruby/object:Provider::Chef::Manifest
      version:
      source:
      issues:
      summary:
      description:
      depends:
        - !ruby/object:Provider::Config::Requirements
        - ...
      operating_systems:
        - !ruby/object:Provider::Config::OperatingSystem
        - ...
    objects: !ruby/object:Api::Resource::HashArray
      <ObjectName>:
        update:
        provider_helpers:
          include:
            - ...
        overrides:
          <api_name>: <new_api_name>
    examples: !ruby/object:Api::Resource::HashArray
      <ObjectName>:
        - <object_name>.rb
        - delete_<object_name>.rb
        - ...
    files: !ruby/object:Provider::Config::Files
      copy:
        <target_file>: <source_file>
        ...
      compile:
        <target_file>: <source_file>
        ...
    test_data: !ruby/object:Provider::Config::TestData
      network: !ruby/object:Api::Resource::HashArray
        <ObjectName>:
          - <file>
          - ...
    style:
      - !ruby/object:Provider::Config::StyleException
        name:
        pinpoints:
          - function:
            exceptions:
              - ...
          - class:
            exceptions:
              - ...
          - cookbook:
            exceptions:
              - ...
          - test: matrix > ...
            exceptions:
              - ...
          - ...

## Template Files

  - provider/chef/common\~operating_systems.yaml
  - provider/chef/common\~copy.yaml
  - provider/chef/common\~compile\~before.yaml
  - provider/chef/common\~compile\~after.yaml

When creating examples you can inlcude the credential block that defines, and
explains how to use the credentials:

  - templates/chef/example\~auth.rb.erb

See [Compute instance.rb example][compute-instance-example] as reference.

## Features

TODO(nelsonjr): Get a list of features here and how they relate to the settings.


## Types


### `Api::Resource::HashArray`

A hash array is the same as a hash with the exception that the key should match
an existing object. If the object is not defined in the `api.yaml` the compiler
will fail with error.

This is to ensure deleted objects are not left behind and litters the
`chef.yaml` files.


## Settings


### `manifest` (`Provider::Chef::Manifest`)

Every Chef cookbook contains a manifest, which provides the basic information
for the cookbook. This information is parsed by Chef (and Chef Supermarket) to
reason about the cookbook (name, version, etc).

Example:

    manifest: !ruby/object:Provider::Chef::Manifest
      version: '0.1.0'
      ...

#### `version` (string)

The version of the cookbook, in [SemVer][semver-home] format. Example:

    version: '0.1.0'

#### `source` (URL)

The URL of the source code for the cookbook. Example:

    source: 'https://github.com/GoogleCloudPlatform/chef-google-compute'

#### `issues` (URL)

The URL of the issues tracker for the cookbook. Example:

    homepage: 'https://github.com/GoogleCloudPlatform/chef-google-compute/issues'

#### `summary` (string)

A short description of the cookbook. Example:

    summary: 'A Chef cookbook to manage Google Compute Engine resources'

#### `description` (string)

A longer description of the cookbook. Example:

    description: |
      This cookbook provides the built-in types and services for Chef
      to manage Google Cloud DNS resources, as native Chef types.

#### `requires` (list(`Provider::Config::Requirements`))

A list of cookbooks that this product cookbook depends on. Example:

    requires:
      - !ruby/object:Provider::Config::Requirements
        name: 'google-gauth'
        versions: '< 0.2.0'

> Note: `google-gauth` should always be listed in the required cookbooks.

#### `operating_systems` (list(`Provider::Config::OperatingSystem`))

The list of operating systems (and versions) that this cookbook fully supports.
Example:

    operating_systems:
      - !ruby/object:Provider::Config::OperatingSystem
        name: RedHat
        versions:
          - '6'
          - '7'

Note that we have a file that can be included with the common operating systems
we run Chef cookbooks on. Then the definition should be:

      operating_systems:
    <%= indent(include('provider/chef/common~operating_systems.yaml'), 4) %>


### `objects` (`Api::Resource::HashArray`)

Provides the root of customizations for a resource defined in the `api.yaml`.

Example ([Compute][compute-chef] `Disk` cannot be changed via API updates):

    objects: !ruby/object:Api::Resource::HashArray
      Disk:
        flush: raise 'Disk cannot be edited'

#### `objects/<object>/create` (code)

Creates a custom `create` function. Chef uses `create` when it detects the
object does not exist and user specified `ensure => present` in the manifest.

Example: [products/compute/chef.yaml][dns-chef]

#### `objects/<object>/delete` (code)

Creates a custom `delete` function. Chef uses `delete` when it detects the
object exists and user specified `ensure => absent` in the manifest.

Example: [products/compute/chef.yaml][dns-chef]

#### `objects/<object>/flush` (code)

Creates a custom `flush` function. Chef uses `flush` to update values for the
resources if they mismatch with the catalog.

Example: [products/compute/chef.yaml][dns-chef]

> Note that `flush` is called after `create` or `delete`. Usually you should not
> execute `flush` when `create` or `delete` is called. Those functions by
> default assign true to a boolean to `@created` and `@deleted` respectively. So
> a guard like this is in the default `flush` function:
>
> `return if @created || @deleted || !@dirty`

#### `objects/<object>/resource_to_request_patch` (code)

Code provided in this section will be added to `resource_to_request` method.
This allows performing object specific changes to the request before it is send
to the API. This is useful when there are special changes to the object being
submitted aside from the properties defined in the `api.yaml` file.

Example ([Compute][compute-chef] `BackendService` requires that we provide the
original `fingerprint` value to avoid concurrent changes since our last read):

    objects: !ruby/object:Api::Resource::HashArray
      BackendService:
        resource_to_request_patch: |
          unless @fetched.nil?
            # Convert to pure JSON
            request = JSON.parse(request.to_json)
            request['fingerprint'] = @fetched['fingerprint']
          end

#### `objects/<object>/resource_to_query` (code)

A function that is used to filter objects returned from the API. When the API
does not provide a `GET` operation to fetch the object (e.g. fetching a VM by
name and zone), but instead it provides a `LIST` operation that returns all
objects in a group (e.g. all DNS records for a zone) we need to filter the
results to find the object user is determine to operate on.

Example: [products/compute/chef.yaml][dns-chef]

#### `objects/<object>/provider_helpers`

Defines provider specific changes to the provider.

#### `objects/<object>/provider_helpers/visible`

Emits or suppresses specific functions from being generated in the object. This
is useful when the provider needs to have a special 'snowflake' version of the
function because the API is not straightforward.

Functions that can be overriden:

  - unwrap_resource
  - resource_to_request
  - return_if_object

Example:

    objects: !ruby/object:Api::Resource::HashArray
      ResourceRecordSet:
        provider_helpers:
          visible:
            unwrap_resource: false

##### `objects/<object>/provider_helpers/include` (list(string))

Appends the files listed into the provider generated file

Example ([Compute][compute-chef] includes a file with helper functions to be
used by the `Disk` object):

    objects: !ruby/object:Api::Resource::HashArray
      Disk:
        provider_helpers:
          include:
            - 'products/compute/helpers/provider_disk_type.rb'

#### `objects/<object>/overrides` (hash)

Maps a field to another name. This is useful when the field name conflicts with
a reserved platform specific keyword. For example the keyword `deprecated` is
reserved on Chef and cannot be used whereas GCE MachineType has a deprecated
field.

Example:

    objects: !ruby/object:Api::Resource::HashArray
      MachineType:
        overrides:
          deprecated: _deprecated

### `tests` (`Api::Resource::HashArray`)

Specifies extra tests to be added to the [rspec][rspec-home] unit tests.

This will display code that should override auto-generated testing code. View
the [DNS tests][dns-chef-test] for the most comprehensive view of altered tests.

### `examples` (`Api::Resource::HashArray`)

Lists the file names of each example included for each object. All examples
must be located in a product's files/ directory with the prefix
'examples\~cookbook\~'. Each file will be included in the cookbook's recipes
folder with the prefix 'examples\~'. The convention is to include a
<resource>.rb file and a delete_<resource>.rb file

Example:

    examples: !ruby/object:Api::Resource::HashArray
      Instance:
        - instance.rb # refers to files/examples~cookbook~instance.rb
        - delete_instance.rb


### `files` (`Provider::Config::Files`)

Lists files to be copied or compiled into the final product. We can add files on
a per-product basis at our whim. If the file is simply copied use the `copy`
section. If you need to execute code to build the file use the `compile`
section instead.

Example:

    files: !ruby/object:Provider::Config::Files
      copy:
        libraries/google/copy_file.rb: copy_file.rb
      compile:
        libraries/google/object_store.rb: google/object_store.rb
        libraries/chef/functions/gcompute_image_family.rb:
          products/compute/files/function~image_family.rb.erb

### `files/copy` (`Hash`)

Lists files that should be copied directly to the cookbook without any
compilation. List the location in the cookbook as the key (relative to the cookbook
root) and the location in the code generator (relative to Magic Module root) as
the value.

### `files/compile` (`Hash`)

Lists files that should be compiled and then copied to the cookbook. List the
location in the cookbook as the key (relative to the cookbook root) and the
location in the code generator (relative to Magic Module root) as the value.

### `test_data` (`Provider::Config::TestData`)

Data files that will be used to generate canned responses from the API. During
tests we trap all API calls and return canned data to avoid the need to hit GCP
with a real request.

This is useful to both make the tests faster and simple, as well as avoid
charges and operating setup to developers helping us improve the cookbooks.

Example:

    test_data:
      network: !ruby/object:Api::Resource::HashArray
        ResourceRecordSet:
          - create~name
          - create~title

May be of type Provider::Config::TestData::NONE with a reason. This means that
all test data will be automatically generated.

Example:

    test_data: !ruby/object:Provider::Config::TestData::NONE
      reason: 'Test Data will be automatically generated'

#### `test_data/network` (`Api::Resource::HashArray`)

A list of objects that have manually written network data files for unit tests.

#### `test_data/network/<object>` (list(`String`))

A list of test data files to be included with the cookbook. These files will be
taken from files/ and are all prefixed with 'spec\~'. They will be copied to
spec/data/network/<object out_name>/


### `style` (list(`Provider::Config::StyleException`))

Applies [rubocop][rubocop-home] exemptions to generated filed. This is useful
when the generated file violates a Rubocop rule and we want a one-off pass
instead of disabling the rule.

For example, if there are too many properties in an object, the
resource_to_request will have too many lines for rubocop taste.

In that case we can whitelist the resource_to_request to be exempt without the
need to allow any other method to have the same leniency.

Example:

    - !ruby/object:Provider::Config::StyleException
      name: lib/google/compute/property/instance_disks.rb
      pinpoints:
        - function: resource_to_request
          exceptions:
            - Metrics/MethodLength

### `style[]/name` (`String`)

The filename where the exception is needed.

### `style[]/pinpoints` (list(`Hash`))

List of places in the file where exceptions are needed.

### `style[]/pinpoints[]/function` (`String`)

The name of the function where exceptions are needed. Prefix with class name
of multiple versions of the function exist and only one needs exceptions

### `style[]/pinpoints[]/class` (`String`)

The name of the class where exceptions are needed. Prefix with cookbook chain
if multiple classes of the same name exist and only one needs exceptions.

### `style[]/pinpoints[]/exceptions` (list(`String`))

A list of necessary Rubocop exceptions.

[compute-chef]: ../products/compute/chef.yaml
[dns-chef]: ../products/dns/chef.yaml
[dns-chef-test]: ../products/dns/test.yaml
[compute-instance-example]: ../products/compute/files/examples~cookbook~instance.rb
[semver-home]: http://semver.org/
[rspec-home]: http://rspec.info/
[rubocop-home]: https://github.com/bbatsov/rubocop
