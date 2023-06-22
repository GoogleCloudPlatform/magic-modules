---
title: "Add a field"
weight: 40
---
# Add a field

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

## Add the field

{{< tabs "field" >}}
{{< tab "MMv1" >}}
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
3. Open the yaml configuration file for the resource you want to add a field to in an editor of your choice.
4. Copy the following template into the resource's `properties` attribute. Be sure to indent appropriately.

```yaml
# Supported types: String, Integer, Boolean, Double, Enum,
# ResourceRef (link to a GCP resource), KeyValuePairs (string -> string map), Array, and NestedObject
- !ruby/object:Api::Type::String
  name: 'API_FIELD_NAME'
  description: |
    MULTILINE_FIELD_DESCRIPTION
  # Mark the field (and any subfields) as beta-only. Ensure a beta version block
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
5. Modify the field configuration according to the API documentation and behavior. These are the most commonly-used fields. For a comprehensive reference, see [ResourceName.yaml field reference]({{<ref "/reference/field-reference.md" >}})



