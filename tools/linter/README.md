# api.yaml linter

The idea of this linter is to lint api.yaml against the Discovery docs.

Because api.yaml is handwritten, it will inevitably veer away from the Discovery docs
(the source of truth) or just be inputted incorrectly.

## How to run
All tests live in tools/linter/tests.rb
Run `ruby tests/linter/tests.rb` to run.
