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

module Provider
  # A set of functions to generate properties of objects being compiled. This is
  # a helper library to be included by Provider::Core.
  # rubocop:disable Metrics/ModuleLength
  module Properties
    private

    def generate_properties(data, properties)
      prop_map = []

      prop_map << generate_base_property(data) unless properties.empty?
      prop_map << generate_primitive_properties(data, properties)
      prop_map << generate_array_properties(data, properties)
      prop_map << generate_nested_object_properties(data, properties)
      prop_map << generate_resourceref_properties(data, properties)
      prop_map << generate_namevalues_properties(data, properties)
      prop_map << generate_enum_properties(data, properties)

      generate_property_files(prop_map, data)
    end

    # Generate the files for the properties
    def generate_property_files(prop_map, data)
      prop_map.flatten.compact.each do |prop|
        compile_file_list(
          data[:output_folder],
          { prop[:target] => prop[:source] },
          {
            product_ns: data[:product_name].camelize(:upper),
            prop_ns_dir: data[:product_name].downcase
          }.merge((prop[:overrides] || {}))
        )
      end
    end

    def generate_primitive_properties(data, properties)
      properties.select { |p| p.is_a?(Api::Type::Primitive) }
                .map { |p| generate_simple_property p.type.downcase, data }
    end

    def generate_enum_properties(_data, _properties)
      # By default, enums don't need properties.  In puppet and chef, which
      # have less sophisticated diffing logic than our other tools, it is
      # necessary to generate properties for enums where the default is
      # 'unset' - this might be a good place to introduce the concept of
      # diffs and diff suppression to puppet / chef.
      []
    end

    # rubocop:disable Metrics/AbcSize
    def generate_array_properties(data, properties)
      prop_map = []

      prop_map << properties.select { |p| p.is_a?(Api::Type::Array) }
                            .select { |p| p.item_type.is_a?(String) }
                            .map { |p| generate_typed_array(data, p) }

      prop_map \
        << properties.select { |p| p.is_a?(Api::Type::Array) }
                     .select { |p| p.item_type.is_a?(Api::Type::NestedObject) }
                     .map { |p| generate_nested_object_array(data, p) }

      prop_map \
        << properties.select { |p| p.is_a?(Api::Type::Array) }
                     .select { |p| p.item_type.is_a?(Api::Type::ResourceRef) }
                     .map { |p| generate_resourceref_array(data, p.item_type) }

      prop_map
    end
    # rubocop:enable Metrics/AbcSize

    def generate_namevalues_properties(data, properties)
      properties.select { |p| p.is_a?(Api::Type::NameValues) }
                .map { |p| generate_simple_property p.type.downcase, data }
    end

    def generate_nested_object_properties(data, properties)
      properties.select { |p| p.is_a?(Api::Type::NestedObject) }
                .map { |p| generate_nested_object(data, p) }
    end

    def generate_nested_object(data, prop)
      prop_map = []

      prop_map << emit_nested_object(
        data.clone.merge(
          emit_array: false,
          api_name: prop.name.underscore,
          property: prop,
          nested_properties: prop.properties,
          obj_name: data[:object].name.underscore
        )
      )

      prop_map << generate_properties(data, prop.properties)

      prop_map
    end

    def generate_nested_object_array(data, prop)
      prop_map = []

      prop_map << emit_nested_object(
        data.clone.merge(
          emit_array: true,
          api_name: prop.name.underscore,
          property: prop,
          nested_properties: prop.item_type.properties,
          obj_name: data[:object].name.underscore
        )
      )

      prop_map << generate_properties(data, prop.item_type.properties)

      prop_map
    end

    def generate_resourceref_properties(data, properties)
      properties.select { |p| p.is_a?(Api::Type::ResourceRef) }
                .map { |p| generate_resourceref_object(data, p) }
    end

    def generate_resourceref_object(data, prop)
      resource = prop.resource_ref.name.underscore
      imports = prop.imports.underscore
      return if resourceref_tracker.key?([resource, imports])
      resourceref_tracker[[resource, imports]] = false

      emit_resourceref_object(
        data.clone.merge(
          emit_array: false,
          property: prop,
          resource: resource,
          imports: imports
        )
      )
    end

    def generate_resourceref_array(data, prop)
      resource = prop.resource_ref.name.underscore
      imports = prop.imports.underscore
      return if resourceref_tracker.key?([resource, imports]) \
        && resourceref_tracker[[resource, imports]] == true
      resourceref_tracker[[resource, imports]] = true

      emit_resourceref_object(
        data.clone.merge(
          emit_array: true,
          property: prop,
          resource: resource,
          imports: imports
        )
      )
    end

    def resourceref_tracker
      @resourceref ||= {}
    end
  end
  # rubocop:enable Metrics/ModuleLength
end
