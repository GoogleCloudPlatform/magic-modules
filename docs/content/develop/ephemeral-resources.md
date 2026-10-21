---
title: "Generate ephemeral resources (experimental)"
weight: 70
---

# Generate ephemeral resources

MMv1 can generate an independent Plugin Framework ephemeral resource alongside
a managed resource. The experimental generator currently supports GET-based
reads and synchronous actions with empty request bodies. General action inputs,
asynchronous operations, renewal, and cleanup require further generator support.

Enable generation in the product YAML:

```yaml
ephemeral_experimental:
  generate: true
  samples:
    - name: basic
      primary_resource_id: test
      steps:
        - name: my_ephemeral_example
          ephemeral_assertions:
            has_value: 'length(ephemeral.google_example.test.value) > 0'
```

Place the configuration in
`mmv1/templates/terraform/samples/services/<service>/my_ephemeral_example.tf.tmpl`.
Use an `ephemeral` block. The existing sample mechanisms for template variables,
test environment variables, steps, and minimum versions apply. For reads that
need a managed resource created first, use a setup step before the ephemeral
read, or an explicit Terraform dependency where appropriate.

Each test step must include named **boolean HCL expressions** in
`ephemeral_assertions`. Generated tests evaluate these through the existing echo
test provider and assert that each is true. Only assertion results are passed to
the echo provider; avoid passing credential values into its persisted state.
Ephemeral test steps do not generate import tests.

These samples are independent from the managed resource's `samples`: they produce
`ephemeral_<resource>_generated_test.go` and examples in
`website/docs/ephemeral-resources/<resource>.html.markdown`. Version filtering and
`exclude_test` / `exclude_basic_doc` work independently for the ephemeral samples.
An enabled resource with no test samples must explicitly set `exclude_test: true`;
document the reason beside that setting.

Documentation classifies URL identifiers as arguments and remaining fields as
computed results, matching the generated ephemeral schema. It does not copy
managed-resource import, update, or deletion behavior.

## Custom response transformations

An SDK `custom_flatten` function can depend on state and cannot automatically be
used by an ephemeral resource. Provide `custom_flatten_ephemeral` naming a Go
function in the service package or `fwutils`:

```yaml
custom_flatten_ephemeral: flattenEphemeralExampleValue
```

The signature is `func(interface{}) (interface{}, error)`. Return API-shaped data
that the recursive Framework converter can consume. The hook runs on missing
values too and applies to nested fields. Generation rejects opted-in resources
with SDK custom flatteners that lack an explicit ephemeral equivalent.

## Verification

Generate and build the complete provider after engine-template changes. Run the
generated `TestAccEphemeral...Generated` tests with current application-default
credentials and a test project. Also use a local Terraform configuration with a
development override when checking the complete plan/apply lifecycle and state
non-persistence.
