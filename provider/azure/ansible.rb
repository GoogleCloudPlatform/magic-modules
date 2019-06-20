require 'provider/azure/ansible/helpers'
require 'provider/azure/ansible/sub_template'
require 'provider/azure/ansible/sdk/sdk_marshal_descriptor'
require 'provider/azure/ansible/sdk/property_normalize_descriptor'
require 'provider/azure/ansible/sdk/helpers'
require 'provider/azure/ansible/module/sub_template'
require 'provider/azure/ansible/sdk/sub_template'
require 'provider/azure/ansible/example/helpers'
require 'provider/azure/ansible/example/sub_template'

require 'provider/azure/ansible/resource_override'
require 'provider/azure/ansible/property_override'

module Provider
  module Azure
    module Ansible
      include Provider::Azure::Ansible::Helpers
      include Provider::Azure::Ansible::SDK::Helpers
      include Provider::Azure::Ansible::SubTemplate
      include Provider::Azure::Ansible::Module::SubTemplate
      include Provider::Azure::Ansible::SDK::SubTemplate
      include Provider::Azure::Ansible::Example::Helpers
      include Provider::Azure::Ansible::Example::SubTemplate

      def initialize
        @provider = 'ansible'
      end

      def azure_python_type(prop)
        return 'raw' if prop.is_a? Api::Azure::Type::ResourceReference
        return 'list' if prop.is_a? Api::Azure::Type::Tags
        nil
      end

      def azure_module_name(object)
        "azure_rm_#{object.name.downcase}"
      end

      def azure_generate_resource(data)
        path = File.join(target_folder, "lib/ansible/modules/cloud/azure/#{azure_module_name(data.object)}.py")
        # TODO: Implement this
      end

      def azure_generate_resource_tests(data)
        prod_name = data[:object].name.underscore
        path = ["products/#{data[:product_name]}",
                "examples/ansible/#{prod_name}.yaml"].join('/')

        return unless data[:object].has_tests
        return if data[:object].inttests.empty?

        target_folder = data[:output_folder]
        FileUtils.mkpath target_folder

        name = module_name(data[:object])
        target_folder = File.join(target_folder, "test/integration/targets/#{name}")

        generate_resource_file data.clone.merge(
          default_template: 'templates/ansible/integration_test.erb',
          out_file: File.join(target_folder, 'tasks/main.yml')
        )
        generate_resource_file data.clone.merge(
          default_template: 'templates/azure/ansible/test/meta.erb',
          out_file: File.join(target_folder, 'meta/main.yml')
        )
        generate_resource_file data.clone.merge(
          default_template: 'templates/azure/ansible/test/aliases.erb',
          out_file: File.join(target_folder, 'aliases')
        )
      end

      def azure_compile_datasource(data)
        name = "#{module_name(data.object)}_info"
        path = File.join(target_folder, "lib/ansible/modules/cloud/azure/#{name}.py")
        # TODO: Implement this
      end

    end
  end
end
