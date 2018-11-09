# api.yaml linter

The idea of this linter is to lint api.yaml against the Discovery docs.

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
