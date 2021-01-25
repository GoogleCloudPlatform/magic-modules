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

module Level1
  module Level2
    class MyType < Api::Type
      attr_reader :myproperty
    end
  end
end

describe Api::Type do
  context 'requires name' do
    subject { -> { Api::Type.new.validate } }
    it { is_expected.to raise_error(StandardError, /Missing 'name'/) }
  end

  context 'requires description' do
    subject do
      lambda do
        Google::YamlValidator.parse('--- !ruby/object:Api::Type
                   name: "test"').validate
      end
    end
    it { is_expected.to raise_error(StandardError, /Missing 'description'/) }
  end

  context 'type name' do
    subject do
      Google::YamlValidator.parse('--- !ruby/object:Level1::Level2::MyType
                 name: "test"
                 description: "mydescription"
                 myproperty: 10')
    end
    it { is_expected.to have_attributes(type: 'MyType') }
  end

  context 'allows specific properties on children' do
    subject do
      lambda do
        Google::YamlValidator.parse('--- !ruby/object:Level1::Level2::MyType
                   name: "test"
                   description: "mydescription"
                   myproperty: 10').validate
      end
    end
    it { is_expected.not_to raise_error }
  end

  context 'prevents extraneous values' do
    subject do
      lambda do
        Google::YamlValidator.parse('--- !ruby/object:Level1::Level2::MyType
                   name: "test"
                   description: "mydescription"
                   myproperty: 10
                   a_bad_property: 20').validate
      end
    end

    it do
      is_expected.to raise_error(StandardError, /Extraneous .*a_bad_property/)
    end
  end

  context 'retrieves simple type name' do
    subject do
      Google::YamlValidator.parse('--- !ruby/object:Level1::Level2::MyType
                 name: "MyTypeName"
                 description: "mydescription"
                 myproperty: 10')
    end

    it { is_expected.to be_instance_of Level1::Level2::MyType }
  end
end

describe Api::Type::Array do
  context 'requires underlying type' do
    subject do
      lambda do
        Google::YamlValidator.parse('--- !ruby/object:Api::Type::Array
                   description: some description
                   name: "test"').validate
      end
    end

    it { is_expected.to raise_error(StandardError, /Missing 'item_type'/) }
  end

  context 'requires underlying type to exist' do
    subject do
      lambda do
        Google::YamlValidator.parse('--- !ruby/object:Api::Type::Array
                   name: myname
                   description: some description
                   item_type: Level1::Level2::MyType').validate
      end
    end

    it { is_expected.not_to raise_error }
  end

  context 'requires underlying type to exist' do
    subject do
      lambda do
        Google::YamlValidator.parse('--- !ruby/object:Api::Type::Array
                   name: myname
                   description: some description
                   item_type: ATypeThatDoesNotExist').validate
      end
    end
    it do
      is_expected.to \
        raise_error(StandardError,
                    /uninitialized constant .*ATypeThatDoesNotExist/)
    end
  end

  context 'requires name' do
    subject do
      lambda do
        Google::YamlValidator.parse('--- !ruby/object:Api::Type::Array
                   item_type: Api::String').validate
      end
    end
    it { is_expected.to raise_error(StandardError, /Missing 'name'/) }
  end

  context 'requires description' do
    subject do
      lambda do
        Google::YamlValidator.parse('--- !ruby/object:Api::Type::Array
                   item_type: Api::String
                   name: "test"').validate
      end
    end
    it { is_expected.to raise_error(StandardError, /Missing 'description'/) }
  end
end

describe Api::Type::Enum do
  context 'requires values' do
    subject do
      lambda do
        Google::YamlValidator.parse('--- !ruby/object:Api::Type::Enum
                   description: some description
                   name: "test"').validate
      end
    end
    it { is_expected.to raise_error(StandardError, /Missing 'values'/) }
  end

  context 'values is an array' do
    subject do
      Google::YamlValidator.parse('--- !ruby/object:Api::Type::Enum
                 name: "name"
                 description: "description"
                 values:
                   - :A
                   - "b"
                   - 3')
    end
    it { is_expected.to have_attributes(values: [:A, 'b', 3]) }
  end

  context 'requires name' do
    subject do
      lambda do
        Google::YamlValidator.parse('--- !ruby/object:Api::Type::Enum
                   values:
                     - :A').validate
      end
    end

    it { is_expected.to raise_error(StandardError, /Missing 'name'/) }
  end

  context 'requires description' do
    subject do
      lambda do
        Google::YamlValidator.parse('--- !ruby/object:Api::Type::Enum
                   name: "test"
                   values:
                     - :A').validate
      end
    end

    it { is_expected.to raise_error(StandardError, /Missing 'description'/) }
  end
end

describe Api::Type::ResourceRef do
  context 'requires valid resource' do
    let(:spec_location) do
      File.join(File.dirname(__FILE__), 'data',
                'resourceref-missingresource.yaml')
    end
    let(:spec) do
      File.open(spec_location, 'r')
    end

    after(:each) { spec.close }
    subject do
      lambda do
        Google::YamlValidator.parse(spec.read).validate
      end
    end

    it { is_expected.to raise_error(StandardError, /Missing 'resource'/) }
  end

  context 'requires valid imports' do
    let(:spec_location) do
      File.join(File.dirname(__FILE__), 'data',
                'resourceref-missingimports.yaml')
    end
    let(:spec) do
      File.open(spec_location, 'r')
    end
    let(:error) do
      /'missing' does not exist on 'ReferencedResource'/
    end

    after(:each) { spec.close }
    subject do
      lambda do
        Google::YamlValidator.parse(spec.read).validate
      end
    end

    it { is_expected.to raise_error(StandardError, error) }
  end
end
