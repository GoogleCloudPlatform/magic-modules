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

describe Google::IntegerUtils do
  context '#underscore' do
    context '0' do
      subject { described_class.underscore(0) }
      it { is_expected.to eq '0' }
    end

    context 'full groups' do
      subject { described_class.underscore(123_456_789) }
      it { is_expected.to eq '123_456_789' }
    end

    context 'non-full groups' do
      subject { described_class.underscore(1_234_567_890) }
      it { is_expected.to eq '1_234_567_890' }
    end

    context 'middle zeros' do
      subject { described_class.underscore(1_034_007_000) }
      it { is_expected.to eq '1_034_007_000' }
    end
  end
end
