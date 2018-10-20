require 'yaml'
require 'provider/resource_overrides'
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
      @override.consume_config(@product, DiscoveryConfig)
      @override.run
    end
  end

  module OverrideProperties
  end

  class DiscoveryConfig
    def self.resource_override
      ResourceOverride
    end

    def self.property_override
      PropertyOverride
    end
  end

  # Product specific overriden properties for Ansible
  class ResourceOverride < Provider::ResourceOverride
    include OverrideProperties
    def validate
    end

    def properties
      @properties || {}
    end

    private

    def overriden
      DiscoveryOverride::OverrideProperties
    end
  end

  class PropertyOverride < Provider::PropertyOverride
    include OverrideProperties
    def validate
    end

    def properties
      @properties || {}
    end

    private

    def overriden
      DiscoveryOverride::OverrideProperties
    end
  end
end
