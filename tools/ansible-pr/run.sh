#!/bin/bash
set -x

declare -a ignored_modules=(
  gcp_backend_service
  gcp_forwarding_rule
  gcp_healthcheck
  gcp_target_proxy
  gcp_url_map
)

get_all_modules() {
  remote_name=$1
  file_name=$remote_name
  git fetch $remote_name
  git checkout $remote_name/devel
  git ls-files -- lib/ansible/modules/cloud/google/gcp_* | cut -d/ -f 6 | cut -d. -f 1 > $file_name
  
  for i in "${ignored_modules[@]}"; do
    sed -i "/$i/d" $file_name
  done
}

if ! which hub; then
  echo "Please install the hub CLI"
  exit 1
fi

if ! git remote -v | grep "origin"; then
  echo "Please set the origin remote"
  exit 1
fi

set -e

git remote add upstream git@github.com:ansible/ansible.git

# Use HTTPS endpoint so we don't have to setup SSH keys.
git remote add magician git@github.com:modular-magician/ansible.git
git fetch magician devel
git fetch upstream devel
echo "Remotes setup properly"

# Create files with list of modules in a given branch.
get_all_modules "upstream"
get_all_modules "magician"

# Split existing modules into sets of 23
# Max 50 files per PR and a module can have 2 files (module + test)
# 23 = 50/2 - 2 (to account for module_util files)
split -l 23 upstream mm-bug

for filename in mm-bug*; do
  echo "Building a Bug Fix PR for $filename"
  # Checkout all files that file specifies and create a commit.
  git checkout upstream/devel
  git checkout -b bug_fixes$filename


  while read p; do
    git checkout magician/devel -- "lib/ansible/modules/cloud/google/$p.py"
    if [[ $p != *"facts"* ]]; then
      git checkout magician/devel -- "test/integration/targets/$p"
    fi
  done < $filename

  git checkout magician/devel -- "lib/ansible/module_utils/gcp_utils.py"

  git commit -m "Bug fixes for GCP modules"

  # Create a PR message + save to file
  ruby ../../tools/ansible-pr/generate_template.rb > bug_fixes$filename

  # Create PR
  git push origin HEAD --force
  hub pull-request -b ansible/ansible:devel -F bug_fixes$filename
  echo "Bug Fix PR built for $filename"
done

## Get list of new modules (in magician, not in upstream)
comm -3 <(sort magician) <(sort upstream) > new_modules

while read module; do
  echo "Building a New Module PR for $module"
  git checkout upstream/devel
  git checkout -b $module

  git checkout magician/devel -- "lib/ansible/modules/cloud/google/$module.py"
  if [[ $module != *"facts"* ]]; then
    git checkout magician/devel -- "test/integration/targets/$module"
  fi

  git checkout magician/devel -- "lib/ansible/module_utils/gcp_utils.py"

  # Create a PR message + save to file
  git commit -m "New Module: $module"
  ruby ../../tools/ansible-pr/generate_template.rb --new-module-name $module > bug_fixes$filename

  # Create PR
  git push origin HEAD --force
  hub pull-request -b ansible/ansible:devel -F bug_fixes$filename
  echo "New Module PR built for $module"
done < new_modules
