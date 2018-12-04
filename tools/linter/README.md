# api.yaml linter

This linter lints api.yaml against the Discovery Docs.

Because api.yaml is handwritten, it will inevitably veer away from the Discovery docs
(the source of truth) or just be inputted incorrectly.

## Running the linter
The linter is powered by RSpec.

```
  rspec tools/linter/run.rb
```

Tests are divided into two parts: property tests and resource tests. Each is
tagged with :property or :resource To run only property tests, do the
following:

```
  rspec tools/linter/run.rb --tag property
```

## Getting Results as a CSV
RSpec uses formatters to create the output.
We have a custom formatter that works only with property tests to show which tests do/do not exist.
The formatter is located at `tools/linter/spreadsheet/csv_formatter.rb`

To get the property tests as a CSV, do the following:

```
  rspec tools/linter/spreadsheet.rb
```

The file will be outputted at `output.csv`

NOTE: The first line of this CSV will be RSpec formatting info and should be deleted.

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
