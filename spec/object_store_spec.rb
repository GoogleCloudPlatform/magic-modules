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

describe Google::ObjectStore do
  context 'Singleton' do
    subject { -> { described_class.new } }
    it { is_expected.to raise_error NoMethodError }
    it { expect(described_class.instance).to be_a Singleton }
  end

  context 'Object persisted' do
    before(:all) do
      described_class.instance.add SomeTypeForKey, SomeTypeForValue.new(1)
      described_class.instance.add SomeTypeForKey, SomeTypeForValue.new(2)
    end

    subject { described_class.instance[SomeTypeForKey.class] }

    it { is_expected.to all(be_a(SomeTypeForValue)) }
    it { described_class.instance[SomeTypeForKey].size == 2 }
  end

  context 'Mixed objects persisted' do
    before(:all) do
      described_class.instance.add SomeTypeForKey, SomeTypeForValue.new(1)
      described_class.instance.add SomeOtherTypeForKey, SomeTypeForValue.new(2)
    end

    context 'SomeTypeForKey' do
      subject { described_class.instance[SomeTypeForKey.class] }
      it { is_expected.to all(be_a(SomeTypeForValue)) }
      it { !described_class.instance[SomeTypeForKey].empty? }
    end

    context 'SomeOtherTypeForKey' do
      subject { described_class.instance[SomeOtherTypeForKey.class] }
      it { is_expected.to all(be_a(SomeTypeForValue)) }
      it { !described_class.instance[SomeOtherTypeForKey].empty? }
    end
  end

  private

  module Puppet
    def self.debug(msg)
      puts "Puppet(debug): #{msg}"
    end
  end

  class SomeTypeForKey
  end

  class SomeOtherTypeForKey
  end

  class SomeTypeForValue
    attr_reader :seed

    def initialize(seed)
      @seed = seed
    end
  end
end
