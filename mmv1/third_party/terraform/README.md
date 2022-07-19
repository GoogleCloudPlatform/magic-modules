# Handwritten

## Overview

The Google providers for Terraform have a large number of handwritten go files, primarily for resources written before Magic Modules was used with them. Most handwritten files are expected to stay handwritten indefinitely, although conversion to a generator may be possible for a limited subset of them.

Handwritten resources like [google_container_cluster](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/container_cluster) can be identified if they have source code present under the [mmv1/third_party/terraform/resources](./resources) folder or by the absence of the `AUTO GENERATED CODE` header in their [Go source](https://github.com/hashicorp/terraform-provider-google/blob/main/google/resource_container_cluster.go) in the downstream repositories. Handwritten datasources should be under the [mmv1/third_party/terraform/data_sources](./data_sources) folder, tests under the [mmv1/third_party/terraform/tests](./tests) folder and web documentaion under the [mmv1/third_party/terraform/website](./website) folder.

## Table of Contents
- [Contributing](#contributing)
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

### Resource

### Datasource

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
