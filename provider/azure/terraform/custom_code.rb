require 'provider/terraform/custom_code'

module Provider
  module Azure
    module Terraform

      class CustomCode < Provider::Terraform::CustomCode
        # This code is run after the Read call succeeds and before setting
        # schema attributes. It's placed in the Read function directly
        # without modification.
        attr_reader :post_read

        # This code snippet will be put after all CRUD and expand/flatten functions
        # of a Terraform resource without modification.
        attr_reader :extra_functions

        def validate
          super
          check_optional_property :post_read, String
          check_optional_property :extra_functions, String
        end
      end

    end
  end
end
