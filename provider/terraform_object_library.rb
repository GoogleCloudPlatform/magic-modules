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

require 'provider/terraform_example'

module Provider
  # Code generator for a library converting terraform state to gcp objects.
  class TerraformObjectLibrary < Provider::Terraform
    def generate(output_folder, types, version_name, _product_path, _dump_yaml)
      version = @api.version_obj_or_default(version_name)
      generate_objects(output_folder, types, version)
    end

    def generate_resource(data)
      target_folder = data[:output_folder]
      product_ns = data[:object].__product.name

      generate_resource_file data.clone.merge(
        object: data[:object],
        default_template: 'templates/terraform/objectlib/base.go.erb',
        out_file: File.join(target_folder,
                            "google/#{product_ns.downcase}_#{data[:object].name.underscore}.go")
      )
    end

    def compile_common_files(output_folder, version_name = 'ga')
      Google::LOGGER.info 'Compiling common files.'
      compile_file_list(output_folder, [
                          ['google/config.go', 'third_party/terraform/utils/config.go.erb']
                        ], version: version_name)
    end

    def copy_common_files(output_folder, _version_name)
      Google::LOGGER.info 'Copying common files.'
      copy_file_list(output_folder, [
                       ['google/constants.go',
                        'third_party/validator/constants.go'],
                       ['google/cai.go',
                        'third_party/validator/cai.go'],
                       ['google/common_operation.go',
                        'third_party/validator/common_operation.go'],
                       ['google/compute_instance_helpers.go',
                        'third_party/validator/compute_instance_helpers.go'],
                       ['google/compute_operation.go',
                        'third_party/validator/compute_operation.go'],
                       ['google/compute_shared_operation.go',
                        'third_party/validator/compute_shared_operation.go'],
                       ['google/convert.go',
                        'third_party/validator/convert.go'],
                       ['google/metadata.go',
                        'third_party/validator/metadata.go'],
                       ['google/service_scope.go',
                        'third_party/validator/service_scope.go'],
                       ['google/compute_instance.go',
                        'third_party/validator/compute_instance.go'],
                       ['google/image.go',
                        'third_party/terraform/utils/image.go'],
                       ['google/disk_type.go',
                        'third_party/terraform/utils/disk_type.go'],
                       ['google/validation.go',
                        'third_party/terraform/utils/validation.go'],
                       ['google/regional_utils.go',
                        'third_party/terraform/utils/regional_utils.go'],
                       ['google/field_helpers.go',
                        'third_party/terraform/utils/field_helpers.go'],
                       ['google/self_link_helpers.go',
                        'third_party/terraform/utils/self_link_helpers.go'],
                       ['google/utils.go',
                        'third_party/validator/utils.go'],
                       ['google/transport.go',
                        'third_party/terraform/utils/transport.go'],
                       ['google/bigtable_client_factory.go',
                        'third_party/terraform/utils/bigtable_client_factory.go']
                     ])
    end

    def generate_resource_tests(data) end
  end
end
