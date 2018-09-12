# Copyright 2018 Google Inc.
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

FactoryBot.define do
  # Using YAML parsing short circuits the object lifecycle.
  # We don't have proper initialize methods, because our
  # objects just appear post-YAML parsing with
  # all of the correct values.
  #
  # FactoryBot should create its test objects in the 
  # same manner.
  initialize_with do
    obj = new
    attributes.each { |k, v| obj.instance_variable_set("@#{k.id2name}", v) }
    # TODO(alexstephen): Build out a better default factory that can successfully
    # validate.
    #obj.validate
    obj
  end
  to_create {}
  
  factory :product, class: Api::Product do
    name { "Google TestSpec Engine" }
    prefix { "gspec" }
    objects { [] }
    scopes { [] }
  end

  factory :resource, class: Api::Resource do
    association :__product, factory: :product
  end

  factory :string, class: Api::Type::String do
    name { "string-property" }
    description { "Description for a test property" }
    output { false }
    required { false }
    input { false }

    association :__resource, factory: :resource
  end
end
