#!/bin/bash
set -e
set -x

pushd "magic-modules/build/$PROVIDER/$PRODUCT"
bundle install

if [ $PROVIDER = "chef" ]; then
    # TODO: https://github.com/GoogleCloudPlatform/magic-modules/issues/236
    # Re-enable chef tests by deleting this if block once the tests are fixed.
    echo "Skipping tests... See issue #236"
elif [ -z "$EXCLUDE_PATTERN" ]; then
  # The unit tests are all under e.g. spec/gcompute* - the integration tests
  # are better run using `bundle exec rake spec` instead of
  # `bundle exec parallel_rspec`, but if you just let it run the default
  # set of specs, parallel_rspec will try to run the integration tests
  # in addition.  This is just running all the tests that aren't excluded -
  # or, in the event of an empty exclude list, all the tests with the right
  # names.
  if ls spec/g$PRODUCT* > /dev/null 2&>1; then
    DISABLE_COVERAGE=true bundle exec parallel_rspec spec/g$PRODUCT*
  fi
else
    # parallel_rspec doesn't support --exclude_pattern
    IFS="," read -ra excluded <<< "$EXCLUDE_PATTERN"
    filtered=$(find spec -name "g$PRODUCT*_spec.rb" $(printf "! -wholename %s " ${excluded[@]}))
    DISABLE_COVERAGE=true bundle exec parallel_rspec ${filtered[@]}
fi

popd
