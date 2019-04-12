# frozen_string_literal: false

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

require 'google/iam/property/iam_policy_audit_configs_audit_log_configs'
module GoogleInSpec
  module Iam
    module Property
      class IamPolicyAuditConfigs
        attr_reader :service

        attr_reader :audit_log_configs

        def initialize(args = nil, parent_identifier = nil)
          return if args.nil?
          @parent_identifier = parent_identifier
          @service = args['service']
          @audit_log_configs = GoogleInSpec::Iam::Property::IamPolicyAuditConfigsAuditLogConfigsArray.parse(args['auditLogConfigs'], to_s)
        end

        def to_s
          "#{@parent_identifier} IamPolicyAuditConfigs"
        end
      end

      class IamPolicyAuditConfigsArray
        def self.parse(value, parent_identifier)
          return if value.nil?
          return IamPolicyAuditConfigs.new(value, parent_identifier) unless value.is_a?(::Array)
          value.map { |v| IamPolicyAuditConfigs.new(v, parent_identifier) }
        end
      end
    end
  end
end
