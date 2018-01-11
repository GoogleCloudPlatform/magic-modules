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

describe Provider::TestMatrix do
  let(:provider) { mock('provider') }
  let(:object) { mock('object') }
  let(:expected) do
    {
      exists: {
        changes: %i[pass fail],
        no_change: %i[fail]
      }
    }
  end

  before(:each) do
    provider.stubs(:lines)
    object.stubs(:name).returns('test')
  end

  subject { described_class.new __FILE__, object, provider, expected }

  context 'good matrix' do
    before(:each) do
      subject.push(:ignore, :exists)
      subject.push(:ignore, :exists, :changes)
      subject.push(:ignore, :exists, :changes, :ignore, :pass)
      subject.pop(:ignore, :exists, :changes, :ignore, :pass)
      subject.push(:ignore, :exists, :changes, :ignore, :fail)
      subject.pop(:ignore, :exists, :changes, :ignore, :fail)
      subject.pop(:ignore, :exists, :changes)
      subject.push(:ignore, :exists, :no_change)
      subject.push(:ignore, :exists, :no_change, :ignore, :fail)
      subject.pop(:ignore, :exists, :no_change, :ignore, :fail)
      subject.pop(:ignore, :exists, :no_change)
      subject.pop(:ignore, :exists)
    end

    it { expect { subject.verify }.not_to raise_error }
  end

  context 'cannot have duplicates' do
    context 'sequential' do
      before(:each) do
        subject.push(:ignore, :exists)
      end

      it do
        expect { subject.push(:ignore, :exists) }.to \
          raise_error(RuntimeError, /already exists/)
      end
    end

    context 'out of order' do
      before(:each) do
        subject.push(:ignore, :exists)
        subject.push(:ignore, :exists, :changes)
      end

      it do
        expect { subject.push(:ignore, :exists) }.to \
          raise_error(RuntimeError, /already exists/)
      end
    end
  end

  context 'cannot #pop out of order' do
    before(:each) do
      subject.push(:ignore, :exists)
      subject.push(:ignore, :exists, :changes)
    end

    it do
      expect { subject.pop(:ignore, :exists) }.to \
        raise_error(RuntimeError,
                    /Unexpected pop.*:ignore.*:exists.*:none.*:none.*:none/)
    end
  end

  context 'cannot miss a test' do
    before(:each) do
      subject.push(:ignore, :exists)
      subject.push(:ignore, :exists, :changes)
      subject.push(:ignore, :exists, :changes, :ignore, :pass)
      subject.pop(:ignore, :exists, :changes, :ignore, :pass)
      # we will miss push+pop for :ignore, :exists, :changes, :ignore, :fail
      subject.pop(:ignore, :exists, :changes)
      subject.push(:ignore, :exists, :no_change)
      subject.push(:ignore, :exists, :no_change, :ignore, :fail)
      subject.pop(:ignore, :exists, :no_change, :ignore, :fail)
      subject.pop(:ignore, :exists, :no_change)
      subject.pop(:ignore, :exists)
    end

    it do
      expect { subject.verify }.to \
        raise_error(RuntimeError, /missing.*\[:exists, :changes, :fail\]/)
    end
  end

  context 'cannot miss a pop()' do
    before(:each) do
      subject.push(:ignore, :exists)
      subject.push(:ignore, :exists, :changes)
      subject.push(:ignore, :exists, :changes, :ignore, :pass)
      subject.pop(:ignore, :exists, :changes, :ignore, :pass)
      subject.push(:ignore, :exists, :changes, :ignore, :fail)
      subject.pop(:ignore, :exists, :changes, :ignore, :fail)
      # we will miss pop for :ignore, :exists, :changes
      subject.push(:ignore, :exists, :no_change)
      subject.push(:ignore, :exists, :no_change, :ignore, :fail)
      subject.pop(:ignore, :exists, :no_change, :ignore, :fail)
      subject.pop(:ignore, :exists, :no_change)
    end

    it do
      expect(-> { subject.pop(:ignore, :exists) }).to \
        raise_error(RuntimeError,
                    /Expect.*pop for.*:ignore.*:exists.*:changes.*:none.*:none/)
    end
  end

  context 'matrixes are added to the collection automatically' do
    let(:registry) { Provider::TestMatrix::Registry.instance }

    before(:each) do
      registry.stubs(:add).once # ensures add() is called
    end

    after(:each) { registry.unstub(:add) }

    it { is_expected.not_to be_nil }
  end
end

describe Provider::TestMatrix::Collector do
  context 'matrixes are verified' do
    let(:matrixes) { [mock('matrix-1'), mock('matrix-2')] }

    before(:each) do
      matrixes.each do |m|
        m.stubs(:verify).once # ensures verify() is called
        m.stubs(:name).returns(m)
        subject.add(m, __FILE__, m)
      end
    end

    it { expect { subject.verify_all }.not_to raise_error }
  end
end
