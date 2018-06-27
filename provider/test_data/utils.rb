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
      # rubocop:disable Metrics/MethodLength
      # rubocop:disable Metrics/BlockLength
      def test_resourcerefs_for_properties(props, original_obj)
        rrefs = []
        props.each do |p|
          # We need to recurse on ResourceRefs to get all levels
          # We do not want to recurse on resourcerefs of type self to avoid
          # infinite loop.
          if p.is_a? Api::Type::ResourceRef
            # We want to avoid a circular reference
            # This reference may be the
            # next reference or have some number of refs in between it.
            next if p.resource_refs.first.resource_ref == original_obj
            next if p.resource_refs.first.resource_ref == p.resource_refs.first.__resource
            rrefs << p
            rrefs.concat(test_resourcerefs_for_properties(
                           p.resource_refs.first.resource_ref.required_properties,
                           original_obj
            ))
          elsif p.is_a? Api::Type::NestedObject
            rrefs.concat(test_resourcerefs_for_properties(p.properties,
                                                          original_obj))
          elsif p.is_a? Api::Type::Array
            if p.item_type.is_a? Api::Type::NestedObject
              rrefs.concat(test_resourcerefs_for_properties(
                             p.item_type.properties,
                             original_obj
              ))
            elsif p.item_type.is_a? Api::Type::ResourceRef
              rrefs << p.item_type
              rrefs.concat(test_resourcerefs_for_properties(
                             p.item_type.resource_refs.first.resource_ref
                                                     .required_properties,
                             original_obj
              ))
            end
          end
        end
        rrefs.uniq
      end
      # rubocop:enable Metrics/AbcSize
      # rubocop:enable Metrics/CyclomaticComplexity
      # rubocop:enable Metrics/MethodLength
      # rubocop:enable Metrics/PerceivedComplexity
      # rubocop:enable Metrics/BlockLength

      def variable_type(object, var)
        return Api::Type::String::PROJECT if var == :project
        return Api::Type::String::NAME if var == :name
        v = object.all_user_properties
                  .select do |p|
          p.out_name.to_sym == var || p.name.to_sym == var
        end.first
        return v.resource_refs.first.property if v.is_a?(Api::Type::ResourceRef)
        v
      end
    end
  end
end
