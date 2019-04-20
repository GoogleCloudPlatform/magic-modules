require 'provider/core'
require 'api/object'

module Provider
  module Azure
    class ExampleReference < Api::Object
      attr_reader :product
      attr_reader :example

      def validate
        super
        check_optional_property :product, String
        check_property :example, String
      end
    end

    class Example < Api::Object
      attr_reader :resource
      attr_reader :description
      attr_reader :prerequisites
      attr_reader :properties

      def validate
        super
        check_property :resource, String
        check_optional_property :description, String
        check_optional_property :prerequisites, Array
        check_optional_property_list :prerequisites, ExampleReference
        check_property :properties, Hash
      end
    end
  end
end
