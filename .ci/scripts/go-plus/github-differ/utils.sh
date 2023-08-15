#!/bin/bash

update_package_name() {
  old_package_name="$1"
  new_package_name="$2"

  # Update import statements and references within files
  find . -type f -name "*.go" -exec sed -i.bak "s~$old_package_name~$new_package_name~g" {} +

  # Update go.mod file
  sed -i.bak "s|$old_package_name|$new_package_name|g" go.mod

  # Optional: Update go.sum file
  sed -i.bak "s|$old_package_name|$new_package_name|g" go.sum
}
