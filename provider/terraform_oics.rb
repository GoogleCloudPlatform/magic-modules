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

require 'provider/terraform'

module Provider
  # Code generator for runnable Terraform examples that can be run via an
  # Open in Cloud Shell link.
  class TerraformOiCS < Provider::Terraform
    # We don't want *any* static generation, so we override generate to only
    # generate objects.
    def generate(output_folder, types, _product_path, _dump_yaml)
      generate_objects(output_folder, types)
    end

    # Create a directory of examples per resource
    def generate_resource(data)
      examples = data.object.examples
                     .reject(&:skip_test)
                     .reject { |e| !e.test_env_vars.nil? && e.test_env_vars.any? }
                     .reject { |e| @version < @api.version_obj_or_closest(e.min_version) }

      examples.each do |example|
        target_folder = data.output_folder
        target_folder = File.join(target_folder, example.name)
        FileUtils.mkpath target_folder

        data.example = example

        data.generate(
          'templates/terraform/examples/base_configs/example_file.tf.erb',
          File.join(target_folder, 'main.tf'),
          self
        )

        data.generate(
          'templates/terraform/examples/base_configs/tutorial.md.erb',
          File.join(target_folder, 'tutorial.md'),
          self
        )

        data.generate(
          'templates/terraform/examples/base_configs/example_backing_file.tf.erb',
          File.join(target_folder, 'backing_file.tf'),
          self
        )

        data.generate(
          'templates/terraform/examples/static/motd',
          File.join(target_folder, 'motd'),
          self
        )
      end
    end

    # We don't want to generate anything but the resource.
    def generate_resource_tests(data) end
  end
end
