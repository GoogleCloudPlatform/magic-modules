require 'provider/ansible/config'

module Provider
  module Azure
    module Ansible

      class Config < Provider::Ansible::Config
        attr_reader :author
        attr_reader :version_added

        def provider
          Provider::Ansible::Core
        end
  
        def resource_override
          Provider::Azure::Ansible::ResourceOverride
        end
  
        def property_override
          Provider::Azure::Ansible::PropertyOverride
        end

        def validate
          super
          check :author, type: ::String, required: true
          check :version_added, type: ::String, required: true
        end
      end

    end
  end
end
