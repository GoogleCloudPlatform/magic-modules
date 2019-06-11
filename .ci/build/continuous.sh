#!/bin/bash

# This file is for builds that are running on kokoro. It is primarily to be a kokoro
# specific wrapper around existing magic modules commands

# Fail on any error.
set -e
# Display commands being run.
set -x

# bundle exec compiler -p products/monitoring -e terraform -o /Users/chrisst/work/go/src/github.com/terraform-providers/terraform-provider-google
ruby --version

# bundle -v
echo "done!"

