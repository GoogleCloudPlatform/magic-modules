---
title: "Update handwritten provider documentation"
summary: "Handwritten resources and datasources have handwritten documentation that needs to be updated in PRs."
weight: 24
---

# Update handwritten provider documentation (for handwritten resource or datasource)

{{< hint info >}}
**Note:** If you want to find information about documentation for a generated resource, look at the [MMv1 resource documentation](/magic-modules/docs/how-to/mmv1-resource-documentation) page instead. The information on this page will not be relevant for resources that have generated documentation.

{{< /hint >}}

## How provider documentation works

For general information about how provider documentation works, see [Provider Documentation](/magic-modules/docs/getting-started/provider-documentation).
That page contains information about how documentation should be structured and how you can test changes to documentation.

This page includes only instructions on how to update the documentation for a handwritten resource or data source, with minimal background info.

## Finding the relevant file

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


The subfolder `d` corresponds to data sources, and `r` corresponds to resources, and each file inside creates a page in the official provider documentation. For example, if you needed to update existing documentation for the `google_compute_instance` resource you should search in the `r` folder for a file with a name starting `compute_instance` (i.e. the resource name with the provider name removed from the start). The file [/docs/r/compute_instance.html.markdown](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/third_party/terraform/website/docs/r/compute_instance.html.markdown) is used to produce the [page for the `google_compute_instance` resource](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_instance) in the official provider documentation.

## Making changes

After finding the file you need, make the changes required for the issue you are working on.

## Testing your changes

Next, you should test your changes to the file. To do this, you can copy and paste the markdown into the [Doc Preview Tool](https://registry.terraform.io/tools/doc-preview) on the Registry website. This will help you identify any malformed markdown and check that it is rendered in the way you expect.

Once you are satified, include the markdown changes in your PR in the Magic Modules repo. When the downstream is generated from Magic Modules your handwritten files will be copied into the correct location within the `website/docs` folder in `terraform-provider-google` or `terraform-provider-google-beta`.