require 'api/azure/type'

module Provider
  module Azure
    module Ansible
      module Helpers
        def is_resource_name?(property)
          property.parent.nil? && property.name == 'name'
        end

        def is_tags?(property)
          property.is_a? Api::Azure::Type::Tags
        end

        def is_tags_defined?(object)
          object.all_user_properties.any?{|p| is_tags?(p)}
        end

        def get_tags_property(object)
          object.all_user_properties.find{|p| is_tags?(p)}
        end

        def is_location?(property)
          property.parent.nil? && property.is_a?(Api::Azure::Type::Location)
        end

        def is_location_defined?(object)
          object.all_user_properties.any?{|p| is_location?(p)}
        end

        def is_resource_group?(property)
          property.parent.nil? && property.is_a?(Api::Azure::Type::ResourceGroupName)
        end

        def always_has_value?(property)
          property.required || !property.default_value.nil?
        end
      end
    end
  end
end
