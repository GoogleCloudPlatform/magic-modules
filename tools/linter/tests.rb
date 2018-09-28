# Load everything from MM root.
$LOAD_PATH.unshift File.join(File.dirname(__FILE__), "../../")
Dir.chdir(File.join(File.dirname(__FILE__), "../../"))

require 'tools/linter/discovery'
require 'tools/linter/api'
require 'tools/linter/test_helpers'

require 'set'

COMPUTE_DISCOVERY_URL = 'https://www.googleapis.com/discovery/v1/apis/compute/v1/rest'
discovery = Discovery.new(COMPUTE_DISCOVERY_URL)
api = ProductApi.new('compute')

# Print out all properties in discovery, but not in api.
puts "INFO: Finding all missing top level properties on api.yaml"

loop_resources_in_api(api, discovery) do |api_obj, disc_obj|
  unless disc_obj.exists?
    puts "WARN: #{api_obj.name} - not found in Discovery"
    next
  end
  api_prop_names = api_obj.all_user_properties.map(&:name)
  disc_prop_names = disc_obj.schema['properties'].keys
  missing_props = (Set.new(disc_prop_names) - Set.new(api_prop_names)) - Set.new(%w[kind selfLink])
  unless missing_props.to_a.empty?
    puts "ERROR: #{api_obj.name}- #{missing_props.to_a.join(', ')} are in discovery, but not API"
  end
end
