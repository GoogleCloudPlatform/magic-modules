#!/bin/bash
set -e
set -x

pushd "magic-modules/build/$PROVIDER/$PRODUCT"
bundle install

if [ $PROVIDER = "chef" ]; then
    # TODO: https://github.com/GoogleCloudPlatform/magic-modules/issues/236
    # Re-enable chef tests by deleting this if block once the tests are fixed.
    echo "Skipping tests... See issue #236"
else if [ -z "$EXCLUDE_PATTERN" ]; then
    DISABLE_COVERAGE=true bundle exec parallel_rspec spec/
else
    # parallel_rspec doesn't support --exclude_pattern
    IFS="," read -ra excluded <<< "$EXCLUDE_PATTERN"
    filtered=$(find spec -name '*_spec.rb' $(printf "! -wholename %s " ${excluded[@]}))
    DISABLE_COVERAGE=true bundle exec parallel_rspec ${filtered[@]}
fi

popd
