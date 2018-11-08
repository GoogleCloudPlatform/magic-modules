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
require 'api/compiler'

describe Api::Compiler do
  context 'should fail if file does not exist' do
    subject { -> { Api::Compiler.new('spec/data/somedummyfile').run } }
    it { is_expected.to raise_error(Errno::ENOENT) }
  end

  context 'should use the file provided' do
    let(:reader) { mock('reader') }

    subject { -> { Api::Compiler.new('my-file-to-parse.yaml').run } }

    before do
      IO.expects(:read).with('my-file-to-parse.yaml')
        .returns('--- !ruby/object:Api::Product
                      name: "foo"')
        .twice
    end

    it { is_expected.not_to raise_error }
  end

  context 'parses file' do
    subject { Api::Compiler.new('spec/data/good-file.yaml').run }

    before do
      subject.validate
    end

    it { is_expected.to be_instance_of Api::Product }
    it { is_expected.to have_attributes(name: 'My Product') }
    it { is_expected.to have_attribute_of_length(objects: 4) }
  end

  context 'should only accept product' do
    let(:reader) { mock('reader') }

    subject do
      -> { Api::Compiler.new('my-file-to-parse.yaml').run.validate }
    end

    before do
      IO.expects(:read).with('my-file-to-parse.yaml')
        .returns('something: "else"')
    end

    it do
      is_expected.to raise_error(StandardError, /is .* instead of Api::Product/)
    end
  end
end
