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
        'ansible.module_utils.gcp_utils'
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

      def compile_datasource(pwd, data)
        target_folder = data.output_folder
        name = module_name(data.object)
        data.generate(pwd,
                      'templates/ansible/facts.erb',
                      File.join(target_folder,
                                "lib/ansible/modules/cloud/google/#{name}_info.py"),
                      self)

        # Generate symlink for old `facts` modules.
        return if version_added(data.object, :facts).to_f >= 2.9

        deprecated_facts_path = File.join(target_folder,
                                          "lib/ansible/modules/cloud/google/_#{name}_facts.py")
        return if File.exist?(deprecated_facts_path)

        File.symlink "#{name}_info.py", deprecated_facts_path
      end

      def generate_resource_tests(data)
        prod_name = data.object.name.underscore
        path = ["products/#{data.product.api_name}",
                "examples/ansible/#{prod_name}.yaml"].join('/')

        return if data.object.tests.tests.empty?

        target_folder = data.output_folder

        name = module_name(data.object)

        # Generate the main file with a list of tests.
        path = File.join(target_folder,
                         "test/integration/targets/#{name}/tasks/main.yml")
        data.generate(
          'templates/ansible/tests_main.erb',
          path,
          self
        )

        # Generate each of the tests individually
        data.object.tests.tests.each do |t|
          path = File.join(target_folder,
                           "test/integration/targets/#{name}/tasks/#{t.name}.yml")
          data.generate(
            t.path,
            path,
            self
          )
        end
      end

      def generate_resource_sweepers(data) end

      def compile_common_files(_arg1, _arg2, _arg3) end

      def copy_common_files(output_folder, provider_name = nil)
        # version_name is actually used because all of the variables in scope in this method
        # are made available within the templates by the compile call.
        # TODO: remove version_name, use @target_version_name or pass it in expicitly
        # rubocop:disable Lint/UselessAssignment
        version_name = @target_version_name
        # rubocop:enable Lint/UselessAssignment
        provider_name ||= self.class.name.split('::').last.downcase
        return unless File.exist?("provider/#{provider_name}/common~copy~devel.yaml")

        Google::LOGGER.info "Copying common files for #{provider_name}"
        files = YAML.safe_load(compile("provider/#{provider_name}/common~copy~devel.yaml"))
        copy_file_list(output_folder, files)
      end

      def generate_resource_files(data)
        return unless @config&.files&.resource

        files = @config.files.resource
                       .map { |k, v| [k % module_name(data.object), v] }
                       .to_h

        # Test directory lives in a different place in devel.
        files.transform_keys! { |k| k.gsub('tests/', 'test/') }

        file_template = ProductFileTemplate.new(
          data.output_folder,
          data.name,
          @api,
          data.version,
          build_env
        )
        compile_file_list(data.output_folder, files, file_template)
      end
    end
  end
end
