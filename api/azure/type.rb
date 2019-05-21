require 'api/type'

module Api
  module Azure
    module Type
      class ResourceGroupName < Api::Type::String
        def validate
          @order ||= 550
          super
        end
      end

      class Location < Api::Type::String
        def validate
          @order ||= 600
          super
        end
      end

      class Tags < Api::Type::KeyValuePairs
        def validate
          @order ||= 2000
          super
        end
      end

      class ResourceReference < Api::Type::String
        attr_reader :resource_type_name

        def validate
          super
          check_property :resource_type_name, ::String
        end
      end
    end
  end
end
