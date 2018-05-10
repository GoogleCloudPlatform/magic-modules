#!/bin/bash
set -e
set -x

# Setup symlinks for puppet modules
declare -r module_path="/opt/puppetlabs/puppet/modules"
declare -r root_dir="/opt/magic-modules"

echo "Root dir: ${root_dir}"

for mod in ${root_dir}/build/puppet/*; do
	mod_name=$(basename "${mod}")
	ln -s "$(realpath "${mod}")" "${module_path}/g${mod_name}"
done

exec "$@"