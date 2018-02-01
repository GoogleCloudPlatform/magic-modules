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
<% unless name == "README.md" -%>
<%= compile 'templates/license.erb' -%>

<%= compile 'templates/autogen_notice.erb' -%>

<% end # name == README.md -%>
# To create a new project the authenticated user need the privilege to create
# new projects, either standalone or under an organization. This example uses
# the 'application-default' provider, which draws from 'gcloud' (Google Cloud
# SDK tool) the current user credentials.
#
# To make the example work you have to run this command once before applying the
# manifest and follow its instructions:
#
#     gcloud auth application-default login
#
# Alternatively you can setup a service account and use the *preferred*
# 'serviceaccount' provider instead with a JSON key file.
gauth_credential { 'mycred':
  provider => defaultuserapplication,
  scopes   => [
    'https://www.googleapis.com/auth/cloud-platform',
  ],
}

# Project ID needs to be unique. Add a random suffix so they are always
# unique. You should set FACTER_project_suffix, or use any other Puppet
# supported way, to set a global variable $project_suffix.
#
# For example you can define the fact to be an always increasing value:
#
# $ FACTER_project_suffix=$(date +%s) puppet apply examples/delete_project.pp
#
# To be able to delete the project via Puppet make sure the instance ID matches
# the ID used during creation. If you used the create example and specified the
# 'project_suffix', you should match it as well during deletion.
if !defined('$project_suffix') {
  fail('For this example to run you need to define a fact named
       "project_suffix". Please refer to the documentation inside
       the example file "examples/project.pp"')
}

gresourcemanager_project { 'My Sample Project':
  ensure     => absent,
  id         => "test-project-${project_suffix}",
  credential => 'mycred',
}
