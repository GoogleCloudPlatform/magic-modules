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

describe Api::Product do
  context 'requires name' do
    subject { -> { product('display_name: "Bar"').validate } }
    it { is_expected.to raise_error(StandardError, /Missing 'name'/) }
  end

  context 'requires versions' do
    subject do
      lambda do
        product('name: "foo"',
                'scopes:',
                '  - link/to/scope',
                'name: "Bar"',
                'objects:',
                '  - !ruby/object:Api::Resource',
                '    kind: foo#resource',
                '    base_url: myres/',
                '    description: foo',
                '    name: "res1"',
                '    properties:',
                '      - !ruby/object:Api::Type',
                '        name: var',
                '        description: desc').validate
      end
    end

    it { is_expected.to raise_error(StandardError, /Missing 'versions'/) }
  end

  context 'requires at most one default version' do
    subject do
      lambda do
        product('name: "foo"',
                'scopes:',
                '  - link/to/scope',
                'name: "Bar"',
                'versions:',
                '  - !ruby/object:Api::Product::Version',
                '    name: ga',
                '    base_url: "http://foo/var/v1"',
                '    default: true',
                '  - !ruby/object:Api::Product::Version',
                '    name: beta',
                '    base_url: "http://foo/var/beta"',
                '    default: true',
                'objects:',
                '  - !ruby/object:Api::Resource',
                '    kind: foo#resource',
                '    base_url: myres/',
                '    description: foo',
                '    name: "res1"',
                '    properties:',
                '      - !ruby/object:Api::Type',
                '        name: var',
                '        description: desc').validate
      end
    end

    it do
      is_expected.to raise_error(StandardError,
                                 /must specify at most one default/)
    end
  end

  context 'requires at least one default version' do
    subject do
      lambda do
        product('name: "foo"',
                'scopes:',
                '  - link/to/scope',
                'name: "Bar"',
                'versions:',
                '  - !ruby/object:Api::Product::Version',
                '    name: ga',
                '    base_url: "http://foo/var/v1"',
                '  - !ruby/object:Api::Product::Version',
                '    name: beta',
                '    base_url: "http://foo/var/beta"',
                'objects:',
                '  - !ruby/object:Api::Resource',
                '    kind: foo#resource',
                '    base_url: myres/',
                '    description: foo',
                '    name: "res1"',
                '    properties:',
                '      - !ruby/object:Api::Type',
                '        name: var',
                '        description: desc').validate
      end
    end

    it do
      is_expected.to raise_error(StandardError,
                                 /must specify a default/)
    end
  end

  context 'requires objects' do
    subject do
      lambda do
        product('name: "foo"',
                'name: "Bar"',
                'versions:',
                '  - !ruby/object:Api::Product::Version',
                '    name: ga',
                '    base_url: "baz"').validate
      end
    end
    it { is_expected.to raise_error(StandardError, /Missing 'objects'/) }
  end

  context 'only allows Resources as objects' do
    subject do
      lambda do
        product('name: "foo"',
                'name: "Bar"',
                'versions:',
                '  - !ruby/object:Api::Product::Version',
                '    name: ga',
                '    base_url: "baz"',
                'objects:',
                '  - bah. bad object!').validate
      end
    end

    it do
      is_expected
        .to raise_error(StandardError,
                        /Property.*objects.*instead.*Api::Resource/)
    end
  end

  private

  def product(*data)
    Google::YamlValidator.parse(['--- !ruby/object:Api::Product'].concat(data)
                                                                 .join("\n"))
  end
end
