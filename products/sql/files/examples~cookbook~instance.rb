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

<%= compile 'templates/autogen_notice.erb' -%>

<%= compile 'templates/chef/example~auth.rb.erb' -%>

# Cloud SQL cannot reuse instance names. Add a random suffix so they are always
# unique.
#
# To be able to delete the instance via Chef make sure the instance ID matches
# the ID used during creation. If you used the create example and specified the
# 'sql_instance_suffix', you should match it as well during deletion.
raise ['For this example to run you need to define a env. variable named',
       '"sql_instance_suffix". Please refer to the documentation inside',
       'the example file "<%= name -%>"'].join(' ') \
  unless ENV.key?('sql_instance_suffix')

<% end -%>
<% res_name = 'sql-test-#{ENV[\'sql_instance_suffix\']}' -%>
gsql_instance <%= example_resource_name(res_name) -%> do
  action :create
  database_version 'MYSQL_5_7'
  settings({
    tier: 'db-n1-standard-1',
    ip_configuration:  {
      authorized_networks: [
        # The ACL below is for example only. (do NOT use in production as-is)
        {
          name: 'google dns server',
          value: '8.8.8.8/32'
        }
      ]
    }
  })
  region 'us-central1'
  project 'google.com:graphite-playground'
  credential 'mycred'
end
