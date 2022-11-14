---
title: "Handwritten"
weight: 1
# bookFlatSection: false
# bookToc: true
# bookHidden: false
# bookCollapseSection: false
# bookComments: false
# bookSearchExclude: false
---


# Handwritten

## Overview

The Google providers for Terraform have a large number of handwritten go files, primarily for resources written before Magic Modules was used with them. Most handwritten files are expected to stay handwritten indefinitely, although conversion to a generator may be possible for a limited subset of them.

Handwritten resources like [google_container_cluster](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/container_cluster) can be identified if they have source code present under the [mmv1/third_party/terraform/resources](./resources) folder or by the absence of the `AUTO GENERATED CODE` header in their [Go source](https://github.com/hashicorp/terraform-provider-google/blob/main/google/resource_container_cluster.go) in the downstream repositories. Handwritten datasources should be under the [mmv1/third_party/terraform/data_sources](./data_sources) folder, tests under the [mmv1/third_party/terraform/tests](./tests) folder and web documentaion under the [mmv1/third_party/terraform/website](./website) folder.

## Table of Contents
- [Contributing](#contributing)
  - [Shared Concepts](#shared-concepts)
  - [Resource](#resource)
  - [Datasource](#datasource)
  - [IAM Resources](#iam-resource)
  - [Testing](#testing)
    - [Simple Tests](#simple-tests)
    - [Update tests](#update-tests)
    - [Testing Beta Features](#testing-beta-features)
  - [Documentation](#documentation)
  - [Beta Feature](#beta-feature)
    - [Adding a beta resource](#adding-a-beta-resource)
    - [Adding beta field(s)](#adding-beta-fields)
    - [Tests that use a beta feature](#tests-that-use-a-beta-feature)
    - [Promote a beta feature](#promote-a-beta-feature)


## Contributing

We're glad to accept contributions to handwritten resources. Tutorials and guidance on making changes are available below.

### Shared Concepts

This section will serve as a point of reference for some shared concepts that
all handwritten files share. It's meant to be an introduction to our serialization
strategy and overview.

#### Serialization strategy
The go files within the directory files are copied literally to their respective providers.
Our serialization methodology may seem complicated but for the case of handwritten resources its quite
simple. Editing the file will change its counterpart downstream.

#### go and go.erb
Within the third party library you'll notice `go` and `go.erb` files.
Go files are native golang code while go.erb pass through ruby before
being serialized. The reason `go.erb` files exist are to protect certain
properties or fields from entering the `ga` provider. Thus you'll often see
lines like `<% unless version == 'ga' -%>` within the file. These blocks
will omit the enclosure from being output to the GA provider. In the
rare case where you are promoting all fields to `ga` and these blocks
are no longer needed you can remove the `.erb` extension.

#### Create, Read, Update, Delete
As far as terraform schema is concerned these are the functions we
need to provide for terraform to be able to provision and delete
resources. In editing any fields you'll likely be adding functionality to
these functions or implementing them wholesale.


#### Expanders and Flatteners
Expanders and flatteners are concepts created to simplify common patterns
and add conformity/code consistency. Essentially expanders are functions
used to segregate some translation from terraform representation to api representation.
We will use these to encapsulate this translation for blocks and/or complicated fields.
This allows our code to be concise and functionality to be readable and easily
apparent by separating these into their own functions. While expanders are used
for terraform to api, flatteners do just the opposite. Converting api to terraform.

Thus
* expanders - helper functions used for translating tf -> api representation
* flatteners - helper functions used for translating api -> tf representation

### Resource
While we no longer accept new handwritten resources except in rare cases. Understanding
how to edit and add to existing resources may be important for implementing new fields
or changing existing behavior.

To edit an existing resource to add a field there are four steps you'll go through.
1. Add the new field to the schema
1. Implement the respective flattener and/or expander for the new field
1. Add a testcase for the field or extend an existing one
1. Add documentation for the field to the respective markdown file

##### Adding a new field to the schema
To add a new field you will have to compare an existing resource
to it's respective rest api documentation. Dependant on how the api implements
the field we will in almost all cases mirror the structure. For example if there is
an `enabled` field nested under a ``IdentityServiceConfig`` block we will mirror
this within the schema.

Thus the block for terraform to utilize this field would then be
```terraform
resource "x" "y" {
  identity_service_config{
    enabled = true
  }
}
```

You might think it convoluted to provide such a structure. Why not
simply provide a single `enable_identity_service_config`. One constant
has echoed through our mind as terraform developers through the years.
Api's are ever evolving. Mirroring the api gives us the best chance to stay
in step with that evolution. Therefore if `IdentityServiceConfig` is extended with new
parameters in the future we can cleanly encapsulate those into the existing block(s).

As far as providing the field itself, it's fairly straightforward. Mirror the field
from the api and look to the other fields and the schema type in the SDK to see
what's available and how to structure it. For the documentation, copying the documentation
from the rest api will be the usual practice.

If you are adding a field that is an ENUM from the api standpoint its best practice
to provide it as as a string to the provider. This field will likely have values
added to it by the api and this future proofs our provider to support new values without
haven't to make new additions. There will be rare exceptions, but generally its a good
practice.

#### Implement the respective flattener and/or expander for the new field
Once you've added the field to the schema you will implement the corresponding
expander/flattener. See [expanders and flatters](#expanders-and-flatteners) for
more context on what these fields are used for. Essentially we will be editing the
read, create, and update operations to parse the schema and call the api to make
the changes to the state of the resource. Following existing patterns to create
this operation will be the best way to implement this. As there are many unique ways
to implement a given field we won't get into specifics.

For example a field in bigtable `google_sheets_options` containers two nested properties.
`range` and `skip_leading_rows`.

```golang
						// GoogleSheetsOptions: [Optional] Additional options if sourceFormat is set to GOOGLE_SHEETS.
						"google_sheets_options": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: `Additional options if source_format is set to "GOOGLE_SHEETS".`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									// Range: [Optional] Range of a sheet to query from. Only used when non-empty.
									// Typical format: !:
									"range": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Range of a sheet to query from. Only used when non-empty. At least one of range or skip_leading_rows must be set. Typical format: "sheet_name!top_left_cell_id:bottom_right_cell_id" For example: "sheet1!A1:B20"`,
										AtLeastOneOf: []string{
											"external_data_configuration.0.google_sheets_options.0.skip_leading_rows",
											"external_data_configuration.0.google_sheets_options.0.range",
										},
									},
									// SkipLeadingRows: [Optional] The number of rows at the top
									// of the sheet that BigQuery will skip when reading the data.
									"skip_leading_rows": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: `The number of rows at the top of the sheet that BigQuery will skip when reading the data. At least one of range or skip_leading_rows must be set.`,
										AtLeastOneOf: []string{
											"external_data_configuration.0.google_sheets_options.0.skip_leading_rows",
											"external_data_configuration.0.google_sheets_options.0.range",
										},
									},
								},
							},
						},
```

To simplify the implementation management of these fields can be delegated to expanders and
flatteners.

```golang
func expandGoogleSheetsOptions(configured interface{}) *bigquery.GoogleSheetsOptions {
	if len(configured.([]interface{})) == 0 {
		return nil
	}

	raw := configured.([]interface{})[0].(map[string]interface{})
	opts := &bigquery.GoogleSheetsOptions{}

	if v, ok := raw["range"]; ok {
		opts.Range = v.(string)
	}

	if v, ok := raw["skip_leading_rows"]; ok {
		opts.SkipLeadingRows = int64(v.(int))
	}
	return opts
}

func flattenGoogleSheetsOptions(opts *bigquery.GoogleSheetsOptions) []map[string]interface{} {
	result := map[string]interface{}{}

	if opts.Range != "" {
		result["range"] = opts.Range
	}

	if opts.SkipLeadingRows != 0 {
		result["skip_leading_rows"] = opts.SkipLeadingRows
	}

	return []map[string]interface{}{result}
}

```

#### Add a testcase for the field or extend an existing one
Once your field has been implemented, go to the corresponding test file for
your resource and extend it. If your field is updatable it's good practice to
have a two step apply to ensure that the field *can* be updated. You'll notice
a lot of our tests have a import state verify directly after apply. These
steps are important as they will essentially attempt to import the resource
you just provisioned and *verify* that the field values are consistent with the
applied state. Please test all fields you've added to the provider. It's important
for us to ensure all fields are usable and workable.

#### Add documentation for the field to the respective markdown file
See [Documentation](#documentation) for more information. Essentially you will
just be opening the corresponding markdown file and adding documentation, likely
copied from the rest api to the markdown file. Follow the existing patterns there-in.

### Datasource
**Note** : only handwritten datasources are currently supported

Datasources are like terraform resources except they don't *create* anything.
They are simply read-only operations that will expose some sort of values needed
for subsequent resource operations. If you're adding a field to an existing
datasource, check the [Resource](#resource) section. Everything there will
be mostly consistent with the type of change you'll need to make. For adding
a new datasource there are 5 steps to doing so.

1. Create a new datasource declaration file and a corresponding test file
1. Add Schema and Read operation implementation
1. Add the datasource to the `provider.go.erb` index
1. Implement a test which will create and resources and read the corresponding
  datasource.
1. Add documentation.

For creating a datasource based off an existing resource you can [make use of the
schema directly](https://github.com/GoogleCloudPlatform/magic-modules/blob/1d293f7bfadacaa20580874c8e8634827fb99a14/mmv1/third_party/terraform/data_sources/data_source_cloud_run_service.go).
Otherwise [implementing the schema directly](https://github.com/GoogleCloudPlatform/magic-modules/blob/1d293f7bfadacaa20580874c8e8634827fb99a14/mmv1/third_party/terraform/data_sources/data_source_google_compute_address.go),
similar to normal resource creation, is the desired path.

### IAM Resource

Handwritten IAM support is only recommended for resources that cannot be managed
using [MMv1](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1#iam-resource),
including for handwritten resources, due to the need to manage tests and
documentation by hand. This guidance goes through the motions of adding support
for new handwritten IAM resources, but does not go into the details of the
implementation as any new handwritten IAM resources are expected to be
exceptional.

IAM resources are implemented using an IAM framework, where you implement an
interface for each parent resource supporting `getIamPolicy`/`setIamPolicy` and
the associated IAM resources that target that parent resource- `_member`,
`_binding`, and `_policy`- are created by the framework.

To add support for a new target, create a new file in
`mmv1/third_party/terraform/utils` called `iam_{{resource}}.go`, and implement
the `ResourceIamUpdater`, `newResourceIamUpdaterFunc`, `iamPolicyModifyFunc`,
`resourceIdParserFunc` interfaces from
https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/third_party/terraform/utils/iam.go.erb
in public types, alongside a public `map[string]*schema.Schema` containing all
fields referenced in the resource.

Once your implementation is complete, add the IAM resources to `provider.go`
inside the `START non-generated IAM resources` block, creating the concrete
resource types using the `ResourceIamMember`, `ResourceIamBinding`, and
`ResourceIamPolicy` functions. For example:

```go
				"google_bigtable_instance_iam_binding":         ResourceIamBinding(IamBigtableInstanceSchema, NewBigtableInstanceUpdater, BigtableInstanceIdParseFunc),
				"google_bigtable_instance_iam_member":          ResourceIamMember(IamBigtableInstanceSchema, NewBigtableInstanceUpdater, BigtableInstanceIdParseFunc),
				"google_bigtable_instance_iam_policy":          ResourceIamPolicy(IamBigtableInstanceSchema, NewBigtableInstanceUpdater, BigtableInstanceIdParseFunc),
```

Following that, write a test for each resource exercising create and update for
both `_policy` and `_binding`, and create for `_member`. No special
accommodations are needed for the IAM test compared to a normal Terraform
resource test.

Documentation for IAM resources is done using single page per target resource,
rather than a distinct page for each IAM resource level. As most of the page is
standard, you can generally copy and edit an existing handwritten page such as
https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/third_party/terraform/website/docs/r/bigtable_instance_iam.html.markdown
to write the documentation.

### Testing

For handwritten resources and generated resources that need to test update,
handwritten tests must be added.

Tests are made up of a templated Terraform configuration where unique values
like GCE names are passed in as arguments, and boilerplate to exercise that
configuration.

The test boilerplate effectively does the following:

1.  Run `terraform apply` on the configuration, waiting for it to succeed and
    recording the results in Terraform state
2.  Run `terraform plan`, and fail if Terraform detects any drift
3.  Clear the resource from state and run `terraform import` on it
4.  Deeply compare the original state from `terraform apply` and the `terraform
    import` results, returning an error if any values are not identical
5.  Destroy all resources in the configuration using `terraform destroy`,
    waiting for the destroy command to succeed
6.  Call `GET` on the resource, and fail the test if it is still present

#### Simple Tests

Terraform configurations are stored as string constants wrapped in Go functions
like the following:

```go
func testAccComputeFirewall_basic(network, firewall string) string {
    return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_firewall" "foobar" {
  name        = "%s"
  description = "Resource created for Terraform acceptance testing"
  network     = google_compute_network.foobar.name
  source_tags = ["foo"]
  allow {
    protocol = "icmp"
  }
}
`, network, firewall)
}
```

For the most part, you can copy and paste a preexisting test case and modify it.
For example, the following test case is a good reference:

```go
func TestAccComputeFirewall_noSource(t *testing.T) {
    t.Parallel()

    networkName := fmt.Sprintf("tf-test-firewall-%s", randString(t, 10))
    firewallName := fmt.Sprintf("tf-test-firewall-%s", randString(t, 10))

    vcrTest(t, resource.TestCase{
        PreCheck:     func() { testAccPreCheck(t) },
        Providers:    testAccProviders,
        CheckDestroy: testAccCheckComputeFirewallDestroyProducer(t),
        Steps: []resource.TestStep{
            {
                Config: testAccComputeFirewall_noSource(networkName, firewallName),
            },
            {
                ResourceName:      "google_compute_firewall.foobar",
                ImportState:       true,
                ImportStateVerify: true,
            },
        },
    })
}

func testAccComputeFirewall_noSource(network, firewall string) string {
    return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "%s"
  auto_create_subnetworks = false
}
resource "google_compute_firewall" "foobar" {
  name        = "%s"
  description = "Resource created for Terraform acceptance testing"
  network     = google_compute_network.foobar.name
  allow {
    protocol = "tcp"
    ports    = [22]
  }
}
`, network, firewall)
}
```

#### Update tests

Inside of a test, additional steps can be added in order to transition between
Terraform configurations, updating the stored state as it progresses. This
allows you to exercise update behaviour. This modifies the flow from before:

1.  Start with an empty Terraform state
1.  For each `Config` and `ImportState` pair:
    1.  Run `terraform apply` on the configuration, waiting for it to succeed
        and recording the results in Terraform state
    1.  Run `terraform plan`, and fail if Terraform detects any drift
    1.  Clear the resource from state and run `terraform import` on it
    1.  Deeply compare the original state from `terraform apply` and the
        `terraform import` results, returning an error if any values are not
        identical
1.  Destroy all resources in the configuration using `terraform destroy`,
    waiting for the destroy command to succeed
1.  Call `GET` on the resource, and fail the test if it is still present

For example:

```go
func TestAccComputeFirewall_disabled(t *testing.T) {
    t.Parallel()

    networkName := fmt.Sprintf("tf-test-firewall-%s", randString(t, 10))
    firewallName := fmt.Sprintf("tf-test-firewall-%s", randString(t, 10))

    vcrTest(t, resource.TestCase{
        PreCheck:     func() { testAccPreCheck(t) },
        Providers:    testAccProviders,
        CheckDestroy: testAccCheckComputeFirewallDestroyProducer(t),
        Steps: []resource.TestStep{
            {
                Config: testAccComputeFirewall_disabled(networkName, firewallName),
            },
            {
                ResourceName:      "google_compute_firewall.foobar",
                ImportState:       true,
                ImportStateVerify: true,
            },
            {
                Config: testAccComputeFirewall_basic(networkName, firewallName),
            },
            {
                ResourceName:      "google_compute_firewall.foobar",
                ImportState:       true,
                ImportStateVerify: true,
            },
        },
    })
}

func testAccComputeFirewall_basic(network, firewall string) string {
    return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_firewall" "foobar" {
  name        = "%s"
  description = "Resource created for Terraform acceptance testing"
  network     = google_compute_network.foobar.name
  source_tags = ["foo"]
  allow {
    protocol = "icmp"
  }
}
`, network, firewall)
}

func testAccComputeFirewall_disabled(network, firewall string) string {
    return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_firewall" "foobar" {
  name        = "%s"
  description = "Resource created for Terraform acceptance testing"
  network     = google_compute_network.foobar.name
  source_tags = ["foo"]
  allow {
    protocol = "icmp"
  }
  disabled = true
}
`, network, firewall)
}
```

#### Testing Beta Features

See [Tests that use a beta feature](#tests-that-use-a-beta-feature)

### Documentation

### Beta Feature

When the underlying API of a feature is not final (i.e. a `vN` version like
`v1` or `v2`), is in preview, or the API has no SLO we add it to the
`google-beta` provider rather than the `google `provider, allowing users to
self-select for the stability level they are comfortable with.

Both the `google` and `google-beta` providers operate off of a shared codebase,
including for handwritten code. MMv1 allows us to write Go source files as
`.go.erb` templated source, and renders them as `.go` files in the downstream
repo.

The sole generator feature you need to be aware of is a "version guard", what is
effectively a preprocessor directive implemented using Embedded Ruby (ERB). A
version guard is a snippet used across this codebase by convention guarding
versioned code on an `unless` clause in a version check. For example:

```
	networkInterfaces, err := expandNetworkInterfaces(d, config)
	if err != nil {
		return nil, fmt.Errorf("Error creating network interfaces: %s", err)
	}
<% unless version == 'ga' -%>
	networkPerformanceConfig, err := expandNetworkPerformanceConfig(d, config)
	if err != nil {
		return nil, fmt.Errorf("Error creating network performance config: %s", err)
	}
<% end -%>
```

In the snippet above, the `networkInterfaces` field is generally available and
is not guarded. The `networkPerformanceConfig` field is only available at beta,
and is guarded by `unless version == ga`, and the guarded block is terminated by
an `end` statement.

If a service includes handwritten resources and mixed features or resources at
different versions, the client libraries used by each provider must be switched
using guards so that the stability level of the client library matches that of
the provider. For example, all handwritten Google Compute Engine (GCE) files
have the following guarded import:

```
<% if version == "ga" -%>
	"google.golang.org/api/compute/v1"
<% else -%>
	compute "google.golang.org/api/compute/v0.beta"
<% end -%>
```

This is not necessary for beta-only services, or for services that are generally
available in their entirety.

#### Adding a beta resource

MMv1 doesn't selectively generate files, and any file that is beta-only must
have all of its contents guarded. When writing a resource that's available at
beta, start with the following snippet:

```

<% autogen_exception -%>
package google

<% unless version == 'ga' -%>

// Add the implementation of the file here

<% end ->
```

This will generate a blank file in the `google` provider. The resource file,
resource test file, and any service or resource specific utility files should be
guarded in this way.

Documentation **should not** be guarded. Instead, write it as normal including
the following snippet above the first example.

```
~> **Warning:** This resource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta resources.
```

When registering the resource in
[`provider.go.erb`](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/third_party/terraform/utils/provider.go.erb),
the entry should be guarded:

```diff
				"google_monitoring_dashboard":                  resourceMonitoringDashboard(),
+				<% unless version == 'ga' -%>
+				"google_project_service_identity":              resourceProjectServiceIdentity(),
+				<% end -%>
				"google_service_networking_connection":         resourceServiceNetworkingConnection(),
```

If this is a new service entirely, all service-specific entries like client
factory initialization should be guarded as well. However, new services should
generally be implemented using an alternate engine- either MMv1 or tpgtools/DCL.

#### Adding beta field(s)

By contrast to beta resources, adding support for a beta field is much more
involved as small snippets of code throughout a resource file must be annotated.

To begin with, add the field to the `Schema` of the resource with guards, i.e.:

```
<% unless version == 'ga' -%>
			"network_performance_config": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: `Configures network performance settings for the instance. If not specified, the instance will be created with its default network performance configuration.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_egress_bandwidth_tier": {
							// ...
						},
					},
				},
			},
<% end -%>
```

Next, implement the d.Get/d.Set calls (for top level fields) or
expanders/flatteners for nested fields within guards.

Even if there are other guarded fields, it's recommended that you add distinct
guards per feature- that way, promotion (covered below) will be simpler as
you'll only need to remove lines rather that move them around.

#### Tests that use a beta feature

If you worked with a beta feature and had to use beta version guards in a
handwritten resource or set `min_version: beta` in a generated resource, you'll
want to version guard both the test case and configuration by enclosing them in
ERB tags like below. Additionally, if the filename ends in `.go`, rename it to
end in `.go.erb`.

```
<% unless version == 'ga' -%>
// test case + config here
<% end -%>
```

Otherwise, tests using a beta feature are written exactly the same as tests
using a GA one. Normally to use the beta provider, it's necessary to specify
`provider = google-beta`, as Terraform maps any resources prefixed with
`google_` to the `google` provider by default. However, inside the test
framework, the `google-beta` provider has been aliased as the `google` provider
and that is not necessary.

Note: You _may_ use version guards to test different configurations between the
GA and beta provider tests, but it's strongly recommended that you write
different test cases instead, even if they're slightly duplicative.

#### Promote a beta feature

"Promoting" a beta feature- making it available in the GA `google` provider when
the underlying feature or service has gone GA- requires removing the version
guards placed previously, so that the previously beta-only code is generated in
the `google` provider as well.

For all promotions, ensure that you remove the guards in:

* The documentation for the resource or field
* The test(s) for the resource or field.

For whole resource promotions, you'll generally only need to remove the file-level
guards and the guards on the resource registration in `provider.go.erb`.

For field promotions ensure that you remove the guards in:

* The Resource schema
* The Resource CRUD methods (for top level fields)
* The Resource Expanders and Flatteners (for nested fields)

When writing a changelog entry for a promotion, write it as if it was a new
field or resource, and suffix it with `(ga only)`. For example, if the
`google_container_cluster` resource was promoted to GA in your change:

```
\`\`\`release-note:new-resource
`google_container_cluster` (ga only)
\`\`\`
```

Alternatively, for field promotions, you may use "{{service}}: promoted
{{field}} in {{resource}} to GA", i.e.

```
\`\`\`release-note:enhancement
container: promoted `node_locations` field in google_container_cluster` to GA
\`\`\`
```
