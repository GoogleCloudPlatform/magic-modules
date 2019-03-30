module Provider
  module Azure
    module Ansible
      module SubTemplate
        def build_multiline_method_call(call_prefix, args, call_postfix, indentation = 0)
          multiline_args = args.join(",\n#{' ' * call_prefix.length}")
          result = "#{call_prefix}#{multiline_args}#{call_postfix}"
          indent result, indentation
        end
      end
    end
  end
end