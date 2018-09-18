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

# ----------------------------------------------------------------------------
#
#     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
#
# ----------------------------------------------------------------------------
#
#     This file is automatically generated by Magic Modules and manual
#     changes will be clobbered when the file is regenerated.
#
#     Please read more about how to change this file in README.md and
#     CONTRIBUTING.md located at the root of this package.
#
# ----------------------------------------------------------------------------

require 'google/compute/property/array'
module Google
  module Compute
    module Property
      # A class to manage data for CustomMetricUtilizations for region_autoscaler.
      class RegionAutoscalerCustomMetricUtilizations
        include Comparable

        attr_reader :metric
        attr_reader :utilization_target
        attr_reader :utilization_target_type


        def initialize(args = nil) 
          return nil if args.nil?
          @metric = args['metric']
          @utilization_target = args['utilizationTarget']
          @utilization_target_type = args['utilizationTargetType']
        end
      end


      class RegionAutoscalerCustomMetricUtilizationsArray < Google::Compute::Property::Array

        def self.parse(value)
          return if value.nil?
          return RegionAutoscalerCustomMetricUtilizations.new(value) unless value.is_a?(::Array)
          value.map { |v| RegionAutoscalerCustomMetricUtilizations.new(v) }
        end
      end
    end
  end
end
