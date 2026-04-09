---
title: "Add a resource"
weight: 20
aliases:
  - /docs/how-to/add-mmv1-resource
  - /how-to/add-mmv1-resource
  - /develop/add-mmv1-resource
  - /docs/how-to/mmv1-resource-documentation
  - /how-to/mmv1-resource-documentation
  - /develop/mmv1-resource-documentation
  - /docs/how-to/add-mmv1-iam
  - /how-to/add-mmv1-iam
  - /develop/add-mmv1-iam
  - /docs/how-to/update-handwritten-resource
  - /how-to/update-handwritten-resource
  - /develop/update-handwritten-resource
  - /docs/how-to/update-handwritten-documentation
  - /how-to/update-handwritten-documentation
  - /develop/update-handwritten-documentation
  - /docs/how-to
  - /how-to
  - /docs/getting-started/provider-documentation
  - /getting-started/provider-documentation
  - /develop/resource
---

# Add a resource

This page describes how to add a new resource to the `google` or `google-beta` Terraform provider using MMv1 and/or handwritten code.

For more information about types of resources and the generation process overall, see [How Magic Modules works]({{< ref "/" >}}).

## Before you begin

1. Complete the steps in [Set up your development environment]({{< ref "/develop/set-up-dev-environment" >}}) to set up your environment and your Google Cloud project.
1. Ensure that your `magic-modules`, `terraform-provider-google`, and `terraform-provider-google-beta` repositories are up to date.
    ```bash
    cd ~/magic-modules
    git checkout main && git clean -f . && git checkout -- . && git pull
    cd $GOPATH/src/github.com/hashicorp/terraform-provider-google
    git checkout main && git clean -f . && git checkout -- . && git pull
    cd $GOPATH/src/github.com/hashicorp/terraform-provider-google-beta
    git checkout main && git clean -f . && git checkout -- . && git pull
    ```

## Add a resource

{{% tabs "resource" %}}
{{< tab "MMv1" >}}
1. Using an editor of your choice, in the appropriate [product folder]({{<ref "/#mmv1" >}}), create a file called `RESOURCE_NAME.yaml`. Replace `RESOURCE_NAME` with the name of the API resource you are adding support for. For example, a configuration file for [NatAddress](https://cloud.google.com/apigee/docs/reference/apis/apigee/rest/v1/organizations.instances.natAddresses) would be called `NatAddress.yaml`.
2. Copy the following template into the new file:
   ```yaml
   # Copyright {{< now >}} Google Inc.
   # Licensed under the Apache License, Version 2.0 (the "License");
   # you may not use this file except in compliance with the License.
   # You may obtain a copy of the License at
   #
   #     http://www.apache.org/licenses/LICENSE-2.0
   #
   # Unless required by applicable law or agreed to in writing, software
   # distributed under the License is distributed on an "AS IS" BASIS,
   # WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   # See the License for the specific language governing permissions and
   # limitations under the License.

   ---
   # API resource name
   name: 'ResourceName'
   # Resource description for the provider documentation.
   description: |
     RESOURCE_DESCRIPTION
   references:
     guides:
      # Link to quickstart in the API's Guides section. For example:
      # 'Create and connect to a database': 'https://cloud.google.com/alloydb/docs/quickstart/create-and-connect'
       'QUICKSTART_TITLE': 'QUICKSTART_URL'
     # Link to the REST API reference for the resource. For example,
     # https://cloud.google.com/alloydb/docs/reference/rest/v1/projects.locations.backups
     api: 'API_REFERENCE_URL'
   # Marks the resource as beta-only. Ensure a beta version block is present in
   # provider.yaml.
   # min_version: beta

   # URL for the resource's standard List method. https://google.aip.dev/132
   # Terraform field names enclosed in double curly braces are replaced with
   # the field values from the resource at runtime.
   base_url: 'projects/{{project}}/locations/{{location}}/resourcenames'
   # URL for the resource's standard Get method. https://google.aip.dev/131
   # Terraform field names enclosed in double curly braces are replaced with
   # the field values from the resource at runtime.
   self_link: 'projects/{{project}}/locations/{{location}}/resourcenames/{{name}}'

   # If true, the resource and all its fields are considered immutable - that is,
   # only creatable, not updatable. Individual fields can override this if they
   # have a custom update method in the API.
   immutable: true

   # URL for the resource's standard Create method, including query parameters.
   # https://google.aip.dev/133
   # Terraform field names enclosed in double curly braces are replaced with
   # the field values from the resource at runtime.
   create_url: 'projects/{{project}}/locations/{{location}}/resourcenames?resourceId={{name}}'

   # Overrides the URL for the resource's standard Update method. (If unset, the
   # self_link URL is used by default.) https://google.aip.dev/134
   # Terraform field names enclosed in double curly braces are replaced with
   # the field values from the resource at runtime.
   # update_url: 'projects/{{project}}/locations/{{location}}/resourcenames/{{name}}'
   # The HTTP verb used to update a resource. Allowed values: :POST, :PUT, :PATCH. Default: :PUT.
   update_verb: 'PATCH'
   # If true, the resource sets an `updateMask` query parameter listing modified
   # fields when updating the resource. If false, it does not.
   update_mask: true

   # If true, code for handling long-running operations is generated along with
   # the resource. If false, that code is not generated.
   autogen_async: true
   # Sets parameters for handling operations returned by the API.
   async:
     # Overrides which API calls return operations. Default: ['create',
     # 'update', 'delete']
     # actions: ['create', 'update', 'delete']
     operation:
       base_url: '{{op_id}}'

   parameters:
     - name: 'location'
       type: String
       required: true
       immutable: true
       url_param_only: true
       description: |
         LOCATION_DESCRIPTION
     - name: 'name'
       type: String
       required: true
       immutable: true
       url_param_only: true
       description: |
         NAME_DESCRIPTION

   properties:
     # Fields go here
   ```

3. Modify the template as needed to match the API resource's documented behavior.
4. Delete all remaining comments in the resource configuration (including attribute descriptions) that were copied from the above template.

> **Note:** The template includes the most commonly-used fields. For a comprehensive reference, see [MMv1 resource reference ↗]({{<ref "/reference/resource" >}}).
{{< /tab >}}
{{< tab "Handwritten" >}}
> **Warning:** Handwritten resources are more difficult to develop and maintain. New handwritten resources will only be accepted if implementing the resource in MMv1 would require entirely overriding two or more CRUD methods.

1. Add the resource in MMv1.
2. [Generate the beta provider]({{< ref "/develop/generate-providers" >}})
3. From the beta provider, copy the files generated for the resource to the following locations:
   - Resource: Copy to the appropriate service folder inside [`magic-modules/mmv1/third_party/terraform/services`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/services)
   - Documentation: [`magic-modules/mmv1/third_party/terraform/website/docs/r`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/website/docs/r)
   - Tests: Copy to the appropriate service folder inside [`magic-modules/mmv1/third_party/terraform/services`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/services), and remove `_generated` from the filename
   - Sweepers: Put to the appropriate service folder inside [`magic-modules/mmv1/third_party/terraform/services`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/services), and add `_sweeper` suffix to the filename
   - Metadata: Copy `*_meta.yaml` to the appropriate service folder inside [`magic-modules/mmv1/third_party/terraform/services`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/services), and remove `_generated` from the filename. For more information, see [Metadata (meta.yaml) reference]({{< ref "/reference/metadata" >}}).
4. Modify the Go code as needed.
   - Replace all occurrences of `github.com/hashicorp/terraform-provider-google-beta/google-beta` with `github.com/hashicorp/terraform-provider-google/google`
   - Remove the `Example` suffix from all test function names.
   - Remove the comments at the top of the file.
   - If any of the added Go code (including any imports) is beta-only, change the file suffix to `.go.tmpl` and wrap the beta-only code in a version guard: `{{- if ne $.TargetVersionName "ga" -}}...{{- else }}...{{- end }}`.
     - If the whole resource is beta-only, wrap everything except package declarations. Otherwise, individually wrap each logically-related block of code in a version guard (field, test, etc) rather than grouping adjacent version-guarded sections - it's easier to read and easier to modify as things move out of beta.
5. Register the resource `handwrittenResources` in [`magic-modules/mmv1/third_party/terraform/provider/provider_mmv1_resources.go.tmpl`](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/third_party/terraform/provider/provider_mmv1_resources.go.tmpl)
   - Add a version guard for any beta-only resources.
6. Optional: Complete other handwritten tasks that require the MMv1 configuration file.
    - [Add resource tests]({{< ref "/test/test" >}})
    - [Add IAM support]({{<ref "/develop/add-iam-support" >}})
7. Delete the MMv1 configuration file.
{{< /tab >}}
{{% /tabs %}}

## What's next?

+ [Add a field to an existing resource]({{< ref "/develop/add-fields" >}})
+ [Add IAM support]({{< ref "/develop/add-iam-support" >}})
+ [Add documentation]({{< ref "/document/add-documentation" >}})
+ [Add custom resource code]({{< ref "/develop/custom-code" >}})
+ [Add tests]({{< ref "/test/test" >}})
+ [Run tests]({{< ref "/test/run-tests" >}})
