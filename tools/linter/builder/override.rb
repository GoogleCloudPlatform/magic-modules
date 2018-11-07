require 'yaml'
require 'provider/overrides/runner'
require 'provider/resource_override'
require 'provider/property_override'
module DiscoveryOverride
  class Runner
    attr_reader :product

    def initialize(product, override_filename)
      @product = product
      return unless override_filename
      @override = YAML.load(File.read(override_filename))
    end

    def run
      return unless @override
      runner = Provider::Overrides::Runner.new(@product, @override)
      runner.run
    end
  end
end
