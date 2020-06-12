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

module Provider
  # Code generator for a library converting terraform state to gcp objects.
  class TerraformObjectLibrary < Provider::Terraform
    def generate(output_folder, types, _product_path, _dump_yaml)
      @base_url = @version.base_url
      generate_objects(output_folder, types)
    end

    def generate_object(object, output_folder, version_name)
      if object.exclude_validator
        Google::LOGGER.info "Skipping fine-grained resource #{object.name}"
        return
      end

      super(object, output_folder, version_name)
    end

    def generate_resource(pwd, data)
      target_folder = data.output_folder
      product_ns = data.object.__product.name

      data.generate(pwd,
                    'templates/terraform/objectlib/base.go.erb',
                    File.join(target_folder,
                              "google/#{product_ns.downcase}_#{data.object.name.underscore}.go"),
                    self)
    end

    def compile_common_files(output_folder, products, _common_compile_file)
      Google::LOGGER.info 'Compiling common files.'
      file_template = ProviderFileTemplate.new(
        output_folder,
        @target_version_name,
        build_env,
        products
      )
      compile_file_list(output_folder, [
                          ['google/config.go',
                           'third_party/terraform/utils/config.go.erb'],
                          ['google/utils.go',
                           'third_party/terraform/utils/utils.go.erb'],
                          ['google/provider_handwritten_endpoint.go',
                           'third_party/terraform/utils/provider_handwritten_endpoint.go.erb']
                        ],
                        file_template)
    end

    def copy_common_files(output_folder)
      Google::LOGGER.info 'Copying common files.'
      copy_file_list(output_folder, [
                       ['google/constants.go',
                        'third_party/validator/constants.go'],
                       ['google/cai.go',
                        'third_party/validator/cai.go'],
                       ['google/cai_test.go',
                        'third_party/validator/cai_test.go'],
                       ['google/json_map.go',
                        'third_party/validator/json_map.go'],
                       ['google/project.go',
                        'third_party/validator/project.go'],
                       ['google/compute_instance.go',
                        'third_party/validator/compute_instance.go'],
                       ['google/sql_database_instance.go',
                        'third_party/validator/sql_database_instance.go'],
                       ['google/storage_bucket.go',
                        'third_party/validator/storage_bucket.go'],
                       ['google/storage_bucket_iam.go',
                        'third_party/validator/storage_bucket_iam.go'],
                       ['google/iam_helpers.go',
                        'third_party/validator/iam_helpers.go'],
                       ['google/iam_helpers_test.go',
                        'third_party/validator/iam_helpers_test.go'],
                       ['google/organization_iam.go',
                        'third_party/validator/organization_iam.go'],
                       ['google/project_iam.go',
                        'third_party/validator/project_iam.go'],
                       ['google/folder_iam.go',
                        'third_party/validator/folder_iam.go'],
                       ['google/container.go',
                        'third_party/validator/container.go'],
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
                       ['google/transport.go',
                        'third_party/terraform/utils/transport.go'],
                       ['google/bigtable_client_factory.go',
                        'third_party/terraform/utils/bigtable_client_factory.go'],
                       ['google/common_operation.go',
                        'third_party/terraform/utils/common_operation.go'],
                       ['google/compute_operation.go',
                        'third_party/terraform/utils/compute_operation.go'],
                       ['google/compute_shared_operation.go',
                        'third_party/terraform/utils/compute_shared_operation.go'],
                       ['google/compute_instance_helpers.go',
                        'third_party/terraform/utils/compute_instance_helpers.go.erb'],
                       ['google/convert.go',
                        'third_party/terraform/utils/convert.go'],
                       ['google/metadata.go',
                        'third_party/terraform/utils/metadata.go'],
                       ['google/service_scope.go',
                        'third_party/terraform/utils/service_scope.go'],
                       ['google/kms_utils.go',
                        'third_party/terraform/utils/kms_utils.go'],
                       ['google/batcher.go',
                        'third_party/terraform/utils/batcher.go'],
                       ['google/retry_utils.go',
                        'third_party/terraform/utils/retry_utils.go'],
                       ['google/source_repo_utils.go',
                        'third_party/terraform/utils/source_repo_utils.go'],
                       ['google/retry_transport.go',
                        'third_party/terraform/utils/retry_transport.go'],
                       ['google/error_retry_predicates.go',
                        'third_party/terraform/utils/error_retry_predicates.go']
                     ])
    end

    def generate_resource_tests(pwd, data) end

    def generate_iam_policy(pwd, data) end

    def generate_resource_sweepers(pwd, data) end
  end
end
