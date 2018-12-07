require 'api/type'

module Api
  module Azure
    module Type

      class ResourceGroupName < Api::Type::String
      end

      class Location < Api::Type::String
      end

      class Tags < Api::Type::KeyValuePairs
      end

    end
  end
end
