<% if false # the license inside this if block assertains to this file -%>
# Copyright 2017 Google Inc.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
<% end -%>
<% unless name == 'README.md' -%>
<%= compile 'templates/license.erb' -%>

<%= lines(autogen_notice :puppet) -%>

<%= compile 'templates/puppet/examples~credential.pp.erb' -%>

# Cloud SQL cannot reuse instance names. Add a random suffix so they are always
# unique. You should set FACTER_sql_instance_suffix, or use any other Puppet
# supported way, to set a global variable $sql_instance_suffix.
#
# For example you can define the fact to be an always increasing value:
#
# $ FACTER_sql_instance_suffix=$(date +%s) puppet apply examples/instance.pp
#
# To be able to delete the instance via Puppet make sure the instance ID matches
# the ID used during creation. If you used the create example and specified the
# 'sql_instance_suffix', you should match it as well during deletion.
if !defined('$sql_instance_suffix') {
  fail('For this example to run you need to define a fact named
       "sql_instance_suffix". Please refer to the documentation inside
       the example file "<%= name -%>"')
}

<% end -%>
<% instance_name = 'sql-test-${sql_instance_suffix}' -%>
gsql_instance { <%= example_resource_name(instance_name) -%>:
  ensure     => absent,
  project    => 'google.com:graphite-playground',
  credential => 'mycred',
}
