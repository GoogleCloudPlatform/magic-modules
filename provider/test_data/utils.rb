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
module Provider
  module TestData
    # These are functions needed for generating test data. These should only
    # be used in the context of generating unit tests.
    # These functions will need to be exposed to templates.
    module TestUtils
      # Given an array of properties, return all ResourceRefs contained within
      # Requires:
      #   props- a list of props
      #   original_object - the original object containing props. This is to
      #                     avoid self-referencing objects.
      # For ResourceRefs with multiple refs, only the first will be used.
      # This is a assumption used throughout the generation of unit tests.
      # rubocop:disable Metrics/AbcSize
      # rubocop:disable Metrics/CyclomaticComplexity
      # rubocop:disable Metrics/PerceivedComplexity
      def test_resourcerefs_for_properties(props, original_obj)
        rrefs = []
        props.each do |p|
          # We need to recurse on ResourceRefs to get all levels
          # We do not want to recurse on resourcerefs of type self to avoid
          # infinite loop.
          if p.is_a? Api::Type::ResourceRef
            # We want to avoid a circular reference
            # This reference may be the next reference or have some number of refs
            # in between it.
            next if p.resources[0].resource_ref == original_obj
            next if p.resources[0].resource_ref == p.resources[0].__resource
            rrefs << p
            rrefs.concat(test_resourcerefs_for_properties(p.resources[0].resource_ref
                                                      .required_properties,
                                                     original_obj))
          elsif p.is_a? Api::Type::NestedObject
            rrefs.concat(test_resourcerefs_for_properties(p.properties, original_obj))
          elsif p.is_a? Api::Type::Array
            if p.item_type.is_a? Api::Type::NestedObject
              rrefs.concat(test_resourcerefs_for_properties(p.item_type.properties,
                                                       original_obj))
            elsif p.item_type.is_a? Api::Type::ResourceRef
              rrefs << p.item_type
              rrefs.concat(test_resourcerefs_for_properties(p.item_type.resources[0].resource_ref
                                                        .required_properties,
                                                       original_obj))
            end
          end
        end
        rrefs.uniq
      end
      # rubocop:enable Metrics/AbcSize
      # rubocop:enable Metrics/CyclomaticComplexity
      # rubocop:enable Metrics/MethodLength
      # rubocop:enable Metrics/PerceivedComplexity
    end
  end
end
