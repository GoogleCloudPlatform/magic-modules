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

describe Provider::Core do
  context '#format' do
    subject { described_class.new(mock('config'), mock('api'), Time.now) }

    it 'does not fail if cannot fit' do
      expect(
        subject.format([['x' * 41]], 0, 0, 20)
      ).to include('rubocop:disable Metrics/LineLength')
    end

    it 'does not fail if cannot fit any' do
      expect(
        subject.format([['x' * 31], ['y' * 31], ['z' * 30]], 0, 0, 20)
      ).to include 'rubocop:disable Metrics/LineLength'
    end

    it 'fits 100 chars' do
      expect(
        subject.format([['x' * 100]])
      ).to eq('x' * 100)
    end

    context 'fits 100 chars' do
      subject do
        described_class.new(nil, nil, nil).format([
                                                    ['x' * 100],
                                                    ['y' * 100],
                                                    ['z' * 100]
                                                  ])
      end

      it { is_expected.to include 'x' }
    end

    context '#format(ident)' do
      it 'fits' do
        expect(
          subject.format([['x' * 74]], 6)
        ).to eq((' ' * 6) + ('x' * 74))
      end

      it 'does not fit' do
        expect(
          subject.format([['x' * 95]], 6)
        ).to include 'rubocop:disable Metrics/LineLength'
      end
    end

    context '#format(start)' do
      it 'fits' do
        expect(
          subject.format([['x' * 74]], 0, 6)
        ).to eq('x' * 74)
      end

      it 'does not fit' do
        expect(
          subject.format([['x' * 115]], 0, 6)
        ).to include 'rubocop:disable Metrics/LineLength'
      end
    end

    context '#format(start, indent)' do
      it 'fits' do
        expect(
          subject.format([['x' * 66]], 8, 6)
        ).to eq((' ' * 8) + ('x' * 66))
      end

      it 'does not fit' do
        expect(
          subject.format([['x' * 87]], 8, 6)
        ).to include 'rubocop:disable Metrics/LineLength'
      end
    end

    context 'selects second option' do
      subject do
        described_class.new(nil, nil, nil).format([
                                                    ['x' * 101],
                                                    ['y' * 80],
                                                    ['z' * 80]
                                                  ])
      end

      it { is_expected.to include 'y' }
    end
  end
end
