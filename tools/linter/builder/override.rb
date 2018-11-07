require 'yaml'
require 'provider/overrides/runner'
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
      runner.build
    end
  end
end
