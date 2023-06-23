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

require 'provider/terraform_oics'
require 'fileutils'

module Provider
  class TerraformCai2hclProvider < Provider::Terraform
    def generate(output_folder, types, _product_path, _dump_yaml, generate_code, generate_docs)
      resources_folder = File.join(output_folder, 'converters/google/resources')
      FileUtils.mkdir_p(resources_folder)

      @base_url = @version.cai_base_url || @version.base_url
      generate_objects(
        output_folder,
        types,
        generate_code,
        generate_docs
      )
    end

    def generate_resource(pwd, data, _generate_code, _generate_docs)
      target_folder = data.output_folder
      product_name = data.object.__product.name.downcase
      object_name = data.object.name.underscore
      data.generate(pwd,
                    'templates/cai2hcl/resource_converter.go.erb',
                    File.join(target_folder,
                              "converters/google/resources/#{product_name}_#{object_name}.go"),
                    self)
    end

    def copy_common_files(output_folder, generate_code, _generate_docs)
      Google::LOGGER.info 'Copying common files. -X'

      copy_file_list(output_folder, [
        ['converters/google/resources/cai2hcl.go', 'templates/cai2hcl/cai2hcl.go'],
        ['converters/google/resources/common.go', 'templates/cai2hcl/common.go'],
        ['converters/google/resources/helper.go', 'templates/cai2hcl/helper.go'],
        ['converters/google/resources/third_party.go', 'templates/cai2hcl/third_party.go'],
      ])
    end
  end
end
