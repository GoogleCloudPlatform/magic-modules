---
title: "Update a handwritten resource"
summary: "The Google providers for Terraform have a large number of handwritten go files, primarily for resources written before Magic Modules was used with them. Most handwritten files are expected to stay handwritten indefinitely, although conversion to a generator may be possible for a limited subset of them."
weight: 20
---

# Update a handwritten resource

The Google providers for Terraform have a large number of handwritten go files, primarily for resources written before Magic Modules was used with them. Most handwritten files are expected to stay handwritten indefinitely, although conversion to a generator may be possible for a limited subset of them.

We no longer accept new handwritten resources except in rare cases. However, understanding
how to edit and add to existing resources may be important for implementing new fields
or changing existing behavior.

To edit an existing resource to add a field there are four steps you'll go through.

1. Add the new field to the schema
1. Implement the respective flattener and/or expander for the new field
1. Add a testcase for the field or extend an existing one
1. Add documentation for the field to the respective markdown file


## Shared concepts

This section will serve as a point of reference for some shared concepts that
all handwritten files share. It's meant to be an introduction to our serialization
strategy and overview.

### Serialization strategy
The go files within the directory files are copied literally to their respective providers.
Our serialization methodology may seem complicated but for the case of handwritten resources its quite
simple. Editing the file will change its counterpart downstream.

### go and go.erb
Within the third party library you'll notice `go` and `go.erb` files.
Go files are native golang code while go.erb pass through ruby before
being serialized. The reason `go.erb` files exist are to protect certain
properties or fields from entering the `ga` provider. Thus you'll often see
lines like `<% unless version == 'ga' -%>` within the file. These blocks
will omit the enclosure from being output to the GA provider. In the
rare case where you are promoting all fields to `ga` and these blocks
are no longer needed you can remove the `.erb` extension.

### Create, Read, Update, Delete
As far as terraform schema is concerned these are the functions we
need to provide for terraform to be able to provision and delete
resources. In editing any fields you'll likely be adding functionality to
these functions or implementing them wholesale.


### Expanders and Flatteners
Expanders and flatteners are concepts created to simplify common patterns
and add conformity/code consistency. Essentially expanders are functions
used to segregate some translation from terraform representation to api representation.
We will use these to encapsulate this translation for blocks and/or complicated fields.
This allows our code to be concise and functionality to be readable and easily
apparent by separating these into their own functions. While expanders are used
for terraform to api, flatteners do just the opposite. Converting api to terraform.

Thus
* expanders - helper functions used for translating tf -> api representation
* flatteners - helper functions used for translating api -> tf representation

## Adding a new field to the schema
To add a new field you will have to compare an existing resource
to it's respective rest api documentation. Dependant on how the api implements
the field we will in almost all cases mirror the structure. For example if there is
an `enabled` field nested under a ``IdentityServiceConfig`` block we will mirror
this within the schema.

Thus the block for terraform to utilize this field would then be
```terraform
resource "x" "y" {
  identity_service_config{
    enabled = true
  }
}
```

You might think it convoluted to provide such a structure. Why not
simply provide a single `enable_identity_service_config`. One constant
has echoed through our mind as terraform developers through the years.
Api's are ever evolving. Mirroring the api gives us the best chance to stay
in step with that evolution. Therefore if `IdentityServiceConfig` is extended with new
parameters in the future we can cleanly encapsulate those into the existing block(s).

As far as providing the field itself, it's fairly straightforward. Mirror the field
from the api and look to the other fields and the schema type in the SDK to see
what's available and how to structure it. For the documentation, copying the documentation
from the rest api will be the usual practice.

If you are adding a field that is an ENUM from the api standpoint its best practice
to provide it as as a string to the provider. This field will likely have values
added to it by the api and this future proofs our provider to support new values without
haven't to make new additions. There will be rare exceptions, but generally its a good
practice.

## Implement the respective flattener and/or expander for the new field
Once you've added the field to the schema you will implement the corresponding
expander/flattener. See [expanders and flatters](#expanders-and-flatteners) for
more context on what these fields are used for. Essentially we will be editing the
read, create, and update operations to parse the schema and call the api to make
the changes to the state of the resource. Following existing patterns to create
this operation will be the best way to implement this. As there are many unique ways
to implement a given field we won't get into specifics.

For example a field in bigtable `google_sheets_options` containers two nested properties.
`range` and `skip_leading_rows`.

```golang
						// GoogleSheetsOptions: [Optional] Additional options if sourceFormat is set to GOOGLE_SHEETS.
						"google_sheets_options": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: `Additional options if source_format is set to "GOOGLE_SHEETS".`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									// Range: [Optional] Range of a sheet to query from. Only used when non-empty.
									// Typical format: !:
									"range": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Range of a sheet to query from. Only used when non-empty. At least one of range or skip_leading_rows must be set. Typical format: "sheet_name!top_left_cell_id:bottom_right_cell_id" For example: "sheet1!A1:B20"`,
										AtLeastOneOf: []string{
											"external_data_configuration.0.google_sheets_options.0.skip_leading_rows",
											"external_data_configuration.0.google_sheets_options.0.range",
										},
									},
									// SkipLeadingRows: [Optional] The number of rows at the top
									// of the sheet that BigQuery will skip when reading the data.
									"skip_leading_rows": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: `The number of rows at the top of the sheet that BigQuery will skip when reading the data. At least one of range or skip_leading_rows must be set.`,
										AtLeastOneOf: []string{
											"external_data_configuration.0.google_sheets_options.0.skip_leading_rows",
											"external_data_configuration.0.google_sheets_options.0.range",
										},
									},
								},
							},
						},
```

To simplify the implementation management of these fields can be delegated to expanders and
flatteners.

```golang
func expandGoogleSheetsOptions(configured interface{}) *bigquery.GoogleSheetsOptions {
	if len(configured.([]interface{})) == 0 {
		return nil
	}

	raw := configured.([]interface{})[0].(map[string]interface{})
	opts := &bigquery.GoogleSheetsOptions{}

	if v, ok := raw["range"]; ok {
		opts.Range = v.(string)
	}

	if v, ok := raw["skip_leading_rows"]; ok {
		opts.SkipLeadingRows = int64(v.(int))
	}
	return opts
}

func flattenGoogleSheetsOptions(opts *bigquery.GoogleSheetsOptions) []map[string]interface{} {
	result := map[string]interface{}{}

	if opts.Range != "" {
		result["range"] = opts.Range
	}

	if opts.SkipLeadingRows != 0 {
		result["skip_leading_rows"] = opts.SkipLeadingRows
	}

	return []map[string]interface{}{result}
}

```

## Add a testcase for the field or extend an existing one
Once your field has been implemented, go to the corresponding test file for
your resource and extend it. If your field is updatable it's good practice to
have a two step apply to ensure that the field *can* be updated. You'll notice
a lot of our tests have a import state verify directly after apply. These
steps are important as they will essentially attempt to import the resource
you just provisioned and *verify* that the field values are consistent with the
applied state. Please test all fields you've added to the provider. It's important
for us to ensure all fields are usable and workable.

## Add documentation for the field to the respective markdown file
See [Documentation](#documentation) for more information. Essentially you will
just be opening the corresponding markdown file and adding documentation, likely
copied from the rest api to the markdown file. Follow the existing patterns there-in.

## Beta features

When the underlying API of a feature is not final (i.e. a `vN` version like
`v1` or `v2`), is in preview, or the API has no SLO we add it to the
`google-beta` provider rather than the `google `provider, allowing users to
self-select for the stability level they are comfortable with.

Both the `google` and `google-beta` providers operate off of a shared codebase,
including for handwritten code. MMv1 allows us to write Go source files as
`.go.erb` templated source, and renders them as `.go` files in the downstream
repo.

The sole generator feature you need to be aware of is a "version guard", what is
effectively a preprocessor directive implemented using Embedded Ruby (ERB). A
version guard is a snippet used across this codebase by convention guarding
versioned code on an `unless` clause in a version check. For example:

```
	networkInterfaces, err := expandNetworkInterfaces(d, config)
	if err != nil {
		return nil, fmt.Errorf("Error creating network interfaces: %s", err)
	}
<% unless version == 'ga' -%>
	networkPerformanceConfig, err := expandNetworkPerformanceConfig(d, config)
	if err != nil {
		return nil, fmt.Errorf("Error creating network performance config: %s", err)
	}
<% end -%>
```

In the snippet above, the `networkInterfaces` field is generally available and
is not guarded. The `networkPerformanceConfig` field is only available at beta,
and is guarded by `unless version == ga`, and the guarded block is terminated by
an `end` statement.

If a service includes handwritten resources and mixed features or resources at
different versions, the client libraries used by each provider must be switched
using guards so that the stability level of the client library matches that of
the provider. For example, all handwritten Google Compute Engine (GCE) files
have the following guarded import:

```
<% if version == "ga" -%>
	"google.golang.org/api/compute/v1"
<% else -%>
	compute "google.golang.org/api/compute/v0.beta"
<% end -%>
```

This is not necessary for beta-only services, or for services that are generally
available in their entirety.

### Adding a beta resource

MMv1 doesn't selectively generate files, and any file that is beta-only must
have all of its contents guarded. When writing a resource that's available at
beta, start with the following snippet:

```

<% autogen_exception -%>
package google

<% unless version == 'ga' -%>

// Add the implementation of the file here

<% end ->
```

This will generate a blank file in the `google` provider. The resource file,
resource test file, and any service or resource specific utility files should be
guarded in this way.

Documentation **should not** be guarded. Instead, write it as normal including
the following snippet above the first example.

```
~> **Warning:** This resource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta resources.
```

When registering the resource in
[`provider.go.erb`](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/third_party/terraform/utils/provider.go.erb),
the entry should be guarded:

```diff
				"google_monitoring_dashboard":                  resourceMonitoringDashboard(),
+				<% unless version == 'ga' -%>
+				"google_project_service_identity":              resourceProjectServiceIdentity(),
+				<% end -%>
				"google_service_networking_connection":         resourceServiceNetworkingConnection(),
```

If this is a new service entirely, all service-specific entries like client
factory initialization should be guarded as well. However, new services should
generally be implemented using an alternate engine- either MMv1 or tpgtools/DCL.

### Adding beta field(s)

By contrast to beta resources, adding support for a beta field is much more
involved as small snippets of code throughout a resource file must be annotated.

To begin with, add the field to the `Schema` of the resource with guards, i.e.:

```
<% unless version == 'ga' -%>
			"network_performance_config": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: `Configures network performance settings for the instance. If not specified, the instance will be created with its default network performance configuration.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_egress_bandwidth_tier": {
							// ...
						},
					},
				},
			},
<% end -%>
```

Next, implement the d.Get/d.Set calls (for top level fields) or
expanders/flatteners for nested fields within guards.

Even if there are other guarded fields, it's recommended that you add distinct
guards per feature- that way, promotion (covered below) will be simpler as
you'll only need to remove lines rather that move them around.

### Promoting a beta feature

"Promoting" a beta feature- making it available in the GA `google` provider when
the underlying feature or service has gone GA- requires removing the version
guards placed previously, so that the previously beta-only code is generated in
the `google` provider as well.

For all promotions, ensure that you remove the guards in:

* The documentation for the resource or field
* The test(s) for the resource or field.

For whole resource promotions, you'll generally only need to remove the file-level
guards and the guards on the resource registration in `provider.go.erb`.

For field promotions ensure that you remove the guards in:

* The Resource schema
* The Resource CRUD methods (for top level fields)
* The Resource Expanders and Flatteners (for nested fields)

When writing a changelog entry for a promotion, write it as if it was a new
field or resource, and suffix it with `(ga only)`. For example, if the
`google_container_cluster` resource was promoted to GA in your change:

```
\`\`\`release-note:new-resource
`google_container_cluster` (ga only)
\`\`\`
```

Alternatively, for field promotions, you may use "{{service}}: promoted
{{field}} in {{resource}} to GA", i.e.

```
\`\`\`release-note:enhancement
container: promoted `node_locations` field in google_container_cluster` to GA
\`\`\`
```
