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

require 'compile/core'
require 'dependencies/dependency_graph'
require 'fileutils'
require 'google/logger'
require 'pathname'
require 'provider/properties'
require 'provider/end2end/core'
require 'provider/shared'
require 'provider/test_matrix'
require 'provider/test_data/spec_formatter'
require 'provider/test_data/constants'
require 'provider/test_data/property'
require 'provider/test_data/create_data'
require 'provider/test_data/expectations'

module Provider
  DEFAULT_FORMAT_OPTIONS = {
    indent: 0,
    start_indent: 0,
    max_columns: 80,
    quiet: false
  }.freeze

  # Basic functionality for code generator providers. Provides basic services,
  # such as compiling and including files, formatting data, etc.
  class Core
    include Compile::Core
    include Provider::Properties
    include Provider::End2End::Core
    include Provider::Shared

    attr_reader :test_data

    def initialize(config, api)
      @config = config
      @api = api
      @property = Provider::TestData::Property.new(self)
      @constants = Provider::TestData::Constants.new(self)
      @data_gen = Provider::TestData::Generator.new
      @create_data = Provider::TestData::CreateData.new(self, @data_gen)
      @prop_data = Provider::TestData::Expectations.new(self, @data_gen)
      @generated = []
      @sourced = []
      @max_columns = 80
    end

    # Main entry point for the compiler. As this method is simply invoking other
    # generators, it is okay to ignore Rubocop warnings about method size and
    # complexity.
    #
    # rubocop:disable Metrics/AbcSize
    # rubocop:disable Metrics/CyclomaticComplexity
    # rubocop:disable Metrics/PerceivedComplexity
    def generate(output_folder, types, version_name)
      version = @api.version_obj_or_default(version_name)
      generate_objects(output_folder, types, version)
      generate_client_functions(output_folder) unless @config.functions.nil?
      copy_files(output_folder) \
        unless @config.files.nil? || @config.files.copy.nil?
      compile_examples(output_folder) unless @config.examples.nil?
      compile_end2end_tests(output_folder) unless @config.examples.nil?
      compile_network_data(output_folder) \
        unless @config.test_data.nil? || @config.test_data.network.nil?
      compile_changelog(output_folder) unless @config.changelog.nil?
      # Compilation has to be the last step, as some files (e.g.
      # CONTRIBUTING.md) may depend on the list of all files previously copied
      # or compiled.
      compile_files(output_folder) \
        unless @config.files.nil? || @config.files.compile.nil?
      apply_file_acls(output_folder) \
        unless @config.files.nil? || @config.files.permissions.nil?
      verify_test_matrixes
    end
    # rubocop:enable Metrics/AbcSize
    # rubocop:enable Metrics/CyclomaticComplexity
    # rubocop:enable Metrics/PerceivedComplexity

    def copy_files(output_folder)
      @config.files.copy.each do |target, source|
        target_file = File.join(output_folder, target)
        target_dir = File.dirname(target_file)
        @sourced << relative_path(target_file, output_folder)
        Google::LOGGER.info "Copying #{source} => #{target}"
        FileUtils.mkpath target_dir unless Dir.exist?(target_dir)
        FileUtils.cp source, target_file
      end
    end

    def compile_files(output_folder)
      compile_file_list(output_folder, @config.files.compile)
    end

    def compile_examples(output_folder)
      compile_file_map(
        output_folder,
        @config.examples,
        lambda do |_object, file|
          ["examples/#{file}",
           "products/#{@api.prefix[1..-1]}/files/examples~#{file}"]
        end
      )
    end

    def compile_network_data(output_folder)
      compile_file_map(
        output_folder,
        @config.test_data.network,
        lambda do |object, file|
          type = Google::StringUtils.underscore(object.name)
          ["spec/data/network/#{object.out_name}/#{file}.yaml",
           "products/#{@api.prefix[1..-1]}/files/spec~#{type}~#{file}.yaml"]
        end
      )
    end

    # Generate the CHANGELOG.md file with the history of the module.
    def compile_changelog(output_folder)
      FileUtils.mkpath output_folder
      generate_file(
        changes: @config.changelog,
        template: 'templates/CHANGELOG.md.erb',
        output_folder: output_folder,
        out_file: File.join(output_folder, 'CHANGELOG.md')
      )
    end

    def apply_file_acls(output_folder)
      @config.files.permissions.each do |perm|
        Google::LOGGER.info "Permission #{perm.path} => #{perm.acl}"
        FileUtils.chmod perm.acl, File.join(output_folder, perm.path)
      end
    end

    def compile_file_map(output_folder, section, mapper)
      create_object_list(section, mapper).each do |o|
        compile_file_list(
          output_folder,
          o
        )
      end
    end

    # Creates an object list by calling a lambda
    # This can be useful for converting a list of config values to something
    # less human-centric.
    def create_object_list(section, mapper)
      @api.objects
          .select { |o| section.key?(o.name) }
          .map do |o|
            Hash[section[o.name].map { |file| mapper.call(o, file) }]
          end
    end

    def list_manual_network_data
      test_data = @config&.test_data&.network || {}
      create_object_list(
        test_data,
        lambda do |object, file|
          type = Google::StringUtils.underscore(object.name)
          ["spec/data/network/#{object.out_name}/#{file}.yaml",
           "products/#{@api.prefix[1..-1]}/files/spec~#{type}~#{file}.yaml"]
        end
      )
    end

    # rubocop:disable Metrics/MethodLength
    # rubocop:disable Metrics/AbcSize
    def compile_file_list(output_folder, files, data = {})
      files.each do |target, source|
        Google::LOGGER.info "Compiling #{source} => #{target}"
        target_file = File.join(output_folder, target)
                          .gsub('{{product_name}}', @api.prefix[1..-1])

        manifest = @config.respond_to?(:manifest) ? @config.manifest : {}
        generate_file(
          data.clone.merge(
            name: target,
            product: @api,
            object: {},
            config: {},
            scopes: @api.scopes,
            manifest: manifest,
            tests: '',
            template: source,
            generated_files: @generated,
            sourced_files: @sourced,
            compiler: compiler,
            output_folder: output_folder,
            out_file: target_file,
            prop_ns_dir: @api.prefix[1..-1].downcase,
            product_ns: Google::StringUtils.camelize(@api.prefix[1..-1], :upper)
          )
        )

        %x(goimports -w #{target_file}) if File.extname(target_file) == '.go'
      end
    end
    # rubocop:enable Metrics/MethodLength
    # rubocop:enable Metrics/AbcSize

    # rubocop:disable Metrics/CyclomaticComplexity
    # rubocop:disable Metrics/PerceivedComplexity
    def generate_objects(output_folder, types, version)
      @api.set_properties_based_on_version(version)
      @api.objects.each do |object|
        if !types.empty? && !types.include?(object.name)
          Google::LOGGER.info "Excluding #{object.name} per user request"
        elsif types.empty? && object.exclude
          Google::LOGGER.info "Excluding #{object.name} per API catalog"
        elsif types.empty? && object.exclude_if_not_in_version(version)
          Google::LOGGER.info "Excluding #{object.name} per API version"
        else
          generate_object object, output_folder, version
        end
      end
    end
    # rubocop:enable Metrics/CyclomaticComplexity
    # rubocop:enable Metrics/PerceivedComplexity

    def generate_object(object, output_folder, version)
      data = build_object_data(object, output_folder, version)

      generate_resource data
      generate_resource_tests data
      generate_properties data, object.all_user_properties
      generate_network_datas data, object
    end

    # Generates all 6 network data files for a object.
    # This includes all combinations of seeds [0-2] and title == / != name
    # Each data file is a YAML file with all properties possible on an object.
    #
    # @config.test_data lists all files that are written by hand and will not
    # be generated.
    #
    # Requires:
    #  object: The Api::Resource used as basis for the network data.
    #  data: A hash with values:
    #    output_folder: root folder for generated module
    def generate_network_datas(data, object)
      target_folder = File.join(data[:output_folder],
                                'spec', 'data', 'network', object.out_name)
      FileUtils.mkpath target_folder

      # Create list of compiled network data
      manual = list_manual_network_data
      3.times.each do |id|
        %w[name title].each do |name|
          out_file = File.join(target_folder, "success#{id + 1}~#{name}.yaml")
          next if manual.include? out_file
          next if true?(data[:object].manual)

          generate_network_data data.clone.merge(
            out_file: File.join(target_folder, "success#{id + 1}~#{name}.yaml"),
            id: id,
            title: name,
            object: object
          )
        end
      end
    end
    # rubocop:enable Metrics/MethodLength

    # Generates a single network data file for unit testing.
    # Required values in data:
    #   out_file: path of data file to create
    #   id: a seed value
    #   title: The name of object who is unit tested with this spec file
    #   object: The Api::Resource used as basis for the network data
    def generate_network_data(data)
      formatter = Provider::TestData::SpecFormatter.new(self)

      name = "title#{data[:id]}" if data[:title] == 'title'
      name = "test name##{data[:id]} data" if data[:title] == 'name'
      generate_file data.clone.merge(
        template: 'templates/network_spec.yaml.erb',
        test_data: formatter.generate(data[:object], '', data[:object].kind,
                                      data[:id],
                                      name: name)
      )
    end

    def build_object_data(object, output_folder, version)
      {
        name: object.out_name,
        object: object,
        config: (@config.objects || {}).select { |o, _v| o == object.name }
                                       .fetch(object.name, {}),
        tests: (@config.tests || {}).select { |o, _v| o == object.name }
                                    .fetch(object.name, {}),
        output_folder: output_folder,
        product_name: object.__product.prefix[1..-1],
        version: version
      }
    end

    def generate_resource_file(data)
      product_ns = if @config.name.nil?
                     Google::StringUtils.camelize(data[:object].__product
                       .prefix[1..-1], :upper)
                   else
                     @config.name
                   end
      generate_file(data.clone.merge(
        # Override with provider specific template for this object, if needed
        template: Google::HashUtils.navigate(data[:config], ['template',
                                                             data[:type]],
                                             data[:default_template]),
        product_ns: product_ns
      ))
    end

    def generate_file(data)
      file_folder = File.dirname(data[:out_file])
      file_relative = relative_path(data[:out_file], data[:output_folder]).to_s
      FileUtils.mkpath file_folder unless Dir.exist?(file_folder)
      @generated << relative_path(data[:out_file], data[:output_folder])
      ctx = binding
      data.each { |name, value| ctx.local_variable_set(name, value) }
      generate_file_write ctx, data
    end

    def generate_file_write(ctx, data)
      enforce_file_expectations data[:out_file] do
        Google::LOGGER.info "Generating #{data[:name]} #{data[:type]}"
        write_file data[:out_file], compile_file(ctx, data[:template])
      end
    end
  end
end
