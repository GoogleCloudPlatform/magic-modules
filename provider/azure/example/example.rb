require 'provider/core'
require 'api/object'

module Provider
  module Azure
    class ExampleReference < Api::Object
      attr_reader :product
      attr_reader :example

      def validate
        super
        check :product, type: ::String
        check :example, type: ::String, required: true
      end
    end

    class Example < Api::Object
      attr_reader :resource
      attr_reader :description
      attr_reader :prerequisites
      attr_reader :properties

      def validate
        super
        check :resource, type: ::String, required: true
        check :description, type: ::String
        check :prerequisites, type: ::Array, item_type: ExampleReference
        check :properties, type: ::Hash, required: true
      end
    end
  end
end
