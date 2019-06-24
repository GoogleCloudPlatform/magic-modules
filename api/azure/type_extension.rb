module Api
  module Azure
    module Type

      module Fields
        attr_reader :order
        attr_reader :sample_value
        attr_reader :azure_sdk_references
      end

      module TypeExtension
        def azure_validate
          default_order = 10
          default_order = 1 if @name == "name"
          default_order = 0 if @name == "id"
          check :order, type: ::Integer, default: default_order
          check :azure_sdk_references, type: ::Array, item_type: ::String, required: true
        end
      end

    end
  end
end
