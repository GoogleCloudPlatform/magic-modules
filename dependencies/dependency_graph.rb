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

require 'dependencies/graph'

module Dependencies
  # Responsible for holding all objects being used throughout the creation
  # of unit tests.
  #
  # These objects include the object being tested and all of the needed
  # dependencies to create that object.
  class DependencyGraph
    def initialize(datagen)
      @datagen = datagen
    end

    def add(object, seed, kind, extra)
      if @graph.nil?
        @graph = Graph.new(object, seed)
      else
        @graph.add_object(object, seed)
      end
      collect_refs(object.all_user_properties, kind, seed, extra)
    end

    def each
      @graph.each { |node| yield node }
    end

    def map
      @graph.map { |node| yield(node) }
    end

    private

    def collect_refs(props, kind, seed, extra)
      props = select_properties(props, kind, extra)
      collect_refs_child(props, seed)
    end

    def collect_refs_child(props, seed)
      props.each do |p|
        emit_manifest_assign(p, seed)
      end
    end

    # prop_name_method is a method that returns the proper name of the object
    # rref_value: If true, return the value being exported by the ref'd block
    #             If false, return the title of the block being referenced
    def emit_manifest_assign(prop, seed)
      if prop.is_a?(Api::Type::Array)
        emit_manifest_array(prop, seed)
      elsif prop.is_a?(Api::Type::NestedObject)
        emit_nested(prop, seed)
      elsif prop.is_a?(Api::Type::ResourceRef)
        emit_resource(prop, seed)
      end
    end

    # rubocop:enable Metrics/AbcSize
    # rubocop:enable Metrics/CyclomaticComplexity
    # rubocop:enable Metrics/MethodLength
    # rubocop:enable Metrics/PerceivedComplexity
    # prop_name_method should be a valid method on a Api::Type::*
    # Typically, this will be "out_name" or "field_name"
    def emit_manifest_array(prop, seed)
      if prop.item_type.is_a?(Api::Type::NestedObject)
        size = @datagen.object_size(prop, seed, true)
        (1..size).each do |index|
          collect_refs(prop.item_type.properties, :title, seed + index - 1,
                       ensure: 'present')
        end
      elsif prop.item_type.is_a?(Api::Type::ResourceRef)
        size = @datagen.object_size(prop, seed, true)
        (1..size).each do |index|
          emit_resource(prop.item_type, (seed + index - 1) % 3)
        end
      end
    end
    # rubocop:enable Metrics/MethodLength

    def emit_nested(prop, seed)
      collect_refs_child(prop.properties, seed)
    end

    def emit_resource(prop, seed)
      # % 3 because only 3 different network test data files per Resource
      @graph.add_ref(prop, seed % 3)

      # Because this is testing, we're always going to use the first
      # ResourceRef in a list of ResourceRefs.

      # Recurse through referenced object for more resourcerefs
      # Don't recurse on resourceref of same type.
      return if prop.resources[0].resource_ref == prop.resources[0].__resource

      # When building resourcerefs in manifests/catalogs, we use the
      # smallest set of properties possible. When looking recursively, we
      # should only look through this smallest set of properties. Otherwise,
      # recursive resource refs may be created that do not get referenced.
      collect_refs(prop.resources[0].resource_ref.required_properties,
                   :name, seed, {})
    end

    # rubocop:disable Metrics/AbcSize
    # rubocop:disable Metrics/CyclomaticComplexity
    # rubocop:disable Metrics/PerceivedComplexity
    def select_properties(props, kind, extra)
      props = props.reject(&:output)

      props = props.select { |p| p.required || p.name == 'name' } \
        if (extra.key?(:ensure) && extra[:ensure].to_sym == :absent) ||
           (extra.key?(:action) && extra[:action] == ':delete')

      if kind == :resource
        props.select { |p| p.required || p.name == 'name' }
      elsif kind == :title
        props.reject { |p| p.name == 'name' }
      else
        props
      end
    end
    # rubocop:enable Metrics/PerceivedComplexity
    # rubocop:enable Metrics/CyclomaticComplexity
    # rubocop:enable Metrics/AbcSize
    # rubocop:enable Metrics/ClassLength
  end
end
