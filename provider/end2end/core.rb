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

require 'binding_of_caller'
require 'tools/end2end/constants'

module Provider
  module End2End
    # Provides functionality for generating a set of parallelizable end-to-end
    # tests.
    module Core
      include ::End2End::Constants
      def compile_end2end_tests(output_folder)
        compile_file_map(
          output_folder,
          @config.examples,
          lambda do |_object, file|
            # Tests go into hidden folder because we don't need to expose
            # to regular users.
            ["#{TEST_FOLDER}/#{file}",
             "products/#{@api.prefix[1..-1]}/files/examples~#{file}"]
          end
        )
      end

      # Returns a parallelizable name for end-to-end test manifests
      # Returns a normal name for example manifests
      # Requires:
      #   * name - The name of the resource without any prefix
      #   * The prior stack frame with a variable "name" that represents
      #     the file name.
      def example_resource_name(name)
        filename = binding.of_caller(1).local_variable_get(:name)

        res_name = name
        res_name = "#{provider_name}-e2e-#{name}" \
          if TEST_FILE_REGEX.any? { |f| filename =~ f }

        quote_string(res_name)
      end
    end
  end
end
