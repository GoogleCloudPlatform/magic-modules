# Copyright 2023 Google Inc.
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
require 'openapi_generate/parser'
require 'active_support/inflector'

describe OpenAPIGenerate::Parser do
  context 'should run on sample file for Post' do
    subject do
      OpenAPIGenerate::Parser.new('fake', 'fake').build_resource(
        'spec/data/test-openapi-spec.json',
        '/v1alpha5/projects/{projectsId}/locations/{locationsId}/posts',
        'Post'
      )
    end
    it { is_expected.to have_attributes(name: 'Post') }
    it { is_expected.to have_attribute_of_length(properties: 12) }
  end

  context 'should run on sample file for Comment' do
    subject do
      OpenAPIGenerate::Parser.new('fake', 'fake').build_resource(
        'spec/data/test-openapi-spec.json',
        '/v1alpha5/projects/{projectsId}/locations/{locationsId}/posts/{postsId}/comments',
        'Comment'
      )
    end
    it { is_expected.to have_attributes(name: 'Comment') }
    it { is_expected.to have_attribute_of_length(properties: 9) }
  end
end
