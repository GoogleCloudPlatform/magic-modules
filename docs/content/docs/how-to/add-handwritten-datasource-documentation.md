---
title: "Add documentation for a handwritten data source"
summary: "New handwritten datasources require new handwritten documentation to be created."
weight: 25
---

# Add documentation for a handwritten data source

{{< hint info >}}
**Note:** If you want to find information about documentation for a generated resource, look at the [MMv1 resource documentation](/magic-modules/docs/how-to/mmv1-resource-documentation) page instead. The information on this page will not be relevant for resources that have generated documentation.
{{< /hint >}}

## How provider documentation works

For general information about how provider documentation works, see [Provider Documentation](/magic-modules/docs/getting-started/provider-documentation).
That page contains information about how documentation should be structured and how you can test changes to documentation.

This page includes only instructions on how to add documentation for a new handwritten data source, with minimal background info.

## Finding the relevant file location

Handwritten documentation is located in the `website/docs` folder, shown below.

```
mmv1/third_party/terraform/website/docs/
├─ guides/
│  ├─ ...
├─ d/
│  ├─ ...
├─ r/
│  ├─ ...
├─ index.html.markdown
```

The subfolder `d` corresponds to data sources, and inside there is a file for each data source.

## Creating the new markdown file

Next you need to add the file for the new data source's documentation. Create a new file inside the `d` folder. The filename should be the name of the data source with the `google_` prefix removed. For example if you were adding a new data source `google_foobar` then you would need to create a new file with the name `foobar.html.markdown`.

## Adding contents

Pages within the documentation need to be consistent and contain the sections that users expect. Below is some guidance about the different sections to include and what their contents should be.

### Frontmatter

The top of the file needs to contain frontmatter, which is used to create the new documentation page's title and manage how the page is linked to in the documentation's sidebar navigation. It it not rendered but it is very important for the new page in the documentation to be generated and available to users.

You need to make sure your file's frontmatter includes:
- `subcategory` - This sets which section in the left-side navigation menu the new page is categorised into.
- `page_title` - This frontmatter is specific to files in the Guides section (`/website/docs/guides`) but we set it for all markdown files.
- `description` - A decription of the page.

For example, here is the frontmatter from `/website/docs/d/container_cluster.html.markdown` ([link to generated page](https://registry.terraform.io/providers/hashicorp/google/latest/docs/data-sources/container_cluster)):

```markdown
---
subcategory: "Kubernetes (Container) Engine"
page_title: "Google: google_container_cluster"
description: |-
  Get info about a Google Kubernetes Engine cluster.
---
```

### Page title and description

The next section in the markdown file is rendered as the first part of the page body.

It should contain:
- the page title, as an H1 header
- a description general information about the data source

The description can be as long or as short as necessary. The minimum information that's included in this section are links to official documentation and the API reference pages. Other guidance, warnings, or explanations of concepts can be included here. To create pronounced warning or info sections, see the [provider documentation](/magic-modules/docs/getting-started/provider-documentation/#what-formatting-is-available) page for info.

For example, here's the title and opening description of `/website/docs/d/cloud_run_locations.html.markdown` ([link to generated page](https://registry.terraform.io/providers/hashicorp/google/latest/docs/data-sources/cloud_run_locations)):

```markdown
# google\_cloud\_run\_locations

Get Cloud Run locations available for a project. 

To get more information about Cloud Run, see:

* [API documentation](https://cloud.google.com/run/docs/reference/rest/v1/projects.locations)
* How-to Guides
    * [Official Documentation](https://cloud.google.com/run/docs/)
```

### Example usage

The next section includes one or more examples showing how a user may use the data source in their Terrform configuration. The examples can be basic and show what the minimum number of arguments are required, or they can be used to demonstrate more complex usage of the data source if necessary.

For example, here's the example usage section from `/website/docs/d/kms_crypto_key.html.markdown` ([link to generated page](https://registry.terraform.io/providers/hashicorp/google/latest/docs/data-sources/kms_crypto_key)):

```markdown
## Example Usage

```hcl
data "google_kms_key_ring" "my_key_ring" {
  name     = "my-key-ring"
  location = "us-central1"
}

data "google_kms_crypto_key" "my_crypto_key" {
  name     = "my-crypto-key"
  key_ring = data.google_kms_key_ring.my_key_ring.id
}
```

### Argument reference

The argument reference section tells user about the fields in the schema which they can set via their Terraform configuration. The fields are listed using bullet points, and each field is marked with whether it is required or optional.

For example, this is the argument reference section from `/website/docs/d/compute_backend_bucket.html.markdown`:

```markdown
## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the resource.

- - -

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.
```

#### Nested blocks

If the data source contains a **nested block** we include the name of the nested block in the list of arguments and then link to a dedicated section below that describes all arguments within that block. If there are multiple levels of nesting in a block then this approach should be repeated.

For example in `/website/docs/d/cloud_identity_group_membership.html.markdown` the `roles` attribute is documented using this approach. To see it in action, view [the attribute in the documentation](https://registry.terraform.io/providers/hashicorp/google/latest/docs/data-sources/cloud_identity_group_membership#roles).

Inside the list of arguments there is this entry, which contains a link to a target elsewhere on the page:

```markdown
* `roles` - The MembershipRoles that apply to the Membership. Structure is [documented below](#nested_roles).
```

Under the list of arguments there are sections like this that contain an anchor tag that defines the target of the hyperlink above using the [name attribute](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/a#attr-name). These sections should be in the same order that the blocks' names are listed in the documentation.

```markdown
<a name="nested_roles"></a>The `roles` block supports:

* `name` - The name of the MembershipRole. One of OWNER, MANAGER, MEMBER.
```

### Attribute reference

Attributes are exported values that can be accessed from a data source (or resource) and are not set by the users Terraform configuration. They could be computed values, or values read from the API.

If the data source has the same attributes as its equivalent resource in the provider, you can just link to the documentation for the resource. This avoids duplicating information that already exists elsewhere.

For example the documentation for the `google_storage_bucket` data source links to the `google_storage_bucket` resource documentation ([see here](https://registry.terraform.io/providers/hashicorp/google/latest/docs/data-sources/storage_bucket#attributes-reference))


If there isn't an equivalent resource in the provider, or the attributes are different, then document the attributes in a bulleted list as usual. Nested blocks are documented in the way previously described above.