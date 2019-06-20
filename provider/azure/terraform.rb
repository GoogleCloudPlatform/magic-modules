require 'provider/azure/terraform/custom_code'
require 'provider/azure/terraform/helpers'
require 'provider/azure/terraform/schema'
require 'provider/azure/terraform/sub_template'
require 'provider/azure/terraform/sdk/sdk_type_definition_descriptor'
require 'provider/azure/terraform/sdk/sdk_marshal_descriptor'
require 'provider/azure/terraform/sdk/expand_flatten_descriptor'
require 'provider/azure/terraform/sdk/sub_template'
require 'provider/azure/terraform/sdk/helpers'
require 'provider/azure/terraform/example/sub_template'
require 'provider/azure/terraform/example/helpers'
require 'provider/azure/terraform/acctest/sub_template'

require 'provider/azure/terraform/resource_override'
require 'provider/azure/terraform/property_override'

module Provider
  module Azure
    module Terraform
      include Provider::Azure::Terraform::Helpers
      include Provider::Azure::Terraform::Schema
      include Provider::Azure::Terraform::SubTemplate
      include Provider::Azure::Terraform::SDK::SubTemplate
      include Provider::Azure::Terraform::SDK::Helpers
      include Provider::Azure::Terraform::Example::SubTemplate
      include Provider::Azure::Terraform::Example::Helpers
      include Provider::Azure::Terraform::AccTest::SubTemplate

      def initialize
        @provider = 'terraform'
      end

      def azure_tf_types(map)
        map[Api::Azure::Type::ResourceReference] = 'schema.TypeString'
        map
      end

      def azure_generate_resource(data)
        dir = "azurerm"
        filepath = File.join(target_folder, "resource_arm_#{name}.go")
        # TODO: Implement this
      end

      def azure_generate_documentation(data)
        filepath = File.join(target_folder, "#{name}.html.markdown")
        # TODO: Implement this
      end

      def azure_generate_resource_tests(data)
        dir = "azurerm"
        filepath = File.join(target_folder, "resource_arm_#{name}_test.go")
        # TODO: Implement this
      end

      def compile_datasource(data)
        dir = 'azurerm'
        target_folder = File.join(data[:output_folder], dir)
        FileUtils.mkpath target_folder
        name = data[:object].name.underscore
        product_name = data[:product_name].underscore

        filepath = File.join(target_folder, "data_source_#{name}.go")
        generate_resource_file data.clone.merge(
          default_template: 'templates/terraform/datasource.erb',
          out_file: filepath
        )

        filepath = File.join(target_folder, "data_source_#{name}_test.go")
        generate_resource_file data.clone.merge(
          default_template: 'templates/terraform/examples/base_configs/datasource_test.go.erb',
          out_file: filepath
        )

        target_folder = File.join(data[:output_folder], 'website', 'docs', 'd')
        FileUtils.mkpath target_folder
        filepath = File.join(target_folder, "#{name}.html.markdown")
        generate_resource_file data.clone.merge(
          default_template: 'templates/terraform/datasource.html.markdown.erb',
          out_file: filepath
        )
      end

    end
  end
end
