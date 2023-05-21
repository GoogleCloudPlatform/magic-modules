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
  # Code generator for a library converting gcp objects to terraform state.
  class TerraformCai2hclProvider < Provider::Terraform
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
      target_folder = data.output_folder
      product_name = data.object.__product.name.downcase
      object_name = data.object.name.underscore
      data.generate(pwd,
                    'templates/cai2hcl/resource_converter.go.erb',
                    File.join(target_folder,
                              "converters/google/resources/#{product_name}_#{object_name}.go"),
                    self)
    end

    def retrieve_list_of_manually_defined_tests_from_file(file)
      content = File.read(file)
      content.scan(/\s*name\s*:\s*"([^,]+)"/).flatten(1)
    end

    def retrieve_list_of_manually_defined_tests
      m1 =
        retrieve_list_of_manually_defined_tests_from_file(
          'third_party/validator/tests/tgc-source/cli_test.go.erb'
        )
      m2 =
        retrieve_list_of_manually_defined_tests_from_file(
          'third_party/validator/tests/tgc-source/read_test.go.erb'
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
      # file_template = ProviderFileTemplate.new(
      #   output_folder,
      #   @target_version_name,
      #   build_env,
      #   products
      # )

      # @non_defined_tests = retrieve_full_manifest_of_non_defined_tests
      # files = retrieve_full_list_of_test_files
      # @tests = files.map { |file| file.split('.')[0] } | []

      # test_source = retrieve_test_source_code_with_location('[b]').map do |location|
      #   [location[0].sub('go.erb', 'go'), location[1]]
      # end

      # compile_file_list(
      #   output_folder,
      #   test_source,
      #   file_template
      # )

      # compile_file_list(output_folder, [
      #                     [['converters/google/resources/config.go',
      #                       'third_party/terraform/utils/config.go.erb'],]
      #                   ],
      #                   file_template)
    end

    def copy_common_files(output_folder, generate_code, _generate_docs)
      Google::LOGGER.info 'Copying common files. -X'
      # return unless generate_code

      # copy_file_list(
      #   output_folder,
      #   retrieve_full_list_of_test_files_with_location
      # )

      # copy_file_list(
      #   output_folder,
      #   retrieve_test_source_code_with_location('[^b]')
      # )

      copy_file_list(output_folder, [
        ['converters/google/resources/cai2hcl.go', 'templates/cai2hcl/cai2hcl.go'],
        ['converters/google/resources/common.go', 'templates/cai2hcl/common.go'],
        ['converters/google/resources/helper.go', 'templates/cai2hcl/helper.go'],
        ['converters/google/resources/third_party.go', 'templates/cai2hcl/third_party.go'],
      ])
    end

    # def generate_resource_tests(pwd, data)
    #   product_whitelist = []

    #   return unless product_whitelist.include?(data.product.name.downcase)
    #   return if data.object.examples
    #                 .reject(&:skip_test)
    #                 .reject do |e|
    #               @api.version_obj_or_closest(data.version) \
    #             < @api.version_obj_or_closest(e.min_version)
    #             end
    #                 .empty?

    #   FileUtils.mkpath folder_name(data.version)
    #   data.generate(
    #     pwd,
    #     'templates/validator/examples/base_configs/test_file.go.erb',
    #     "test/resource_#{full_resource_name(data)}_generated_test.go",
    #     self
    #   )
    # end
  end
end
