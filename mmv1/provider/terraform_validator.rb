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
  # Code generator for a library converting terraform state to gcp objects.
  class TerraformValidator < Provider::Terraform
    def generating_hashicorp_repo?
      # This code is not used when generating TPG/TPGB
      false
    end

    def generate(output_folder, types, _product_path, _dump_yaml, generate_code, generate_docs)
      # Temporary shim to generate the missing resources directory. Can be removed
      # once the folder exists downstream.
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

    def generate_object(object, output_folder, version_name, generate_code, generate_docs)
      if object.exclude_validator
        Google::LOGGER.info "Skipping fine-grained resource #{object.name}"
        return
      end

      super(object, output_folder, version_name, generate_code, generate_docs)
    end

    def generate_resource(pwd, data, _generate_code, _generate_docs)
      product_name = data.object.__product.name.downcase
      output_folder = File.join(
        data.output_folder,
        'converters/google/resources/services',
        product_name
      )
      object_name = data.object.name.underscore
      target = "#{product_name}_#{object_name}.go"
      data.generate(pwd,
                    'templates/validator/resource_converter.go.erb',
                    File.join(output_folder, target),
                    self)
      replace_import_path(output_folder, target)
    end

    def retrieve_list_of_manually_defined_tests_from_file(file)
      content = File.read(file)
      content.scan(/\s*name\s*:\s*"([^,]+)"/).flatten(1)
    end

    def retrieve_list_of_manually_defined_tests
      m1 =
        retrieve_list_of_manually_defined_tests_from_file(
          'third_party/validator/tests/source/cli_test.go.erb'
        )
      m2 =
        retrieve_list_of_manually_defined_tests_from_file(
          'third_party/validator/tests/source/read_test.go.erb'
        )
      m1 | m2 # union of manually defined tests
    end

    def validate_non_defined_tests(file_set, non_defined_tests)
      if non_defined_tests.any? { |test| !file_set.member?("#{test}.json") }
        raise "test file named #{test}.json expected but found none"
      end

      if non_defined_tests.any? { |test| !file_set.member?("#{test}.tfplan.json") }
        raise "test file named #{test}.tfplan.json expected but found none"
      end

      return unless non_defined_tests.any? { |test| !file_set.member?("#{test}.tf") }

      raise "test file named #{test}.tf expected but found none"
    end

    def retrieve_full_list_of_test_files
      files = Dir['third_party/validator/tests/data/*']
      files = files.map { |file| file.split('/')[-1] }
      files.sort
    end

    def retrieve_full_list_of_test_files_with_location
      files = retrieve_full_list_of_test_files
      files.map do |file|
        ["testdata/templates/#{file}", "third_party/validator/tests/data/#{file}"]
      end
    end

    def retrieve_full_manifest_of_non_defined_tests
      files = retrieve_full_list_of_test_files
      tests = files.map { |file| file.split('.')[0] } | []
      non_defined_tests = tests - retrieve_list_of_manually_defined_tests
      non_defined_tests = non_defined_tests.reject do |file|
        file.end_with?('_without_default_project')
      end
      validate_non_defined_tests(files.to_set, non_defined_tests)
      non_defined_tests
    end

    def retrieve_test_source_files(path, suffix)
      files = Dir["#{path}**#{suffix}"]
      files = files.map { |file| file.split(path)[-1] }
      files.sort
    end

    def retrieve_test_source_code_with_location(suffix)
      path = 'third_party/validator/tests/source/'
      files = retrieve_test_source_files(path, suffix)
      files.map do |file|
        ["test/#{file}", path + file]
      end
    end

    def compile_common_files(output_folder, products, _common_compile_file)
      Google::LOGGER.info 'Compiling common files.'
      file_template = ProviderFileTemplate.new(
        output_folder,
        @target_version_name,
        build_env,
        products
      )

      @non_defined_tests = retrieve_full_manifest_of_non_defined_tests
      files = retrieve_full_list_of_test_files
      @tests = files.map { |file| file.split('.')[0] } | []

      test_source = retrieve_test_source_code_with_location('[b]').map do |location|
        [location[0].sub('go.erb', 'go'), location[1]]
      end

      compile_file_list(
        output_folder,
        test_source,
        file_template
      )

      compile_file_list(
        output_folder,
        [
          ['converters/google/resources/compute_operation.go',
           'third_party/terraform/utils/compute_operation.go.erb'],
          ['converters/google/resources/services/compute/compute_operation.go',
           'third_party/terraform/services/compute/compute_operation.go.erb'],
          ['converters/google/resources/transport/config.go',
           'third_party/terraform/transport/config.go.erb'],
          ['converters/google/resources/config_test_utils.go',
           'third_party/terraform/utils/config_test_utils.go.erb'],
          ['converters/google/resources/iam.go',
           'third_party/terraform/utils/iam.go.erb'],
          ['converters/google/resources/tpgiamresource/iam.go',
           'third_party/terraform/tpgiamresource/iam.go.erb'],
          ['converters/google/resources/tpgresource/common_diff_suppress.go',
           'third_party/terraform/tpgresource/common_diff_suppress.go.erb'],
          ['converters/google/resources/services/compute/compute_instance_helpers.go',
           'third_party/terraform/services/compute/compute_instance_helpers.go.erb'],
          ['converters/google/resources/transport/provider_handwritten_endpoint.go',
           'third_party/terraform/transport/provider_handwritten_endpoint.go.erb'],
          ['converters/google/resources/resource_converters.go',
           'templates/validator/resource_converters.go.erb'],
          ['converters/google/resources/mappers.go',
           'templates/validator/mappers/mappers.go.erb'],
          ['converters/google/resources/services/kms/iam_kms_key_ring.go',
           'third_party/terraform/services/kms/iam_kms_key_ring.go.erb'],
          ['converters/google/resources/services/kms/iam_kms_crypto_key.go',
           'third_party/terraform/services/kms/iam_kms_crypto_key.go.erb'],
          ['converters/google/resources/services/compute/metadata.go',
           'third_party/terraform/services/compute/metadata.go.erb'],
          ['converters/google/resources/services/compute/compute_instance.go',
           'third_party/validator/compute_instance.go.erb']
        ],
        file_template
      )
    end

    def copy_common_files(output_folder, generate_code, _generate_docs)
      Google::LOGGER.info 'Copying common files.'
      return unless generate_code

      copy_file_list(
        output_folder,
        retrieve_full_list_of_test_files_with_location
      )

      copy_file_list(
        output_folder,
        retrieve_test_source_code_with_location('[^b]')
      )

      copy_file_list(output_folder, [
                       ['converters/google/resources/tpgresource/constants.go',
                        'third_party/validator/tpgresource/constants.go'],
                       ['converters/google/resources/constants.go',
                        'third_party/validator/constants.go'],
                       ['converters/google/resources/cai.go',
                        'third_party/validator/cai.go'],
                       ['converters/google/resources/tpgresource/cai.go',
                        'third_party/validator/tpgresource/cai.go'],
                       ['converters/google/resources/tpgresource/cai_test.go',
                        'third_party/validator/tpgresource/cai_test.go'],
                       ['converters/google/resources/org_policy_policy.go',
                        'third_party/validator/org_policy_policy.go'],
                       ['converters/google/resources/getconfig.go',
                        'third_party/validator/getconfig.go'],
                       ['converters/google/resources/folder.go',
                        'third_party/validator/folder.go'],
                       ['converters/google/resources/getconfig_test.go',
                        'third_party/validator/getconfig_test.go'],
                       ['converters/google/resources/transport/config_test_utils.go',
                        'third_party/terraform/transport/config_test_utils.go'],
                       ['converters/google/resources/tpgresource/json_map.go',
                        'third_party/validator/tpgresource/json_map.go'],
                       ['converters/google/resources/json_map.go',
                        'third_party/validator/json_map.go'],
                       ['converters/google/resources/project.go',
                        'third_party/validator/project.go'],
                       ['converters/google/resources/sql_database_instance.go',
                        'third_party/validator/sql_database_instance.go'],
                       ['converters/google/resources/storage_bucket.go',
                        'third_party/validator/storage_bucket.go'],
                       ['converters/google/resources/cloudfunctions_function.go',
                        'third_party/validator/cloudfunctions_function.go'],
                       ['converters/google/resources/transport/dcl.go',
                        'third_party/validator/dcl.go'],
                       ['converters/google/resources/cloudfunctions_cloud_function.go',
                        'third_party/validator/cloudfunctions_cloud_function.go'],
                       ['converters/google/resources/bigquery_table.go',
                        'third_party/validator/bigquery_table.go'],
                       ['converters/google/resources/bigtable_cluster.go',
                        'third_party/validator/bigtable_cluster.go'],
                       ['converters/google/resources/bigtable_instance.go',
                        'third_party/validator/bigtable_instance.go'],
                       ['converters/google/resources/tpgiamresource/iam_helpers.go',
                        'third_party/validator/tpgiamresource/iam_helpers.go'],
                       ['converters/google/resources/iam_helpers.go',
                        'third_party/validator/iam_helpers.go'],
                       ['converters/google/resources/tpgiamresource/iam_helpers_test.go',
                        'third_party/validator/tpgiamresource/iam_helpers_test.go'],
                       ['converters/google/resources/organization_iam.go',
                        'third_party/validator/organization_iam.go'],
                       ['converters/google/resources/project_iam.go',
                        'third_party/validator/project_iam.go'],
                       ['converters/google/resources/project_organization_policy.go',
                        'third_party/validator/project_organization_policy.go'],
                       ['converters/google/resources/folder_organization_policy.go',
                        'third_party/validator/folder_organization_policy.go'],
                       ['converters/google/resources/folder_iam.go',
                        'third_party/validator/folder_iam.go'],
                       ['converters/google/resources/container.go',
                        'third_party/validator/container.go'],
                       ['converters/google/resources/project_service.go',
                        'third_party/validator/project_service.go'],
                       ['converters/google/resources/services/monitoring/monitoring_slo_helper.go',
                        'third_party/validator/monitoring_slo_helper.go'],
                       ['converters/google/resources/service_account.go',
                        'third_party/validator/service_account.go'],
                       ['converters/google/resources/services/compute/image.go',
                        'third_party/terraform/services/compute/image.go'],
                       ['converters/google/resources/import.go',
                        'third_party/terraform/utils/import.go'],
                       ['converters/google/resources/tpgresource/import.go',
                        'third_party/terraform/tpgresource/import.go'],
                       ['converters/google/resources/services/compute/disk_type.go',
                        'third_party/terraform/services/compute/disk_type.go'],
                       ['converters/google/resources/verify/validation.go',
                        'third_party/terraform/verify/validation.go'],
                       ['converters/google/resources/validation.go',
                        'third_party/terraform/utils/validation.go'],
                       ['converters/google/resources/tpgresource/regional_utils.go',
                        'third_party/terraform/tpgresource/regional_utils.go'],
                       ['converters/google/resources/regional_utils.go',
                        'third_party/terraform/utils/regional_utils.go'],
                       ['converters/google/resources/tpgresource/field_helpers.go',
                        'third_party/terraform/tpgresource/field_helpers.go'],
                       ['converters/google/resources/field_helpers.go',
                        'third_party/terraform/utils/field_helpers.go'],
                       ['converters/google/resources/self_link_helpers.go',
                        'third_party/terraform/utils/self_link_helpers.go'],
                       ['converters/google/resources/tpgresource/self_link_helpers.go',
                        'third_party/terraform/tpgresource/self_link_helpers.go'],
                       ['converters/google/resources/transport/header_transport.go',
                        'third_party/terraform/transport/header_transport.go'],
                       ['converters/google/resources/transport/bigtable_client_factory.go',
                        'third_party/terraform/transport/bigtable_client_factory.go'],
                       ['converters/google/resources/tpgresource/common_operation.go',
                        'third_party/terraform/tpgresource/common_operation.go'],
                       ['converters/google/resources/common_operation.go',
                        'third_party/terraform/utils/common_operation.go'],
                       ['converters/google/resources/convert.go',
                        'third_party/terraform/utils/convert.go'],
                       ['converters/google/resources/tpgresource/convert.go',
                        'third_party/terraform/tpgresource/convert.go'],
                       ['converters/google/resources/services/logging/extract.go',
                        'third_party/terraform/services/logging/extract.go'],
                       ['converters/google/resources/service_scope.go',
                        'third_party/terraform/utils/service_scope.go'],
                       ['converters/google/resources/tpgresource/service_scope.go',
                        'third_party/terraform/tpgresource/service_scope.go'],
                       ['converters/google/resources/services/kms/kms_utils.go',
                        'third_party/terraform/services/kms/kms_utils.go'],
                       ['converters/google/resources/transport/batcher.go',
                        'third_party/terraform/transport/batcher.go'],
                       ['converters/google/resources/transport/retry_utils.go',
                        'third_party/terraform/transport/retry_utils.go'],
                       ['converters/google/resources/retry_utils.go',
                        'third_party/terraform/utils/retry_utils.go'],
                       ['converters/google/resources/services/sourcerepo/source_repo_utils.go',
                        'third_party/terraform/services/sourcerepo/source_repo_utils.go'],
                       ['converters/google/resources/transport/retry_transport.go',
                        'third_party/terraform/transport/retry_transport.go'],
                       ['converters/google/resources/transport/transport.go',
                        'third_party/terraform/transport/transport.go'],
                       ['converters/google/resources/transport/error_retry_predicates.go',
                        'third_party/terraform/transport/error_retry_predicates.go'],
                       ['converters/google/resources/error_retry_predicates.go',
                        'third_party/terraform/utils/error_retry_predicates.go'],
                       ['converters/google/resources/services/pubsub/pubsub_utils.go',
                        'third_party/terraform/services/pubsub/pubsub_utils.go'],
                       ['converters/google/resources/sqladmin_operation.go',
                        'third_party/terraform/utils/sqladmin_operation.go'],
                       ['converters/google/resources/services/sql/sqladmin_operation.go',
                        'third_party/terraform/services/sql/sqladmin_operation.go'],
                       ['converters/google/resources/verify/path_or_contents.go',
                        'third_party/terraform/verify/path_or_contents.go'],
                       ['converters/google/resources/path_or_contents.go',
                        'third_party/terraform/utils/path_or_contents.go'],
                       ['converters/google/resources/transport/mutexkv.go',
                        'third_party/terraform/transport/mutexkv.go'],
                       ['converters/google/resources/hashcode.go',
                        'third_party/terraform/utils/hashcode.go'],
                       ['converters/google/resources/tpgresource/hashcode.go',
                        'third_party/terraform/tpgresource/hashcode.go'],
                       ['converters/google/resources/iam_organization.go',
                        'third_party/terraform/utils/iam_organization.go'],
                       ['converters/google/resources/iam_folder.go',
                        'third_party/terraform/utils/iam_folder.go'],
                       ['converters/google/resources/iam_project.go',
                        'third_party/terraform/utils/iam_project.go'],
                       ['converters/google/resources/services/privateca/privateca_utils.go',
                        'third_party/terraform/services/privateca/privateca_utils.go'],
                       ['converters/google/resources/utils.go',
                        'third_party/terraform/utils/utils.go'],
                       ['converters/google/resources/tpgresource/utils.go',
                        'third_party/terraform/tpgresource/utils.go'],
                       ['converters/google/resources/iam_bigquery_dataset.go',
                        'third_party/terraform/utils/iam_bigquery_dataset.go'],
                       ['converters/google/resources/services/bigquery/iam_bigquery_dataset.go',
                        'third_party/terraform/services/bigquery/iam_bigquery_dataset.go'],
                       ['converters/google/resources/bigquery_dataset_iam.go',
                        'third_party/validator/bigquery_dataset_iam.go'],
                       ['converters/google/resources/compute_security_policy.go',
                        'third_party/validator/compute_security_policy.go'],
                       ['converters/google/resources/kms_key_ring_iam.go',
                        'third_party/validator/kms_key_ring_iam.go'],
                       ['converters/google/resources/kms_crypto_key_iam.go',
                        'third_party/validator/kms_crypto_key_iam.go'],
                       ['converters/google/resources/project_iam_custom_role.go',
                        'third_party/validator/project_iam_custom_role.go'],
                       ['converters/google/resources/organization_iam_custom_role.go',
                        'third_party/validator/organization_iam_custom_role.go'],
                       ['converters/google/resources/iam_pubsub_subscription.go',
                        'third_party/terraform/utils/iam_pubsub_subscription.go'],
                       ['converters/google/resources/pubsub_subscription_iam.go',
                        'third_party/validator/pubsub_subscription_iam.go'],
                       ['converters/google/resources/iam_spanner_database.go',
                        'third_party/terraform/utils/iam_spanner_database.go'],
                       ['converters/google/resources/spanner_database_iam.go',
                        'third_party/validator/spanner_database_iam.go'],
                       ['converters/google/resources/iam_spanner_instance.go',
                        'third_party/terraform/utils/iam_spanner_instance.go'],
                       ['converters/google/resources/spanner_instance_iam.go',
                        'third_party/validator/spanner_instance_iam.go'],
                       ['converters/google/resources/storage_bucket_iam.go',
                        'third_party/validator/storage_bucket_iam.go'],
                       ['converters/google/resources/organization_policy.go',
                        'third_party/validator/organization_policy.go'],
                       ['converters/google/resources/iam_storage_bucket.go',
                        'third_party/validator/iam_storage_bucket.go']
                     ])
    end

    def generate_resource_tests(pwd, data)
      product_whitelist = []

      return unless product_whitelist.include?(data.product.name.downcase)
      return if data.object.examples
                    .reject(&:skip_test)
                    .reject do |e|
                  @api.version_obj_or_closest(data.version) \
                < @api.version_obj_or_closest(e.min_version)
                end
                    .empty?

      FileUtils.mkpath folder_name(data.version)
      data.generate(
        pwd,
        'templates/validator/examples/base_configs/test_file.go.erb',
        "test/resource_#{full_resource_name(data)}_generated_test.go",
        self
      )
    end

    # Generate the IAM policy for this object. This is used to query and test
    # IAM policies separately from the resource itself
    # Docs are generated for the terraform provider, not here.
    def generate_iam_policy(pwd, data, generate_code, _generate_docs)
      return unless generate_code
      return if data.object.iam_policy.exclude_validator

      name = data.object.filename_override || data.object.name.underscore
      product_name = data.product.name.downcase
      product_name_underscore = product_name.underscore
      output_folder = File.join(
        data.output_folder,
        'converters/google/resources/services',
        product_name
      )

      FileUtils.mkpath output_folder
      target = "#{product_name_underscore}_#{name}_iam.go"
      data.generate(pwd,
                    'templates/validator/resource_converter_iam.go.erb',
                    File.join(output_folder, target),
                    self)
      replace_import_path(output_folder, target)

      target = "iam_#{product_name_underscore}_#{name}.go"
      data.generate(pwd,
                    'templates/terraform/iam_policy.go.erb',
                    File.join(output_folder, target),
                    self)
      replace_import_path(output_folder, target)

      # Don't generate tests - we can rely on the terraform provider
      # to test these.
    end

    def generate_resource_sweepers(pwd, data) end

    def replace_import_path(output_folder, target)
      # Replace import paths to reference the resources dir instead of the google provider
      data = File.read("#{output_folder}/#{target}")
      # rubocop:disable Layout/LineLength
      data = data.gsub(
        %r{(?<!provider ")github.com/hashicorp/terraform-provider-google/google},
        'github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources'
      )
      # rubocop:enable Layout/LineLength
      File.write("#{output_folder}/#{target}", data)
    end
  end
end
