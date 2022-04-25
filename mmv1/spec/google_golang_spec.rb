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
require 'google/golang_utils'

class Test
  include Google::GolangUtils
end

describe Google::GolangUtils do
  context '#go_literal' do
    let(:golang) { Test.new }

    describe 'string' do
      subject { golang.go_literal('foo') }
      it { is_expected.to eq '"foo"' }
    end

    describe 'integer' do
      subject { golang.go_literal(123) }
      it { is_expected.to eq '123' }
    end

    describe 'float' do
      subject { golang.go_literal(0.987) }
      it { is_expected.to eq '0.987' }
    end

    describe 'symbol' do
      subject { golang.go_literal(:NONE) }
      it { is_expected.to eq '"NONE"' }
    end

    describe 'empty_array' do
      subject { golang.go_literal([]) }
      it { is_expected.to eq '[]string{}' }
    end

    describe 'string_array_single' do
      subject { golang.go_literal(['abc']) }
      it { is_expected.to eq '[]string{"abc"}' }
    end

    describe 'string_array_multiple' do
      subject { golang.go_literal(%w[abc def]) }
      it { is_expected.to eq '[]string{"abc", "def"}' }
    end

    describe 'int_array' do
      subject { -> { golang.go_literal([1, 2]) } }
      it { is_expected.to raise_error(/Unsupported/) }
    end

    describe 'unknown type' do
      subject { -> { golang.go_literal(Class.new) } }
      it { is_expected.to raise_error(/Unsupported/) }
    end
  end
end
