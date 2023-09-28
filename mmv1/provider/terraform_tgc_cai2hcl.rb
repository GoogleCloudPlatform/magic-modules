# Copyright 2023 Google Inc.
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
  # Code generator for a library converting GCP CAI objects to Terraform state.
  class CaiToTerraformConversion < Provider::Terraform
    def generating_hashicorp_repo?
      # This code is not used when generating TPG/TPGB
      false
    end

    def generate(output_folder, types, _product_path, _dump_yaml, generate_code, generate_docs)
      Google::LOGGER.info('Generating cai2hcl converters')
      Google::LOGGER.info('NOTE: Cai2hcl converters are IN DEVELOPMENT and subject to change.')

      @base_url = @version.cai_base_url || @version.base_url
      generate_objects(
        output_folder,
        types,
        generate_code,
        generate_docs
      )
    end

    def generate_resource(pwd, data, _generate_code, _generate_docs)
      product_name = data.object.__product.name.downcase

      generators_folder = File.join(data.output_folder, 'services', product_name)
      FileUtils.mkdir_p(generators_folder)

      object_name = data.object.name.underscore

      converter_file_name = "#{product_name}_#{object_name}.go"
      data.generate(pwd,
                    'templates/cai2hcl/resource_converter.go.erb',
                    File.join(generators_folder, converter_file_name),
                    self)

      converter_test_file_name = "#{product_name}_#{object_name}_test.go"
      data.generate(pwd,
                    'templates/cai2hcl/resource_converter_test.go.erb',
                    File.join(generators_folder, converter_test_file_name),
                    self)
    end

    def compile_common_files(output_folder, products, _common_compile_file) end

    def copy_common_files(output_folder, generate_code, _generate_docs)
      return unless generate_code

      Google::LOGGER.info('Coping cai2hcl common files')

      FileUtils.mkdir_p(output_folder)

      FileUtils.cp_r('third_party/cai2hcl/.', output_folder)
    end

    def generate_resource_tests(pwd, data)
      # Generated at "generate_resource" stage.
    end

    def generate_iam_policy(pwd, data, generate_code, _generate_docs) end

    def generate_resource_sweepers(pwd, data) end
  end
end
