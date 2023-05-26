---
title: "Provider documentation"
weight: 60
---

# Provider documentation

The provider is [documented on HashiCorp's Terraform Registry](https://registry.terraform.io/providers/hashicorp/google/latest/docs), which includes information about individual resources and datasources, and includes guides to help users configure or upgrade the provider in their projects.

This document includes details about how provider documentation is used by the [Terraform Registry](https://registry.terraform.io/providers), how it is made in the Magic Modules repo, and tools you can use when editing documentation.

There are other pages under [How To](/magic-modules/docs/how-to) that describe _how_ to approach making additions to the documentation.

## How documentation is used by the Terraform Registry

The provider's documentation is rendered in the Terraform Registry using [markdown files](https://github.com/hashicorp/terraform-provider-google/tree/main/website/docs) that are packaged into each release. The Registry allows users to browse past versions of the documentation, for example [the documentation for v3.0.0](https://registry.terraform.io/providers/hashicorp/google/3.0.0/docs/guides/getting_started).

There are 4 types of documentation page. There's the [index page](https://github.com/hashicorp/terraform-provider-google/blob/main/website/docs/index.html.markdown), documentation for [resources](https://github.com/hashicorp/terraform-provider-google/tree/main/website/docs/r), documentation for [data sources](https://github.com/hashicorp/terraform-provider-google/tree/main/website/docs/d), and finally [guide pages](https://github.com/hashicorp/terraform-provider-google/tree/main/website/docs/guides).

For the Registry to successfully render documentation page, the markdown files in each provider release need to follow some requirements, described below.

### Directory structure

Files need to be saved in a [specific directory](https://developer.hashicorp.com/terraform/registry/providers/docs#directory-structure).

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
Note that the Google provider uses a legacy version of this requirement - a `website/docs/` folder.

### YAML frontmatter

Each file must include specific [YAML frontmatter](https://developer.hashicorp.com/terraform/registry/providers/docs#yaml-frontmatter).

- `subcategory` - for resource/data source pages -  determines where the link to the page is located in the left-side navigation.
- `page_title` - for guide pages -  sets the page title (as there isn't a resource to name it after). Here's [an example](https://github.com/hashicorp/terraform-provider-google/blob/46b96dcaec4e1563a5a0aff412e47896a3b72ea7/website/docs/guides/getting_started.html.markdown?plain=1#L2).
## What information documentation needs to include

[HashiCorp advice](https://developer.hashicorp.com/terraform/registry/providers/docs#headers) is to include these sections:

- Title and description
- Example Usage section
- Argument Reference section
- Attribute Reference section

In the Google provider we also include:
- Timeouts, describing configurable timeouts for a resource: [see example](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_address#timeouts)
- Import, how to import a resource into Terraform state: [see example](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_address#import)
- User Project Overrides, whether or not direct user project overrides are supported: [see example](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_address#user-project-overrides)

## How do you test documentation changes?

You can copy and paste markdown into the Registry's [Doc Preview Tool](https://registry.terraform.io/tools/doc-preview) to see how it will be rendered.

There currently isn't a way to preview how frontmatter will be used to create the left-side navigation menu.


## What formatting is available

You can expect markdown to be rendered in a similar way to READMEs in GitHub. When in doubt, you can test how some markdown will be rendered using the testing tool in the Doc Preview Tool mentioned above.

Something useful to be aware of are [callouts](https://developer.hashicorp.com/terraform/registry/providers/docs#callouts), that allow blue, yellow and red (warning) sections to be used for important information. To see how they're rendered, paste the markdown below into the Doc Preview Tool.

```markdown
-> **Note** This callout is blue

~> **Note** This callout is yellow

!> **Warning** This callout is red
```

## How to contribute to the provider documentation

### Handwritten documentation

Terraform Provider Google (TPG) contains handwritten documentation for handwritten resources and data sources. For guidance on updating handwritten documentation, see:
- [Update handwritten provider documentation](/magic-modules/docs/how-to/update-handwritten-documentation) 
- [Add documentation for a handwritten data source](/magic-modules/docs/how-to/add-handwritten-datasource-documentation)

### Generated documentation (mmv1)

The majority of resources in TPG are generated, and the information used to generate provider code is also used to generate documentation. For information about how documentation is generated, see:
- [Add and update MMv1 resource documentation](/magic-modules/docs/how-to/mmv1-resource-documentation)
