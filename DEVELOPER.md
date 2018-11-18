# Developer Guide

All the information required for Magic Modules (MM) to compile a specific
product is usually contained within its folder under the [`products`](products/)
folder. For example all the sources for the Google Compute Engine are inside the
[`products/compute`](products/compute) folder.

When compiling a product you specify the path to the products folder with the
`-p` parameter.

Please refer to these documents for general guidelines on how to write Magic
Modules code:

  - [GOVERNANCE][governance]
  - [Template SDK][template-sdk]


## Anatomy of a Product

A product definition for a specific provider contains a few basic components:

  1. Product definition: `api.yaml`
  2. Provider dependent product definitions: `<provider>.yaml`
  3. Examples
  4. Tests

### Product definition: api.yaml

The `api.yaml` contains all the object definitions for the GCP product. It also
contains the relationships between objects.

It also includes other product specific information, such as product version and
API endpoint.

#### Resource

A resource is an object defined and supported by the product. For example a
virtual machine, a disk, a network, a container, a container cluster are all
resources from MM's point of view.

A resource defines some basic properties about the resource such as name,
endpoint name, as well as its properties and parameters.

    - !ruby/object:Api::Resource
      name: 'Address'
      kind: 'compute#address'
      base_url: projects/{{project}}/regions/{{region}}/addresses
      exports:
				...
      description: |
				...
      parameters:
				...
      properties:
				...

##### Resource / [ parameters | properties ]

Specifies fields that the user can specify to define the object and its desired
state.

> The main difference between a parameter and a property is that a parameter is
> not part of the object, and will be used by the user to convey data to the
> generated code, usually to help locate the object. Every property will be
> eventually persisted as fields of the remote object in GCP.

Required fields:

-  `type`: Specifies the type of the parameter / property
-  `name`: The user facing name of the parameter / property
-  `description`: The description for the parameter / property

Optional fields:

-  `required`: true|false indicating if the field if required to be specified by
   the user
-  `input`: true|false indicating if the field will be used as "input", which
   means that the field will be used only during the _creation_ of the object
-  `output`: true|false indicating that the field is produced and controlled by
   the server. The user can use an output field to ensure the value matches
   what's expected.
-  `field`: In case we want the user facing name to be different from the
   corresponding API property we can use `field` to map the user facing name
   (specified by the `name` parameter) to the backend API (specified by the
   `field` parameter)
-  `resource`: A resource this resource is dependent upon. See
   [Api::Type::ResourceRef](#resource-ref).
-  `imports`: An imported property from the dependent resource specified by
   `resource`. See [Api::Type::ResourceRef](#resource-ref).

Example:

    - !ruby/object:Api::Type::ResourceRef
      name: 'region'
      resource: 'Region'
      imports: 'name'
      description: |
        URL of the region where the regional address resides.
        This field is not applicable to global addresses.
      required: true

> Please describe these fields from the user's perspective and not from the API
> perspective. That will allow the generated documentation be useful for the
> provider user that should not care about how the internals work. This is a
> core strength of the MM paradigm: the user should neither need to know nor
> care how the backend works. All he cares is how to describe in a high-level,
> uniform and elegant way his dependencies.
>
> For example a in a virtual machine disk, do not say "disk: array or full URI
> (self link) to the disks". Instead you should say "disks: reference to a list
> of disks to attach to the machine".
>
> Also avoid language or provider specific lingo in the product definitions.
>
> For example instead of saying "a hash array of name to datetime" say "a map
> from name to timestamps". While the former may be correct for a provider in
> Java it will not be in a Ruby or Python based output.


### Provider dependent product definitions: <provider>.yaml

Each provider has their own product specific definitions that are relevant and
necessary to properly build the product for such provider.

For example a Puppet provider requires to specify the Puppet Forge URL for the
module being created. That goes into the product `puppet.yaml`. Similarly a Chef
provider requires metadata about dependencies of other cookbooks. That goes into
the product `chef.yaml`.

Detailed provider documentation:

- [`puppet.yaml`][puppet-yaml] on how to create Puppet product definitions
- [`chef.yaml`][chef-yaml] on how to create Chef product definitions

### Examples

It is strongly encouraged to provide live, working examples to end users. It is
also good practice to test the examples to make sure they do not become stale.

To test examples you can use our `[end2end](tools/end2end)` mini framework.


## Types

When defining a property you have to specify its type. Although some products
may not care about types and may convert from/to strings, it is important to
know the type so the compiler can do the best job possible to validate the
input, ensure consistency, etc.

Currently MM supports the following types:

-  `Api::Type::Array`: Represents an array of values. The type of the values is
   identified by the `item\_type`property.
-  `Api::Type::Boolean`: A boolean (true or false) value.
-  `Api::Type::Constant`: A constant that will be passed to the API.
-  `Api::Type::Double`: A double number.
-  `Api::Type::Enum`: Input is allowed only from a fixed list of values,
   specified by the `values` property.
-  `Api::Type::Integer`: An integer number.
-  `Api::Type::Long`: A long number
-  `Api::Type::KeyValuePairs`: A string -> string key -> value pair such as
   labels
-  `Api::Type::Map`: A string -> `Api::Type::NestedObject` map.
-  `Api::Type::NestedObject`: A composite field, composed of inner fields. This
   is used for structures that are nested.
-  <a id="resource-ref"></a>`Api::Type::ResourceRef`: A reference to another object described in the
   product. This is used to create strong relationship binding between the
   objects, where the generated code will make sure the object depended upon
   exists. A `ResourceRef` also specifies which property from the dependent
   object we are interested to fetch, by specifying the `resource` and `imports`
   fields.
-  `Api::Type::String`: A string field.
-  `Api::Type::Time`: An RFC 3339 representation of a time stamp.


## Exports

When data needs to be read from another object that we depend on we use exports.

For example a virtual machine has many disks. When specifying the disks we do
not want the user to worry about which representation Google developers chose
for that object (is it the self\_link? is it just the disk name? is it the
zone+name?).

Instead we describe the relationship in `api.yaml` and let the generated code
deal with that. The user will only say "vm X uses disks [A, B, C]". As GCP VM
requires the disk self\_link the generated code will compute their values and
pass along.

All fields need to be explicitly exported before they can be imported and
dependent upon by other objects. This allows us to track and verify the
dependencies more reliably.

The following exports are allowed:

-  `<property-name>`: By specifying the `name` of an existing property of the
   object, we're allowing other objects to import it.
-  `Api::Type::FetchedExternal`: Exports the version of the data that was
   returned from the API. That's useful when the user specification of the
   property is different from the representation coming from the API
-  `Api::Type::SelfLink`: URI that represents the object

> To avoid overexposure and inefficient code, only export the fields you
> actively require to access from dependent objects, and practice housecleaning
> by removing exports if you change or remove an object that depends on it.


[puppet-yaml]: docs/puppet.yaml.md
[chef-yaml]: docs/chef.yaml.md
[governance]: GOVERNANCE.md
[template-sdk]: TEMPLATE_SDK.md
