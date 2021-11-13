# tpgtools

`tpgtools` is the generator responsible for creating DCL-based resources in the
Terraform Google Provider (TPG). The DCL provides [OpenAPI schema objects](https://swagger.io/specification/#schema-object)
to describe the available types, and `tpgtools` uses those to construct
Terraform resource schemas.

## Usage

`tpgtools` expects to run targeting a "root service directory", a dir-of-dirs
where the child dirs contain OpenApi specs for resources such as the `api/` path
above. Additionally, overrides are expected in a similar structure (as seen in
the `overrides/` path. For example:

```
go run . --path "api" --overrides "overrides"
```

This will output the file contents of each resource to stdout, for fast
iteration. You can filter by service and resource to make it more useful:

```
go run . --path "api" --overrides "overrides" --service redis --resource instance
```

To persist the output(s) to files you can specify an output path. This is the
most familiar experience for MMv1 developers. For example:

```
go run . --path "api" --overrides "overrides" --output ~/tpg-fork
```

If generation fails, an error should be logged showing what went wrong. The raw
output will be returned, and line numbers (if available in the error) will
correspond to the line numbers in the output.

### Version

You can specify a version such as `beta` using the `--version`:


```
go run . --path "api" --overrides "overrides" --output ~/tpg-fork --version "beta"
```

### Accessory Code

To generate accessory code such as `serializarion`, you can specify the `--mode`:

```
go run . --path "api" --overrides "overrides" --output ~/some/dir --mode "serialization"
```

## Development

`tpgtools` builds resources using Go Templates, with the templates stored under
the `templates/` directory. They're fed the `Resource` type, which contains
resource-level metadata and a list of `Property` types which represent top-level
fields in a Terraform resource. Additionally, `Property`s can contain other
`Property` objects to create a nested structure of fields.

`main.go` parses the OpenAPI schemas you've selected into Go structs and then
transforms them into `Resource`/`Property` before running them through Go
Templates.

### Overrides

Overrides are specified per-resource, with a directory structure parallel to the
OpenAPI specs. Inside each resource file is an unordered array of overrides made
up of an override type (like `CUSTOM_DESCRIPTION` or `VIRTUAL_FIELD`) as well as
a field they affect (if a field is omitted, they affect the resource) and an
optional `details` object that will be parsed for additional metadata.

For example, override entries will look like the following:

```yaml
- type: CUSTOM_DESCRIPTION
  field: lifecycle.rule.condition.age
  details:
    description: Custom description here.
```

#### Samples

We will autoingest samples from the dcl, however we currently must
manually fill the substitutions for these samples.

You may need to re-serialize `serialization.go` if you are adding newer resources.
To do this
```
cd tpgtools
go get -u github.com/GoogleCloudPlatform/declarative-resource-client-library
make serialize
```

To do so, first create a folder in the relevant product
```
$ cd overrides/{{product}}/
$ mkdir samples
$ cd samples
```

Then make a folder for the resource you would like to add samples for.
```
$ mkdir {{resource}}
$ cd resource
```

Create a meta.yaml file. This file will merge with any tests you create
providing sustitutions and other relevant test data (injections, hiding, ect..)

Provide the relevant sustitutions needed. See the referenced variables in the dcl
jsons. They should surrounded by `{{}}`
```
substitutions:
  - substitution: "project"
    value: ":PROJECT"
  - substitution: "region"
    value: ":REGION"
  - substitution: "name"
    value: "trigger"
  - substitution: "topic"
    value: "topic"
  - substitution: "event_arc_service"
    value: "service-eventarc"
  - substitution: "service_account"
    value: "sa"
```

If you need to hide sample from doc or hide a sample from docs you can do so here as well.
```
doc_hide:
  - basic.tf.tmpl
test_hide:
  - basic_trigger.yaml
```

Any files with a `.tf.tmpl` (terraform template) extension located in the `override/{{product}}samples/{{resource}}` directory
and without `_update` in the name are considered to be tests independently.
These are normal terraform files with the desired sustituted variables surrounded by `{{}}`.
These tests will also use the substitutions defined in the `meta.yaml`. If you want to provide test specific
rules (updates, ect), you can create a yaml file with the same name as the `.tf.tmpl` file. Here you can supply updates
or version requirements specific to this test. If version is ommitted the sample assumed to run against all versions.
```
updates:
  - resource: basic_update_transport.tf.tmpl
  - resource: basic_update_transport_2.tf.tmpl
version:
  - beta
```

If you want to add additional rules the following options are currently supported within `meta.yaml`
```
  - substitution: "name"
    value: "workload"
  - substitution: "region"
    value: ":REGION"
 ....
ignore_read:
  - "billing_account"
  - "kms_settings"
  - "resource_settings"
  - "provisioned_resources_parent"
check:
  -  deleteAssuredWorkloadProvisionedResources(t)
extra_dependencies:
  - "time"
  - "log"
  - "strconv"
  - "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
code_inject:
  - delete_assured_workload_provisioned_resources.go
doc_hide:
  - basic.tf.tmpl # basic_update.tf.tmpl auto hides
  - full.tf.tmpl
test_hide:
  - basic_workload.yaml
```

#### New Resource Guide (tpgtools only)

This guide is written to document the process for adding a resource to the Google Terraform Provider (TPG) after it has been added to the [DCL](https://github.com/GoogleCloudPlatform/declarative-resource-client-library).

## Adding Resource Overrides

Every resource added via tpgtools needs an override file for every version it is available at. This file may be empty, but must exist. A resource available at GA (TPG) must also exist at beta (TPGB) and needs a corresponding override file at beta. These override files are often identical between versions. This file should exist at tpgtools/overrides/$PRODUCT_NAME/$VERSION/$RESOURCE.yaml. For example, [this override](https://github.com/GoogleCloudPlatform/magic-modules/blob/master/tpgtools/overrides/assuredworkloads/beta/workload.yaml) exists for the product assuredworkloads resource workload at beta version.

# What should I add to the overrides file?

For the first pass try nothing. If behavior is needed to change the resource's behavior in Terraform from the default DCL resource, then you can add these at a later date. Generally no overrides will be needed for a new resource, but they will be needed for resources that exist in mmv1 that are being converted to use tpgtools.

## Adding Samples

See [guide above](#samples) for adding samples.

## Adding Additional Information to Handwritten/MMv1 Generated Files

There are a few files within mmv1 that use Ruby templates to generate shared files for the provider. One of these is provider.go (generated by [provider.go.erb](https://github.com/GoogleCloudPlatform/magic-modules/blob/master/mmv1/third_party/terraform/utils/provider.go.erb)) which contains the map of resources to their names within the provider. Currently these need to be updated by hand to add a new resource. There are a couple of other places that need to be modified as well to add endpoint information if the resource being added is part of a new product.

# Adding to provider.go

Add the resource definition within the `ResourceMapWithErrors` map. There is currently a block of tpgtools resources, so adding to that block is preferred. If this resource is available only at a specific version (beta, private) then add version tags around the resource definition.

Additionally you may need to add a few line to provider.go for the endpoint if this resource is in a new product that doesn't have any other resources currently in the provider. Find the correct lines in the diff on `provider_dcl_endpoints.go` within the provider once it has been generated with the new resource.

## Adding Documentation

Documentation will be automatically generated based on the resource description in the DCL and examples will be generated from the samples for this resource.