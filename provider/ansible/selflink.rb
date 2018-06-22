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

require 'api/object'
require 'compile/core'
require 'provider/config'
require 'provider/core'
require 'provider/ansible/manifest'

module Provider
  module Ansible
    # Responsible for building out all logic for checking + assembling selflinks
    # on Ansible virtual properties.
    #
    # Some virtual properties (regions, zones) will probably be entered as
    # names by users (that's most intuitive). If a self-link is required,
    # the module should assemble that itself.
    module SelfLink
      # Returns all rrefs that are virtual and require selflinks.
      def virtual_selflink_rrefs(object)
        resourcerefs_for_properties(object.all_user_properties,
                                    object,
                                    virtual: 'only')
          .select { |prop| prop.resources.first.imports == 'selfLink' }
      end

      # Build out functions that will check + create selflinks.
      def selflink_functions(object)
        virtuals = virtual_selflink_rrefs(object).map { |x| x.resources.first }
                                                 .map(&:resource_ref)
                                                 .uniq
        virtuals.map do |virt|
          if virt == virtuals.last
            lines(selflink_function(virt))
          else
            lines(selflink_function(virt), 1)
          end
        end
      end

      def selflink_function(resource)
        url = self_link_url(resource).gsub('{project}', '.*')
                                     .gsub('{name}', '[a-z1-9\-]*')
        lines([
                method_decl(
                  "#{Google::StringUtils.underscore(resource.name)}_selflink",
                  %w[name params]
                ),
                indent(
                  [
                    'if name is None:',
                    indent('return', 4),
                    "url = r#{url}",
                    'if not re.match(url, name):',
                    # '%s' confuses Rubocop (it's Python code, not Ruby)
                    indent([
                      "name = #{self_link_url(resource).gsub('{name}', '%s')}",
                      '.format(**params) % name'
                    ].join, 4),
                    # rubocop:enable Style/FormatStringToken
                    'return name'
                  ], 4
                )
              ])
      end
      # rubocop:enable Metrics/MethodLength

      private

      def path_for_rref(rref)
        past_values = []
        until rref.nil?
          past_values << Google::StringUtils.underscore(rref.name)
          # TODO(alexstephen): Investigate a better way to handle parent
          # pointers on Arrays of NestedObjects
          rref = if rref.is_a?(Api::Type::NestedObject) && \
                    rref.parent.is_a?(Api::Type::Array)
                   rref.parent.parent
                 else
                   rref.parent
                 end
        end
        "[#{past_values.reverse.map { |x| quote_string(x) }.join(', ')}]"
      end
    end
  end
end
