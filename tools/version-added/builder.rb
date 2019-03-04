$LOAD_PATH.unshift(File.dirname(__FILE__) + '/../..')
Dir.chdir(File.dirname(__FILE__) + '/../..')

require 'api/compiler'
require 'api/product'
require 'api/resource'
require 'provider/config'
require 'provider/ansible'
require 'overrides/runner'


CURRENT_ANSIBLE_VERSION = '2.8'

def version(path, struct)
  path = path.map(&:to_sym) + [:version_added]
  struct.dig(*path) || CURRENT_ANSIBLE_VERSION
end

# Builds out property information (with nesting)
def property(prop, path, struct)
  property_hash = {
    version_added: version(path + [prop.name], struct)
  }

  prop.nested_properties.each do |nested_p|
    property_hash[nested_p.name.to_sym] = property(nested_p, path + [prop.name], struct)
  end
  property_hash
end

products = Dir["products/**/ansible.yaml"].map { |x| x.split('/')[1] }
for product in products
  # Get api.yaml
  product_yaml_path = "products/#{product}/api.yaml"
  provider_yaml_path = "products/#{product}/ansible.yaml"
  version_yaml_path = "products/#{product}/ansible_version_added.yaml"
  product_api = Api::Compiler.new(product_yaml_path).run
  product_api.validate

  # Build overrides.
  product_api, _ = Provider::Config.parse(provider_yaml_path, product_api, 'ga')


  versions = if File.exists?(version_yaml_path)
               YAML.load(File.read(version_yaml_path))
             else
               {}
             end

  struct = {
    facts: {},
    regular: {}
  }

  # Build out paths for regular modules.
  product_api.objects.each do |obj|
    resource = {
      version_added: version([:regular, obj.name], versions)
    }

    # Add properties.
    obj.all_user_properties.each do |prop|
      resource[prop.name.to_sym] = property(prop, [:regular, obj.name], versions)
    end
    struct[:regular][obj.name.to_sym] = resource

    # Add facts modules from facts datasources.
    struct[:facts][obj.name.to_sym] = {
      version_added: version([:facts, obj.name], versions)
    }
  end

  # Add facts modules.
  File.write("products/#{product}/ansible_version_added.yaml", struct.to_yaml)
end
