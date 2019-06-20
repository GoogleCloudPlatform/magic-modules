# api.yaml linter

This linter lints api.yaml against the Discovery Docs. It iterates through all
apis and Resources in api.yaml files and attempts to find the correspond
Discovery Doc to determine if Resource fields are different or missing.

Because api.yaml is handwritten, it will inevitably drift from the Discovery
Docs (the source of truth).

## Running the linter
The linter is powered by RSpec.

```
  rspec tools/linter/run.rb
```

Tests are divided into two parts: property tests and resource tests. The linter uses [rspec tags](https://relishapp.com/rspec/rspec-core/v/3-8/docs/command-line/tag-option) to filter which tests are run. The tests are broken down into `property` and `resource` tests and can be run with `--tag property` or `--tag resource`. It is also possible to filter by product or resource `--tag product:<name>` or `--tag resource:<name>`.

To run only property tests, do the following:
```
  rspec tools/linter/run.rb --tag property
```

To run only cloudbuild tests, do the following:
```
  rspec tools/linter/run.rb --tag product:cloudbuild
```

To run only cloudbuild Trigger tests, do the following:
```
  rspec tools/linter/run.rb --tag resource:Trigger
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
dns ManagedZone labels should exist
```
Each property is tested to see if it exists.
If a property does not exist, RSpec will print the following message.
The message is formatted as "<product> <resource> <property> should exist"
