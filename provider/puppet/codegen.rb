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

require 'google/hash_utils'
require 'provider/core'

module Provider
  class Puppet < Provider::Core
    # Functions that generate code for Puppet modules
    module Codegen
      private

      def generate_type(data)
        target_folder = File.join(data[:output_folder], 'lib', 'puppet', 'type')
        FileUtils.mkpath target_folder
        generate_resource_file data.clone.merge(
          type: 'type',
          default_template: 'templates/puppet/type.erb',
          out_file: File.join(target_folder, "#{data[:name]}.rb")
        )
      end

      def generate_provider(data)
        target_folder = File.join(data[:output_folder], 'lib', 'puppet',
                                  'provider', data[:name])
        FileUtils.mkpath target_folder
        generate_resource_file data.clone.merge(
          type: 'provider',
          default_template: provider_template_source(data),
          out_file: File.join(target_folder, 'google.rb')
        )
      end

      def generate_resource_tests(data)
        return if true?(data[:object].manual)
        generate_provider_tests data
      end

      def generate_provider_tests(data)
        generate_resource_file data.clone.merge(
          type: 'provider_spec',
          default_template: 'templates/puppet/provider_spec.erb',
          out_file: File.join(data[:output_folder], 'spec',
                              "#{data[:name]}_provider_spec.rb")
        )
      end

      def provider_template_source(data)
        object_name = data[:object].name.underscore
        if true?(data[:object].manual)
          File.join('products', data[:product_name], 'files',
                    "provider~#{object_name}.rb")
        else
          'templates/puppet/resource.erb'
        end
      end
    end
  end
end
