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
    - [Add or update a beta feature](#add-or-update-a-beta-feature)
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

#### Add or update a beta feature

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

#### Promote a beta feature
