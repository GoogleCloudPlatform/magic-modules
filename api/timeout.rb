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

module  Api

  # Provides timeout information for the different operation types
  class Timeouts < Api::Object
    # Default timeout for all operation types is 4 minutes. This can be
    # overridden for each resource.
    DEFAULT_INSERT_TIMEOUT_SEC = 4 * 60
    DEFAULT_UPDATE_TIMEOUT_SEC = 4 * 60
    DEFAULT_DELETE_TIMEOUT_SEC = 4 * 60

    attr_reader :insert_sec
    attr_reader :update_sec
    attr_reader :delete_sec

    def initialize
      validate
    end

    def validate
      super

      @insert_sec ||= DEFAULT_INSERT_TIMEOUT_SEC
      @update_sec ||= DEFAULT_UPDATE_TIMEOUT_SEC
      @delete_sec ||= DEFAULT_DELETE_TIMEOUT_SEC

      check_property :insert_sec, Integer
      check_property :update_sec, Integer
      check_property :delete_sec, Integer
    end
  end

end