$LOAD_PATH.unshift(File.dirname(__FILE__) + '/../..')
Dir.chdir(File.dirname(__FILE__) + '/../..')

require 'api/compiler'
require 'api/product'
require 'api/resource'
require 'provider/config'
require 'provider/ansible'
require 'overrides/runner'

products = [
  'bigquery',
  'cloudbuild',
  'compute',
  'container',
  'dns',
  'iam',
  'pubsub',
  'redis',
  'resourcemanager',
  'sourcerepo',
  'spanner',
  'sql',
  'storage'
]

CURRENT_ANSBILE_VERSION = '2.8'

class Version
  def initialize(product)
    version_added_path = "products/#{product}/ansible_version_added.yaml"
    if File.exists?(version_added_path)
      @version_added_file = YAML.load(File.read(version_added_path))
    end
  end

  def version(object, higher_level)
    return higher_level unless object
    return object if object.to_f > higher_level.to_f
    higher_level
  end
end

def version(object, higher_level)
  return higher_level unless object
  return object if object.to_f > higher_level.to_f
  higher_level
end

def property(prop, resource_version_added)
  version_for_prop = version(prop.version_added, resource_version_added)
  property_hash = {
    version_added: version_for_prop,
  }

  prop.nested_properties.each do |nested_p|
    property_hash[nested_p.name.to_sym] = property(nested_p, version_for_prop)
  end
  property_hash
end

for product in products
  # Get api.yaml
  product_yaml_path = ("products/#{product}/api.yaml")
  provider_yaml_path = ("products/#{product}/ansible.yaml")
  product_api = Api::Compiler.new(product_yaml_path).run
  product_api.validate

  # Get ansible.yaml
  product_api, provider_config = Provider::Config.parse(provider_yaml_path, product_api, 'ga')

  struct = {
    facts: {
    },
    regular: {}
  }

  # Build out paths for regular modules.
  product_api.objects.each do |obj|
    resource = {
      version_added: provider_config.manifest.get('version_added', obj)
    }

    # Add properties.
    obj.all_user_properties.each do |prop|
      resource[prop.name.to_sym] = property(prop, resource[:version_added])
    end
    struct[:regular][obj.name.to_sym] = resource
    struct[:facts][obj.name.to_sym] = {
      version_added: version(provider_config.datasources.instance_variable_get("@#{obj.name}").instance_variable_get('@version_added'), provider_config.manifest.get('version_added', obj))
    }
  end

  # Add facts modules.
  File.write("products/#{product}/ansible_version_added.yaml", struct.to_yaml)
end
