# Copyright 2017 Google Inc.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

require 'google/logger'
require 'yaml'

module Google
  # A helper class to validate contents coming from YAML files.
  class YamlValidator
    class << self
      def parse(content)
        # TODO(nelsonjr): Allow specifying which symbols to restrict it further.
        # But it requires inspecting all configuration files for symbol sources,
        # such as Enum values. Leaving it as a nice-to-have for the future.
        YAML.safe_load(content, allowed_classes)
      end

      def allowed_classes
        ObjectSpace.each_object(Class).select do |klass|
          klass < Google::YamlValidator
        end.concat([Time, Symbol])
      end
    end

    def validate
      Google::LOGGER.debug "Validating #{self.class} '#{@name}'"
      check_extraneous_properties
    end

    def set_variable(value, property)
      Google::LOGGER.debug "Setting variable of #{value} to #{self}"
      ensure_property_does_not_exist property
      instance_variable_set("@#{property}", value)
    end

    private

    def check_types(objects, type)
      return if objects.nil?
      objects.each do |object|
        log_check_type object
        check_type object, type
      end
    end

    def check_type(name, object, type)
      if type == :boolean
        return unless [TrueClass, FalseClass].find_index(object.class).nil?
      elsif type.is_a? ::Array
        return unless type.find_index(object.class).nil?
      elsif object.is_a?(type)
        return
      end
      raise "Property '#{name}' is '#{object.class}' instead of '#{type}'"
    end

    def log_check_type(object)
      if object.respond_to?(:name)
        Google::LOGGER.debug "Checking object #{object.name}"
      else
        Google::LOGGER.debug "Checking object #{object}"
      end
    end

    def check_property(property, type = nil)
      check_property_value property, instance_variable_get("@#{property}"), type
    end

    def check_optional_property(property, type = nil)
      value = instance_variable_get("@#{property}")
      return if value.nil?
      check_property_value property, value, type
    end

    def check_property_value(property, prop_value, type)
      Google::LOGGER.debug "Checking '#{property}' on #{object_display_name}"
      raise "Missing '#{property}' on #{object_display_name}" if prop_value.nil?
      check_type property, prop_value, type unless type.nil?
      prop_value.validate if prop_value.is_a?(Api::Object)
    end

    def default_value_property(property, def_value)
      value = instance_variable_get("@#{property}")
      instance_variable_set("@#{property}", def_value) if value.nil?
    end

    def check_extraneous_properties
      instance_variables.each do |variable|
        var_name = variable.id2name[1..-1]
        next if var_name.start_with?('__')
        Google::LOGGER.debug "Validating '#{var_name}' on #{object_display_name}"
        raise "Extraneous variable '#{var_name}' in #{object_display_name}" \
          unless methods.include?(var_name.intern)
      end
    end

    def check_property_list(name, type = nil)
      obj_list = instance_variable_get("@#{name}")
      if obj_list.nil?
        Google::LOGGER.debug "No next level @ #{object_display_name}: #{name}"
      else
        Google::LOGGER.debug \
          "Checking next level for #{object_display_name}: #{name}"
        obj_list.each { |o| check_property_value "#{name}:item", o, type }
      end
    end

    def check_property_non_empty_list(name, type = nil)
      check_property_list(name, type)
      obj_list = instance_variable_get("@#{name}")
      raise "Property #{name} is empty" if obj_list.empty?
    end

    def check_optional_property_list(name, type = nil)
      obj_list = instance_variable_get("@#{name}")
      return if obj_list.nil?
      check_property_list(name, type)
    end

    def check_property_hash(name, keyType = nil, valType = nil)
      obj_hash = instance_variable_get("@#{name}")
      if obj_hash.nil?
        Google::LOGGER.debug "No next level @ #{object_display_name}: #{name}"
      else
        Google::LOGGER.debug \
          "Checking next level for #{object_display_name}: #{name}"
          obj_hash.each do |k, v|
            check_property_value "#{name}:key", k, keyType
            check_property_value "#{name}:value", v, valType
          end
      end
    end

    def check_optional_property_hash(name, keyType = nil, valType = nil)
      obj_hash = instance_variable_get("@#{name}")
      return if obj_hash.nil?
      check_property_hash(name, keyType, valType)
    end

    # Verifies if a property is of a given type and its value are one of the
    # valid possibilities.
    def check_property_oneof(property, valid_values, type = nil)
      check_property(property, type)
      prop_value = instance_variable_get("@#{property}")
      raise "Invalid #{property} #{prop_value}" \
        unless valid_values.include?(prop_value)
    end

    # Similar to 'check_property_oneof' but assigns a default value if missing.
    def check_property_oneof_default(prop, valid_values, default, type = nil)
      prop_value = instance_variable_get("@#{prop}")
      instance_variable_set("@#{prop}", if prop_value.nil?
                                          default
                                        else
                                          prop_value
                                        end)
      check_property_oneof prop, valid_values, type
    end

    # Similar to 'check_property_oneof_default' but allows the default value to
    # also be nil with the property being effectively optional.
    def check_optional_property_oneof_default(prop, valid_values, default,
                                              type = nil)
      prop_value = instance_variable_get("@#{prop}")
      return if prop_value.nil? && default.nil?
      check_property_oneof_default(prop, valid_values, default, type)
    end

    def set_variables(objects, property)
      return if objects.nil?
      objects.each do |object|
        object.set_variable(self, property) if object.respond_to?(:set_variable)
      end
    end

    def ensure_property_does_not_exist(property)
      raise "Conflict of property '#{property}' for object '#{self}'" \
        unless instance_variable_get("@#{property}").nil?
    end

    def object_display_name
      "#{@name}<#{self.class.name}>"
    end
  end
end
