#!/usr/bin/env bash

set -x
# Constants + functions
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
  ssh-agent bash -c "ssh-add ~/github_private_key; git fetch $remote_name"
  git checkout $remote_name/devel
  git ls-files -- lib/ansible/modules/cloud/google/gcp_* | cut -d/ -f 6 | cut -d. -f 1 > $file_name
  
  for i in "${ignored_modules[@]}"; do
    sed -i "/$i/d" $file_name
  done
}

# Install dependencies for Template Generator
pushd "magic-modules-gcp"
bundle install

# Setup SSH keys.

# Since these creds are going to be managed externally, we need to pass
# them into the container as an environment variable.  We'll use
# ssh-agent to ensure that these are the credentials used to update.
set +x
echo "$CREDS" > ~/github_private_key
set -x
chmod 400 ~/github_private_key
popd

# Clone ansible/ansible
ssh-agent bash -c "ssh-add ~/github_private_key; git clone git@github.com:modular-magician/ansible.git"

# Setup Git config and remotes.
pushd "ansible"
git config --global user.email "magic-modules@google.com"
git config --global user.name "Modular Magician"

git remote remove origin
git remote add origin git@github.com:modular-magician/ansible.git
git remote add upstream git@github.com:ansible/ansible.git
git remote add magician git@github.com:modular-magician/ansible.git
echo "Remotes setup properly"
popd

# Copy code into ansible/ansible + commit to our fork
# By using the "ansible_devel" provider, we get versions of the resources that work
# with ansible devel.
pushd "magic-modules-gcp"
ruby compiler.rb -a -e ansible -f ansible_devel -o ../ansible/
popd

# Commit code from magic modules into our fork
pushd "ansible"
git add lib/ansible/modules/cloud/google/gcp_* test/integration/targets/gcp_*
git commit -m "Migrating code from collection"
ssh-agent bash -c "ssh-add ~/github_private_key; git push magician devel"

set -e

ssh-agent bash -c "ssh-add ~/github_private_key; git fetch magician devel"
ssh-agent bash -c "ssh-add ~/github_private_key; git fetch upstream devel"

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
    if [[ $p != *"info"* ]]; then
      git checkout magician/devel -- "test/integration/targets/$p"
    fi
  done < $filename

  git checkout magician/devel -- "lib/ansible/module_utils/gcp_utils.py"
  git checkout magician/devel -- "lib/ansible/plugins/doc_fragments/gcp.py"

  # This commit may be empty
  set +e
  git commit -m "Bug fixes for GCP modules"

  # Create a PR message + save to file
  ruby ../magic-modules-gcp/tools/ansible-pr/generate_template.rb > bug_fixes$filename

  # Create PR
  ssh-agent bash -c "ssh-add ~/github_private_key; git push origin bug_fixes$filename --force"
  hub pull-request -b ansible/ansible:devel -F bug_fixes$filename -f
  set -e

  echo "Bug Fix PR built for $filename"
done

## Get list of new modules (in magician, not in upstream)
comm -3 <(sort magician) <(sort upstream) > new_modules

while read module; do
  echo "Building a New Module PR for $module"
  git checkout upstream/devel
  git checkout -b $module

  git checkout magician/devel -- "lib/ansible/modules/cloud/google/$module.py"
  if [[ $module != *"info"* ]]; then
    git checkout magician/devel -- "test/integration/targets/$module"
  fi

  git checkout magician/devel -- "lib/ansible/module_utils/gcp_utils.py"

  # Create a PR message + save to file
  set +e
  git commit -m "New Module: $module"
  ruby ../magic-modules-gcp/tools/ansible-pr/generate_template.rb --new-module-name $module > $module

  # Create PR
  ssh-agent bash -c "ssh-add ~/github_private_key; git push origin $module --force"
  hub pull-request -b ansible/ansible:devel -F $module -f
  set -e

  echo "New Module PR built for $module"
done < new_modules
