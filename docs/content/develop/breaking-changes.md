---
title: "Understanding Breaking Changes"
summary: "This page discusses provider versioning, handling of breaking changes, and rare exceptions within Terraform development."
weight: 12
---


# Breaking Changes and Provider Development

## Provider Versioning
As a provider is developed, resources are added, old resources are updated, and bugs are fixed. These changes are [bundled together as a release](https://github.com/hashicorp/terraform-provider-google/releases/tag/v4.32.0).

Releases are numerically defined with a version number in the form of `MAJOR.MINOR.PATCH`. Here, 'Patch' indicates bug fixes, 'Minor' represents new features, and 'Major' represents significant changes which would be breaking to the customer if committed. Once a release is published, the provider binary is copied to [Hashicorp's provider registry](https://registry.terraform.io/browse/providers).

## Customer Trust
Terraform authors can write modular configurations, aptly named modules. These are shared within organizations and [online](https://registry.terraform.io/browse/modules). Terraform configurations can specify [provider requirements](https://www.terraform.io/language/providers/requirements), including a [version constraint field](https://www.terraform.io/language/providers/requirements#version-constraints).

The configuration will then [tie these version constraints](https://www.terraform.io/language/expressions/version-constraints) to an approximate minor or exact full version. Maintaining trust and consistency on every `MINOR` or `MAJOR` version upgrade is critical.

If breaking changes are allowed within `MINOR` versions, trust in the provider will be eroded and module creators will not have confidence in provider stability. This diminished trust will eventually lead to customers investing or deploying less to GCP.

## Exceptions to Breaking Changes

While we strive to minimize breaking changes, there are certain exceptions where they become unavoidable. Notably, breaking changes are permissible when existing functionality is demonstrably broken due to an API or provider-level issue. In such cases, the change does not impact users negatively, since there is no instance where the Terraform provider is currently using the affected field or resource correctly.

For example, consider a situation involving the Google provider where an API endpoint we depend on changes its behavior or is deprecated. If the current implementation in the Terraform provider cannot adapt to this change and is thus broken, a breaking change would be necessary to restore the functionality.

## Breaking Changes

Having established that we want to avoid breaking changes, let's delve into what exactly constitutes a breaking change. We'll discuss this under four main categories and the rules within each.

### Provider Configuration Level Breakages

* Top-level behavior such as provider configuration and authentication changes.

<h4 id="provider-config-fundamental"> Changing fundamental provider behavior (Undetectable) </h4>

Including, but not limited to, modification of: authentication, environment variable usage, and constricting retry behavior.

### Resource List Level Breakages

* Resource/datasource naming conventions and entry differences.

<h4 id="resource-map-resource-removal-or-rename"> Removing or Renaming a Resource  </h4>

In Terraform, resources should be retained whenever possible. Removal of a resource will result in a configuration breakage wherever a dependency on that resource exists. Renaming or removing resources are functionally equivalent in terms of configuration breakages.

### Resource Level Breakages

* Individual resource breakages like field entry removals or behavior within a resource.

<h4 id="resource-schema-field-removal-or-rename"> Removing or Renaming a field  </h4>

In Terraform, fields should be retained whenever possible. Removal of a field will result in a configuration breakage wherever a dependency on that field exists. Renaming or removing a field are functionally equivalent in terms of configuration breakages.

<h4 id="resource-id"> Changing resource ID format (Undetectable) </h4>

Terraform uses resource ID to read resource state from the API. Modification of the ID format will break the ability to parse the IDs from any deployments.

<h4 id="resource-import-format"> Changing resource ID import format (Undetectable) </h4>

Automation external to our provider may rely on importing resources with a certain format. Removal or modification of existing formats will break this automation.

### Field Level Breakages

* Field-level conventions like attribute changes and naming conventions.

<h4 id="field-changing-type"> Changing Field Type  </h4>

While certain Field Type migrations may be supported at a technical level, it's a practice that we highly discourage. We see little value for these transitions vs the risk they impose.

<h4 id="field-optional-to-required"> Field becoming Required Field  </h4>

A field should not become 'Required' as existing configurations may not have this field defined, leading to broken configurations in sequential plans or applies.. If you are adding 'Required' to a field so a block won't remain empty, this can cause two issues. First, if it's a singular nested field, the block may gain more fields later and it's not clear whether the field is actually required so it may be misinterpreted by future contributors. Second, if users are defining empty blocks in existing configurations, this change will break them. Consider these points in admittance of this type of change.

<h4 id="field-becoming-computed"> Becoming a Computed only Field  </h4>

While a field can transition from 'Optional' to 'Optional+Computed', it should not change from 'Required' or 'Optional' to solely 'Computed'. This transition would effectively make the field read-only, thus breaking configs in sequential plan or applies where this field is defined in a configuration.

<h4 id="field-oc-to-c"> Optional and Computed to Optional  </h4>

A field should not transition from 'Computed + Optional' to 'Optional'. During a sequential apply, the Terraform state retains the previously computed value, which won't match the configuration, thus causing a discrepancy.

<h4 id="field-changing-default-value"> Adding or Changing a Default Value  </h4>

Adding a default value where one was not previously declared can work in a very limited subset of scenarios but is an all around 'not good' practice to engage in. Changing a default value will absolutely cause a breakage. The mechanism of break for both scenarios is current terraform deployments now gain a diff with sequential applies where the diff is the new or changed default value.

<h4 id="field-growing-min"> Growing Minimum Items  </h4>

'MinItems' cannot grow. Otherwise, existing terraform configurations that don't satisfy this rule will break.

<h4 id="field-shrinking-max"> Shrinking Maximum Items  </h4>

'MaxItems' cannot shrink. Otherwise, existing terraform configurations that don't satisfy this rule will break.

<h4 id="field-changing-data-format"> Changing field data format (Undetectable) </h4>

Modification of the data format (either by the API or manually) will cause a diff in subsequent plans if that field is not Computed. This results in a breakage. API breaking changes are out of scope with respect to provider responsibility but we may make changes in response to API breakages in some instances to provide more customer stability.

