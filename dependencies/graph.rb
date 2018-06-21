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

module Dependencies
  # Represents a full graph made up of Api::Resources
  class Graph
    attr_reader :node
    attr_reader :pointers
    attr_reader :start

    def initialize(object, seed)
      @node = ItemListNode.new(object, seed, nil)
      @pointers = { object.name => @node }
      @start = object.name
    end

    # Adds an object to the graph. This object is represented by a ResourceRef
    # property
    def add_ref(prop, seed)
      raise 'Only ResourceRef may be inserted' unless \
        prop.is_a? Api::Type::ResourceRef

      # Since this is testing, we're always going to use the first reference.
      # It doesn't matter which we choose, but we have to choose one.
      reference = prop.resources[0]
      # Add each item to its proper node.

      referenced_by = reference.__resource.name
      object = reference.resource_ref.name

      if @pointers.key?(object)
        @pointers[object].add(reference.resource_ref, seed, prop)
      else
        node = ItemListNode.new(reference.resource_ref, seed, prop)
        @pointers[object] = node
      end

      # Ensure that a link exists between this object and the object
      # that referenced it.
      @pointers[referenced_by].add_child(object)
    end

    def add_object(object, seed)
      @pointers[object.name].add(object, seed, nil)
    end

    def each
      sort.each do |obj_type|
        @pointers[obj_type].objects.each do |obj|
          yield obj
        end
      end
    end

    def map
      sort.map do |obj_type|
        @pointers[obj_type].objects.map do |obj|
          yield(obj)
        end
      end
    end

    private

    # Depth first search through grid
    # See https://en.wikipedia.org/wiki/Topological_sorting for algorithm
    # (listed under Depth First Search)
    def sort
      markers = @pointers.keys.product(['not visited']).to_h
      sorted = []
      visit(markers.keys[0], sorted, markers) until markers.empty?
      sorted.reverse
    end

    def visit(node, sorted, markers)
      return unless markers.key?(node)
      return if markers[node] != 'not visited'
      markers[node] = 'temp'
      @pointers[node].children.each do |child|
        visit(child, sorted, markers)
      end
      markers.delete(node)
      sorted.unshift(node)
    end
  end

  # A list of objects of a certain type.
  # This represents a single node in the graph.
  class ItemListNode
    attr_reader :type
    attr_reader :objects
    attr_reader :children

    def initialize(object, seed, parent)
      @type = object.name
      @objects = []
      @children = []

      add(object, seed, parent)
    end

    def add(object, seed, parent)
      @objects << Item.new(object, seed, parent)
      @objects = @objects.sort_by(&:seed).uniq
    end

    def add_child(name)
      @children << name unless @children.include? name
    end
  end

  # A single instance of an object with parent property and seed value.
  # No knowledge of the graph is necessary when handling an Item.
  class Item
    attr_reader :object
    attr_reader :seed
    attr_reader :parent

    def initialize(object, seed, parent)
      @object = object
      @seed = seed
      @parent = parent
    end
  end
end
