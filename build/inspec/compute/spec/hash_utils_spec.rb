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

describe Google::HashUtils do
  context '#camelize_keys' do
    subject do
      described_class.camelize_keys(a_a_a: 'aaa',
                                    bb_bb_bb: 'bbb',
                                    abc_def_ghi: 'abcdefghi')
    end

    it do
      is_expected.to eq('aAA' => 'aaa',
                        'bbBbBb' => 'bbb',
                        'abcDefGhi' => 'abcdefghi')
    end
  end

  context '#navigate' do
    let(:source) { { a: { b: { c: %i[d e] } } } }
    let(:default) { Object.new }

    context 'find item middle' do
      subject { described_class.navigate(source, %i[a b]) }
      it { is_expected.to eq(c: %i[d e]) }
    end

    context 'find item leaf' do
      subject { described_class.navigate(source, %i[a b c]) }
      it { is_expected.to eq(%i[d e]) }
    end

    context 'item does not exist' do
      subject { described_class.navigate(source, %i[d]) }
      it { is_expected.to be nil }
    end

    context 'returns default' do
      subject { described_class.navigate(source, %i[d], default) }
      it { is_expected.to eq default }
    end
  end
end
