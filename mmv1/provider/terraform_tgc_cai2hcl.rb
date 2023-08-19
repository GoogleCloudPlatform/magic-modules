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
  # Code generator for a library converting GCP CAI objects to Terraform state.
  class CaiToTerraformConversion < Provider::Terraform
    def generating_hashicorp_repo?
      # This code is not used when generating TPG/TPGB
      false
    end

    def generate(output_folder, types, _product_path, _dump_yaml, generate_code, generate_docs)
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
      output_folder = File.join(
        data.output_folder,
        'cai2hcl/generated/converters',
        product_name
      )
      object_name = data.object.name.underscore
      target = "#{product_name}_#{object_name}.go"
      data.generate(pwd,
                    'templates/cai2hcl/resource_converter.go.erb',
                    File.join(output_folder, target),
                    self)
    end

    def compile_common_files(output_folder, products, _common_compile_file) end

    def copy_common_files(output_folder, generate_code, _generate_docs)
      Google::LOGGER.info 'Copying common files.'
      return unless generate_code

      copy_file_list(output_folder, [
                       ['cai2hcl/generated/converters/common/converter.go',
                        'third_party/cai2hcl/converters/common/converter.go'],
                       ['cai2hcl/generated/converters/common/utils.go',
                        'third_party/cai2hcl/converters/common/utils.go'],
                       ['cai2hcl/generated/convert.go',
                        'third_party/cai2hcl/convert.go'],
                       ['cai2hcl/generated/converter_map.go',
                        'third_party/cai2hcl/converter_map.go']
                     ])

      Google::LOGGER.info 'Copying testdata files.'

      copy_file_list(output_folder, [
                       ['cai2hcl/generated/converters/testdata/full_compute_backend_service.json',
                        'third_party/cai2hcl/converters/testdata/full_compute_backend_service.json'],
                       ['cai2hcl/generated/converters/testdata/full_compute_backend_service.tf',
                        'third_party/cai2hcl/converters/testdata/full_compute_backend_service.tf'],

                       ['cai2hcl/generated/converters/testdata/full_compute_forwarding_rule.json',
                        'third_party/cai2hcl/converters/testdata/full_compute_forwarding_rule.json'],
                       ['cai2hcl/generated/converters/testdata/full_compute_forwarding_rule.tf',
                        'third_party/cai2hcl/converters/testdata/full_compute_forwarding_rule.tf'],

                       ['cai2hcl/generated/converters/testdata/full_compute_health_check.json',
                        'third_party/cai2hcl/converters/testdata/full_compute_health_check.json'],
                       ['cai2hcl/generated/converters/testdata/full_compute_health_check.tf',
                        'third_party/cai2hcl/converters/testdata/full_compute_health_check.tf']
                     ])
    end

    def generate_resource_tests(pwd, data) end

    def generate_iam_policy(pwd, data, generate_code, _generate_docs) end

    def generate_resource_sweepers(pwd, data) end
  end
end
