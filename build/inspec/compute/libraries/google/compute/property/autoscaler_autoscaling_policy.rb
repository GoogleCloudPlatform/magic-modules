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

module Google
  module Compute
    module Property
      # A class to manage data for AutoscalingPolicy for autoscaler.
      class AutoscalerAutoscalingPolicy
        include Comparable

        attr_reader :min_num_replicas
        attr_reader :max_num_replicas
        attr_reader :cool_down_period_sec
        attr_reader :cpu_utilization
        attr_reader :custom_metric_utilizations
        attr_reader :load_balancing_utilization


        def initialize(args = nil) 
          return nil if args.nil?
          @min_num_replicas = args['minNumReplicas']
          @max_num_replicas = args['maxNumReplicas']
          @cool_down_period_sec = args['coolDownPeriodSec']
          @cpu_utilization = Google::Compute::Property::AutoscalerCpuUtilization.new(args['cpuUtilization'])
          @custom_metric_utilizations = Google::Compute::Property::AutoscalerCustomMetricUtilizationsArray.parse(args['customMetricUtilizations'])
          @load_balancing_utilization = Google::Compute::Property::AutoscalerLoadBalancingUtilization.new(args['loadBalancingUtilization'])
        end
      end

    end
  end
end
