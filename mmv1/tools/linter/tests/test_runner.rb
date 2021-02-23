# Copyright 2018 Google Inc.
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

require 'tools/linter/tests/tests'

def run_tests(discovery_doc, api, filters, tags = {}, **kwargs)
  product = kwargs[:api_name] || api.api_name
  RSpec.describe product, product: product do
    discovery_doc.resources.each do |disc_resource|
      api_obj = api&.objects&.select { |p| p.name == disc_resource.name }&.first
      # Second context: resource name
      describe disc_resource.name, resource: disc_resource.name do
        # Run all resource tests on this resource
        include_examples 'resource_tests', disc_resource, api_obj, tags if filters[:resource]

        if filters[:property]
          PropertyFetcher.fetch_property_pairs(disc_resource.properties,
                                               api_obj&.all_user_properties) \
                                              do |disc_prop, api_prop, name|
            # Third context: property name
            context name do
              # Run all tests on this property
              include_examples 'property_tests', disc_prop, api_prop, tags
            end
          end
        end
      end
    end
  end
end
