module Azure
  module YamlValidatorExtension

    # Does extended validation checking for a variable
    # options:
    # :default      - the default value for this variable if its nil
    # :key_type     - the allowed types that all keys in a hash should be
    # :item_type    - the allowed types that all values in an array or hash should be
    # :allowed      - the allowed values that this non-array variable should be
    # :item_allowed - the allowed values that all values in an array or hash should be
    # :required     - is the variable required? (defaults: false)
    def check_ext(variable, **opts)
      check variable, opts

      value = instance_variable_get("@#{variable}")

      # Check key_type and item_type
      if value.is_a?(::Hash)
        raise "#{variable} must have key_type and item_type on hashes" unless opts[:key_type] && opts[:item_type]

        value.each do |k, v|
          check_property_value("#{variable} key (#{k})", k, opts[:key_type])
          check_property_value("#{variable}[#{k}]", v, opts[:item_type])
        end
      end

      # Check if item values are allowed
      return unless opts[:item_allowed]
      if value.is_a?(::Array)
        value.each_with_index do |v, i|
          raise "#{v} on #{variable}[#{i}] should be one of #{opts[:item_allowed]}" \
            unless opts[:item_allowed].include?(v)
        end
      elsif value.is_a?(::Hash)
        value.each do |k, v|
          raise "#{v} on #{variable}[#{k}] should be one of #{opts[:item_allowed]}" \
            unless opts[:item_allowed].include?(v)
        end
      end
    end

  end
end
