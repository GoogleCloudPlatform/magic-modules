require 'api/azure/type'

module Provider
  module Azure
    module Ansible
      module Helpers
        def is_tags?(property)
          property.is_a? Api::Azure::Type::Tags
        end

        def is_tags_defined?(object)
          object.all_user_properties.any?{|p| is_tags?(p)}
        end

        def is_location?(property)
          property.is_a? Api::Azure::Type::Location
        end

        def is_location_defined?(object)
          object.all_user_properties.any?{|p| is_location?(p)}
        end
      end
    end
  end
end
