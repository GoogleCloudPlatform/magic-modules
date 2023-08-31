---
title: Handwritten docs style guide
weight: 200
---

# Handwritten documentation style guide

This document describes the style guide for handwritten documentation for resources and data sources. MMv1-based resources will automatically generate documentation that matches this style guide.

## File name and location

Handwritten documentation lives in:

- Data sources: [`magic-modules/third_party/terraform/website/docs/d/`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/website/docs/d)
- Resources: [`magic-modules/third_party/terraform/website/docs/r/`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/website/docs/r)

The name of the file is the name of the resource without a `google_` prefix. For example, for `google_compute_instance`, the file is called `compute_instance.html.markdown`

## YAML frontmatter

Every resource or datasource documentation page must include YAML frontmatter which sets `subcategory` (where the page will be displayed in the left sidebar).

```yaml
---
subcategory: Cloud Foobar
---
```

## Callouts

Use [callouts](https://developer.hashicorp.com/terraform/registry/providers/docs#callouts) for important information.

```markdown
-> **Note** This callout is blue

~> **Note** This callout is yellow

!> **Warning** This callout is red
```

## Sections

Every resource or datasource documentation page must include the following sections as described in Hashicorp's [Documenting Providers: Resource/Data Source Headers](https://developer.hashicorp.com/terraform/registry/providers/docs#resource-data-source-headers)

1. **Title and description.** Include a general description of the resource or data source and links to the official product usage documentation and REST API reference. Example:

   ```markdown
   # google\_cloud\_run\_locations

   Get Cloud Run locations available for a project. 

   To get more information about Cloud Run, see:

   * [API documentation](https://cloud.google.com/run/docs/reference/rest/v1/projects.locations)
   * How-to Guides
       * [Official Documentation](https://cloud.google.com/run/docs/)
   ```

   For beta-only resources or data sources, add the following snippet at the end of this section: 

   {{< tabs "resource-beta-warning" >}}
   {{< tab "Resource" >}}
   ```markdown
   ~> **Warning:** This resource is in beta, and should be used with the terraform-provider-google-beta provider.
   See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta resources.
   ```
   {{< /tab >}}
   {{< tab "Data source" >}}
   ```markdown
   ~> **Warning:** This data source is in beta, and should be used with the terraform-provider-google-beta provider.
   See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta resources.
   ```
   {{< /tab >}}
   {{< /tabs >}}
2. **Example Usage.** Include a minimal set of examples showing how to use the resource or data source.
3. **Argument Reference.** List settable fields on the datasource. For example:
   ```markdown
   ## Argument Reference

   The following arguments are supported:

   * `name` - (Required) Name of the resource.

   - - -

   * `project` - (Optional) The ID of the project in which the resource belongs. If it
       is not provided, the provider project is used.
   * `beta_field` - (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html)) This field is in beta.
   * `roles` - The MembershipRoles that apply to the Membership. Structure is [documented below](#nested_roles).

   <a name="nested_roles"></a>The `roles` block supports:

   * `name` - The name of the MembershipRole. One of OWNER, MANAGER, MEMBER.

   ```
4. **Attribute Reference.** List all output-only fields.
   ```markdown
   ## Attribute Reference

   In addition to the arguments listed above, the following computed attributes are exported:

   * `create_time` - (Output) The time when the repository was created.
   ```

   {{< hint "info" >}}
   **Note:** If a data source is a read-only version of a resource, instead provide a link to the resource documentation to avoid duplicating information:

   ```markdown
   ## Attribute Reference

   See [google_FOOBAR](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/FOOBAR#argument-reference) for details of the available attributes.
   ```
   {{< /hint >}}


If relevant, also include the following sections:

1. Timeouts ([example](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_address#timeouts))
2. Import ([example](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_address#import))
3. User Project Overrides ([example](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_address#user-project-overrides))