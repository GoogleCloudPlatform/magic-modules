# Copyright 2017 Google Inc.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

require 'provider/abstract_core'

module Provider
  class Terraform < Provider::AbstractCore
    # Functions to compile sub-templates.
    module SubTemplate
      def build_schema_property(property, object, indentation = 0)
        compile_template schema_property_template(property),
                         indentation: indentation,
                         prop_name: property.name.underscore,
                         property: property,
                         object: object
      end

      # Transforms a Cloud API representation of a property into a Terraform
      # schema representation.
      def build_flatten_method(property, sdk_marshal)
        compile_template 'templates/terraform/flatten_property_method.erb',
                         property: property,
                         sdk_marshal: sdk_marshal
      end

      # Transforms a Terraform schema representation of a property into a
      # representation used by the Cloud API.
      def build_expand_method(property, sdk_marshal)
        compile_template 'templates/terraform/expand_property_method.erb',
                         property: property,
                         sdk_marshal: sdk_marshal
      end

      def build_expand_resource_ref(var_name, property)
        compile_template 'templates/terraform/expand_resource_ref.erb',
                         var_name: var_name,
                         property: property
      end

      def build_property_documentation(property, is_data_source = false)
        compile_template 'templates/terraform/property_documentation.erb',
                         property: property,
                         is_data_source: is_data_source
      end

      def build_nested_property_documentation(property)
        compile_template(
          'templates/terraform/nested_property_documentation.erb',
          property: property
        )
      end

      private

      def autogen_notice_contrib
        ['Please read more about how to change this file in',
         '.github/CONTRIBUTING.md.']
      end

      def autogen_notice_text(line)
        line&.empty? ? '//' : "// #{line}"
      end

      def compile_template(template_file, data)
        ctx = binding
        data.each { |name, value| ctx.local_variable_set(name, value) }
        result = compile_file(ctx, template_file).join("\n")
        indent result, data[:indentation] || 0
      end
    end
  end
end
