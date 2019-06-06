module Provider
  module Azure
    module Ansible
      module Example
        module SubTemplate
          def build_test_yaml_from_example(example_name, name_postfix = nil, register_name = 'output')
            random_vars = Set.new
            deps, main = build_yaml_from_example(nil, example_name, random_vars, name_postfix, {}, register_name)
            return deps, main, random_vars
          end

          def build_documentation_yaml_from_example(example)
            _, main = build_yaml_from_example(nil, example.example, Set.new, nil, example.resource_name_hints, nil)
            return main
          end

          def build_yaml_from_example(product_name, example_name, random_variables, name_postfix, name_hints, register_name)
            example = get_example_by_names(example_name, product_name)
            yaml_deps = compile 'templates/azure/ansible/example/example_deps_yaml.erb', 1
            yaml_raw = compile 'templates/azure/ansible/example/example_yaml.erb', 1
            context = ExampleContextBinding.new(name_hints, random_variables)
            yaml_deps = compile_string context.get_binding, yaml_deps
            yaml_raw = compile_string context.get_binding, yaml_raw
            return yaml_deps, yaml_raw
          end

          def build_yaml_properties(properties, indentation = 2)
            result = compile 'templates/azure/ansible/example/yaml_properties.erb', 1
            indent result, indentation
          end

          private

          class ExampleContextBinding
            attr_reader :my_binding
            attr_reader :name_hints
            attr_reader :random_variables
  
            def initialize(name_hints, random_vars)
              @my_binding = binding
              @name_hints = name_hints
              @random_variables = random_vars
            end
  
            def get_binding()
              @my_binding
            end
  
            def get_resource_name(name_hint, random_var_name, random_var_prefix = '')
              return "#{name_hints[name_hint]}\n" if name_hints.has_key?(name_hint)
              @random_variables << RandomizedVariable.new(:Standard, random_var_name, random_var_prefix)
              "\"{{ #{random_var_name} }}\"\n"
            end
          end

          class RandomizedVariable
            attr_reader :variable_name
            attr_reader :variable_value
  
            def initialize(type, var_name, prefix)
              case type
              when :Standard
                @variable_name = var_name
                @variable_value = prefix + "{{ resource_group | hash('md5') | truncate(7, True, '') }}{{ 1000 | random }}"
              end
            end

            def hash
              hash = 17 * 31 + @variable_name.hash
              hash * 31 + @variable_value.hash
            end

            def eql?(other)
              return false unless other.is_a?(RandomizedVariable)
              @variable_name.eql?(other.variable_name) && @variable_value.eql?(other.variable_value)
            end
          end
        end
      end
    end
  end
end
