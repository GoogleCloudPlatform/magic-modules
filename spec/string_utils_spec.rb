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

describe Google::StringUtils do
  context '#underscore' do
    subject { described_class.underscore('aStringInCamelCase') }
    it { is_expected.to eq 'a_string_in_camel_case' }
  end

  context '#first_sentence' do
    context 'sentence end with period' do
      subject do
        described_class.first_sentence('Lorem ipsum. Dolor sit amet. Elit')
      end
      it { is_expected.to eq 'Lorem ipsum.' }
    end

    context 'sentence end with question mark' do
      subject do
        described_class.first_sentence('Lorem ipsum? Dolor sit amet. Elit')
      end
      it { is_expected.to eq 'Lorem ipsum?' }
    end

    context 'sentence end with exclamation mark' do
      subject do
        described_class.first_sentence('Lorem ipsum! Dolor sit amet. Elit')
      end
      it { is_expected.to eq 'Lorem ipsum!' }
    end

    context 'no period returns full string' do
      subject { described_class.first_sentence('Lorem ipsum dolor') }
      it { is_expected.to eq 'Lorem ipsum dolor' }
    end
  end
end
