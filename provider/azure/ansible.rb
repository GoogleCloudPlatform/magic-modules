require 'provider/azure/ansible/sub_template'
require 'provider/azure/ansible/sdk/sub_template'

module Provider
  module Azure
    module Ansible
      include Provider::Azure::Ansible::SubTemplate
      include Provider::Azure::Ansible::SDK::SubTemplate
    end
  end
end