---
title: "Add a resource"
weight: 30
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
  - /docs/how-to
  - /how-to
---

# Add a resource

This page describes how to add a new resource to the `google` or `google-beta` Terraform provider using MMv1 and/or handwritten code.

For more information about types of resources and the generation process overall, see [How Magic Modules works]({{< ref "/get-started/how-magic-modules-works.md" >}}).

## Before you begin

1. Complete the [Generate the providers]({{< ref "/get-started/generate-providers" >}}) quickstart to set up your environment and your Google Cloud project.
2. Ensure that your `magic-modules`, `terraform-provider-google`, and `terraform-provider-google-beta` repositories are up to date.
   ```
   cd ~/magic-modules
   git checkout main && git clean -f . && git checkout -- . && git pull
   cd $GOPATH/src/github.com/hashicorp/terraform-provider-google
   git checkout main && git clean -f . && git checkout -- . && git pull
   cd $GOPATH/src/github.com/hashicorp/terraform-provider-google-beta
   git checkout main && git clean -f . && git checkout -- . && git pull
   ```

## Add the resource

{{< tabs "resource" >}}
{{< tab "MMv1" >}}
1. Open the [product folder]({{<ref "/get-started/how-magic-modules-works.md#mmv1" >}}) for the resource.
3. Create a new file for your new resource.

   ```bash
   touch RESOURCE_NAME.yaml
   ```

   Replace RESOURCE_NAME with the name of the API resource you are adding support for. For example, the [NatAddress](https://cloud.google.com/apigee/docs/reference/apis/apigee/rest/v1/organizations.instances.natAddresses) resource would be represented by `NatAddress.yaml`.
4. Open RESOURCE_NAME.yaml in an editor of your choice. Copy in the following template:
   ```yaml
   # Copyright 2023 Google Inc.
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

   --- !ruby/object:Api::Resource
   # API resource name
   name: 'ResourceName'
   # Resource description for the provider documentation.
   description: |
     RESOURCE_DESCRIPTION
   references: !ruby/object:Api::Resource::ReferenceLinks
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

   # Inserts styled markdown into the header of the resource's page in the
   # provider documentation.
   # docs:
   #   warning: WARNING_MARKDOWN
   #   note: NOTE_MARKDOWN

   # URL for the resource within the API domain. This should match the
   # resource's create URL (excluding any query parameters).
   # Terraform field names enclosed in double curly braces will be replaced
   # with the field values from the resource.
   base_url: 'projects/{{project}}/locations/{{location}}/resourcenames'
   # URL for a created resource within the API domain. This should match
   # the URL for getting a single resource.
   # Terraform field names enclosed in double curly braces will be replaced
   # with the field values from the resource.
   self_link: 'projects/{{project}}/locations/{{location}}/resourcenames/{{resource_id}}'
   # URL for importing a resource that already exists in GCP. In general
   # this will be a list containing self_link. If the resource cannot be read
   # from GCP, comment this out and set exclude_import: true instead.
   import_format: ['projects/{{project}}/locations/{{location}}/resourcenames/{{resource_id}}']
   # If true, import code will not be generated for this resource. This should
   # only be used if the resource cannot be read from GCP due to either API or
   # provider constraints.
   # exclude_import: true

   # If true, the resource and all its fields will be considered immutable - that
   # is, only creatable, not updatable. Individual fields can override this if
   # they have a custom update method in the API.
   # immutable: true

   # Overrides one or more timeouts, in minutes. All timeouts default to 20.
   # timeouts: !ruby/object:Api::Timeouts
   #   insert_minutes: 20 
   #   update_minutes: 20 
   #   delete_minutes: 20 

   # URL for creating a new resource, including query parameters.
   # Terraform field names enclosed in double curly braces will be replaced
   # with the field values from the resource.
   create_url: 'projects/{{project}}/locations/{{location}}/resourcenames?resourceId={{resource_id}}'
   # Overrides the HTTP verb used to create a new resource.
   # Allowed values: :POST, :PUT, :PATCH. Default: :POST
   # create_verb: :POST

   # Overrides the update URL for the resource. (Otherwise, the self_link URL
   # will be used.)
   # Terraform field names enclosed in double curly braces will be replaced
   # with the field values from the resource.
   # update_url: 'projects/{{project}}/locations/{{location}}/resourcenames/{{resource_id}}'
   # The HTTP verb used to update a resource. Allowed values: :POST, :PUT, :PATCH. Default: :PUT.
   update_verb: :PATCH
   # True if the resource should use an update mask for updates.
   update_mask: true

   # Overrides the delete URL for the resource. (Otherwise, the self_link URL
   # will be used.)
   # Terraform field names enclosed in double curly braces will be replaced
   # with the field values from the resource.
   # delete_url: 'projects/{{project}}/locations/{{location}}/resourcenames/{{resource_id}}'
   # Overrides the HTTP verb used to delete a resource.
   # Allowed values: :POST, :PUT, :PATCH, :DELETE. Default: :DELETE
   # delete_verb: :DELETE

   # Enables generation of code to handle API calls that return operations.
   autogen_async: true
   # Sets parameters for handling operations returned by the API.
   async: !ruby/object:Api::OpAsync
     # Overrides which API calls return operations. Default: ['create',
     # 'update', 'delete']
     # actions: ['create', 'update', 'delete']
     operation: !ruby/object:Api::OpAsync::Operation
       base_url: '{{op_id}}'
     # If true, the completed operation's returned JSON will be expected to
     # contain a full resource in the "response" field
     # result: !ruby/object:Api::OpAsync::Result
     #   resource_inside_response: true

   # All resources (of all kinds) that share a mutex value will block rather
   # than executing concurrent API requests.
   # Terraform field names enclosed in double curly braces will be replaced
   # with the field values from the resource.
   # mutex: RESOURCE_NAME/{{resource_id}}

   # IAM_GOES_HERE

   # EXAMPLES_GO_HERE

   parameters:
     - !ruby/object:Api::Type::String
       name: 'location'
       required: true
       immutable: true
       url_param_only: true
       description: |
         LOCATION_DESCRIPTION
     - !ruby/object:Api::Type::String
       name: 'resource_id'
       required: true
       immutable: true
       url_param_only: true
       description: |
         RESOURCE_ID_DESCRIPTION

   properties:
     # Fields go here
   ```

5. Modify the template as needed to match the API resource's documented behavior. These are the most commonly-used fields. For a comprehensive reference, see [ResourceName.yaml reference]({{<ref "/reference/resource-reference.md" >}}).
{{< /tab >}}
{{< tab "Handwritten" >}}
> **Warning:** Handwritten resources are much more difficult to develop and maintain. Please make an MMv1 resource instead. If you believe that is not possible, get explicit confirmation from the core team that it is okay to add a new handwritten resource before proceeding.

1. In your cloned `magic-modules` repository, list the folders in `mmv1/products`.
   ```bash
   cd ~/magic-modules
   ls mmv1/third_party/terraform/services
   ```

   Output will look like:

   ```
   accessapproval   containerattached   networksecurity
   alloydb          datalossprevention  privateca
   apigee           dataproc            pubsub
   appengine        dataprocmetastore   redis
   ...
   ```
2. Navigate to the folder your resource belongs to. For example, a new Apigee resource would be added to the `apigee` folder.

   ```bash
   cd PRODUCT
   ```

   Replace `PRODUCT` with the name of the folder.

   > **Tip:** Create a new folder if one does not exist. The name of the folder should match the API subdomain the resource will interact with.
3. Create a file for the resource code.

   ```bash
   touch resource_PRODUCT_RESOURCE_NAME.go
   ```

   Replace `RESOURCE_NAME` with the name of the API resource, split with `_` at any word breaks and lowercased. For example,
   `resource_alloydb_backup.go`.
4. Open the file in the editor of your choice and write the code for the
   resource.

   The `google` and `google-beta` providers use resources based on Terraform Plugin SDK v2. Please consult [Hashicorp's documentation](https://developer.hashicorp.com/terraform/plugin/sdkv2) for guidance on creating new resources.

   Alternately, create an MMv1 resource, [generate the providers]({{< ref "/get-started/generate-providers.md" >}}), and then copy the generated code as a starting point.
{{< /tab >}}
{{< /tabs >}}

## Add fields

{{< tabs "docs" >}}
{{< tab "MMv1" >}}
1. For each API field, copy the following template into the resource's `properties` attribute. Be sure to indent appropriately.

```yaml
# Supported types: String, Integer, Boolean, Double, Enum,
# ResourceRef (link to a GCP resource), KeyValuePairs (string -> string map),
# Array, and NestedObject
- !ruby/object:Api::Type::String
  name: 'API_FIELD_NAME'
  description: |
    MULTILINE_FIELD_DESCRIPTION
  # Marks the field (and any subfields) as beta-only. Ensure a beta version block
  # is present in provider.yaml. Do not use if an ancestor field (or the overall
  # resource) is already marked as beta-only.
  # min_version: beta

  # If true, the field (and any subfields) will be considered immutable - that
  # is, only settable on create. If unset or false, the field will still be
  # considered immutable if any ancestor field (or the overall resource) is
  # immutable, unless `update_url` is set.
  # immutable: true

  # If set, changes to the field's value will trigger a separate call to a
  # specific API method for updating the field's value. The field will not be
  # considered immutable even if an ancestor field (or the overall resource) is
  # immutable.
  # Terraform field names enclosed in double curly braces will be replaced
  # with the field values from the resource.
  # update_url: 'projects/{{project}}/locations/{{location}}/resourcenames/{{resource_id}}/setFieldName'

  # Overrides the verb used to update this specific field. Allowed values:
  # :POST, :PUT, :PATCH. Default: Resource's update_verb (which defaults to :PUT
  # if unset).
  # update_verb: :POST

  # If true, the field is required. If unset or false, the field is optional.
  # required: true

  # If true, the field is output-only - that is, it cannot be configured by the
  # user. If unset or false, the field is configurable.
  # output: true

  # If true, the provider will set the field's value in the resource state based
  # only on the user's configuration. If false or unset, the provider will set
  # the field's value in the resource state based on the API response. This
  # should only be used if the field cannot be read from GCP due to either API or
  # provider constraints.
  # ignore_read: true

  # Sets a client-side default value for the field. This should be used if the
  # API has a default value that applies in all cases and is stable. Removing
  # or changing a default value is a breaking change.
  # default_value: DEFAULT_VALUE

  # If true, and the user has not set a value for this field in their
  # configuration, the provider will accept any value returned from the API as
  # the value for the field. If false, the provider will show a diff if the API
  # returns a value and none is set in the user's config. This is useful for
  # complex or frequently-changed API-side defaults, but provides less useful
  # information at plan time than default_value.
  # default_from_api: true

  # If true, the provider will send values that Terraform core considers "empty"
  # (such as zero, false, or empty strings) to the API if set explicitly in the
  # user's config. If false, values that Terraform core considers "empty" will
  # cause the field to be omitted entirely from the API request. This is useful
  # for fields where the API would behave differently when receiving an "empty"
  # value vs no value for a particular field - for example, boolean fields that
  # have an API-side default of true.
  # send_empty_value: true

  # Specifies a list of fields (excluding the current field) that cannot be
  # specified at the same time as the current field. Must be set separately on
  # all listed fields.
  # conflicts:
  #   - field_one
  #   - nested_object.0.nested_field

  # Specifies a list of fields (including the current field) that cannot be
  # specified at the same time (but at least one of which must be set). Must be
  # set separately on all listed fields.
  # exactly_one_of:
  #   - field_one
  #   - nested_object.0.nested_field

  # Enum only. Sets allowed values as ruby "literal constants" (prefixed with a
  # colon). If the allowed values change frequently, use a String field instead
  # to allow better forwards-compatibility, and link to API documentation
  # stating the current allowed values in the String field's description. Do not
  # include UNSPECIFIED values in this list.
  # values:
  #   - :VALUE_ONE
  #   - :VALUE_TWO

  # Array only. Sets the expected type of the items in the array. Primitives
  # should use the name of the primitive class as a string; other types should
  # define the attributes of the nested type.
  # item_type: Api::Type::String
  # item_type: !ruby/object:Api::Type::Enum
  #   name: 'required but unused'
  #   description: 'required but unused'
  #   values:
  #     - :VALUE_ONE
  #     - :VALUE_TWO

  # NestedObject only. Defines fields nested inside the current field.
  # properties:
  #   - !ruby/object:Api::Type::String
  #     name: 'FIELD_NAME'
  #     description: |
  #       MULTI_LINE_FIELD_DESCRIPTION
```
2. Modify the field configuration according to the API documentation and behavior. These are the most commonly-used fields. For a comprehensive reference, see [Field reference]({{<ref "/reference/field-reference.md" >}})
{{< /tab >}}
{{< tab "Handwritten" >}}
{{< /tab >}}
{{< /tabs >}}


## Add IAM support

This section covers how to add IAM resources in Terraform if they are supported by a particular API resource (indicated by `setIamPolicy` and `getIamPolicy` methods in the API documentation for the resource).

{{< tabs "IAM" >}}
{{< tab "MMv1" >}}
1. Add the following top-level block to `ResourceName.yaml`, replacing `IAM_GOES_HERE`.

```yaml
iam_policy: !ruby/object:Api::Resource::IamPolicy
  # Name of the field on the terraform IAM resources which will reference
  # the parent resource. Update to match the parent resource's name.
  parent_resource_attribute: 'resource_name'
  # Character preceding setIamPolicy in the full URL for the API method.
  # Usually `:`
  method_name_separator: ':'
  # HTTP method for getIamPolicy. Usually :POST.
  # Allowed values: :GET, :POST. Default: :GET
  fetch_iam_policy_verb: :POST
  # Overrides the HTTP method for setIamPolicy.
  # Allowed values: :POST, :PUT. Default: :POST
  # set_iam_policy_verb: :POST

  # Must match the parent resource's import_format, but with the
  # parent_resource_attribute value substituted for the final field.
  import_format: [
    'projects/{{project}}/locations/{{location}}/resourcenames/{{resource_name}}'
  ]
  # Valid IAM role that can be set by generated tests. Default: 'roles/viewer'
  # allowed_iam_role: 'roles/viewer'

  # If IAM conditions are supported, set this attribute to indicate how the
  # conditions should be passed to the API. Allowed values: :QUERY_PARAM,
  # :REQUEST_BODY, :QUERY_PARAM_NESTED. Note: :QUERY_PARAM_NESTED should
  # only be used if the query param field contains a `.`
  # iam_conditions_request_type: :REQUEST_BODY

  # Marks IAM support as beta-only
  # min_version: beta
```

2. Modify the template as needed to match the API resource's documented behavior. These are the most commonly-used fields. For a comprehensive reference, see [IAM policy YAML reference]({{<ref "/reference/iam-policy-reference.md" >}}).
{{< /tab >}}
{{< tab "Handwritten" >}}
Handwritten resources should use MMv1-based IAM support.

1. Follow the MMv1 directions in [Add the resource]({{<ref "#add-the-resource" >}}) to create a skeleton ResourceName.yaml file for the handwritten resource, but set only the following top-level fields: `name`, `base_url` (set to URL of IAM parent resource), `self_link` (set to same value as `base_url`) `description` (required but unused), `id_format`, `import_format`, and `properties`.
2. Follow the MMv1 directions in [Add fields]({{<ref "#add-fields" >}}) to add only the fields used by base_url.
3. Follow the MMv1 directions in this section to add IAM support.
{{< /tab >}}
{{< /tabs >}}

## Add documentation

{{< tabs "docs" >}}
{{< tab "MMv1" >}}
Documentation is autogenerated for MMv1 resources.
{{< /tab >}}
{{< tab "Handwritten" >}}
{{< /tab >}}
{{< /tabs >}}

 It is a good idea to check the markdown changes when you [generate the providers]({{< ref "/get-started/generate-providers.md" >}}), especially if you are making lots of changes.

 You can copy and paste markdown into the Hashicorp Registry's [Doc Preview Tool](https://registry.terraform.io/tools/doc-preview) to see how it will be rendered.
