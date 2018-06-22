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

require 'fileutils'
require 'google/ruby_utils'
require 'google/hash_utils'
require 'google/string_utils'
require 'provider/config'
require 'provider/core'
require 'provider/puppet/codegen'
require 'provider/puppet/manifest'
require 'provider/puppet/resource_override'
require 'provider/puppet/property_override'
require 'provider/puppet/test_manifest'
require 'provider/test_matrix'
require 'provider/test_data/utils'

module Provider
  # A code generator for Puppet modules
  class Puppet < Provider::Core
    STUBS_FOLDER = File.join('spec', 'stubs').freeze
    BOLT_TASK_ACL = 'u=rwx,g=rx,o=rx'.freeze
    BOLT_UNDEF_MAGIC = '<-undef->'.freeze

    include Provider::Puppet::Codegen
    include Google::RubyUtils
    include Provider::TestData::TestUtils

    # Settings for the provider
    class Config < Provider::Config
      attr_reader :manifest
      attr_reader :functions
      attr_reader :bolt_tasks

      def provider
        Provider::Puppet
      end

      def resource_override
        Provider::Puppet::ResourceOverride
      end

      def property_override
        Provider::Puppet::PropertyOverride
      end

      def validate
        super

        check_optional_property :manifest, Provider::Puppet::Manifest
        check_property_list :functions, Provider::Config::Function
        check_property_list :bolt_tasks, Provider::Puppet::BoltTask
      end
    end

    # A Bolt task
    class BoltTask < Provider::Config::Function
      attr_reader :style
      attr_reader :input
      attr_reader :manifest

      INPUTS = %i[stdin].freeze
      STYLES = %i[puppet ruby].freeze

      def description_display
        @description.split("\n").map(&:strip).join(' ')
      end

      def target_file
        File.join('tasks', case @style
                           when :ruby
                             "#{name}.rb"
                           when :puppet
                             "#{name}.sh"
                           else
                             raise "Unknown task style #{task.style}"
                           end)
      end

      def validate
        super
        check_property :style, Symbol
        check_property :input, Symbol
        check_property :manifest, String if @style == 'puppet'
        raise "Unknown style #{@style}. Expected #{STYLES.join(', ')}" \
          unless STYLES.include?(@style)
        raise "Unknown input #{@input}. Expected #{INPUTS.join(', ')}" \
          unless INPUTS.include?(@input)
        raise 'Manifests can only be specified for :puppet style tasks' \
          unless @manifest.nil? || @style != 'puppet'
      end

      # A user provided argument to a Bolt task
      class Argument < Provider::Config::Function::Argument
        attr_reader :required
        attr_reader :default
        attr_reader :comment

        def description_display
          [
            @description.split("\n").map(&:strip).join(' '),
            ("(default: #{@default.display})" if default?)
          ].compact.join(' ')
        end

        def type_metadata(_provider)
          puppet_type = type.split('::').last
          optional_wrapper(puppet_type == 'String' ? 'String[1]' : puppet_type)
        end

        def validate
          super
          check_optional_property :required, :boolean
          check_optional_property :default
          check_optional_property :comment
          @default = Default.new(@default) \
            if default? && !@default.is_a?(Default)
        end

        def default?
          !@default.nil?
        end

        def required?
          !@required.nil? || @required
        end

        # Represents an argument can only accept form a preset list of values.
        class Enum < Argument
          attr_reader :values

          def validate
            @type = Api::Type::Enum.name
            super
            check_property_list :values, Symbol
          end

          def type_metadata(provider)
            vs = values.map { |v| provider.quote_string(v.id2name) }.join(', ')
            optional_wrapper("Enum[#{vs}]")
          end
        end

        # Definitions of the default value for a Bolt task.
        class Default < Api::Object
          attr_reader :code
          attr_reader :display

          def validate
            super
            check_property :display
          end

          private

          def initialize(value)
            @code = value
            @display = value
          end
        end

        private

        def optional_wrapper(metadata)
          required? ? metadata : "Optional[#{metadata}]"
        end
      end
    end

    def generate(output_folder, types, version)
      generate_client_functions output_folder unless @config.functions.nil?
      generate_bolt_tasks output_folder unless @config.bolt_tasks.nil?
      super(output_folder, types, version)
    end

    def compile_examples(output_folder)
      compile_file_map(
        output_folder,
        @config.examples,
        lambda do |_object, file|
          ["examples/#{file}",
           "products/#{@api.prefix[1..-1]}/examples/puppet/#{file}"]
        end
      )
    end

    def compile_end2end_tests(output_folder)
      compile_file_map(
        output_folder,
        @config.examples,
        lambda do |_object, file|
          # Tests go into hidden folder because we don't need to expose
          # to regular users.
          ["#{TEST_FOLDER}/#{file}",
           "products/#{@api.prefix[1..-1]}/examples/puppet/#{file}"]
        end
      )
    end

    def property_body(property)
      lines(
        indent([
          (['newvalue(:true)', 'newvalue(:false)'] \
           if property.is_a? Api::Type::Boolean),
          (generate_enum_body(property) if property.is_a? Api::Type::Enum),
          ("defaultto #{ruby_literal(property.default_value)}" \
         if property.default_value)
        ].compact.flatten, 4)
      )
    end

    def format_description(object, spaces, container, suffix = '')
      description = build_description object, suffix
      # A single line description is of the form:
      # [  newparam(....)]
      # [    desc '<description>']
      # [  end]
      #
      # So the description line has 11 extra characters other than the message
      # itself.
      #
      # The indentation level of 'desc' attribute is 4, spaces=4 due to the
      # template (newparam=2, desc=newparam+2), leaving 7 characters to
      # "compensate" for.
      format([
               ["#{container} #{quote_string(description)}"],
               [
                 "#{container} <<-DOC",
                 wrap_field(description, spaces),
                 'DOC'
               ]
             ], spaces)
    end

    def build_description(object, suffix)
      [object.description, suffix].reject(&:empty?)
                                  .join(' ').tr("\n", ' ')
                                  .gsub('  ', ' ').strip
    end

    def generate_user_agent(product, file_name)
      emit_user_agent(
        product, 'Puppet[:http_user_agent]',
        ['TODO(nelsonjr): Check how to fetch module version.'],
        file_name
      )
    end

    def quote_string(value)
      raise 'Invalid value' if value.nil?
      # Puppet DSL uses '${' to string interpolation while Ruby uses #{
      if value.include?('${')
        ['"', value, '"'].join
      else
        super(value)
      end
    end

    # Generates the documentation for the Puppet client side function to be
    # included in the module. Call this function immediately before the function
    # definition and the code generator will use data from api.yaml to build the
    # documentation comment block.
    #
    # rubocop: Method returns a big array. Easier to read a single block
    # rubocop:disable Metrics/AbcSize
    # rubocop:disable Metrics/MethodLength
    def emit_function_doc(function)
      [
        function.description.strip,
        '',
        'Arguments:',
        indent(function.arguments.map do |arg|
                 [
                   "- #{arg.name}: #{arg.type.split('::').last.downcase}",
                   indent(arg.description.strip.split("\n"), 2)
                 ]
               end, 2),
        (
          unless function.examples.nil?
            [
              '',
              'Examples:',
              indent(function.examples.map do |eg|
                       "- #{expand_function_vars(function, eg)}"
                     end, 2)
            ]
          end
        ),
        (
          unless function.notes.nil?
            [
              '',
              expand_function_vars(function, function.notes.strip)
            ]
          end
        )
      ].compact.flatten.join("\n").split("\n").map { |l| "# #{l}".strip }
    end
    # rubocop:enable Metrics/MethodLength
    # rubocop:enable Metrics/AbcSize

    # Returns true if any copy: or compile: entries there are any spec stub
    # files defined.
    def spec_stubs?
      @config.files.copy.any? { |f, _| f.start_with?(STUBS_FOLDER) } \
        || @config.files.compile.any? { |f, _| f.start_with?(STUBS_FOLDER) }
    end

    # rubocop:disable Metrics/MethodLength
    def emit_bolt_params_ruby(task)
      raise 'Only :stdin style supported' unless task.input == :stdin
      [
        'params = {}',
        'begin',
        indent(
          [
            'Timeout.timeout(3) do',
            indent('params = JSON.parse(STDIN.read)', 2),
            'end'
          ], 2
        ),
        'rescue Timeout::Error',
        indent(
          [
            ['puts(',
             "{ status: 'failure', error: 'Cannot read JSON from stdin' }",
             '.to_json)'].join,
            'exit 1'
          ], 2
        ),
        'end',
        ''
      ].concat(
        task.arguments.map do |arg|
          if arg.default?
            "#{arg.name} = validate(params, :#{arg.name}, #{arg.default.code})"
          else
            "#{arg.name} = validate(params, :#{arg.name})"
          end
        end
      )
    end
    # rubocop:enable Metrics/MethodLength

    private

    def generate_resource(data)
      generate_type data
      generate_provider data
    end

    def generate_provider_tests(data)
      super(data) \
        unless true?(Google::HashUtils.navigate(data[:config], %w[manual]))
    end

    def generate_simple_property(type, data)
      {
        source: File.join('templates', 'puppet', 'property', "#{type}.rb.erb"),
        target: File.join('lib', 'google', data[:product_name], 'property',
                          "#{type}.rb")
      }
    end

    def generate_base_property(data)
      {
        source: File.join('templates', 'puppet', 'property', 'base.rb.erb'),
        target: File.join('lib', 'google', data[:product_name], 'property',
                          'base.rb')
      }
    end

    def generate_typed_array(data, prop)
      type = Module.const_get(prop.item_type).new(prop.name).type
      file = Google::StringUtils.underscore(type)
      prop_map = []
      prop_map << {
        source: File.join('templates', 'puppet', 'property',
                          'array_typed.rb.erb'),
        target: File.join('lib', 'google', data[:product_name], 'property',
                          "#{file}_array.rb"),
        overrides: { type: type }
      }
      prop_map << generate_base_array(data)
      prop_map
    end

    def generate_base_array(data)
      {
        source: File.join('templates', 'puppet', 'property', 'array.rb.erb'),
        target: File.join('lib', 'google', data[:product_name], 'property',
                          'array.rb')
      }
    end

    def emit_nested_object(data)
      target = if data[:emit_array]
                 data[:property].item_type.property_file
               else
                 data[:property].property_file
               end
      result = [
        {
          source: File.join('templates', 'puppet', 'property',
                            'nested_object.rb.erb'),
          target: "lib/#{target}.rb",
          overrides: emit_nested_object_overrides(data)
        }
      ]

      result << generate_simple_property('array', data) if data[:emit_array]

      result
    end

    def emit_nested_object_overrides(data)
      data.clone.merge(
        field_name: Google::StringUtils.camelize(data[:field], :upper),
        object_type: Google::StringUtils.camelize(data[:obj_name], :upper),
        product_ns: Google::StringUtils.camelize(data[:product_name], :upper),
        class_name: if data[:emit_array]
                      data[:property].item_type.property_class.last
                    else
                      data[:property].property_class.last
                    end
      )
    end

    def emit_resourceref_object(data)
      target = data[:property].property_file
      {
        source: File.join('templates', 'puppet', 'property',
                          'resourceref.rb.erb'),
        target: "lib/#{target}.rb",
        overrides: data.clone.merge(
          class_name: data[:property].property_class.last
        )
      }
    end

    def generate_enum_body(property)
      property.values.collect do |value|
        if value.is_a?(Symbol)
          "newvalue(:#{value})"
        elsif value.is_a?(String)
          "newvalue(#{quote_string(value)})"
        else
          "#{value.class}newvalue(#{value})"
        end
      end
    end

    def google_lib_basic(file, product_ns)
      google_lib_basic_files(file, product_ns, 'lib', 'google')
    end

    def google_lib_network(file, product_ns)
      google_lib_network_files(file, product_ns, 'lib', 'google')
    end

    # Emits all the Puppet client functions available for use by end users.
    def generate_client_function(output_folder, func)
      target_folder = File.join(output_folder, 'lib', 'puppet', 'functions')
      {
        fn: func,
        target_folder: target_folder,
        template: 'templates/puppet/function.erb',
        output_folder: output_folder,
        out_file: File.join(target_folder, "#{func.name}.rb")
      }
    end

    # Emits all the Bolt tasks
    def generate_bolt_tasks(output_folder)
      target_folder = File.join(output_folder, 'tasks')
      FileUtils.mkpath target_folder

      generate_bolt_readme(output_folder)

      @config.bolt_tasks.each do |task|
        generate_file(
          name: "Bolt task #{task.name} (json)",
          task: task,
          template: 'templates/puppet/bolt~task.json.erb',
          output_folder: output_folder,
          out_file: File.join(target_folder, "#{task.name}.json")
        )

        generate_bolt_tasks_code(task, output_folder)
      end
    end

    def generate_bolt_readme(output_folder)
      generate_file(
        name: 'Bolt task README.md',
        template: 'templates/puppet/bolt~README.md.erb',
        output_folder: output_folder,
        out_file: File.join(output_folder, 'tasks/README.md')
      )
    end

    def generate_bolt_tasks_code(task, output_folder)
      style = task.style
      template = File.join('templates', 'puppet', case style
                                                  when :ruby
                                                    'bolt~task.rb.erb'
                                                  when :puppet
                                                    'bolt~task.pp.erb'
                                                  else
                                                    raise "No style #{style}"
                                                  end)
      out_file = File.join(output_folder, task.target_file)

      generate_file(
        name: "Bolt task #{task.name}",
        task: task,
        template: template,
        output_folder: output_folder,
        out_file: out_file
      )

      FileUtils.chmod BOLT_TASK_ACL, out_file
    end
    # rubocop:enable Metrics/MethodLength

    def expand_function_vars(fn_config, data)
      data.gsub('{{function:name}}', fn_config.name)
    end
  end
end
