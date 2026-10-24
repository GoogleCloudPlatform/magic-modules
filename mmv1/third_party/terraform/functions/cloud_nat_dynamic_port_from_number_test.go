// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package functions_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccDefaultFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"google": providerserver.NewProtocol6WithError(New("test")()),
	}
)

func TestCloudNatDynamicPortFromNumber_Equals(t *testing.T) {
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.8.0"))),
		},
		ProtoV6ProviderFactories: testAccDefaultFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				output "test" {
					value = provider::google::cloud_nat_dynamic_port_from_number("min", 8192)
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckOutput("test", "8192"),
				),
			},
		},
	})
}

func TestCloudNatDynamicPortFromNumber_MoreMin(t *testing.T) {
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.8.0"))),
		},
		ProtoV6ProviderFactories: testAccDefaultFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				output "test" {
					value = provider::google::cloud_nat_dynamic_port_from_number("min", 65537)
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckOutput("test", "32768"),
				),
			},
		},
	})
}

func TestCloudNatDynamicPortFromNumber_MoreMax(t *testing.T) {
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.8.0"))),
		},
		ProtoV6ProviderFactories: testAccDefaultFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				output "test" {
					value = provider::google::cloud_nat_dynamic_port_from_number("max", 65537)
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckOutput("test", "65536"),
				),
			},
		},
	})
}

func TestCloudNatDynamicPortFromNumber_Error(t *testing.T) {
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.8.0"))),
		},
		ProtoV6ProviderFactories: testAccDefaultFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				output "test" {
					value = provider::google::cloud_nat_dynamic_port_from_number("toto", 42)
				}
				`,
				// The parameter does not enable AllowNullValue
				ExpectError: regexp.MustCompile(`Invalid port_type`),
			},
		},
	})
}
