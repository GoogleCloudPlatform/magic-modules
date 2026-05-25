---
subcategory: "Cloud Storage Control"
description: |-
  Summarize Storage Control Intelligence Findings in a project.
---

# google_storage_control_project_intelligence_findings_summary

Summarizes Cloud Storage intelligence findings in a specified project and location.

## Example Usage

```hcl
data "google_storage_control_project_intelligence_findings_summary" "summary" {
  project = "my-project-id"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Optional) The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

* `location` - (Optional) The location of the intelligence findings summary. Currently default value is global and users cannot use for input for now.

* `filter` - (Optional) The filter expression to apply.

* `resource_scope` - (Optional) The scope of the resources to include in the summary. Possible values are PARENT and PROJECT. Default value is PARENT.

## Attributes Reference

The following attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{location}}/intelligenceFindingsSummary`

* `finding_summaries` - A list of summaries for individual finding types. Structure is documented below.

<a name="nested_finding_summaries"></a>The `finding_summaries` block contains:

* `type` - The finding type.

* `category` - The category of the finding.

* `target_resource` - The target resource of the finding summary.

* `create_time` - The creation time of the finding summary.

* `update_time` - The last update time of the finding summary.

* `severity` - The severity of the finding.

* `summary_details` - Detailed summaries for the finding type. Structure is documented below.

<a name="nested_summary_details"></a>The `summary_details` block contains:

* `count` - The total count of findings for this summary slice.

* `percentage` - The percentage magnitude associated with this summary slice.

* `resource_type` - The resource type associated with the summary slice.

* `description` - A description explaining the summary slice.
