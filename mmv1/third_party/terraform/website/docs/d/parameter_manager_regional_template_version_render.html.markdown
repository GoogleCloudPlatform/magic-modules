---
subcategory: "Parameter Manager"
description: |-
  Get information about a Parameter Manager Regional Template Version Render
---

# google_parameter_manager_regional_template_version_render

Renders a Parameter Manager Regional Template Version against a specified Parameter Version. For more information see the [official documentation](https://cloud.google.com/secret-manager/parameter-manager/docs/parameter-templates-overview) and [API](https://cloud.google.com/secret-manager/parameter-manager/docs/reference/rest/v1/projects.locations.templates.versions/render).

~> **Warning:** This datasource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](../guides/provider_versions.html.markdown) for more details on beta resources.

## Example Usage

```hcl
data "google_parameter_manager_regional_template_version_render" "basic" {
  provider            = google-beta
  template            = "my-regional-template"
  template_version_id = "v1"
  location            = "us-central1"
  parameter_version   = "projects/my-project/locations/us-central1/parameters/my-parameter/versions/v1"
}
```

## Argument Reference

The following arguments are supported:

* `template` - (Required) The Regional Template for obtaining the Regional Template Version.
    This can be either the reference of the template as in `projects/{{project}}/locations/{{location}}/templates/{{template_id}}` or only the name of the template as in `{{template_id}}`.

* `template_version_id` - (Required) The version of the regional template to render.

* `parameter_version` - (Required) The resource name of the ParameterVersion used to render the template version in the format `projects/{{project}}/locations/{{location}}/parameters/{{parameter_id}}/versions/{{parameter_version_id}}`. The parameter version's `location` must match the regional template's `location`.

- - -

* `location` - (Optional) Location of the Regional Template resource. It must be provided when the `template` field consists of only the name of the template.

* `project` - (Optional) The project for retrieving the Regional Template Version. If not specified, the provider project will be used.

## Attributes Reference

The following attributes are exported:

* `id` - An identifier for the data source with format `projects/{{project}}/locations/{{location}}/templates/{{template_id}}/versions/{{template_version_id}}?parameter_version={{parameter_version}}`.

* `name` - The resource name of the RegionalTemplateVersion. Format:
  `projects/{{project}}/locations/{{location}}/templates/{{template_id}}/versions/{{template_version_id}}`

* `template_format` - The format of the Regional Template Version (`TEMPLATE_FORMAT_YAML` or `TEMPLATE_FORMAT_JSON`).

* `template_data` - The template data. **Note**: This property is sensitive and will not be displayed in the plan.

* `rendered_template_data` - The server-generated rendered template payload with variables resolved from the parameter version. **Note**: This property is sensitive and will not be displayed in the plan.
