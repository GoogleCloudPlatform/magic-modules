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

class TestObject < Api::Object::Named
  attr_reader :some_property
end

describe Api::Object do
  context 'requires name' do
    subject { -> { object('some_property: "bar"').validate } }
    it { is_expected.to raise_error(StandardError, /Missing 'name'/) }
  end

  context 'out_name underscore style' do
    subject { object('name: "MyCamelCaseObjectName"').out_name }
    it { is_expected.to eq 'my_camel_case_object_name' }
  end

  private

  def object(*data)
    Google::YamlValidator.parse(['--- !ruby/object:TestObject'].concat(data)
                                                               .join("\n"))
  end
end
