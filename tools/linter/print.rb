$LOAD_PATH.unshift File.join(File.dirname(__FILE__), '../../')
Dir.chdir(File.join(File.dirname(__FILE__), '../../'))

require 'api/resource'
require 'api/product'
require 'api/type'
require 'api/compiler'

def alter_object(obj)
  obj.remove_instance_variable(:@description) if obj.instance_variable_get(:@description)
  obj.instance_variables.each do |inst_var|
    obj.remove_instance_variable(inst_var) unless obj.instance_variable_get(inst_var)
    if obj.is_a?(Api::Resource)
      obj.properties.sort_by! { |a| a.name }
      obj.parameters.sort_by! { |a| a.name }
      obj.all_user_properties.each { |x| alter_object(x) } if obj.is_a?(Api::Resource)
    end
    if obj.is_a?(Api::Type::NestedObject)
      obj.properties.sort_by! { |a| a.name }
      obj.properties.each { |x| alter_object(x) } if obj.is_a?(Api::Type::NestedObject)
    end

    if obj.is_a?(Api::Type::Enum)
      obj.values.sort!
    end
    if obj.is_a?(Api::Type::Array) && obj.item_type.is_a?(Api::Type::NestedObject)
      obj.item_type.properties.sort_by! { |a| a.name }
      obj.item_type.properties.each { |x| alter_object(x) }
    end
  end
end

raise "Must include four file locations" if ARGV.length != 4
file1 = {
  original: ARGV[0],
  new: ARGV[1]
}

file2 = {
  original: ARGV[2],
  new: ARGV[3]
}

[file1, file2].each do |file|
  product_api = Api::Compiler.new(file[:original]).run
  product_api.objects.each { |obj| alter_object(obj) }
  File.write(file[:new], YAML::dump(product_api))
end
