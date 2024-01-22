# Copyright 2023 Google Inc.
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

require 'openapi_parser'

module OpenAPIGenerate
  # Parser to convert from OpenAPI spec to MMv1 YAML
  class Parser
    attr_reader :folder
    attr_reader :output

    def initialize(folder, output)
      @folder = folder
      @output = output
    end

    def run
      Dir[@folder].each do |openapi_file|
        write_yaml(openapi_file, @output)
      end
    end

    def write_object(name, obj, type, url_param)
      field = nil
      case name
      when 'projectsId'
        return field
      when 'locationsId'
        name = 'location'
      end
      additional_description = ''

      # allOf is a workaround for overriding fields on shared objects
      if obj.respond_to?(:all_of) && !obj.all_of&.length().nil?
        obj = obj.all_of[0]
        type = obj.type
      end

      case type
      when 'string'
        field = Api::Type::String.new(name)
        if obj.respond_to?(:enum) && obj.enum
          additional_description = "\n Possible values:\n #{obj.enum.join("\n")}"
        end
      when 'integer'
        field = Api::Type::Integer.new
        field.instance_variable_set(:@name, name)
      when 'number'
        field = Api::Type::Double.new
        field.instance_variable_set(:@name, name)
      when 'boolean'
        field = Api::Type::Boolean.new
        field.instance_variable_set(:@name, name)
      when 'object'
        if name == 'labels'
          # standard labels field handling
          field = Api::Type::KeyValueLabels.new
        elsif name == 'annotations'
          # standard annotations field handling
          field = Api::Type::KeyValueAnnotations.new
        elsif obj.respond_to?(:additional_properties) \
          && obj.additional_properties.respond_to?(:type)
          # additionalProperties.type signifies a string -> string map
          field = Api::Type::KeyValuePairs.new
        else
          field = Api::Type::NestedObject.new
          required_props = obj.required || []

          properties = []
          obj.properties&.each do |prop, i|
            prop = write_object(prop, i, i.type, false)
            prop.instance_variable_set(:@required, true) if required_props.include?(prop.name)
            required_props.delete(prop.name)
            properties.push(prop)
          end
          raise "Unknown required properties #{required_props}" unless required_props.empty?

          field.instance_variable_set(:@properties, properties)
        end
        field.instance_variable_set(:@name, name)

      when 'array'
        field = Api::Type::Array.new
        field.instance_variable_set(:@name, name)
        case obj.items.type
        when 'string'
          field.instance_variable_set(:@item_type, 'Api::Type::String')
        when 'number'
          field.instance_variable_set(:@item_type, 'Api::Type::Double')
        when 'boolean'
          field.instance_variable_set(:@item_type, 'Api::Type::Boolean')
        else
          nested_object = Api::Type::NestedObject.new
          object_properties = build_properties(
            obj.items.properties,
            obj.items.required || []
          )
          nested_object.instance_variable_set(:@properties, object_properties)
          field.instance_variable_set(:@item_type, nested_object)
        end
      else
        raise "Failed to identify field type #{type} #{name}"
      end

      field.instance_variable_set(
        :@description,
        "#{obj.description} #{additional_description}" || 'No description'
      )
      if url_param
        field.instance_variable_set(:@url_param_only, true)
        field.instance_variable_set(:@required, true) if obj.required
      end

      # These methods are only available when the field is set
      if obj.respond_to?(:read_only) && obj.read_only
        field.instance_variable_set(:@output, obj.read_only)
      end

      if (obj.respond_to?(:write_only) && obj.write_only) \
        || obj.instance_variable_get(:@raw_schema)['x-google-immutable']
        field.instance_variable_set(:@immutable, true)
      end

      field
    end

    def find_resources(spec_path)
      resource_paths = []
      root = OpenAPIParser.parse(YAML.load_file(spec_path))
      root.paths.path.each do |path|
        next unless path[1].post

        # Not very clever way of identifying create resource methods
        if path[1].post.operation_id.start_with?('Create')
          resource_paths.push([path[0], path[1].post.operation_id.gsub('Create', '')])
        end
      end
      resource_paths
    end

    def parse_openapi(spec_path, resource_path, resource_name)
      # Write YAML
      root = OpenAPIParser.parse(YAML.load_file(spec_path))
      path = root.paths.path[resource_path]
      parameters = []
      path.post.parameters.each do |param|
        parameter_object = write_object(param.name, param, param.schema.type, true)
        # Ignore standard requestId field
        next if param.name == 'requestId'
        next if parameter_object.nil?

        # All parameters are immutable
        parameter_object.instance_variable_set(:@immutable, true)
        parameters.push(parameter_object)
      end
      properties = build_properties(
        path.post.request_body.content['application/json'].schema.properties,
        path.post.request_body.content['application/json'].schema.required || []
      )

      id_param = path.post.parameters.select do |p|
        p.name.downcase.include?(resource_name.downcase)
      end.last
      raise 'did not find ID param' unless id_param

      [properties, parameters, id_param.name]
    end

    def build_properties(properties, required)
      prop_objects = []
      properties&.each do |prop, i|
        prop_object = write_object(prop, i, i.type, false)
        prop_object.instance_variable_set(:@required, true) if required.include?(prop)

        required.delete(prop)
        prop_objects.push(prop_object)
      end
      raise "Unknown required properties in object #{required}" unless required.empty?

      prop_objects
    end

    def base_url(resource_path)
      base = resource_path.gsub('{', '{{').gsub('}', '}}')

      base = base.gsub('projectsId', 'project')
      base = base.gsub('locationsId', 'location')
      field_names = base.scan(/(?<=\{\{)\w+(?=\}\})/)
      field_names.each do |field_name|
        field_name_in_snake_case = field_name.underscore
        base = base.gsub("{{#{field_name}}}", "{{#{field_name_in_snake_case}}}")
      end
      base = base.gsub('/v1/', '')
      base.gsub('/v1alpha/', '')
    end

    def build_resource(spec_path, resource_path, resource_name)
      properties, parameters, query_param = parse_openapi(spec_path, resource_path, resource_name)

      resource = Api::Resource.new
      base_url = base_url(resource_path)
      resource.base_url = base_url
      resource.create_url = "#{base_url}?#{query_param}={{#{query_param.underscore}}}"
      self_link = "#{base_url}/{{#{query_param.underscore}}}"
      resource.self_link = self_link
      resource.id_format = self_link
      resource.import_format = [self_link]

      # Name is on the Api::NamedObject parent resource, lets not modify that
      resource.instance_variable_set(:@name, resource_name)
      # TODO(slevenick): Get resource description published in OpenAPI spec
      resource.description = 'Description'
      if update?(spec_path, resource_name)
        resource.update_verb = :PATCH
        resource.update_mask = true
      else
        resource.immutable = true
      end

      resource.autogen_async = true
      resource.properties = properties
      resource.parameters = parameters

      # Default operation handling
      op = Api::OpAsync::Operation.new('name', '{{op_id}}', 1000, nil)
      result = Api::OpAsync::Result.new('response', true)
      status = Api::OpAsync::Status.new('done', true, [true, false])
      error = Api::OpAsync::Error.new('error', 'message')
      async = Api::OpAsync.new(op, result, status, error)
      resource.async = async
      resource
    end

    def update?(spec_path, resource_name)
      root = OpenAPIParser.parse(YAML.load_file(spec_path))
      root.paths.path.each do |path|
        # PATCH is the standard update method
        next unless path[1].patch

        return true if path[1].patch.operation_id.start_with?("Update#{resource_name}")
      end
      false
    end

    def build_product(spec_path, output)
      root = OpenAPIParser.parse(YAML.load_file(spec_path))
      version = root.raw_schema['info']['version']
      server = root.raw_schema['servers'][0]['url']
      product_name = spec_path.split('/').last.split('_').first
      product_path = File.join(output, product_name)
      FileUtils.mkdir_p(product_path)
      product = Api::Product.new
      api_version = Api::Product::Version.new
      api_version.base_url = "#{server}/#{version}/"
      # TODO(slevenick) figure out how to tell the API version
      api_version.name = 'ga'
      product.versions = [api_version]
      # Standard titling is "Service Name API"
      display_name = root.raw_schema['info']['title'].sub(' API', '')
      # Name is on the Api::NamedObject parent resource, lets not modify that
      product.instance_variable_set(:@name, display_name.gsub(' ', ''))
      product.display_name = display_name
      # Scopes should be added soon to OpenAPI, until then use global scope
      product.scopes = ['https://www.googleapis.com/auth/cloud-platform']
      File.write(File.join(output, "/#{product_name}/product.yaml"), product.to_yaml)
      product_path
    end

    def write_yaml(spec_path, output)
      resource_paths = find_resources(spec_path)
      product_path = build_product(spec_path, output)
      resource_paths.each do |path_array|
        resource = build_resource(spec_path, path_array[0], path_array[1])
        file_path = File.join(product_path, "#{resource.name}.yaml")
        File.write(file_path, resource.to_yaml)
      end
    end
  end
end
