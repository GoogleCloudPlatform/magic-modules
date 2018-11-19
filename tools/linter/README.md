# api.yaml linter

This linter lints api.yaml against the Discovery Docs.

Because api.yaml is handwritten, it will inevitably veer away from the Discovery docs
(the source of truth) or just be inputted incorrectly.

## Running the linter
The linter is powered by RSpec.

```
  rspec tools/linter/run.rb
```

## Adding new tests
All new tests should be added in `tools/linter/tests.rb`
New tests for properties should be added in `property_tests`.
New tests for resources should be added in `resource_tests`

## Adding new api.yamls
The linter will run against every api.yaml listed in `tools/linter/docs.yaml`
Each doc must contain a `url` key (URL of the discovery doc) and `filename` key
(path to api.yaml from magic-modules root)

## Tests

### Property not found
```
gdns ManagedZone labels should exist
```
Each property is tested to see if it exists.
If a property does not exist, RSpec will print the following message.
The message is formatted as "<product> <resource> <property> should exist"
