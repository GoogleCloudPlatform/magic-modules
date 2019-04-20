require 'provider/azure/example/example'

module Provider
  module Azure
    module Core
      def get_example_by_reference(reference)
        get_example_by_names reference.example, reference.product
      end

      def get_example_by_names(example_name, product_name = nil)
        product_name ||= @api.prefix
        example_yaml = "products/#{product_name}/examples/#{@provider}/#{example_name}.yaml"
        example = Google::YamlValidator.parse(File.read(example_yaml))
        raise "#{example_yaml}(#{example.class}) is not Provider::Azure::Example" unless example.is_a?(Provider::Azure::Example)
        example.validate
        example
      end
    end
  end
end
