module Provider
  module Azure
    module Ansible
      module Example
        module SubTemplate
          def build_yaml_from_example(example)
            yaml = to_yaml({
              'name' => example.description,
              example.resource => example.properties.transform_keys(&:underscore)
            })
            lines = yaml.split("\n")
            lines('- ' + lines[0]) + indent(lines[1..-1], 2)
          end
        end
      end
    end
  end
end
