---
title: "Add a resource"
weight: 10
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

# Add a resource to an existing product

This page contains information about adding new resources to the `google` or `google-beta` Terraform providers using MMv1 and/or handwritten code.

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

### Create the resource configuration

1. In your cloned `magic-modules` repository, list the folders in `mmv1/products`.
   ```bash
   cd ~/magic-modules
   ls mmv1/products
   ```

   Output will look like:

   ```
   accessapproval          firebasehosting
   accesscontextmanager    firebasestorage
   activedirectory         firestore
   alloydb                 gameservices
   apigateway              gkebackup
   apigee                  gkehub
   appengine               gkehub2
   ...
   ```
2. Navigate to the folder your resource belongs to. For example, a new Apigee resource would be added to the `apigee` folder.

   ```bash
   cd PRODUCT
   ```

   Replace `PRODUCT` with the name of the folder.
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
   # Uncomment for beta resources
   # min_version: beta

   # Allows inserting styled markdown into the header of the resource's page
   # in the provider documentation.
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
   # exclude_import: true

   # Uncomment for resources that are primarily immutable (even if some
   # fields can be updated).
   # immutable: true

   # Uncomment to override one or more timeouts.
   # timeouts: !ruby/object:Api::Timeouts
   #   insert_minutes: 20 
   #   update_minutes: 20 
   #   delete_minutes: 20 

   # URL for creating a new resource, including query parameters.
   # Terraform field names enclosed in double curly braces will be replaced
   # with the field values from the resource.
   create_url: 'projects/{{project}}/locations/{{location}}/resourcenames?resourceId={{resource_id}}'
   # Uncomment to override the HTTP verb used to create a new resource.
   # Allowed values: :POST, :PUT, :PATCH. Default: :POST
   # create_verb: :POST

   # Uncomment to override the update URL for the resource. (Otherwise, the
   # self_link URL will be used.)
   # update_url: 'projects/{{project}}/locations/{{location}}/resourcenames/{{resource_id}}'
   # The HTTP verb used to update a resource. Allowed values: :POST, :PUT, :PATCH. Default: :PUT.
   update_verb: :PATCH
   # True if the resource should use an update mask for updates.
   update_mask: true

   # Uncomment to override the delete URL for the resource. (Otherwise, the
   # self_link URL will be used.)
   # delete_url: 'projects/{{project}}/locations/{{location}}/resourcenames/{{resource_id}}'
   # Uncomment to override the HTTP verb used to delete a resource.
   # Allowed values: :POST, :PUT, :PATCH, :DELETE. Default: :DELETE
   # delete_verb: :DELETE

   # Enable generation of code to handle API calls that return operations.
   autogen_async: true
   # Set parameters for handling operations returned by the API.
   async: !ruby/object:Api::OpAsync
     # Uncomment to override which API calls return operations.
     # Default: ['create', 'update', 'delete']
     # actions: ['create', 'update', 'delete']
     operation: !ruby/object:Api::OpAsync::Operation
       base_url: '{{op_id}}'
     # Uncomment if the completed operation's returned JSON will contain
     # a full resource in the "response" field
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

   Modify the template as needed to match the API resource's documented behavior. These are the most commonly-used fields. For a comprehensive reference, see [ResourceName.yaml reference]({{<ref "/reference/iam-policy-reference.md" >}}).

### Add fields

### `ResourceName.yaml`

To add a field, you'll append an entry to `properties` within `ResourceName.yaml`, such
as the following adding support for a `fooBar` field in the API:

```yaml
      - !ruby/object:Api::Type::String
        name: 'fooBar'
        min_version: beta
        immutable: true
        description: |
          The cloud.google.com description of this field.
        default_from_api: true
        custom_expand: 'templates/terraform/custom_expand/shortname_to_url.go.erb'
```

The first line of that snippet is the type of the field in the API, including
primitives like `String`, `Integer`, `Boolean`, `Double`. Additional special
types are detailed below in "Complex Types".

You can configure settings on the field that describe it in the API. Avoid
setting values to `false`, and omit them instead.

*   `description` is the cloud.google.com description of the field, and must be
    filled out manually.
*   `min_version` can be set to `beta` if the field is only available at public
    preview or beta. If it is GA, do not set a `min_version` value.
*   `required: true` indicates that a field is required. New top-level fields
    should not be considered required, as that is a breaking change. Subfields
    of newly-added optional fields can be added as required.
*   `immutable: true` indicates that a field can only be set when the API resource is
    created. Changing the field will force the resource to be recreated.
*   `output: true` indicates that a field is output-only in the API and cannot
    be configured by the user.
*   `default_value: {{value}}` adds a default value for the field. This should
    only be used if the default value is fixed in the API.
*   `send_empty_value: true` indicates that an explicit zero value should be
    sent to the API. This is useful when a value has a nonzero default in the
    API but the zero value for the type can be set. This is extremely common for
    booleans that default to `true`.
*   `update_verb` and `update_url` configure a custom update function for a
    field.`update_verb`should be set to a literal symbol for the type (such
    as :POST for `POST`) and the URL to a templated URL such
    as`projects/{{project}}/global/backendServices/{{name}}/setSecurityPolicy`.
*   `default_from_api: true` indicates that Terraform needs to handle a field
    specially. This is common for fields with complex defaults from the API that
    can't be represented with `default_value`. If a `default_from_api: true`
    field is set in a user's config, Terraform will treat it as an optional
    field, detecting drift and correcting drift. If it is not set, it will be
    treated as an output-only field.
*   `url_param_only: true` indicates that a field is not a part of the resource
    body (i.e. they will never be sent in request bodies or read from response
    bodies), and generally indicates that they are part of the URL. These fields
    can be referenced in template strings or custom code. Typically projects,
    regions, zones, locations, and parent fields will be annotated as parameters. 
*   `ignore_read: true` indicates that a value is not returned from an API, and
    Terraform should not look for it in API responses.
*   `custom_expand` and `custom_flatten` are custom functions to read/write a
    value from state. They refer to files holding function contents under
    [`mmv1/template/terraform/custom_expand`](https://github.com/GoogleCloudPlatform/magic-modules/tree/8728bc89c37d5033b530c7d7157bb43865d9df58/mmv1/templates/terraform/custom_expand)
    and
    [`mmv1/template/terraform/custom_flatten`](https://github.com/GoogleCloudPlatform/magic-modules/tree/8728bc89c37d5033b530c7d7157bb43865d9df58/mmv1/templates/terraform/custom_flatten)
    respectively.

## Field Configuration - Complex Types

### Enum

```yaml
          - !ruby/object:Api::Type::Enum
            name: 'metadata'
            description: |
              Can only be specified if VPC flow logging for this subnetwork is enabled.
              Configures whether metadata fields should be added to the reported VPC
              flow logs.
            values:
              - :EXCLUDE_ALL_METADATA
              - :INCLUDE_ALL_METADATA
              - :CUSTOM_METADATA
            default_value: :INCLUDE_ALL_METADATA
```

Enum values represent enums in the underlying API where it is valuable to
restrict the range of inputs to a fixed set of values. They are strings that
support a `values` key to define the array of possible values specified as
literal constants, and `default_value` should be specified as a literal constant
as well.

Most API enums should be typed as `String` instead- if the value will not be
fixed for >1 year, use a `String`.

### ResourceRef

```yaml
      - !ruby/object:Api::Type::ResourceRef
        name: 'urlMap'
        resource: 'UrlMap'
        imports: 'selfLink'
        description: |
          A reference to the UrlMap resource that defines the mapping from URL
          to the BackendService.
```

ResourceRefs are fields that reference other resource. They're most typical in
GCE, and making a field a `ResourceRef` instead of a `String` will cause
Terraform to allow switching between reference formats and versions safely. If a
field can refer to multiple resource types, use a `String` instead.

In a `ResourceRef`, `resource` and `imports` must be defined but Terraform
ignores those values. `resource` should be set to the resource kind, and
`imports` to `selfLink` within GCE and `name` elsewhere.

### Array

```yaml
      - !ruby/object:Api::Type::Array
        name: scopes
        item_type: Api::Type::String
        description: |
          The list of scopes to be made available for this service
          account.
```

```yaml
      - !ruby/object:Api::Type::Array
        name: 'instances'
        description: |
          A list of virtual machine instances serving this pool.
          They must live in zones contained in the same region as this pool.
        item_type: !ruby/object:Api::Type::ResourceRef
          name: 'instance'
          description: 'The instance being served by this pool.'
          resource: 'Instance'
          imports: 'selfLink'
```

Arrays refer to arrays in the underlying API, with their item being specified
through an `item_type` field. `item_type` accepts any type, although primitives
(String / Integer / Boolean) must be specified differently than other types as
shown above.

### NestedObject

```yaml
      - !ruby/object:Api::Type::NestedObject
        name: 'imageEncryptionKey'
        description: |
          Encrypts the image using a customer-supplied encryption key.
          After you encrypt an image with a customer-supplied key, you must
          provide the same key if you use the image later (e.g. to create a
          disk from the image)
        properties:
          - !ruby/object:Api::Type::String
            name: 'rawKey'
            description: |
              Specifies a 256-bit customer-supplied encryption key, encoded in
              RFC 4648 base64 to either encrypt or decrypt this resource.
          - !ruby/object:Api::Type::String
            name: 'sha256'
            output: true
            description: |
              The RFC 4648 base64 encoded SHA-256 hash of the
              customer-supplied encryption key that protects this resource.
```

NestedObject is an object in the JSON API, and contains a `properties` subfield
where a sub-properties array can be defined (including additional NestedObjects)

### KeyValuePairs (Labels / Annotations)

```yaml
      - !ruby/object:Api::Type::KeyValuePairs
        name: 'labels'
        description: Labels to apply to this address.  A list of key->value pairs.
```

KeyValuePairs is a special type to handle string -> string maps, such as GCE
`labels` fields. No extra configuration is required.

### Exactly One Of

To restrain a parent object to contain exactly one of its nested objects,
use `exactly_one_of` in the affected child objects.

Yaml:

```yaml
objects:
  - !ruby/object:Api::Resource
    name: 'Connection'
    ...
    properties:
      - !ruby/object:Api::Type::NestedObject
        name: 'cloudSql'
        exactly_one_of:
          - cloud_sql
          - aws
        properties:
          ...
      - !ruby/object:Api::Type::NestedObject
        name: aws
        exactly_one_of:
          - cloud_sql
          - aws
        properties:
          ...
```

## Advanced customization

### DiffSuppressFunc

Terraform allows fields to specify a
[DiffSuppressFunc](https://www.terraform.io/plugin/sdkv2/schemas/schema-behaviors#diffsuppressfunc),
which allows you to ignore diffs in cases where the two values are
**functionally identical**. This is generally useful when the API returns a
normalized value - for example by standardizing the case.

Note: The *preferred* behavior for APIs is to always return the value that the
user sent. DiffSuppressFunc is a workaround for APIs that don't.

The Terraform provider comes with a set of
["common diff suppress functions"](https://github.com/hashicorp/terraform-provider-google/blob/main/google/common_diff_suppress.go).
These fit frequent needs like ignoring whitespace at the beginning and end of a
string, or ignoring case differences.

If you need to define a custom diff specifically for your resource, you can do
so in a "constants" file, which is a `.go.erb` file in
[mmv1/templates/terraform/constants](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform/constants)
named `<product>_<resource>.erb`. You can then declare this custom code in
`ResourceName.yaml`:

```yaml
--- !ruby/object:Api::Resource
  name: ResourceName

  custom_code: !ruby/object:Provider::Terraform::CustomCode
    constants: templates/terraform/constants/product_resource_name.go.erb
```

Once you have chosen a DiffSuppressFunc, you can declare it as an override on
your resource:

```yaml
--- !ruby/object:Api::Resource
  name: ResourceName

  properties:
    - !ruby/object:Api::Type::String
      name: "myField"
      diff_suppress_func: 'tpgresource.CaseDiffSuppress'
```

The value of diff_suppress_func can be any valid DiffSuppressFunc, including the
result of a function call. For example:

```yaml
diff_suppress_func: 'tpgresource.OptionalPrefixSuppress("folders/")'
```

Please make sure to add thorough unit tests (in addition to basic integration
tests) for your diff suppress func.

Example: DomainMapping (DomainMappingLabelDiffSuppress)

-   [DomainMapping.yaml resource file](https://github.com/GoogleCloudPlatform/magic-modules/blob/67cef91ee76fc4871566f03e7caee1ef664f8aa0/mmv1/products/cloudrun/DomainMapping.yaml)
    -   [`custom_code`](https://github.com/GoogleCloudPlatform/magic-modules/blob/67cef91ee76fc4871566f03e7caee1ef664f8aa0/mmv1/products/cloudrun/DomainMapping.yaml#L40)
    -   [`diff_suppress_func: 'resourceBigQueryDatasetAccessRoleDiffSuppress'`](https://github.com/GoogleCloudPlatform/magic-modules/blob/67cef91ee76fc4871566f03e7caee1ef664f8aa0/mmv1/products/bigquery/DatasetAccess.yaml#L112)
-   [constants file](https://github.com/GoogleCloudPlatform/magic-modules/blob/67cef91ee76fc4871566f03e7caee1ef664f8aa0/mmv1/templates/terraform/constants/cloud_run_domain_mapping.go.erb)
-   [unit tests](https://github.com/GoogleCloudPlatform/magic-modules/blob/67cef91ee76fc4871566f03e7caee1ef664f8aa0/mmv1/third_party/terraform/tests/resource_cloud_run_domain_mapping_test.go#L9)

## Documentation

When adding a new MMv1 product or resource there are fields that you need to set within `product.yaml` and resource-level `ResourceName.yaml` that are specific to documentation for that resource. To learn more about MMv1 generated documentation and what YAML fields you need to pay attention to, see [Add and update MMv1 resource documentation](/magic-modules/docs/how-to/mmv1-resource-documentation).

# Beta features

When the underlying API of a feature is not final (i.e. a `vN` version like
`v1` or `v2`), is in preview, or the API has no SLO we add it to the
`google-beta` provider rather than the `google` provider, allowing users to
self-select for the stability level they are comfortable with.

In MMv1, a "version tag" can be annotated on resources, fields, resource iam
metadata and examples to control the stability level that a feature is available
at. Version tags are a specification of the minimum version a feature is
available at, written as `min_version: {{version}}`. This is only
specified when a feature is available at `beta`, and omitting a tag indicates
the target is generally available, or available at `ga`.

## Adding a beta resource

To add support for a beta resource in a preexisting product, ensure that a
`beta` level exists in the `versions` map in the `product.yaml` file for the product.
If one doesn't already exist, add it, setting the `base_url` to the appropriate
value. This is generally an API version including `beta`, such as `v1beta`, but
may be the same `base_url` as the `ga` entry for services that mix fields with
different stability levels within a single endpoint.

For example:

```diff
versions:
  - !ruby/object:Api::Product::Version
    name: ga
    base_url: https://compute.googleapis.com/compute/v1/
+  - !ruby/object:Api::Product::Version
+    name: beta
+    base_url: https://compute.googleapis.com/compute/beta/
```

If the product doesn't already exist, it's only necessary to add the `beta`
entry, i.e.:

```
versions:
  - !ruby/object:Api::Product::Version
    name: beta
    base_url: https://runtimeconfig.googleapis.com/v1beta1/
```

## Adding beta field(s)

NOTE: If a resource is already tagged as `min_version: beta`, follow the general
instructions for adding a field instead.

To add support for a beta field to a GA resource, ensure that the `beta` entry
already exists in the `versions` map for the product. See
[above](#adding-a-beta-resource) for details on doing so.

Next, add the field(s) as normal with a `min_version: beta` tag specified. In
the case of nested fields, only the highest-level field must be tagged, as
demonstrated below:

```diff
         - !ruby/object:Api::Type::NestedObject
            name: 'scaleDownControl'
+            min_version: beta
            description: |
              Defines scale down controls to reduce the risk of response latency
              and outages due to abrupt scale-in events
            properties:
            - !ruby/object:Api::Type::Integer
              name: 'timeWindowSec'
              description: |
                How long back autoscaling should look when computing recommendations
                to include directives regarding slower scale down, as described above.
...
```

## Promote a beta feature

In order to promote a beta feature to GA, remove the version tags previously set
on the feature or its tests. This will automatically make it available in the
`google` provider and remove the note that the feature is in beta in the provider
documentation.

For a resource, this typically means ensuring their removal:

* At the resource level, `ResourceName.yaml`
* On all resource examples, `examples` in `ResourceName.yaml` (unless some examples
use other beta resources or fields)
  * Additionally, for any modified examples, all `provider = google-beta`
annotations must be cleared

For a field, this typically means ensuring their removal:

* At the field level, `properties` in `ResourceName.yaml`
* On any resource examples where this was the last beta feature, `examples` in
`ResourceName.yaml`
  * Additionally, for any modified examples, all `provider = google-beta`
annotations must be cleared

If the feature was tested using handwritten tests, the version guards must be
removed, as described in the
[guidance for handwritten resources](third_party/terraform/README.md#promote-a-beta-feature).

When writing a changelog entry for a promotion, write it as if it was a new
field or resource, and suffix it with `(ga only)`. For example, if the
`google_container_cluster` resource was promoted to GA in your change:

~~~
```release-note:new-resource
`google_container_cluster` (ga only)
```
~~~

Alternatively, for field promotions, you may use `{{service}}: promoted
{{field}} in {{resource}} to GA`, i.e.

~~~
```release-note:enhancement
container: promoted `node_locations` field in google_container_cluster` to GA
```
~~~
{{< /tab >}}
{{< tab "Handwritten" >}}
> **Warning:** Handwritten resources are much more difficult to develop and maintain. Please try to make an MMv1 resource first. If you believe that is not possible, get explicit confirmation from the core team that it is okay to add a new handwritten resource before proceeding.

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

## Add IAM support

{{< tabs "IAM" >}}
{{< tab "MMv1" >}}

If the API resource supports IAM policies (indicated with `setIamPolicy` and `getIamPolicy` methods in the API documentation for the resource), add the following top-level block to `ResourceName.yaml`, replacing `IAM_GOES_HERE`.

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
  # Uncomment to override HTTP method for setIamPolicy.
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

  # Uncomment for beta-only IAM support
  # min_version: beta
```

Modify the template as needed to match the API resource's documented behavior. These are the most commonly-used fields. For a comprehensive reference, see [IAM policy YAML reference]({{<ref "/reference/iam-policy-reference.md" >}}).
{{< /tab >}}
{{< tab "Handwritten" >}}
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
