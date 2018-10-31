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

class MyTestObject < Google::YamlValidator
  attr_reader :some_property
end

describe Google::YamlValidator do
  context 'prevents extraneous variables' do
    subject do
      lambda do
        object('some_property: "good"',
               'other_variable: 12345').validate
      end
    end

    it do
      is_expected.to raise_error(StandardError,
                                 /Extraneous variable 'other_variable'/)
    end
  end

  context 'sets variable properties' do
    subject { object('name: "bar"', 'description: "good"') }

    it 'pre-flight test' do
      is_expected.not_to respond_to(:@__custom_property)
    end

    context 'set variable' do
      let(:custom_obj) { mock('custom') }

      before do
        subject.set_variable(custom_obj, :__custom_property)
      end

      it do
        expect(subject.instance_variable_get(:@__custom_property))
          .to be custom_obj
      end
    end
  end

  context 'do not allow unapproved classes deserialized' do
    subject do
      -> { described_class.parse("--- !ruby/object:Digest::SHA256\na: b") }
    end

    it do
      is_expected.to raise_error(Psych::DisallowedClass, /Digest::SHA256/)
    end
  end

  private

  def object(*data)
    described_class.parse(['--- !ruby/object:MyTestObject'].concat(data)
                                                           .join("\n"))
  end
end
