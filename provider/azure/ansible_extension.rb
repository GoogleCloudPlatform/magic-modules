require 'provider/azure/ansible/config'
require 'provider/azure/ansible/helpers'
require 'provider/azure/ansible/sub_template'
require 'provider/azure/ansible/sdk/sdk_marshal_descriptor'
require 'provider/azure/ansible/sdk/helpers'
require 'provider/azure/ansible/module/sub_template'
require 'provider/azure/ansible/sdk/sub_template'
require 'provider/azure/ansible/example/helpers'
require 'provider/azure/ansible/example/sub_template'

require 'provider/azure/ansible/resource_override'
require 'provider/azure/ansible/property_override'

module Provider
  module Azure
    module AnsibleExtension
      include Provider::Azure::Ansible::Helpers
      include Provider::Azure::Ansible::SDK::Helpers
      include Provider::Azure::Ansible::SubTemplate
      include Provider::Azure::Ansible::Module::SubTemplate
      include Provider::Azure::Ansible::SDK::SubTemplate
      include Provider::Azure::Ansible::Example::Helpers
      include Provider::Azure::Ansible::Example::SubTemplate

      def initialize(config, api, start_time)
        super
        @provider = 'ansible'
      end

      def azure_python_type(prop)
        return 'raw' if prop.is_a? Api::Azure::Type::ResourceReference
        return 'list' if prop.is_a? Api::Azure::Type::Tags
        python_type prop
      end

      def azure_module_name(object)
        "azure_rm_#{object.name.downcase}"
      end

      def azure_generate_resource(data)
        target_folder = data.output_folder
        path = File.join(target_folder, "lib/ansible/modules/cloud/azure/#{azure_module_name(data.object)}.py")
        data.generate(
          data.object.template || 'templates/azure/ansible/resource.erb',
          path,
          self
        )
      end

      def azure_generate_resource_tests(data)
        return unless data.object.has_tests
        return if data.object.inttests.empty?

        name = azure_module_name(data.object)
        target_folder = data.output_folder
        target_folder = File.join(target_folder, "test/integration/targets/#{name}")

        data.generate('templates/azure/ansible/integration_test.erb', File.join(target_folder, 'tasks/main.yml'), self)
        data.generate('templates/azure/ansible/test/meta.erb', File.join(target_folder, 'meta/main.yml'), self)
        data.generate('templates/azure/ansible/test/aliases.erb', File.join(target_folder, 'aliases'), self)
      end

      def azure_compile_datasource(data)
        target_folder = data.output_folder
        name = "#{azure_module_name(data.object)}_info"
        path = File.join(target_folder, "lib/ansible/modules/cloud/azure/#{name}.py")
        data.generate('templates/azure/ansible/info.erb', path, self)
      end

    end
  end
end
