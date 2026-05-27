---
name: dev-overrides
description: "Configures Terraform CLI with developer overrides to test locally built provider binaries without full Go acceptance tests, enabling rapid plan/apply cycles."
---

# `dev-overrides`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your current roadblock or required task.

## Prerequisites
* You must have successfully compiled a local provider binary using the `generate-provider` skill.
* You need a directory containing the `.tf` configuration you want to test.

## Execution Steps

### 1. Verification
Ensure the provider binaries exist in the expected output paths.

#### Verify Built Binaries
```bash
# Verify the binaries were built (modify paths as necessary based on your host OS)
ls -lh $GOPATH/bin/terraform-provider-google*
```

### 2. The Core Commands
Create a `.terraformrc` file referencing the local binaries and invoke Terraform with it.

#### Setup Overrides and Test
```bash
# 1. Create the dev override file dynamically
cat << EOF > /tmp/tf-dev-override.tfrc
provider_installation {
  dev_overrides {
      "hashicorp/google"      = "${GOPATH}/bin"
      "hashicorp/google-beta" = "${GOPATH}/bin"
  }
  direct {}
}
EOF

# 2. Navigate to your test directory and run terraform using the override
# cd /path/to/test-directory
TF_CLI_CONFIG_FILE="/tmp/tf-dev-override.tfrc" terraform plan
# TF_CLI_CONFIG_FILE="/tmp/tf-dev-override.tfrc" terraform apply -auto-approve
```

### 3. Verification & Handoff
* Verify Terraform warns that "Provider development overrides are in effect."
* Once the behavior is validated via plan/apply outputs, proceed to write the formal Go acceptance tests.
