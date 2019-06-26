require 'api/type'

module Api
  module Azure
    module Type

      class ResourceGroupName < Api::Type::String
        def validate
          @order ||= 3
          super
        end
      end

      class Location < Api::Type::String
        def validate
          @order ||= 5
          super
        end
      end

      class Tags < Api::Type::KeyValuePairs
        def validate
          @order ||= 20
          super
        end
      end

      class ResourceReference < Api::Type::String
        attr_reader :resource_type_name

        def validate
          super
          check :resource_type_name, type: ::String, required: true
        end
      end

      class ISO8601Duration < Api::Type::String
      end

      class ISO8601DateTime < Api::Type::String
      end

    end
  end
end
