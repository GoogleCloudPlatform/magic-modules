# tpgtools

`tpgtools` is the generator responsible for creating DCL-based resources in the
Terraform provider for Google Cloud (TPG). The DCL provides
[OpenAPI schema objects](https://swagger.io/specification/#schema-object) to
describe the available types, and `tpgtools` uses those to construct Terraform
resource schemas.

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

To generate accessory code such as `serializarion`, you can specify the
`--mode`:

```
go run . --path "api" --overrides "overrides" --output ~/some/dir --mode "serialization"
```

## New Resource Guide

This guide is written to document the process for adding a resource to the
Terraform provider for Google Cloud (TPG) after it has been added to the
[DCL](https://github.com/GoogleCloudPlatform/declarative-resource-client-library).

### Adding Resource Overrides

Every resource added via tpgtools needs an override file for every version it is
available at. This file should be empty, but must exist. A resource available at
GA (TPG) must also exist at beta (TPGB) and needs a corresponding override file
at beta. These override files are often identical between versions. This file
should exist at tpgtools/overrides/$PRODUCT_NAME/$VERSION/$RESOURCE.yaml. For
example,
[this override](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/tpgtools/overrides/assuredworkloads/beta/workload.yaml)
exists for the product assuredworkloads resource workload at beta version.

Override files contain information on how the Terraform representation of the
resource should be different from the DCL's representation. This could be naming
or behavior differences, but for a new resource implemented through tpgtools
there should be no differences from the DCL's representation.

### Adding Samples

For a deeper understanding on test anatomy please read the accompanying
[Tests and Sample Anatomy](#tests-and-sample-anatomy)

#### Create a meta.yaml file

Create a meta.yaml file in the overrides directory for the resource. For
example,
[this meta.yaml](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/tpgtools/overrides/assuredworkloads/samples/workload/meta.yaml)
file exists for the assured workloads resource. This file will merge with any
tests. You can customize behavior of the tests and examples generated dcl
samples data (injections, hiding, ect..). See
[the section of the meta.yaml file](#the-metayaml-file) for a more detailed
dive.

#### Adding DCL Tests

Start by copying the relevant samples from the DCL for your new resource. These
will be added to the
[tpgtools/api folder](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/tpgtools/api)
under the relevant product. These samples can be found under the samples/ folder
within the DCL for the resource being added. For example, assured workloads can
be found
[here](https://github.com/GoogleCloudPlatform/declarative-resource-client-library/tree/main/services/google/assuredworkloads/samples).

Re-serialize `serialization.go` to enable generating configs from samples by
running:

```
make serialize
```

#### Adding a Non DCL Test

In some cases you may need to add a non DCL test when either the current tests
are insufficient or you want to showcase/test some specific behavior not present
in the dcl tests.

If you need to write tests manually you can add terraform templates to the
relevant `./overrides/<product>/samples/<resource>` folder.

A terraform template test has the following anatomy.

* `<test-name>.yaml` - this is the test definition
* `<test-name>.tf.tmpl` - this is the accompanying terraform configuration. A companion to the definition if you will.

The `<test-name>.yaml` test specific configurations. For example it lists the
variables to replace in the template companion. You can also add additional
templates as updates. This will act as sequential applies and are useful for
testing update specific behavior. Make sure any templates added as an update has
the `_update.tf.tmpl` extension.

The following is an example test definition.

```yaml
updates:
- resource: ./basic_update.tf.tmpl
variables:
  - name: "name"
    type: "resource_name"
  - name: "region"
    type: "region"
```

The `<test-name>.tf.tmpl` file is simply a terraform configuration. Any
replacements should be surrounding by double brackets `{{ }}`. The variable name
from the test definition will be used to key into and replace these.

### Adding Documentation

Provided you have added samples for the resource, documentation will be
automatically generated based on the resource description in the DCL and
examples will be generated from the samples for this resource. If you did not
provide samples for the resource, then documentation will need to be written by
hand.

### Handwritten Tests

Sometimes you may need to test unusual resource behavior in a way that does not
fit well with generated tests. In this circumstance you can write a
[handwritten test file](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/third_party/terraform/README.md#testing)
and add it
[here](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/tests).

These tests can be used for more granular testing of specific behavior and add
custom checks. Tests in these files will not have examples generated for them,
so handwritten tests should not be considered a replacement for samples.

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

### Tests and Sample Anatomy

Adding samples is essential for generating tests and documentation. In fact
Documentation won't generate without it!

Tests come from two sources.

* The top level api (`./api`) folder. If you look
in here you'll see a bunch of yaml files and json files. These are DCL tests!
Forked from the dcl library.
* The override folder
(`./overrides/<product>/samples/<resource>`). This contains `meta.yaml` a file
used for managing resource-wide test configurations and custom non-dcl tests.

In either case, DCL or non-DCL, every test definition is a yaml file which takes
in variables.

```yaml
variables:
  - name: "name"
    type: "resource_name"
  - name: "org_id"
    type: "org_id"
```

`type` is inferred from `sample.go`'s translation map to figure out what needs
to be placed in the field. `name` is used for substitution in the file itself
and in the case of `resource_name`, actually used to create the value itself.

#### The meta.yaml file

In the `./overrides/<product>/samples/<resource>` a `meta.yaml` file exists
which controls configuration which applies to multiple tests or hiding/showing
specific tests.

If you need to hide sample from a doc or disable a sample from the tests you can
do so here.

```yaml
doc_hide:
  - basic.tf.tmpl
  - full.tf.tmpl
test_hide:
  - basic_workload.yaml
```

If you want to add additional rules the following options are currently supported within `meta.yaml`

```yaml
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
