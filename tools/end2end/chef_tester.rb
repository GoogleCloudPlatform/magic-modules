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

require 'tester_base'

module Chef
  # Executes all examples against a real Google Cloud Platform project. The
  # account requires to have Owner (or Editor) to all resources being tested.
  class Tester < TesterBase
    def provider
      'chef'
    end

    def header
      <<-HEADER
       _______ _     _ _______ _______
       |       |_____| |______ |______
       |_____  |     | |______ |

      HEADER
    end

    private

    def command(data)
      %w[bundle exec chef-client -z --runlist] \
        << "recipe[#{data['run']}]"
    end

    def variables(env)
      super env
    end
  end
end
