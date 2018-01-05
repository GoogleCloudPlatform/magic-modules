# products/{{product}}/puppet.yaml

The `puppet.yaml` file contains Puppet specific settings, overrides and other
customizations to build the Puppet module for the product.

A `puppet.yaml` file should derive from `Provider::Puppet::Config` object.

> Please note that you may find Ruby code inlined in the `puppet.yaml`
> throughout this doc and in the products/*/puppet.yaml. The main reason for
> that are 2 fold:
>
> 1. Puppet is in Ruby. The code listed will be placed into the final product
>    generated, so it *has* to be in Ruby. If your provider is targeting
>    something on another language, e.g. Go or Python, those scripts will likely
>    be written in Go or Python respectively.
> 2. We decided to inline the Ruby functions because they were usually small and
>    managing many small files. That was an arbitrary decision, but localized to
>    the Puppet provider. If you write another provider you are free to do in a
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

    --- !ruby/object:Provider::Puppet::Config
    manifest: !ruby/object:Provider::Puppet::Manifest
      version: '0.1.0'
      source: 'https://github.com/GoogleCloudPlatform/puppet-google-compute'
      ...
      ...


## Skeleton

    --- !ruby/object:Provider::Puppet::Config
    manifest: !ruby/object:Provider::Puppet::Manifest
      version:
      source:
      homepage:
      issues:
      summary:
      tags:
        - ...
      requires:
        - !ruby/object:Provider::Config::Requirements
        - ...
      operating_systems:
        - !ruby/object:Provider::Config::OperatingSystem
        - ...
    objects: !ruby/object:Api::Resource::HashArray
      <ObjectName>:
        create:
        delete:
        flush:
        resource_to_request_patch:
        resource_to_query:
        provider_helpers:
          visible:
            unwrap_resource:
            resource_to_request:
            return_if_object:
          include:
            - <file>
            - ...
    functions:
      - !ruby/object:Provider::Puppet::Function
        name:
        description:
        arguments:
          - !ruby/object:Provider::Puppet::Function::Argument
            name:
            type:
            description:
          - ...
        examples:
          - ...
        notes:
      - ...
    examples: !ruby/object:Api::Resource::HashArray
      <ObjectName>:
        - <object_name>.pp
        - delete_<object_name>.pp
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
          - module:
            exceptions:
              - ...
          - test: matrix > ...
            exceptions:
              - ...
          - ...

## Template Files

  - provider/puppet/common\~operating_systems.yaml
  - provider/puppet/common\~copy.yaml
  - provider/puppet/common\~compile\~before.yaml
  - provider/puppet/common\~compile\~after.yaml

When creating examples you can inlcude the credential block that defines, and
explains how to use the credentials:

  - templates/puppet/examples\~credential.pp.erb

See [Compute instance.pp example][compute-instance-example] as reference.

## Features

TODO(nelsonjr): Get a list of features here and how they relate to the settings.


## Types


### `Api::Resource::HashArray`

A hash array is the same as a hash with the exception that the key should match
an existing object. If the object is not defined in the `api.yaml` the compiler
will fail with error.

This is to ensure deleted objects are not left behind and litters the
`puppet.yaml` files.


## Settings


### `manifest` (`Provider::Puppet::Manifest`)

Every Puppet module contains a manifest, which provides the basic information
for the module. This information is parsed by Puppet (and Puppet Forge) to
reason about the module (name, version, etc).

Example:

    manifest: !ruby/object:Provider::Puppet::Manifest
      version: '0.1.0'
      ...

#### `version` (string)

The version of the module, in [SemVer][semver-home] format. Example:

    version: '0.1.0'

#### `source` (URL)

The URL of the source code for the module. Example:

    source: 'https://github.com/GoogleCloudPlatform/puppet-google-compute'

#### `homepage` (URL)

The URL of the project page for the module. Example:

    homepage: 'https://github.com/GoogleCloudPlatform/puppet-google-compute'

#### `issues` (URL)

The URL where bugs, issues and feature requests are tracked. Example:

    issues:
      'https://github.com/GoogleCloudPlatform/puppet-google-compute/issues'

#### `summary` (string)

A short description of the module. Example:

    summary: 'A Puppet module to manage Google Compute Engine resources'

#### `tags` (list(string))

A list of tags for the module. Used by Puppet Forge to build keyword search.
Example:

    tags:
      - google
      - cloud
      - compute
      - engine
      - gce

#### `requires` (list(`Provider::Config::Requirements`))

A list of modules that this product module depends on. Example:

    requires:
      - !ruby/object:Provider::Config::Requirements
        name: 'google-gauth'
        versions: '< 0.2.0'

> Note: `google-gauth` should always be listed in the required modules.

#### `operating_systems` (list(`Provider::Config::OperatingSystem`))

The list of operating systems (and versions) that this module fully supports.
Example:

    operating_systems:
      - !ruby/object:Provider::Config::OperatingSystem
        name: RedHat
        versions:
          - '6'
          - '7'

Note that we have a file that can be included with the common operating systems
we run Puppet modules on. Then the definition should be:

      operating_systems:
    <%= indent(include('provider/puppet/common~operating_systems.yaml'), 4) %>


### `objects` (`Api::Resource::HashArray`)

Provides the root of customizations for a resource defined in the `api.yaml`.

Example ([Compute][compute-puppet] `Disk` cannot be changed via API updates):

    objects: !ruby/object:Api::Resource::HashArray
      Disk:
        flush: raise 'Disk cannot be edited'

#### `objects/<object>/create` (code)

Creates a custom `create` function. Puppet uses `create` when it detects the
object does not exist and user specified `ensure => present` in the manifest.

Example: [products/compute/puppet.yaml][dns-puppet]

#### `objects/<object>/delete` (code)

Creates a custom `delete` function. Puppet uses `delete` when it detects the
object exists and user specified `ensure => absent` in the manifest.

Example: [products/compute/puppet.yaml][dns-puppet]

#### `objects/<object>/flush` (code)

Creates a custom `flush` function. Puppet uses `flush` to update values for the
resources if they mismatch with the catalog.

Example: [products/compute/puppet.yaml][dns-puppet]

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

Example ([Compute][compute-puppet] `BackendService` requires that we provide the
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

Example: [products/compute/puppet.yaml][dns-puppet]

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

Example ([Compute][compute-puppet] includes a file with helper functions to be
used by the `Disk` object):

    objects: !ruby/object:Api::Resource::HashArray
      Disk:
        provider_helpers:
          include:
            - 'products/compute/helpers/provider_disk_type.rb'


### `properties` (list(string))

> Deprecated: Properties are automatically inferred and included into the
> providers and types automatically.


### `tests` (`Api::Resource::HashArray`)

Specifies extra tests to be added to the [rspec][rspec-home] unit tests.

This will display code that should override auto-generated testing code. View
the [DNS tests][dns-puppet-test] for the most comprehensive view of altered tests.


### `functions` (list(`Provider::Puppet::Function`))

Describes Puppet client functions provided by this module.

> To add the same documentation to the function code use the provider function
> `emit_function_doc(name)`.
> For example `emit_function_doc('gcompute_image_family')`.

TODO(nelsonjr): Make the `functions:` actually generate the function and remove
the requirement above about files:compile:

> Note that this definition does not include the actual function code (yet). It
> generates the documentation and README updates, but it is not involved in the
> function generation, so it is still required to have the function listed in
> the files:compile: section.

#### `name` (string)

The name of the function

#### `description` (string)

The description of the function

#### `arguments` (list(`Provider::Puppet::Function::Argument`))

The arguments that the function requires. Example:

    arguments:
      - !ruby/object:Provider::Puppet::Function::Argument
        name: image_family
        type: Api::Type::String
        description: 'the name of the family, e.g. ubuntu-1604-lts'

#### `examples` (list(string))

Examples to be included in the documentation

    examples:
      - gcompute_image_family('ubuntu-1604-lts', 'ubuntu-os-cloud')
      - gcompute_image_family('my-web-server', 'my-project')

#### `notes` (string)

Additional notes to be added at the end of the documentation.

    notes: |
      Note: In the case of private images, your credentials will need to have
      the proper permissions to access the image.

      To get a list of supported families you can use the gcloud utility:

        gcloud compute images list


### `examples` (`Api::Resource::HashArray`)

Lists the file names of each example included for each object. All examples
must be located in a product's files/ directory. Each file will be included in
the examples folder with the prefix 'examples\~'. The convention is to include a
<resource>.pp file and a delete_<resource>.pp file

Example:

    examples: !ruby/object:Api::Resource::HashArray
      Instance:
        - instance.rb # refers to files/examples~cookb
        - delete_instance.rb

### `files` (`Provider::Config::Files`)

Lists files to be copied or compiled into the final product. We can add files on
a per-product basis at our whim. If the file is simply copied use the `copy`
section. If you need to execute code to build the file use the `compile`
section instead.

Example:

    files: !ruby/object:Provider::Config::Files
      copy:
        lib/google/copy_file.rb: copy_file.rb
      compile:
        lib/google/object_store.rb: google/object_store.rb
        lib/puppet/functions/gcompute_image_family.rb:
          products/compute/files/function~image_family.rb.erb

### `files/copy` (`Hash`)

Lists files that should be copied directly to the module without any
compilation. List the location in the module as the key (relative to the module
root) and the location in the code generator (relative to Magic Module root) as
the value.

### `files/compile` (`Hash`)

Lists files that should be compiled and then copied to the module. List the
location in the module as the key (relative to the module root) and the
location in the code generator (relative to Magic Module root) as the value.


### `test_data` (`Provider::Config::TestData`)

Data files that will be used to generate canned responses from the API. During
tests we trap all API calls and return canned data to avoid the need to hit GCP
with a real request.

This is useful to both make the tests faster and simple, as well as avoid
charges and operating setup to developers helping us improve the modules.

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

A list of test data files to be included with the module. These files will be
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

The name of the class where exceptions are needed. Prefix with module chain
if multiple classes of the same name exist and only one needs exceptions.

### `style[]/pinpoints[]/exceptions` (list(`String`))

A list of necessary Rubocop exceptions.

TODO(nelsonjr): Add a link to each example pointing to the other `puppet.yaml`
files that exist in the products.


[compute-puppet]: ../products/compute/puppet.yaml
[dns-puppet]: ../products/dns/puppet.yaml
[dns-puppet-test]: ../products/dns/test.yaml
[compute-instance-example]: ../products/compute/files/examples~instance.pp
[semver-home]: http://semver.org/
[rspec-home]: http://rspec.info/
[rubocop-home]: https://github.com/bbatsov/rubocop
