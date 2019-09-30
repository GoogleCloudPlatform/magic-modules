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

require 'provider/config'
require 'provider/core'
require 'provider/ansible/config'
require 'provider/ansible/documentation'
require 'provider/ansible/example'
require 'provider/ansible/module'
require 'provider/ansible/request'
require 'provider/ansible/resourceref'
require 'provider/ansible/version_added'
require 'provider/ansible/facts_override'
require 'overrides/ansible/resource_override'
require 'overrides/ansible/property_override'

module Provider
  module Ansible
    # Ansible code needs to go into devel as well as
    # collections. This provider handles the creation of
    # devel code (identical to regular code, except in a couple places)
    class Devel < Provider::Ansible::Core
      def module_utils_import_path
        'ansible_collections.google.cloud.plugins.module_utils.gcp_utils'
      end

      def generate_resource(data)
        target_folder = data.output_folder
        name = module_name(data.object)
        path = File.join(target_folder,
                         "lib/ansible/modules/cloud/google/#{name}.py")
        data.generate(
          data.object.template || 'templates/ansible/resource.erb',
          path,
          self
        )
      end

      def compile_datasource(data)
        target_folder = data.output_folder
        name = module_name(data.object)
        data.generate('templates/ansible/facts.erb',
                      File.join(target_folder,
                                "lib/ansible/modules/cloud/google/#{name}_info.py"),
                      self)

        # Generate symlink for old `facts` modules.
        return if version_added(data.object, :facts) >= '2.9'

        deprecated_facts_path = File.join(target_folder,
                                          "lib/ansible/modules/cloud/google/_#{name}_facts.py")
        return if File.exist?(deprecated_facts_path)

        File.symlink "#{name}_info.py", deprecated_facts_path
      end
    end
  end
end
