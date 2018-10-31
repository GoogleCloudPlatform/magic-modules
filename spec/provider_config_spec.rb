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

module Module1
  module Module2
    class MyConfig < Provider::Config
      attr_reader :p1
      attr_reader :p2

      def validate
        # Not validating or invoking full validation for these tests
      end
    end
  end
end

describe Provider::Config do
  it 'requires override provider' do
    expect { Provider::Config.new.provider }
      .to raise_error(StandardError, /provider not implemented/)
  end

  context 'parsing validation' do
    it 'fails if not a Provider::Config class' do
      IO.expects(:read).with('foo/bar').returns([
        'a: A',
        'b: B'
      ].join("\n"))

      expect do
        Provider::Config.parse('foo/bar')
      end.to raise_error(StandardError, /is not a Provider::Config/)
    end
  end
end
