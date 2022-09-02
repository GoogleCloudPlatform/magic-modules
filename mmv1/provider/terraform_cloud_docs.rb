# Copyright 2021 Google Inc.
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
  # Code generator for Terraform samples meant to be displayed in cloud.google.com
  class TerraformCloudDocs < Provider::Terraform
    # We don't want *any* static generation, so we override generate to only
    # generate objects.
    def generate(output_folder, types, _product_path, _dump_yaml, generate_code, generate_docs)
      generate_objects(
        output_folder,
        types,
        generate_code,
        generate_docs
      )
    end

    # Create a directory of examples per resource
    def generate_resource(pwd, data, _generate_code, generate_docs) end

    # We don't want to generate anything but the resource.
    def generate_resource_tests(pwd, data) end

    def generate_resource_sweepers(pwd, data) end

    def compile_common_files(output_folder, products, common_compile_file) end

    def copy_common_files(output_folder, generate_code, generate_docs) end

    def generate_iam_policy(pwd, data, generate_code, generate_docs) end
  end
end
