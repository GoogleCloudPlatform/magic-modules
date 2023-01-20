---
title: "Add and update MMv1 resource documentation"
summary: "Generated resources have generated documentation. This page describes the generation process and what YAML inputs are used."
weight: 13
---

# MMv1 resource documentation

A majority of the provider's documentation is generated using the same information that's used for generating the provider's Go code. For example, when adding a new field to a resource in the relevant `api.yaml` file we include a `description` field. This is used to set a field description in the resource's schema and is also used to document the field in generated markdown files.

## Updating an existing MMv1 resource's documentation

As a result of code generation, you often do not need to explicitly think about making updates to documentation. However it is a good idea to check the markdown changes when you generate the provider, especially if you are making lots of changes.

## Adding new documentation for a new MMv1 resource

More thought is required when adding a new resource from scratch. See the relevant section below for guidance about what YAML fields are needed to create complete documentation.

## How generated documentation is made

The Magic Modules compiler parses YAML files, creates Ruby objects that are populated with the data from those files, and then uses that data within template files to create produce the final markdown files.

The main template used for documentation is `mmv1/templates/terraform/resource.html.markdown.erb`. As an example of how it works, you can see that the main title of documentation pages are created by some [processing of name data](https://github.com/GoogleCloudPlatform/magic-modules/blob/a69f1150de76f2b2cd9d37faa6bd44c1fb8a460a/mmv1/templates/terraform/resource.html.markdown.erb#L41) for a given resource and then using Ruby string methods to [print an escaped version of the name](https://github.com/GoogleCloudPlatform/magic-modules/blob/a69f1150de76f2b2cd9d37faa6bd44c1fb8a460a/mmv1/templates/terraform/resource.html.markdown.erb#L58) into an H1 header.


## What YAML fields are used in the documentation

The YAML files in Magic Modules are used to generate both the provider's Go code and the markdown documentation. As a result of this, changes to documentation happen automatically while you make a change to a resource's implementation in the provider. Often it's possible to address an issue without ever needing to think about documentation changes.

However, if you are implementing a new product or resource from scratch, or making non-routine changes to documentation, then you will need to be aware of how the YAML fields are used. Especially as some are specific to documentation.

Below are descriptions of fields that are directly referenced in the documentation templates.

### Top level fields for a product

These fields are found at the top of `api.yaml` files, and describe an overall product (ruby/object:Api::Product).

| Field | Type | Relation to documentation | Example value |
| ----- | ---- |------------------------- | ------------- |
| `display_name` | string | Controls the value of `subcategory` in YAML frontmatter; determines where the link to the page appears in the left-side navigation menu. | [See example](https://github.com/GoogleCloudPlatform/magic-modules/blob/39bded78e3032328e972f8e5b5f37796a451440b/mmv1/products/accesscontextmanager/product.yaml#L16) |

### Top level fields within resources

These can be top-level properties of a resource (ruby/object:Api::Resource) or overrides for the resource (ruby/object:Overrides::Terraform::ResourceOverride), i.e. in `api.yaml` or `terraform.yaml`.

| Field | Type | Relation to documentation | Example value |
| ----- | ---- |------------------------- | ------------- |
| `has_self_link` | boolean | Boolean to indicate if a resource has a `self_link` attribute. If true, the attribute is included in the templated markdown | [See example](https://github.com/hashicorp/magic-modules/blob/44d348dc92c279992febd7132a88656417a2a86f/mmv1/products/bigquery/api.yaml#L32-L33) and the resulting [self_link attribute in docs](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/bigquery_dataset#self_link) |
| `import_format`| array of strings | Sets the identifiers that can be used to import a resource into Terraform state. Used to add multiple entries in the 'Import' section of documentation. | [See example](https://github.com/hashicorp/magic-modules/blob/44d348dc92c279992febd7132a88656417a2a86f/mmv1/products/apigateway/terraform.yaml#L23) that results in [an import section listing multiple options](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/api_gateway_api#import) |
| `description`| string | Sets the description of the resource at the top of the page | [See example](https://github.com/hashicorp/magic-modules/blob/44d348dc92c279992febd7132a88656417a2a86f/mmv1/products/bigquery/api.yaml#L34-L35) |
| `docs.note` | string | Text is templated into a callout block at the top of the page which is titled "Note" | [See example](https://github.com/hashicorp/magic-modules/blob/dc463cb5b459044bf6bb37a1d502ae8bb14e2127/mmv1/products/iamworkforcepool/terraform.yaml#L19-L21) and the resulting [callout in the docs](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/iam_workforce_pool)| 
| `docs.warning` | string | Text is templated into a callout block at the top of the page which is titled "Warning" | [See example](https://github.com/hashicorp/magic-modules/blob/dc463cb5b459044bf6bb37a1d502ae8bb14e2127/mmv1/products/bigquery/terraform.yaml#L94-L98) and the resulting [callout in the docs](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/bigquery_dataset)|
| * `docs.optional_properties` | string | Used to append extra content to the bulleted list describing optional properties for a resource. | [See example](https://github.com/GoogleCloudPlatform/magic-modules/blob/1589b882611cceafdf2615ca74cc215c327ef141/mmv1/products/cloudiot/terraform.yaml#L28-L61) |
| * `docs.required_properties` | string | Used to append extra content to the bulleted list describing required properties for a resource. | There are no examples of this currently in use. |
| * `docs.attributes` | string | Used to append extra content to the bulleted list describing attributes for a resource. | There is currently only one example of this field in use, [here](https://github.com/hashicorp/magic-modules/blob/dacfb793fec55a9a2929be00b0cfa8f6cc5f1f88/mmv1/products/iap/terraform.yaml#L224-L226). |
| `min_version`| string | If set to `beta`, the template includes a beta warning at the start of the documentation | [See example in the official docs](https://registry.terraform.io/providers/hashicorp/google-beta/latest/docs/resources/alloydb_backup) |
| `references.api`| string | Sets the URL used in generated links to documentation | [See example](https://github.com/hashicorp/magic-modules/blob/44d348dc92c279992febd7132a88656417a2a86f/mmv1/products/bigquery/api.yaml#L39) |
| `references.guides`| hash | A set of key-value pairs where the key is text to be rendered and the value is the URL the text links to | [See example](https://github.com/hashicorp/magic-modules/blob/44d348dc92c279992febd7132a88656417a2a86f/mmv1/products/bigquery/api.yaml#L38)  |
| `supports_indirect_user_project_override`| boolean | This is the explicit way to make sure the 'User Project Overrides' section is shown in the resource's documentation | [See example](https://github.com/hashicorp/magic-modules/blob/44d348dc92c279992febd7132a88656417a2a86f/mmv1/products/datacatalog/terraform.yaml#L43) and the [resulting section in the docs](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/data_catalog_entry#user-project-overrides)  |

\* = Avoid unless absolutely necessary.

### Fields set within resource properties

The main thing to focus on for properties is to provide an adequate description. Typically we copy and paste the descriptions from the official API reference, making sure to pick the richest description and include any links.

Updates to other fields, like `output` and `required`, are easier to troubleshoot via acceptance testing but are also relevant to the produced documentation.

| Field | Type | Relation to documentation |
| ----- | ---- | ------------------------- |
| `deprecated` | boolean | Controls if `(Deprecated)` will be shown next to the argument name |
| `description` | string | The description shown in for a given argument or attribute |
| `name` | string | The name displayed in a list of arguments or attributes |
| `output` | boolean | Used to prevent a field being presented as an argument in the docs |
| `required` | boolean | Controls if `(Required)` will be shown next to the argument name |
| `sensitive` | boolean | Is used to help create a list of sensitive values. This list is used in a warning callout about sensitive values in state at the top of the page |
| `skip_docs_values` | boolean | Controls if a default value will be shown, or not, in the docs |
