require 'chef/resource'
require_relative 'google/authorization'

module Google
  module Auth
    # Chef Resource to authenticate to GCP
    class Credential < Chef::Resource::LWRPBase
      resource_name :gauth_credential

      default_action :nothing

      property :name, String, identity: true, desired_state: false
      property :path, String, desired_state: false
      property :scopes, Array, desired_state: false
      property :__auth, ::Google::Authorization, desired_state: false

      action :serviceaccount do
        if new_resource.path.nil?
          raise ["Missing 'path' parameter in",
                 "gauth_credential[#{new_resource.name}]"].join(' ')
        end

        if new_resource.scopes.nil?
          raise ["Missing 'scopes' parameter in",
                 "gauth_credential[#{new_resource.name}]"].join(' ')
        end

        # TODO: How do define a private property, or better, how to store a
        # variable only for this instance
        new_resource.__auth ::Google::Authorization.new.for!(
          new_resource.scopes
        ).from_service_account_json!(
          new_resource.path
        )
      end

      action :defaultuseraccount do
        __auth ::Google::Authorization.new.from_user_credential!
      end

      action :nothing do
        raise 'An action for a provider is required: service_account'
      end

      def authorization
        raise "Failed to authenticate gauth_credential[#{name}]" if __auth.nil?
        __auth
      end
    end
  end
end
