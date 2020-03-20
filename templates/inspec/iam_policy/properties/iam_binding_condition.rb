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

module GoogleInSpec
  module Iam
    module Property
      class IamBindingCondition
        attr_reader :title

        attr_reader :description

        attr_reader :expression

        def initialize(args = nil, parent_identifier = nil)
          return if args.nil?
          @parent_identifier = parent_identifier
          @title = args['title']
          @description = args['description']
          @expression = args['expression']
        end

        def to_s
          "#{@parent_identifier} IamBindingCondition"
        end
      end
    end
  end
end