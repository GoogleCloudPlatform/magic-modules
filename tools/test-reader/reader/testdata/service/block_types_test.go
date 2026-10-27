package service_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/acctest"
)

func TestAccBlockTypes(t *testing.T) {
	acctest.VcrTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccBlockTypes(),
			},
		},
	})
}

// Configs can contain any top-level block type Terraform supports, including
// ones added after this tool was written. None of them should prevent the
// blocks in the same config from being read.
//
// Every block type here uses the same type label on purpose: a data, ephemeral
// or list block frequently shares its type label with a resource block without
// sharing its schema, so the attributes of each have to be recorded separately.
func testAccBlockTypes() string {
	return `
terraform {
  required_version = ">= 1.14.0"
}

variable "var_one" {
  type = string
}

locals {
  local_one = "value-one"
}

resource "block_types_resource" "resource" {
  field_one = var.var_one
}

data "block_types_resource" "data" {
  field_two = "value-two"
}

ephemeral "block_types_resource" "ephemeral" {
  field_three = "value-three"
}

list "block_types_resource" "list_query" {
  provider = google
  config {}
}

check "check_one" {
  assert {
    condition     = local.local_one != ""
    error_message = "local_one must be set"
  }
}

output "output_one" {
  value = local.local_one
}
`
}
