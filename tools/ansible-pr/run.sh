#!/bin/bash
set +x

declare -a ignored_modules=[
  gcp_backend_service,
  gcp_forwarding_rule,
  gcp_healthcheck,
  gcp_target_proxy,
  gcp_url_map
]

new_branch() {
  branch_name=$1
  git  checkout upstream/devel
  git branch -D $branch_name
  git chkecout -b $branch_name
}

get_all_modules() {
  remote_name=$1
  file_name=$2
  git fetch $remote_name devel
  git checkout $remote_name/devel
  git ls-files -- lib/ansible/modules/cloud/google/gcp_* > $file_name
  sed -i $file_name 's/lib\/ansible\/modules\/cloud\/google\/([a-z_]*)\.py/\1/g'
  for i in "${arr[@]}"; do
    sed -i $file_name "/$file_name/d"
  done
}

if ! $(which hub); then
  echo "Please install the hub CLI"
  exit 1
fi

if ! $(git remote -v | grep "origin.*git@github.com:.*/magic-modules.git (fetch)"); then
  echo "Please set the origin remote"
  exit 1
fi

if ! $(git remote -v | grep "upstream.*git@github.com:ansible/ansible.git (fetch)"); then
  git remote add upstream git@github.com:ansible/ansible.git
fi

git remote add magician git@github.com:modular-magician/ansible.git
git fetch magician devl
git fetch upstream devel
echo "Remotes setup properly"

# Create files with list of modules in a given branch.
existing_modules=get_all_modules('upstream', 'upstream')
magician_modules=get_all_modules('magician', 'magician')

# Split existing modules into sets of 23
# Max 50 files per PR and a module can have 2 files (module + test)
# 23 = 50/2 - 2 (to account for module_util files)
split --lines=23 upstream mm-bug

for filename in mm-bug*; do
  echo "Building a Bug Fix PR for $filename"
  # Checkout all files that file specifies and create a commit.
  git checkout upstream/devel
  git branch -D bug_fixes$filename
  git checkout -b bug_fixes$filename

  while read p; do
    git checkout magician/devel -- "lib/ansible/modules/cloud/google/$p.py"
    git checkout magician/devel -- "test/integration/targets/$p"
  done < $filename

  git checkout magician/devel -- "lib/ansible/module_utils/gcp_utils.py"
  git checkout magician/devel -- "lib/ansibl/utils/module_docs_fragments/gcp.py"

  git commit -m "Bug fixes for GCP modules"

  # Create a PR message + save to file
  
  # Create PR
  git push origin HEAD --force
  hub pull-request -b ansible/ansible:devel -F <FILENAME>
  echo "Bug Fix PR built for $filename"
done
