require 'provider/azure/example/example'

module Provider
  module Azure
    module Core
      def get_example_by_reference(reference)
        get_example_by_names reference.example, reference.product
      end

      def get_example_by_names(example_name, product_name = nil)
        spec_dir = File.dirname(@config.cfg_file)
        product_name ||= File.basename(spec_dir)
        example_yaml = File.join(File.dirname(spec_dir), product_name, 'examples', @provider, "#{example_name}.yaml")
        example = Google::YamlValidator.parse(File.read(example_yaml))
        raise "#{example_yaml}(#{example.class}) is not Provider::Azure::Example" unless example.is_a?(Provider::Azure::Example)
        example.validate
        example
      end

      def get_custom_template_path(template_path)
        return nil if template_path.nil?
        spec_dir = File.dirname(@config.cfg_file)
        File.join(spec_dir, template_path)
      end
    end
  end
end
