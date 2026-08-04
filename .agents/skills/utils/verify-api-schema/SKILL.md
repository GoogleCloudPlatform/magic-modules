---
name: verify-api-schema
description: "Verify if an MMv1 YAML property matches the live GCP REST API JSON payload using the Discovery Documents."
---

# `verify-api-schema`

> **Note to AI Agents:** You MUST read the YAML frontmatter above first. Only read the rest of this file if the `description` matches your current roadblock or required task.

## Prerequisites
* You must be in the `magic-modules` root directory.
* You must have `curl` and `jq` installed.

## Execution Steps

### 1. Verification
Provide the exact bash commands the agent should run to verify the prerequisites.

#### Verify jq installation
```bash
jq --version
```

### 2. The Core Commands
Query the live GCP Discovery Documents directly. This is the ultimate source of truth for exact JSON payload names (including weird edge cases like double-plurals or typos in the API).

#### Query exact schema in v1 or beta
```bash
# For v1 Compute API schema
curl -s https://www.googleapis.com/discovery/v1/apis/compute/v1/rest | jq '.schemas.SecurityPolicyRule'

# For beta Compute API schema
curl -s https://www.googleapis.com/discovery/v1/apis/compute/beta/rest | jq '.schemas.SecurityPolicyRule'
```

Replace `SecurityPolicyRule` with the specific schema name you are tracking (e.g., `UrlMap`, `BackendService`).

### 3. Verification & Handoff
Instructions on how the agent should verify the command succeeded, and what workflow or rule it should return to next.

* Compare the JSON payload attributes (types, names, pluralization) with the Magic Modules YAML properties.
* If you find a discrepancy (e.g., `exclusions` in JSON vs `exclusion` in YAML), verify if `api_name` is being used to bridge the gap.
