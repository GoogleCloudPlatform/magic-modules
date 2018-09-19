# Copyright 2018 Google Inc.
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

module Api
  class Resource < Api::Object::Named
    class Example < Api::Object
      # An Example Resource is a resource block that will appear in an example
      # config alongside product documentation
      #
      # Example Resources have a name, a type, and a set of properties that will
      # be set in their config.
      class Resource < Api::Object
        include Compile::Core

        # The name/id of a provider resource. This is NOT the GCP `name` field.
        # This field is the Terraform URI name / the PCA equivalent.
        attr_reader :name

        # The type of the resource in the format {{mm-product}}/{{mm-resource}}
        # This form is weird and artificial. When we enter "gcompute/Address" we
        # want "compute_address" across TPCA. But we preserve the MM product and
        # name this way in case we need that information to do overrides or
        # something later on.
        attr_reader :type

        # properties is a Hash from field names to values.
        # Most values are just strings that will be read out verbatim.
        # Inputs starting with an @ sign will be read as interpolations of the
        # form:
        # @{{mm-product}}/{{mm-resource}}/{{name}}/{{field}}
        # For example:
        # @gcompute/Subnetwork/default/selfLink"
        # This should? contain enough information to allow referencing across
        # all 4 providers
        attr_reader :properties

        def validate
          super
          @properties ||= {}

          check_property :name, String
          check_property :type, String
          check_property :properties, Hash
        end
      end
    end
  end
end
