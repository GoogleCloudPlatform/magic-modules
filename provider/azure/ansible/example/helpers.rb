module Provider
  module Azure
    module Ansible
      module Example
        module Helpers
          def generate_info_assert_list(example_name)
            example = get_example_by_names(example_name)
            asserts = ["- output.items[0]['id'] != None"]
            example.properties.each_key do |p|
              asserts << "- output.items[0]['#{p.underscore}'] != None"
            end
            asserts
          end
        end
      end
    end
  end
end
