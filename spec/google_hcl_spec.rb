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

require 'spec_helper'
require 'google/hcl_utils'

class Test
  include Google::HclUtils
end

describe Google::GolangUtils do
  context '#go_literal' do
    let(:hcl) { Test.new }

    describe 'sample case' do
      let(:original) {
        {
          "variable" => [{
            "image" => [{
              "my-image" => [{
                "description": "the Image to use"
              }]
            }]
          }]
        }
      }

      let(:final) {
        [
          "variable \"image\" \"my-image\" {",
          "  description = \"the Image to use\"",
          '}'
        ].join("\n")
      }

      subject { hcl.hcl(original) }
      it { is_expected.to eq final }
    end

    describe 'correct spacing' do
      let(:original) {
        {
          "variable" => [{
            "image" => [{
              "my-image" => [{
                "description": "the Image to use",
                "a really long name": "long name",
                "short": "a short name"
              }]
            }]
          }]
        }
      }

      let(:final) {
        [
          "variable \"image\" \"my-image\" {",
          "  description        = \"the Image to use\"",
          "  a really long name = \"long name\"",
          "  short              = \"a short name\"",
          '}'
        ].join("\n")
      }

      subject { hcl.hcl(original) }
      it { is_expected.to eq final }
    end
  end
end
