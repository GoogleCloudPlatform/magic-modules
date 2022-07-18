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
