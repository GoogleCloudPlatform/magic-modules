module Overrides
  module Azure

    module ResourceOverride
      def azure_apply(_resource)
        convert_properties_to_datasource(_resource.all_user_properties, _resource.azure_sdk_definition) if @__is_data_source
        update_overriden_azure_sdk_definition(_resource)
        update_properties_default_sort_order(_resource)
      end

      private

      def convert_properties_to_datasource(properties, azure_sdk_definition)
        properties.each do |p|
          if p.is_a? Api::Azure::Type::ResourceGroupName
            p.instance_variable_set('@custom_schema_definition', 'templates/azure/terraform/schemas/datasource_resource_group_name.erb')
          elsif p.is_a? Api::Azure::Type::Location
            p.instance_variable_set('@custom_schema_definition', 'templates/azure/terraform/schemas/datasource_location.erb')
          elsif p.is_a? Api::Azure::Type::Tags
            p.instance_variable_set('@custom_schema_definition', 'templates/azure/terraform/schemas/datasource_tags.erb')
          end
          p.instance_variable_set('@input', false)
          unless p.azure_sdk_references.any?{|r| azure_sdk_definition.read.request.has_key?(r)}
            p.instance_variable_set('@required', false)
            p.instance_variable_set('@output', true)
            convert_properties_to_datasource(p.properties, azure_sdk_definition) if p.respond_to?(:properties) && !p.properties.nil?
          end
        end
      end

      def update_overriden_azure_sdk_definition(_resource)
        unless _resource.azure_sdk_definition.nil?
          _resource.azure_sdk_definition.filter_language! @azure_sdk_language
          override = instance_variable_get('@azure_sdk_definition')
          _resource.azure_sdk_definition.merge_overrides!(override) unless override.nil?
        end
      end
  
      def update_properties_default_sort_order(_resource)
        name = _resource.all_user_properties.find{|p| p.name == 'name'}
        name.instance_variable_set('@order', @name_default_order) unless name.nil?
        id = _resource.all_user_properties.find{|p| p.name == 'id'}
        id.instance_variable_set('@order', @id_default_order) unless id.nil?
      end
    end

  end
end
