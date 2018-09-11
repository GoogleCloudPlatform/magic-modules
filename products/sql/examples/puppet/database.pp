<% if false # the license inside this if block pertains to this file -%>
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
<% instance_name = 'sql-test-${sql_instance_suffix}' -%>
<% if name != 'README.md' -%>
<%= compile 'templates/license.erb' -%>

<%= lines(autogen_notice :puppet) -%>

<%= compile 'templates/puppet/examples~credential.pp.erb' -%>

# TODO(alexstephen): Change this warning and remove the "requires to to exists"
# once a resource reference is added (it will enforce that automatically).
#
# This example requires an instance to exist. You should set
# FACTER_sql_instance_suffix, or use any other Puppet # supported way, to set a
# global variable $sql_instance_suffix.
#
# For example you can define the fact to be an always increasing value:
#
# $ FACTER_sql_instance_suffix=100 puppet apply examples/database.pp
#
# If that instance does not exist in your project run the examples/instance.pp
# to create it, with the same $sql_instance_suffix.
if !defined('$sql_instance_suffix') {
  fail('For this example to run you need to define a fact named
       "sql_instance_suffix". Please refer to the documentation inside
       the example file "<%= name -%>"')
}

gsql_instance { <%= example_resource_name(instance_name) -%>:
  ensure     => present,
  project    => $project, # e.g. 'my-test-project'
  credential => 'mycred',
}

<% else -%>
# Tip: Remember to define gsql_instance to match the 'instance' property.
<% end -%>
gsql_database { <%= example_resource_name('webstore') -%>:
  ensure     => present,
  charset    => 'utf8',
  instance   => <%= example_resource_name(instance_name) -%>,
  project    => $project, # e.g. 'my-test-project'
  credential => 'mycred',
}
