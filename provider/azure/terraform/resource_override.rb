require 'provider/terraform/resource_override'

module Provider
  module Azure
    module Terraform
      module OverrideProperties
        attr_reader :acctests
        include Provider::Terraform::OverrideProperties
      end

      class ResourceOverride < Provider::Terraform::ResourceOverride
        include Provider::Azure::Terraform::OverrideProperties

        def validate
          super
          check_optional_property :acctests, Hash
          check_optional_property_hash :acctests, String, AccTestDefinition
          post_initialization
        end

        def apply(resource)
          super
        end

        class AccTestDefinition < Api::Object
          attr_reader :based_on
          attr_reader :steps
          attr_reader :hcl_parameters
          attr_reader :check_import
          attr_reader :custom_dependencies_code
          attr_reader :properties

          def validate
            super
            @initialized = false

            check_property :steps, Array
            check_property_list :steps, AccTestStepDefinition

            check_optional_property :based_on, String
            check_optional_property :check_import, :boolean
            @check_import = true if @check_import.nil?
            check_optional_property :custom_dependencies_code, String

            check_optional_property :properties, Hash
            @properties ||= Hash.new
            check_optional_property_hash :properties, String, String

            check_optional_property :hcl_parameters, Array
            @hcl_parameters ||= []
            convert_hcl_parameter_strings
            check_optional_property_list :hcl_parameters, AccTestHCLParametersDefinition
          end

          def post_initialization(resource)
            merge_hcl_parameters(resource) unless @initialized
            @initialized = true
          end

          private

          def convert_hcl_parameter_strings()
            @hcl_parameters.each_index do |i|
              @hcl_parameters[i] = AccTestHCLParametersDefinition.new(@hcl_parameters[i]) if @hcl_parameters[i].instance_of? String
            end
          end

          def merge_hcl_parameters(resource)
            parent_def_id = @based_on
            until parent_def_id.nil?
              parent_def = resource.acctests[parent_def_id]
              parent_def.post_initialization(resource)
              parent_def.hcl_parameters.reverse_each{|param| @hcl_parameters.insert(0, param)}
              parent_def_id = parent_def.based_on
            end
          end
        end

        class AccTestHCLParametersDefinition < Api::Object
          attr_reader :variable_name
          attr_reader :go_type
          attr_reader :create_expression

          def validate
            super
            check_property :variable_name, String
            check_property :go_type, String
            check_property :create_expression, String
          end

          def initialize(str)
            case str
            when "AccDefaultInt"
              @variable_name = "rInt"
              @go_type = "int"
              @create_expression = "tf.AccRandTimeInt()"
            when "AccLocation"
              @variable_name = "location"
              @go_type = "string"
              @create_expression = "testLocation()"
            else
              raise "#{str} is not a valid pre-defined AccTestHCLParametersDefinition"
            end
          end
        end

        class AccTestStepDefinition < Api::Object
          attr_reader :config_name
          attr_reader :hcl_reference

          def validate
            super
            check_property :config_name, String
            check_property :hcl_reference, String
          end
        end

        private

        def overriden
          Provider::Azure::Terraform::OverrideProperties
        end

        def post_initialization
          @acctests.each_value{|acctest_def| acctest_def.post_initialization(self)}
        end

      end
    end
  end
end
