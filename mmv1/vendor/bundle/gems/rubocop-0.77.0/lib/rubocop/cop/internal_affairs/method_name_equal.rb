# frozen_string_literal: true

module RuboCop
  module Cop
    module InternalAffairs
      # Checks that method names are checked using `method?` method.
      #
      # @example
      #   # bad
      #   node.method_name == :do_something
      #
      #   # good
      #   node.method?(:do_something)
      #
      class MethodNameEqual < Cop
        include RangeHelp

        MSG = 'Use `method?(%<method_name>s)` instead of ' \
              '`method_name == %<method_name>s`.'

        def_node_matcher :method_name?, <<~PATTERN
          (send
            $(send
              (...) :method_name) :==
            $...)
        PATTERN

        def on_send(node)
          method_name?(node) do |method_name_node, method_name_arg|
            message = format(MSG, method_name: method_name_arg.first.source)

            range = range(method_name_node, node)

            add_offense(node, location: range, message: message)
          end
        end

        def autocorrect(node)
          lambda do |corrector|
            method_name?(node) do |method_name_node, method_name_arg|
              corrector.replace(
                range(method_name_node, node),
                "method?(#{method_name_arg.first.source})"
              )
            end
          end
        end

        private

        def range(method_name_node, node)
          range_between(
            method_name_node.loc.selector.begin_pos, node.source_range.end_pos
          )
        end
      end
    end
  end
end
