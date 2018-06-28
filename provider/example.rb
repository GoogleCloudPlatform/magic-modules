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

require 'provider/config'
require 'provider/core'

module Provider
  # Code generator for Example Cookbooks that manage Google Cloud Platform
  # resources.
  class Example < Provider::Core
    # Settings for the provider
    class Config < Provider::Config
      attr_reader :manifest
      def provider
        Provider::Example
      end
    end

    private

    # This function uses the resource.erb template to create one file
    # per resource. The resource.erb template forms the basis of a single
    # GCP Resource on Example.
    def generate_resource(data)
      target_folder = data[:output_folder]
      FileUtils.mkpath target_folder
      name = Google::StringUtils.underscore(data[:object].name)
      generate_resource_file data.clone.merge(
        default_template: 'templates/example/resource.erb',
        out_file: File.join(target_folder, "#{name}.rb")
      )
    end

    # This function would generate unit tests using a template
    def generate_resource_tests(data) end

    # This function would automatically generate the files used for verifying
    # network calls in unit tests. If you comment out the following line,
    # a bunch of YAML files will be created under the spec/ folder.
    def generate_network_datas(data, object) end

    # We build a lot of property classes to help validate + coerce types.
    # The following functions would generate all of these properties.
    # Some of these property classes help us handle Strings, Times, etc.
    #
    # Others (nested objects) ensure that all Hashes contain proper values +
    # types for its nested properties.
    #
    # ResourceRefs properties help ensure that links between different objects
    # (Addresses + Instances for example) work properly, are abstracted away,
    # and don't require the user to have a large knowledge base of how GCP
    # works.
    # rubocop:disable Layout/EmptyLineBetweenDefs
    def generate_base_property(data) end
    def generate_simple_property(type, data) end
    def emit_nested_object(data) end
    def emit_resourceref_object(data) end
    # rubocop:enable Layout/EmptyLineBetweenDefs
  end
end
