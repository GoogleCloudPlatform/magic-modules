# Load everything from MM root.
$LOAD_PATH.unshift File.join(File.dirname(__FILE__), "../../")
Dir.chdir(File.join(File.dirname(__FILE__), "../../"))

require 'tools/linter/discovery'

COMPUTE_DISCOVERY_URL = 'https://www.googleapis.com/discovery/v1/apis/compute/v1/rest'
discovery = Discovery.new(COMPUTE_DISCOVERY_URL)
puts discovery.results
