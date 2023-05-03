require 'openapi_parser'

def writeObject(name, obj, type, url_param)
  res = nil
  case type
  when "string"
    res = Api::Type::String.new(name)
  when "integer"
    res = Api::Type::Integer.new(name)
  when "number"
    res = Api::Type::Double.new(name)
  when "boolean"
    res = Api::Type::Boolean.new()
    res.instance_variable_set(:@name, name)
  when "object"
    if name == "labels"
      # standard labels field handling
      res = Api::Type::KeyValuePairs.new()
    else
      res = Api::Type::NestedObject.new()
      req = obj.required || []

      pr = []
      if obj.properties
        obj.properties.each do |prop, i|
          prop = writeObject(prop, i, i.type, false)
          prop.instance_variable_set(:@required, req.include?(prop.name))
          pr.push(prop)
        end
      end
      res.instance_variable_set(:@properties, pr)
    end
    res.instance_variable_set(:@name, name)

  when "array"
    res = Api::Type::Array.new()
    res.instance_variable_set(:@name, name)
    if obj.items.type == "string"
        res.instance_variable_set(:@item_type, "Api::Type::String")
    else
      no = Api::Type::NestedObject.new()
      pr = []
      if obj.items.properties
        obj.items.properties.each do |prop, i|
          prop = writeObject(prop, i, i.type, false)
          pr.push(prop)
        end
      end
      no.instance_variable_set(:@properties, pr)
      res.instance_variable_set(:@item_type, no)
    end
  else
    raise "Failed to identify field type #{type}"
    return nil
  end
  res.instance_variable_set(:@description, obj.description || "No description")
  if url_param
    res.instance_variable_set(:@url_param_only, true)
    res.instance_variable_set(:@required, obj.required)
  end

  # These methods are only available when the field is set
  if obj.respond_to?(:read_only) && obj.read_only
    res.instance_variable_set(:@output, obj.read_only)
  end

  if obj.respond_to?(:write_only) && obj.write_only
    res.instance_variable_set(:@immutable, obj.write_only)
  end

  return res
end

def find_resources(spec_path)
  resource_paths = []
  root = OpenAPIParser.parse(YAML.load_file(spec_path))
  root.paths.path.each do |path|
    if path[1].post
      # Not very clever way of identifying create resource methods
      if path[1].post.operation_id.start_with?("Create")
        resource_paths.push([path[0], path[1].post.operation_id.gsub("Create", "")])
      end
    end
  end
  return resource_paths
end

def parse_openapi(spec_path, resource_path)
  # Write YAML
  root = OpenAPIParser.parse(YAML.load_file(spec_path))
  op = root.request_operation(:post, resource_path)
  path = root.paths.path[resource_path]
  parameters = []
  path.post.parameters.each do |param|
    parameters.push(writeObject(param.name, param, param.schema.type, true))
  end
  properties = []
  path.post.request_body.content["application/json"].schema.properties.each do |prop, i|
    properties.push(writeObject(prop, i, i.type, false))
  end
  return properties, parameters, path.post.parameters.last.name
end

def base_url(resource_path)
  base = resource_path.gsub("{", "{{").gsub("}", "}}")
  field_names = base.scan(/(?<=\{\{)\w+(?=\}\})/)
  field_names.each do |field_name|
    field_name_in_snake_case = field_name.underscore
    base = base.gsub("{{#{field_name}}}", "{{#{field_name_in_snake_case}}}")
  end
  return base
end

def write_resource(spec_path, resource_path, resource_name)
  properties, parameters, query_param = parse_openapi(spec_path, resource_path)

  resources = []
  resource = Api::Resource.new()
  base_url = base_url(resource_path)
  resource.base_url = base_url
  resource.create_url = "#{base_url}?#{query_param}={{#{query_param.underscore}}}"
  resource.self_link = "#{base_url}/{{#{query_param.underscore}}}"

  # Name is on the Api::Object::Named parent resource, lets not modify that
  resource.instance_variable_set(:@name, resource_name)
  # TODO(slevenick): Get resource description published in OpenAPI spec
  resource.description = "Description"
  resource.update_verb = :PATCH
  resource.update_mask = true
  resource.autogen_async = true
  resource.properties = properties
  resource.parameters = parameters


  # Default operation handling
  op = Api::OpAsync::Operation.new("name", "{{op_id}}", 1000, nil)
  result = Api::OpAsync::Result.new("response", true)
  status = Api::OpAsync::Status.new("done", true, [true, false])
  error = Api::OpAsync::Error.new("error", "message")
  async = Api::OpAsync.new(op, result, status, error)
  resource.async = async

  resources.push(resource)
  file_path = File.join("products/demo", "#{resource_name}.yaml")
  File.open(file_path, 'w') { |file| file.write(resource.to_yaml) }
end

def write_yaml(spec_path)
  resource_paths = find_resources(spec_path)
  resource_paths.each do |path_array|
    write_resource(spec_path, path_array[0], path_array[1])
  end
end