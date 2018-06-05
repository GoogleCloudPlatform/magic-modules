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
require 'provider/ansible/manifest'
require 'provider/ansible/example'
require 'provider/ansible/documentation'
require 'provider/ansible/module'
require 'provider/ansible/request'
require 'provider/ansible/resourceref'
require 'provider/ansible/resource_override'
require 'provider/ansible/property_override'
require 'provider/ansible/selflink'

module Provider
  module Ansible
    module Facts
      # Handles Configuration for Ansible Facts
      class Config < Provider::Config
        attr_reader :manifest

        def provider
          Provider::Ansible::Facts::Core
        end

        def resource_override
          Provider::Ansible::ResourceOverride
        end

        def property_override
          Provider::Ansible::PropertyOverride
        end

        def validate
          super
          check_optional_property :manifest, Provider::Ansible::Manifest
        end
      end

      # Provider code for Ansible Facts.
      # Has full access to all functions from regular Ansible provider
      class Core < Provider::Ansible::Core
        def list_kind(object)
          "#{object.kind}List"
        end

        private

        def generate_resource(data)
          target_folder = data[:output_folder]
          FileUtils.mkpath target_folder
          name = "#{module_name(data[:object])}_facts"
          generate_resource_file data.clone.merge(
            default_template: 'templates/ansible/facts.erb',
            out_file: File.join(target_folder,
                                "lib/ansible/modules/cloud/google/#{name}.py")
          )
        end

        def generate_resource_tests(data) end

        def generate_network_datas(data, object) end

        def generate_base_property(data) end

        def generate_simple_property(type, data) end

        def generate_typed_array(type, data) end

        def emit_nested_object(data) end

        def emit_resourceref_object(data) end
      end
    end
    # rubocop:enable Metrics/ClassLength
  end
end
