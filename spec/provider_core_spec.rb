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
    subject { described_class.new(mock('config'), mock('api')) }

    it 'fails if cannot fit' do
      expect do
        subject.format [['x' * 21]], 0, 0, 20
      end.to raise_error ArgumentError, /No code fits/
    end

    it 'fails if cannot fit any' do
      expect do
        subject.format [['x' * 21], ['y' * 21], ['z' * 30]], 0, 0, 20
      end.to raise_error ArgumentError, /No code fits/
    end

    it 'fits 80 chars' do
      subject.format [['x' * 80]]
    end

    context 'fits 80 chars' do
      subject do
        described_class.new(nil, nil).format([
                                               ['x' * 80],
                                               ['y' * 80],
                                               ['z' * 80]
                                             ])
      end

      it { is_expected.to include 'x' }
    end

    context '#format(ident)' do
      it 'fits' do
        subject.format [['x' * 74]], 6
      end

      it 'does not fit' do
        expect do
          subject.format [['x' * 75]], 6
        end.to raise_error(ArgumentError, /No code fits/)
      end
    end

    context '#format(start)' do
      it 'fits' do
        subject.format [['x' * 74]], 0, 6
      end

      it 'does not fit' do
        expect do
          subject.format [['x' * 75]], 0, 6
        end.to raise_error(ArgumentError, /No code fits/)
      end
    end

    context '#format(start, indent)' do
      it 'fits' do
        subject.format [['x' * 66]], 8, 6
      end

      it 'does not fit' do
        expect do
          subject.format [['x' * 67]], 8, 6
        end.to raise_error(ArgumentError, /No code fits/)
      end
    end

    context 'selects second option' do
      subject do
        described_class.new(nil, nil).format([
                                               ['x' * 81],
                                               ['y' * 80],
                                               ['z' * 80]
                                             ])
      end

      it { is_expected.to include 'y' }
    end
  end
end
